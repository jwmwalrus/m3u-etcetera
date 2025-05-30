package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/status"
)

var (
	newPlaybackSvcClient = m3uetcpb.NewPlaybackSvcClient
)

// Playback playback task.
func Playback() *cli.Command {
	return &cli.Command{
		Name:        "playback",
		Aliases:     []string{"pb"},
		Category:    "Control",
		Usage:       "Controls the playback",
		Description: "Control the application's playback according with the given subcommand. If no subcommand is given, display current status.",
		Before:      checkServerStatus,
		Action:      playbackAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "Output JSON",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "play",
				Aliases:     []string{"pl"},
				Usage:       "Plays or resumes",
				ArgsUsage:   "LOCATION|ID ...",
				Description: "Play the given payload (i.e., list of LOCATION or ID) or resume a paused playback.",
				Action:      playbackPlayAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Usage: "add to playback immediately instead of queueing",
					},
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "use IDs instead of LOCATIONs",
					},
				},
			},
			{
				Name:        "pause",
				Aliases:     []string{"pa"},
				Usage:       "Pauses playback",
				Description: "Pause the playback.",
				Action:      playbackPlayAction,
			},
			{
				Name:        "stop",
				Usage:       "Stops playback",
				Description: "Stop the playback.",
				Action:      playbackPlayAction,
			},
			{
				Name:        "next",
				Usage:       "Next in playback",
				Description: "Play next track.",
				Action:      playbackPlayAction,
			},
			{
				Name:        "previous",
				Aliases:     []string{"prev"},
				Usage:       "Previous in playback",
				Description: "Play previous track.",
				Action:      playbackPlayAction,
			},
			{
				Name:        "seek",
				Usage:       "Seeks in playback",
				ArgsUsage:   "POSITION",
				Description: "Seek `POSITION` (seconds) in the current playback stream.",
				Action:      playbackSeekAction,
			},
			{
				Name:        "list",
				Aliases:     []string{"l"},
				Usage:       "Lists playback",
				Description: "Show the contents of the playback.",
				Action:      playbackListAction,
			},
		},
	}
}

func playbackAction(ctx context.Context, c *cli.Command) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybackSvcClient(cc)
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

	artist := res.Track.Artist
	if artist == "" {
		artist = res.Track.Albumartist
	}

	dur := time.Duration(res.Track.Duration) * time.Nanosecond

	tbl := table.New("ID", "Title", "Artist", "Album", "Duration")
	tbl.AddRow(
		res.Track.Id,
		res.Track.Title,
		artist,
		res.Track.Album,
		dur.Truncate(time.Second).String(),
	)
	tbl.Print()

	return
}

func playbackPlayAction(ctx context.Context, c *cli.Command) (err error) {
	const actionPrefix = "PB_"

	action := m3uetcpb.PlaybackAction_value[strings.ToUpper(actionPrefix+c.Name)]
	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction(action),
	}

	rest := c.Args().Slice()
	if c.Name == "play" {
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

	cl := newPlaybackSvcClient(cc)
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK\n")
	return
}

func playbackSeekAction(ctx context.Context, c *cli.Command) error {
	rest := c.Args().Slice()
	if len(rest) < 1 {
		return fmt.Errorf("I need one POSITION to seek")
	}
	if len(rest) > 1 {
		return fmt.Errorf("Too many values in command")
	}

	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_SEEK,
	}

	var err error
	if req.Seek, err = parseSeconds(rest[0]); err != nil {
		return err
	}

	cc, err := getClientConn()
	if err != nil {
		return err
	}
	defer cc.Close()

	cl := newPlaybackSvcClient(cc)
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf(s.Message())
	}

	fmt.Printf("OK\n")
	return nil
}

func playbackListAction(ctx context.Context, c *cli.Command) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybackSvcClient(cc)
	res, err := cl.GetPlaybackList(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
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

	tbl := table.New("ID", "TrackID", "Location")
	for _, e := range res.PlaybackEntries {
		un, _ := url.QueryUnescape(e.Location)
		if un == "" {
			un = e.Location
		}
		tbl.AddRow(e.Id, e.TrackId, un)
	}
	tbl.Print()
	return
}
