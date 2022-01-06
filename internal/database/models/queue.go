package models

import (
	"encoding/json"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// GetPerspectiveQueue returns the queue associated to the perspective index
func (idx PerspectiveIndex) GetPerspectiveQueue() (q *Queue, err error) {
	q = &Queue{}
	err = db.Preload("Perspective").Joins("JOIN perspective ON queue.perspective_id = perspective.id AND perspective.idx = ?", int(idx)).First(q).Error
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
	err = db.
		Preload("Perspective").
		Joins("JOIN perspective ON queue.perspective_id = perspective.id").
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
		err = q.appendTo(&qt)
		onerror.Log(err)
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

	for i := range s {
		if s[i].Position == position {
			s[i].Played = true
			break
		}
	}

	s = reasignQueueTrackPositions(s)

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
		err := q.insertInto(&qt)
		onerror.Log(err)
	}
	for i := len(locations) - 1; i >= 0; i-- {
		qt := QueueTrack{Location: locations[i], Position: position}
		err := q.insertInto(&qt)
		onerror.Log(err)
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

	var moved, afterPiv []QueueTrack
	var piv *QueueTrack
	for i := range s {
		if s[i].Position == from {
			piv = &s[i]
		} else if s[i].Position < to {
			moved = append(moved, s[i])
		} else if s[i].Position > to {
			afterPiv = append(afterPiv, s[i])
		} else if s[i].Position == to {
			if from < to {
				moved = append(moved, s[i])
			} else {
				afterPiv = append(afterPiv, s[i])
			}
		}
	}

	if piv != nil {
		moved = append(moved, *piv)
	}
	moved = append(moved, afterPiv...)

	moved = reasignQueueTrackPositions(moved)

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

	for i := range s {
		if s[i].Position == 1 {
			s[i].Played = true
			qt = &s[i]
			break
		}
	}
	if qt == nil {
		return
	}
	s = reasignQueueTrackPositions(s)

	log.Info("Found location to pop from queue:", qt.Location)
	err = db.Save(&s).Error

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

	s = append(s, *qt)
	s = reasignQueueTrackPositions(s)
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

	if qt.Position <= 1 {
		aux := s
		s = []QueueTrack{*qt}
		s = append(s, aux...)
	} else if qt.Position > 1 && qt.Position <= len(s) {
		aux := s
		piv := int(qt.Position) - 1
		s = aux[:piv]
		s = append(s, *qt)
		s = append(s, aux[piv:]...)
	} else {
		err = q.appendTo(qt)
		return
	}

	s = reasignQueueTrackPositions(s)
	if err = db.Save(&s).Error; err != nil {
		return
	}

	go findQueueTrack(qt)
	return
}

// QueueTrack defines a track in the queue
type QueueTrack struct { // too transient
	ID        int64  `json:"id" gorm:"primaryKey"`
	Position  int    `json:"position"`
	Played    bool   `json:"played"`
	Location  string `json:"location" gorm:"not null"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:nano"`
	TrackID   int64  `json:"trackId" gorm:"index:idx_queue_track_track_id"`
	QueueID   int64  `json:"queueId" gorm:"index:idx_queue_track_queue_id,not null"`
	Queue     Queue  `json:"queue" gorm:"foreignKey:QueueID"`
}

// Create implements the DataCreator interface
func (qt *QueueTrack) Create() (err error) {
	err = db.Create(qt).Error
	return
}

// Save implements the DataUpdater interface
func (qt *QueueTrack) Save() (err error) {
	err = db.Save(qt).Error
	return
}

// ToProtobuf implements the ProtoOut interface
func (qt *QueueTrack) ToProtobuf() proto.Message {
	bv, err := json.Marshal(qt)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.QueueTrack{}
	}

	out := &m3uetcpb.QueueTrack{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)

	// Unmatched
	q := Queue{}
	err = q.Read(qt.QueueID)
	onerror.Log(err)

	out.Perspective = m3uetcpb.Perspective(q.Perspective.Idx)
	out.TrackId = qt.TrackID
	out.CreatedAt = qt.CreatedAt
	out.UpdatedAt = qt.UpdatedAt
	return out
}

// AfterCreate is a GORM hook
func (qt *QueueTrack) AfterCreate(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode && !base.IsAppBusyBy(base.IdleStatusEngineLoop) {
			PlaybackChanged <- struct{}{}
		}
	}()
	return nil
}

// GetAllQueueTracks returns all queue tracks for the given perspective, constrained by a limit
func GetAllQueueTracks(idx PerspectiveIndex, limit int) (qts []*QueueTrack, ts []*Track) {
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

	if limit > 0 {
		tx.Limit(limit)
	}

	qts = []*QueueTrack{}
	ts = []*Track{}
	list := []QueueTrack{}
	if err = tx.Find(&list).Error; err != nil {
		log.Error(err)
		return
	}

	ids := []int64{}
	locations := []string{}
	for i := range list {
		if list[i].TrackID > 0 {
			ids = append(ids, list[i].TrackID)
		} else {
			locations = append(locations, list[i].Location)
		}
		qts = append(qts, &list[i])
	}

	ts, _ = FindTracksIn(ids)
	for _, l := range locations {
		t := &Track{}
		if t, err = ReadTagsForLocation(l); err != nil {
			continue
		}
		ts = append(ts, t)
	}
	return
}

// GetQueueStore returns all queue tracks for all perspectives
func GetQueueStore() (qs []*QueueTrack, ts []*Track) {
	log.Info("Getting queue store")

	for _, idx := range PerspectiveIndexList() {
		qout, tout := GetAllQueueTracks(idx, 0)
		for i := range qout {
			qs = append(qs, qout[i])
		}
		for i := range tout {
			ts = append(ts, tout[i])
		}
	}
	return
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
	err := qt.Save()
	onerror.Log(err)
	return
}

func reasignQueueTrackPositions(s []QueueTrack) []QueueTrack {
	pos := 0
	for i := range s {
		if s[i].Played {
			continue
		}
		pos++
		s[i].Position = pos
	}
	return s
}
