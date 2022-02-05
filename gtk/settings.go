package gtkui

import (
	"fmt"
	"path/filepath"
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
		addDlg           *gtk.Dialog
		loc              *gtk.FileChooserButton
		name, descr      *gtk.Entry
		persp            *gtk.ComboBoxText
		disabled, remote *gtk.CheckButton
		addBtn           *gtk.Button

		editDlg       *gtk.Dialog
		discoverBtn   *gtk.ToggleToolButton
		updateTagsBtn *gtk.ToggleToolButton
	}
	pg struct {
		addDlg, editDlg *gtk.Dialog
		name, descr     *gtk.Entry
		persp           *gtk.ComboBoxText
		addBtn          *gtk.Button
	}
	pm *gtk.PopoverMenu
}

func (osm *onSettingsMenu) addCollection(btn *gtk.Button) {
	osm.hide()

	osm.coll.loc.SetFilename("")
	osm.coll.name.SetText("")
	osm.coll.descr.SetText("")
	osm.coll.persp.SetActive(0)
	osm.coll.disabled.SetActive(false)
	osm.coll.remote.SetActive(false)
	osm.coll.addBtn.SetSensitive(false)

	osm.coll.loc.Connect("file-set", func(fcb *gtk.FileChooserButton) {
		u := fcb.GetURI()
		if u == "" {
			osm.coll.addBtn.SetSensitive(false)
			return
		}
		osm.coll.name.SetText("local:" + filepath.Base(u))
	})

	osm.coll.name.Connect("changed", func(e *gtk.Entry) {
		loc := osm.coll.loc.GetURI()
		name, _ := e.GetText()
		if loc == "" || name == "" {
			osm.coll.addBtn.SetSensitive(false)
			return
		}

		osm.coll.addBtn.SetSensitive(
			!store.CollectionAlreadyExists(loc, name),
		)
	})

	res := osm.coll.addDlg.Run()
	defer osm.coll.addDlg.Hide()

	switch res {
	case gtk.RESPONSE_APPLY:
		loc := osm.coll.loc.GetURI()
		name, _ := osm.coll.name.GetText()
		descr, _ := osm.coll.descr.GetText()
		perspid := osm.coll.persp.GetActive()
		persp := m3uetcpb.Perspective_MUSIC
		if perspid != 0 {
			persp = m3uetcpb.Perspective_AUDIOBOOKS
		}
		req := &m3uetcpb.AddCollectionRequest{
			Name:        name,
			Description: descr,
			Location:    loc,
			Perspective: persp,
			Disabled:    osm.coll.disabled.GetActive(),
			Remote:      osm.coll.remote.GetActive(),
		}

		_, err := store.AddCollection(req)
		if err != nil {
			log.Error(err)
			return
		}
	case gtk.RESPONSE_CANCEL:
	default:
	}
}

func (osm *onSettingsMenu) addPlaylistGroup(btn *gtk.Button) {
	osm.hide()

	osm.pg.name.SetText("")
	osm.pg.descr.SetText("")
	osm.pg.persp.SetActive(0)
	osm.pg.addBtn.SetSensitive(false)

	osm.pg.name.Connect("changed", func(e *gtk.Entry) {
		name, _ := e.GetText()
		if name == "" {
			osm.pg.addBtn.SetSensitive(false)
			return
		}
		osm.pg.addBtn.SetSensitive(!store.PlaylistGroupAlreadyExists(name))
	})

	res := osm.pg.addDlg.Run()
	defer osm.pg.addDlg.Hide()

	switch res {
	case gtk.RESPONSE_APPLY:
		name, err := osm.pg.name.GetText()
		if err != nil {
			log.Error(err)
			return
		}
		descr, err := osm.pg.descr.GetText()
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

func (osm *onSettingsMenu) createCollectionDialogs() (err error) {
	log.Info("Setting up collections dialog")

	osm.coll.addDlg, err = builder.GetDialog("collections_add_dialog")
	if err != nil {
		return
	}

	osm.coll.loc, err = builder.GetFileChooserButton("collections_add_dialog_location")
	if err != nil {
		return
	}

	osm.coll.name, err = builder.GetEntry("collections_add_dialog_name")
	if err != nil {
		return
	}

	osm.coll.descr, err = builder.GetEntry("collections_add_dialog_descr")
	if err != nil {
		return
	}

	osm.coll.persp, err = builder.GetComboBoxText("collections_add_dialog_perspective")
	if err != nil {
		return
	}

	osm.coll.disabled, err = builder.GetCheckButton("collections_add_dialog_disabled")
	if err != nil {
		return
	}

	osm.coll.remote, err = builder.GetCheckButton("collections_add_dialog_remote")
	if err != nil {
		return
	}

	osm.coll.addBtn, err = builder.GetButton("collections_add_dialog_btn_apply")
	if err != nil {
		return
	}

	osm.coll.editDlg, err = builder.GetDialog("collections_dialog")
	if err != nil {
		return
	}

	osm.coll.discoverBtn, err = builder.GetToggleToolButton("collections_dialog_toggle_discover")
	if err != nil {
		return
	}

	osm.coll.updateTagsBtn, err = builder.GetToggleToolButton("collections_dialog_toggle_update_tags")
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
		{store.CColPerspective, textro},
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

	model, err := store.CreateCollectionModel()
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

	osm.pg.addBtn, err = builder.GetButton("playlist_group_add_dialog_btn_apply")
	if err != nil {
		err = fmt.Errorf("Unable to get playlist_group_add_dialog_btn_apply: %v", err)
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

	res := osm.coll.editDlg.Run()
	defer osm.coll.editDlg.Hide()

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

func (osm *onSettingsMenu) getCollectionsToggles() (opts store.CollectionOptions) {
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
		gtk.FILE_CHOOSER_ACTION_OPEN,
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
	dlg.SetSelectMultiple(true)
	dlg.AddFilter(filter)
	dlg.ShowAll()
	res := dlg.Run()
	defer dlg.Destroy()

	switch res {
	case gtk.RESPONSE_APPLY:

		locs, err := dlg.GetURIs()
		if err != nil {
			log.Error(err)
			return
		}

		req := &m3uetcpb.ImportPlaylistsRequest{
			Locations: locs,
		}

		msgList, err := store.ImportPlaylists(req)
		if err != nil {
			log.Error(err)
			return
		}

		for _, msg := range msgList {
			log.Error(msg)
		}
	case gtk.RESPONSE_CANCEL:
	default:
	}
}

func (osm *onSettingsMenu) openFiles(btn *gtk.Button) {
	osm.hide()

	dlg, err := gtk.FileChooserDialogNewWith2Buttons(
		"Open files",
		osm.window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
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
	filter.SetName("Music files")

	for _, v := range base.SupportedFileExtensions {
		filter.AddPattern("*" + v)
	}
	dlg.SetSelectMultiple(true)
	dlg.AddFilter(filter)
	dlg.ShowAll()
	res := dlg.Run()
	switch res {
	case gtk.RESPONSE_APPLY:
		locs, err := dlg.GetURIs()
		if err != nil {
			log.Error(err)
			return
		}

		plID := playlists.GetFocused(m3uetcpb.Perspective_MUSIC)

		if plID > 0 {
			action := m3uetcpb.PlaylistTrackAction_PT_APPEND
			req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
				PlaylistId: plID,
				Action:     action,
				Locations:  locs,
			}

			err := store.ExecutePlaylistTrackAction(req)
			onerror.Log(err)
		} else {
			action := m3uetcpb.QueueAction_Q_APPEND
			req := &m3uetcpb.ExecuteQueueActionRequest{
				Action:    action,
				Locations: locs,
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
	opts := store.CollectionOptions{}
	opts.SetDefaults()

	osm.coll.discoverBtn.SetActive(opts.Discover)
	osm.coll.updateTagsBtn.SetActive(opts.UpdateTags)
}
