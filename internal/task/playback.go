package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// Playback playback task
func Playback() *cli.Command {
	return &cli.Command{
		Name:        "playback",
		Aliases:     []string{"pb"},
		Category:    "Control",
		Usage:       "Process the playback task",
		UsageText:   "playback [subcommand] ...",
		Description: "Control the application's playback according with the given subcommand. If no subcommand is given, display current status",
		Subcommands: []*cli.Command{
			{
				Name: "play",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Usage: "add to playback instead of queueing",
					},
				},
				Usage:       "playback play [ [--force] location ... ]",
				Description: "Plays the given payload or resumes a paused playback",
				Action:      playbackPlayAction,
			},
			{
				Name:        "pause",
				Usage:       "playback pause",
				Description: "pauses the playback",
				Action:      playbackPlayAction,
			},
			{
				Name:        "stop",
				Usage:       "playback stop",
				Description: "stops the playback",
				Action:      playbackPlayAction,
			},
			{
				Name:        "next",
				Usage:       "playback next",
				Description: "plays next track",
				Action:      playbackPlayAction,
			},
			{
				Name:        "previous",
				Usage:       "playback previous",
				Description: "plays previous track",
				Action:      playbackPlayAction,
			},
			{
				Name:        "jump",
				Usage:       "playback jump POS",
				Description: "jumps to a position in the current playback",
				Action:      playbackJumpAction,
			},
		},
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "playback",
		BashComplete: func(c *cli.Context) {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
		Before: checkServerStatus,
		Action: playbackAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "Output JSON",
			},
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "for shame\n")
			return err
		},
	}
}

func playbackAction(c *cli.Context) (err error) {
	opts := getGrpcOpts()

	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts...)
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	res, err := cl.GetPlayback(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		return
	}

	if _, ok := res.Playing.(*m3uetcpb.GetPlaybackResponse_Empty); ok {
		fmt.Printf("\nThere is no active playback\n")
		return
	}

	if c.Bool("json") {
		var bv []byte
		bv, err = json.MarshalIndent(res.Playing, "", "  ")
		if err != nil {
			return
		}
		fmt.Printf("\n%v\n", string(bv))
	} else {
		type jt struct {
			TrackID  int64
			Location string
		}
		j := jt{}
		switch res.Playing.(type) {
		case *m3uetcpb.GetPlaybackResponse_Playback:
			pb := res.GetPlayback()
			j = jt{pb.Id, pb.Location}
		case *m3uetcpb.GetPlaybackResponse_Track:
			t := res.GetTrack()
			j = jt{t.Id, t.Location}
		default:
		}

		tbl := table.New("ID", "Location")
		tbl.AddRow(j.TrackID, j.Location)
		tbl.Print()
	}

	return
}

func playbackPlayAction(c *cli.Context) (err error) {
	const actionPrefix = "PB_"
	opts := getGrpcOpts()

	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts...)
	if err != nil {
		return
	}
	defer cc.Close()

	action := m3uetcpb.PlaybackAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]
	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction(action),
	}

	rest := c.Args().Slice()
	if c.Command.Name == "play" {
		if len(rest) > 0 {
			req.Locations, err = parseLocations(rest)
			if err != nil {
				return
			}
		}
		req.Force = c.Bool("force")
	} else {
		if len(rest) > 0 {
			err = errors.New("Too many values in command")
			return
		}
	}

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}

func playbackJumpAction(c *cli.Context) (err error) {
	// TODO: implement
	fmt.Printf("TODO\n")
	return
}
