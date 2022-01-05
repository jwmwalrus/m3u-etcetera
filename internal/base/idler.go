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

	// IdleStatusRequest A client's request is being processed
	IdleStatusRequest

	// IdleStatusSubscription A client subscription is active
	IdleStatusSubscription

	// IdleStatusDbOperations A DB-related operation is in progress
	IdleStatusDbOperations

	// IdleStatusFileOperations A file-related operation is in progress
	IdleStatusFileOperations
)

func (is IdleStatus) String() string {
	return []string{
		"idle",
		"engine-loop",
		"request",
		"subscription",
		"db-operations",
		"file-operations",
	}[is]
}

var (
	idleCancel context.CancelFunc
	idleCtx    context.Context

	forceExit       = false
	exitNow         = false
	idleStatusStack = []IdleStatus{IdleStatusIdle}
	doneEmmitted    = 0
	idleGotCalled   = false

	// InterruptSignal -
	InterruptSignal chan os.Signal
)

// DoTerminate forces immediate termination of the application
func DoTerminate(force bool) {
	exitNow = true
	if force || IsAppIdling() {
		forceExit = true
	}
}

// GetBusy registers a process as busy, to prevent idle timeout
func GetBusy(is IdleStatus) {
	if is == IdleStatusIdle {
		return
	}

	log.WithField("is", is).
		Debug("server got a lot busier")

	idleStatusStack = append(idleStatusStack, is)

	if idleGotCalled {
		idleCancel()
	}
}

// GetFree registers a process as less busy
func GetFree(is IdleStatus) {
	if is != IdleStatusIdle {
		log.WithField("is", is).
			Debug("server got a little less busy")

		for i := len(idleStatusStack) - 1; i >= 0; i-- {
			if is == idleStatusStack[i] {
				idleStatusStack[i] = idleStatusStack[len(idleStatusStack)-1]
				idleStatusStack = idleStatusStack[:len(idleStatusStack)-1]
				break
			}
		}
	}

	log.Debugf(
		"Topmost idle status is %v",
		idleStatusStack[len(idleStatusStack)-1],
	)

	if len(idleStatusStack) == 1 {
		if !idleGotCalled {
			idleCtx, idleCancel = context.WithCancel(context.Background())
			go Idle(idleCtx)
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
		if IsAppBusy() || idleGotCalled {
			log.Info("Server is busy or already idling, so cancelling request")
			<-ctx.Done()
			return
		}

		idleGotCalled = true
		log.Info("Entering Idle state")

		select {
		case <-time.After(time.Duration(ServerIdleTimeout) * time.Second):
			if IsAppBusy() {
				log.Info("Server is busy, so cancelling timeout")
				<-ctx.Done()
				return
			}
			break
		case <-ctx.Done():
			log.Info("idleCancel got called explicitly")
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
