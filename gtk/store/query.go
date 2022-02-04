package store

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/pkg/qparams"
	log "github.com/sirupsen/logrus"
)

type queryTreeModel struct {
	model       *gtk.TreeStore
	filterVal   string
	initialMode bool
}

var (
	queryTree         queryTreeModel
	queryResultsModel *gtk.ListStore

	// QYData query store
	QYData struct {
		subscriptionID string
		Mu             sync.Mutex
		Query          []*m3uetcpb.Query
		tracks         []*m3uetcpb.Track
	}
)

// ClearQueryResults -
func ClearQueryResults() {
	model := queryResultsModel

	_, ok := model.GetIterFirst()
	if ok {
		model.Clear()
	}
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

// CreateQueryResultsModel creates a query model
func CreateQueryResultsModel() (model *gtk.ListStore, err error) {
	log.Info("Creating query model")

	queryResultsModel, err = gtk.ListStoreNew(TColumns.getTypes()...)
	if err != nil {
		return
	}

	model = queryResultsModel
	return
}

// FilterQueryTreeBy filters the queryTree by the given value
func FilterQueryTreeBy(val string) {
	queryTree.filterVal = val
	queryTree.update()
}

// GetQuery returns the query for the gven id
func GetQuery(id int64) *m3uetcpb.Query {
	QYData.Mu.Lock()
	defer QYData.Mu.Unlock()

	for _, v := range QYData.Query {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// GetQueryTreeModel returns the current query model
func GetQueryTreeModel() *gtk.TreeStore {
	return queryTree.model
}

// GetQueryResultsSelections returns the list of selected query results
func GetQueryResultsSelections() (ids []int64, err error) {
	model := queryResultsModel
	if model == nil {
		return
	}

	if model.GetNColumns() == 0 {
		return
	}

	iter, ok := model.GetIterFirst()
	for ok {
		var values map[ModelColumn]interface{}
		values, err = GetListStoreModelValues(model, iter, []ModelColumn{TColTrackID, TColToggleSelect})
		if err != nil {
			log.Error(err)
			return
		}
		selected := values[TColToggleSelect].(bool)
		if selected {
			ids = append(ids, values[TColTrackID].(int64))
		}
		ok = model.IterNext(iter)
	}
	return
}

func (qyt *queryTreeModel) update() bool {
	log.Info("Updating query model")

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
			list = append(list, strings.Split(strings.ToLower(qy.Description), " ")...)
		}
		if qy.Params != "" {
			if qp, err := qparams.ParseParams(qy.Params); err == nil {
				for _, x := range qp {
					list = append(list, strings.Split(strings.ToLower(x.Val), " ")...)
				}
			}
		}
		return strings.Join(list, ",")
	}

	all := []queryInfo{}

	QYData.Mu.Lock()
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
	QYData.Mu.Unlock()

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

func updateQueryResults() bool {
	log.Info("Updating query results")

	model := queryResultsModel
	if model == nil {
		return false
	}
	if model.GetNColumns() == 0 {
		return false
	}

	_, ok := model.GetIterFirst()
	if ok {
		model.Clear()
	}

	QYData.Mu.Lock()
	var iter *gtk.TreeIter
	for i, t := range QYData.tracks {
		iter = model.Append()
		dur := time.Duration(t.Duration) * time.Nanosecond
		err := model.Set(
			iter,
			[]int{
				int(TColTrackID),
				int(TColCollectionID),
				int(TColFormat),
				int(TColType),
				int(TColTitle),
				int(TColAlbum),
				int(TColArtist),
				int(TColAlbumartist),
				int(TColComposer),
				int(TColGenre),

				int(TColYear),
				int(TColTracknumber),
				int(TColTracktotal),
				int(TColDiscnumber),
				int(TColDisctotal),
				int(TColLyrics),
				int(TColComment),
				int(TColPlaycount),

				int(TColRating),
				int(TColDuration),
				int(TColRemote),
				int(TColLastplayed),
				int(TColNumber),
				int(TColToggleSelect),
			},
			[]interface{}{
				t.Id,
				t.CollectionId,
				t.Format,
				t.Type,
				t.Title,
				t.Album,
				t.Artist,
				t.Albumartist,
				t.Composer,
				t.Genre,

				int(t.Year),
				int(t.Tracknumber),
				int(t.Tracktotal),
				int(t.Discnumber),
				int(t.Disctotal),
				t.Lyrics,
				t.Comment,
				int(t.Playcount),

				int(t.Rating),
				fmt.Sprint(dur.Truncate(time.Second)),
				t.Remote,
				t.Lastplayed,
				i + 1,
				false,
			},
		)
		if err != nil {
			log.Error(err)
			return false
		}
	}
	QYData.Mu.Unlock()
	return false
}