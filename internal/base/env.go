package base

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/jwmwalrus/bnp/ing2"
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

	// ServerWaitTimeout Maximum amount of seconds the server will wait
	// for an event to sync up
	ServerWaitTimeout = 30

	// ClientWaitTimeout Maximum amount of seconds a client should wait
	// for an event to sync up
	ClientWaitTimeout = 30
)

var (
	configFile string
	lockFile   string
	logFile    *lumberjack.Logger

	logFilename = "app.log"

	// OS Operating system's name
	OS string

	// AppInstance Application's instance
	AppInstance string

	// CacheDir Home directory for cache
	CacheDir string

	// ConfigDir Home directory for config
	ConfigDir string

	// DataDir Home directory for data
	DataDir string

	// RuntimeDir Run (volatile) directory
	RuntimeDir string

	// CoversDir directory where covers are stored
	CoversDir string

	flagHelp bool

	// FlagDry Dry run
	FlagDry bool

	// FlagDebugMode Start in debug mode
	FlagDebugMode bool

	// FlagTestingMode Start in testing mode
	FlagTestingMode bool

	// FlagVerbose Logger severity on
	FlagVerbose bool

	// FlagSeverity Logger severity level
	FlagSeverity string

	// FlagEchoLogging Echo logs to stderr
	FlagEchoLogging bool

	// InstanceSuffix suffix used for the running instance
	InstanceSuffix string

	// Conf global configuration
	Conf config.Config
)

// SetTestingMode -
func SetTestingMode() {
	FlagTestingMode = true
}

// UnsetTestingMode -
func UnsetTestingMode() {
	FlagTestingMode = false
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
		if FlagDebugMode {
			FlagSeverity = "debug"
		} else if FlagTestingMode {
			FlagSeverity = "debug"
		} else if FlagVerbose {
			FlagSeverity = "info"
		} else {
			FlagSeverity = "error"
		}
	} else {
		if _, err := log.ParseLevel(givenSeverity); err != nil {
			fmt.Printf("Unsupported severity level: %v\n", givenSeverity)
			FlagSeverity = "error"
		} else {
			FlagSeverity = givenSeverity
		}
	}

	level, _ := log.ParseLevel(FlagSeverity)
	log.SetLevel(level)
	// log.SetReportCaller(true)
	log.SetReportCaller(FlagSeverity == "debug")

	return
}

func init() {
	OS = runtime.GOOS

	InstanceSuffix = ing2.GetRandomString(8)

	// XDG-related
	DataDir = filepath.Join(xdg.DataHome, AppDirName)
	ConfigDir = filepath.Join(xdg.ConfigHome, AppDirName)
	CacheDir = filepath.Join(xdg.CacheHome, AppDirName)
	RuntimeDir = filepath.Join(xdg.RuntimeDir, AppDirName)
	CoversDir = filepath.Join(DataDir, "covers")

	// Define global flags
	getopt.FlagLong(&flagHelp, "help", 'h', "Display this help")
	getopt.FlagLong(&FlagDry, "dry", 'n', "Dry run")
	getopt.FlagLong(&FlagDebugMode, "debug", 0, "Start in debug mode")
	getopt.FlagLong(&FlagTestingMode, "testing", 0, "Start in testing mode")
	getopt.FlagLong(&FlagVerbose, "verbose", 'v', "Bump logging severity")
	getopt.FlagLong(&FlagSeverity, "severity", 0, "Logging severity (debug|info|warn|error|fatal)")
	getopt.FlagLong(&FlagEchoLogging, "echo-logging", 'e', "Echo logs to stderr")

	// log-related
	AppInstance = filepath.Base(os.Args[0])
	logFilename = AppInstance + ".log"
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
