package gtkui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/store"
	log "github.com/sirupsen/logrus"
)

func setupCollectionsDialog() (dlg *gtk.Dialog, err error) {
	log.Info("Setting up collections dialog")

	dlg, err = builder.GetDialog("collections_dialog")
	if err != nil {
		log.Error(err)
		return
	}

	view, err := builder.GetTreeView("collections_dialog_view")
	if err != nil {
		log.Error(err)
		return
	}

	textro, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	togglero, err := gtk.CellRendererToggleNew()
	if err != nil {
		return
	}

	namerw, err := store.GetCollectionRenderer(store.CColName)
	if err != nil {
		return
	}

	descriptionrw, err := store.GetCollectionRenderer(store.CColDescription)
	if err != nil {
		return
	}

	remotelocationrw, err := store.GetCollectionRenderer(store.CColRemoteLocation)
	if err != nil {
		return
	}

	disabledrw, err := store.GetCollectionRenderer(store.CColDisabled)
	if err != nil {
		return
	}

	remoterw, err := store.GetCollectionRenderer(store.CColRemote)
	if err != nil {
		return
	}

	rescanrw, err := store.GetCollectionRenderer(store.CColRescan)
	if err != nil {
		return
	}

	cols := []struct {
		idx store.StoreModelColumn
		r   gtk.ICellRenderer
	}{
		{store.CColName, namerw},
		{store.CColDescription, descriptionrw},
		{store.CColLocation, textro},
		{store.CColScanned, textro},
		{store.CColTracksView, textro},
		{store.CColHidden, togglero},
		{store.CColDisabled, disabledrw},
		{store.CColRemote, remoterw},
		{store.CColRescan, rescanrw},
		{store.CColRemoteLocation, remotelocationrw},
	}

	for _, v := range cols {
		var col *gtk.TreeViewColumn
		if renderer, ok := v.r.(*gtk.CellRendererToggle); ok {
			col, err = gtk.TreeViewColumnNewWithAttribute(
				store.CColumns[v.idx].Name,
				renderer,
				"active",
				int(v.idx),
			)
		} else if renderer, ok := v.r.(*gtk.CellRendererText); ok {
			col, err = gtk.TreeViewColumnNewWithAttribute(
				store.CColumns[v.idx].Name,
				renderer,
				"text",
				int(v.idx),
			)
		} else {
			log.Error("¿Cómo sabré si es pez o iguana?")
			continue
		}
		if err != nil {
			return
		}
		view.InsertColumn(col, -1)
	}

	model, err := store.CreateCollectionsModel()
	if err != nil {
		return
	}
	view.SetModel(model)

	return
}
