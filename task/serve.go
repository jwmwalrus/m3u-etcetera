package task

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/urfave/cli/v2"
)

// Serve serve task.
func Serve() *cli.Command {
	return &cli.Command{
		Name:        "serve",
		Category:    "Server",
		Usage:       "Controls the server",
		Description: "Starts or stops the m3uetc-server.",
		Subcommands: []*cli.Command{
			{
				Name:        "off",
				Usage:       "Terminates the server",
				Description: "Terminate the server.",
				Action:      serveOffAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Usage: "force server termination",
					},
				},
			},
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

	err = alive.Serve()
	return
}

func serveOffAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	err = alive.Serve(alive.WithTurnOff(),
		alive.WithForceOff(c.Bool("force")),
	)
	return
}
