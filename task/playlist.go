package task

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/jwmwalrus/bnp/urlstr"
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
			{
				Name:        "merge",
				Usage:       "playlist merge ID1 ID2",
				Description: "Merge playlists identified by `ID1` and `ID2`. The merge playlist will be identified by ID1, and the playlist identified by ID2 will be deleted.",
				Action:      playlistExecuteAction,
			},
			{
				Name:    "import",
				Aliases: []string{"imp"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "persp",
						Usage: "Applies to `PERSPECTIVE`",
						Value: "music",
					},
				},
				Usage:       "playlist import [<flags> ...] locations ...",
				Description: "Imports a playlist from a supported file format (e.g., M3U)",
				Action:      playlistImportAction,
			},
			{
				Name:    "export",
				Aliases: []string{"exp"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"f"},
						Usage:   "Playlist `FORMAT` (M3U|PLS)",
						Value:   "M3U",
					},
				},
				Usage:       "playlist export [--format FORMAT] ID location",
				Description: "Exports a playlist to a supported file format (e.g., M3U)",
				Action:      playlistExportAction,
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
				Name:  "persp",
				Usage: "Applies to `PERSPECTIVE`",
				Value: "music",
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

	req := &m3uetcpb.GetAllPlaylistsRequest{
		Perspective: getPerspective(c),
		Limit:       int32(c.Int("limit")),
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

	var id, id2 int64
	if c.Command.Name == "create" {
		err = mustNotParseExtraArgs(c)
		if err != nil {
			return
		}
	} else if c.Command.Name == "merge" {
		rest := c.Args().Slice()
		if len(rest) != 2 {
			err = fmt.Errorf("I need two playlist IDs")
			return
		}
		var s []int64
		if s, err = parseIDs(rest); err != nil {
			return
		}
		id, id2 = s[0], s[1]
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
		Description: c.String("descr"),
	}

	if c.Command.Name == "merge" {
		req.Id2 = id2
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

func playlistImportAction(c *cli.Context) (err error) {

	rest := c.Args().Slice()
	if len(rest) < 1 {
		err = fmt.Errorf("I need a list of locations to playlists")
		return
	}

	req := &m3uetcpb.ImportPlaylistsRequest{
		Perspective: getPerspective(c),
	}

	if req.Locations, err = parseLocations(rest); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	stream, err := cl.ImportPlaylists(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	for {
		var res *m3uetcpb.ImportPlaylistsResponse
		res, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		fmt.Printf("ID: %d\n", res.Id)
		if len(res.ImportErrors) > 0 {
			fmt.Printf("\tThere were errors during import:\n")
			for _, ie := range res.ImportErrors {
				fmt.Printf("\t\t%s\n", ie)
			}
			fmt.Printf("\n")
		}

	}

	return
}

func playlistExportAction(c *cli.Context) (err error) {

	rest := c.Args().Slice()
	if len(rest) != 2 {
		err = fmt.Errorf("I need an ID and a location")
		return
	}

	format, ok := m3uetcpb.PlaylistExportFormat_value["PLEF_"+c.String("format")]
	if !ok {
		err = fmt.Errorf("Unknown format: %v", c.String("format"))
		return
	}

	ids, err := parseIDs([]string{rest[0]})
	if err != nil {
		return
	}

	path, err := urlstr.PathToURLUnchecked(rest[1])
	if err != nil {
		return
	}

	req := &m3uetcpb.ExportPlaylistRequest{
		Id:       ids[0],
		Location: path,
		Format:   m3uetcpb.PlaylistExportFormat(format),
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	_, err = cl.ExportPlaylist(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK\n")
	return
}
