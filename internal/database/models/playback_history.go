package models

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
