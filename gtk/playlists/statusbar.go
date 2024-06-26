package playlists

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
)

var (
	statusBar       *gtk.Statusbar
	statusBarDigest uint
)

func setupStatusbar() (err error) {
	slog.Info("Setting up status bar")

	statusBar, err = builder.GetStatusBar("status_bar")
	if err != nil {
		return
	}

	statusBarDigest = statusBar.ContextID("Active Digest")
	statusBar.Push(statusBarDigest, "Ready")
	return
}

// UpdateStatusBar updates the status bar with a message suitable to the given context.
func UpdateStatusBar(context uint) {
	time.Sleep(1 * time.Second)
	switch context {
	case statusBarDigest:
		glib.IdleAdd(pushToStatusBarDigest)
	default:
	}
}

func pushToStatusBarDigest() bool {
	p := store.GetActivePerspective()
	id := GetFocused(p)
	var showing, duration int64
	if id > 0 {
		showing = store.BData.PlaylistTracksCount(id)
		pl := store.BData.GetOpenPlaylist(id)
		duration = pl.Duration
	} else {
		showing = store.QData.QueueTracksCount(p)
		dig := store.QData.QueueDigest(p)
		duration = dig.Duration
	}

	nano := time.Duration(duration) * time.Nanosecond
	collTracks := store.CData.TracksTotalCount()
	msg := fmt.Sprintf("%v showing (%v), %v in collections", showing, nano.Truncate(time.Second), collTracks)

	statusBar.RemoveAll(statusBarDigest)
	statusBar.Push(statusBarDigest, msg)

	return false
}
