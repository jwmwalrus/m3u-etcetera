package gtkui

import (
	"fmt"
	"strings"

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
	window *gtk.ApplicationWindow
	coll   struct {
		dlg           *gtk.Dialog
		discoverBtn   *gtk.ToggleToolButton
		updateTagsBtn *gtk.ToggleToolButton
	}
	pg struct {
		addDlg      *gtk.Dialog
		name, descr *gtk.Entry
		persp       *gtk.ComboBoxText

		editDlg *gtk.Dialog
	}
	pm *gtk.PopoverMenu
}

func (osm *onSettingsMenu) addCollection(btn *gtk.Button) {
	osm.hide()

	dlg, err := gtk.FileChooserDialogNewWith2Buttons(
		"Add collection",
		osm.window,
		gtk.FILE_CHOOSER_ACTION_SAVE,
		"Add",
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

	// TODO: directories only
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

func (osm *onSettingsMenu) addPlaylistGroup(btn *gtk.Button) {
	osm.hide()

	var name, descr string

	osm.pg.name.SetText(name)
	osm.pg.descr.SetText(descr)
	osm.pg.persp.SetActive(0)

	res := osm.pg.addDlg.Run()
	defer osm.pg.addDlg.Hide()

	switch res {
	case gtk.RESPONSE_APPLY:
		name, err := osm.pg.name.GetText()
		if err != nil {
			log.Error(err)
			return
		}
		descr, err = osm.pg.descr.GetText()
		if err != nil {
			log.Error(err)
			return
		}
		ptext := osm.pg.persp.GetActiveText()
		newPersp := m3uetcpb.Perspective_value[strings.ToUpper(ptext)]
		action := m3uetcpb.PlaylistGroupAction_PG_CREATE
		req := &m3uetcpb.ExecutePlaylistGroupActionRequest{
			Action:      action,
			Name:        name,
			Description: descr,
			Perspective: m3uetcpb.Perspective(newPersp),
		}
		if descr == "" {
			req.ResetDescription = true
		}
		_, err = store.ExecutePlaylistGroupAction(req)
		onerror.Log(err)
	case gtk.RESPONSE_CANCEL:
	default:
	}
	return
}

func (osm *onSettingsMenu) createCollectionsDialog() (err error) {
	log.Info("Setting up collections dialog")

	osm.coll.dlg, err = builder.GetDialog("collections_dialog")
	if err != nil {
		return
	}

	view, err := builder.GetTreeView("collections_dialog_view")
	if err != nil {
		return
	}

	textro, err := gtk.CellRendererTextNew()
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

func (osm *onSettingsMenu) createPlaylistGroupDialogs() (err error) {
	osm.pg.addDlg, err = builder.GetDialog("playlist_group_add_dialog")
	if err != nil {
		err = fmt.Errorf("Unable to get playlist_group_add_dialog: %v", err)
		return
	}

	osm.pg.name, err = builder.GetEntry("playlist_group_add_dialog_name")
	if err != nil {
		err = fmt.Errorf("Unable to get playlist_group_add_dialog_name: %v", err)
		return
	}

	osm.pg.descr, err = builder.GetEntry("playlist_group_add_dialog_descr")
	if err != nil {
		err = fmt.Errorf("Unable to get playlist_group_add_dialog_descr: %v", err)
		return
	}

	osm.pg.persp, err = builder.GetComboBoxText("playlist_group_add_dialog_persp")
	if err != nil {
		err = fmt.Errorf("Unable to get playlist_group_add_dialog_persp: %v", err)
		return
	}

	osm.pg.editDlg, err = builder.GetDialog("playlist_groups_dialog")
	if err != nil {
		err = fmt.Errorf("Unable to get playlist_groups_dialog: %v", err)
		return
	}

	view, err := builder.GetTreeView("playlist_groups_dialog_view")
	if err != nil {
		return
	}

	textro, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	namerw, err := store.GetPlaylistGroupRenderer(store.PGColName)
	if err != nil {
		return
	}

	descriptionrw, err := store.GetPlaylistGroupRenderer(store.PGColDescription)
	if err != nil {
		return
	}

	cols := []struct {
		idx store.ModelColumn
		r   gtk.ICellRenderer
	}{
		{store.PGColName, namerw},
		{store.PGColDescription, descriptionrw},
		{store.PGColPerspective, textro},
	}

	for _, v := range cols {
		var col *gtk.TreeViewColumn
		if renderer, ok := v.r.(*gtk.CellRendererToggle); ok {
			col, err = gtk.TreeViewColumnNewWithAttribute(
				store.PGColumns[v.idx].Name,
				renderer,
				"active",
				int(v.idx),
			)
		} else if renderer, ok := v.r.(*gtk.CellRendererText); ok {
			col, err = gtk.TreeViewColumnNewWithAttribute(
				store.PGColumns[v.idx].Name,
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

	model, err := store.CreatePlaylistGroupsModel()
	if err != nil {
		return
	}
	view.SetModel(model)
	return
}

func (osm *onSettingsMenu) editCollections(btn *gtk.Button) {
	osm.hide()

	osm.resetCollectionsToggles()

	res := osm.coll.dlg.Run()
	defer osm.coll.dlg.Hide()

	switch res {
	case gtk.RESPONSE_APPLY:
		store.ApplyCollectionChanges(osm.getCollectionsToggles())
	case gtk.RESPONSE_CANCEL:
	default:
	}
}

func (osm *onSettingsMenu) editPlaylistGroups(btn *gtk.Button) {
	osm.hide()

	res := osm.pg.editDlg.Run()
	defer osm.pg.editDlg.Hide()

	switch res {
	case gtk.RESPONSE_APPLY:
		store.ApplyPlaylistGroupChanges()
	case gtk.RESPONSE_CANCEL:
	default:
	}
	return
}

func (osm *onSettingsMenu) getCollectionsToggles() (opts store.CollectionsOptions) {
	opts.Discover = osm.coll.discoverBtn.GetActive()
	opts.UpdateTags = osm.coll.updateTagsBtn.GetActive()
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
	osm.hide()

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

	store.SetForceExit()
	osm.window.Destroy()
}

func (osm *onSettingsMenu) resetCollectionsToggles() {
	opts := store.CollectionsOptions{}
	opts.SetDefaults()

	osm.coll.discoverBtn.SetActive(opts.Discover)
	osm.coll.updateTagsBtn.SetActive(opts.UpdateTags)
}
