package models

import (
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
	log.WithFields(log.Fields{
		"locations": locations,
		"ids":       ids,
	}).
		Info("Adding payload to queue")

	for _, v := range locations {
		qt := QueueTrack{Location: v}

		if err := q.appendTo(&qt); err != nil {
			log.Error(err)
			continue
		}
	}
	for _, v := range ids {
		var err error
		t := Track{}

		if err = t.Read(v); err != nil {
			log.Error(err)
			continue
		}

		qt := QueueTrack{Location: t.Location, TrackID: v}
		onerror.Log(q.appendTo(&qt))
	}

	subscription.Broadcast(subscription.ToQueueStoreEvent)
	return
}

// Clear removes all entries from the queue
func (q *Queue) Clear() {
	log.WithField("q", *q).
		Info("Clearing queue")

	s := []QueueTrack{}
	err := db.Where("queue_id = ? AND played = 0", q.ID).Find(&s).Error
	if err != nil {
		return
	}

	for i := 0; i < len(s); i++ {
		s[i].Played = true
	}

	if err = db.Save(&s).Error; err != nil {
		log.Error(err)
		return
	}

	subscription.Broadcast(subscription.ToQueueStoreEvent)
	return
}

// DeleteAt deletes the given position from the queue
func (q *Queue) DeleteAt(position int) {
	log.WithFields(log.Fields{
		"q":        *q,
		"position": position,
	}).
		Info("Deleting entry from queue")

	s := []QueueTrack{}
	err := db.Where("queue_id = ? AND played = 0", q.ID).Find(&s).Error
	if err != nil {
		return
	}

	s = poser.DeleteAt(s, position)

	if err := db.Save(&s).Error; err != nil {
		log.Error(err)
		return
	}

	subscription.Broadcast(subscription.ToQueueStoreEvent)
}

// InsertAt inserts the given locations and IDs into the queue
func (q *Queue) InsertAt(position int, locations []string, ids []int64) {
	log.WithFields(log.Fields{
		"q":         q,
		"position":  position,
		"locations": locations,
		"ids":       ids,
	}).
		Info("Inserting entry into queue")

	for i := len(ids) - 1; i >= 0; i-- {
		t := Track{}
		if err := t.Read(ids[i]); err != nil {
			log.Error(err)
			continue
		}

		qt := QueueTrack{Location: t.Location, TrackID: ids[i], Position: position}
		onerror.Log(q.insertInto(&qt))
	}
	for i := len(locations) - 1; i >= 0; i-- {
		qt := QueueTrack{Location: locations[i], Position: position}
		onerror.Log(q.insertInto(&qt))
	}

	subscription.Broadcast(subscription.ToQueueStoreEvent)
}

// MoveTo moves one queue track from one position to another
func (q *Queue) MoveTo(to, from int) {
	if from == to || from < 1 {
		return
	}

	log.WithFields(log.Fields{
		"from": from,
		"to":   to,
	}).
		Info("Moving queue tracks")

	s := []QueueTrack{}
	db.Where("queue_id = ? AND played = 0", q.ID).Order("position").Find(&s)
	if len(s) == 0 || from > len(s) {
		return
	}

	moved := poser.MoveTo(s, to, from)

	if err := db.Save(&moved).Error; err != nil {
		return
	}
	subscription.Broadcast(subscription.ToQueueStoreEvent)
}

// Pop returns the next entry to be played from the queue
func (q *Queue) Pop() (qt *QueueTrack) {
	log.WithField("q", *q).
		Debug("Popping from queue")

	s := []QueueTrack{}
	err := db.Where("queue_id = ? AND played = 0", q.ID).Find(&s).Error
	if err != nil {
		return
	}

	s, x := poser.Pop(s)

	if x.ID == 0 {
		return
	}
	qt = &x

	log.Info("Found location to pop from queue:", qt.Location)
	onerror.Log(db.Save(&s).Error)

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
	err = db.Where("queue_id = ? AND played = 0", q.ID).Find(&s).Error
	if err != nil {
		return
	}

	s = poser.AppendTo(s, *qt)
	if err = db.Save(&s).Error; err != nil {
		return
	}

	go findQueueTrack(qt)
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
	s = poser.InsertInto(s, qt.Position, *qt)
	if err = db.Save(&s).Error; err != nil {
		return
	}

	go findQueueTrack(qt)
	return
}
