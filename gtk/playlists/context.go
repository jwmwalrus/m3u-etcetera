package playlists

import (
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

type playlistType int

const (
	queuePlaylist playlistType = iota
	tabPlaylist
)

type onContext struct {
	id          int64
	perspective m3uetcpb.Perspective
	ctxMenu     *gtk.Menu
	view        *gtk.TreeView

	selection interface{}
}

func (oc *onContext) Context(tv *gtk.TreeView, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	values := oc.getSelection(true)
	if len(values) == 0 {
		return
	}

	var pos, lastpos int
	if oc.id > 0 {
		pos = values[store.TColPosition].(int)
		lastpos = values[store.TColLastPosition].(int)
	} else {
		pos = values[store.QColPosition].(int)
		lastpos = values[store.QColLastPosition].(int)
	}

	atTop := pos == 1
	atBottom := pos == lastpos

	oc.ctxMenu.GetChildren().Foreach(func(item interface{}) {
		w, ok := item.(*gtk.Widget)
		if !ok {
			return
		}

		l, _ := w.GetName()
		if strings.Contains(l, "top") ||
			strings.Contains(l, "up") {
			w.SetSensitive(!atTop)
		} else if strings.Contains(l, "down") ||
			strings.Contains(l, "bottom") {
			w.SetSensitive(!atBottom)
		}
	})
	oc.ctxMenu.PopupAtPointer(event)
}

func (oc *onContext) ContextClear(mi *gtk.MenuItem) {
	if oc.id > 0 {
		req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
			PlaylistId: oc.id,
			Action:     m3uetcpb.PlaylistTrackAction_PT_CLEAR,
		}

		onerror.Log(store.ExecutePlaylistTrackAction(req))
		return
	}

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Perspective: oc.perspective,
		Action:      m3uetcpb.QueueAction_Q_CLEAR,
	}

	onerror.Log(store.ExecuteQueueAction(req))
	return
}

func (oc *onContext) ContextDelete(mi *gtk.MenuItem) {
	values := oc.getSelection()
	if len(values) == 0 {
		return
	}

	if oc.id > 0 {
		req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
			PlaylistId: oc.id,
			Action:     m3uetcpb.PlaylistTrackAction_PT_DELETE,
			Position:   int32(values[store.TColPosition].(int)),
		}

		onerror.Log(store.ExecutePlaylistTrackAction(req))
		return
	}

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action:   m3uetcpb.QueueAction_Q_DELETE,
		Position: int32(values[store.QColPosition].(int)),
	}

	onerror.Log(store.ExecuteQueueAction(req))
	return
}

func (oc *onContext) ContextEnqueue(mi *gtk.MenuItem) {
	values := oc.getSelection()
	if len(values) == 0 {
		return
	}

	var id int64
	var loc string
	if oc.id > 0 {
		id = values[store.TColTrackID].(int64)
		loc = values[store.TColLocation].(string)
	} else {
		id = values[store.QColTrackID].(int64)
		loc = values[store.QColLocation].(string)
	}

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

func (oc *onContext) ContextMove(mi *gtk.MenuItem) {
	values := oc.getSelection()
	if len(values) == 0 {
		return
	}

	var rowPosition, lastPosition store.ModelColumn

	if oc.id > 0 {
		rowPosition = store.TColPosition
		lastPosition = store.TColLastPosition
	} else {
		rowPosition = store.QColPosition
		lastPosition = store.QColLastPosition
	}

	label := mi.GetLabel()
	fromPos := values[rowPosition].(int)
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
		pos = values[lastPosition].(int)
		if fromPos == pos {
			return
		}
	} else {
		log.Error("Invalid/unsupported queue move")
		return
	}

	if oc.id > 0 {
		req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
			PlaylistId:   oc.id,
			Action:       m3uetcpb.PlaylistTrackAction_PT_MOVE,
			Position:     int32(pos),
			FromPosition: int32(fromPos),
		}

		onerror.Log(store.ExecutePlaylistTrackAction(req))
		return
	}

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action:       m3uetcpb.QueueAction_Q_MOVE,
		Position:     int32(pos),
		FromPosition: int32(fromPos),
	}

	onerror.Log(store.ExecuteQueueAction(req))
	return
}

func (oc *onContext) ContextPlayNow(mi *gtk.MenuItem) {
	values := oc.getSelection()
	if len(values) == 0 {
		return
	}

	if oc.id > 0 {
		pos := values[store.TColPosition].(int)

		req := &m3uetcpb.ExecutePlaybarActionRequest{
			Action:   m3uetcpb.PlaybarAction_BAR_ACTIVATE,
			Position: int32(pos),
			Ids:      []int64{oc.id},
		}

		onerror.Log(store.ExecutePlaybarAction(req))
		return
	}

	id := values[store.QColTrackID].(int64)
	loc := values[store.QColLocation].(string)
	pos := values[store.QColPosition].(int)

	reqbar := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_PLAY,
		Force:  true,
	}
	if id > 0 {
		reqbar.Ids = []int64{id}
	} else {
		reqbar.Locations = []string{loc}
	}

	if err := store.ExecutePlaybackAction(reqbar); err != nil {
		log.Error(err)
		return
	}

	reqq := &m3uetcpb.ExecuteQueueActionRequest{
		Action:   m3uetcpb.QueueAction_Q_DELETE,
		Position: int32(pos),
	}

	onerror.Log(store.ExecuteQueueAction(reqq))
	return
}

func (oc *onContext) SelChanged(sel *gtk.TreeSelection) {
	var err error
	if oc.id > 0 {
		oc.selection, err = store.GetTreeSelectionValues(
			sel,
			[]store.ModelColumn{
				store.TColPosition,
				store.TColTrackID,
				store.TColLocation,
				store.TColLastPosition,
			},
		)
	} else {
		oc.selection, err = store.GetTreeSelectionValues(
			sel,
			[]store.ModelColumn{
				store.QColPosition,
				store.QColTrackID,
				store.QColLocation,
				store.QColLastPosition,
			},
		)
	}
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Selected context entres: %v", oc.selection)
}

func (oc *onContext) getSelection(keep ...bool) (
	values map[store.ModelColumn]interface{}) {

	if oc.selection == nil {
		sel, err := oc.view.GetSelection()
		if err != nil {
			log.Error(err)
			return
		}
		oc.SelChanged(sel)
	}

	values, ok := oc.selection.(map[store.ModelColumn]interface{})
	if !ok {
		log.Debug("There is no selection available for context")
		values = map[store.ModelColumn]interface{}{}
		return
	}

	reset := true
	if len(keep) > 0 {
		reset = !keep[0]
	}

	if reset {
		oc.selection = nil
	}
	return
}
