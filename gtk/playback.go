package gtkui

import (
	"log/slog"

	"github.com/diamondburned/gotk4/pkg/gdk/v3"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
)

func setupPlayback(signals *builder.Signals) (err error) {
	slog.Info("Setting up playback")

	(*signals).AddDetail(
		"control_prev",
		"clicked",
		func(btn *gtk.ToolButton) {
			go onControlClicked(btn, m3uetcpb.PlaybackAction_PB_PREVIOUS)
		},
	)
	(*signals).AddDetail(
		"control_play",
		"clicked",
		func(btn *gtk.ToolButton) {
			action := m3uetcpb.PlaybackAction_PB_PLAY

			name := btn.IconName()
			if name == "media-playback-pause" {
				action = m3uetcpb.PlaybackAction_PB_PAUSE
			}
			go onControlClicked(btn, action)
		},
	)
	(*signals).AddDetail(
		"control_stop",
		"clicked",
		func(btn *gtk.ToolButton) {
			go onControlClicked(btn, m3uetcpb.PlaybackAction_PB_STOP)
		},
	)
	(*signals).AddDetail(
		"control_next",
		"clicked",
		func(btn *gtk.ToolButton) {
			go onControlClicked(btn, m3uetcpb.PlaybackAction_PB_NEXT)
		},
	)
	(*signals).AddDetail(
		"progress_eb",
		"button-press-event",
		onProgressBarClicked,
	)
	return
}

func onControlClicked(btn *gtk.ToolButton, action m3uetcpb.PlaybackAction) {
	logw := slog.With("action", action.String())
	logw.Info("ToolButton clicked fot playback action")

	req := &m3uetcpb.ExecutePlaybackActionRequest{Action: action}
	if err := dialer.ExecutePlaybackAction(req); err != nil {
		logw.Error("Failed to execute playback action", "error", err)
	}
}

// onProgressBarClicked is the signal handler for the button-press-event on
// the event-box that wraps the progress bar.
func onProgressBarClicked(eb *gtk.EventBox, event *gdk.Event) {
	_, _, duration, status := store.PbData.GetCurrentPlayback()

	if !status["is-streaming"] {
		return
	}

	btn := event.AsButton()
	x := btn.X()
	width := eb.Widget.AllocatedWidth()
	seek := int64(x * float64(duration) / float64(width))

	go func() {
		req := &m3uetcpb.ExecutePlaybackActionRequest{
			Action: m3uetcpb.PlaybackAction_PB_SEEK,
			Seek:   seek,
		}
		onerror.Log(dialer.ExecutePlaybackAction(req))
	}()
}
