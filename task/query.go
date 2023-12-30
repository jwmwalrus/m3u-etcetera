package task

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	newQuerySvcClient = m3uetcpb.NewQuerySvcClient
)

// Query defines the query-related tasks.
func Query() *cli.Command {
	return &cli.Command{
		Name:        "query",
		Aliases:     []string{"search", "s"},
		Category:    "Control",
		Usage:       "Queries the database",
		Description: "Performs query-related actions. When no subcommand is given, display list of queries.",
		Before:      checkServerStatus,
		Action:      queryAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "output JSON",
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "limit output count",
				Value: 0,
			},
			&cli.Int64SliceFlag{
				Name:    "collection-id",
				Aliases: []string{"coll-id"},
				Usage:   "bound to collection ID",
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:        "add",
				Usage:       "Adds query",
				Description: "Add query, according to the given options.",
				Action:      queryAddAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "query `NAME`",
					},
					&cli.StringFlag{
						Name:  "descr",
						Usage: "query `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:  "random",
						Usage: "query is random",
					},
					&cli.IntFlag{
						Name:  "rating",
						Usage: "query `RATING`",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "query `LIMIT`",
					},
					&cli.StringFlag{
						Name:  "params",
						Usage: "query `PARAMS` for title,artist,album,genre (e.g.: \"title=thing and genre=[sh]ome or genre=some*other\").",
					},
					&cli.Int64Flag{
						Name:  "from",
						Usage: "query's start `TIMESTAMP` (i.e., from the date the track was issued)",
					},
					&cli.Int64Flag{
						Name:  "to",
						Usage: "query's end `TIMESTAMP` (i.e., to the date the track was issued)",
					},
					&cli.Int64SliceFlag{
						Name:    "collection-id",
						Aliases: []string{"coll-id"},
						Usage:   "`ID` for the collection bounding the query (can appear more than once in the command line)",
					},
				},
			},
			{
				Name:        "info",
				Aliases:     []string{"i"},
				Usage:       "Shows query info",
				ArgsUsage:   "ID",
				Description: "Show information for the query identified by `ID`.",
				Action:      queryInfoAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "output JSON",
					},
				},
			},
			{
				Name:        "remove",
				Aliases:     []string{"rem"},
				Usage:       "Removes query",
				ArgsUsage:   "ID",
				Description: "Remove the query identified by `ID`.",
				Action:      queryRemoveAction,
			},
			{
				Name:        "update",
				Aliases:     []string{"upd"},
				Usage:       "Updates query",
				Description: "Update query according to the given options.",
				Action:      queryUpdateAction,
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "id",
						Usage:    "query's existing `ID`",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "name",
						Usage: "query `NAME`",
					},
					&cli.StringFlag{
						Name:  "descr",
						Usage: "query `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:  "random",
						Usage: "query is random",
					},
					&cli.BoolFlag{
						Name:  "no-random",
						Usage: "query is not random",
					},
					&cli.IntFlag{
						Name:  "rating",
						Usage: "query `RATING`",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "query `LIMIT`",
					},
					&cli.StringFlag{
						Name:  "params",
						Usage: "query `PARAMS` for title,artist,album,genre (e.g.: \"title=thing and genre=[sh]ome or genre=some*other\"",
					},
					&cli.Int64Flag{
						Name:  "from",
						Usage: "query's start `TIMESTAMP` (i.e., from the date the track was issued)",
					},
					&cli.Int64Flag{
						Name:  "to",
						Usage: "query's end `TIMESTAMP` (i.e., to the date the track was issued)",
					},
					&cli.Int64SliceFlag{
						Name:    "collection-id",
						Aliases: []string{"coll-id"},
						Usage:   "`ID` for the collection bounding the query (can appear more than once in the command line)",
					},
				},
			},
			{
				Name:        "inplaylist",
				Aliases:     []string{"inpl"},
				Usage:       "Query in playlist",
				ArgsUsage:   "ID",
				Description: "Add the results of the query identified by `ID` to a playlist.",
				Action:      queryInPlaylistAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "output JSON",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "limit the number of tracks shown",
					},
					&cli.BoolFlag{
						Name:  "play",
						Usage: "add all playlist tracks to playback instead of listing them",
					},
					&cli.BoolFlag{
						Name:  "force",
						Usage: "force playback",
					},
					&cli.IntFlag{
						Name:  "pl",
						Usage: "playlist ID",
					},
				},
			},
			{
				Name:        "inqueue",
				Aliases:     []string{"inq"},
				Usage:       "Query in queue",
				ArgsUsage:   "ID",
				Description: "Add the results of the query identified by `ID` to the current queue.",
				Action:      queryInQueueAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "output JSON",
					},
				},
			},
			{
				Name:        "by",
				Usage:       "Query by params",
				ArgsUsage:   "PARAMS",
				Description: "Perform a query by the given `PARAMS`.",
				Action:      queryByAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "persist-as",
						Usage: "persist/save query as `NAME`",
					},
					&cli.StringFlag{
						Name:  "descr",
						Usage: "query `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:  "random",
						Usage: "query is random",
					},
					&cli.IntFlag{
						Name:  "rating",
						Usage: "query `RATING`",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "query `LIMIT`",
					},
					&cli.Int64Flag{
						Name:  "from",
						Usage: "query's start `TIMESTAMP` (i.e., from the date the track was issued)",
					},
					&cli.Int64Flag{
						Name:  "to",
						Usage: "query's end `TIMESTAMP` (i.e., to the date the track was issued)",
					},
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "output JSON",
					},
					&cli.BoolFlag{
						Name:    "play",
						Aliases: []string{"pl"},
						Usage:   "add tracks to playback instead of listing them",
					},
					&cli.BoolFlag{
						Name:  "force",
						Usage: "force playback",
					},
				},
			},
		},
	}
}

func queryAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	req := &m3uetcpb.GetQueriesRequest{
		CollectionIds: c.Int64Slice("collection-id"),
		Limit:         int32(c.Int("limit")),
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newQuerySvcClient(cc)
	res, err := cl.GetQueries(context.Background(), req)
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

	tbl := table.New("ID", "Name", "Params", "Limit", "Random", "Bounds")
	for _, q := range res.Queries {
		b := ""
		if len(q.CollectionIds) > 0 {
			b = "C"
		}
		name := q.Name
		if q.ReadOnly {
			name = q.Description
		}
		tbl.AddRow(q.Id, name, q.Params, q.Limit, q.Random, b)
	}
	tbl.Print()

	return
}

func queryInfoAction(c *cli.Context) (err error) {
	var id int64
	if id, err = mustParseSingleID(c); err != nil {
		return
	}

	req := &m3uetcpb.GetQueryRequest{
		Id: id,
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newQuerySvcClient(cc)
	res, err := cl.GetQuery(context.Background(), req)
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

	tbl := table.New("ID", "Name", "Params", "Limit", "Random", "Bounds")
	q := res.Query
	b := ""
	if len(q.CollectionIds) > 0 {
		b = "C"
	}
	tbl.AddRow(q.Id, q.Name, q.Params, q.Limit, q.Random, b)

	tbl.Print()
	return
}

func queryAddAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	var from, to *timestamppb.Timestamp

	ts := c.Int64("from")
	if ts > 0 {
		from = timestamppb.New(time.Unix(0, ts))
	}

	ts = c.Int64("to")
	if ts > 0 {
		to = timestamppb.New(time.Unix(0, ts))
	}

	q := &m3uetcpb.Query{
		Name:          c.String("name"),
		Description:   c.String("descr"),
		Random:        c.Bool("random"),
		Rating:        int32(c.Int("rating")),
		Limit:         int32(c.Int("limit")),
		Params:        c.String("params"),
		From:          from,
		To:            to,
		CollectionIds: c.Int64Slice("collection-id"),
	}
	req := &m3uetcpb.AddQueryRequest{
		Query: q,
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newQuerySvcClient(cc)
	res, err := cl.AddQuery(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("ID: %v\n", res.Id)
	return
}

func queryRemoveAction(c *cli.Context) (err error) {
	var id int64
	if id, err = mustParseSingleID(c); err != nil {
		return
	}

	req := &m3uetcpb.RemoveQueryRequest{
		Id: id,
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newQuerySvcClient(cc)
	_, err = cl.RemoveQuery(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}

func queryUpdateAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	if c.Int64("id") < 1 {
		err = fmt.Errorf("I need an ID greater than zero")
		return
	}

	req0 := &m3uetcpb.GetQueryRequest{Id: c.Int64("id")}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newQuerySvcClient(cc)
	res, err := cl.GetQuery(context.Background(), req0)
	if err != nil {
		return
	}

	q := res.Query

	if c.String("name") != "" {
		q.Name = c.String("name")
	}

	if c.String("descr") != "" {
		q.Name = c.String("descr")
	}

	if c.Bool("random") {
		q.Random = true
	}

	if c.Bool("no-random") {
		q.Random = false
	}

	if c.Int("rating") > 0 {
		q.Rating = int32(c.Int("rating"))
	}

	if c.Int("limit") > 0 {
		q.Limit = int32(c.Int("limit"))
	}

	if c.String("params") != "" {
		q.Params = c.String("params")
	}

	ts := c.Int64("from")
	if ts > 0 {
		q.From = timestamppb.New(time.Unix(0, ts))
	}

	ts = c.Int64("to")
	if ts > 0 {
		q.To = timestamppb.New(time.Unix(0, ts))
	}

	if len(c.Int64Slice("collection-id")) > 0 {
		q.CollectionIds = c.Int64Slice("collection-id")
	}
	req := &m3uetcpb.UpdateQueryRequest{
		Query: q,
	}

	_, err = cl.UpdateQuery(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}

func queryByAction(c *cli.Context) (err error) {
	rest := c.Args().Slice()

	var from, to *timestamppb.Timestamp

	ts := c.Int64("from")
	if ts > 0 {
		from = timestamppb.New(time.Unix(0, ts))
	}

	ts = c.Int64("to")
	if ts > 0 {
		to = timestamppb.New(time.Unix(0, ts))
	}

	q := &m3uetcpb.Query{
		Name:        c.String("persist-as"),
		Description: c.String("descr"),
		Random:      c.Bool("random"),
		Rating:      int32(c.Int("rating")),
		Limit:       int32(c.Int("limit")),
		Params:      strings.Join(rest, " "),
		From:        from,
		To:          to,
	}
	req := &m3uetcpb.QueryByRequest{
		Query: q,
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newQuerySvcClient(cc)
	res, err := cl.QueryBy(context.Background(), req)
	if err != nil {
		return
	}

	if c.Bool("play") {
		if err = playTracks(cc, res.Tracks, c.Bool("force")); err != nil {
			return
		}
		fmt.Printf("Tracks added to playback!\n")
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

	tbl := table.New("#", "ID", "Title", "Artist", "Album")
	for i, t := range res.Tracks {
		artist := t.Artist
		if artist == "" {
			artist = t.Albumartist
		}
		tbl.AddRow(i+1, t.Id, t.Title, artist, t.Album)
	}
	tbl.Print()
	return
}

func queryInPlaylistAction(c *cli.Context) (err error) {
	var id int64
	if id, err = mustParseSingleID(c); err != nil {
		return
	}

	playlistID := c.Int("pl")

	req := &m3uetcpb.QueryInPlaylistRequest{
		Id:         id,
		PlaylistId: int64(playlistID),
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newQuerySvcClient(cc)
	res, err := cl.QueryInPlaylist(context.Background(), req)
	if err != nil {
		return
	}

	if c.Bool("play") {
		err = playTracksFromPlaylist(cc, res.PlaylistId, c.Bool("force"), c.Int("limit"))
		if err != nil {
			return
		}
		fmt.Printf("Tracks added to playback!\n")
		return
	}

	showPlaylist(c, res.PlaylistId)
	return
}

func queryInQueueAction(c *cli.Context) (err error) {
	var id int64
	if id, err = mustParseSingleID(c); err != nil {
		return
	}

	req := &m3uetcpb.QueryInQueueRequest{
		Id: id,
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newQuerySvcClient(cc)
	_, err = cl.QueryInQueue(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}

func playTracks(cc iClientConn, ts []*m3uetcpb.Track, force bool) (err error) {
	ids := []int64{}
	for _, v := range ts {
		ids = append(ids, v.Id)
	}

	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_PLAY,
		Force:  force,
		Ids:    ids,
	}

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	return
}

func playTracksFromPlaylist(cc iClientConn, playlistID int64, force bool, limit int) (err error) {
	cl := newPlaybarSvcClient(cc)

	req := &m3uetcpb.GetPlaylistRequest{Id: playlistID, Limit: int32(limit)}

	res, err := cl.GetPlaylist(context.Background(), req)
	if err != nil {
		return
	}

	err = playTracks(cc, res.Tracks, force)
	return
}
