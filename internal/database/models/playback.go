package models

import (
	"encoding/json"
	"fmt"

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

// AfterCreate is a GORM hook
func (pb *Playback) AfterCreate(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode &&
			!base.IsAppBusyBy(base.IdleStatusEngineLoop) {
			PlaybackChanged <- struct{}{}
		}
	}()
	return nil
}

// ClearPending removes all pending playback entries
func (pb *Playback) ClearPending() {
	err := db.Model(&Playback{}).
		Where("played = 0 AND id <> ?", pb.ID).
		Update("played", 1).
		Error
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
	err := db.Where("location = ?", pb.Location).First(&t).Error
	if err != nil {
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

// GetAllPlayback returns all the playback entries
func GetAllPlayback() (pbs []*Playback) {
	log.Info("Obtaining all playback")

	pblist := []Playback{}
	err := db.Where("played = 0").Find(&pblist).Error
	if err != nil {
		log.Error(err)
		return
	}

	pbs = []*Playback{}
	for i := range pblist {
		pbs = append(pbs, &pblist[i])
	}
	return
}
