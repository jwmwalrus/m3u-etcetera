package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
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
				Name:    "play",
				Aliases: []string{"pl"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Usage: "add to playback instead of queueing",
					},
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "Use IDs instead of locations",
					},
				},
				Usage:       "playback play [ [--force] [location ... | --ids id ... ]",
				Description: "Plays the given payload or resumes a paused playback",
				Action:      playbackPlayAction,
			},
			{
				Name:        "pause",
				Aliases:     []string{"pa"},
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
				Aliases:     []string{"prev"},
				Usage:       "playback previous",
				Description: "plays previous track",
				Action:      playbackPlayAction,
			},
			{
				Name:        "jump",
				Usage:       "playback jump POSITION",
				Description: "jumps to a `POSITION` in the current playback",
				Action:      playbackJumpAction,
			},
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
	}
}

func playbackAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	res, err := cl.GetPlayback(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		return
	}

	if !res.IsStreaming {
		fmt.Printf("\nThere is no active playback\n")
		return
	}

	if c.Bool("json") {
		var bv []byte
		bv, err = json.MarshalIndent(res, "", "  ")
		if err != nil {
			return
		}
		fmt.Printf("\n%v\n", string(bv))
		return
	}

	if res.Track.Title == "" {

		tbl := table.New("ID", "Location")
		un, _ := url.QueryUnescape(res.Playback.Location)
		if un == "" {
			un = res.Playback.Location
		}
		tbl.AddRow(res.Playback.Id, un)
		tbl.Print()
	}

	tbl := table.New("ID", "Title", "Artist", "Album", "Duration")
	artist := res.Track.Artist
	if artist == "" {
		artist = res.Track.Albumartist
	}
	tbl.AddRow(res.Track.Id, res.Track.Title, artist, res.Track.Album, res.Track.Duration)
	tbl.Print()

	return
}

func playbackPlayAction(c *cli.Context) (err error) {
	const actionPrefix = "PB_"

	action := m3uetcpb.PlaybackAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]
	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction(action),
	}

	rest := c.Args().Slice()
	if c.Command.Name == "play" {
		if len(rest) > 0 {
			if c.Bool("ids") {
				if req.Ids, err = parseIDs(rest); err != nil {
					return
				}
			} else {
				if req.Locations, err = parseLocations(rest); err != nil {
					return
				}
			}
		}
		req.Force = c.Bool("force")
	} else {
		if len(rest) > 0 {
			err = fmt.Errorf("Too many values in command")
			return
		}
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
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