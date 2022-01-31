package musicpane

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
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
	selection interface{}
}

func (oc *onContext) getSelection(keep ...bool) (ids []int64) {
	value, ok := oc.selection.(string)
	if !ok {
		log.Debug("There is no selection available for collection context")
		return
	}

	ids, err := store.StringToIDList(value)
	if err != nil {
		log.Errorf("Error parsing selection value for collection context: %v", err)
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
		oc.selection, err = store.GetTreeSelectionValue(sel, store.PLColTreeIDList)
	case queryContext:
		oc.selection, err = store.GetTreeSelectionValue(sel, store.QYColTreeIDList)
	}
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Selected context entry: %v", oc.selection)
}
