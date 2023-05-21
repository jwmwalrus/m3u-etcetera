package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/seeds"
	"gorm.io/gorm"
)

func m20230514080315863_add_read_only_queries() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230514080315863",

		Migrate: func(tx *gorm.DB) error {
			err := tx.Exec("ALTER TABLE query ADD COLUMN idx INTEGER NOT NULL DEFAULT 0").Error
			if err != nil {
				return err
			}

			return seeds.SeedQuery(tx)
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn("query", "idx")
		},
	}
}
