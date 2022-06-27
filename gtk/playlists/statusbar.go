package playlists

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	log "github.com/sirupsen/logrus"
)

var (
	statusBar       *gtk.Statusbar
	StatusBarDigest uint
)

func setupStatusbar() (err error) {
	log.Info("Setting up status bar")

	statusBar, err = builder.GetStatusBar("status_bar")
	if err != nil {
		return
	}

	StatusBarDigest = statusBar.GetContextId("Active Digest")
	statusBar.Push(StatusBarDigest, "Ready")
	return
}

func UpdateStatusBar(context uint) {
	time.Sleep(1 * time.Second)
	switch context {
	case StatusBarDigest:
		glib.IdleAdd(pushToStatusBarDigest)
	default:
	}
}

func pushToStatusBarDigest() bool {
	p := m3uetcpb.Perspective_MUSIC
	id := GetFocused(p)
	var showing, duration int64
	if id > 0 {
		showing = store.GetPlaylistTracksCount(id)
		pl := store.GetOpenPlaylist(id)
		duration = pl.Duration
	} else {
		showing = store.GetQueueTracksCount(p)
		dig := store.GetQueueDigest(p)
		duration = dig.Duration
	}

	nano := time.Duration(duration) * time.Nanosecond
	collTracks := store.GetCollectionTracksTotalCount()
	msg := fmt.Sprintf("%v showing (%v), %v in collections", showing, nano.Truncate(time.Second), collTracks)

	statusBar.RemoveAll(StatusBarDigest)
	statusBar.Push(StatusBarDigest, msg)

	return false
}
