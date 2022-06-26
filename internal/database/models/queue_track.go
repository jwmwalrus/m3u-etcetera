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
		if !base.FlagTestingMode &&
			!base.IsAppBusyBy(base.IdleStatusEngineLoop) {
			PlaybackChanged <- struct{}{}
		}
	}()
	return nil
}

// GetPosition implements the Poser interface
func (qt *QueueTrack) GetPosition() int {
	return qt.Position
}

// SetPosition implements the Poser interface
func (qt *QueueTrack) SetPosition(pos int) {
	qt.Position = pos
}

// GetIgnore implements the Poser interface
func (qt *QueueTrack) GetIgnore() bool {
	return qt.Played
}

// SetIgnore implements the Poser interface
func (qt *QueueTrack) SetIgnore(ignore bool) {
	qt.Played = ignore
}

// GetAllQueueTracks returns all queue tracks for the given perspective,
// constrained by a limit
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
	s := []QueueTrack{}
	if err = tx.Find(&s).Error; err != nil {
		log.Error(err)
		return
	}

	ids := []int64{}
	locations := []string{}
	for i := range s {
		if s[i].TrackID > 0 {
			ids = append(ids, s[i].TrackID)
		} else {
			locations = append(locations, s[i].Location)
		}
		qts = append(qts, &s[i])
	}

	ts, _ = FindTracksIn(ids)
	for _, l := range locations {
		var t *Track
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
		qsaux, tsaux := GetAllQueueTracks(idx, 0)
		for i := range qsaux {
			qs = append(qs, qsaux[i])
		}
		for i := range tsaux {
			ts = append(ts, tsaux[i])
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
}
