package task

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

var (
	newCollectionSvcClient = m3uetcpb.NewCollectionSvcClient
)

// Collection defines the collection-related tasks.
func Collection() *cli.Command {
	return &cli.Command{
		Name:        "collection",
		Aliases:     []string{"coll"},
		Category:    "Organization",
		Usage:       "Handles collections",
		Description: "Processes collection-related subcommands. When no subcommand is given, displays all collections currently defined.",
		Before:      checkServerStatus,
		Action:      collectionAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "output JSON",
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:        "info",
				Aliases:     []string{"i"},
				Usage:       "Shows collection(s) info",
				ArgsUsage:   "ID",
				Description: "Show all the fields/properties for the collection defined by the given `ID`.",
				Action:      collectionInfoAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "output JSON",
					},
				},
			},
			{
				Name:        "add",
				Usage:       "Adds a collection",
				ArgsUsage:   "NAME LOCATION",
				Description: "Adds a collection with the given `NAME`. If collection is not remote and not disabled, it will start scanning the `LOCATION` immediately.",
				Action:      collectionAddAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "disabled",
						Usage: "disable collection",
					},
					&cli.BoolFlag{
						Name:    "remote",
						Aliases: []string{"r"},
						Usage:   "collection is remote",
					},
					&cli.StringFlag{
						Name:  "descr",
						Usage: "collection's `DESCRIPTION`",
					},
					&cli.StringFlag{
						Name:  "persp",
						Usage: "applies to `PERSPECTIVE` (music|audiobooks)",
						Value: "music",
					},
				},
			},
			{
				Name:        "remove",
				Aliases:     []string{"rem"},
				Usage:       "Removes a collection",
				ArgsUsage:   "ID",
				Description: "Remove the collection identified by the given `ID`.",
				Action:      collectionRemoveAction,
			},
			{
				Name:        "update",
				Aliases:     []string{"upd"},
				Usage:       "Updates a collection",
				ArgsUsage:   "ID",
				Description: "Updates values in the collection identified by `ID`, according to the  the given options.",
				Action:      collectionUpdateAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "enable",
						Aliases: []string{"en"},
						Usage:   "enable collection",
					},
					&cli.BoolFlag{
						Name:    "disable",
						Aliases: []string{"dis"},
						Usage:   "disable collection",
					},
					&cli.BoolFlag{
						Name:    "local",
						Aliases: []string{"l"},
						Usage:   "collection is local",
					},
					&cli.BoolFlag{
						Name:    "remote",
						Aliases: []string{"r"},
						Usage:   "collection is remote",
					},
					&cli.StringFlag{
						Name:  "name",
						Usage: "rename collection to the given `NAME` (shall be unique)",
					},
					&cli.StringFlag{
						Name:  "descr",
						Usage: "change collection's `DESCRIPTION` for the given one",
					},
				},
			},
			{
				Name:        "scan",
				Usage:       "Scans a collection",
				ArgsUsage:   "ID",
				Description: "Scan collection with the given `ID` for new and deleted tracks.",
				Action:      collectionScanAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "update-tags",
						Usage: "update the tags of tracks already in the collection",
					},
				},
			},
			{
				Name:        "discover",
				Aliases:     []string{"dis"},
				Usage:       "Discovers collections",
				Description: "Scan all collections for new tracks.",
				Action:      collectionDiscoverActiion,
			},
		},
	}
}

func collectionAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newCollectionSvcClient(cc)
	res, err := cl.GetAllCollections(context.Background(), &m3uetcpb.Empty{})
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

	tbl := table.New("ID", "Name", "Disabled", "Remote", "Tracks", "Location")
	for _, i := range res.Collections {
		var st string
		if i.Scanned != 100 {
			st = strconv.Itoa(int(i.Scanned)) + "%"
		} else {
			st = strconv.FormatInt(i.Tracks, 10)
		}
		tbl.AddRow(i.Id, i.Name, i.Disabled, i.Remote, st, i.Location)
	}
	tbl.Print()

	return
}

func collectionInfoAction(c *cli.Context) (err error) {
	var id int64
	if id, err = mustParseSingleID(c); err != nil {
		return
	}

	req := &m3uetcpb.GetCollectionRequest{Id: id}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newCollectionSvcClient(cc)
	res, err := cl.GetCollection(context.Background(), req)
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

	tbl := table.New("ID", "Name", "Disabled", "Remote", "Tracks", "Location")
	coll := res.Collection
	var st string
	if coll.Scanned != 100 {
		st = strconv.Itoa(int(coll.Scanned)) + "%"
	} else {
		st = strconv.FormatInt(coll.Tracks, 10)
	}
	tbl.AddRow(coll.Id, coll.Name, coll.Disabled, coll.Remote, st, coll.Location)
	tbl.Print()

	return
}

func collectionAddAction(c *cli.Context) (err error) {
	rest := c.Args().Slice()
	if len(rest) != 2 {
		err = fmt.Errorf("I need name and path")
		return
	}

	persp := getPerspective(c)
	if !slices.Contains([]m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC, m3uetcpb.Perspective_AUDIOBOOKS}, persp) {
		err = fmt.Errorf("Invalid perspective provided")
		return
	}

	req := &m3uetcpb.AddCollectionRequest{
		Name:        rest[0],
		Location:    rest[1],
		Disabled:    c.Bool("disabled"),
		Remote:      c.Bool("remote"),
		Description: c.String("descr"),
	}

	req.Location, err = urlstr.PathToURL(rest[1])
	if err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newCollectionSvcClient(cc)
	res, err := cl.AddCollection(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("ID: %v\n", res.Id)
	return
}

func collectionRemoveAction(c *cli.Context) (err error) {
	var id int64
	if id, err = mustParseSingleID(c); err != nil {
		return
	}

	req := &m3uetcpb.RemoveCollectionRequest{Id: id}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newCollectionSvcClient(cc)
	_, err = cl.RemoveCollection(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}

func collectionUpdateAction(c *cli.Context) (err error) {
	var id int64
	if id, err = mustParseSingleID(c); err != nil {
		return
	}

	req := &m3uetcpb.UpdateCollectionRequest{Id: id}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	if c.String("name") != "" {
		req.NewName = c.String("name")
	}
	if c.String("descr") != "" {
		req.NewDescription = c.String("descr")
	}
	if c.Bool("enable") {
		req.Enable = true
	}
	if c.Bool("disable") {
		req.Disable = true
	}
	if c.Bool("local") {
		req.MakeLocal = true
	}
	if c.Bool("remote") {
		req.MakeRemote = true
	}

	cl := newCollectionSvcClient(cc)
	_, err = cl.UpdateCollection(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}

func collectionScanAction(c *cli.Context) (err error) {
	var id int64
	if id, err = mustParseSingleID(c); err != nil {
		return
	}

	req := &m3uetcpb.ScanCollectionRequest{
		Id:         id,
		UpdateTags: c.Bool("update-tags"),
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newCollectionSvcClient(cc)
	_, err = cl.ScanCollection(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}

func collectionDiscoverActiion(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newCollectionSvcClient(cc)
	_, err = cl.DiscoverCollections(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}
