package base

import (
	"os"
	"path/filepath"

	"github.com/jwmwalrus/bnp/env"
	"github.com/jwmwalrus/bnp/onerror"
	log "github.com/sirupsen/logrus"
)

// Unloader defines a method to be called when unloading the application
type Unloader struct {
	Description string
	Callback    UnloaderCallback
}

// UnloaderCallback defines the signature of the method to be called when unloading the application
type UnloaderCallback func() error

var unloadRegistry []Unloader

// Load Loads application's configuration
func Load(noParseArgs ...bool) (args []string) {
	noParse := false
	if len(noParseArgs) > 0 {
		noParse = noParseArgs[0]
	}

	if !noParse {
		args = parseArgs()
	}

	configFile = filepath.Join(ConfigDir, configFilename)
	lockFile = filepath.Join(RuntimeDir, lockFilename)

	err := env.SetDirs(
		CacheDir,
		ConfigDir,
		DataDir,
		RuntimeDir,
	)
	onerror.Panic(err)

	if _, err = os.Stat(CoversDir); os.IsNotExist(err) {
		err = os.MkdirAll(CoversDir, 0755)
		if err != nil {
			return
		}
	}

	err = Conf.Load(configFile, lockFile)
	onerror.Panic(err)
	return
}

// RegisterUnloader registers an Unloader, to be invoked before stopping the app
func RegisterUnloader(u Unloader) {
	unloadRegistry = append(unloadRegistry, u)
}

// Unload Cleans up server before exit
func Unload() {
	log.Info("Unloading application")

	if Conf.FirstRun {
		Conf.FirstRun = false

		err := Conf.Save(configFile, lockFile)
		onerror.Log(err)
	}

	for _, u := range unloadRegistry {
		log.Infof("Calling unloader: %v", u.Description)
		err := u.Callback()
		onerror.Log(err)
	}
}
