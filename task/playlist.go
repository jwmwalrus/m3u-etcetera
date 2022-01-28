package task

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

// Playlist playlist task
func Playlist() *cli.Command {
	return &cli.Command{
		Name:        "playlist",
		Aliases:     []string{"pl"},
		Category:    "Organization",
		Usage:       "Performs playlist-related actions",
		UsageText:   "playlist [subcommand] ...",
		Description: "Perform the playlist-related action given by the subcommand. When no subcommand is given, return all the playlists",
		Subcommands: []*cli.Command{
			{
				Name:    "info",
				Aliases: []string{"i"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "Output JSON",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "Limit the number of playlists shown",
					},
				},
				Usage:       "playlist info ID",
				Description: "Show information for the playlist with the given `ID`",
				Action:      playlistInfoAction,
			},
			{
				Name:    "create",
				Aliases: []string{"new"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Playlist name",
					},
					&cli.StringFlag{
						Name:    "descr",
						Aliases: []string{"d"},
						Usage:   "Playlist `DESCRIPTION`",
					},
				},
				Usage:       "playlist create [<flags> ...]",
				Description: "Creates a new playlist with values from the given flags",
				Action:      playlistExecuteAction,
			},
			{
				Name:    "update",
				Aliases: []string{"upd"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Playlist name",
					},
					&cli.StringFlag{
						Name:    "descr",
						Aliases: []string{"d"},
						Usage:   "Playlist `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:    "reset-descr",
						Aliases: []string{"rd"},
						Usage:   "Reset the description to an empty string",
					},
				},
				Usage:       "playlist update [<flags> ...] ID",
				Description: "Update the playlist identified by `ID` with values from the given flags",
				Action:      playlistExecuteAction,
			},
			{
				Name:        "destroy",
				Aliases:     []string{"del", "delete"},
				Usage:       "playlist delete ID",
				Description: "Delete playlist identified by the given `ID`",
				Action:      playlistExecuteAction,
			},
		},
		Before: checkServerStatus,
		Action: playlistAction,
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

func playlistAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)

	req := &m3uetcpb.GetAllPlaylistsRequest{Limit: int32(c.Int("limit"))}

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

	res, err := cl.GetAllPlaylists(context.Background(), req)
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

	tbl := table.New("ID", "Name", "Open", "Active", "Transient")
	for _, pl := range res.Playlists {
		tbl.AddRow(pl.Id, pl.Name, pl.Open, pl.Active, pl.Transient)
	}
	tbl.Print()
	return
}

func playlistInfoAction(c *cli.Context) (err error) {
	id, err := mustParseSingleID(c)
	if err != nil {
		return
	}

	err = showPlaylist(c, id)
	return
}

func playlistExecuteAction(c *cli.Context) (err error) {
	const actionPrefix = "PL_"

	var id int64
	if c.Command.Name == "create" {
		err = mustNotParseExtraArgs(c)
		if err != nil {
			return
		}
	} else {
		id, err = mustParseSingleID(c)
		if err != nil {
			return
		}
	}

	action := m3uetcpb.PlaylistAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]

	req := &m3uetcpb.ExecutePlaylistActionRequest{
		Action:      m3uetcpb.PlaylistAction(action),
		Id:          id,
		Name:        c.String("name"),
		Description: c.String("description"),
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	res, err := cl.ExecutePlaylistAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK, ID: %v\n", res.Id)
	return
}

func showPlaylist(c *cli.Context, id int64) (err error) {
	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)

	req := &m3uetcpb.GetPlaylistRequest{Id: id, Limit: int32(c.Int("limit"))}

	res, err := cl.GetPlaylist(context.Background(), req)
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

	pl := res.Playlist
	fmt.Printf(
		"\nPlaylist: %v\nActive: %v\nOpen: %v\nTransient:%v\n\n",
		pl.Name,
		pl.Active,
		pl.Open,
		pl.Transient,
	)

	tbl := table.New("Position", "Title", "Artist", "Album", "Dynamic")
	for _, pt := range res.PlaylistTracks {
		t := &m3uetcpb.Track{}
		for i := range res.Tracks {
			if res.Tracks[i].Id == pt.TrackId {
				t = res.Tracks[i]
				break
			}
		}
		artist := t.Albumartist
		if artist == "" {
			artist = t.Artist
		}
		tbl.AddRow(pt.Position, t.Title, artist, t.Album, pt.Dynamic)
	}
	tbl.Print()
	return
}
