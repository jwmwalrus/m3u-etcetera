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

	if err = createMusicQueue(); err != nil {
		return
	}
	(*signals)["on_music_queue_sel_changed"] = onMusicQueueSelChanged

	if err = createMusicCollections(); err != nil {
		return
	}
	(*signals)["on_collections_sel_changed"] = onMusicCollectionsSelChanged

	return
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
		store.CColTree,
	}

	for _, i := range qcols {
		var col *gtk.TreeViewColumn
		col, err = gtk.TreeViewColumnNewWithAttribute(
			store.TreeColumn[i].Name,
			renderer,
			"text",
			i,
		)
		if err != nil {
			return
		}
		view.InsertColumn(col, -1)
	}

	model, err := store.CreateCollectionsModel(store.ArtistYearAlbumTree)
	if err != nil {
		return
	}
	view.SetModel(model)
	return
}

func createMusicQueue() (err error) {
	log.Info("Creating music queue view and model")

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

func onMusicCollectionsSelChanged(sel *gtk.TreeSelection) {
	model, iter, ok := sel.GetSelected()
	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, store.CColIDList)
		if err != nil {
			log.Error(err)
			return
		}
		_, err = value.GoValue()
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func onMusicQueueSelChanged(sel *gtk.TreeSelection) {
	model, _, ok := sel.GetSelected()
	if ok {
		rows := sel.GetSelectedRows(model)

		for l := rows; l != nil; l = l.Next() {
			path := l.Data().(*gtk.TreePath)
			iter, err := model.(*gtk.TreeModel).GetIter(path)
			if err != nil {
				log.Error(err)
				return
			}
			value, err := model.(*gtk.TreeModel).GetValue(iter, 0)
			if err != nil {
				log.Error(err)
				return
			}
			_, err = value.GoValue()
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}
