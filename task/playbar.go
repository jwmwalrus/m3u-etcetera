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

var (
	newPlaybarSvcClient = m3uetcpb.NewPlaybarSvcClient
)

// Playbar playbar task.
func Playbar() *cli.Command {
	return &cli.Command{
		Name:        "playbar",
		Aliases:     []string{"bar"},
		Category:    "Control",
		Usage:       "Controls the playbar",
		Description: "The playbar command controls active playlist and the number of open ones as well.",
		Before:      checkServerStatus,
		Action:      playbarAction,
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
		},
		Subcommands: []*cli.Command{
			{
				Name:        "open",
				Usage:       "Opens playlist",
				ArgsUsage:   "ID",
				Description: "Open the playlist identified by the given `ID`.",
				Action:      playbarExecuteAction,
			},
			{
				Name:        "activate",
				Aliases:     []string{"act"},
				Usage:       "Activates playlist",
				ArgsUsage:   "ID",
				Description: "(Open and) Activate the playlist identified by the given `ID`.",
				Action:      playbarExecuteAction,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "pos",
						Aliases:     []string{"p"},
						Usage:       "activate playlist at the given `POSITION`",
						Value:       1,
						DefaultText: "1",
					},
				},
			},
			{
				Name:        "deactivate",
				Aliases:     []string{"deact"},
				Usage:       "Deactivates playlist",
				ArgsUsage:   "ID",
				Description: "Deactivate the playlist identified by the given `ID`.",
				Action:      playbarExecuteAction,
			},
			{
				Name:        "close",
				Aliases:     []string{"clo"},
				Usage:       "Closes playlist",
				ArgsUsage:   "ID",
				Description: "Close the playlist identified by the given `ID`.",
				Action:      playbarExecuteAction,
			},
		},
	}
}

func playbarAction(c *cli.Context) (err error) {
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

func playbarExecuteAction(c *cli.Context) (err error) {
	const actionPrefix = "BAR_"

	rest := c.Args().Slice()
	if len(rest) != 1 {
		if c.Command.Name == "activate" ||
			c.Command.Name == "deactivate" &&
				len(rest) > 1 {
			err = fmt.Errorf("I need one ID to activate/deactuvate")
			return
		}
		err = fmt.Errorf("I need a list of IDs")
		return
	}

	if c.Command.Name == "activate" && c.Int("pos") < 1 {
		err = fmt.Errorf("I need a position to activate")
		return
	}

	action := m3uetcpb.PlaybarAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]

	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Action: m3uetcpb.PlaybarAction(action),
	}

	if req.Ids, err = parseIDs(rest); err != nil {
		return
	}

	if req.Action == m3uetcpb.PlaybarAction_BAR_ACTIVATE {
		req.Position = int32(c.Int("pos"))
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPlaybarSvcClient(cc)
	_, err = cl.ExecutePlaybarAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK\n")
	return
}
