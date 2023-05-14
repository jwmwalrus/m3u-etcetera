package alive

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/jwmwalrus/bnp/env"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	rtc "github.com/jwmwalrus/rtcycler"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// ServerCheckInterval Amount of seconds between checks
	ServerCheckInterval = 180
)

var (
	serverAliveFilename = "server-alive"

	// LastCheck UTC timestamp for last check
	LastCheck atomic.Int64

	lastStatus error
)

func init() {
	readServerAlive()
}

// CheckServerStatus If ServerCheckInterval is up, starts the server
func CheckServerStatus() error {
	if lastStatus == nil || (time.Now().Unix()-LastCheck.Load() > ServerCheckInterval) {
		lastStatus = Serve()
	}

	return lastStatus
}

// Serve starts or stops the server
func Serve(o ...ServeOptions) (err error) {
	options := ServeOptions{}
	if len(o) > 0 {
		options = o[0]
	}
	turnOff := options.TurnOff
	force := options.Force
	if force {
		turnOff = true
	}
	noWait := options.NoWait

	if isServerAlive() {
		err = &ServerAlreadyRunning{}
		if turnOff {
			err = stopServer(force, noWait)
		}
		return
	}

	err = &ServerNotRunning{}
	if !turnOff {
		err = startServer()
	}

	return
}

func isServerAlive() bool {
	log.Debug("Checking if server is alive")

	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts...)
	if err != nil {
		s := status.Convert(err)
		if s.Code() != codes.Unavailable {
			log.WithError(err).
				WithField("authority", auth).
				Info("Failed to dial server")
		}
		return false
	}
	defer cc.Close()

	log.WithField("authority", auth).Debug("Dialing was successful")

	c := m3uetcpb.NewRootSvcClient(cc)
	res, err := c.Status(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Debugf("Failed to obtain server status: %v", err)
		return false
	}

	return res.GetAlive()
}

func readServerAlive() {
	log.Debug("Reading server status from file")

	// Last alive check for server
	info, err := os.Stat(filepath.Join(rtc.DataDir(), serverAliveFilename))
	if !os.IsNotExist(err) {
		LastCheck.Store(info.ModTime().Unix())
	}
}

func startServer() (err error) {
	log.Debug("Starting server")

	var dir, full string
	bin := "m3uetc-server"

	if rtc.FlagTestMode() {
		if full = env.FindExec(bin); full == "" {
			err = fmt.Errorf("failed to find binary `%s` for server", bin)
			return
		}
	} else {
		log.WithField("PATH", os.Getenv("PATH")).
			Debug("Using $PATH to find binary")
		full, err = exec.LookPath(bin)
		if err != nil {
			log.WithError(err).
				WithField("bin", bin).
				Debug("Error finding binary from path")
			if full = env.FindExec(bin); full == "" {
				err = fmt.Errorf("failed to find binary `%s` for server", bin)
				return
			}
		}
	}

	dir = filepath.Dir(full)

	args := []string{}
	if rtc.FlagTestMode() {
		args = append(args, "--testing")
	}

	cmd := exec.Command(filepath.Join(dir, bin), args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err = cmd.Start()
	if err != nil {
		log.WithError(err).
			WithField("cmd", cmd.String()).
			Debug("Error during command execution")
		return
	}

	pid := cmd.Process.Pid

	err = &ServerStarted{PID: pid}

	alive := false
	for i := 0; i < base.ClientWaitTimeout; i++ {
		if alive = isServerAlive(); !alive {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	if !alive {
		aux, _ := err.(*ServerStarted)
		aux.Desc = "server might not be ready yet"
	}

	writeServerAliveFile()
	readServerAlive()

	return
}

func stopServer(force, noWait bool) (err error) {
	log.WithFields(log.Fields{
		"force":  force,
		"noWait": noWait,
	}).Debug("Stopping server")

	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts...)
	if err != nil {
		log.WithError(err).
			WithField("authority", auth).
			Debug("Error while dialing server")
		return
	}
	defer cc.Close()

	log.WithField("authority", auth).Debug("Dialing was successful")

	c := m3uetcpb.NewRootSvcClient(cc)
	res, err := c.Off(context.Background(), &m3uetcpb.OffRequest{Force: force})
	if err != nil {
		log.WithError(err).
			WithField("authority", auth).
			Debug("Error requesting server off")
		return
	}

	err = &ServerStopped{}
	if !res.GoingOff {
		aux, _ := err.(*ServerStopped)
		aux.Desc = res.Reason
	}

	if noWait {
		aux, _ := err.(*ServerStopped)
		aux.Desc = "unconfirmed"
		return
	}

	alive := true
	for i := 0; i < base.ClientWaitTimeout; i++ {
		if alive = isServerAlive(); alive {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	if alive {
		aux, _ := err.(*ServerStopped)
		aux.Desc = "server might still be running"
	}

	return
}

// writeServerAliveFile Updates the server alive flag file
func writeServerAliveFile() {
	log.Debug("Writting server alive file")

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
