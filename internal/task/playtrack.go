package task

import (
	"context"
	"fmt"
	"strings"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

// Playtrack playtrack task
func Playtrack() *cli.Command {
	return &cli.Command{
		Name:        "playtrack",
		Aliases:     []string{"pt"},
		Category:    "Organization",
		Usage:       "Performs playtrack-related actions",
		UsageText:   "playtrack [subcommand] ...",
		Description: "Perform the playtrack-related action on playlist tracks, using the ID given by the --playlist flag, and the subcommand. When no subcommand is given, display all the tracks in the playlist",
		Subcommands: []*cli.Command{
			{
				Name:    "append",
				Aliases: []string{"app"},
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "Use playlist id, `PL`",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "Use IDs instead of locations",
					},
				},
				Usage:       "playtrack append --playlist PL ID|location ...",
				Description: "Add track(s) at the end of playlist",
				Action:      playtrackExecuteAction,
			},
			{
				Name:    "preppend",
				Aliases: []string{"prep"},
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "Use playlist id, `PL`",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "Use IDs instead of locations",
					},
				},
				Usage:       "playtrack preppend --playlist PL ID|location ...",
				Description: "Add track(s) at the beginning of playlist",
				Action:      playtrackExecuteAction,
			},
			{
				Name:    "insert",
				Aliases: []string{"i"},
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "Use playlist id, `PL`",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "pos",
						Aliases: []string{"p"},
						Usage:   "Insert at the given `POSITION`",
					},
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "Use IDs instead of locations",
					},
				},
				Usage:       "playtrack insert --playlist PL --pos POSITION ID|location ...",
				Description: "Insert track(s) at the given `POSITION` in playlist",
				Action:      playtrackExecuteAction,
			},
			{
				Name:    "move",
				Aliases: []string{"m"},
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "Use playlist id, `PL`",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "from-pos",
						Aliases: []string{"from"},
						Usage:   "Move the given `POSITION`",
					},
					&cli.IntFlag{
						Name:    "to-pos",
						Aliases: []string{"to"},
						Usage:   "Move to the given `POSITION`",
					},
				},
				Usage:       "playtrack move --playlist PL --from-pos FROM --to-pos TO",
				Description: "Move playlist track from one position to another",
				Action:      playtrackExecuteAction,
			},
			{
				Name:    "delete",
				Aliases: []string{"del"},
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "Use playlist id, `PL`",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "pos",
						Aliases: []string{"p"},
						Usage:   "Remove from the given `POSITION`",
					},
				},
				Usage:       "playtrack remove --playlist PL --pos POSITION",
				Description: "Remove track at the given `POSITION` in playlist",
				Action:      playtrackExecuteAction,
			},
			{
				Name: "clear",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "Use playlist id, `PL`",
						Required: true,
					},
				},
				Usage:       "playtrack clear --playlist PL",
				Action:      playtrackExecuteAction,
				Description: "Remove all tracks in playlist",
			},
		},
		Before: checkServerStatus,
		Action: playtrackAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "Output JSON",
			},
			&cli.StringFlag{
				Name:    "perspective",
				Aliases: []string{"persp"},
				Usage:   "Applies to perspective",
				Value:   "music",
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit the number of playlists shown",
			},
		},
	}
}

func playtrackAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)

	req := &m3uetcpb.GetPlaybarRequest{}

	persp := strings.ToLower(c.String("perspective"))
	if strings.HasPrefix("radio", persp) {
		req.Perspective = m3uetcpb.Perspective_RADIO
	} else if strings.HasPrefix("podcasts", persp) {
		req.Perspective = m3uetcpb.Perspective_PODCASTS
	} else if strings.HasPrefix("audiobooks", persp) {
		req.Perspective = m3uetcpb.Perspective_AUDIOBOOKS
	} else {
		req.Perspective = m3uetcpb.Perspective_MUSIC
	}
	res, err := cl.GetPlaybar(context.Background(), req)
	if err != nil {
		return
	}

	var id int64
	for _, pl := range res.Playlists {
		if pl.Active {
			id = pl.Id
			break
		}
	}

	if id == 0 {
		fmt.Printf("There's no active playlist")
	}

	err = showPlaylist(c, id)
	return
}

func playtrackExecuteAction(c *cli.Context) (err error) {
	const actionPrefix = "PT_"

	rest := c.Args().Slice()
	if (c.Command.Name == "move" ||
		c.Command.Name == "delete" ||
		c.Command.Name == "clear") &&
		len(rest) > 0 {
		err = fmt.Errorf("Too many values in command")
		return
	}
	if (c.Command.Name == "append" ||
		c.Command.Name == "preppend" ||
		c.Command.Name == "insert") && len(rest) < 1 {
		err = fmt.Errorf("I need a list of locations or IDs")
		return
	}

	var frompos, topos int
	if c.Command.Name == "insert" ||
		c.Command.Name == "delete" {
		if c.Int("pos") < 1 {
			err = fmt.Errorf("I need a position to insert|delete")
			return
		}
		topos = c.Int("pos")
	}

	if c.Command.Name == "move" {
		if c.Int("from-pos") < 1 ||
			c.Int("to-pos") < 1 {
			err = fmt.Errorf("I need a valid positions to move")
			return
		}
		frompos = c.Int("from-pos")
		topos = c.Int("to-pos")
	}

	action := m3uetcpb.PlaylistTrackAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]

	req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
		Action:       m3uetcpb.PlaylistTrackAction(action),
		PlaylistId:   c.Int64("playlist"),
		Position:     int32(topos),
		FromPosition: int32(frompos),
	}

	if c.Bool("ids") {
		if req.TrackIds, err = parseIDs(rest); err != nil {
			return
		}
	} else {
		if req.Locations, err = parseLocations(rest); err != nil {
			return
		}
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	_, err = cl.ExecutePlaylistTrackAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK\n")
	return
}
