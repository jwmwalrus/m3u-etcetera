package musicpane

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/gtk/util"
	log "github.com/sirupsen/logrus"
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
		sel, err := oc.view.GetSelection()
		if err != nil {
			log.Error(err)
			return
		}
		oc.selChanged(sel)
	}

	value, ok := oc.selection.(string)
	if !ok {
		log.Debug("There is no selection available for context")
		return
	}

	ids, err := util.StringToIDList(value)
	if err != nil {
		log.Errorf("Error parsing selection IDs for context: %v", err)
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
		sel, err := oc.view.GetSelection()
		if err != nil {
			log.Error(err)
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
		oc.selection, err = store.GetTreeSelectionValue(sel, store.CColTreeIDList)
	case playlistContext:
		oc.selection, err = store.GetTreeSelectionValues(
			sel,
			[]store.ModelColumn{store.PLColTreeIDList, store.PLColTreeIsGroup},
		)
	case queryContext:
		oc.selection, err = store.GetTreeSelectionValue(sel, store.QYColTreeIDList)
	}
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Selected context entry: %v", oc.selection)
}
