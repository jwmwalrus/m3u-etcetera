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

// Playgroup playgroup task.
func Playgroup() *cli.Command {
	return &cli.Command{
		Name:        "playgroup",
		Aliases:     []string{"pg"},
		Category:    "Organization",
		Usage:       "Performs playtrack-related actions",
		UsageText:   "playtrack [subcommand] ...",
		Description: "Perform the playtrack-related action on playlist tracks, using the ID given by the --playlist flag, and the subcommand. When no subcommand is given, display all the tracks in the playlist",
		Subcommands: []*cli.Command{
			{
				Name:        "info",
				Aliases:     []string{"i"},
				Usage:       "playgroup info ID",
				Description: "Show information for the playlist group with the given `ID`",
				Action:      playgroupInfoAction,
			},
			{
				Name: "create",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Playlist group name",
					},
					&cli.StringFlag{
						Name:    "descr",
						Aliases: []string{"d"},
						Usage:   "Playlist group `DESCRIPTION`",
					},
				},
				Usage:       "playgroup create [<flags> ...]",
				Description: "Creates a new playlist group with values from the given flags",
				Action:      playgroupExecuteAction,
			},
			{
				Name: "update",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Playlist group name",
					},
					&cli.StringFlag{
						Name:    "descr",
						Aliases: []string{"d"},
						Usage:   "Playlist group `DESCRIPTION`",
					},
					&cli.BoolFlag{
						Name:    "reset-descr",
						Aliases: []string{"rd"},
						Usage:   "Reset the description to an empty string",
					},
				},
				Aliases:     []string{"upd"},
				Usage:       "playgroup update [<flags> ...] ID",
				Description: "Update the playlist group identified by `ID` with values from the given flags",
				Action:      playgroupExecuteAction,
			},
			{
				Name:        "destroy",
				Aliases:     []string{"del", "delete"},
				Usage:       "playgroup destroy ID",
				Description: "Delete playlist group identified by the given `ID`",
				Action:      playgroupExecuteAction,
			},
		},
		Before: checkServerStatus,
		Action: playgroupAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "Output JSON",
			},
		},
	}
}

func playgroupAction(c *cli.Context) (err error) {
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

func playgroupInfoAction(c *cli.Context) (err error) {
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

func playgroupExecuteAction(c *cli.Context) (err error) {
	const actionPrefix = "PG_"

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

	action := m3uetcpb.PlaylistGroupAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]

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
