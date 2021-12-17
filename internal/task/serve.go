package task

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/urfave/cli/v2"
)

// Serve serve task
func Serve() *cli.Command {
	return &cli.Command{
		Name:        "serve",
		Category:    "Server",
		Usage:       "Controls the server",
		UsageText:   "serve [--off]",
		Description: "Starts or stops the m3uetc-server",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "off"},
		},
		Action: serveAction,
	}
}

func checkServerStatus(c *cli.Context) (err error) {
	err = alive.CheckServerStatus()
	switch err.(type) {
	case *alive.ServerAlreadyRunning,
		*alive.ServerStarted:
		err = nil
	default:
	}
	return
}

func serveAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	err = alive.Serve(c.Bool("off"))
	return
}
