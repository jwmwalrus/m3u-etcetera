package alive

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jwmwalrus/bnp/env"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	rtc "github.com/jwmwalrus/rtcycler"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type aliveSrv struct {
	turnOff  bool
	forceOff bool
	noWait   bool
}

func (a *aliveSrv) serve() (err error) {
	if isServerAlive() {
		err = &ServerAlreadyRunning{}
		if a.turnOff {
			err = a.stopServer()
		}
		return
	}

	err = &ServerNotRunning{}
	if !a.turnOff {
		err = a.startServer()
	}

	return
}

func (a *aliveSrv) startServer() (err error) {
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
		args = append(args, "--test")
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

	if a.noWait {
		aux, _ := err.(*ServerStarted)
		aux.Desc = "unconfirmed"
		return
	}

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

func (a *aliveSrv) stopServer() (err error) {
	log.WithFields(log.Fields{
		"forceOff": a.forceOff,
		"noWait":   a.noWait,
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
	res, err := c.Off(context.Background(), &m3uetcpb.OffRequest{Force: a.forceOff})
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

	if a.noWait {
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
