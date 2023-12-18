package store

import (
	"fmt"
	"log/slog"

	"github.com/diamondburned/gotk4/pkg/gtk/v3"
)

// GetMultipleTreeSelectionValues returns the values of the tree selection for the
// given columns.
func GetMultipleTreeSelectionValues(sel *gtk.TreeSelection, tv *gtk.TreeView, cols []ModelColumn) (
	values []map[ModelColumn]interface{}, paths []*gtk.TreePath, err error) {

	model, paths := sel.SelectedRows()
	for _, p := range paths {
		iter, ok := model.Iter(p)
		if !ok {
			slog.Error("Failed to get tree-iter")
			continue
		}

		var m map[ModelColumn]interface{}
		m, err = GetTreeModelValues(model, iter, cols)
		if err != nil {
			slog.Error("Failed to get tree-model", "error", err)
			return
		}

		values = append(values, m)
	}
	return
}

// GetSingleTreeSelectionValue returns the value of the tree selection for the given
// column.
func GetSingleTreeSelectionValue(sel *gtk.TreeSelection, col ModelColumn) (
	value interface{}, err error) {

	model, iter, ok := sel.Selected()
	slog.With(
		"ok", ok,
		"model-is-nil", model == nil,
		"iter-is-nil", iter == nil,
		"col", col,
	).Debug("GetSingleTreeSelectionValue")
	if ok {
		value, err = GetTreeModelValue(model, iter, col)
	}
	return
}

// GetSingleTreeSelectionValues returns the values of the tree selection for the
// given columns.
func GetSingleTreeSelectionValues(sel *gtk.TreeSelection, cols []ModelColumn) (
	values map[ModelColumn]interface{}, err error) {

	model, iter, ok := sel.Selected()
	slog.With(
		"ok", ok,
		"model-is-nil", model == nil,
		"iter-is-nil", iter == nil,
		"cols", cols,
	).Debug("GetSingleTreeSelectionValues")
	if ok {
		values, err = GetTreeModelValues(model, iter, cols)
	}
	return
}

// GetTreeModel given a gtk.ITreeModel returns the *gtk.TreeModel.
func GetTreeModel(imodel gtk.TreeModeller) *gtk.TreeModel {
	model, ok := imodel.(*gtk.TreeModel)
	if ok {
		return model
	}
	list, ok := imodel.(*gtk.ListStore)
	if ok {
		return &list.TreeModel
	}
	tree, ok := imodel.(*gtk.TreeStore)
	if ok {
		return &tree.TreeModel
	}
	return nil
}

// GetTreeModelValue returns the value of the tree model for the given
// column at the given tree-iter.
func GetTreeModelValue(model *gtk.TreeModel, iter *gtk.TreeIter,
	col ModelColumn) (value interface{}, err error) {

	gval := model.Value(iter, int(col))
	value = gval.GoValue()
	if value == nil {
		err = fmt.Errorf("failed to get tree-model-value")
		return
	}
	return
}

// GetTreeModelValues returns the values of the list store for the given
// columns at the given tree-iter.
func GetTreeModelValues(model *gtk.TreeModel, iter *gtk.TreeIter,
	cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {

	m := map[ModelColumn]interface{}{}
	for _, c := range cols {
		var v interface{}
		v, err = GetTreeModelValue(model, iter, c)
		if err != nil {
			return
		}
		m[c] = v
	}
	values = m
	return
}

// GetTreeViewTreePathValue returns the value of the list view for the given column at
// the given tree-path.
func GetTreeViewTreePathValue(tv *gtk.TreeView, path *gtk.TreePath,
	col ModelColumn) (value interface{}, err error) {

	model := tv.Model()
	if model == nil {
		err = fmt.Errorf("Unable to get model from treeview")
		return
	}

	iter, ok := model.Iter(path)
	if !ok {
		err = fmt.Errorf("failed to get model iterator")
		return
	}

	value, err = GetTreeModelValue(model, iter, col)
	return

}

// GetTreeViewTreePathValues returns the values of the list view for the given columns
// at the given tree-path.
func GetTreeViewTreePathValues(tv *gtk.TreeView, path *gtk.TreePath,
	cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {

	model := tv.Model()
	if model == nil {
		err = fmt.Errorf("failed to get model from treeview")
		return
	}

	iter, ok := model.Iter(path)
	if !ok {
		err = fmt.Errorf("failed to get model iterator")
		return
	}

	values, err = GetTreeModelValues(model, iter, cols)
	return
}
