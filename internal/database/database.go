package database

import (
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/migrations"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/onerror"

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

	// Unloader -
	Unloader = base.Unloader{
		Description: "CloseDatabase",
		Callback: func() error {
			Close()
			return nil
		},
	}
)

// Close closes the application database
func Close() {
	if conn == nil {
		return
	}

	var err error
	var sqlDB *sql.DB
	if sqlDB, err = conn.DB(); err != nil {
		panic(err)
	}
	sqlDB.Close()
}

// DSN returns the application's DSN
func DSN() string {
	return getDatabasePath() + getConnectionOptions()
}

// Open creates the application database, if it doesn't exist
func Open() *gorm.DB {
	var err error

	backupDatabase()

	conn, err = gorm.Open(sqlite.Open(DSN()), &gorm.Config{
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
		conn.Where("played = 1").Delete(models.Queue{})
	}()

	log.WithField("dsn", DSN()).
		Info("Database loaded")

	return conn
}

func backupDatabase() {
	if !base.FlagTestingMode &&
		base.Conf.Server.Database.Backup {
		path := getDatabasePath()
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			if _, err := exec.LookPath("sqlite3"); err != nil {
				log.Error(err)
				return
			}

			if out, err := exec.Command("sqlite3", path, ".backup "+path+".bak").CombinedOutput(); err != nil {
				log.WithField("output", out).Error(err)
				return
			}
		}
	}
}

// getConnectionOptions returns the database directory
func getConnectionOptions() string {
	if !base.FlagTestingMode {
		return connectionOptions
	}
	return "?_loc=Local"
}

// getDatabaseDir returns the database directory
func getDatabaseDir() string {
	if !base.FlagTestingMode {
		return base.DataDir
	}
	return os.TempDir()
}

// getDatabasePath returns the full path to the database {
func getDatabasePath() string {
	return filepath.Join(getDatabaseDir(), getDatabaseFilename())
}

// getDatabaseFilename returns the database filename
func getDatabaseFilename() string {
	if !base.FlagTestingMode {
		return base.DatabaseFilename
	}
	return "m3uetc-test-music-" + base.InstanceSuffix + ".db"
}
