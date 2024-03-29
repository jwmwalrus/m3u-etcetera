package database

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/migrations"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	rtc "github.com/jwmwalrus/rtcycler"

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
	return Path() + ConnectionOptions()
}

// Instance returns the database instance.
func Instance() *gorm.DB {
	return conn
}

// Open creates the application database, if it doesn't exist.
func Open() *rtc.Unloader {
	logw := slog.With("dsn", DSN())
	onerrorw := onerror.NewRecorder(logw)

	var err error

	backupDatabase()

	conn, err = gorm.Open(sqlite.Open(DSN()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DryRun: rtc.FlagDry(),
		Logger: logger.Default.LogMode(logger.Silent),
	})
	onerrorw.Fatal(err)

	modelsCtx, modelsCancel = context.WithCancel(context.Background())
	models.SetUp(modelsCtx, conn)

	// Migrate the schema
	m := gormigrate.New(conn, gormigrate.DefaultOptions, migrations.All())

	m.InitSchema(migrations.InitSchema)
	onerrorw.Fatal(m.Migrate())

	go models.DoInitialCleanup()

	logw.Info("Database loaded")

	return unloader
}

// ConnectionOptions returns the database directory.
func ConnectionOptions() string {
	if !rtc.FlagTestMode() {
		return connectionOptions
	}
	return "?_loc=Local"
}

// Dir returns the database directory.
func Dir() string {
	if !rtc.FlagTestMode() {
		return rtc.DataDir()
	}
	return os.TempDir()
}

// Path returns the full path to the database {.
func Path() string {
	return filepath.Join(Dir(), Filename())
}

// Filename returns the database filename.
func Filename() string {
	if !rtc.FlagTestMode() {
		return base.DatabaseFilename
	}
	return "m3uetc-test-music-" + rtc.InstanceSuffix() + ".db"
}

func backupDatabase() {
	if !rtc.FlagTestMode() &&
		base.Conf.Server.Database.Backup {
		path := Path()
		_, err := os.Stat(path)
		if err == nil {
			if _, err := exec.LookPath("sqlite3"); err != nil {
				slog.Error("Failed to find sqlite3", "error", err)
				return
			}

			out, err := exec.
				Command("sqlite3", path, ".backup "+path+".bak").
				CombinedOutput()
			if err != nil {
				slog.With(
					"output", out,
					"error", err,
				).Error("Failed to backup database")
				return
			}
		}
	}
}
