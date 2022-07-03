package store

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

// GetTreeSelectionValue returns the value of the tree selection for the given
// column
func GetTreeSelectionValue(sel *gtk.TreeSelection, col ModelColumn) (
	value interface{}, err error) {

	model, iter, ok := sel.GetSelected()
	if ok {
		var gval *glib.Value

		gval, err = model.(*gtk.TreeModel).GetValue(iter, int(col))
		if err != nil {
			log.Error(err)
			return
		}

		value, err = gval.GoValue()
		if err != nil {
			log.Error(err)
			return
		}

		if value == nil {
			err = fmt.Errorf("Unable to get tree-selection value")
			return
		}
	}
	return
}

// GetTreeSelectionValues returns the values of the tree selection for the
// given columns
func GetTreeSelectionValues(sel *gtk.TreeSelection, cols []ModelColumn) (
	values map[ModelColumn]interface{}, err error) {

	m := map[ModelColumn]interface{}{}
	for _, c := range cols {
		var v interface{}
		v, err = GetTreeSelectionValue(sel, c)
		if err != nil {
			return
		}
		m[c] = v
	}
	values = m
	return
}

// GetListStoreModelValue returns the value of the list store for the given
// column at the given tree-iter
func GetListStoreModelValue(model *gtk.ListStore, iter *gtk.TreeIter,
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

// GetListStoreModelValues returns the values of the list store for the given
// columns at the given tree-iter
func GetListStoreModelValues(model *gtk.ListStore, iter *gtk.TreeIter,
	cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {

	m := map[ModelColumn]interface{}{}
	for _, c := range cols {
		var v interface{}
		v, err = GetListStoreModelValue(model, iter, c)
		if err != nil {
			return
		}
		m[c] = v
	}
	values = m
	return
}

// GetListStoreValue returns the value of the list view for the given column at
// the given tree-path
func GetListStoreValue(tv *gtk.TreeView, path *gtk.TreePath,
	col ModelColumn) (value interface{}, err error) {

	imodel, err := tv.GetModel()
	if err != nil {
		return
	}

	model, ok := imodel.(*gtk.ListStore)
	if !ok {
		err = fmt.Errorf("Unable to get model from treeview")
		return
	}

	iter, err := model.GetIter(path)
	if err != nil {
		log.Error(err)
		return
	}

	gval, err := model.GetValue(iter, int(col))
	if err != nil {
		return
	}

	value, err = gval.GoValue()
	if err != nil {
		return
	}

	if value == nil {
		err = fmt.Errorf("Unable to get list-store value")
		return
	}
	return
}

// GetListStoreValues returns the values of the list view for the given columns
// at the given tree-path
func GetListStoreValues(tv *gtk.TreeView, path *gtk.TreePath,
	cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {

	m := map[ModelColumn]interface{}{}
	for _, c := range cols {
		var v interface{}
		v, err = GetListStoreValue(tv, path, c)
		if err != nil {
			return
		}
		m[c] = v
	}
	values = m
	return
}

// GetTreeStoreValue returns the value of the list view for the given column at
// the given tree-path
func GetTreeStoreValue(tv *gtk.TreeView, path *gtk.TreePath,
	col ModelColumn) (value interface{}, err error) {

	imodel, err := tv.GetModel()
	if err != nil {
		return
	}

	model, ok := imodel.(*gtk.TreeStore)
	if !ok {
		err = fmt.Errorf("Unable to get model from treeview")
		return
	}

	iter, err := model.GetIter(path)
	if err != nil {
		log.Error(err)
		return
	}

	gval, err := model.GetValue(iter, int(col))
	if err != nil {
		return
	}

	value, err = gval.GoValue()
	if err != nil {
		return
	}

	if value == nil {
		err = fmt.Errorf("Unable to get tree-store value")
		return
	}
	return
}

// GetTreeStoreValues returns the values of the list view for the given columns
// at the given tree-path
func GetTreeStoreValues(tv *gtk.TreeView, path *gtk.TreePath,
	cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {

	m := map[ModelColumn]interface{}{}
	for _, c := range cols {
		var v interface{}
		v, err = GetTreeStoreValue(tv, path, c)
		if err != nil {
			return
		}
		m[c] = v
	}
	values = m
	return
}
