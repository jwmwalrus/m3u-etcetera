package pane

import (
	"context"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/store"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	collsFilter *gtk.TreeModelFilter
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
	(*signals)["on_collections_view_row_activated"] = onMusicCollectionsDblClicked
	(*signals)["on_collections_filter_search_changed"] = onMusicCollectionsFiltered

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

func onMusicCollectionsDblClicked(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	imodel, err := tv.GetModel()
	if err != nil {
		log.Error(err)
		return
	}
	model, ok := imodel.(*gtk.TreeStore)
	if !ok {
		log.Error("Unable to get model from treeview")
		return
	}
	iter, err := model.GetIter(path)
	if err != nil {
		log.Error(err)
		return
	}
	value, err := model.GetValue(iter, store.CColTree)
	if err != nil {
		log.Error(err)
		return
	}
	goval, err := value.GoValue()
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Doouble-clicked column value: %v", goval)

	value, err = model.GetValue(iter, store.CColTreeIDList)
	if err != nil {
		log.Error(err)
		return
	}
	goval, err = value.GoValue()
	if err != nil {
		log.Error(err)
		return
	}

	ids, err := store.StringToIDList(goval.(string))
	if err != nil {
		log.Error(err)
		return
	}

	var cc *grpc.ClientConn
	opts := alive.GetGrpcDialOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		return
	}
	defer cc.Close()

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: m3uetcpb.QueueAction_Q_APPEND,
		Ids:    ids,
	}
	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.ExecuteQueueAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		log.Error(s.Message())
		return
	}
}

func onMusicCollectionsFiltered(se *gtk.SearchEntry) {
	text, err := se.GetText()
	if err != nil {
		log.Error(err)
		return
	}
	store.FilterCollectionsBy(text)

}

func onMusicCollectionsSelChanged(sel *gtk.TreeSelection) {
	model, iter, ok := sel.GetSelected()
	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, store.CColTreeIDList)
		if err != nil {
			log.Error(err)
			return
		}
		goval, err := value.GoValue()
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("Selected collection entry: %v", goval)
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
			goval, err := value.GoValue()
			if err != nil {
				log.Error(err)
				return
			}
			log.Infof("Selected queue row: %v", goval)
		}
	}
}
