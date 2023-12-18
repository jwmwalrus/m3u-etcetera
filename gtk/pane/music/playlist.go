package musicpane

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/diamondburned/gotk4/pkg/gdk/v3"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/gtk/util"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
)

type onMusicPlaylist struct {
	*onContext

	export struct {
		dlg    *gtk.Dialog
		loc    *gtk.FileChooserButton
		name   *gtk.Entry
		format *gtk.ComboBoxText
		btn    *gtk.Button
	}
}

var (
	cbtid2format = map[string]struct {
		val m3uetcpb.PlaylistExportFormat
		ext string
	}{
		"playlist-out-format-m3u": {
			val: m3uetcpb.PlaylistExportFormat_PLEF_M3U,
			ext: base.SupportedPlaylistExtensionM3U,
		},
		"playlist-out-format-pls": {
			val: m3uetcpb.PlaylistExportFormat_PLEF_PLS,
			ext: base.SupportedPlaylistExtensionPLS,
		},
	}
)

func createMusicPlaylists() (ompl *onMusicPlaylist, err error) {
	slog.Info("Creating music playlists view and model")

	ompl = &onMusicPlaylist{
		onContext: &onContext{ct: playlistContext},
	}

	if err = builder.AddFromFile("ui/pane/playlist-export.ui"); err != nil {
		err = fmt.Errorf("Unable to add playlist-export file to builder: %v", err)
		return
	}

	ompl.export.dlg, err = builder.GetDialog("playlist_export")
	if err != nil {
		return
	}

	ompl.export.loc, err = builder.GetFileChooserButton("playlist_export_location")
	if err != nil {
		return
	}

	ompl.export.name, err = builder.GetEntry("playlist_export_name")
	if err != nil {
		return
	}

	ompl.export.format, err = builder.GetComboBoxText("playlist_export_format")
	if err != nil {
		return
	}

	ompl.export.btn, err = builder.GetButton("playlist_export_btn_apply")
	if err != nil {
		return
	}

	ompl.view, err = builder.GetTreeView("music_playlists_view")
	if err != nil {
		return
	}

	renderer := gtk.NewCellRendererText()

	plcols := []int{
		int(store.PLColTree),
	}

	for _, i := range plcols {
		col := gtk.NewTreeViewColumn()
		col.SetTitle(store.PLTreeColumn[i].Name)
		col.PackStart(renderer, true)
		col.AddAttribute(renderer,
			"text",
			i,
		)
		ompl.view.InsertColumn(col, -1)
	}

	model, err := store.CreatePlaylistsTreeModel(m3uetcpb.Perspective_MUSIC)
	if err != nil {
		return
	}

	ompl.view.SetModel(model)
	return
}

func (ompl *onMusicPlaylist) context(tv *gtk.TreeView, event *gdk.Event) {
	btn := event.AsButton()
	if btn.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	ids, isGroup := ompl.getPlaylistSelections(true)
	if len(ids) != 1 {
		return
	}

	if isGroup {
		return
	}

	if len(ids) != 1 {
		return
	}

	menu, err := builder.GetMenu("music_playlists_view_context")
	if err != nil {
		slog.With(
			"menu", "music_playlists_view_context",
			"error", err,
		).Error("Failed to get menu from builder")
		return
	}

	openmi, err := builder.GetMenuItem("music_playlists_view_context_open")
	if err != nil {
		slog.With(
			"menu-item", "music_playlists_view_context_open",
			"error", err,
		).Error("Failed to get menu item from builder")
		return
	}
	deletemi, err := builder.GetMenuItem("music_playlists_view_context_delete")
	if err != nil {
		slog.With(
			"menu-item", "music_playlists_view_context_delete",
			"error", err,
		).Error("Failed to get menu item from builder")
		return
	}

	pl := store.BData.GetPlaylist(ids[0])
	if pl == nil {
		slog.With("IDs", ids).
			Error("Playlist unavailable during context")
		return
	}
	openmi.SetSensitive(!pl.Open)
	deletemi.SetSensitive(!pl.Open)

	menu.PopupAtPointer(event)
}

func (ompl *onMusicPlaylist) contextDelete(mi *gtk.MenuItem) {
	ids, _ := ompl.getPlaylistSelections()
	if len(ids) != 1 {
		return
	}

	req := &m3uetcpb.ExecutePlaylistActionRequest{
		Action: m3uetcpb.PlaylistAction_PL_DESTROY,
		Id:     ids[0],
	}

	_, err := dialer.ExecutePlaylistAction(req)
	onerror.Log(err)
}

func (ompl *onMusicPlaylist) contextEdit(mi *gtk.MenuItem) {
	ids, _ := ompl.getPlaylistSelections()
	if len(ids) != 1 {
		slog.Error("Playlist selection vanished?")
		return
	}

	onerror.Log(playlists.EditPlaylist(ids[0]))
}

func (ompl *onMusicPlaylist) contextExport(mi *gtk.MenuItem) {
	ids, _ := ompl.getPlaylistSelections()
	if len(ids) != 1 {
		slog.Error("Playlist selection vanished?")
		return
	}

	validateExportBtn := func(u, name string) {
		format := cbtid2format[ompl.export.format.ActiveID()]
		dir, _ := urlstr.URLToPath(u)
		_, err := os.Stat(filepath.Join(dir, name+format.ext))
		if u == "" || name == "" || !os.IsNotExist(err) {
			ompl.export.btn.SetSensitive(false)
			return
		}
		ompl.export.btn.SetSensitive(true)
	}

	pl := store.BData.GetPlaylist(ids[0])

	ompl.export.loc.SetFilename("")
	ompl.export.name.SetText(pl.Name)
	ompl.export.format.SetActive(0)
	ompl.export.btn.SetSensitive(false)

	ompl.export.loc.Connect("file-set", func(fcb *gtk.FileChooserButton) {
		u := fcb.URI()
		name := ompl.export.name.Text()
		validateExportBtn(u, name)
	})

	ompl.export.name.Connect("changed", func(e *gtk.Entry) {
		u := ompl.export.loc.URI()
		name := e.Text()
		validateExportBtn(u, name)
	})

	res := ompl.export.dlg.Run()
	defer ompl.export.dlg.Hide()

	switch gtk.ResponseType(res) {
	case gtk.ResponseApply:
		format := cbtid2format[ompl.export.format.ActiveID()]
		u := ompl.export.loc.URI()
		name := ompl.export.name.Text()
		if name == "" {
			slog.Warn("Faild to get text from export name")
			return
		}
		dir, err := urlstr.URLToPath(u)
		if err != nil {
			slog.Error("Failed to convert URL to path", "error", err)
			return
		}
		loc, err := urlstr.PathToURLUnchecked(filepath.Join(dir, name+format.ext))
		if err != nil {
			slog.Error("Failed to convert path to URL", "error", err)
			return
		}
		req := &m3uetcpb.ExportPlaylistRequest{
			Id:       ids[0],
			Location: loc,
			Format:   format.val,
		}

		err = dialer.ExportPlaylist(req)
		if err != nil {
			slog.Error("Failed to export playlist", "error", err)
			return
		}
	case gtk.ResponseCancel:
	default:
	}
}

func (ompl *onMusicPlaylist) contextOpen(mi *gtk.MenuItem) {
	ids, _ := ompl.getPlaylistSelections()
	if len(ids) != 1 {
		slog.Error("Query selection vanished?")
		return
	}

	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Ids:    []int64{ids[0]},
		Action: m3uetcpb.PlaybarAction_BAR_OPEN,
	}

	if err := dialer.ExecutePlaybarAction(req); err != nil {
		slog.Error("Failed to execute playbar action", "error", err)
		return
	}

	playlists.RequestFocus(m3uetcpb.Perspective_MUSIC, ids[0])
}

func (ompl *onMusicPlaylist) dblClicked(tv *gtk.TreeView,
	path *gtk.TreePath, col *gtk.TreeViewColumn) {

	values, err := store.GetTreeViewTreePathValues(
		tv,
		path,
		[]store.ModelColumn{store.PLColTree, store.PLColTreeIDList, store.PLColTreeIsGroup},
	)
	if err != nil {
		slog.Error("Failed to get tree-view's tree-path values", "error", err)
		return
	}
	slog.Debug("Doouble-clicked column value", "value", values[store.CColTree])

	if values[store.PLColTreeIsGroup].(bool) {
		return
	}

	ids, err := util.StringToIDList(values[store.PLColTreeIDList].(string))
	if err != nil {
		slog.Error("Failed to convert string to ID list", "error", err)
		return
	}

	if len(ids) != 1 {
		slog.Error("Length of IDs is different from 1", "IDs", ids)
		return
	}

	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Ids:    []int64{ids[0]},
		Action: m3uetcpb.PlaybarAction_BAR_OPEN,
	}

	if err := dialer.ExecutePlaybarAction(req); err != nil {
		slog.Error("Failed to execute playback action", "error", err)
		return
	}

	playlists.RequestFocus(m3uetcpb.Perspective_MUSIC, ids[0])
}

func (ompl *onMusicPlaylist) filtered(se *gtk.SearchEntry) {
	text := se.Text()
	store.FilterPlaylistTreeBy(m3uetcpb.Perspective_MUSIC, text)
}

func (ompl *onMusicPlaylist) getPlaylistSelections(keep ...bool) (
	ids []int64, isGroup bool) {

	values := ompl.getSelectionValues(keep...)
	if len(values) == 0 {
		return
	}

	idstr, ok := values[store.PLColTreeIDList].(string)
	if !ok {
		slog.Error("This should not happen!!!", "values", values)
	}

	isGroup, ok = values[store.PLColTreeIsGroup].(bool)
	if !ok {
		slog.Error("This should not happen!!!", "values", values)
	}

	ids, err := util.StringToIDList(idstr)
	onerror.Log(err)
	return
}
