package playlists

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
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

	statusBarDigest = statusBar.GetContextId("Active Digest")
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
		showing = store.BData.GetPlaylistTracksCount(id)
		pl := store.BData.GetOpenPlaylist(id)
		duration = pl.Duration
	} else {
		showing = store.QData.GetQueueTracksCount(p)
		dig := store.QData.GetQueueDigest(p)
		duration = dig.Duration
	}

	nano := time.Duration(duration) * time.Nanosecond
	collTracks := store.CData.GetTracksTotalCount()
	msg := fmt.Sprintf("%v showing (%v), %v in collections", showing, nano.Truncate(time.Second), collTracks)

	statusBar.RemoveAll(statusBarDigest)
	statusBar.Push(statusBarDigest, msg)

	return false
}
