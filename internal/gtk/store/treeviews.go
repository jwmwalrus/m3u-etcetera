package store

import (
	"errors"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

func GetTreeSelectionValue(sel *gtk.TreeSelection, col ModelColumn) (value interface{}, err error) {
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
	}
	return
}

func GetTreeSelectionValues(sel *gtk.TreeSelection, cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {
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

func GetListStoreModelValue(model *gtk.ListStore, iter *gtk.TreeIter, col ModelColumn) (value interface{}, err error) {
	gval, err := model.GetValue(iter, int(col))
	if err != nil {
		return
	}
	value, err = gval.GoValue()
	return
}

func GetListStoreModelValues(model *gtk.ListStore, iter *gtk.TreeIter, cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {
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

func GetListStoreValue(tv *gtk.TreeView, path *gtk.TreePath, col ModelColumn) (value interface{}, err error) {
	imodel, err := tv.GetModel()
	if err != nil {
		return
	}
	model, ok := imodel.(*gtk.ListStore)
	if !ok {
		err = errors.New("Unable to get model from treeview")
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
	return
}

func GetListStoreValues(tv *gtk.TreeView, path *gtk.TreePath, cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {
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

func GetTreeStoreValue(tv *gtk.TreeView, path *gtk.TreePath, col ModelColumn) (value interface{}, err error) {
	imodel, err := tv.GetModel()
	if err != nil {
		return
	}
	model, ok := imodel.(*gtk.TreeStore)
	if !ok {
		err = errors.New("Unable to get model from treeview")
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
	return
}

func GetTreeStoreValues(tv *gtk.TreeView, path *gtk.TreePath, cols []ModelColumn) (values map[ModelColumn]interface{}, err error) {
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
