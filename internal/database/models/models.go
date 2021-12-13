package models

import "gorm.io/gorm"

var (
	db *gorm.DB

	// make sure we only do heavy collection tasks one at a time
	storageGuard chan struct{}

	// DbgChan is a debug channel
	DbgChan chan string
)

// SetConnection sets the database connection for the whole package
func SetConnection(conn *gorm.DB) {
	db = conn
}

func init() {
	DbgChan = make(chan string)
	storageGuard = make(chan struct{}, 1)
}
