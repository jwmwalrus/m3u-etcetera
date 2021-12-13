package alive

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
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
	LastCheck int64
)

// CheckServerStatus If ServerCheckInterval is up, starts the server
func CheckServerStatus() (err error) {
	if time.Now().Unix()-LastCheck > ServerCheckInterval {
		err = Serve(false)
	}

	return
}

// Serve starts or stops the server
func Serve(turnOff bool) (err error) {
	if isServerAlive() {
		err = &ServerAlreadyRunning{}
		if turnOff {
			err = stopServer()
		}
		return
	}

	err = &ServerNotRunning{}
	if !turnOff {
		err = startServer()
	}

	return
}

func findBinary(bin string) (path string, err error) {
	path, err = filepath.Abs(filepath.Join(".", bin))
	rel := bin
	if err != nil {
		return
	}

	for {
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			var s os.FileInfo
			if s, err = os.Stat(".git"); os.IsExist(err) && s.IsDir() {
				err = errors.New("Reached .git without finding binary")
				return
			}
			rel = filepath.Join("..", rel)
			path, err = filepath.Abs(rel)
			if err != nil {
				return
			}
			if path == string(filepath.Separator) || path == "" {
				err = errors.New("Reached root without finding binary")
				return
			}
		} else {
			break
		}
	}
	return
}

func isServerAlive() bool {
	opts := grpc.WithInsecure()

	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts)
	if err != nil {
		s := status.Convert(err)
		if s.Code() != codes.Unavailable {
			log.Info(err)
		}
		return false
	}
	defer cc.Close()

	c := m3uetcpb.NewRootSvcClient(cc)
	res, err := c.Status(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Info(err)
		return false
	}

	return res.GetAlive()
}

func readServerAlive() {
	// Last alive check for server
	if info, err := os.Stat(filepath.Join(base.DataDir, serverAliveFilename)); !os.IsNotExist(err) {
		LastCheck = info.ModTime().Unix()
	}
}

func startServer() (err error) {
	var dir, full string
	bin := "m3uetc-server"

	if !base.FlagTestingMode {
		if full, err = findBinary(bin); err != nil {
			return
		}
	} else {
		full, err = exec.LookPath("m3uetc-server")
		if err == nil {
			if full, err = findBinary(bin); err != nil {
				return
			}
		}
	}

	dir = filepath.Dir(full)

	args := []string{}
	if base.FlagTestingMode {
		args = append(args, "--testing")
	}

	cmd := exec.Command(filepath.Join(dir, bin), args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err = cmd.Start()
	if err != nil {
		log.Error(err)
		return
	}

	pid := cmd.Process.Pid

	err = &ServerStarted{PID: pid}

	alive := false
	for i := 0; i < 12; i++ {
		if alive = isServerAlive(); !alive {
			time.Sleep(5 * time.Second)
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

func stopServer() (err error) {
	opts := grpc.WithInsecure()

	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts)
	if err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()

	c := m3uetcpb.NewRootSvcClient(cc)
	res, err := c.Off(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Error(err)
		return
	}

	err = &ServerStopped{}
	if !res.GetGoingOff() {
		aux, _ := err.(*ServerStopped)
		aux.Desc = res.GetReason()
	}

	alive := true
	for i := 0; i < 12; i++ {
		if alive = isServerAlive(); alive {
			time.Sleep(5 * time.Second)
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
	f, err := os.OpenFile(filepath.Join(base.DataDir, serverAliveFilename), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()

	_, err = f.WriteString("1")

	return
}

func init() {
	readServerAlive()
}
