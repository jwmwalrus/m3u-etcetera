package discover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jwmwalrus/bnp/env"
	log "github.com/sirupsen/logrus"
)

// Execute invokes m3uetc-discover for the given location
func Execute(location string) (*Info, error) {
	log.WithField("location", location).Info("Executing discover")

	app := "m3u-etcetera"
	bin := "m3uetc-discover"
	path := env.FindLibExec(bin, app)
	if path == "" {
		bin = "discover"
		path = env.FindLibExec(bin, app)
		if path == "" {
			return nil, fmt.Errorf("failed to find the discover binary")
		}
	}
	cmd := exec.Command(path, "-l", location)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		str := strings.Replace(errb.String(), "\n", " | ", -1)
		str = strings.Replace(str, "|  |", "|", -1)
		err = fmt.Errorf(strings.TrimSpace(str))
	}

	is := &Info{}
	if len(outb.String()) > 0 {
		err := json.Unmarshal(outb.Bytes(), is)
		if err != nil {
			return nil, err
		}
	}

	return is, err
}
