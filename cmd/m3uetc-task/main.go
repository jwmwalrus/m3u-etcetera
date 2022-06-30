package main

import (
	"fmt"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/task"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

func main() {
	args := base.Load()

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "usage",
		Usage: "Print application usage",
	}
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "Print only the version",
	}

	app := &cli.App{
		Name:    "m3uetc-task",
		Version: "v0.20.0",
		Authors: []*cli.Author{
			{
				Name:  "John M",
				Email: "jwmwalrus@gmail.com",
			},
		},
		Copyright:   "(c) 2021 WalrusAhead Solutions",
		HelpName:    "m3uetc-task",
		Usage:       "Task interface for M3U Etcétera",
		UsageText:   "m3uetc-task command [subcommand [--flags...] [args...]]",
		Description: "A playlist-centric music player",
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				s, ok := status.FromError(err)
				if !ok {
					log.Error(err)
					fmt.Fprintf(c.App.ErrWriter, err.Error()+"\n")
					return
				}

				log.WithFields(log.Fields{
					"code":    s.Code(),
					"details": s.Details(),
				}).Error(s.Message())
				fmt.Fprintf(c.App.ErrWriter, s.Message()+"\n")
			}
		},
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			task.Serve(),
			task.Playback(),
			task.Queue(),
			task.Collection(),
			task.Query(),
			task.Playbar(),
			task.Playlist(),
			task.Playtrack(),
			task.Playgroup(),
			task.Perspective(),
		},
	}

	app.Run(args)

	base.Unload()
}
