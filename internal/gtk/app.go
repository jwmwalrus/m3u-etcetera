package gtkui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/playlists"
	log "github.com/sirupsen/logrus"
)

var (
	settingsMenuSignals = &onSettingsMenu{}
)

// Setup sets the whole GTK UI
func Setup(w *gtk.ApplicationWindow, signals *map[string]interface{}) (err error) {
	settingsMenuSignals.window = w

	if err = builder.AddFromFile("data/ui/collections-dialog.ui"); err != nil {
		err = fmt.Errorf("Unable to add collections-dialog file to builder: %v", err)
		return
	}

	if err = settingsMenuSignals.createCollectionsDialog(); err != nil {
		err = fmt.Errorf("Unable to setup collections-dialog: %v", err)
		return
	}
	(*signals)["on_settings_quit_all_clicked"] = settingsMenuSignals.quitAll
	(*signals)["on_settings_collections_edit_clicked"] = settingsMenuSignals.editCollections

	settingsMenuSignals.discoverBtn, err = builder.GetToggleToolButton("collections_dialog_toggle_discover")
	if err != nil {
		log.Error(err)
		return
	}

	settingsMenuSignals.updateTagsBtn, err = builder.GetToggleToolButton("collections_dialog_toggle_update_tags")
	if err != nil {
		log.Error(err)
		return
	}

	if err = setupPlayback(signals); err != nil {
		return
	}

	if err = AddPerspectives(signals); err != nil {
		return
	}

	if err = playlists.Setup(signals); err != nil {
		return
	}
	return
}

func init() {
	perspectivesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_RADIO,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}

}
