package api

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

const (
	testDataDir = "../data/testing"
)

var (
	conn        *gorm.DB
	fixtures    *testfixtures.Loader
	watchingDbg = false
)

type testCase struct {
	name        string
	fixturesDir string
	startEngine bool
	sleep       int
	req         proto.Message
	res         proto.Message
	err         bool
}

func closeTestDatabase() {
	database.Close()
	/*
		var err error
		var sqlDB *sql.DB
		if sqlDB, err = conn.DB(); err != nil {
			panic(err)
		}
		sqlDB.Close()
	*/
}

func openTestDatabase(fixturesDir string) (db *gorm.DB) {
	var err error
	/*
		db, err = gorm.Open(sqlite.Open(database.DSN()), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: logger.Default.LogMode(logger.Silent),
		})
		onerror.Panic(err)

		models.SetConnection(conn)

		m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.All())

		m.InitSchema(migrations.InitSchema)
		err = m.Migrate()
		if err != nil {
			panic(err)
		}
	*/

	db = database.Open()

	sqlDB, _ := db.DB()
	fixtures, err = testfixtures.New(
		testfixtures.Database(sqlDB),
		testfixtures.Dialect("sqlite"),
		testfixtures.Paths(filepath.Join(testDataDir, "fixtures", fixturesDir)),
	)
	if err != nil {
		panic(err)
	}

	err = fixtures.Load()
	if err != nil {
		panic(err)
	}

	return
}

func setupTest(t *testing.T, tc testCase) {
	base.SetTestingMode()
	conn = openTestDatabase(tc.fixturesDir)

	if tc.startEngine {
		playback.SetMode(playback.TestMode)
		playback.StartEngine()
	}

	watchingDbg = true
	go watchDebugChannel(t)
	return
}

func teardownTest(t *testing.T) {
	watchingDbg = false

	playback.StopEngine()
	closeTestDatabase()

	if _, err := os.Stat(database.DSN()); !os.IsNotExist(err) {
		if err = os.Remove(database.DSN()); err != nil {
			panic(err)
		}
	}
}

func watchDebugChannel(t *testing.T) {
	for watchingDbg {
		msg := <-models.DbgChan
		t.Log(msg)
	}
}
