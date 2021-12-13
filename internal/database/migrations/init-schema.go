package migrations

import (
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/seeds"
	"gorm.io/gorm"
)

// InitSchema initializes schema
func InitSchema(db *gorm.DB) (err error) {
	err = db.AutoMigrate(
		// no foreign keys
		&models.Track{},
		&models.Perspective{},

		// soft reference
		&models.Playback{},
		&models.PlaybackHistory{},

		// one foreign key
		&models.Queue{},

		// foreign key in previous group
		&models.QueueTrack{},
	)
	onerror.Panic(err)

	seeds.All(db)

	return
}
