package base

import (
	"io"
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

	// DatabaseFilename defines the name of the music database
	DatabaseFilename = "music.db"

	// ServerWaitTimeout Maximum amount of seconds the server will wait for an event to sync up
	ServerWaitTimeout = 60
)

var (
	configFile     string
	lockFile       string
	logFile        *lumberjack.Logger
	unloadRegistry []Unloader

	logFilename = "app.log"

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

	flagHelp bool

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

// Unloader defines a method to be called when unloading the application
type Unloader struct {
	Description string
	Callback    UnloaderCallback
}

// UnloaderCallback defines the signature of the method to be called when unloading the application
type UnloaderCallback func() error

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

		// NOTE: ignore error
		Conf.Save(configFile, lockFile)
	}

	for _, u := range unloadRegistry {
		log.Infof("Calling finalizer: %v", u.Description)
		err := u.Callback()
		onerror.Log(err)
	}
}

func parseArgs() (args []string) {
	getopt.Parse()
	args = getopt.Args()
	arg0 := []string{os.Args[0]}
	args = append(arg0, args...)

	if flagHelp {
		getopt.Usage()
		os.Exit(0)
	}

	resolveSeverity()

	if FlagEchoLogging {
		mw := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(mw)
	}

	return
}

func resolveSeverity() {
	givenSeverity := FlagSeverity

	if givenSeverity == "" {
		if FlagVerbose {
			FlagSeverity = "info"
		} else {
			FlagSeverity = "error"
		}
	} else {
		if _, err := log.ParseLevel(givenSeverity); err != nil {
			FlagSeverity = "error"
		} else {
			FlagSeverity = givenSeverity
		}
	}

	level, _ := log.ParseLevel(FlagSeverity)
	log.SetLevel(level)
	log.SetReportCaller(FlagSeverity == "debug")

	return
}

func init() {
	OS = runtime.GOOS

	// XDG-related
	DataDir = filepath.Join(xdg.DataHome, AppDirName)
	ConfigDir = filepath.Join(xdg.ConfigHome, AppDirName)
	CacheDir = filepath.Join(xdg.CacheHome, AppDirName)
	RuntimeDir = filepath.Join(xdg.RuntimeDir, AppDirName)

	// Define global flags
	getopt.FlagLong(&flagHelp, "help", 'h', "Display this help")
	getopt.FlagLong(&FlagDry, "dry", 'n', "Dry run")
	getopt.FlagLong(&FlagTestingMode, "testing", 0, "Start in testing mode")
	getopt.FlagLong(&FlagVerbose, "verbose", 'v', "Bump logging severity")
	getopt.FlagLong(&FlagSeverity, "severity", 0, "Logging severity (trace|debug|info|warn|error|fatal|panic)")
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
