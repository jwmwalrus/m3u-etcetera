package task

import (
	"errors"
	"fmt"

	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/urfave/cli/v2"
)

// Serve serve task
func Serve() *cli.Command {
	return &cli.Command{
		Name:            "serve",
		Category:        "Server",
		Usage:           "Controls the server",
		UsageText:       "serve [--off]",
		Description:     "Starts or stops the m3uetc-server",
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "doo!",
		BashComplete: func(c *cli.Context) {
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "off"},
		},
		Action: serveAction,
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			fmt.Fprintf(c.App.Writer, "for shame\n")
			return err
		},
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
	rest := c.Args().Slice()
	if len(rest) > 0 {
		err = errors.New("Too many values in command")
		return
	}

	err = alive.Serve(c.Bool("off"))
	return
}
