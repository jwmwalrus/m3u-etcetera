package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
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
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:nano"`
	TrackID   int64  `json:"trackId"`
}

// Create implements the DataCreator interface
func (pb *Playback) Create() (err error) {
	err = db.Create(pb).Error
	return
}

// Read implements the DataReader interface
func (pb *Playback) Read(id int64) (err error) {
	err = db.First(pb, id).Error
	return
}

// Save implements the DataUpdater interface
func (pb *Playback) Save() (err error) {
	err = db.Save(pb).Error
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

	// Unmatched
	out.TrackId = pb.TrackID
	out.CreatedAt = pb.CreatedAt
	out.UpdatedAt = pb.UpdatedAt
	return out
}

// AddToHistory adds unplayed playback to history and marks it as played
func (pb *Playback) AddToHistory(position, duration int64, freeze bool) {
	log.WithFields(log.Fields{
		"pb":       *pb,
		"position": position,
		"duration": duration,
	}).
		Info("Adding playback to history")

	if !freeze {
		pb.Played = true
	}
	pb.Skip = position
	pb.Save()
	go func() {
		h := PlaybackHistory{
			Location: pb.Location,
			TrackID:  pb.TrackID,
			Duration: position,
		}
		onerror.Log(h.Create())
		if pb.TrackID > 0 {
			t := &Track{}
			t, err := pb.GetTrack()
			if err != nil {
				log.Error(err)
				return
			}

			if time.Duration(position) >= time.Duration(base.Conf.Server.Playback.PlayedThreshold)*time.Second {
				t.Lastplayed = h.CreatedAt
				t.Playcount++
			}
			if t.Duration == 0 {
				t.Duration = duration
			}
			onerror.Warn(t.Save())
		}
	}()
	return
}

// BeforeCreate is a GORM hook
func (pb *Playback) BeforeCreate(tx *gorm.DB) error {
	fmt.Println("BeforeCreate location:", pb.Location)
	return nil
}

// BeforeSave is a GORM hook
func (pb *Playback) BeforeSave(tx *gorm.DB) error {
	fmt.Println("BeforeSave location:", pb.Location)
	return nil
}

// AfterCreate is a GORM hook
func (pb *Playback) AfterCreate(tx *gorm.DB) error {
	fmt.Println("AfterCreate location:", pb.Location)
	go func() {
		if !base.FlagTestingMode && !base.IsAppBusyBy(base.IdleStatusEngineLoop) {
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
	onerror.Log(pb.Save())
	return
}

// GetNextToPlay returns the next playback entry to play
func (pb *Playback) GetNextToPlay() (err error) {
	err = db.Where("played = 0").First(pb).Error
	if err == nil && pb.ID == 0 {
		err = fmt.Errorf("There's nothing in the playback queue")
	}
	return
}

// GetTrack returns the track for the given playback
func (pb *Playback) GetTrack() (t *Track, err error) {
	t = &Track{}
	err = db.First(t, pb.TrackID).Error
	return
}

// PlaybackHistory defines a playback_history row
type PlaybackHistory struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	Location  string `json:"location" gorm:"index:idx_playback_history_location, not null"`
	Duration  int64  `json:"duration"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:nano,index:idx_playback_history_created_at"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:nano"`
	TrackID   int64  `json:"trackId" gorm:"index:idx_playback_history_track_id"`
}

// Create implements the DataCreator interface
func (h *PlaybackHistory) Create() (err error) {
	err = db.Create(h).Error
	return
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

	pb = &Playback{Location: location}
	if err := pb.Create(); err != nil {
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
	onerror.Log(pb.Create())
	return
}
