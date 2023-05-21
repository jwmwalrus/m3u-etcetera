package store

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

// Renderer defines a gtk.ICellRenderer generator.
type Renderer struct {
	Model   *gtk.ListStore
	Columns storeColumns
}

// GetActivatable returns an activatable cell renderer.
func (r *Renderer) GetActivatable(col ModelColumn) (gtk.ICellRenderer, error) {
	if !slices.Contains(r.Columns.GetActivatableColumns(), col) {
		return nil, fmt.Errorf("The provided column is not activatable: %v", col)
	}

	renderer, err := gtk.CellRendererToggleNew()
	if err != nil {
		return nil, err
	}
	err = renderer.Set("activatable", true)
	if err != nil {
		return nil, err
	}
	renderer.Connect(
		"toggled",
		func(cell *gtk.CellRendererToggle, pathString string) {
			onBoolColumnToggled(r.Model, col, cell, pathString)
		},
	)

	return renderer, nil
}

// GetEditable returns an editable cell renderer.
func (r *Renderer) GetEditable(col ModelColumn) (gtk.ICellRenderer, error) {
	if !slices.Contains(r.Columns.GetEditableColumns(), col) {
		return nil, fmt.Errorf("The provided column is not editable: %v", col)
	}

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}
	err = renderer.Set("editable", true)
	if err != nil {
		return nil, err
	}
	renderer.Connect(
		"edited",
		func(cell *gtk.CellRendererText, pathString, newText string) {
			onTextColumnEdited(r.Model, col, cell, pathString, newText)
		},
	)
	return renderer, nil
}

func onBoolColumnToggled(model *gtk.ListStore, col ModelColumn,
	cell *gtk.CellRendererToggle, pathString string) {

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

func onTextColumnEdited(model *gtk.ListStore, col ModelColumn,
	cell *gtk.CellRendererText, pathString, newText string) {

	iter, err := model.GetIterFromString(pathString)
	if err != nil {
		log.Error(err)
		return
	}

	model.SetValue(iter, int(col), newText)
}
