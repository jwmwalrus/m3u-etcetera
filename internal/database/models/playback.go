package models

import (
	"encoding/json"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// Playback defines a playback row
type Playback struct { // too transient
	ID        int64  `json:"id" gorm:"primaryKey"`
	Location  string `json:"location" gorm:"not null"`
	Played    bool   `json:"played"`
	Skip      int64  `json:"skip"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime"`
	TrackID   int64  `json:"trackId"`
}

// Read implements the DataReader interface
func (pb *Playback) Read(id int64) (err error) {
	err = db.First(pb, id).Error
	return
}

// ToProtobuf implments ProtoOut interface
func (pb *Playback) ToProtobuf() proto.Message {
	bv, err := json.Marshal(pb)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.Playback{}
	}

	out := &m3uetcpb.Playback{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)
	return out
}

// AddToHistory adds unplayed playback to history and marks it as played
func (pb *Playback) AddToHistory(duration int64) {
	log.WithFields(log.Fields{
		"pb":       *pb,
		"duration": duration,
	}).
		Info("Adding playback to history")

	pb.Played = true
	db.Save(&pb)
	go func() {
		h := PlaybackHistory{
			Location: pb.Location,
			TrackID:  pb.TrackID,
		}
		err := db.Create(&h).Error
		onerror.Log(err)
		if pb.TrackID > 0 {
			t := Track{}
			if err := db.First(&t, pb.TrackID).Error; err != nil {
				log.Error(err)
				return
			}
			t.Lastplayed = h.CreatedAt
			t.Playcount++
			err = db.Save(&t).Error
			onerror.Warn(err)
		}
	}()
	return
}

// AfterCreate is a GORM hook
func (pb *Playback) AfterCreate(tx *gorm.DB) error {
	go func() {
		if !base.IsAppBusyBy(base.IdleStatusEngineLoop) {
			PlaybackChanged <- struct{}{}
		}
	}()
	return nil
}

// ClearPending removes all pending playback entries
func (pb *Playback) ClearPending() {
	err := db.Model(&Playback{}).Where("played = 0 AND id <> ?", pb.ID).Update("played", 1).Error
	onerror.Warn(err)
}

// FindTrack attempts to find track from location
func (pb *Playback) FindTrack() {
	if pb.TrackID > 0 {
		return
	}

	log.WithField("pb", *pb).
		Info("Finding track for current playback")

	t := Track{}
	if err := db.Where("location = ?", pb.Location).First(&t).Error; err != nil {
		return
	}
	pb.TrackID = t.ID
	err := db.Save(pb).Error
	onerror.Log(err)
	return
}

// GetNextToPlay returns the next playback entry to play
func (pb *Playback) GetNextToPlay() (err error) {
	err = db.Where("played = 0").First(pb).Error
	return
}

// GetTrack returns the track for the given playback
func (pb *Playback) GetTrack() (t *Track, err error) {
	t = &Track{}
	err = db.First(&t, pb.TrackID).Error
	return

}

// PlaybackHistory defines a playback_history row
type PlaybackHistory struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	Location  string `json:"location" gorm:"index:idx_playback_history_location, not null"`
	Duration  int64  `json:"duration"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime,index:idx_playback_history_created_at"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime"`
	TrackID   int64  `json:"trackId" gorm:"index:idx_playback_history_track_id"`
}

// FindLastBy returns the newest entry in the playback history, according to the given query
func (h *PlaybackHistory) FindLastBy(query interface{}) (err error) {
	err = db.Where(query).Last(h).Error
	return
}

// ReadLast returns the newest entry in the playback history
func (h *PlaybackHistory) ReadLast() (err error) {
	err = db.Last(&h).Error
	return
}

// AddPlaybackLocation adds a playback entry by location
func AddPlaybackLocation(location string) (pb *Playback) {
	log.WithField("location", location).
		Info("Adding playback entry by location")

	if !base.IsSupportedURL(location) {
		log.WithField("location", location).
			Error("The given location is unsupported for playback")

		return
	}

	pb = &Playback{Location: location}
	if err := db.Create(pb).Error; err != nil {
		log.Error(err)
		return
	}
	go pb.FindTrack()
	return
}

// AddPlaybackTrack adds a playback entry by track
func AddPlaybackTrack(t *Track) (pb *Playback) {
	log.WithField("t", t).
		Info("Adding playback entry by track")

	pb = &Playback{Location: t.Location, TrackID: t.ID}
	err := db.Create(pb).Error
	onerror.Log(err)
	return
}
