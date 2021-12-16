package main

import (
	"fmt"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/task"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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
		Name:     "m3uetc-task",
		Version:  "v0.20.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "John M",
				Email: "jwmwalrus@gmail.com",
			},
		},
		Copyright: "(c) 2021 WalrusAhead Solutions",
		// HelpName:  "m3uetc-task",
		Usage:     "Task interface for M3U Etcétera",
		UsageText: "m3uetc-task - control M3U Etcétera from the terminal",
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				log.Error(err)
				fmt.Fprintf(c.App.ErrWriter, err.Error()+"\n")
			}
		},
		Commands: []*cli.Command{
			task.Serve(),
			task.Playback(),
			task.Queue(),
			task.Collection(),
			task.Query(),
		},
	}

	app.Run(args)

	base.Unload()
}
