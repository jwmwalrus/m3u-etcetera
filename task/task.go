package task

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/gear-pieces/middleware"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
)

type iClientConn interface {
	grpc.ClientConnInterface
	Close() error
}

var (
	getClientConn = getClientConnDefault
)

func DefaultAction(ctx context.Context, c *cli.Command) error {
	argsSlice := c.Args().Slice()
	if len(argsSlice) == 0 {
		cli.ShowAppHelpAndExit(c, 1)
		return nil
	}

	var cmd *cli.Command
	for i := range c.Commands {
		if c.Commands[i].Name == "playback" {
			cmd = c.Commands[i]
			break
		}
	}

	args := []string{"playback", "play"}
	if base.Conf.Task.ForceDefaultAction {
		args = append(args, "--force")
	}

	args = append(args, c.Args().Slice()...)

	// FIXME: how does this work now?
	// c2 := cli.NewContext(c.App, nil, c)
	return cmd.Run(context.Background(), args)
}

func getClientConnDefault() (iClientConn, error) {
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	return grpc.Dial(auth, opts...)
}

func getPerspective(c *cli.Command) m3uetcpb.Perspective {
	return getPerspectiveFromString(c.String("persp"))
}

func getPerspectiveFromString(s string) (out m3uetcpb.Perspective) {
	persp := strings.ToLower(s)
	if strings.HasPrefix(persp, "radio") {
		out = m3uetcpb.Perspective_RADIO
	} else if strings.HasPrefix(persp, "podcasts") {
		out = m3uetcpb.Perspective_PODCASTS
	} else if strings.HasPrefix(persp, "audiobooks") {
		out = m3uetcpb.Perspective_AUDIOBOOKS
	} else {
		out = m3uetcpb.Perspective_MUSIC
	}
	return
}

func mustNotParseExtraArgs(c *cli.Command) (err error) {
	rest := c.Args().Slice()
	if len(rest) > 0 {
		err = fmt.Errorf("Too many values in command")
		return
	}
	return
}

func mustParseSingleID(c *cli.Command) (id int64, err error) {
	rest := c.Args().Slice()
	if len(rest) != 1 {
		err = fmt.Errorf("I need one ID")
		return
	}

	id, err = strconv.ParseInt(rest[0], 10, 64)
	if err != nil {
		return
	}
	if id < 1 {
		err = fmt.Errorf("I need one ID greater than zero")
		return
	}
	return
}

func parseIDs(ids []string) (parsed []int64, err error) {
	for _, v := range ids {
		var aux int64
		if aux, err = strconv.ParseInt(v, 10, 64); err != nil {
			return
		}
		if aux < 1 {
			err = fmt.Errorf("Found invalid ID: %v", aux)
			return
		}
		parsed = append(parsed, aux)
	}
	return
}

func parseLocations(locations []string) (parsed []string, err error) {
	for _, v := range locations {
		var u string
		if u, err = urlstr.PathToURL(v); err != nil {
			return
		}
		parsed = append(parsed, u)
	}
	return
}

func parseSeconds(secs string) (parsed int64, err error) {
	var aux float64
	if aux, err = strconv.ParseFloat(secs, 64); err != nil {
		return
	}

	parsed = int64(aux * 1e9)
	return
}
