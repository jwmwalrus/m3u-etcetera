package gtkui

import (
	"context"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func setupPlayback(signals *map[string]interface{}) (err error) {
	log.Info("Setting up playback")

	(*signals)["on_control_prev_clicked"] = func(btn *gtk.ToolButton) {
		go onControlClicked(btn, m3uetcpb.PlaybackAction_PB_PREVIOUS)
	}
	(*signals)["on_control_play_clicked"] = func(btn *gtk.ToolButton) {
		action := m3uetcpb.PlaybackAction_PB_PLAY

		name := btn.GetIconName()
		if name == "media-playback-pause" {
			action = m3uetcpb.PlaybackAction_PB_PAUSE
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

func onControlClicked(btn *gtk.ToolButton, action m3uetcpb.PlaybackAction) {
	log.WithField("action", action.String()).
		Info("ToolButton clicked")

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
	req := &m3uetcpb.ExecutePlaybackActionRequest{Action: action}
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	if err != nil {
		log.Error(err)
		return
	}
}
