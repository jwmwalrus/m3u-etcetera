package musicpane

import (
	"log/slog"

	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/gtk/util"
)

type contextType int

const (
	collectionContext contextType = iota
	playlistContext
	queryContext
)

type onContext struct {
	ct        contextType
	view      *gtk.TreeView
	selection interface{}
}

func (oc *onContext) getSelection(keep ...bool) (ids []int64) {
	if oc.selection == nil {
		sel := oc.view.Selection()
		if sel == nil {
			slog.Error("Failed to get selection")
			return
		}
		oc.selChanged(sel)
	}

	value, ok := oc.selection.(string)
	if !ok {
		slog.Debug("There is no selection available for context")
		return
	}

	ids, err := util.StringToIDList(value)
	if err != nil {
		slog.Error("Failed to parse selection IDs for context", "error", err)
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

func (oc *onContext) getSelectionValues(keep ...bool) (
	values map[store.ModelColumn]interface{}) {

	if oc.selection == nil {
		sel := oc.view.Selection()
		if sel == nil {
			slog.Error("Failed to get selection")
			return
		}
		oc.selChanged(sel)
	}

	values, ok := oc.selection.(map[store.ModelColumn]interface{})
	if !ok {
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

func (oc *onContext) selChanged(sel *gtk.TreeSelection) {
	var err error
	switch oc.ct {
	case collectionContext:
		oc.selection, err = store.GetSingleTreeSelectionValue(
			sel,
			store.CColTreeIDList,
		)
	case playlistContext:
		oc.selection, err = store.GetSingleTreeSelectionValues(
			sel,
			[]store.ModelColumn{store.PLColTreeIDList, store.PLColTreeIsGroup},
		)
	case queryContext:
		oc.selection, err = store.GetSingleTreeSelectionValue(
			sel,
			store.QYColTreeIDList,
		)
	}
	if err != nil {
		slog.Error("Failed to get single tree selection", "error", err)
		return
	}
	slog.Debug("Selected context entry", "entry", oc.selection)
}
