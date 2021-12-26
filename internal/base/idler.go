package base

import (
	"context"
	"os"
	"time"

	"github.com/jwmwalrus/bnp/slice"
	log "github.com/sirupsen/logrus"
)

const (
	// ServerIdleTimeout Amount of idle seconds before server exits
	ServerIdleTimeout = 300
)

// IdleStatus defines the idle status type
type IdleStatus int

const (
	// IdleStatusIdle Server is idle
	IdleStatusIdle IdleStatus = iota

	// IdleStatusEngineLoop The engine loop is working
	IdleStatusEngineLoop

	// IdleStatusSubscription A client subscription is active
	IdleStatusSubscription

	// IdleStatusDbOperations A DB-related operation is in progress
	IdleStatusDbOperations

	// IdleStatusFileOperations A file-related operation is in progress
	IdleStatusFileOperations
)

func (is IdleStatus) String() string {
	return []string{"idle", "engine-loop", "db-operations", "file-operations"}[is]
}

var (
	forceExit       = false
	idleStatusStack = []IdleStatus{IdleStatusIdle}
	doneEmmitted    = 0
	idleGotCalled   = false

	// InteruptSignal -
	InterruptSignal chan os.Signal
)

// GetBusy registers a process as busy, to prevent idle timeout
func GetBusy(is IdleStatus) {
	if is == IdleStatusIdle {
		return
	}

	log.Info("server got a lot busier")
	idleStatusStack = append(idleStatusStack, is)
}

// GetFree registers a process as less busy
func GetFree(is IdleStatus) {
	if is == IdleStatusIdle {
		return
	}

	log.WithField("is", is).
		Info("server got a little less busy")

	for i := len(idleStatusStack) - 1; i >= 0; i-- {
		if is == idleStatusStack[i] {
			idleStatusStack[i] = idleStatusStack[len(idleStatusStack)-1]
			idleStatusStack = idleStatusStack[:len(idleStatusStack)-1]
			break
		}
	}
}

// Idle exits the server if it has been idle for a while and no long-term processes are pending
func Idle(ctx context.Context) {
	log.WithFields(log.Fields{
		"forceExit":            forceExit,
		"len(idleStatusStack)": len(idleStatusStack) - 1,
	}).
		Info("Stating Idle checks")

	if !forceExit {
		if len(idleStatusStack) > 1 || idleGotCalled {
			log.Info("Server is busy or already idling, so cancelling request")
			<-ctx.Done()
			return
		}

		idleGotCalled = true
		log.Info("Entering Idle state")

		select {
		case <-time.After(time.Duration(ServerIdleTimeout) * time.Second):
			break
		case <-ctx.Done():
			idleGotCalled = false
			return
		}
	}

	if doneEmmitted > 0 {
		log.WithField("doneEmmitted", doneEmmitted).
			Warn("ignoring further attempt at ctx.Done()")

		doneEmmitted++
		return
	}

	doneEmmitted++

	log.Info("Server seems to have been Idle for a while, and that's gotta stop!")
	InterruptSignal <- os.Interrupt
}

// IsAppBusy returns true if some process has registered as busy
func IsAppBusy() bool {
	return len(idleStatusStack) > 1
}

// IsAppBusyBy returns true if some process has registered as busy
func IsAppBusyBy(is IdleStatus) bool {
	return slice.Contains(idleStatusStack, is)
}

// IsAppIdling returns true if the Idle method is active
func IsAppIdling() bool {
	return idleGotCalled //&& len(idleStatusStack) == 1
}

// DoTerminate forces immediate termination of the application
func DoTerminate() {
	forceExit = true
}
