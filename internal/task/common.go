package task

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/urfave/cli/v2"
)

func mustNotParseExtraArgs(c *cli.Context) (err error) {
	rest := c.Args().Slice()
	if len(rest) > 0 {
		err = errors.New("Too many values in command")
		return
	}
	return
}
func mustParseSingleID(c *cli.Context) (id int64, err error) {
	rest := c.Args().Slice()
	if len(rest) != 1 {
		err = errors.New("I need one ID")
		return
	}

	id, err = strconv.ParseInt(rest[0], 10, 64)
	if err != nil {
		return
	}
	if id < 1 {
		err = errors.New("I need one ID greater than zero")
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
		var u *url.URL
		if u, err = url.Parse(v); err != nil {
			return
		}
		if u.Scheme == "" {
			u.Scheme = "file"
		}
		parsed = append(parsed, u.String())
	}
	return
}
