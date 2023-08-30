package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/internal/database"
	rtc "github.com/jwmwalrus/rtcycler"
	"gorm.io/gorm"
)

const (
	testData = "data/testing/fixtures"
)

type TestCase interface {
	FixturesDir() string
}

func SetupTest(t *testing.T, tc TestCase) *gorm.DB {
	rtc.ResetInstanceSuffix()
	rtc.SetTestMode()
	return openTestDatabase(tc.FixturesDir())
}

func TeardownTest(t *testing.T) {
	closeTestDatabase()

	if _, err := os.Stat(database.Path()); err == nil {
		if err = os.Remove(database.Path()); err != nil {
			panic(err)
		}
	}
}

func getTestDataDir() string {
	path := filepath.Join(".", testData)
	for {
		if _, err := os.Stat(path); err == nil {
			path, _ = filepath.Abs(path)
			return path
		}
		full, _ := filepath.Abs(path)
		if full == filepath.Join(string(filepath.Separator), testData) {
			break
		}
		path = filepath.Join("..", path)
	}
	return ""
}
