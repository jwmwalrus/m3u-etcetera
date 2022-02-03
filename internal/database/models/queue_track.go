package models

import (
	"encoding/json"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

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
	onerror.Log(q.Read(qt.QueueID))

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
	onerror.Log(qt.Save())
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
