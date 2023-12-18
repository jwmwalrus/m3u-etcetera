package gtkui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
)

var (
	interruptSignal chan os.Signal = make(chan os.Signal, 1)

	settingsMenuSignals = &onSettingsMenu{}
)

func init() {
	signal.Notify(interruptSignal, os.Interrupt, syscall.SIGTERM)

	perspectivesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_RADIO,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}
}

// Setup sets the whole GTK UI.
func Setup(w *gtk.ApplicationWindow, signals *builder.Signals) (err error) {
	settingsMenuSignals.window = w

	err = builder.AddFromFile("ui/collections-add-dialog.ui")
	if err != nil {
		err = fmt.Errorf(
			"Unable to add collections-add-dialog file to builder: %v",
			err,
		)
		return
	}

	if err = builder.AddFromFile("ui/collections-dialog.ui"); err != nil {
		err = fmt.Errorf(
			"Unable to add collections-dialog file to builder: %v",
			err,
		)
		return
	}

	err = builder.AddFromFile("ui/playlist-group-add-dialog.ui")
	if err != nil {
		err = fmt.Errorf(
			"Unable to add playlist-group-add-dialog file to builder: %v",
			err,
		)
		return
	}

	err = builder.AddFromFile("ui/playlist-groups-dialog.ui")
	if err != nil {
		err = fmt.Errorf(
			"Unable to add playlist-groups-dialog file to builder: %v",
			err,
		)
		return
	}

	if err = settingsMenuSignals.createCollectionDialogs(); err != nil {
		err = fmt.Errorf("Unable to setup collections-dialog: %v", err)
		return
	}
	(*signals).AddDetail(
		"settings_collections_add",
		"clicked",
		settingsMenuSignals.addCollection,
	)
	(*signals).AddDetail(
		"settings_collections_edit",
		"clicked",
		settingsMenuSignals.editCollections,
	)

	if err = settingsMenuSignals.createPlaylistGroupDialogs(); err != nil {
		err = fmt.Errorf("Unable to setup playlist group dialogs: %v", err)
		return
	}
	(*signals).AddDetail(
		"settings_playlist_groups_add",
		"clicked",
		settingsMenuSignals.addPlaylistGroup,
	)
	(*signals).AddDetail(
		"settings_playlist_groups_edit",
		"clicked",
		settingsMenuSignals.editPlaylistGroups,
	)
	(*signals).AddDetail(
		"settings_quit_all",
		"clicked",
		settingsMenuSignals.quitAll,
	)
	(*signals).AddDetail(
		"settings_open_files",
		"clicked",
		settingsMenuSignals.openFiles,
	)
	(*signals).AddDetail(
		"settings_open_url",
		"clicked",
		settingsMenuSignals.openURL,
	)
	(*signals).AddDetail(
		"settings_import_playlist",
		"clicked",
		settingsMenuSignals.importPlaylist,
	)

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

	dialer.SetForceExit()
	settingsMenuSignals.window.Destroy()
}
