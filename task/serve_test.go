package task

import (
	"reflect"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/urfave/cli/v2"
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
			func() { alive.Serve(alive.ServeOptions{TurnOff: true}) },
			func() { alive.Serve(alive.ServeOptions{TurnOff: true}) },
		},
		{
			"Start already running",
			[]string{"", "serve"},
			&alive.ServerAlreadyRunning{},
			func() { alive.Serve() },
			func() { alive.Serve(alive.ServeOptions{TurnOff: true}) },
		},
		{
			"Stop non-running",
			[]string{"", "serve", "off"},
			&alive.ServerNotRunning{},
			func() { alive.Serve(alive.ServeOptions{TurnOff: true}) },
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

	base.Load(true)

	base.SetTestingMode()
	defer base.UnsetTestingMode()

	app := &cli.App{
		Commands: []*cli.Command{
			Serve(),
		},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			err := app.Run(tc.command)
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expected) {
				t.Errorf("Expected type %T but got %T", tc.expected, err)
			}
			if tc.teardown != nil {
				tc.teardown()
			}
		})
	}
}
