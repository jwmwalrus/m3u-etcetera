package api

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

const (
	testDataDir = "../data/testing"
)

var (
	conn       *gorm.DB
	fixtures   *testfixtures.Loader
	dbUnloader *base.Unloader
)

type testCase struct {
	name        string
	fixturesDir string
	req         proto.Message
	res         proto.Message
	err         error
}

func closeTestDatabase() {
	dbUnloader.Callback()
}

func openTestDatabase(fixturesDir string) (db *gorm.DB) {
	var err error

	dbUnloader = database.Open()
	db = database.Instance()

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
	base.FlagTestingMode = true
	conn = openTestDatabase(tc.fixturesDir)
	return
}

func teardownTest(t *testing.T) {
	closeTestDatabase()

	if _, err := os.Stat(database.DSN()); !os.IsNotExist(err) {
		if err = os.Remove(database.DSN()); err != nil {
			panic(err)
		}
	}
}
