package store

import (
	"fmt"
	"log/slog"

	"github.com/gotk3/gotk3/gtk"
)

// GetMultipleTreeSelectionValues returns the values of the tree selection for the
// given columns.
func GetMultipleTreeSelectionValues(sel *gtk.TreeSelection, tv *gtk.TreeView, cols []ModelColumn) (
	values []map[ModelColumn]interface{}, paths []*gtk.TreePath, err error) {

	paths = []*gtk.TreePath{}

	imodel, err := tv.GetModel()
	if err != nil {
		return
	}
	model := GetTreeModel(imodel)

	glist := sel.GetSelectedRows(imodel)
	glist.Foreach(func(i interface{}) {
		p, ok := i.(*gtk.TreePath)
		if !ok {
			slog.Error("failed to get tree-path from interface")
			return
		}

		iter, err := model.GetIter(p)
		if err != nil {
			slog.Error("Failed to get tree-iter")
			return
		}
		paths = append(paths, p)

		var m map[ModelColumn]interface{}
		m, err = GetTreeModelValues(model, iter, cols)
		if err != nil {
			slog.Error("Failed to get tree-model", "error", err)
			return
		}

		values = append(values, m)
	})
	return
}

// GetSingleTreeSelectionValue returns the value of the tree selection for the given
// column.
func GetSingleTreeSelectionValue(sel *gtk.TreeSelection, col ModelColumn) (
	value interface{}, err error) {

	model, iter, ok := sel.GetSelected()
	if ok {
		value, err = GetTreeModelValue(model.(*gtk.TreeModel), iter, col)
	}
	return
}

// GetSingleTreeSelectionValues returns the values of the tree selection for the
// given columns.
func GetSingleTreeSelectionValues(sel *gtk.TreeSelection, cols []ModelColumn) (
	values map[ModelColumn]interface{}, err error) {

	model, iter, ok := sel.GetSelected()
	if ok {
		values, err = GetTreeModelValues(model.(*gtk.TreeModel), iter, cols)
	}
	return
}

// GetTreeModel given a gtk.ITreeModel returns the *gtk.TreeModel.
func GetTreeModel(imodel gtk.ITreeModel) *gtk.TreeModel {
	model, ok := imodel.(*gtk.TreeModel)
	if ok {
		return model
	}
	list, ok := imodel.(*gtk.ListStore)
	if ok {
		return list.ToTreeModel()
	}
	tree, ok := imodel.(*gtk.TreeStore)
	if ok {
		return tree.ToTreeModel()
	}
	return nil
}

// GetTreeModelValue returns the value of the tree model for the given
// column at the given tree-iter.
func GetTreeModelValue(model *gtk.TreeModel, iter *gtk.TreeIter,
	col ModelColumn) (value interface{}, err error) {

	gval, err := model.GetValue(iter, int(col))
	if err != nil {
		return
	}

	value, err = gval.GoValue()
	if err != nil {
		return
	}

	if value == nil {
		err = fmt.Errorf("Unable to get list-store-model value")
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

	imodel, err := tv.GetModel()
	if err != nil {
		return
	}

	model := GetTreeModel(imodel)
	if model == nil {
		err = fmt.Errorf("Unable to get model from treeview")
		return
	}

	iter, err := model.GetIter(path)
	if err != nil {
		return
	}

	value, err = GetTreeModelValue(model, iter, col)
	return

}

// GetTreeViewTreePathValues returns the values of the list view for the given columns
// at the given tree-path.
func GetTreeViewTreePathValues(tv *gtk.TreeView, path *gtk.TreePath,
	cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {

	imodel, err := tv.GetModel()
	if err != nil {
		return
	}

	model := GetTreeModel(imodel)
	if model == nil {
		err = fmt.Errorf("Unable to get model from treeview")
		return
	}

	iter, err := model.GetIter(path)
	if err != nil {
		return
	}

	values, err = GetTreeModelValues(model, iter, cols)
	return
}
