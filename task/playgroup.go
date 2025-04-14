package task

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/status"
)

// Playgroup playgroup task.
func Playgroup() *cli.Command {
	return &cli.Command{
		Name:        "playgroup",
		Aliases:     []string{"pg"},
		Category:    "Organization",
		Usage:       "Performs playlist-group-related actions",
		Description: "Perform the playlist-group-related action. When no subcommand is given, display all the playlist groups.",
		Before:      checkServerStatus,
		Action:      playgroupAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "output JSON",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "info",
				Aliases:     []string{"i"},
				Usage:       "Shows playlist group info",
				ArgsUsage:   "ID",
				Description: "Show information for the playlist group identified by the given `ID`.",
				Action:      playgroupInfoAction,
			},
			{
				Name:        "create",
				Usage:       "Creates playlist group",
				Description: "Creates a new playlist group, according to the given options.",
				Action:      playgroupExecuteAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "playlist group `NAME`",
					},
					&cli.StringFlag{
						Name:    "descr",
						Aliases: []string{"d"},
						Usage:   "playlist group `DESCRIPTION`",
					},
				},
			},
			{
				Name:        "update",
				Aliases:     []string{"upd"},
				Usage:       "Updates playlist group",
				ArgsUsage:   "ID",
				Description: "Update the playlist group identified by `ID`, according to the given options.",
				Action:      playgroupExecuteAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "playlist group `NAME`",
					},
					&cli.StringFlag{
						Name:    "descr",
						Aliases: []string{"d"},
						Usage:   "playlist group `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:    "reset-descr",
						Aliases: []string{"rd"},
						Usage:   "reset the description to an empty string",
					},
				},
			},
			{
				Name:        "destroy",
				Aliases:     []string{"del", "delete"},
				Usage:       "Deletes playlist group",
				ArgsUsage:   "ID",
				Description: "Delete playlist group identified by the given `ID`",
				Action:      playgroupExecuteAction,
			},
		},
	}
}

func playgroupAction(ctx context.Context, c *cli.Command) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybarSvcClient(cc)

	req := &m3uetcpb.GetAllPlaylistGroupsRequest{
		Perspective: getPerspective(c),
		Limit:       int32(c.Int("limit")),
	}

	res, err := cl.GetAllPlaylistGroups(context.Background(), req)
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

	tbl := table.New("ID", "Name", "Perspective")
	for _, pg := range res.PlaylistGroups {
		tbl.AddRow(pg.Id, pg.Name, pg.Perspective)
	}
	tbl.Print()
	return
}

func playgroupInfoAction(ctx context.Context, c *cli.Command) (err error) {
	id, err := mustParseSingleID(c)
	if err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybarSvcClient(cc)

	req := &m3uetcpb.GetPlaylistGroupRequest{Id: id}

	res, err := cl.GetPlaylistGroup(context.Background(), req)
	if err != nil {
		return
	}

	var bv []byte
	bv, err = json.MarshalIndent(res, "", "  ")
	if err != nil {
		return
	}
	fmt.Printf("\n%v\n", string(bv))
	return
}

func playgroupExecuteAction(ctx context.Context, c *cli.Command) (err error) {
	const actionPrefix = "PG_"

	var id int64
	if c.Name == "create" {
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

	action := m3uetcpb.PlaylistGroupAction_value[strings.ToUpper(actionPrefix+c.Name)]

	req := &m3uetcpb.ExecutePlaylistGroupActionRequest{
		Action:      m3uetcpb.PlaylistGroupAction(action),
		Id:          id,
		Name:        c.String("name"),
		Description: c.String("descr"),
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybarSvcClient(cc)
	res, err := cl.ExecutePlaylistGroupAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK, ID: %v\n", res.Id)
	return
}
