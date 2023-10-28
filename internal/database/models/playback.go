package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/pointers"
	"github.com/jwmwalrus/gear-pieces/idler"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	rtc "github.com/jwmwalrus/rtcycler"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// Playback defines a playback row.
type Playback struct { // too transient
	ID        int64  `json:"id" gorm:"primaryKey"`
	Location  string `json:"location" gorm:"not null"`
	Played    bool   `json:"played"`
	Skip      int64  `json:"skip"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:nano"`
	TrackID   int64  `json:"trackId"`
}

func (pb *Playback) Create() error {
	return pb.CreateTx(db)
}

func (pb *Playback) CreateTx(tx *gorm.DB) error {
	return tx.Create(pb).Error
}

func (pb *Playback) Read(id int64) error {
	return pb.ReadTx(db, id)
}

func (pb *Playback) ReadTx(tx *gorm.DB, id int64) error {
	return tx.First(pb, id).Error
}

func (pb *Playback) Save() error {
	return pb.SaveTx(db)
}

func (pb *Playback) SaveTx(tx *gorm.DB) error {
	return tx.Save(pb).Error
}

func (pb *Playback) ToProtobuf() proto.Message {
	bv, err := json.Marshal(pb)
	if err != nil {
		slog.Error("Failed to marshal playback", "error", err)
		return &m3uetcpb.Playback{}
	}

	out := &m3uetcpb.Playback{}
	err = jsonUnmarshaler.Unmarshal(bv, out)
	onerror.Log(err)

	return out
}

// AfterCreate is a GORM hook.
func (pb *Playback) AfterCreate(tx *gorm.DB) error {
	go func() {
		if !rtc.FlagTestMode() &&
			!idler.IsAppBusyBy(idler.StatusEngineLoop) {
			TriggerPlaybackChange()
		}
	}()
	return nil
}

func (pb *Playback) Blacklist() {
	err := db.Model(&Playback{}).
		Where("played = 0 AND id = ?", pb.ID).
		Update("played", 1).
		Error
	onerror.Warn(err)

	if pb.TrackID > 0 {
		DeleteLocalTrackIfDangling(pb.TrackID, pb.Location)
	}
}

// ClearPending removes all pending playback entries.
func (pb *Playback) ClearPending() {
	err := db.Model(&Playback{}).
		Where("played = 0 AND id <> ?", pb.ID).
		Update("played", 1).
		Error
	onerror.Warn(err)
}

// GetNextToPlay returns the next playback entry to play.
func (pb *Playback) GetNextToPlay() (err error) {
	err = db.Where("played = 0").First(pb).Error
	if err == nil && pb.ID == 0 {
		err = fmt.Errorf("There's nothing in the playback queue")
	}
	return
}

// GetTrack returns the track for the given playback.
func (pb *Playback) GetTrack() (t *Track, err error) {
	t = &Track{}
	err = db.First(t, pb.TrackID).Error
	return
}

// AddPlaybackLocation adds a playback entry by location.
func AddPlaybackLocation(location string) (pb *Playback) {
	logw := slog.With("location", location)
	logw.Info("Adding playback entry by location")

	pb = &Playback{Location: location}
	if err := pb.Create(); err != nil {
		logw.Error("Failed to create playback", "error", err)
		return
	}
	go func() {
		playbackTrackNeeded <- pb.ID
	}()
	return
}

// AddPlaybackTrack adds a playback entry by track.
func AddPlaybackTrack(t *Track) (pb *Playback) {
	logw := slog.With("track", t)
	logw.Info("Adding playback entry by track")

	pb = &Playback{Location: t.Location, TrackID: t.ID}
	onerror.NewRecorder(logw).Log(pb.Create())
	return
}

// GetAllPlayback returns all the playback entries.
func GetAllPlayback() []*Playback {
	slog.Info("Obtaining all playback")

	pbs := []Playback{}
	err := db.Where("played = 0").Find(&pbs).Error
	if err != nil {
		slog.Error("Failed to find unplayed placybacks in database", "error", err)
		return []*Playback{}
	}

	return pointers.FromSlice(pbs)
}

func findPlaybackTrack(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case id := <-playbackTrackNeeded:
			go func(id int64) {
				pb := Playback{}
				err := pb.Read(id)
				if err != nil {
					slog.With(
						"id", id,
						"error", err,
					).Error("Failed to read playback")
					return
				}
				if pb.TrackID > 0 {
					return
				}

				logw := slog.With("pb", pb)
				logw.Info("Finding track for current playback")

				t := Track{}
				err = db.Where("location = ?", pb.Location).First(&t).Error
				if err != nil {
					return
				}
				pb.TrackID = t.ID
				onerror.NewRecorder(logw).Log(pb.Save())
			}(id)

		}
	}
}
