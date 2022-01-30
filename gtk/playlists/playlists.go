package playlists

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	log "github.com/sirupsen/logrus"
)

var (
	tabsList        []onTab
	perspToNotebook map[m3uetcpb.Perspective]string
	playlistDlg     *gtk.Dialog
)

// GetFocused returns the ID of the focused playlist
// for the given perspective
func GetFocused(p m3uetcpb.Perspective) int64 {
	nb, err := builder.GetNotebook(perspToNotebook[p])
	if err != nil {
		log.Error(err)
		return 0
	}

	page, _ := nb.GetNthPage(nb.GetCurrentPage())
	header, _ := nb.GetTabLabel(page)
	pageName, _ := header.ToWidget().GetName()
	for _, t := range tabsList {
		if t.headerName == pageName {
			return t.id
		}
	}
	return 0
}

// Setup kickstarts playlists
func Setup(signals *map[string]interface{}) (err error) {
	store.SetUpdatePlaybarViewFn(updatePlaybarView)

	(*signals)["on_music_playlist_new_clicked"] = func(btn *gtk.Button) {
		createPlaylist(btn, m3uetcpb.Perspective_MUSIC)
	}

	if err = builder.AddFromFile("data/ui/pane/playlist-dialog.ui"); err != nil {
		err = fmt.Errorf("Unable to add playlist-dialog file to builder: %v", err)
		return
	}

	playlistDlg, err = builder.GetDialog("playlist_dialog")
	if err != nil {
		err = fmt.Errorf("Unable to get playlist_dialog: %v", err)
		return
	}
	return
}

func createPlaylist(btn *gtk.Button, p m3uetcpb.Perspective) {
	req := &m3uetcpb.ExecutePlaylistActionRequest{
		Action:      m3uetcpb.PlaylistAction_PL_CREATE,
		Perspective: p,
	}
	_, err := store.ExecutePlaylistAction(req)
	if err != nil {
		log.Error(err)
		return
	}
}

func init() {
	perspToNotebook = map[m3uetcpb.Perspective]string{
		m3uetcpb.Perspective_MUSIC:      "music_playbar",
		m3uetcpb.Perspective_RADIO:      "radio_playbar",
		m3uetcpb.Perspective_PODCASTS:   "podcasts_playbar",
		m3uetcpb.Perspective_AUDIOBOOKS: "audiobooks_playbar",
	}
}
