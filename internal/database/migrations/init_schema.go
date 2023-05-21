package migrations

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/seeds"
	"github.com/jwmwalrus/onerror"
	"gorm.io/gorm"
)

// InitSchema initializes schema.
func InitSchema(db *gorm.DB) (err error) {
	err = db.AutoMigrate(
		// no foreign keys
		&models.Collection{},
		&models.Query{},
		&models.Perspective{},

		// soft reference
		&models.Playback{},
		&models.PlaybackHistory{},

		// one foreign key
		&models.Track{},
		&models.Playbar{},
		&models.PlaylistGroup{},
		&models.Queue{},

		// one foreign key, one soft reference
		&models.QueueTrack{},

		// two foreign keys
		&models.CollectionQuery{},
		&models.Playlist{},
		&models.PlaylistQuery{},
		&models.PlaylistTrack{},
	)
	onerror.Panic(err)

	seeds.All(db)

	return
}
