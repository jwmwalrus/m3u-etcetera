package gtkui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
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
	(*signals)["on_progress_eb_button_press_event"] = onProgressBarClicked
	return
}

func onControlClicked(btn *gtk.ToolButton, action m3uetcpb.PlaybackAction) {
	entry := log.WithField("action", action.String())
	entry.Info("ToolButton clicked fot playback action")

	req := &m3uetcpb.ExecutePlaybackActionRequest{Action: action}
	if err := dialer.ExecutePlaybackAction(req); err != nil {
		entry.Error(err)
	}
}

// onProgressBarClicked is the signal handler for the button-press-event on
// the event-box that wraps the progress bar
func onProgressBarClicked(eb *gtk.EventBox, event *gdk.Event) {
	_, _, duration, status := store.PbData.GetCurrentPlayback()

	if !status["is-streaming"] {
		return
	}

	btn := gdk.EventButtonNewFromEvent(event)
	x, _ := btn.MotionVal()
	width := eb.Widget.GetAllocatedWidth()
	seek := int64(x * float64(duration) / float64(width))

	go func() {
		req := &m3uetcpb.ExecutePlaybackActionRequest{
			Action: m3uetcpb.PlaybackAction_PB_SEEK,
			Seek:   seek,
		}
		onerror.Log(dialer.ExecutePlaybackAction(req))
	}()
}
