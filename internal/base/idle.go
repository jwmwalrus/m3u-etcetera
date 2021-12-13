package base

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	// ServerIdleTimeout Amount of idle seconds before server exits
	ServerIdleTimeout = 300
)

var (
	serverIsBusy = 0
)

// GetBusy registers a process as busy, to prevent idle timeout
func GetBusy() {
	log.Info("server got a lot busier")
	serverIsBusy++
}

// GetFree registers a process as less busy
func GetFree() {
	log.Info("server got a little less busy")
	serverIsBusy--
}

// Idle exits the server if it has been idle for a while and no long-term processes are pending
func Idle(force bool) {
	log.WithFields(log.Fields{
		"force":        force,
		"serverIsBusy": serverIsBusy,
	}).
		Info("Server has been idle for a while, and that's gotta stop!")

	if force || serverIsBusy <= 0 {
		log.Info("Server seems to have been idle for a while, and that's gotta stop!")
		Unload()
		fmt.Printf("Bye %v from %v\n", OS, filepath.Base(os.Args[0]))
		os.Exit(0)
	}
}

// IsItBusy returns true if some process has registered as busy
func IsItBusy() bool {
	return serverIsBusy > 0
}
