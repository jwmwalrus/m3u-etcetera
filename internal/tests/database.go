package tests

import (
	"path/filepath"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jwmwalrus/m3u-etcetera/internal/database"
	rtc "github.com/jwmwalrus/rtcycler"
	"gorm.io/gorm"
)

var (
	fixtures   *testfixtures.Loader
	dbUnloader *rtc.Unloader
)

func closeTestDatabase() {
	dbUnloader.Callback()
}

func openTestDatabase(fixturesDir string) (db *gorm.DB) {
	var err error

	dbUnloader = database.Open()
	db = database.Instance()

	dir := getTestDataDir()
	if dir == "" {
		panic("failed to find test data directory")
	}

	sqlDB, _ := db.DB()
	fixtures, err = testfixtures.New(
		testfixtures.Database(sqlDB),
		testfixtures.Dialect("sqlite"),
		testfixtures.Paths(filepath.Join(dir, fixturesDir)),
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
