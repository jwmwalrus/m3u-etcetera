package task

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func getClientConn() (*grpc.ClientConn, error) {
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	return grpc.Dial(auth, opts...)
}

func getPerspective(c *cli.Context) (p m3uetcpb.Perspective) {
	persp := strings.ToLower(c.String("persp"))
	if strings.HasPrefix(persp, "radio") {
		p = m3uetcpb.Perspective_RADIO
	} else if strings.HasPrefix(persp, "podcasts") {
		p = m3uetcpb.Perspective_PODCASTS
	} else if strings.HasPrefix(persp, "audiobooks") {
		p = m3uetcpb.Perspective_AUDIOBOOKS
	} else {
		p = m3uetcpb.Perspective_MUSIC
	}
	return
}

func mustNotParseExtraArgs(c *cli.Context) (err error) {
	rest := c.Args().Slice()
	if len(rest) > 0 {
		err = fmt.Errorf("Too many values in command")
		return
	}
	return
}

func mustParseSingleID(c *cli.Context) (id int64, err error) {
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
