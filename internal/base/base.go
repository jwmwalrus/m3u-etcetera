package base

import (
	"path/filepath"

	"github.com/jwmwalrus/m3u-etcetera/internal/config"
	rtc "github.com/jwmwalrus/rtcycler"
)

const (
	// AppDirName Application's directory name.
	AppDirName = "m3u-etcetera"

	// AppName Application name.
	AppName = "M3U Etcétera"

	// DatabaseFilename defines the name of the music database.
	DatabaseFilename = "music.db"

	// ServerWaitTimeout Maximum amount of seconds the server will wait
	// for an event to sync up.
	ServerWaitTimeout = 30

	// ClientWaitTimeout Maximum amount of seconds a client should wait
	// for an event to sync up.
	ClientWaitTimeout = 30

	// PlaybackPlayedThreshold -.
	PlaybackPlayedThreshold = 30

	// CoversDirname directory where covers are stored.
	CoversDirname = "covers"
)

var (
	// Conf global configuration.
	Conf config.Config
)

// CoversDir returns the covers directory.
func CoversDir() string {
	return filepath.Join(rtc.DataDir(), CoversDirname)
}
