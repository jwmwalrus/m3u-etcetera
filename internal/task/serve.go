package task

import (
	"fmt"
	"net/url"

	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
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

func getGrpcOpts() (opts []grpc.DialOption) {
	opts = append(opts, grpc.WithInsecure())
	return
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

func parseLocations(locations []string) (parsed []string, err error) {
	for _, v := range locations {
		var u *url.URL
		if u, err = url.Parse(v); err != nil {
			return
		}
		if u.Scheme == "" {
			u.Scheme = "file"
		}
		parsed = append(parsed, u.String())
	}
	return
}

func serveAction(c *cli.Context) error {
	return alive.Serve(c.Bool("off"))
}
