package main

import (
	"fmt"
	"log/slog"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/task"
	rtc "github.com/jwmwalrus/rtcycler"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

func main() {
	args := rtc.Load(&base.Conf,
		base.AppDirName,
		rtc.WithAppName(base.AppName),
	)

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "usage",
		Usage: "print application usage",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "print only the version",
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
		Usage:       "CLI interface for M3U Etcétera",
		Description: "A playlist-centric music player.",
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				s, ok := status.FromError(err)
				if !ok {
					slog.Error("Command finished with error status", "error", err)
					fmt.Fprintf(c.App.ErrWriter, err.Error()+"\n")
					return
				}

				slog.With(
					"code", s.Code(),
					"details", s.Details(),
				).Error(s.Message())
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
		Action: task.DefaultAction,
	}

	app.Run(args)

	rtc.Unload()
}
