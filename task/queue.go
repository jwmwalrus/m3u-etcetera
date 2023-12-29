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

var (
	newQueueSvcClient = m3uetcpb.NewQueueSvcClient
)

// Queue queue task.
func Queue() *cli.Command {
	return &cli.Command{
		Name:        "queue",
		Aliases:     []string{"q"},
		Category:    "Control",
		Usage:       "Controls the queue",
		Description: "Control the application's queue according to the given subcommand. When no subcommand is given, display current queue.",
		Before:      checkServerStatus,
		Action:      queueAction,
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
				Usage: "limit output count",
				Value: 0,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:        "append",
				Aliases:     []string{"app", "add"},
				Usage:       "Append to queue",
				ArgsUsage:   "LOCATION|ID ...",
				Description: "Append the given payload (i.e., list of `LOCATION` or `ID`) to the queue.",
				Action:      queueCreateAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "use IDs instead of LOCATIONs",
					},
				},
			},
			{
				Name:        "clear",
				Usage:       "Clears the queue",
				Description: "Clear the queue.",
				Action:      queueDestroyAction,
			},
			{
				Name:        "delete",
				Aliases:     []string{"del", "remove", "rem"},
				Usage:       "Delete from queue",
				Description: "Delete position from the queue.",
				Action:      queueDestroyAction,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "pos",
						Aliases:  []string{"p"},
						Usage:    "delete the given (1-based) `POSITION`",
						Required: true,
					},
				},
			},
			{
				Name:        "insert",
				Aliases:     []string{"ins"},
				Usage:       "Inserts into queue",
				Description: "Insert into queue at the given position.",
				Action:      queueCreateAction,
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
						Usage: "use IDs instead of LOCATIONs",
					},
				},
			},
			{
				Name:        "prepend",
				Aliases:     []string{"prep", "top"},
				Usage:       "Prepends to queue",
				Description: "Prepend to queue (i.e., insert at position 1).",
				Action:      queueCreateAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "ids",
						Usage: "use IDs instead of LOCATIONs",
					},
				},
			},
			{
				Name:        "move",
				Aliases:     []string{"mv"},
				Usage:       "Moves track",
				Description: "Move track from one position to another.",
				Action:      queueMoveAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "persp",
						Usage: "applies to `PERSPECTIVE`",
						Value: "music",
					},
					&cli.Int64Flag{
						Name:     "from-pos",
						Aliases:  []string{"from"},
						Usage:    "move track at this `POSITION`",
						Required: true,
					},
					&cli.Int64Flag{
						Name:     "to-pos",
						Aliases:  []string{"to"},
						Usage:    "move to this `POSITION`",
						Required: true,
					},
				},
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

	cl := newQueueSvcClient(cc)
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
				s = fmt.Sprintf(
					"%s -- by: %v (from: %v)",
					t.Title,
					t.Artist,
					t.Album,
				)
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

	cl := newQueueSvcClient(cc)
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

	cl := newQueueSvcClient(cc)
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

	cl := newQueueSvcClient(cc)
	_, err = cl.ExecuteQueueAction(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}
