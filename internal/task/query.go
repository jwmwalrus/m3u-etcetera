package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// Query defines the query-related tasks
func Query() *cli.Command {
	return &cli.Command{
		Name:        "query",
		Aliases:     []string{"search", "s"},
		Category:    "Control",
		Usage:       "Process the queue task",
		UsageText:   "queue [subtask] ...",
		Description: "Control the application's queue according to the given subcommand. When no subcommand is given, display current queue",
		Subcommands: []*cli.Command{
			{
				Name:   "add",
				Action: queryAddAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "Query `NAME`",
					},
					&cli.StringFlag{
						Name:  "description",
						Usage: "Query `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:  "random",
						Usage: "Query is random",
					},
					&cli.IntFlag{
						Name:  "rating",
						Usage: "Query `RATING`",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "Query `LIMIT`",
					},
					&cli.StringFlag{
						Name:  "params",
						Usage: "Query `PARAMS` for title,artist,album,genre (e.g.: \"title=thing and genre=[sh]ome or genre=some*other\"",
					},
					&cli.Int64Flag{
						Name:  "from",
						Usage: "Query's start `TIMESTAMP` (i.e., from the date the track was issued)",
					},
					&cli.Int64Flag{
						Name:  "to",
						Usage: "Query's end `TIMESTAMP` (i.e., to the date the track was issued)",
					},
					&cli.Int64SliceFlag{
						Name:    "collection-id",
						Aliases: []string{"coll-id"},
						Usage:   "`ID` for the collection bounding the query (can appear more than once in the command)",
					},
				},
			},
			{
				Name:    "info",
				Aliases: []string{"i"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "Output JSON",
					},
				},
				Action: queryInfoAction,
			},
			{
				Name:    "remove",
				Aliases: []string{"rem"},
				Action:  queryRemoveAction,
			},
			{
				Name:    "update",
				Aliases: []string{"upd"},
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "id",
						Usage:    "Query's existing `ID`",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "name",
						Usage: "Query `NAME`",
					},
					&cli.StringFlag{
						Name:  "description",
						Usage: "Query `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:  "random",
						Usage: "Query is random",
					},
					&cli.IntFlag{
						Name:  "rating",
						Usage: "Query `RATING`",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "Query `LIMIT`",
					},
					&cli.StringFlag{
						Name:  "params",
						Usage: "Query `PARAMS` for title,artist,album,genre (e.g.: \"title=thing and genre=[sh]ome or genre=some*other\"",
					},
					&cli.Int64Flag{
						Name:  "from",
						Usage: "Query's start `TIMESTAMP` (i.e., from the date the track was issued)",
					},
					&cli.Int64Flag{
						Name:  "to",
						Usage: "Query's end `TIMESTAMP` (i.e., to the date the track was issued)",
					},
					&cli.Int64SliceFlag{
						Name:    "collection-id",
						Aliases: []string{"coll-id"},
						Usage:   "`ID` for the collection bounding the query (can appear more than once in the command)",
					},
				},
				Action: queryUpdateAction,
			},
			{
				Name:    "tracks",
				Aliases: []string{"t"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "Output JSON",
					},
				},
				Action: queryTracksAction,
			},
			{
				Name: "by",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "persist-as",
						Usage: "Query `NAME`",
					},
					&cli.StringFlag{
						Name:  "description",
						Usage: "Query `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:  "random",
						Usage: "Query is random",
					},
					&cli.IntFlag{
						Name:  "rating",
						Usage: "Query `RATING`",
					},
					&cli.IntFlag{
						Name:  "limit",
						Usage: "Query `LIMIT`",
					},
					&cli.StringFlag{
						Name:  "params",
						Usage: "Query `PARAMS` for title,artist,album,genre (e.g.: \"title=thing and genre=[sh]ome or genre=some*other\"",
					},
					&cli.Int64Flag{
						Name:  "from",
						Usage: "Query's start `TIMESTAMP` (i.e., from the date the track was issued)",
					},
					&cli.Int64Flag{
						Name:  "to",
						Usage: "Query's end `TIMESTAMP` (i.e., to the date the track was issued)",
					},
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "Output JSON",
					},
				},
				Action: queryByAction,
			},
		},
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "doo!",
		BashComplete: func(c *cli.Context) {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "--better\n")
		},
		Before: checkServerStatus,
		Action: queryAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "Output JSON",
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit output count",
				Value: 0,
			},
			&cli.Int64SliceFlag{
				Name:    "collection-id",
				Aliases: []string{"coll-id"},
				Usage:   "Bound to collection ID",
			},
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "for shame\n")
			return err
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

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(getAuthority(), getGrpcOpts()...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
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

	tbl := table.New("ID", "Name", "Params", "Limit", "Bounds")
	for _, q := range res.Queries {
		b := ""
		if len(q.CollectionIds) > 0 {
			b = "C"
		}
		tbl.AddRow(q.Id, q.Name, q.Params, q.Limit, b)
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

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(getAuthority(), getGrpcOpts()...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
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

	tbl := table.New("ID", "Name", "Params", "Limit", "Bounds")
	q := res.Query
	b := ""
	if len(q.CollectionIds) > 0 {
		b = "C"
	}
	tbl.AddRow(q.Id, q.Name, q.Params, q.Limit, b)

	tbl.Print()
	return
}

func queryAddAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	q := &m3uetcpb.Query{
		Name:          c.String("name"),
		Description:   c.String("description"),
		Random:        c.Bool("random"),
		Rating:        int32(c.Int("rating")),
		Limit:         int32(c.Int("limit")),
		Params:        c.String("params"),
		From:          c.Int64("from"),
		To:            c.Int64("from"),
		CollectionIds: c.Int64Slice("collection-id"),
	}
	req := &m3uetcpb.AddQueryRequest{
		Query: q,
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(getAuthority(), getGrpcOpts()...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
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

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(getAuthority(), getGrpcOpts()...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
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
		err = errors.New("I need an ID greater than zero")
		return
	}

	req0 := &m3uetcpb.GetQueryRequest{Id: c.Int64("id")}

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(getAuthority(), getGrpcOpts()...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	res, err := cl.GetQuery(context.Background(), req0)
	if err != nil {
		return
	}

	q := res.Query

	if c.String("name") != "" {
		q.Name = c.String("name")
	}

	if c.String("description") != "" {
		q.Name = c.String("description")
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

	if c.Int64("from") > 0 {
		q.From = c.Int64("from")
	}

	if c.Int64("to") > 0 {
		q.To = c.Int64("to")
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

func queryTracksAction(c *cli.Context) (err error) {
	var id int64
	if id, err = mustParseSingleID(c); err != nil {
		return
	}

	req := &m3uetcpb.ApplyQueryRequest{
		Id: id,
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(getAuthority(), getGrpcOpts()...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	res, err := cl.ApplyQuery(context.Background(), req)
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

func queryByAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	q := &m3uetcpb.Query{
		Name:        c.String("name"),
		Description: c.String("description"),
		Random:      c.Bool("random"),
		Rating:      int32(c.Int("rating")),
		Limit:       int32(c.Int("limit")),
		Params:      c.String("params"),
		From:        c.Int64("from"),
		To:          c.Int64("from"),
	}
	req := &m3uetcpb.QueryByRequest{
		Query: q,
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(getAuthority(), getGrpcOpts()...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	res, err := cl.QueryBy(context.Background(), req)
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
