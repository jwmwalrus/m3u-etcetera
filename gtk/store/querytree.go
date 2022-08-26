package store

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/pkg/qparams"
	log "github.com/sirupsen/logrus"
)

type queryTreeModel struct {
	model       *gtk.TreeStore
	filterVal   string
	initialMode bool

	mu sync.Mutex
}

func (qyt *queryTreeModel) update() bool {
	log.Info("Updating query model")

	qyt.mu.Lock()
	defer qyt.mu.Unlock()

	model := qyt.model
	if model == nil {
		return false
	}

	if model.GetNColumns() == 0 {
		return false
	}

	type queryInfo struct {
		id         int64
		name       string
		kw         string
		hasCBounds bool
	}

	_, ok := model.GetIterFirst()
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

	QYData.mu.Lock()
	for _, qy := range QYData.Query {
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

		qi := queryInfo{id: qy.Id, name: qy.Name, kw: getKeywords(qy)}
		if len(qy.CollectionIds) > 0 {
			qi.hasCBounds = true
		}
		all = append(all, qi)
	}
	QYData.mu.Unlock()

	sort.SliceStable(all, func(i, j int) bool {
		return all[i].name < all[j].name
	})

	for _, qi := range all {
		iter := model.Append(nil)
		name := qi.name
		if qi.hasCBounds {
			name += " (C)"
		}
		model.SetValue(iter, int(QYColTree), qi.name)
		model.SetValue(iter, int(QYColTreeIDList), strconv.FormatInt(qi.id, 10))
		model.SetValue(iter, int(QYColTreeKeywords), qi.kw)
	}
	return false
}

// CreateQueryTreeModel creates a query model
func CreateQueryTreeModel() (model *gtk.TreeStore, err error) {
	log.Info("Creating query model")

	queryTree.model, err = gtk.TreeStoreNew(QYTreeColumn.getTypes()...)
	if err != nil {
		return
	}

	model = queryTree.model
	return
}

// FilterQueryTreeBy filters the queryTree by the given value
func FilterQueryTreeBy(val string) {
	queryTree.filterVal = val
	queryTree.update()
}

// GetQueryTreeModel returns the current query model
func GetQueryTreeModel() *gtk.TreeStore {
	return queryTree.model
}
