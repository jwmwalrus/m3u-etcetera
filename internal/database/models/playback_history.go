package models

import (
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

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

// AddPlaybackToHistory adds unplayed playback to history and marks it as played
func AddPlaybackToHistory(id, position, duration int64, freeze bool) {
	pb := Playback{}
	if err := pb.Read(id); err != nil {
		log.Errorf("Error adding playback to history: %v", err)
		return
	}

	log.WithFields(log.Fields{
		"location": pb.Location,
		"trackId":  pb.TrackID,
		"position": position,
		"duration": duration,
	}).
		Info("Adding playback to history")

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
	onerror.Log(h.Create())
	if pb.TrackID > 0 {
		t := &Track{}
		t, err := pb.GetTrack()
		if err != nil {
			log.Error(err)
			return
		}

		if time.Duration(position)*time.Nanosecond >=
			time.Duration(base.Conf.Server.Playback.PlayedThreshold)*time.Second {
			t.Lastplayed = h.CreatedAt
			t.Playcount++
		}
		if t.Duration == 0 {
			t.Duration = duration
		}
		onerror.Warn(t.Save())
	}
}
