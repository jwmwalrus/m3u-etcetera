package base

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/jwmwalrus/bnp/env"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/pkg/config"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	configFilename = "config.json"
	lockFilename   = configFilename + ".lock"

	// AppDirName Application's directory name
	AppDirName = "m3u-etcetera"

	// AppName Application name
	AppName = "M3U Etcétera"
)

var (
	configFile string
	lockFile   string
	logFile    *lumberjack.Logger

	logFilename = "app.log"

	serverIsBusy = 0

	// OS Operating system's name
	OS string

	// CacheDir Home directory for cache
	CacheDir string

	// ConfigDir Home directory for config
	ConfigDir string

	// DataDir Home directory for data
	DataDir string

	// RuntimeDir Run (volatile) directory
	RuntimeDir string

	// FlagDry Dry run
	FlagDry bool

	// FlagTestingMode Start in testing mode
	FlagTestingMode bool

	// FlagVerbose Logger severity on
	FlagVerbose bool

	// FlagSeverity Logger severity level
	FlagSeverity string

	// FlagEchoLogging Echo logs to stderr
	FlagEchoLogging bool

	// Conf global configuration
	Conf config.Config
)

// Idle exits the server if it has been idle for a while and no long-term processes are pending
func Idle(force bool) {
	if force || serverIsBusy <= 0 {
		log.Info("Server has been idle for a while, and that's gotta stop!")
		Unload()
		fmt.Printf("Bye %v from %v\n", OS, filepath.Base(os.Args[0]))
		os.Exit(0)
	}
}

// Load Loads application's configuration
func Load() (args []string) {
	args = env.ParseArgs(
		logFile,
		&FlagEchoLogging,
		&FlagVerbose,
		&FlagSeverity,
	)

	configFile = filepath.Join(ConfigDir, configFilename)
	lockFile = filepath.Join(RuntimeDir, lockFilename)

	err := env.SetDirs(
		CacheDir,
		ConfigDir,
		DataDir,
		RuntimeDir,
	)
	onerror.Panic(err)

	err = Conf.Load(configFile, lockFile)
	onerror.Panic(err)

	return
}

// Unload Cleans up server before exit
func Unload() {}

func init() {
	OS = runtime.GOOS

	// XDG-related
	DataDir = filepath.Join(xdg.DataHome, AppDirName)
	ConfigDir = filepath.Join(xdg.ConfigHome, AppDirName)
	CacheDir = filepath.Join(xdg.CacheHome, AppDirName)
	RuntimeDir = filepath.Join(xdg.RuntimeDir, AppDirName)

	// Define global flags
	getopt.FlagLong(&FlagDry, "dry", 'n', "Dry run")
	getopt.FlagLong(&FlagTestingMode, "testing", 't', "Start in testing mode")
	getopt.FlagLong(&FlagVerbose, "verbose", 'v', "Bump logging severity")
	getopt.FlagLong(&FlagSeverity, "severity", 0, "Logging severity")
	getopt.FlagLong(&FlagEchoLogging, "echo-logging", 'e', "Echo logs to stderr")

	// log-related
	logFilename = filepath.Base(os.Args[0]) + ".log"
	logFilePath := filepath.Join(DataDir, logFilename)
	logFile = &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     30,    //days
		Compress:   false, // disabled by default
	}
	log.SetOutput(logFile)
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.ErrorLevel)
}
