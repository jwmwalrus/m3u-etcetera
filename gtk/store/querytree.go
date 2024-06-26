package store

import (
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/pkg/qparams"
)

type queryTreeModel struct {
	model       *gtk.TreeStore
	filterVal   string
	initialMode bool

	mu sync.RWMutex
}

func (qyt *queryTreeModel) update() bool {
	slog.Info("Updating query model")

	qyt.mu.Lock()
	defer qyt.mu.Unlock()

	model := qyt.model
	if model == nil {
		return false
	}

	if model.NColumns() == 0 {
		return false
	}

	type queryInfo struct {
		id         int64
		name       string
		kw         string
		hasCBounds bool
		sort       int
	}

	_, ok := model.IterFirst()
	if ok {
		model.Clear()
	}

	getKeywords := func(qy *m3uetcpb.Query) string {
		list := strings.Split(qy.Name, " ")
		if qy.Description != "" {
			list = append(
				list,
				strings.Split(strings.ToLower(qy.Description), " ")...,
			)
		}
		if qy.Params != "" {
			if qp, err := qparams.ParseParams(qy.Params); err == nil {
				for _, x := range qp {
					list = append(
						list,
						strings.Split(strings.ToLower(x.Val), " ")...,
					)
				}
			}
		}
		return strings.Join(list, ",")
	}

	all := []queryInfo{}

	QYData.mu.RLock()
	for i, qy := range QYData.Query {
		if qyt.filterVal != "" {
			kw := getKeywords(qy)
			match := false
			for _, s := range strings.Split(qyt.filterVal, " ") {
				match = match || strings.Contains(kw, s)
				if match {
					break
				}
			}
			if !match {
				continue
			}
		}

		name := qy.Name
		sortVal := i
		if qy.ReadOnly {
			name = qy.Description
			sortVal -= 1000
		}

		qi := queryInfo{id: qy.Id, name: name, kw: getKeywords(qy), sort: sortVal}
		if len(qy.CollectionIds) > 0 {
			qi.hasCBounds = true
		}
		all = append(all, qi)
	}
	QYData.mu.RUnlock()

	sort.SliceStable(all, func(i, j int) bool {
		return all[i].name < all[j].name
	})

	for _, qi := range all {
		iter := model.Append(nil)
		name := qi.name
		if qi.hasCBounds {
			name += " (C)"
		}
		model.SetValue(iter, int(QYColTree), glib.NewValue(name))
		model.SetValue(iter, int(QYColTreeIDList), glib.NewValue(strconv.FormatInt(qi.id, 10)))
		model.SetValue(iter, int(QYColTreeKeywords), glib.NewValue(qi.kw))
		model.SetValue(iter, int(QYColTreeSort), glib.NewValue(qi.sort))
	}

	model.SetSortColumnID(int(QYColTreeSort), gtk.SortAscending)

	return false
}

// CreateQueryTreeModel creates a query model.
func CreateQueryTreeModel() (model *gtk.TreeStore, err error) {
	slog.Info("Creating query model")

	queryTree.model = gtk.NewTreeStore(QYTreeColumn.getTypes())
	if queryTree.model == nil {
		err = fmt.Errorf("failed to create tree-store")
		return
	}

	model = queryTree.model
	return
}

// FilterQueryTreeBy filters the queryTree by the given value.
func FilterQueryTreeBy(val string) {
	queryTree.filterVal = val
	queryTree.update()
}

// GetQueryTreeModel returns the current query model.
func GetQueryTreeModel() *gtk.TreeStore {
	return queryTree.model
}
