package database

import (
	"context"
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/migrations"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/onerror"
	rtc "github.com/jwmwalrus/rtcycler"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	connectionOptions = "?_foreign_keys=1&_loc=Local"
)

var (
	conn *gorm.DB

	modelsCtx    context.Context
	modelsCancel context.CancelFunc

	unloader = &rtc.Unloader{
		Description: "CloseDatabase",
		Callback: func() error {
			Close()
			return nil
		},
	}
)

// Close closes the application database.
func Close() {
	if conn == nil {
		return
	}

	modelsCancel()

	models.TearDown()

	var err error
	var sqlDB *sql.DB
	if sqlDB, err = conn.DB(); err != nil {
		panic(err)
	}
	sqlDB.Close()
}

// DSN returns the application's DSN.
func DSN() string {
	return getDatabasePath() + getConnectionOptions()
}

// Instance returns the database instance.
func Instance() *gorm.DB {
	return conn
}

// Open creates the application database, if it doesn't exist.
func Open() *rtc.Unloader {
	entry := log.WithField("dsn", DSN())

	var err error

	backupDatabase()

	conn, err = gorm.Open(sqlite.Open(DSN()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DryRun: rtc.FlagDry(),
		Logger: logger.Default.LogMode(logger.Silent),
	})
	onerror.WithEntry(entry).Panic(err)

	// TODO: connect with logrus

	modelsCtx, modelsCancel = context.WithCancel(context.Background())
	models.SetUp(modelsCtx, conn)

	// Migrate the schema
	m := gormigrate.New(conn, gormigrate.DefaultOptions, migrations.All())

	m.InitSchema(migrations.InitSchema)
	onerror.WithEntry(entry).Panic(m.Migrate())

	go models.DoInitialCleanup()

	log.WithField("dsn", DSN()).
		Info("Database loaded")

	return unloader
}

func backupDatabase() {
	if !rtc.FlagTestMode() &&
		base.Conf.Server.Database.Backup {
		path := getDatabasePath()
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			if _, err := exec.LookPath("sqlite3"); err != nil {
				log.Error(err)
				return
			}

			out, err := exec.
				Command("sqlite3", path, ".backup "+path+".bak").
				CombinedOutput()
			if err != nil {
				log.WithField("output", out).Error(err)
				return
			}
		}
	}
}

// getConnectionOptions returns the database directory
func getConnectionOptions() string {
	if !rtc.FlagTestMode() {
		return connectionOptions
	}
	return "?_loc=Local"
}

// getDatabaseDir returns the database directory
func getDatabaseDir() string {
	if !rtc.FlagTestMode() {
		return rtc.DataDir()
	}
	return os.TempDir()
}

// getDatabasePath returns the full path to the database {
func getDatabasePath() string {
	return filepath.Join(getDatabaseDir(), getDatabaseFilename())
}

// getDatabaseFilename returns the database filename
func getDatabaseFilename() string {
	if !rtc.FlagTestMode() {
		return base.DatabaseFilename
	}
	return "m3uetc-test-music-" + rtc.InstanceSuffix() + ".db"
}
