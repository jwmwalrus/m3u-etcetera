package models

import (
	"encoding/json"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	log "github.com/sirupsen/logrus"
)

// GetPerspectiveQueue returns the queue associated to the perspective index
func (idx PerspectiveIndex) GetPerspectiveQueue() (q *Queue, err error) {
	q = &Queue{}
	err = db.Joins("JOIN perspective ON queue.perspective_id = perspective.id AND perspective.idx = ?", int(idx)).First(q).Error
	return

}

// Queue defines a queue
type Queue struct { // too transient
	ID            int64       `json:"id" gorm:"primaryKey"`
	CreatedAt     int64       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     int64       `json:"updatedAt" gorm:"autoUpdateTime"`
	PerspectiveID int64       `json:"perspectiveId" gorm:"uniqueIndex:unique_idx_queue_perspective_id,not null"`
	Perspective   Perspective `json:"perspective" gorm:"foreignKey:PerspectiveID"`
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
		t := Track{}
		if err := db.First(&t, v).Error; err != nil {
			log.Error(err)
			continue
		}

		qt := QueueTrack{Location: t.Location, TrackID: v}
		if err := q.appendTo(&qt); err != nil {
			log.Error(err)
			continue
		}
	}
	return
}

// Clear removes all entries from the queue
func (q *Queue) Clear() {
	log.WithField("q", *q).
		Info("Clearing queue")

	s := []QueueTrack{}
	if err := db.Where("queue_id = ? AND played = 0", q.ID).Find(&s).Error; err != nil {
		return
	}

	for i := 0; i < len(s); i++ {
		s[i].Played = true
	}

	err := db.Where("id > 0").Save(&s).Error
	onerror.Log(err)
	return
}

// DeleteAt deletes the given position from the queue
func (q *Queue) DeleteAt(position int) {
	log.WithFields(log.Fields{
		"q":        *q,
		"position": position,
	}).
		Info("Deleting entry from queue")

	qt := QueueTrack{}
	if err := db.Where("queue_id = ? AND position = ?", q.ID, position).First(&qt).Error; err != nil {
		return
	}

	err := q.deleteFrom(&qt)
	onerror.Log(err)
	return
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
		if err := db.First(&t, ids[i]).Error; err != nil {
			log.Error(err)
			continue
		}

		qt := QueueTrack{Location: t.Location, TrackID: ids[i], Position: position}
		err := q.insertInto(&qt)
		onerror.Log(err)
	}
	for i := len(locations) - 1; i >= 0; i-- {
		qt := QueueTrack{Location: locations[i], Position: position}
		err := q.insertInto(&qt)
		onerror.Log(err)
	}
}

// Pop returns the next entry to be played from the queue
func (q *Queue) Pop() (qt QueueTrack, err error) {
	log.WithField("q", *q).
		Info("Popping from queue")

	q.reorder()

	if err = db.Where("queue_id = ? AND played = 0 AND position = 1", q.ID).First(&qt).Error; err != nil {
		log.Info("Nothing to pop!")
		return
	}
	log.Info("Found location to pop from queue:", qt.Location)
	qt.Played = true
	err = db.Save(&qt).Error

	go q.reorder()
	return
}

func (q *Queue) appendTo(qt *QueueTrack) (err error) {
	log.WithFields(log.Fields{
		"q":  *q,
		"qt": *qt,
	}).
		Info("Appending track to queue")

	var qs int64
	db.Model(&QueueTrack{}).Where("queue_id = ? AND played = 0", q.ID).Count(&qs)

	qt.QueueID = q.ID
	qt.Position = int(qs) + 1
	err = db.Create(qt).Error
	go findQueueTrack(qt)
	return
}

func (q *Queue) deleteFrom(qt *QueueTrack) (err error) {
	log.WithFields(log.Fields{
		"q":  *q,
		"qt": *qt,
	}).
		Info("Deleting track from queue")

	qt.Played = true
	if err = db.Save(qt).Error; err != nil {
		return
	}
	go q.reorder()
	return
}

func (q *Queue) insertInto(qt *QueueTrack) (err error) {
	log.WithFields(log.Fields{
		"q":  *q,
		"qt": *qt,
	}).
		Info("Inserting track into queue")

	var qs int64
	if err = db.Model(&QueueTrack{}).Where("queue_id = ? AND played = 0", q.ID).Count(&qs).Error; err != nil {
		return
	}
	if qt.Position <= 1 {
		qt.Position = 0
	} else if qt.Position > 1 && qt.Position <= int(qs) {
		s := []QueueTrack{}
		if err = db.Where("queue_id = ? AND played = 0 AND position >= ?", q.ID, qt.Position).Find(&s).Error; err != nil {
			return
		}
		for i := 0; i < len(s); i++ {
			s[i].Position++
		}
		if err = db.Where("id > 0").Save(&s).Error; err != nil {
			return
		}
	} else {
		err = q.appendTo(qt)
		return
	}
	qt.QueueID = q.ID
	if err = db.Create(qt).Error; err != nil {
		return
	}
	q.reorder()
	go findQueueTrack(qt)
	return
}

func (q *Queue) reorder() {
	log.WithField("q", *q).
		Info("Reordering queue")

	s := []QueueTrack{}
	db.Where("queue_id = ? AND played = 0", q.ID).Order("position").Find(&s)
	if len(s) == 0 {
		return
	}

	for i := 0; i < len(s); i++ {
		s[i].Position = i + 1
	}

	err := db.Where("id > 0").Save(&s).Error
	onerror.Log(err)
	return
}

// QueueTrack defines a track in the queue
type QueueTrack struct { // too transient
	ID        int64  `json:"id" gorm:"primaryKey"`
	Position  int    `json:"position"`
	Played    bool   `json:"played"`
	Location  string `json:"location" gorm:"not null"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime"`
	TrackID   int64  `json:"trackId" gorm:"index:idx_queue_track_track_id"`
	QueueID   int64  `json:"queueId" gorm:"index:idx_queue_track_queue_id,not null"`
	Queue     Queue  `json:"queue" gorm:"foreignKey:QueueID"`
}

// GetAllQueueTracks returns all queue tracks for the given perspective, constrained by a limit
func GetAllQueueTracks(idx PerspectiveIndex, limit int) (s []QueueTrack) {
	log.WithFields(log.Fields{
		"idx":   idx,
		"limit": limit,
	}).
		Info("Getting all queue tracks")

	q, err := idx.GetPerspectiveQueue()
	if err != nil {
		log.Error(err)
		return
	}

	tx := db.
		Where("played = 0 AND queue_id = ?", q.ID).
		Order("position ASC")

	if limit > 9 {
		tx.Limit(limit)
	}

	tx.Find(&s)
	return
}

// ToProtobuf converter
func (qt *QueueTrack) ToProtobuf() *m3uetcpb.QueueTrack {
	bv, err := json.Marshal(qt)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.QueueTrack{}
	}

	out := &m3uetcpb.QueueTrack{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)
	return out
}

// findQueueTrack attempts to find track from location
func findQueueTrack(qt *QueueTrack) {
	if qt.TrackID > 0 {
		return
	}

	log.WithField("qt", *qt).
		Info("Finding track for queue entry")

	t := Track{}
	if err := db.Where("location = ?", qt.Location).First(&t).Error; err != nil {
		return
	}
	qt.TrackID = t.ID
	err := db.Save(qt).Error
	onerror.Log(err)
	return
}
