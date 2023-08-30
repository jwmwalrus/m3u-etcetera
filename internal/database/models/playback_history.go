package models

import (
	"log/slog"
	"time"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"gorm.io/gorm"
)

// PlaybackHistory defines a playback_history row.
type PlaybackHistory struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	Location  string `json:"location" gorm:"index:idx_playback_history_location, not null"`
	Duration  int64  `json:"duration"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime:nano"`
	TrackID   int64  `json:"trackId" gorm:"index:idx_playback_history_track_id"`
}

// Create implements the Creator interface.
func (h *PlaybackHistory) Create() error {
	return h.CreateTx(db)
}

// CreateTx implements the Creator interface.
func (h *PlaybackHistory) CreateTx(tx *gorm.DB) error {
	return tx.Create(h).Error
}

// FindLastBy returns the newest entry in the playback history,
// according to the given query.
func (h *PlaybackHistory) FindLastBy(query interface{}) (err error) {
	err = db.Where(query).Last(h).Error
	return
}

// ReadLast returns the newest entry in the playback history.
func (h *PlaybackHistory) ReadLast() (err error) {
	err = db.Last(&h).Error
	return
}

// AddPlaybackToHistory adds unplayed playback to history and marks it as played.
func AddPlaybackToHistory(id, position, duration int64, freeze bool) {
	storageGuard <- struct{}{}
	defer func() { <-storageGuard }()

	logw := slog.With(
		"position", position,
		"duration", duration,
	)

	pb := Playback{}
	if err := pb.Read(id); err != nil {
		logw.Error("Failed add playback to history: failed to read playback", "error", err)
		return
	}

	logw = logw.With(
		"location", pb.Location,
		"track_id", pb.TrackID,
	)
	logw.Info("Adding playback to history")

	if !freeze {
		pb.Played = true
	}
	pb.Skip = position
	pb.Save()

	h := PlaybackHistory{
		Location: pb.Location,
		TrackID:  pb.TrackID,
		Duration: position,
	}

	onerrorw := onerror.NewRecorder(logw)

	onerrorw.Log(h.Create())
	if pb.TrackID > 0 {
		t, err := pb.GetTrack()
		if err != nil {
			logw.Error("Failed to get playback track", "error", err)
			return
		}

		if time.Duration(position)*time.Nanosecond >=
			time.Duration(base.PlaybackPlayedThreshold)*time.Second {

			t.Lastplayed = h.CreatedAt
			t.Playcount++
		}
		if t.Duration == 0 {
			t.Duration = duration
		}
		onerrorw.Warn(t.Save())
	}
}
