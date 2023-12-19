package gtkui

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
)

const (
	modifiableColumnHeaderSuffix = " (*)"
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
		u := fcb.URI()
		if u == "" {
			osm.coll.addBtn.SetSensitive(false)
			return
		}
		osm.coll.name.SetText("local:" + filepath.Base(u))
	})

	osm.coll.name.Connect("changed", func(e *gtk.Entry) {
		loc := osm.coll.loc.URI()
		name := e.Text()
		if loc == "" || name == "" {
			osm.coll.addBtn.SetSensitive(false)
			return
		}

		osm.coll.addBtn.SetSensitive(
			!store.CData.CollectionAlreadyExists(loc, name),
		)
	})

	res := osm.coll.addDlg.Run()
	defer osm.coll.addDlg.Hide()

	switch gtk.ResponseType(res) {
	case gtk.ResponseApply:
		loc := osm.coll.loc.URI()
		name := osm.coll.name.Text()
		descr := osm.coll.descr.Text()
		perspid := osm.coll.persp.Active()
		persp := m3uetcpb.Perspective_MUSIC
		if perspid != 0 {
			persp = m3uetcpb.Perspective_AUDIOBOOKS
		}
		req := &m3uetcpb.AddCollectionRequest{
			Name:        name,
			Description: descr,
			Location:    loc,
			Perspective: persp,
			Disabled:    osm.coll.disabled.Active(),
			Remote:      osm.coll.remote.Active(),
		}

		_, err := dialer.AddCollection(req)
		if err != nil {
			slog.Error("Failed to add collection", "error", err)
			return
		}
	case gtk.ResponseCancel:
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
		name := e.Text()
		if name == "" {
			osm.pg.addBtn.SetSensitive(false)
			return
		}
		osm.pg.addBtn.SetSensitive(!store.BData.PlaylistGroupAlreadyExists(name))
	})

	res := osm.pg.addDlg.Run()
	defer osm.pg.addDlg.Hide()

	switch gtk.ResponseType(res) {
	case gtk.ResponseApply:
		name := osm.pg.name.Text()
		descr := osm.pg.descr.Text()
		ptext := osm.pg.persp.ActiveText()
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
		_, err := dialer.ExecutePlaylistGroupAction(req)
		onerror.Log(err)
	case gtk.ResponseCancel:
	default:
	}
}

func (osm *onSettingsMenu) createCollectionDialogs() (err error) {
	slog.Info("Setting up collections dialog")

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

	osm.coll.discoverBtn, err = builder.GetToggleToolButton(
		"collections_dialog_toggle_discover",
	)
	if err != nil {
		return
	}

	osm.coll.updateTagsBtn, err = builder.GetToggleToolButton(
		"collections_dialog_toggle_update_tags",
	)
	if err != nil {
		return
	}

	view, err := builder.GetTreeView("collections_dialog_view")
	if err != nil {
		return
	}

	model, err := store.CreateCollectionModel()
	if err != nil {
		return
	}

	textro := gtk.NewCellRendererText()

	cr := store.Renderer{Model: model, Columns: store.CColumns}

	namerw, err := cr.NewEditable(store.CColName)
	if err != nil {
		return
	}

	descriptionrw, err := cr.NewEditable(store.CColDescription)
	if err != nil {
		return
	}

	remotelocationrw, err := cr.NewEditable(store.CColRemoteLocation)
	if err != nil {
		return
	}

	disabledrw, err := cr.NewActivatable(store.CColDisabled)
	if err != nil {
		return
	}

	remoterw, err := cr.NewActivatable(store.CColRemote)
	if err != nil {
		return
	}

	rescanrw, err := cr.NewActivatable(store.CColActionRescan)
	if err != nil {
		return
	}

	removerw, err := cr.NewActivatable(store.CColActionRemove)
	if err != nil {
		return
	}

	cols := []struct {
		idx       store.ModelColumn
		r         gtk.CellRendererer
		canModify bool
	}{
		{store.CColName, namerw, true},
		{store.CColPerspective, textro, false},
		{store.CColTracksView, textro, false},
		{store.CColActionRescan, rescanrw, true},
		{store.CColActionRemove, removerw, true},
		{store.CColDisabled, disabledrw, true},
		{store.CColRemote, remoterw, true},
		{store.CColDescription, descriptionrw, true},
		{store.CColLocation, textro, false},
		{store.CColRemoteLocation, remotelocationrw, true},
	}

	for _, v := range cols {
		var suffix string
		if v.canModify {
			suffix = modifiableColumnHeaderSuffix
		}
		col := gtk.NewTreeViewColumn()
		col.SetTitle(store.CColumns[v.idx].Name + suffix)
		if renderer, ok := v.r.(*gtk.CellRendererToggle); ok {
			col.PackStart(renderer, true)
			col.AddAttribute(
				renderer,
				"active",
				int(v.idx),
			)
		} else if renderer, ok := v.r.(*gtk.CellRendererText); ok {
			col.PackStart(renderer, true)
			col.AddAttribute(
				renderer,
				"text",
				int(v.idx),
			)
		} else {
			slog.Error("¿Cómo sabré si es pez o iguana?")
			continue
		}
		view.InsertColumn(col, -1)
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

	model, err := store.CreatePlaylistGroupsModel()
	if err != nil {
		return
	}

	pgr := store.Renderer{Model: model, Columns: store.PGColumns}

	textro := gtk.NewCellRendererText()

	namerw, err := pgr.NewEditable(store.PGColName)
	if err != nil {
		return
	}

	descriptionrw, err := pgr.NewEditable(store.PGColDescription)
	if err != nil {
		return
	}

	removerw, err := pgr.NewActivatable(store.PGColActionRemove)
	if err != nil {
		return
	}

	cols := []struct {
		idx       store.ModelColumn
		r         gtk.CellRendererer
		canModify bool
	}{
		{store.PGColName, namerw, true},
		{store.PGColPerspective, textro, false},
		{store.PGColActionRemove, removerw, true},
		{store.PGColDescription, descriptionrw, true},
	}

	for _, v := range cols {
		var suffix string
		if v.canModify {
			suffix = modifiableColumnHeaderSuffix
		}
		col := gtk.NewTreeViewColumn()
		col.SetTitle(store.PGColumns[v.idx].Name + suffix)
		if renderer, ok := v.r.(*gtk.CellRendererToggle); ok {
			col.PackStart(renderer, true)
			col.AddAttribute(
				renderer,
				"active",
				int(v.idx),
			)
		} else if renderer, ok := v.r.(*gtk.CellRendererText); ok {
			col.PackStart(renderer, true)
			col.AddAttribute(
				renderer,
				"text",
				int(v.idx),
			)
		} else {
			slog.Error("¿Cómo sabré si es pez o iguana?")
			continue
		}
		view.InsertColumn(col, -1)
	}

	view.SetModel(model)
	return
}

func (osm *onSettingsMenu) editCollections(btn *gtk.Button) {
	osm.hide()

	osm.resetCollectionsToggles()

	res := osm.coll.editDlg.Run()
	defer osm.coll.editDlg.Hide()

	switch gtk.ResponseType(res) {
	case gtk.ResponseApply:
		dialer.ApplyCollectionChanges(osm.getCollectionsToggles())
	case gtk.ResponseCancel:
	default:
	}
}

func (osm *onSettingsMenu) editPlaylistGroups(btn *gtk.Button) {
	osm.hide()

	res := osm.pg.editDlg.Run()
	defer osm.pg.editDlg.Hide()

	switch gtk.ResponseType(res) {
	case gtk.ResponseApply:
		dialer.ApplyPlaylistGroupChanges()
	case gtk.ResponseCancel:
	default:
	}
}

func (osm *onSettingsMenu) getCollectionsToggles() (
	opts store.CollectionOptions) {

	opts.Discover = osm.coll.discoverBtn.Active()
	opts.UpdateTags = osm.coll.updateTagsBtn.Active()
	return
}

func (osm *onSettingsMenu) hide() {
	if osm.pm == nil {
		var err error
		osm.pm, err = builder.GetPopoverMenu("settings_menu")
		if err != nil {
			slog.With(
				"menu", "settings_menu",
				"error", err,
			).Error("Failed to get popover menu")
			return
		}
	}
	osm.pm.Popdown()
}

func (osm *onSettingsMenu) importPlaylist(btn *gtk.Button) {
	osm.hide()

	dlg := gtk.NewFileChooserNative(
		"Import playlist",
		&osm.window.Window,
		gtk.FileChooserActionOpen,
		"Import",
		"Cancel",
	)
	if dlg == nil {
		slog.Error("Failed to create file-chooser-dialog")
		return
	}

	filter := gtk.NewFileFilter()
	filter.SetName("Playlist files")

	for _, v := range base.SupportedPlaylistExtensions {
		filter.AddPattern("*" + v)
	}
	dlg.SetSelectMultiple(true)
	dlg.AddFilter(filter)
	res := dlg.Run()
	defer dlg.Destroy()

	switch gtk.ResponseType(res) {
	case gtk.ResponseApply, gtk.ResponseAccept:
		locs := dlg.URIs()

		req := &m3uetcpb.ImportPlaylistsRequest{
			Locations: locs,
		}

		msgList, err := dialer.ImportPlaylists(req)
		if err != nil {
			slog.Error("Failed to import playlists", "error", err)
			return
		}

		for _, msg := range msgList {
			slog.Error(msg)
		}
	case gtk.ResponseCancel:
	default:
	}
}

func (osm *onSettingsMenu) openFiles(btn *gtk.Button) {
	osm.hide()

	dlg := gtk.NewFileChooserNative(
		"Open files",
		&osm.window.Window,
		gtk.FileChooserActionOpen,
		"Open",
		"Cancel",
	)
	if dlg == nil {
		slog.Error("Failed to create file-chooser-dialog")
		return
	}
	defer dlg.Destroy()

	filter := gtk.NewFileFilter()
	filter.SetName("Music files")

	fileExts := base.SupportedFileExtensions
	fileExts = append(fileExts, base.SupportedPlaylistExtensions...)

	for _, v := range fileExts {
		filter.AddPattern("*" + v)
	}

	dlg.SetSelectMultiple(true)
	dlg.AddFilter(filter)
	res := dlg.Run()

	switch gtk.ResponseType(res) {
	case gtk.ResponseApply, gtk.ResponseAccept:
		locs := dlg.URIs()

		var mmLocs, plLocs []string
		for i := range locs {
			if base.IsSupportedPlaylistURL(locs[i]) {
				plLocs = append(plLocs, locs[i])
				continue
			}
			mmLocs = append(mmLocs, locs[i])
		}

		if len(mmLocs) > 0 {
			plID := playlists.GetFocused(m3uetcpb.Perspective_MUSIC)

			if plID > 0 {
				action := m3uetcpb.PlaylistTrackAction_PT_APPEND
				req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
					PlaylistId: plID,
					Action:     action,
					Locations:  mmLocs,
				}

				err := dialer.ExecutePlaylistTrackAction(req)
				onerror.Log(err)
			} else {
				action := m3uetcpb.QueueAction_Q_APPEND
				req := &m3uetcpb.ExecuteQueueActionRequest{
					Action:    action,
					Locations: mmLocs,
				}

				err := dialer.ExecuteQueueAction(req)
				onerror.Log(err)
			}
		}

		if len(plLocs) > 0 {
			req := &m3uetcpb.ImportPlaylistsRequest{
				Locations:   plLocs,
				AsTransient: true,
			}

			msgList, err := dialer.ImportPlaylists(req)
			if err != nil {
				slog.Error("Failed to import playlists", "error", err)
				return
			}

			for _, msg := range msgList {
				slog.Error(msg)
			}
		}
	case gtk.ResponseCancel:
	default:
	}
}

func (osm *onSettingsMenu) openURL(btn *gtk.Button) {
}

func (osm *onSettingsMenu) quitAll(btn *gtk.Button) {
	osm.hide()

	dialer.SetForceExit()
	osm.window.Destroy()
}

func (osm *onSettingsMenu) resetCollectionsToggles() {
	opts := store.CollectionOptions{}
	opts.SetDefaults()

	osm.coll.discoverBtn.SetActive(opts.Discover)
	osm.coll.updateTagsBtn.SetActive(opts.UpdateTags)
}
