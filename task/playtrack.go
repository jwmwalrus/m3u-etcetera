package task

import (
	"context"
	"fmt"
	"strings"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/status"
)

// Playtrack playtrack task.
func Playtrack() *cli.Command {
	return &cli.Command{
		Name:        "playtrack",
		Aliases:     []string{"pt"},
		Category:    "Organization",
		Usage:       "Performs playtrack-related actions",
		ArgsUsage:   "ID",
		Description: "Perform the playtrack-related action on playlist tracks. When no subcommand is given, display all the tracks in the playlist identified by `ID`.",
		Before:      checkServerStatus,
		Action:      playtrackAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "output JSON",
			},
			&cli.StringFlag{
				Name:  "persp",
				Usage: "applies to `PERSPECTIVE`",
				Value: "music",
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "limit the `NUMBER` of playlists shown",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "append",
				Aliases:     []string{"app"},
				Usage:       "Appends tracks",
				ArgsUsage:   "ID|LOCATION ...",
				Description: "Add track(s) at the end of playlist.",
				Action:      playtrackExecuteAction,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "use playlist id, `PL`",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "use IDs instead of LOCATIONs",
					},
				},
			},
			{
				Name:        "prepend",
				Aliases:     []string{"prep"},
				Usage:       "Prepends tracks",
				ArgsUsage:   "ID|LOCATION ...",
				Description: "Add track(s) at the beginning of playlist.",
				Action:      playtrackExecuteAction,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "use playlist id, `PL`",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "use IDs instead of LOCATIONs",
					},
				},
			},
			{
				Name:        "insert",
				Aliases:     []string{"i"},
				Usage:       "Inserts tracks",
				ArgsUsage:   "ID|LOCATION ...",
				Description: "Insert track(s) at the given `POSITION` in playlist.",
				Action:      playtrackExecuteAction,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "use playlist id, `PL`",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "pos",
						Aliases: []string{"p"},
						Usage:   "insert at the given `POSITION`",
					},
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "use IDs instead of LOCATIONs",
					},
				},
			},
			{
				Name:        "move",
				Aliases:     []string{"m"},
				Usage:       "Moves track",
				Description: "Move playlist track from one position to another.",
				Action:      playtrackExecuteAction,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "use playlist id, `PL`",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "from-pos",
						Aliases: []string{"from"},
						Usage:   "move the given `POSITION`",
					},
					&cli.IntFlag{
						Name:    "to-pos",
						Aliases: []string{"to"},
						Usage:   "move to the given `POSITION`",
					},
				},
			},
			{
				Name:        "delete",
				Aliases:     []string{"del"},
				Usage:       "Delete track",
				Description: "Delete track at the given `POSITION` in playlist.",
				Action:      playtrackExecuteAction,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "use playlist id, `PL`",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "pos",
						Aliases: []string{"p"},
						Usage:   "delete at the given `POSITION`",
					},
				},
			},
			{
				Name:        "clear",
				Usage:       "Clears tracks",
				Description: "Remove all tracks in playlist.",
				Action:      playtrackExecuteAction,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "playlist",
						Aliases:  []string{"pl"},
						Usage:    "use playlist id, `PL`",
						Required: true,
					},
				},
			},
		},
	}
}

func playtrackAction(ctx context.Context, c *cli.Command) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybarSvcClient(cc)

	req := &m3uetcpb.GetPlaybarRequest{
		Perspective: getPerspective(c),
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

	err = showPlaylist(ctx, c, id)
	return
}

func playtrackExecuteAction(ctx context.Context, c *cli.Command) (err error) {
	const actionPrefix = "PT_"

	rest := c.Args().Slice()
	if (c.Name == "move" ||
		c.Name == "delete" ||
		c.Name == "clear") &&
		len(rest) > 0 {
		err = fmt.Errorf("Too many values in command")
		return
	}
	if (c.Name == "append" ||
		c.Name == "prepend" ||
		c.Name == "insert") && len(rest) < 1 {
		err = fmt.Errorf("I need a list of locations or IDs")
		return
	}

	var frompos, topos int
	if c.Name == "insert" ||
		c.Name == "delete" {
		if c.Int("pos") < 1 {
			err = fmt.Errorf("I need a position to insert|delete")
			return
		}
		topos = int(c.Int("pos"))
	}

	if c.Name == "move" {
		if c.Int("from-pos") < 1 ||
			c.Int("to-pos") < 1 {
			err = fmt.Errorf("I need a valid positions to move")
			return
		}
		frompos = int(c.Int("from-pos"))
		topos = int(c.Int("to-pos"))
	}

	action := m3uetcpb.PlaylistTrackAction_value[strings.ToUpper(actionPrefix+c.Name)]

	req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
		Action:       m3uetcpb.PlaylistTrackAction(action),
		PlaylistId:   c.Int("playlist"),
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

	cl := newPlaybarSvcClient(cc)
	_, err = cl.ExecutePlaylistTrackAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK\n")
	return
}
