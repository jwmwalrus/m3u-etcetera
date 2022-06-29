package gtkui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	log "github.com/sirupsen/logrus"
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
	(*signals)["on_progress_eb_button_press_event"] = store.OnProgressBarClicked
	return
}

func onControlClicked(btn *gtk.ToolButton, action m3uetcpb.PlaybackAction) {
	log.WithField("action", action.String()).
		Info("ToolButton clicked fot playback action")

	req := &m3uetcpb.ExecutePlaybackActionRequest{Action: action}
	if err := store.ExecutePlaybackAction(req); err != nil {
		log.Error(err)
	}
}
