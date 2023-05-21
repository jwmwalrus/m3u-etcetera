package base

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

const (
	// ServerIdleTimeout Amount of idle seconds before server exits.
	ServerIdleTimeout = 300
)

// IdleStatus defines the idle status type.
type IdleStatus int

const (
	// IdleStatusIdle Server is idle.
	IdleStatusIdle IdleStatus = iota

	// IdleStatusEngineLoop The engine loop is working.
	IdleStatusEngineLoop

	// IdleStatusRequest A client's request is being processed.
	IdleStatusRequest

	// IdleStatusSubscription A client subscription is active.
	IdleStatusSubscription

	// IdleStatusDbOperations A DB-related operation is in progress.
	IdleStatusDbOperations

	// IdleStatusFileOperations A file-related operation is in progress.
	IdleStatusFileOperations
)

func init() {
	signal.Notify(InterruptSignal, os.Interrupt, syscall.SIGTERM)
	idleStatusStack.s = []IdleStatus{IdleStatusIdle}
}

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

	forceExit       atomic.Bool
	doneEmmitted    atomic.Int32
	idleGotCalled   atomic.Bool
	idleStatusStack struct {
		s  []IdleStatus
		mu sync.RWMutex
	}

	// InterruptSignal -.
	InterruptSignal chan os.Signal = make(chan os.Signal, 1)
)

// DoTerminate forces immediate termination of the application.
func DoTerminate(force bool) {
	forceExit.Store(force || IsAppIdling())

	log.WithFields(log.Fields{
		"force":     force,
		"forceExit": forceExit.Load(),
	}).
		Debug("Immediate termination status")
}

// GetBusy registers a process as busy, to prevent idle timeout.
func GetBusy(is IdleStatus) {
	if is == IdleStatusIdle {
		return
	}

	log.WithField("is", is).
		Debug("server got a lot busier")

	idleStatusStack.mu.Lock()
	idleStatusStack.s = append(idleStatusStack.s, is)
	idleStatusStack.mu.Unlock()

	if idleGotCalled.Load() {
		idleCancel()
	}
}

// GetFree registers a process as less busy.
func GetFree(is IdleStatus) {
	entry := log.WithField("is", is)

	idleStatusStack.mu.Lock()
	defer idleStatusStack.mu.Unlock()

	if is != IdleStatusIdle {
		entry.Debug("server got a little less busy")

		for i := len(idleStatusStack.s) - 1; i >= 0; i-- {
			if is == idleStatusStack.s[i] {
				idleStatusStack.s[i] = idleStatusStack.s[len(idleStatusStack.s)-1]
				idleStatusStack.s = idleStatusStack.s[:len(idleStatusStack.s)-1]
				break
			}
		}
	}

	entry.Debugf(
		"Topmost idle status is %v",
		idleStatusStack.s[len(idleStatusStack.s)-1],
	)

	if len(idleStatusStack.s) == 1 {
		if !idleGotCalled.Load() {
			idleCtx, idleCancel = context.WithCancel(context.Background())
			go Idle(idleCtx)
		}
	}
}

// Idle exits the server if it has been idle for a while and no long-term
// processes are pending.
func Idle(ctx context.Context) {
	idleStatusStack.mu.RLock()
	entry := log.WithFields(log.Fields{
		"forceExit":            forceExit.Load(),
		"len(idleStatusStack)": len(idleStatusStack.s) - 1,
	})
	idleStatusStack.mu.RUnlock()

	entry.Info("Starting Idle checks")

	if !forceExit.Load() {
		if IsAppBusy() || idleGotCalled.Load() {
			entry.Info("Server is busy or already idling, so cancelling request")
			<-ctx.Done()
			return
		}

		idleGotCalled.Store(true)
		entry.Info("Entering Idle state")

		select {
		case <-time.After(time.Duration(ServerIdleTimeout) * time.Second):
			if IsAppBusy() {
				entry.Info("Server is busy, so cancelling timeout")
				<-ctx.Done()
				return
			}
			break
		case <-ctx.Done():
			entry.Info("idleCancel got called explicitly")
			idleGotCalled.Store(false)
			return
		}
	}

	if doneEmmitted.Load() > 0 {
		entry.WithField("doneEmmitted", doneEmmitted.Load()).
			Warn("ignoring further attempt at ctx.Done()")

		doneEmmitted.Add(1)
		return
	}

	doneEmmitted.Add(1)

	entry.Info("Server seems to have been Idle for a while, and that's gotta stop!")
	InterruptSignal <- os.Interrupt
}

// IsAppBusy returns true if some process has registered as busy.
func IsAppBusy() bool {
	idleStatusStack.mu.RLock()
	defer idleStatusStack.mu.RUnlock()

	return len(idleStatusStack.s) > 1
}

// IsAppBusyBy returns true if some process has registered as busy.
func IsAppBusyBy(is IdleStatus) bool {
	idleStatusStack.mu.RLock()
	defer idleStatusStack.mu.RUnlock()

	return slices.Contains(idleStatusStack.s, is)
}

// IsAppIdling returns true if the Idle method is active.
func IsAppIdling() bool {
	idleStatusStack.mu.RLock()
	defer idleStatusStack.mu.RUnlock()

	return idleGotCalled.Load() || len(idleStatusStack.s) == 1
}

// StartIdler -.
func StartIdler() {
	GetFree(IdleStatusIdle)
}
