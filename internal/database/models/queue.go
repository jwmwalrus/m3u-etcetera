package models

import (
	"github.com/jwmwalrus/bnp/pointers"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/m3u-etcetera/pkg/poser"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

// GetPerspectiveQueue returns the queue associated to the perspective index
func (idx PerspectiveIndex) GetPerspectiveQueue() (q *Queue, err error) {
	q = &Queue{}
	err = db.Preload("Perspective").
		Joins(
			"JOIN perspective ON queue.perspective_id = perspective.id AND perspective.idx = ?",
			int(idx),
		).
		First(q).
		Error
	return

}

// Queue defines a queue
type Queue struct { // too transient
	ID            int64       `json:"id" gorm:"primaryKey"`
	CreatedAt     int64       `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt     int64       `json:"updatedAt" gorm:"autoUpdateTime:nano"`
	PerspectiveID int64       `json:"perspectiveId" gorm:"uniqueIndex:unique_idx_queue_perspective_id,not null"`
	Perspective   Perspective `json:"perspective" gorm:"foreignKey:PerspectiveID"`
}

// Read implements the DataReader interface
func (q *Queue) Read(id int64) (err error) {
	err = db.Joins("Perspective").
		// Joins("JOIN perspective ON queue.perspective_id = perspective.id").
		First(q, id).Error
	return
}

// Add adds the given locations/IDs to the end of the queue
func (q *Queue) Add(locations []string, ids []int64) {
	entry := log.WithFields(log.Fields{
		"locations": locations,
		"ids":       ids,
	})
	entry.Info("Adding payload to queue")

	for _, v := range locations {
		qt := QueueTrack{Location: v}

		if err := q.appendTo(&qt); err != nil {
			entry.Error(err)
			continue
		}
	}
	for _, v := range ids {
		var err error
		t := Track{}

		if err = t.Read(v); err != nil {
			entry.Error(err)
			continue
		}

		qt := QueueTrack{Location: t.Location, TrackID: v}
		onerror.WithEntry(entry).Log(q.appendTo(&qt))
	}

	subscription.Broadcast(subscription.ToQueueStoreEvent)
}

// Clear removes all entries from the queue
func (q *Queue) Clear() {
	entry := log.WithField("q", *q)
	entry.Info("Clearing queue")

	s := []QueueTrack{}
	err := db.Where("queue_id = ? AND played = 0", q.ID).Find(&s).Error
	if err != nil {
		return
	}

	for i := 0; i < len(s); i++ {
		s[i].Played = true
	}

	if err = db.Save(&s).Error; err != nil {
		entry.Error(err)
		return
	}

	subscription.Broadcast(subscription.ToQueueStoreEvent)
}

// DeleteAt deletes the given position from the queue
func (q *Queue) DeleteAt(position int) {
	entry := log.WithFields(log.Fields{
		"q":        *q,
		"position": position,
	})
	entry.Info("Deleting entry from queue")

	s := []QueueTrack{}
	err := db.Where("queue_id = ? AND played = 0", q.ID).Order("position ASC").
		Find(&s).
		Error
	if err != nil {
		return
	}

	list, qt := poser.DeleteAt(pointers.FromSlice(s), position)
	s = pointers.ToValues(list)

	if qt != nil && qt.ID > 0 {
		if err := qt.Save(); err != nil {
			entry.Error(err)
			return
		}
	}

	if err := db.Save(&s).Error; err != nil {
		entry.Error(err)
		return
	}

	subscription.Broadcast(subscription.ToQueueStoreEvent)
}

// InsertAt inserts the given locations and IDs into the queue
func (q *Queue) InsertAt(position int, locations []string, ids []int64) {
	entry := log.WithFields(log.Fields{
		"q":         q,
		"position":  position,
		"locations": locations,
		"ids":       ids,
	})
	entry.Info("Inserting entry into queue")

	for i := len(ids) - 1; i >= 0; i-- {
		t := Track{}
		if err := t.Read(ids[i]); err != nil {
			entry.Error(err)
			continue
		}

		qt := QueueTrack{Location: t.Location, TrackID: ids[i], Position: position}
		onerror.WithEntry(entry).Log(q.insertInto(&qt))
	}
	for i := len(locations) - 1; i >= 0; i-- {
		qt := QueueTrack{Location: locations[i], Position: position}
		onerror.WithEntry(entry).Log(q.insertInto(&qt))
	}

	subscription.Broadcast(subscription.ToQueueStoreEvent)
}

// IsEmpty returns true if there are no tracks in the queue
func (q *Queue) IsEmpty() bool {
	s := []QueueTrack{}
	db.Where("queue_id = ? AND played = 0", q.ID).Order("position").Find(&s)
	return len(s) == 0
}

// MoveTo moves one queue track from one position to another
func (q *Queue) MoveTo(to, from int) {
	if from == to || from < 1 {
		return
	}

	entry := log.WithFields(log.Fields{
		"from": from,
		"to":   to,
	})
	entry.Info("Moving queue tracks")

	s := []QueueTrack{}
	err := db.Where("queue_id = ? AND played = 0", q.ID).Order("position").
		Find(&s).
		Error
	if err != nil {
		entry.Error(err)
		return
	}
	if len(s) == 0 || from > len(s) {
		return
	}

	list := poser.MoveTo(pointers.FromSlice(s), to, from)
	s = pointers.ToValues(list)

	if err := db.Save(&s).Error; err != nil {
		entry.Error(err)
		return
	}
	subscription.Broadcast(subscription.ToQueueStoreEvent)
}

// Pop returns the next entry to be played from the queue
func (q *Queue) Pop() (qt *QueueTrack) {
	entry := log.WithField("q", *q)
	entry.Debug("Popping from queue")

	s := []QueueTrack{}
	err := db.Where("queue_id = ? AND played = 0", q.ID).Order("position ASC").
		Find(&s).
		Error
	if err != nil {
		entry.Error(err)
		return
	}

	list, qt := poser.Pop(pointers.FromSlice(s))
	s = pointers.ToValues(list)

	if qt == nil {
		return
	}

	entry.Info("Found location to pop from queue:", qt.Location)
	if len(s) > 0 {
		onerror.WithEntry(entry).Log(db.Save(&s).Error)
	}
	onerror.WithEntry(entry).Log(qt.Save())

	subscription.Broadcast(subscription.ToQueueStoreEvent)
	return
}

func (q *Queue) appendTo(qt *QueueTrack) (err error) {
	log.WithFields(log.Fields{
		"q":  *q,
		"qt": *qt,
	}).
		Info("Appending track to queue")

	qt.QueueID = q.ID
	qt.Played = true
	if err = qt.Create(); err != nil {
		return
	}
	qt.Played = false

	s := []QueueTrack{}
	err = db.Where("queue_id = ? AND played = 0", q.ID).Order("position ASC").
		Find(&s).
		Error
	if err != nil {
		return
	}

	list := poser.AppendTo(pointers.FromSlice(s), qt)
	s = pointers.ToValues(list)
	if err = db.Save(&s).Error; err != nil {
		return
	}

	go func() {
		queueTrackNeeded <- qt.ID
	}()
	return
}

func (q *Queue) insertInto(qt *QueueTrack) (err error) {
	log.WithFields(log.Fields{
		"q":  *q,
		"qt": *qt,
	}).
		Info("Inserting track into queue")

	qt.QueueID = q.ID
	qt.Played = true
	if err = qt.Create(); err != nil {
		return
	}
	qt.Played = false

	s := []QueueTrack{}
	err = db.Where("queue_id = ? AND played = 0", q.ID).
		Order("position ASC").
		Find(&s).
		Error
	if err != nil {
		return
	}
	list := poser.InsertInto(pointers.FromSlice(s), qt.Position, qt)
	s = pointers.ToValues(list)
	if err = db.Save(&s).Error; err != nil {
		return
	}

	go func() {
		queueTrackNeeded <- qt.ID
	}()

	return
}
