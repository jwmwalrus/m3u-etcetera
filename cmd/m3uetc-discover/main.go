package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/internal/discover"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/pbutils"
)

func main() {

	args := os.Args

	gst.Init(nil)

	var out []string
	var errcount int

	withLocationArgs := false
	for _, path := range args[1:] {
		if path == "-l" {
			withLocationArgs = true
			continue
		}

		location := path
		if !withLocationArgs {
			var err error
			location, err = urlstr.PathToURL(location)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to encode path `%s`: %v\n", path, err)
				errcount++
				continue
			}
		}
		i, err := run(location)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			errcount++
			continue
		}
		bv, err := json.Marshal(i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to marshal JSON for URI `%s`: %v\n", path, err)
			errcount++
			continue
		}

		out = append(out, string(bv))
	}

	for _, o := range out {
		fmt.Fprintln(os.Stdout, o)
	}

	os.Exit(errcount)
}

// run invokes pbutils.DiscoverURI for the given location.
func run(location string) (*discover.Info, error) {
	discoverer, err := pbutils.NewDiscoverer(time.Second * 15)
	if err != nil {
		return nil, fmt.Errorf("failed to create discoverer: %w", err)
	}

	info, err := discoverer.DiscoverURI(location)
	if err != nil {
		return nil, fmt.Errorf("failed to discover URI `%s`: %w", location, err)
	}

	// info.GetTags()
	// info.GetAudioStreams()
	// info.GetContainerStreams()
	// info.GetStreamInfo()
	// info.GetStreamList()
	return &discover.Info{
		Duration: int64(info.GetDuration() * time.Nanosecond),
		Live:     info.GetLive(),
		Seekable: info.GetSeekable(),
		URI:      info.GetURI(),
	}, nil
}
