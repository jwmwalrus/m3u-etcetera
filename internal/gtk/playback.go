package gtkui

import (
	"context"
	"net/url"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
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
	subscribed                            = false
	status                                playbackData
	playbackID, trackID                   int64
	location, title, artist, album, extra string
)

func onControlClicked(btn *gtk.ToolButton, action m3uetcpb.PlaybackAction) {
	var cc *grpc.ClientConn
	var err error
	opts := alive.GetGrpcDialOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	req := &m3uetcpb.ExecutePlaybackActionRequest{Action: action}
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	if err != nil {
		log.Error(err)
		return
	}
}

func setupPlayback(signals *map[string]interface{}) (err error) {
	(*signals)["on_control_prev_clicked"] = func(btn *gtk.ToolButton) {
		go onControlClicked(btn, m3uetcpb.PlaybackAction_PB_PREVIOUS)
	}
	(*signals)["on_control_play_clicked"] = func(btn *gtk.ToolButton) {
		action := m3uetcpb.PlaybackAction_PB_PLAY
		iconName := "media-playback-pause"

		name := btn.GetIconName()
		if name == "media-playback-pause" {
			action = m3uetcpb.PlaybackAction_PB_PAUSE
			iconName = "media-playback-start"
		}
		btn, err := builder.GetToolButton("control_play")
		onerror.Warn(err)
		if btn != nil {
			btn.SetIconName(iconName)
		}
		go onControlClicked(btn, action)
	}
	(*signals)["on_control_stop_clicked"] = func(btn *gtk.ToolButton) {
		go onControlClicked(btn, m3uetcpb.PlaybackAction_PB_STOP)
	}
	(*signals)["on_control_next_clicked"] = func(btn *gtk.ToolButton) {
		go onControlClicked(btn, m3uetcpb.PlaybackAction_PB_NEXT)
	}
	return
}

func subscribeToPlayback() {
	subscribed = true
	defer func() { subscribed = false }()

	var cc *grpc.ClientConn
	var err error
	opts := alive.GetGrpcDialOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	stream, err := cl.SubscribeToPlayback(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Error(err)
			break
		}

		status.mu.Lock()
		status.data = res
		status.mu.Unlock()
		glib.IdleAdd(updatePlayback)
	}
}

func updatePlayback() bool {
	status.mu.Lock()
	playing := status.data.Playing
	if playing {
		playbackID = status.data.Playback.Id
		location = status.data.Playback.Location
		trackID = status.data.Track.Id
		title = status.data.Track.Title
		artist = status.data.Track.Artist
		album = status.data.Track.Album
	} else {
		if location != "" {
			playbackID, trackID = 0, 0
			location = ""
			title, artist, album = "", "", ""
		}
	}
	status.mu.Unlock()

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

	btn, err := builder.GetToolButton("control_play")
	onerror.Warn(err)
	if btn != nil {
		if playing {
			btn.SetIconName("media-playback-pause")
		} else {
			btn.SetIconName("media-playback-start")
		}
	}
	return false
}
