package task

import (
	"context"
	"fmt"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/status"
)

var (
	newPerspectiveSvcClient = m3uetcpb.NewPerspectiveSvcClient
)

// Perspective perspective task.
func Perspective() *cli.Command {
	return &cli.Command{
		Name:        "perspective",
		Aliases:     []string{"persp"},
		Category:    "Control",
		Usage:       "Gets or sets the active perspective",
		Description: "Control the application's perspective according with the given subcommand. If no subcommand is given, display the active perspective.",
		Commands: []*cli.Command{
			{
				Name:        "activate",
				Aliases:     []string{"a"},
				Usage:       "Activates perspective",
				ArgsUsage:   "PERSPECTIVE",
				Description: "Activates the given PERSPECTIVE (music|radio|podcasts|audiobooks).",
				Action:      perspectiveActivateAction,
			},
		},
		Before: checkServerStatus,
		Action: perspectiveAction,
	}
}

func perspectiveAction(ctx context.Context, c *cli.Command) (err error) {
	if err = mustNotParseExtraArgs(c); err != nil {
		return
	}

	cc, err := getClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := newPerspectiveSvcClient(cc)
	res, err := cl.GetActivePerspective(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		return
	}

	fmt.Printf("Active: %v\n", res.Perspective.String())

	return
}

func perspectiveActivateAction(ctx context.Context, c *cli.Command) error {
	rest := c.Args().Slice()
	if len(rest) < 1 {
		return fmt.Errorf("I need one PERSPECTIVE to activate")
	}
	if len(rest) > 1 {
		return fmt.Errorf("Too many values in command")
	}

	req := &m3uetcpb.SetActivePerspectiveRequest{
		Perspective: getPerspectiveFromString(rest[0]),
	}

	cc, err := getClientConn()
	if err != nil {
		return nil
	}
	defer cc.Close()

	cl := newPerspectiveSvcClient(cc)
	_, err = cl.SetActivePerspective(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf(s.Message())
	}

	fmt.Printf("OK\n")
	return nil
}
