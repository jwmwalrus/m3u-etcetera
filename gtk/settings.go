package gtkui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

type onSettingsMenu struct {
	window        *gtk.ApplicationWindow
	dlg           *gtk.Dialog
	discoverBtn   *gtk.ToggleToolButton
	updateTagsBtn *gtk.ToggleToolButton
	pm            *gtk.PopoverMenu
}

func (osm *onSettingsMenu) createCollectionsDialog() (err error) {
	log.Info("Setting up collections dialog")

	osm.dlg, err = builder.GetDialog("collections_dialog")
	if err != nil {
		log.Error(err)
		return
	}

	view, err := builder.GetTreeView("collections_dialog_view")
	if err != nil {
		log.Error(err)
		return
	}

	textro, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	togglero, err := gtk.CellRendererToggleNew()
	if err != nil {
		return
	}

	namerw, err := store.GetCollectionRenderer(store.CColName)
	if err != nil {
		return
	}

	descriptionrw, err := store.GetCollectionRenderer(store.CColDescription)
	if err != nil {
		return
	}

	remotelocationrw, err := store.GetCollectionRenderer(store.CColRemoteLocation)
	if err != nil {
		return
	}

	disabledrw, err := store.GetCollectionRenderer(store.CColDisabled)
	if err != nil {
		return
	}

	remoterw, err := store.GetCollectionRenderer(store.CColRemote)
	if err != nil {
		return
	}

	rescanrw, err := store.GetCollectionRenderer(store.CColRescan)
	if err != nil {
		return
	}

	cols := []struct {
		idx store.ModelColumn
		r   gtk.ICellRenderer
	}{
		{store.CColName, namerw},
		{store.CColDescription, descriptionrw},
		{store.CColLocation, textro},
		{store.CColScanned, textro},
		{store.CColTracksView, textro},
		{store.CColHidden, togglero},
		{store.CColDisabled, disabledrw},
		{store.CColRemote, remoterw},
		{store.CColRescan, rescanrw},
		{store.CColRemoteLocation, remotelocationrw},
	}

	for _, v := range cols {
		var col *gtk.TreeViewColumn
		if renderer, ok := v.r.(*gtk.CellRendererToggle); ok {
			col, err = gtk.TreeViewColumnNewWithAttribute(
				store.CColumns[v.idx].Name,
				renderer,
				"active",
				int(v.idx),
			)
		} else if renderer, ok := v.r.(*gtk.CellRendererText); ok {
			col, err = gtk.TreeViewColumnNewWithAttribute(
				store.CColumns[v.idx].Name,
				renderer,
				"text",
				int(v.idx),
			)
		} else {
			log.Error("¿Cómo sabré si es pez o iguana?")
			continue
		}
		if err != nil {
			return
		}
		view.InsertColumn(col, -1)
	}

	model, err := store.CreateCollectionsModel()
	if err != nil {
		return
	}
	view.SetModel(model)

	return
}

func (osm *onSettingsMenu) editCollections(btn *gtk.Button) {
	osm.hide()

	osm.resetCollectionsToggles()

	res := osm.dlg.Run()
	switch res {
	case gtk.RESPONSE_APPLY:
		store.ApplyCollectionChanges(osm.getCollectionsToggles())
	case gtk.RESPONSE_CANCEL:
	default:
	}
	osm.dlg.Hide()
}

func (osm *onSettingsMenu) getCollectionsToggles() (opts store.CollectionsOptions) {
	opts.Discover = osm.discoverBtn.GetActive()
	opts.UpdateTags = osm.updateTagsBtn.GetActive()
	return
}

func (osm *onSettingsMenu) hide() {
	if osm.pm == nil {
		var err error
		osm.pm, err = builder.GetPopoverMenu("settings_menu")
		if err != nil {
			log.Error(err)
			return
		}
	}
	osm.pm.Popdown()
}

func (osm *onSettingsMenu) importPlaylist(btn *gtk.Button) {
	dlg, err := gtk.FileChooserDialogNewWith2Buttons(
		"Import playlist",
		osm.window,
		gtk.FILE_CHOOSER_ACTION_SAVE,
		"Import",
		gtk.RESPONSE_APPLY,
		"Cancel",
		gtk.RESPONSE_CANCEL,
	)
	if err != nil {
		log.Error(err)
		return
	}

	filter, err := gtk.FileFilterNew()
	if err != nil {
		log.Error(err)
		return
	}

	for _, v := range base.SupportedPlaylistExtensions {
		filter.AddPattern("*" + v)
	}
	dlg.AddFilter(filter)
	dlg.ShowAll()
	res := dlg.Run()
	switch res {
	case gtk.RESPONSE_APPLY:
		// TODO: implement
	case gtk.RESPONSE_CANCEL:
	default:
	}
	dlg.Destroy()
}

func (osm *onSettingsMenu) openFile(btn *gtk.Button) {
	osm.hide()

	dlg, err := gtk.FileChooserDialogNewWith2Buttons(
		"Open file",
		osm.window,
		gtk.FILE_CHOOSER_ACTION_SAVE,
		"Open",
		gtk.RESPONSE_APPLY,
		"Cancel",
		gtk.RESPONSE_CANCEL,
	)
	if err != nil {
		log.Error(err)
		return
	}
	defer dlg.Destroy()

	filter, err := gtk.FileFilterNew()
	if err != nil {
		log.Error(err)
		return
	}

	for _, v := range base.SupportedFileExtensions {
		filter.AddPattern("*" + v)
	}
	dlg.AddFilter(filter)
	dlg.ShowAll()
	res := dlg.Run()
	switch res {
	case gtk.RESPONSE_APPLY:
		loc := dlg.GetURI()

		plID := playlists.GetFocused(m3uetcpb.Perspective_MUSIC)

		if plID > 0 {
			action := m3uetcpb.PlaylistTrackAction_PT_APPEND
			req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
				PlaylistId: plID,
				Action:     action,
				Locations:  []string{loc},
			}

			err := store.ExecutePlaylistTrackAction(req)
			onerror.Log(err)
		} else {
			action := m3uetcpb.QueueAction_Q_APPEND
			req := &m3uetcpb.ExecuteQueueActionRequest{
				Action:    action,
				Locations: []string{loc},
			}

			err := store.ExecuteQueueAction(req)
			onerror.Log(err)
		}
	case gtk.RESPONSE_CANCEL:
	default:
	}
}

func (osm *onSettingsMenu) openURL(btn *gtk.Button) {
}

func (osm *onSettingsMenu) quitAll(btn *gtk.Button) {
	osm.hide()

	// window.Emit("delete-event", gdk.EVENT_DELETE)
	store.SetForceExit()
	osm.window.Destroy()
}

func (osm *onSettingsMenu) resetCollectionsToggles() {
	opts := store.CollectionsOptions{}
	opts.SetDefaults()

	osm.discoverBtn.SetActive(opts.Discover)
	osm.updateTagsBtn.SetActive(opts.UpdateTags)
}
