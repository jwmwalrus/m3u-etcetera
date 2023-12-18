package store

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
)

// Renderer defines a gtk.ICellRenderer generator.
type Renderer struct {
	Model   *gtk.ListStore
	Columns storeColumns
}

// GetActivatable returns an activatable cell renderer.
func (r *Renderer) GetActivatable(col ModelColumn) (gtk.CellRendererer, error) {
	if !slices.Contains(r.Columns.GetActivatableColumns(), col) {
		return nil, fmt.Errorf("The provided column is not activatable: %v", col)
	}

	renderer := gtk.NewCellRendererToggle()
	renderer.SetActivatable(true)
	renderer.Connect(
		"toggled",
		func(cell *gtk.CellRendererToggle, pathString string) {
			onBoolColumnToggled(r.Model, col, cell, pathString)
		},
	)

	return renderer, nil
}

// GetEditable returns an editable cell renderer.
func (r *Renderer) GetEditable(col ModelColumn) (gtk.CellRendererer, error) {
	if !slices.Contains(r.Columns.GetEditableColumns(), col) {
		return nil, fmt.Errorf("The provided column is not editable: %v", col)
	}

	renderer := gtk.NewCellRendererText()
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

	iter, ok := model.IterFromString(pathString)
	if !ok {
		slog.Error("failed to get iter from string")
		return
	}

	gval := model.Value(iter, int(col))
	value := gval.GoValue()
	if value == nil {
		slog.Error("Failed to get Go value")
		return
	}

	model.SetValue(iter, int(col), glib.NewValue(!value.(bool)))
}

func onTextColumnEdited(model *gtk.ListStore, col ModelColumn,
	cell *gtk.CellRendererText, pathString, newText string) {

	iter, ok := model.IterFromString(pathString)
	if !ok {
		slog.Error("Failed to get iter from string")
		return
	}

	model.SetValue(iter, int(col), glib.NewValue(newText))
}
