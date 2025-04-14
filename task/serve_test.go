package task

import (
	"reflect"
	"testing"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	rtc "github.com/jwmwalrus/rtcycler"
	"github.com/urfave/cli/v3"
)

func TestServe(t *testing.T) {
	table := []struct {
		name     string
		command  []string
		expected error
		setup    func()
		teardown func()
	}{
		{
			"Start non-running",
			[]string{"", "serve"},
			&alive.ServerStarted{},
			func() { alive.Serve(alive.WithForceOff()) },
			func() { alive.Serve(alive.WithForceOff()) },
		},
		{
			"Start already running",
			[]string{"", "serve"},
			&alive.ServerAlreadyRunning{},
			func() { alive.Serve() },
			func() { alive.Serve(alive.WithForceOff()) },
		},
		{
			"Stop non-running",
			[]string{"", "serve", "off"},
			&alive.ServerNotRunning{},
			func() { alive.Serve(alive.WithForceOff()) },
			nil,
		},
		{
			"Stop already running",
			[]string{"", "serve", "off"},
			&alive.ServerStopped{},
			func() { alive.Serve() },
			nil,
		},
	}

	rtc.Load(rtc.RTCycler{
		NoParseArgs: true,
		AppDirName:  base.AppDirName,
		AppName:     base.AppName,
		Config:      &base.Conf,
	})

	rtc.SetTestMode()
	defer rtc.UnsetTestMode()

	app := &cli.App{
		Commands: []*cli.Command{
			Serve(),
		},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
				time.Sleep(5 * time.Second)
			}
			if tc.teardown != nil {
				t.Cleanup(tc.teardown)
			}
			err := app.Run(tc.command)
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expected) {
				t.Errorf("Expected type %T but got %T", tc.expected, err)
			}
		})
	}
}
