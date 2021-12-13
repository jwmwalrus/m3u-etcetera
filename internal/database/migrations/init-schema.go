package migrations

import (
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

// InitSchema initializes schema
func InitSchema(db *gorm.DB) (err error) {
	err = db.AutoMigrate(
		// no foreign keys
		&models.Track{},

		// soft reference
		&models.Playback{},
		&models.PlaybackHistory{},
	)
	onerror.Panic(err)

	return
}
