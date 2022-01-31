package musicpane

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

type onMusicPlaylist struct {
	*onContext

	dlg                               *gtk.Dialog
	name, id, descr, params, from, to *gtk.Entry
	rating, limit                     *gtk.SpinButton
	random                            *gtk.CheckButton
	resultsLabel                      *gtk.Label
}

func createMusicPlaylists() (ompl *onMusicPlaylist, err error) {
	log.Info("Creating music playlists view and model")

	ompl = &onMusicPlaylist{
		onContext: &onContext{ct: playlistContext},
	}

	view, err := builder.GetTreeView("music_playlists_view")
	if err != nil {
		return
	}

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	plcols := []int{
		int(store.PLColTree),
	}

	for _, i := range plcols {
		var col *gtk.TreeViewColumn
		col, err = gtk.TreeViewColumnNewWithAttribute(
			store.PLTreeColumn[i].Name,
			renderer,
			"text",
			i,
		)
		if err != nil {
			return
		}
		view.InsertColumn(col, -1)
	}

	model, err := store.CreatePlaylistsTreeModel(m3uetcpb.Perspective_MUSIC)
	if err != nil {
		return
	}

	view.SetModel(model)
	return
}

func (ompl *onMusicPlaylist) context(tv *gtk.TreeView, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() == gdk.BUTTON_SECONDARY {
		ids := ompl.getSelection(true)
		if len(ids) != 1 {
			return
		}

		menu, err := builder.GetMenu("music_playlists_view_context")
		if err != nil {
			log.Error(err)
			return
		}

		pl := store.GetPlaylist(ids[0])
		if pl == nil {
			log.WithField("ids", ids).
				Error("Playlist unavailable during context")
			return
		}

		openmi, err := builder.GetMenuItem("music_playlists_view_context_open")
		if err != nil {
			log.Error(err)
			return
		}
		deletemi, err := builder.GetMenuItem("music_playlists_view_context_delete")
		if err != nil {
			log.Error(err)
			return
		}

		openmi.SetSensitive(!pl.Open)
		deletemi.SetSensitive(!pl.Open)

		menu.PopupAtPointer(event)
	}
}

func (ompl *onMusicPlaylist) contextDelete(mi *gtk.MenuItem) {
	ids := ompl.getSelection()
	if len(ids) != 1 {
		log.Error("Query selection vanished?")
		return
	}

	req := &m3uetcpb.ExecutePlaylistActionRequest{
		Action: m3uetcpb.PlaylistAction_PL_DESTROY,
		Id:     ids[0],
	}

	if _, err := store.ExecutePlaylistAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (ompl *onMusicPlaylist) contextEdit(mi *gtk.MenuItem) {
	ids := ompl.getSelection()
	if len(ids) != 1 {
		log.Error("Query selection vanished?")
		return
	}
	onerror.Log(playlists.EditPlaylist(ids[0]))
}

func (ompl *onMusicPlaylist) contextOpen(mi *gtk.MenuItem) {
	ids := ompl.getSelection()
	if len(ids) != 1 {
		log.Error("Query selection vanished?")
		return
	}

	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Ids:    []int64{ids[0]},
		Action: m3uetcpb.PlaybarAction_BAR_OPEN,
	}

	if err := store.ExecutePlaybarAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (ompl *onMusicPlaylist) dblClicked(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	values, err := store.GetTreeStoreValues(tv, path, []store.ModelColumn{store.QYColTree, store.QYColTreeIDList})
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Doouble-clicked column value: %v", values[store.CColTree])

	ids, err := store.StringToIDList(values[store.QYColTreeIDList].(string))
	if err != nil {
		log.Error(err)
		return
	}

	if len(ids) != 1 {
		log.Errorf("Length of ids is different from 1: %+v", ids)
		return
	}

	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Ids:    []int64{ids[0]},
		Action: m3uetcpb.PlaybarAction_BAR_OPEN,
	}

	if err := store.ExecutePlaybarAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (ompl *onMusicPlaylist) filtered(se *gtk.SearchEntry) {
	text, err := se.GetText()
	if err != nil {
		log.Error(err)
		return
	}
	store.FilterPlaylistTreeBy(m3uetcpb.Perspective_MUSIC, text)
}
