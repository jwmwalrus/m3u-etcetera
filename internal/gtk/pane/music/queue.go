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

type onMusicQueue struct {
	selection interface{}
}

var (
	musicQueueSignals = &onMusicQueue{}
)

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
		view.InsertColumn(col, -1)
	}

	model, err := store.CreateQueueModel(m3uetcpb.Perspective_MUSIC)
	if err != nil {
		return
	}
	view.SetModel(model)
	return
}

func (omq *onMusicQueue) context(tv *gtk.TreeView, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() == gdk.BUTTON_SECONDARY {
		menu, err := builder.GetMenu("music_queue_view_context")
		if err != nil {
			log.Error(err)
			return
		}
		menu.PopupAtPointer(event)
	}
}

func (omq *onMusicQueue) contextDelete(mi *gtk.MenuItem) {
	values, ok := omq.selection.(map[store.TreeModelColumn]interface{})
	if !ok {
		log.Error("There is no selection available for music queue context")
		return
	}
	omq.selection = nil

	req2 := &m3uetcpb.ExecuteQueueActionRequest{
		Action:   m3uetcpb.QueueAction_Q_DELETE,
		Position: int32(values[store.QColPosition].(int)),
	}

	if err := store.ExecuteQueueAction(req2); err != nil {
		log.Error(err)
		return
	}
}

func (omq *onMusicQueue) contextEnqueue(mi *gtk.MenuItem) {
	values, ok := omq.selection.(map[store.TreeModelColumn]interface{})
	if !ok {
		log.Error("There is no selection available for music queue context")
		return
	}
	omq.selection = nil

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

func (omq *onMusicQueue) contextMove(mi *gtk.MenuItem) {
	values, ok := omq.selection.(map[store.TreeModelColumn]interface{})
	if !ok {
		log.Error("There is no selection available for music queue MOVE context")
		return
	}
	omq.selection = nil

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

	req2 := &m3uetcpb.ExecuteQueueActionRequest{
		Action:       m3uetcpb.QueueAction_Q_MOVE,
		Position:     int32(pos),
		FromPosition: int32(fromPos),
	}

	if err := store.ExecuteQueueAction(req2); err != nil {
		log.Error(err)
		return
	}
}

func (omq *onMusicQueue) contextPlayNow(mi *gtk.MenuItem) {
	values, ok := omq.selection.(map[store.TreeModelColumn]interface{})
	if !ok {
		log.Error("There is no selection available for music queue context")
		return
	}
	omq.selection = nil

	id := values[store.QColTrackID].(int64)
	loc := values[store.QColLocation].(string)
	pos, ok := values[store.QColPosition].(int)

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

func (omq *onMusicQueue) dblClicked(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	values, err := store.GetListStoreValues(
		tv,
		path,
		[]store.TreeModelColumn{
			store.QColPosition,
			store.QColTrackID,
			store.QColLocation,
		},
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Doouble-clicked column values: %v", values[store.CColTree])

	id := values[store.QColTrackID].(int64)
	pos := values[store.QColPosition].(int)
	loc := values[store.QColLocation].(string)

	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action:    m3uetcpb.PlaybackAction_PB_PLAY,
		Force:     true,
		Ids:       []int64{},
		Locations: []string{},
	}

	if id > 0 {
		req.Ids = append(req.Ids, id)
	} else {
		req.Locations = append(req.Locations, loc)
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

func (omq *onMusicQueue) selChanged(sel *gtk.TreeSelection) {
	var err error
	omq.selection, err = store.GetTreeSelectionValues(
		sel,
		[]store.TreeModelColumn{
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
	log.Infof("Selected collection entres: %v", omq.selection)
	/*
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
				log.Debugf("Selected queue row: %v", goval)
			}
		}
	*/
}
