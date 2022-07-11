package playlists

import (
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

type onContext struct {
	id          int64
	perspective m3uetcpb.Perspective
	ctxMenu     *gtk.Menu
	view        *gtk.TreeView

	selection struct {
		values interface{}
		keep   bool
		paths  []*gtk.TreePath
	}
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

	sensitive := map[string]bool{
		"top":     false,
		"up":      false,
		"down":    false,
		"bottom":  false,
		"playnow": false,
	}

	if len(values) == 1 {
		var pos, lastpos int
		if oc.id > 0 {
			pos = values[0][store.TColPosition].(int)
			lastpos = values[0][store.TColLastPosition].(int)
		} else {
			pos = values[0][store.QColPosition].(int)
			lastpos = values[0][store.QColLastPosition].(int)
		}

		sensitive["top"] = !(pos == 1)
		sensitive["up"] = !(pos == 1)
		sensitive["down"] = !(pos == lastpos)
		sensitive["bottom"] = !(pos == lastpos)
		sensitive["playnow"] = true
	}
	oc.ctxMenu.GetChildren().Foreach(func(item interface{}) {
		w, ok := item.(*gtk.Widget)
		if !ok {
			return
		}

		l, _ := w.GetName()
		for k, v := range sensitive {
			if strings.Contains(l, k) {
				w.SetSensitive(v)
				break
			}
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

		onerror.Log(dialer.ExecutePlaylistTrackAction(req))
		return
	}

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Perspective: oc.perspective,
		Action:      m3uetcpb.QueueAction_Q_CLEAR,
	}

	onerror.Log(dialer.ExecuteQueueAction(req))
}

func (oc *onContext) ContextDelete(mi *gtk.MenuItem) {
	values := oc.getSelection()
	if len(values) == 0 {
		return
	}

	oc.DeleteRows(values)
}

func (oc *onContext) ContextEnqueue(mi *gtk.MenuItem) {
	values := oc.getSelection()
	if len(values) == 0 {
		return
	}

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: m3uetcpb.QueueAction_Q_APPEND,
	}

	if oc.id > 0 {
		for _, m := range values {
			req.Ids = append(req.Ids, m[store.TColTrackID].(int64))
		}
	} else {
		for _, m := range values {
			req.Locations = append(req.Locations, m[store.QColLocation].(string))
		}
	}

	if err := dialer.ExecuteQueueAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (oc *onContext) ContextHide(m *gtk.Menu) {
	oc.resetSelection()
}

func (oc *onContext) ContextMove(mi *gtk.MenuItem) {
	values := oc.getSelection()
	if len(values) != 1 {
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
	fromPos := values[0][rowPosition].(int)
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
		pos = values[0][lastPosition].(int)
		if fromPos == pos {
			return
		}
	} else {
		log.Error("Invalid playlist/queue move")
		return
	}

	if oc.id > 0 {
		req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
			PlaylistId:   oc.id,
			Action:       m3uetcpb.PlaylistTrackAction_PT_MOVE,
			Position:     int32(pos),
			FromPosition: int32(fromPos),
		}

		onerror.Log(dialer.ExecutePlaylistTrackAction(req))
		return
	}

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action:       m3uetcpb.QueueAction_Q_MOVE,
		Position:     int32(pos),
		FromPosition: int32(fromPos),
	}

	onerror.Log(dialer.ExecuteQueueAction(req))
}

func (oc *onContext) ContextPlayNow(mi *gtk.MenuItem) {
	values := oc.getSelection()
	if len(values) != 1 {
		return
	}

	if oc.id > 0 {
		pos := values[0][store.TColPosition].(int)

		req := &m3uetcpb.ExecutePlaybarActionRequest{
			Action:   m3uetcpb.PlaybarAction_BAR_ACTIVATE,
			Position: int32(pos),
			Ids:      []int64{oc.id},
		}

		onerror.Log(dialer.ExecutePlaybarAction(req))
		return
	}

	id := values[0][store.QColTrackID].(int64)
	loc := values[0][store.QColLocation].(string)
	pos := values[0][store.QColPosition].(int)

	reqbar := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_PLAY,
		Force:  true,
	}
	if id > 0 {
		reqbar.Ids = []int64{id}
	} else {
		reqbar.Locations = []string{loc}
	}

	if err := dialer.ExecutePlaybackAction(reqbar); err != nil {
		log.Error(err)
		return
	}

	reqq := &m3uetcpb.ExecuteQueueActionRequest{
		Action:   m3uetcpb.QueueAction_Q_DELETE,
		Position: int32(pos),
	}

	onerror.Log(dialer.ExecuteQueueAction(reqq))
}

func (oc *onContext) ContextPoppedUp(m *gtk.Menu) {
	glib.IdleAdd(func() bool {
		sel, _ := oc.view.GetSelection()
		if oc.selection.keep && oc.selection.paths != nil {
			for i := range oc.selection.paths {
				if sel.PathIsSelected(oc.selection.paths[i]) {
					continue
				}
				sel.SelectPath(oc.selection.paths[i])
			}
		}
		return false
	})
}

func (oc *onContext) DeleteRows(values []map[store.ModelColumn]interface{}) {
	colPosition := store.QColPosition
	if oc.id > 0 {
		colPosition = store.TColPosition
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i][colPosition].(int) <
			values[j][colPosition].(int)
	})

	if oc.id > 0 {
		for i := len(values) - 1; i >= 0; i-- {
			req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
				PlaylistId: oc.id,
				Action:     m3uetcpb.PlaylistTrackAction_PT_DELETE,
				Position:   int32(values[i][store.TColPosition].(int)),
			}

			onerror.Log(dialer.ExecutePlaylistTrackAction(req))
		}
		return
	}

	for i := len(values) - 1; i >= 0; i-- {
		req := &m3uetcpb.ExecuteQueueActionRequest{
			Action:   m3uetcpb.QueueAction_Q_DELETE,
			Position: int32(values[i][store.QColPosition].(int)),
		}

		onerror.Log(dialer.ExecuteQueueAction(req))
	}
}

func (oc *onContext) Key(tv *gtk.TreeView, event *gdk.Event) {
	key := gdk.EventKeyNewFromEvent(event)
	if key.KeyVal() != gdk.KEY_Delete {
		return
	}

	values := oc.getSelection(true)
	if len(values) == 0 {
		return
	}

	oc.DeleteRows(values)
}

func (oc *onContext) SelChanged(sel *gtk.TreeSelection) {
	if oc.selection.keep {
		log.Info("Ignoring selection change")
		return
	}
	var err error
	cols := []store.ModelColumn{
		store.QColPosition,
		store.QColTrackID,
		store.QColLocation,
		store.QColLastPosition,
	}
	if oc.id > 0 {
		cols = []store.ModelColumn{
			store.TColPosition,
			store.TColTrackID,
			store.TColLocation,
			store.TColLastPosition,
		}
	}

	oc.selection.values, oc.selection.paths, err = store.GetMultipleTreeSelectionValues(sel, oc.view, cols)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("Selected context entres: %v", oc.selection)
}

func (oc *onContext) getSelection(keep ...bool) (
	values []map[store.ModelColumn]interface{}) {

	if oc.selection.values == nil {
		sel, err := oc.view.GetSelection()
		if err != nil {
			log.Error(err)
			return
		}
		oc.SelChanged(sel)
	}

	values, ok := oc.selection.values.([]map[store.ModelColumn]interface{})
	if !ok {
		log.Debug("There is no selection available for context")
		values = []map[store.ModelColumn]interface{}{}
		return
	}

	oc.selection.keep = false
	if len(keep) > 0 {
		oc.selection.keep = keep[0]
	}

	if !oc.selection.keep {
		oc.resetSelection()
	}
	return
}

func (oc *onContext) resetSelection() {
	oc.selection.keep = false
	oc.selection.values = nil
	oc.selection.paths = nil
}
