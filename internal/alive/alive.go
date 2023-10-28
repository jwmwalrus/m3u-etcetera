package alive

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jwmwalrus/bnp/env"
	"github.com/jwmwalrus/gear-pieces/middleware"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	rtc "github.com/jwmwalrus/rtcycler"
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
	slog.Debug("Starting server")

	var dir, full string
	bin := "m3uetc-server"

	if rtc.FlagTestMode() {
		if full = env.FindExec(bin); full == "" {
			err = fmt.Errorf("failed to find binary `%s` for server", bin)
			return
		}
	} else {
		slog.Debug("Using $PATH to find binary", "PATH", os.Getenv("PATH"))
		full, err = exec.LookPath(bin)
		if err != nil {
			slog.With(
				"bin", bin,
				"error", err,
			).Debug("Error finding binary from path")
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
		slog.With(
			"cmd", cmd,
			"error", err,
		).Debug("Command exited with an error status")
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
	slog.With(
		"forceOff", a.forceOff,
		"noWait", a.noWait,
	).Debug("Stopping server")

	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts...)
	if err != nil {
		slog.With(
			"authority", auth,
			"error", err,
		).Debug("Error while dialing server")
		return
	}
	defer cc.Close()

	slog.Debug("Dialing was successful", "authority", auth)

	c := m3uetcpb.NewRootSvcClient(cc)
	res, err := c.Off(context.Background(), &m3uetcpb.OffRequest{Force: a.forceOff})
	if err != nil {
		slog.With(
			"authority", auth,
			"error", err,
		).Debug("Error requesting server off")
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
	slog.Debug("Checking if server is alive")

	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts...)
	if err != nil {
		s := status.Convert(err)
		if s.Code() != codes.Unavailable {
			slog.With(
				"authority", auth,
				"error", err,
			).Info("Failed to dial server")
		}
		return false
	}
	defer cc.Close()

	slog.With("authority", auth).Debug("Dialing was successful")

	c := m3uetcpb.NewRootSvcClient(cc)
	_, err = c.Status(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		slog.Debug("Failed to obtain server status", "error", err)
		return false
	}

	return true
}
