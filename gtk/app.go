package gtkui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	log "github.com/sirupsen/logrus"
)

var (
	interruptSignal chan os.Signal

	settingsMenuSignals = &onSettingsMenu{}
)

// Setup sets the whole GTK UI
func Setup(w *gtk.ApplicationWindow, signals *map[string]interface{}) (err error) {
	settingsMenuSignals.window = w

	if err = builder.AddFromFile("data/ui/collections-dialog.ui"); err != nil {
		err = fmt.Errorf("Unable to add collections-dialog file to builder: %v", err)
		return
	}

	if err = builder.AddFromFile("data/ui/playlist-group-add-dialog.ui"); err != nil {
		err = fmt.Errorf("Unable to add playlist-group-add-dialog file to builder: %v", err)
		return
	}

	if err = builder.AddFromFile("data/ui/playlist-groups-dialog.ui"); err != nil {
		err = fmt.Errorf("Unable to add playlist-groups-dialog file to builder: %v", err)
		return
	}

	if err = settingsMenuSignals.createCollectionsDialog(); err != nil {
		err = fmt.Errorf("Unable to setup collections-dialog: %v", err)
		return
	}
	(*signals)["on_settings_collections_add_clicked"] = settingsMenuSignals.addCollection
	(*signals)["on_settings_collections_edit_clicked"] = settingsMenuSignals.editCollections

	settingsMenuSignals.coll.discoverBtn, err = builder.GetToggleToolButton("collections_dialog_toggle_discover")
	if err != nil {
		log.Error(err)
		return
	}

	settingsMenuSignals.coll.updateTagsBtn, err = builder.GetToggleToolButton("collections_dialog_toggle_update_tags")
	if err != nil {
		log.Error(err)
		return
	}

	if err = settingsMenuSignals.createPlaylistGroupDialogs(); err != nil {
		err = fmt.Errorf("Unable to setup playlist group dialogs: %v", err)
		return
	}
	(*signals)["on_settings_playlist_groups_add_clicked"] = settingsMenuSignals.addPlaylistGroup
	(*signals)["on_settings_playlist_groups_edit_clicked"] = settingsMenuSignals.editPlaylistGroups

	(*signals)["on_settings_quit_all_clicked"] = settingsMenuSignals.quitAll
	(*signals)["on_settings_open_file_clicked"] = settingsMenuSignals.openFile
	(*signals)["on_settings_open_url_clicked"] = settingsMenuSignals.openURL
	(*signals)["on_settings_import_playlist_clicked"] = settingsMenuSignals.importPlaylist

	if err = setupPlayback(signals); err != nil {
		return
	}

	if err = addPerspectives(signals); err != nil {
		return
	}

	if err = playlists.Setup(signals); err != nil {
		return
	}

	go onInterruptSignal()

	return
}

func onInterruptSignal() {
	<-interruptSignal
	store.SetForceExit()
	settingsMenuSignals.window.Destroy()

}

func init() {
	interruptSignal = make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGTERM)

	perspectivesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_RADIO,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}

}
