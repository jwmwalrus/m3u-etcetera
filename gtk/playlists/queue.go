package playlists

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/bnp/stringing"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	log "github.com/sirupsen/logrus"
)

// OnQueue handles queue signals
type OnQueue struct {
	*onContext

	view *gtk.TreeView
}

// CreateQueue returns a queue signals handler
func CreateQueue(p m3uetcpb.Perspective, queueID, contextMenuID string) (oq *OnQueue, err error) {
	log.WithFields(log.Fields{
		"perspective": p,
		"queueID":     queueID,
		"queueMenuID": contextMenuID,
	}).
		Info("Creating queue view and model")

	ctxMenu, err := builder.GetMenu(contextMenuID)
	if err != nil {
		log.Error(err)
		return
	}

	miSuffix := stringing.GetRandomString(6)
	for _, l := range []string{"top", "up", "down", "bottom"} {
		mi, err := builder.GetMenuItem(contextMenuID + "_" + l)
		if err != nil {
			log.Error(err)
			continue
		}
		mi.SetName(fmt.Sprintf("menuitem-%s-%s", l, miSuffix))
	}

	oq = &OnQueue{
		onContext: &onContext{
			id:          0,
			perspective: p,
			ctxMenu:     ctxMenu,
		},
	}

	oq.view, err = builder.GetTreeView(queueID)
	if err != nil {
		return
	}

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	qcols := []int{
		int(store.QColPosition),
		int(store.QColTitle),
		int(store.QColArtist),
		int(store.QColAlbum),
		int(store.QColDuration),
		int(store.QColTrackID),
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
		oq.view.InsertColumn(col, -1)
	}

	model, err := store.CreateQueueModel(oq.perspective)
	if err != nil {
		return
	}
	oq.view.SetModel(model)
	return
}

// DblClicked handles double-click
func (oq *OnQueue) DblClicked(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	values, err := store.GetListStoreValues(
		tv,
		path,
		[]store.ModelColumn{
			store.QColPosition,
			store.QColTrackID,
			store.QColLocation,
		},
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Doouble-clicked column values: %v", values)

	id := values[store.QColTrackID].(int64)
	pos := values[store.QColPosition].(int)
	loc := values[store.QColLocation].(string)

	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_PLAY,
		Force:  true,
	}

	if id > 0 {
		req.Ids = []int64{id}
	} else {
		req.Locations = []string{loc}
	}

	time.Sleep(200 * time.Millisecond)

	if err := store.ExecutePlaybackAction(req); err != nil {
		log.Error(err)
		return
	}

	req2 := &m3uetcpb.ExecuteQueueActionRequest{
		Action:   m3uetcpb.QueueAction_Q_DELETE,
		Position: int32(pos),
	}

	if err := store.ExecuteQueueAction(req2); err != nil {
		log.Error(err)
		return
	}
}
