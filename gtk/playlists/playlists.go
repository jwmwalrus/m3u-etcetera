package playlists

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
)

var (
	tabsList        []onTab
	perspToNotebook map[m3uetcpb.Perspective]string
	playlistDlg     *gtk.Dialog
	focusRequest    struct {
		p  m3uetcpb.Perspective
		id int64
	}
)

func init() {
	perspToNotebook = map[m3uetcpb.Perspective]string{
		m3uetcpb.Perspective_MUSIC:      "music_playbar",
		m3uetcpb.Perspective_RADIO:      "radio_playbar",
		m3uetcpb.Perspective_PODCASTS:   "podcasts_playbar",
		m3uetcpb.Perspective_AUDIOBOOKS: "audiobooks_playbar",
	}
}

// Setup kickstarts playlists.
func Setup(signals *builder.Signals) (err error) {
	store.SetUpdatePlaybarViewFn(updatePlaybarView)

	(*signals).AddDetail(
		"music_playlist_new",
		"clicked",
		func(btn *gtk.Button) {
			createPlaylist(btn, m3uetcpb.Perspective_MUSIC)
		},
	)
	(*signals).AddDetail(
		"music_playbar",
		"switch-page",
		func(nb *gtk.Notebook) {
			go UpdateStatusBar(statusBarDigest)
		},
	)

	if err = builder.AddFromFile("ui/pane/playlist-dialog.ui"); err != nil {
		err = fmt.Errorf("Unable to add playlist-dialog file to builder: %v", err)
		return
	}

	playlistDlg, err = builder.GetDialog("playlist_dialog")
	if err != nil {
		err = fmt.Errorf("Unable to get playlist_dialog: %v", err)
		return
	}

	if err = setupStatusbar(); err != nil {
		return
	}
	return
}

// GetFocused returns the ID of the focused playlist
// for the given perspective.
func GetFocused(p m3uetcpb.Perspective) int64 {
	nb, err := builder.GetNotebook(perspToNotebook[p])
	if err != nil {
		slog.Error("Failed to get notebook from builder", "error", err)
		return 0
	}

	page := nb.NthPage(nb.CurrentPage())
	header := nb.TabLabel(page)
	pageName := gtk.BaseWidget(header).Name()
	for _, t := range tabsList {
		if t.headerName == pageName {
			return t.id
		}
	}
	return 0
}

// RequestFocus registers a focus request for the given playlist ID on the
// given perspective.
func RequestFocus(p m3uetcpb.Perspective, id int64) {
	for waitCount := 0; waitCount < 15; waitCount++ {
		if store.BData.PlaylistIsOpen(id) {
			break
		}

		slog.With(
			"id", id,
			"wait-count", waitCount+1,
		).Debug("Waiting for playlist to be open")
		time.Sleep(1 * time.Second)
	}

	focusRequest.p = p
	focusRequest.id = id
}

func createPlaylist(btn *gtk.Button, p m3uetcpb.Perspective) {
	req := &m3uetcpb.ExecutePlaylistActionRequest{
		Action:      m3uetcpb.PlaylistAction_PL_CREATE,
		Perspective: p,
	}
	res, err := dialer.ExecutePlaylistAction(req)
	if err != nil {
		slog.Error("Failed to execute playlist action", "error", err)
		return
	}

	RequestFocus(p, res.Id)
}

func setFocused() {
	if focusRequest.id < 0 {
		return
	}

	nb, err := builder.GetNotebook(perspToNotebook[focusRequest.p])
	if err != nil {
		slog.Error("Failed to get notebook from builder", "error", err)
		return
	}

	if focusRequest.id == 0 {
		nb.SetCurrentPage(int(focusRequest.id))
		focusRequest.id = -1
		return
	}

	var headerName string
	for _, t := range tabsList {
		if t.id == focusRequest.id {
			headerName = t.headerName
			break
		}
	}

	if headerName == "" {
		slog.Warn("Playlist is not open, so cannot be focused", "id", focusRequest.id)
		return
	}

	for ipage := 0; ipage < nb.NPages(); ipage++ {
		page := nb.NthPage(ipage)
		if page == nil {
			slog.Warn("Failed to get page from notebook", "page", ipage)
			continue
		}
		header := nb.TabLabel(page)
		pageName := gtk.BaseWidget(header).Name()
		if headerName == pageName {
			nb.SetCurrentPage(ipage)
			focusRequest.id = -1
			break
		}
	}
}
