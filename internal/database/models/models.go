package models

import "gorm.io/gorm"

var (
	db *gorm.DB

	// make sure we only do heavy collection tasks one at a time
	storageGuard chan struct{}

	// DbgChan is a debug channel
	DbgChan chan map[string]interface{}

	// PlaybackChanged is the AfterCreate-hook channel for QueueTrack and Playback
	PlaybackChanged chan struct{}
)

// SetConnection sets the database connection for the whole package
func SetConnection(conn *gorm.DB) {
	db = conn
}

func init() {
	DbgChan = make(chan map[string]interface{})
	storageGuard = make(chan struct{}, 1)
	PlaybackChanged = make(chan struct{}, 1)
}
