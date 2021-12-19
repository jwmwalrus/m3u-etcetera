package database

import (
	"os"
	"os/exec"

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

const (
	connectionOptions = "?_foreign_keys=1&_loc=Local"
)

var (
	conn *gorm.DB
)

// DSN returns the application's DSN
func DSN() string {
	return base.GetDatabasePath() + connectionOptions
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
		path := base.GetDatabasePath()
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
