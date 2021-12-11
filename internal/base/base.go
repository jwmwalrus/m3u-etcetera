package base

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
)

var (
	serverIsBusy = 0

	// OS Operating system's name
	OS string
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

// Unload Cleans up server before exit
func Unload() {}

func init() {
	OS = runtime.GOOS
}
