package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	log "github.com/sirupsen/logrus"
)

type queryData struct {
	subscriptionID string
	Query          []*m3uetcpb.Query
	tracks         []*m3uetcpb.Track

	mu sync.Mutex
}

var (
	// QYData query store
	QYData = &queryData{}

	queryResultsModel *gtk.ListStore

	queryTree queryTreeModel
)

// GetQuery returns the query for the gven id
func (qyd *queryData) GetQuery(id int64) *m3uetcpb.Query {
	qyd.mu.Lock()
	defer qyd.mu.Unlock()

	for _, v := range qyd.Query {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (qyd *queryData) GetSubscriptionID() string {
	qyd.mu.Lock()
	defer qyd.mu.Unlock()

	return qyd.subscriptionID
}

func (qyd *queryData) ProcessSubscriptionResponse(
	res *m3uetcpb.SubscribeToQueryStoreResponse) {

	appendItem := func(res *m3uetcpb.SubscribeToQueryStoreResponse) {
		for _, qy := range qyd.Query {
			if qy.Id == res.Query.Id {
				return
			}
		}
		qyd.Query = append(
			qyd.Query,
			res.Query,
		)
	}

	changeItem := func(res *m3uetcpb.SubscribeToQueryStoreResponse) {
		qy := res.Query
		for i := range qyd.Query {
			if qyd.Query[i].Id == qy.Id {
				qyd.Query[i] = qy
				break
			}
		}
	}

	removeItem := func(res *m3uetcpb.SubscribeToQueryStoreResponse) {
		n := len(qyd.Query)
		for i := range qyd.Query {
			if qyd.Query[i].Id == res.Query.Id {
				qyd.Query[i] = qyd.Query[n-1]
				qyd.Query = qyd.Query[:n-1]
				break
			}
		}
	}

	qyd.mu.Lock()
	defer qyd.mu.Unlock()

	if qyd.subscriptionID == "" {
		qyd.subscriptionID = res.SubscriptionId
	}

	switch res.Event {
	case m3uetcpb.QueryEvent_QYE_INITIAL:
		queryTree.initialMode = true
		qyd.Query = []*m3uetcpb.Query{}
	case m3uetcpb.QueryEvent_QYE_INITIAL_ITEM:
		appendItem(res)
	case m3uetcpb.QueryEvent_QYE_INITIAL_DONE:
		queryTree.initialMode = false
	case m3uetcpb.QueryEvent_QYE_ITEM_ADDED:
		appendItem(res)
	case m3uetcpb.QueryEvent_QYE_ITEM_CHANGED:
		changeItem(res)
	case m3uetcpb.QueryEvent_QYE_ITEM_REMOVED:
		removeItem(res)
	}

	if !queryTree.initialMode {
		glib.IdleAdd(queryTree.update)
	}
}

func (qyd *queryData) UpdateQueryByResults(res *m3uetcpb.QueryByResponse) int {
	qyd.mu.Lock()
	defer qyd.mu.Unlock()

	qyd.tracks = res.Tracks
	count := len(res.Tracks)

	glib.IdleAdd(qyd.updateQueryResults)
	return count
}

func (qyd *queryData) updateQueryResults() bool {
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

	qyd.mu.Lock()
	defer qyd.mu.Unlock()

	var iter *gtk.TreeIter
	for i, t := range qyd.tracks {
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
	return false
}

// ClearQueryResults -
func ClearQueryResults() {
	model := queryResultsModel

	_, ok := model.GetIterFirst()
	if ok {
		model.Clear()
	}
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
		values, err = GetListStoreModelValues(
			model,
			iter,
			[]ModelColumn{TColTrackID, TColToggleSelect},
		)
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

// ToggleQueryResultsSelection inverts the query results selection
func ToggleQueryResultsSelection() {
	log.Info("Toggling query results selection")

	model := queryResultsModel

	iter, ok := model.GetIterFirst()

	for ok {
		gval, err := model.GetValue(iter, int(TColToggleSelect))
		if err != nil {
			log.Error(err)
			return
		}

		value, err := gval.GoValue()
		if err != nil {
			log.Error(err)
			return
		}

		model.SetValue(iter, int(TColToggleSelect), !value.(bool))

		ok = model.IterNext(iter)
	}
}
