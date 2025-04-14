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
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/status"
)

// Playlist playlist task.
func Playlist() *cli.Command {
	return &cli.Command{
		Name:        "playlist",
		Aliases:     []string{"pl"},
		Category:    "Organization",
		Usage:       "Performs playlist-related actions",
		Description: "Perform the playlist-related action given by the subcommand. When no subcommand is given, return all the playlists.",
		Before:      checkServerStatus,
		Action:      playlistAction,
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
				Name:        "info",
				Aliases:     []string{"i"},
				Usage:       "Shows playlist info",
				ArgsUsage:   "ID",
				Description: "Show information for the playlist with the given `ID`.",
				Action:      playlistInfoAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "output JSON",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "limit the `NUMBER` of tracks shown",
					},
				},
			},
			{
				Name:        "create",
				Aliases:     []string{"new"},
				Usage:       "Creates playlist",
				Description: "Create a new playlist, according to the given options.",
				Action:      playlistExecuteAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "playlist `NAME`",
					},
					&cli.StringFlag{
						Name:    "descr",
						Aliases: []string{"d"},
						Usage:   "playlist `DESCRIPTION`",
					},
				},
			},
			{
				Name:        "update",
				Aliases:     []string{"upd"},
				Usage:       "Updates playlist",
				ArgsUsage:   "ID",
				Description: "Update the playlist identified by `ID`, according to the given options.",
				Action:      playlistExecuteAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "playlist `NAME`",
					},
					&cli.StringFlag{
						Name:    "descr",
						Aliases: []string{"d"},
						Usage:   "Playlist `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:    "reset-descr",
						Aliases: []string{"rd"},
						Usage:   "reset the description to an empty string",
					},
					&cli.BoolFlag{
						Name:  "bucket",
						Usage: "set the playlist as bucket",
					},
					&cli.BoolFlag{
						Name:  "no-bucket",
						Usage: "unset the playlist as bucket",
					},
				},
			},
			{
				Name:        "destroy",
				Aliases:     []string{"del", "delete"},
				ArgsUsage:   "ID",
				Usage:       "Deletes playlist",
				Description: "Delete playlist identified by the given `ID`.",
				Action:      playlistExecuteAction,
			},
			{
				Name:        "merge",
				Usage:       "Merge playlists",
				UsageText:   "ID1 ID2",
				Description: "Merge playlists identified by `ID1` and `ID2`. The merge playlist will be identified by ID1, and the playlist identified by ID2 will be deleted.",
				Action:      playlistExecuteAction,
			},
			{
				Name:        "import",
				Aliases:     []string{"imp"},
				Usage:       "Imports playlist",
				Description: "Import a playlist from a supported file format (M3U|PLS).",
				Action:      playlistImportAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "persp",
						Usage: "applies to `PERSPECTIVE`",
						Value: "music",
					},
				},
			},
			{
				Name:        "export",
				Aliases:     []string{"exp"},
				Usage:       "Exports playlist",
				ArgsUsage:   "ID LOCATION",
				Description: "Export the playlist identified by `ID` to the given `LOCATION`.",
				Action:      playlistExportAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"f"},
						Usage:   "playlist `FORMAT` (M3U|PLS)",
						Value:   "M3U",
					},
				},
			},
			{
				Name:        "open",
				Usage:       "Open playlist",
				ArgsUsage:   "LOCATION ...",
				Description: "Open the playlist at the given `LOCATION`. Supported formats: M3U|PlS",
				Action:      playlistOpenFromLocationsAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "persp",
						Usage: "applies to `PERSPECTIVE`",
						Value: "music",
					},
				},
			},
		},
	}
}

func playlistAction(ctx context.Context, c *cli.Command) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybarSvcClient(cc)

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

func playlistInfoAction(ctx context.Context, c *cli.Command) (err error) {
	id, err := mustParseSingleID(c)
	if err != nil {
		return
	}

	err = showPlaylist(ctx, c, id)
	return
}

func playlistExecuteAction(ctx context.Context, c *cli.Command) (err error) {
	const actionPrefix = "PL_"

	var id, id2 int64
	if c.Name == "create" {
		err = mustNotParseExtraArgs(c)
		if err != nil {
			return
		}
	} else if c.Name == "merge" {
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

	action := m3uetcpb.PlaylistAction_value[strings.ToUpper(actionPrefix+c.Name)]

	var bucket int
	if c.Bool("bucket") {
		bucket = 1
	} else if c.Bool("bucket") {
		bucket = 2
	}

	req := &m3uetcpb.ExecutePlaylistActionRequest{
		Action:      m3uetcpb.PlaylistAction(action),
		Id:          id,
		Name:        c.String("name"),
		Description: c.String("descr"),
		Bucket:      int32(bucket),
	}

	if c.Name == "merge" {
		req.Id2 = id2
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybarSvcClient(cc)
	res, err := cl.ExecutePlaylistAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK, ID: %v\n", res.Id)
	return
}

func showPlaylist(ctx context.Context, c *cli.Command, id int64) (err error) {
	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybarSvcClient(cc)

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

func playlistImportAction(ctx context.Context, c *cli.Command) error {
	rest := c.Args().Slice()
	if len(rest) < 1 {
		return fmt.Errorf("I need a list of locations to playlists")
	}

	req := &m3uetcpb.ImportPlaylistsRequest{
		Perspective: getPerspective(c),
	}

	var err error
	if req.Locations, err = parseLocations(rest); err != nil {
		return err
	}

	return istImportPlaylists(req)
}

func playlistExportAction(ctx context.Context, c *cli.Command) (err error) {

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

	cl := newPlaybarSvcClient(cc)
	_, err = cl.ExportPlaylist(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK\n")
	return
}

func playlistOpenFromLocationsAction(ctx context.Context, c *cli.Command) error {
	rest := c.Args().Slice()
	if len(rest) < 1 {
		return fmt.Errorf("I need a list of locations to playlists")
	}

	req := &m3uetcpb.ImportPlaylistsRequest{
		Perspective: getPerspective(c),
		AsTransient: true,
	}

	var err error
	if req.Locations, err = parseLocations(rest); err != nil {
		return err
	}

	return istImportPlaylists(req)
}

func istImportPlaylists(req *m3uetcpb.ImportPlaylistsRequest) (err error) {
	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybarSvcClient(cc)
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
