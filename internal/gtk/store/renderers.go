package store

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

// GetCollectionRenderer returns a renderer for any editable collection column
func GetCollectionRenderer(col StoreModelColumn) (gtk.ICellRenderer, error) {

	switch col {
	case CColName, CColDescription, CColRemoteLocation:
		renderer, err := gtk.CellRendererTextNew()
		if err != nil {
			return nil, err
		}
		err = renderer.Set("editable", true)
		if err != nil {
			return nil, err
		}
		renderer.Connect("edited", func(cell *gtk.CellRendererText, pathString, newText string) {
			onTextColumnEdited(col, cell, pathString, newText)
		})
		return renderer, nil

	case CColDisabled, CColRemote, CColScanned, CColRescan:
		renderer, err := gtk.CellRendererToggleNew()
		if err != nil {
			return nil, err
		}
		err = renderer.Set("activatable", true)
		if err != nil {
			return nil, err
		}
		renderer.Connect("toggled", func(cell *gtk.CellRendererToggle, pathString string) {
			onBoolColumnToggled(col, cell, pathString)
		})

		return renderer, nil
	}

	return nil, fmt.Errorf("The provided column is not editable: %v", col)
}

func onBoolColumnToggled(col StoreModelColumn, cell *gtk.CellRendererToggle, pathString string) {
	model := collectionsModel

	iter, err := model.GetIterFromString(pathString)
	if err != nil {
		log.Error(err)
		return
	}

	gval, err := model.GetValue(iter, int(col))
	if err != nil {
		log.Error(err)
		return
	}

	value, err := gval.GoValue()
	if err != nil {
		log.Error(err)
		return
	}

	model.SetValue(iter, int(col), !value.(bool))
}

func onTextColumnEdited(col StoreModelColumn, cell *gtk.CellRendererText, pathString, newText string) {
	model := collectionsModel

	iter, err := model.GetIterFromString(pathString)
	if err != nil {
		log.Error(err)
		return
	}

	model.SetValue(iter, int(col), newText)
}
