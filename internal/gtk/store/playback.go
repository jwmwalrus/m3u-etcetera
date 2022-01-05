package store

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type playbackData struct {
	mu   sync.Mutex
	data *m3uetcpb.SubscribeToPlaybackResponse
}

var (
	pbres                                 playbackData
	playbackID, trackID                   int64
	location, title, artist, album, extra string
)

func subscribeToPlayback() {
	log.Info("Subscribing to playback")

	defer wgplayback.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	stream, err := cl.SubscribeToPlayback(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Errorf("Error subscribing to playback: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		pbres.mu.Lock()
		pbres.data = res
		pbres.mu.Unlock()
		glib.IdleAdd(updatePlayback)
		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromPlayback() {
	log.Info("Unsubscribing from playback")

	pbres.mu.Lock()
	id := pbres.data.SubscriptionId
	pbres.mu.Unlock()

	var cc *grpc.ClientConn
	var err error
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.UnsubscribeFromPlayback(
		context.Background(),
		&m3uetcpb.UnsubscribeFromPlaybackRequest{
			SubscriptionId: id,
		},
	)
	if err != nil {
		log.Error(err)
		return
	}
}

func updatePlayback() bool {
	log.Debug("Updating playback")

	pbres.mu.Lock()
	iconName := "media-playback-pause"
	if pbres.data.IsPaused {
		iconName = "media-playback-start"
	}
	btn, err := builder.GetToolButton("control_play")
	onerror.Warn(err)
	if btn != nil {
		btn.SetIconName(iconName)
	}

	var duration, position int64
	if pbres.data.IsStreaming {
		playbackID = pbres.data.Playback.Id
		location = pbres.data.Playback.Location
		trackID = pbres.data.Playback.TrackId

		title = pbres.data.Track.Title
		artist = pbres.data.Track.Artist
		album = pbres.data.Track.Album
		duration = pbres.data.Track.Duration
		position = pbres.data.Playback.Skip
	} else {
		playbackID, trackID = 0, 0
		location = ""
		title, artist, album = "", "", ""
	}
	pbres.mu.Unlock()

	prog, err := builder.GetProgressBar("progress")
	if duration > 0 {
		prog.SetFraction(float64(position) / float64(duration))
		prog.SetText(
			fmt.Sprintf(
				"%v / %v",
				time.Duration(position)*time.Nanosecond,
				time.Duration(duration)*time.Nanosecond,
			),
		)
	} else {
		prog.SetFraction(float64(0))
		prog.SetText("")
	}

	un, err := url.QueryUnescape(location)
	if err != nil {
		location = un
	}

	ltitle, lartist, lsource := title, artist, location
	if title == "" {
		ltitle = "Not Playing"
	}
	if artist != "" {
		lartist = "by " + artist
	}
	if album != "" {
		lsource = "from " + album
	}
	err = builder.SetTextView("playback_title", ltitle)
	onerror.Warn(err)
	err = builder.SetTextView("playback_artist", lartist)
	onerror.Warn(err)
	err = builder.SetTextView("playback_source", lsource)
	onerror.Warn(err)

	return false
}
