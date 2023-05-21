package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m20230515223654066_add_lastplayedfor_to_playlist_track() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230515223654066",

		Migrate: func(tx *gorm.DB) error {
			return tx.Exec("ALTER TABLE playlist_track ADD COLUMN lastplayedfor INTEGER DEFAULT 0").Error
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn("playlist_track", "lastplayedfor")
		},
	}
}
