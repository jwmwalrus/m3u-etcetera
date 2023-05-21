package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m20230515200631346_add_query_id_to_playlist() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230515200631346",

		Migrate: func(tx *gorm.DB) error {
			return tx.Exec("ALTER TABLE playlist ADD COLUMN query_id INTEGER DEFAULT 0").Error
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn("playlist", "query_id")
		},
	}
}
