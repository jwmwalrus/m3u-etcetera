package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
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
				Name:        "append",
				Usage:       "queue append location ...",
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
				Name: "delete",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "position",
						Aliases:  []string{"p"},
						Usage:    "delete the (1-based) position given by `POS`",
						Required: true,
					},
				},
				Usage:       "queue delete --position POS",
				Description: "Delete position in queue",
				Action:      queueDestroyAction,
			},
			{
				Name: "insert",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "position",
						Aliases:     []string{"p"},
						Usage:       "insert at the (1-based) position given by `POS`",
						Value:       1,
						DefaultText: "1",
					},
				},
				Usage:       "queue insert --position POS location|payload ...",
				Description: "Insert into queue at the queven position",
				Action:      queueCreateAction,
			},
			{
				Name:        "preppend",
				Usage:       "queue preppend location|payload ...",
				Description: "Preppend to queue",
				Action:      queueCreateAction,
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
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			// TODO: complete
			fmt.Fprintf(c.App.Writer, "for shame\n")
			return err
		},
	}
}

func queueAction(c *cli.Context) (err error) {
	opts := getGrpcOpts()

	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts...)
	if err != nil {
		return
	}
	defer cc.Close()

	req := &m3uetcpb.GetQueueRequest{}

	persp := strings.ToLower(c.String("perspective"))
	if strings.HasPrefix("radio", persp) {
		req.Perspective = m3uetcpb.Perspective_RADIO
	} else if strings.HasPrefix("podcasts", persp) {
		req.Perspective = m3uetcpb.Perspective_PODCASTS
	} else if strings.HasPrefix("audiobooks", persp) {
		req.Perspective = m3uetcpb.Perspective_AUDIOBOOKS
	} else {
		req.Perspective = m3uetcpb.Perspective_MUSIC
	}

	if c.Int("limit") > 0 {
		req.Limit = int32(c.Int("limit"))
	}

	cl := m3uetcpb.NewQueueSvcClient(cc)
	res, err := cl.GetQueue(context.Background(), req)
	if err != nil {
		return
	}

	if c.Bool("json") {
		var bv []byte
		bv, err = json.MarshalIndent(res.QueueTracks, "", "  ")
		if err != nil {
			return
		}
		fmt.Printf("\n%v\n", string(bv))
	} else {
		tbl := table.New("Position", "Location")
		for _, qt := range res.QueueTracks {
			tbl.AddRow(qt.Position, qt.Location)
		}
		tbl.Print()
	}

	return
}

func queueCreateAction(c *cli.Context) (err error) {
	const actionPrefix = "Q_"

	rest := c.Args().Slice()
	if len(rest) < 1 {
		err = errors.New("I need a list of locations")
		return
	}

	opts := getGrpcOpts()

	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts...)
	if err != nil {
		return
	}
	defer cc.Close()

	action := m3uetcpb.QueueAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]
	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: m3uetcpb.QueueAction(action),
	}

	req.Locations, err = parseLocations(rest)

	if req.Action == m3uetcpb.QueueAction_Q_INSERT {
		req.Position = int32(c.Int("position"))
	}

	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.ExecuteQueueAction(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}

func queueDestroyAction(c *cli.Context) (err error) {
	const actionPrefix = "Q_"

	rest := c.Args().Slice()
	if len(rest) > 0 {
		err = errors.New("Too many values in command")
		return
	}

	opts := getGrpcOpts()

	auth := base.Conf.Server.GetAuthority()
	cc, err := grpc.Dial(auth, opts...)
	if err != nil {
		return
	}
	defer cc.Close()

	action := m3uetcpb.QueueAction_value[strings.ToUpper(actionPrefix+c.Command.Name)]
	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: m3uetcpb.QueueAction(action),
	}

	if req.Action == m3uetcpb.QueueAction_Q_DELETE {
		req.Position = int32(c.Int("position"))
	}

	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.ExecuteQueueAction(context.Background(), req)
	if err != nil {
		return
	}

	fmt.Printf("OK\n")
	return
}
