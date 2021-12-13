package database

import (
	"path/filepath"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/migrations"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"

	// _ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var conn *gorm.DB

// DSN returns the application's DSN
func DSN(prod bool) string {
	if prod {
		return filepath.Join(base.DataDir, base.DatabaseFilename) + "?_foreign_keys=1&_loc=Local"
	}
	return "/tmp/m3uetc-test-music.db?_foreign_keys=1&_loc=Local"
}

// Open creates the application database, if it doesn't exist
func Open() *gorm.DB {
	var err error

	conn, err = gorm.Open(sqlite.Open(DSN(true)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DryRun: base.FlagDry,
		Logger: logger.Default.LogMode(logger.Silent),
	})
	onerror.Panic(err)

	models.SetConnection(conn)

	// TODO: connect with logrus

	// Migrate the schema
	m := gormigrate.New(conn, gormigrate.DefaultOptions, migrations.All())

	m.InitSchema(migrations.InitSchema)
	err = m.Migrate()
	onerror.Panic(err)

	go func() {
		conn.Where("played = 1").Delete(models.Playback{})
	}()

	log.Info("Database loaded")

	return conn
}
