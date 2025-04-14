package task

import (
	"context"

	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/urfave/cli/v3"
)

// Serve serve task.
func Serve() *cli.Command {
	return &cli.Command{
		Name:        "serve",
		Category:    "Server",
		Usage:       "Controls the server",
		Description: "Starts or stops the m3uetc-server.",
		Commands: []*cli.Command{
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

func checkServerStatus(ctx context.Context, c *cli.Command) (context.Context, error) {
	err := alive.CheckServerStatus()
	switch err.(type) {
	case *alive.ServerAlreadyRunning,
		*alive.ServerStarted:
		err = nil
	default:
	}
	return ctx, err
}

func serveAction(ctx context.Context, c *cli.Command) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	err = alive.Serve()
	return
}

func serveOffAction(ctx context.Context, c *cli.Command) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	err = alive.Serve(alive.WithTurnOff(),
		alive.WithForceOff(c.Bool("force")),
	)
	return
}
