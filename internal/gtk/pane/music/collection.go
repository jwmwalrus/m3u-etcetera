package musicpane

import (
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/store"
	log "github.com/sirupsen/logrus"
)

type onMusicCollections struct {
	selection interface{}
}

func createMusicCollections() (err error) {
	log.Info("Creating music collections view and model")

	view, err := builder.GetTreeView("collections_view")
	if err != nil {
		return
	}

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	qcols := []int{
		int(store.CColTree),
	}

	for _, i := range qcols {
		var col *gtk.TreeViewColumn
		col, err = gtk.TreeViewColumnNewWithAttribute(
			store.CTreeColumn[i].Name,
			renderer,
			"text",
			i,
		)
		if err != nil {
			return
		}
		view.InsertColumn(col, -1)
	}

	model, err := store.CreateCollectionTreeModel(store.ArtistYearAlbumTree)
	if err != nil {
		return
	}

	view.SetModel(model)
	return
}

func (omc *onMusicCollections) context(tv *gtk.TreeView, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() == gdk.BUTTON_SECONDARY {
		menu, err := builder.GetMenu("collections_view_context")
		if err != nil {
			log.Error(err)
			return
		}
		menu.PopupAtPointer(event)
	}
}

func (omc *onMusicCollections) contextAppend(mi *gtk.MenuItem) {
	ids := omc.getSelection()
	if len(ids) == 0 {
		return
	}

	action := m3uetcpb.QueueAction_Q_APPEND
	if strings.Contains(mi.GetLabel(), "Preppend") {
		action = m3uetcpb.QueueAction_Q_PREPPEND
	}
	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: action,
		Ids:    ids,
	}

	if err := store.ExecuteQueueAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (omc *onMusicCollections) contextPlayNow(mi *gtk.MenuItem) {
	ids := omc.getSelection()
	if len(ids) == 0 {
		return
	}

	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_PLAY,
		Force:  true,
		Ids:    ids,
	}

	if err := store.ExecutePlaybackAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (omc *onMusicCollections) dblClicked(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	values, err := store.GetTreeStoreValues(tv, path, []store.ModelColumn{store.CColTree, store.CColTreeIDList})
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Doouble-clicked column value: %v", values[store.CColTree])

	ids, err := store.StringToIDList(values[store.CColTreeIDList].(string))
	if err != nil {
		log.Error(err)
		return
	}

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: m3uetcpb.QueueAction_Q_APPEND,
		Ids:    ids,
	}

	if err := store.ExecuteQueueAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (omc *onMusicCollections) filtered(se *gtk.SearchEntry) {
	text, err := se.GetText()
	if err != nil {
		log.Error(err)
		return
	}
	store.FilterCollectionsBy(text)

}

func (omc *onMusicCollections) getSelection(keep ...bool) (ids []int64) {
	value, ok := omc.selection.(string)
	if !ok {
		log.Debug("There is no selection available for collection context")
		return
	}

	ids, err := store.StringToIDList(value)
	if err != nil {
		log.Errorf("Error parsing selection value for collection context: %v", err)
	}

	reset := true
	if len(keep) > 0 {
		reset = !keep[0]
	}
	if reset {
		omc.selection = nil
	}
	return
}

func (omc *onMusicCollections) selChanged(sel *gtk.TreeSelection) {
	var err error
	omc.selection, err = store.GetTreeSelectionValue(sel, store.CColTreeIDList)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Selected collection entry: %v", omc.selection)
}
