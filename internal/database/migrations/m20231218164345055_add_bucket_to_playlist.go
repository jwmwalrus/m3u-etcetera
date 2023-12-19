package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

func m20231218164345055_add_bucket_to_playlist() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231218164345055",

		Migrate: func(tx *gorm.DB) error {
			err := tx.Migrator().AddColumn(&models.Playlist{}, "Bucket")
			if err != nil {
				return err
			}

			return tx.Exec("UPDATE playlist SET bucket = 0 WHERE id > 0").Error
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn("playlist", "bucket")
		},
	}
}
