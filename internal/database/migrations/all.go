package migrations

import "github.com/go-gormigrate/gormigrate/v2"

// All -.
func All() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		// you migrations here
		// {
		// 	ID:"datetime",
		// 	Migrate: func(db *gorm.DB) err,
		// 	Rollback: func(db *gorm.DB) err,
		// },
	}
}
