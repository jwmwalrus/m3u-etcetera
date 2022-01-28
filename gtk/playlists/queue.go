package playlists

import (
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	log "github.com/sirupsen/logrus"
)

type OnQueue struct {
	selection              interface{}
	perspective            m3uetcpb.Perspective
	view                   *gtk.TreeView
	queueID, contextMenuID string
}

func CreateQueue(p m3uetcpb.Perspective, queueID, contextMenuID string) (oq *OnQueue, err error) {
	log.WithFields(log.Fields{
		"perspective": p,
		"queueID":     queueID,
		"queueMenuID": contextMenuID,
	}).
		Info("Creating queue view and model")

	oq = &OnQueue{}
	oq.perspective = p
	oq.queueID = queueID
	oq.contextMenuID = contextMenuID

	oq.view, err = builder.GetTreeView(oq.queueID)
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

func (oq *OnQueue) Context(tv *gtk.TreeView, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() == gdk.BUTTON_SECONDARY {
		menu, err := builder.GetMenu(oq.contextMenuID)
		if err != nil {
			log.Error(err)
			return
		}
		menu.PopupAtPointer(event)
	}
}

func (oq *OnQueue) ContextClear(mi *gtk.MenuItem) {
	req := &m3uetcpb.ExecuteQueueActionRequest{
		Perspective: oq.perspective,
		Action:      m3uetcpb.QueueAction_Q_CLEAR,
	}

	if err := store.ExecuteQueueAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (oq *OnQueue) ContextDelete(mi *gtk.MenuItem) {
	values := oq.getSelection()
	if len(values) == 0 {
		return
	}

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action:   m3uetcpb.QueueAction_Q_DELETE,
		Position: int32(values[store.QColPosition].(int)),
	}

	if err := store.ExecuteQueueAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (oq *OnQueue) ContextEnqueue(mi *gtk.MenuItem) {
	values := oq.getSelection()
	if len(values) == 0 {
		return
	}

	id := values[store.QColTrackID].(int64)
	loc := values[store.QColLocation].(string)

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: m3uetcpb.QueueAction_Q_APPEND,
	}
	if id > 0 {
		req.Ids = []int64{id}
	} else {
		req.Locations = []string{loc}
	}

	if err := store.ExecuteQueueAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (oq *OnQueue) ContextMove(mi *gtk.MenuItem) {
	values := oq.getSelection()
	if len(values) == 0 {
		return
	}

	label := mi.GetLabel()
	fromPos := values[store.QColPosition].(int)
	var pos int
	if strings.Contains(label, "top") {
		if fromPos == 1 {
			return
		}
		pos = 1
	} else if strings.Contains(label, "up") {
		pos = fromPos - 1
	} else if strings.Contains(label, "down") {
		pos = fromPos + 1
	} else if strings.Contains(label, "bottom") {
		pos = values[store.QColLastPosition].(int)
		if fromPos == pos {
			return
		}
	} else {
		log.Error("Invalid/unsupported queue move")
		return
	}

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action:       m3uetcpb.QueueAction_Q_MOVE,
		Position:     int32(pos),
		FromPosition: int32(fromPos),
	}

	if err := store.ExecuteQueueAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (oq *OnQueue) ContextPlayNow(mi *gtk.MenuItem) {
	values := oq.getSelection()
	if len(values) == 0 {
		return
	}

	id := values[store.QColTrackID].(int64)
	loc := values[store.QColLocation].(string)
	pos := values[store.QColPosition].(int)

	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_PLAY,
		Force:  true,
	}
	if id > 0 {
		req.Ids = []int64{id}
	} else {
		req.Locations = []string{loc}
	}

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

func (oq *OnQueue) getSelection(keep ...bool) (values map[store.ModelColumn]interface{}) {
	values, ok := oq.selection.(map[store.ModelColumn]interface{})
	if !ok {
		log.Debug("There is no selection available for queue context")
		values = map[store.ModelColumn]interface{}{}
		return
	}

	reset := true
	if len(keep) > 0 {
		reset = !keep[0]
	}

	if reset {
		oq.selection = nil
	}
	return
}

func (oq *OnQueue) SelChanged(sel *gtk.TreeSelection) {
	var err error
	oq.selection, err = store.GetTreeSelectionValues(
		sel,
		[]store.ModelColumn{
			store.QColPosition,
			store.QColTrackID,
			store.QColLocation,
			store.QColLastPosition,
		},
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Selected collection entres: %v", oq.selection)
}
