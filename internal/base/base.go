package base

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/adrg/xdg"
	"github.com/jwmwalrus/bnp/env"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/bnp/urlstr"
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

	// SupportedFileExtensionMP3 supported mp3
	SupportedFileExtensionMP3 = ".mp3"

	// SupportedFileExtensionM4A supported m4a
	SupportedFileExtensionM4A = ".m4a"

	// SupportedFileExtensionOGG supported ogg
	SupportedFileExtensionOGG = ".ogg"

	// SupportedFileExtensionFLAC supported flac
	SupportedFileExtensionFLAC = ".flac"
)

var (
	configFile     string
	lockFile       string
	logFile        *lumberjack.Logger
	unloadRegistry []Unloader
	randomSeed     *rand.Rand
	dbTmpSuffix    string

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

	// Conf global configuration
	Conf config.Config

	// SupportedFileExtensions supported file extensons
	SupportedFileExtensions = []string{
		SupportedFileExtensionMP3,
		SupportedFileExtensionM4A,
		SupportedFileExtensionOGG,
		SupportedFileExtensionFLAC,
	}

	// IgnoredFileExtensions supported file extensons
	IgnoredFileExtensions = []string{
		".bmp",
		".db",
		".gif",
		".jpeg",
		".jpg",
		".png",
	}
)

// Unloader defines a method to be called when unloading the application
type Unloader struct {
	Description string
	Callback    UnloaderCallback
}

// UnloaderCallback defines the signature of the method to be called when unloading the application
type UnloaderCallback func() error

// GetDatabaseDir returns the database directory
func GetDatabaseDir() string {
	if !FlagTestingMode {
		return DataDir
	}
	return os.TempDir()
}

// GetDatabaseFilename returns the database filename
func GetDatabaseFilename() string {
	if !FlagTestingMode {
		return DatabaseFilename
	}
	return "m3uetc-test-music-" + dbTmpSuffix + ".db"
}

func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[randomSeed.Intn(len(charset))]
	}
	return string(b)
}

// GetDatabasePath returns the full path to the database {
func GetDatabasePath() string {
	return filepath.Join(GetDatabaseDir(), GetDatabaseFilename())
}

// IsSupportedURL returns true if the path is supported
func IsSupportedURL(s string) bool {
	path, err := urlstr.URLToPath(s)
	if err != nil {
		return false
	}

	return IsSupportedFile(path)
}

// IsSupportedFile returns true if the path is supported
func IsSupportedFile(path string) bool {
	return slice.Contains(SupportedFileExtensions, filepath.Ext(path))
}

// IsIgnoredFile returns true if the path should be ignored
func IsIgnoredFile(path string) bool {
	return slice.Contains(IgnoredFileExtensions, filepath.Ext(path))
}

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

// SetTestingMode -
func SetTestingMode() {
	FlagTestingMode = true
}

// UnsetTestingMode -
func UnsetTestingMode() {
	FlagTestingMode = false
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
		log.Infof("Calling unloader: %v", u.Description)
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
	log.SetReportCaller(FlagSeverity == "debug")

	return
}

func init() {
	OS = runtime.GOOS

	randomSeed = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	dbTmpSuffix = randomString(8)

	// XDG-related
	DataDir = filepath.Join(xdg.DataHome, AppDirName)
	ConfigDir = filepath.Join(xdg.ConfigHome, AppDirName)
	CacheDir = filepath.Join(xdg.CacheHome, AppDirName)
	RuntimeDir = filepath.Join(xdg.RuntimeDir, AppDirName)

	// Define global flags
	getopt.FlagLong(&flagHelp, "help", 'h', "Display this help")
	getopt.FlagLong(&FlagDry, "dry", 'n', "Dry run")
	getopt.FlagLong(&FlagDebugMode, "debug", 0, "Start in debug mode")
	getopt.FlagLong(&FlagTestingMode, "testing", 0, "Start in testing mode")
	getopt.FlagLong(&FlagVerbose, "verbose", 'v', "Bump logging severity")
	getopt.FlagLong(&FlagSeverity, "severity", 0, "Logging severity (debug|info|warn|error|fatal)")
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
