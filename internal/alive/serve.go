package alive

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/jwmwalrus/onerror"
	rtc "github.com/jwmwalrus/rtcycler"
	log "github.com/sirupsen/logrus"
)

const (
	// ServerCheckInterval Amount of seconds between checks.
	ServerCheckInterval = 180
)

var (
	serverAliveFilename = "server-alive"

	lastStatus error

	// LastCheck UTC timestamp for last check.
	LastCheck atomic.Int64
)

// CheckServerStatus if ServerCheckInterval is up, starts the server. Otherwise,
// it returns the status since the last check.
func CheckServerStatus() error {
	if lastStatus == nil || (time.Now().Unix()-LastCheck.Load() > ServerCheckInterval) {
		lastStatus = Serve()
	}

	return lastStatus
}

// Serve starts or stops the server, according to the given options.
func Serve(o ...Option) error {
	a := &aliveSrv{}
	for i := range o {
		o[i](a)
	}

	if a.forceOff {
		a.turnOff = true
	}

	return a.serve()
}

// readServerAlive reads the server alive flag file.
func readServerAlive() {
	log.Debug("Reading server status from file")

	// Last alive check for server
	info, err := os.Stat(filepath.Join(rtc.DataDir(), serverAliveFilename))
	if !os.IsNotExist(err) {
		LastCheck.Store(info.ModTime().Unix())
	}
}

// writeServerAliveFile updates the server alive flag file.
func writeServerAliveFile() {
	log.Debug("Writing server alive file")

	f, err := os.OpenFile(
		filepath.Join(rtc.DataDir(), serverAliveFilename),
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()

	_, err = f.WriteString("1")
	onerror.Log(err)
}
