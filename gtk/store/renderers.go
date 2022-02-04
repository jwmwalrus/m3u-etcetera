package store

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

// GetCollectionRenderer returns a renderer for any editable collection column
func GetCollectionRenderer(col ModelColumn) (gtk.ICellRenderer, error) {

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
			onTextColumnEdited(collectionModel, col, cell, pathString, newText)
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
			onBoolColumnToggled(collectionModel, col, cell, pathString)
		})

		return renderer, nil
	}

	return nil, fmt.Errorf("The provided column is not editable: %v", col)
}

// GetPlaylistGroupRenderer returns a renderer for any editable collection column
func GetPlaylistGroupRenderer(col ModelColumn) (gtk.ICellRenderer, error) {

	switch col {
	case PGColName, PGColDescription:
		renderer, err := gtk.CellRendererTextNew()
		if err != nil {
			return nil, err
		}
		err = renderer.Set("editable", true)
		if err != nil {
			return nil, err
		}
		renderer.Connect("edited", func(cell *gtk.CellRendererText, pathString, newText string) {
			onTextColumnEdited(playlistGroupsModel, col, cell, pathString, newText)
		})
		return renderer, nil
	}

	return nil, fmt.Errorf("The provided column is not editable: %v", col)
}

// GetQueryResultsRenderer returns a renderer for any editable query column
func GetQueryResultsRenderer(col ModelColumn) (gtk.ICellRenderer, error) {

	switch col {
	case TColToggleSelect:
		renderer, err := gtk.CellRendererToggleNew()
		if err != nil {
			return nil, err
		}
		err = renderer.Set("activatable", true)
		if err != nil {
			return nil, err
		}
		renderer.Connect("toggled", func(cell *gtk.CellRendererToggle, pathString string) {
			onBoolColumnToggled(queryResultsModel, col, cell, pathString)
		})

		return renderer, nil
	}

	return nil, fmt.Errorf("The provided column is not editable: %v", col)
}

func onBoolColumnToggled(model *gtk.ListStore, col ModelColumn, cell *gtk.CellRendererToggle, pathString string) {
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

func onTextColumnEdited(model *gtk.ListStore, col ModelColumn, cell *gtk.CellRendererText, pathString, newText string) {
	iter, err := model.GetIterFromString(pathString)
	if err != nil {
		log.Error(err)
		return
	}

	model.SetValue(iter, int(col), newText)
}
