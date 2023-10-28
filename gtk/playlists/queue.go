package playlists

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/bnp/chars"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
)

// OnQueue handles queue signals.
type OnQueue struct {
	*onContext
}

// CreateQueue returns a queue signals handler.
func CreateQueue(p m3uetcpb.Perspective, queueID, contextMenuID string) (
	oq *OnQueue, err error) {

	logw := slog.With(
		"perspective", p,
		"queueID", queueID,
		"queueMenuID", contextMenuID,
	)
	logw.Info("Creating queue view and model")

	ctxMenu, err := builder.GetMenu(contextMenuID)
	if err != nil {
		logw.With(
			"menu-id", contextMenuID,
			"error", err,
		).Error("Failed to get context menu")
		return
	}

	miSuffix, _ := chars.GetRandomLetters(6)
	for _, l := range []string{"top", "up", "down", "bottom"} {
		mi, err := builder.GetMenuItem(contextMenuID + "_" + l)
		if err != nil {
			logw.With(
				"menu", contextMenuID,
				"item", l,
				"error", err,
			).Error("Failed to get menu item")
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
		col.SetSizing(gtk.TREE_VIEW_COLUMN_AUTOSIZE)
		col.SetResizable(true)
		oq.view.InsertColumn(col, -1)
	}

	model, err := store.CreateQueueModel(oq.perspective)
	if err != nil {
		return
	}
	oq.view.SetModel(model)
	return
}

// DblClicked handles double-click.
func (oq *OnQueue) DblClicked(tv *gtk.TreeView,
	path *gtk.TreePath, col *gtk.TreeViewColumn) {

	values, err := store.GetTreeViewTreePathValues(
		tv,
		path,
		[]store.ModelColumn{
			store.QColPosition,
			store.QColTrackID,
			store.QColLocation,
		},
	)
	if err != nil {
		slog.Error("Failed to get tree-view's tree-path values", "error", err)
		return
	}
	slog.Debug("Doouble-clicked column values", "values", values)

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

	if err := dialer.ExecutePlaybackAction(req); err != nil {
		slog.Error("Failed to execute playback action", "error", err)
		return
	}

	req2 := &m3uetcpb.ExecuteQueueActionRequest{
		Action:   m3uetcpb.QueueAction_Q_DELETE,
		Position: int32(pos),
	}

	if err := dialer.ExecuteQueueAction(req2); err != nil {
		slog.Error("Failed to execute queue action", "error", err)
		return
	}
}
