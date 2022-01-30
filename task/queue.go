package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

// Queue queue task
func Queue() *cli.Command {
	return &cli.Command{
		Name:        "queue",
		Aliases:     []string{"q"},
		Category:    "Control",
		Usage:       "Process the queue task",
		UsageText:   "queue [subtask] ...",
		Description: "Control the application's queue according to the given subcommand. When no subcommand is given, display current queue",
		Subcommands: []*cli.Command{
			{
				Name:    "append",
				Aliases: []string{"app", "add"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "Use IDs instead of locations",
					},
				},
				Usage:       "queue append [location ... | --ids id ...]",
				Description: "Append to queue",
				Action:      queueCreateAction,
			},
			{
				Name:        "clear",
				Usage:       "queue clear",
				Description: "Clear queue",
				Action:      queueDestroyAction,
			},
			{
				Name:    "delete",
				Aliases: []string{"del", "remove", "rem"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "pos",
						Aliases:  []string{"p"},
						Usage:    "delete the (1-based) `POSITION` given by `POS`",
						Required: true,
					},
				},
				Usage:       "queue delete --pos POSITION",
				Description: "Delete `POSITION` in queue",
				Action:      queueDestroyAction,
			},
			{
				Name:    "insert",
				Aliases: []string{"ins"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "pos",
						Aliases:     []string{"p"},
						Usage:       "insert at the (1-based) `POSITION`",
						Value:       1,
						DefaultText: "1",
					},
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "Use IDs instead of locations",
					},
				},
				Usage:       "queue insert --pos POSITION [location ... | --ids id ...]",
				Description: "Insert into queue at the queven position",
				Action:      queueCreateAction,
			},
			{
				Name:    "prepend",
				Aliases: []string{"prep", "top"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "Use IDs instead of locations",
					},
				},
				Usage:       "queue prepend [location ... | --ids id ...]",
				Description: "Prepend to queue",
				Action:      queueCreateAction,
			},
			{
				Name:    "move",
				Aliases: []string{"mv"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "perspective",
						Aliases: []string{"persp"},
						Usage:   "Applies to perspective",
						Value:   "music",
					},
					&cli.Int64Flag{
						Name:     "from-pos",
						Aliases:  []string{"from"},
						Usage:    "Move this `POSITION`",
						Required: true,
					},
					&cli.Int64Flag{
						Name:     "to-pos",
						Aliases:  []string{"to"},
						Usage:    "Move to this `POSITION`",
						Required: true,
					},
				},
				Usage:       "queue move [--flags ...]",
				Description: "Move track from one position to another",
				Action:      queueMoveAction,
			},
		},
		Before: checkServerStatus,
		Action: queueAction,
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
				Usage: "Limit output count",
				Value: 0,
			},
		},
	}
}

func queueAction(c *cli.Context) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	req := &m3uetcpb.GetQueueRequest{
		Perspective: getPerspective(c),
	}

	if c.Int("limit") > 0 {
		req.Limit = int32(c.Int("limit"))
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	res, err := cl.GetQueue(context.Background(), req)
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

	tbl := table.New("Position", "Track|Location")
	for _, qt := range res.QueueTracks {
		if qt.TrackId > 0 {
			s := ""
			for _, t := range res.Tracks {
				if t.Id != qt.TrackId {
					continue
				}
				s = fmt.Sprintf("%s -- by: %v (from: %v)", t.Title, t.Artist, t.Album)
				break
			}
			if s != "" {
				tbl.AddRow(qt.Position, s)
				continue
			}
		}
		un, _ := url.QueryUnescape(qt.Location)
		if un == "" {
			un = qt.Location
		}
		tbl.AddRow(qt.Position, un)
	}
	tbl.Print()

	return
}

func queueCreateAction(c *cli.Context) (err error) {
	const actionPrefix = "Q_"

	rest := c.Args().Slice()
	if len(rest) < 1 {
		err = fmt.Errorf("I need a list of locations or IDs")
		return
	}

	if c.Command.Name == "insert" && c.Int("pos") < 1 {
		err = fmt.Errorf("I need a position to insert")
		return
	}

	action := m3uetcpb.QueueAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]
	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: m3uetcpb.QueueAction(action),
	}

	if c.Bool("ids") {
		if req.Ids, err = parseIDs(rest); err != nil {
			return
		}
	} else {
		if req.Locations, err = parseLocations(rest); err != nil {
			return
		}
	}

	if req.Action == m3uetcpb.QueueAction_Q_INSERT {
		req.Position = int32(c.Int("pos"))
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.ExecuteQueueAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	fmt.Printf("OK\n")
	return
}

func queueDestroyAction(c *cli.Context) (err error) {
	const actionPrefix = "Q_"

	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	if c.Command.Name == "delete" && c.Int("pos") < 1 {
		err = fmt.Errorf("I need a position to delete")
		return
	}

	action := m3uetcpb.QueueAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]
	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: m3uetcpb.QueueAction(action),
	}

	if req.Action == m3uetcpb.QueueAction_Q_DELETE {
		req.Position = int32(c.Int("pos"))
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.ExecuteQueueAction(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}

func queueMoveAction(c *cli.Context) (err error) {
	const actionPrefix = "Q_"

	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	action := m3uetcpb.QueueAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]
	req := &m3uetcpb.ExecuteQueueActionRequest{
		Perspective:  getPerspective(c),
		Action:       m3uetcpb.QueueAction(action),
		FromPosition: int32(c.Int("from-pos")),
		Position:     int32(c.Int("to-pos")),
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.ExecuteQueueAction(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}
