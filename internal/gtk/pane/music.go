package pane

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/store"
	log "github.com/sirupsen/logrus"
)

func setupMusic(signals *map[string]interface{}) (err error) {
	log.Info("Setting up music")

	if err = createMusicViewAndModel(); err != nil {
		return
	}

	(*signals)["on_music_queue_sel_changed"] = onMusicQueueSelChanged
	return
}

func createMusicViewAndModel() (err error) {
	log.Info("Creating view and model")

	view, err := builder.GetTreeView("music_queue_view")
	if err != nil {
		return
	}

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	qcols := []int{
		store.QColPosition,
		store.QColTitle,
		store.QColArtist,
		store.QColAlbum,
		store.QColDuration,
		store.QColTrackID,
	}

	for _, i := range qcols {
		var col *gtk.TreeViewColumn
		col, err = gtk.TreeViewColumnNewWithAttribute(
			store.QColumns[i].Name,
			renderer,
			"text",
			i,
		)
		if err != nil {
			return
		}
		view.InsertColumn(col, -1)
	}

	model, err := store.CreateQueueModel(m3uetcpb.Perspective_MUSIC)
	if err != nil {
		return
	}
	view.SetModel(model)
	return
}

func onMusicQueueSelChanged(tv *gtk.TreeSelection) {
	model := store.GetQueueModel(m3uetcpb.Perspective_MUSIC)
	rows := tv.GetSelectedRows(model)

	items := make([]string, 0, rows.Length())

	for l := rows; l != nil; l = l.Next() {
		path := l.Data().(*gtk.TreePath)
		iter, _ := model.GetIter(path)
		value, _ := model.GetValue(iter, 0)
		str, _ := value.GetString()
		items = append(items, str)
	}
}
