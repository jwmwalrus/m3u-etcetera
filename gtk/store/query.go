package store

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
)

type queryData struct {
	subscriptionID string
	Query          []*m3uetcpb.Query
	tracks         []*m3uetcpb.Track

	mu sync.RWMutex
}

var (
	// QYData query store.
	QYData = &queryData{}

	queryResultsModel *gtk.ListStore

	queryTree queryTreeModel
)

// GetQuery returns the query for the gven id.
func (qyd *queryData) GetQuery(id int64) *m3uetcpb.Query {
	qyd.mu.RLock()
	defer qyd.mu.RUnlock()

	for _, v := range qyd.Query {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (qyd *queryData) SubscriptionID() string {
	qyd.mu.RLock()
	defer qyd.mu.RUnlock()

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
	slog.Info("Updating query results")

	model := queryResultsModel
	if model == nil {
		return false
	}
	if model.NColumns() == 0 {
		return false
	}

	_, ok := model.IterFirst()
	if ok {
		model.Clear()
	}

	qyd.mu.Lock()
	defer qyd.mu.Unlock()

	var iter *gtk.TreeIter
	for i, t := range qyd.tracks {
		iter = model.Append()
		dur := time.Duration(t.Duration) * time.Nanosecond
		model.Set(
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
			[]glib.Value{
				*glib.NewValue(t.Id),
				*glib.NewValue(t.CollectionId),
				*glib.NewValue(t.Format),
				*glib.NewValue(t.Type),
				*glib.NewValue(t.Title),
				*glib.NewValue(t.Album),
				*glib.NewValue(t.Artist),
				*glib.NewValue(t.Albumartist),
				*glib.NewValue(t.Composer),
				*glib.NewValue(t.Genre),

				*glib.NewValue(int(t.Year)),
				*glib.NewValue(int(t.Tracknumber)),
				*glib.NewValue(int(t.Tracktotal)),
				*glib.NewValue(int(t.Discnumber)),
				*glib.NewValue(int(t.Disctotal)),
				*glib.NewValue(t.Lyrics),
				*glib.NewValue(t.Comment),
				*glib.NewValue(int(t.Playcount)),

				*glib.NewValue(int(t.Rating)),
				*glib.NewValue(fmt.Sprint(dur.Truncate(time.Second))),
				*glib.NewValue(t.Remote),
				*glib.NewValue(time.Unix(0, t.Lastplayed).Format(lastPlayedLayout)),
				*glib.NewValue(i + 1),
				*glib.NewValue(false),
			},
		)
	}
	return false
}

// ClearQueryResults -.
func ClearQueryResults() {
	model := queryResultsModel

	_, ok := model.IterFirst()
	if ok {
		model.Clear()
	}
}

// CreateQueryResultsModel creates a query model.
func CreateQueryResultsModel() (model *gtk.ListStore, err error) {
	slog.Info("Creating query model")

	queryResultsModel = gtk.NewListStore(TColumns.getTypes())
	if queryResultsModel == nil {
		err = fmt.Errorf("failed to create list-store")
		return
	}

	model = queryResultsModel
	return
}

// GetQueryResultsSelections returns the list of selected query results.
func GetQueryResultsSelections() (ids []int64, err error) {
	model := queryResultsModel
	if model == nil {
		return
	}

	if model.NColumns() == 0 {
		return
	}

	iter, ok := model.IterFirst()
	for ok {
		var values map[ModelColumn]interface{}
		values, err = GetTreeModelValues(
			&model.TreeModel,
			iter,
			[]ModelColumn{TColTrackID, TColToggleSelect},
		)
		if err != nil {
			slog.Error("Failed to get tree-model values", "error", err)
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

// ToggleQueryResultsSelection inverts the query results selection.
func ToggleQueryResultsSelection() {
	slog.Info("Toggling query results selection")

	model := queryResultsModel

	iter, ok := model.IterFirst()

	for ok {
		gval := model.Value(iter, int(TColToggleSelect))
		value := gval.GoValue()

		model.SetValue(iter, int(TColToggleSelect), glib.NewValue(!value.(bool)))

		ok = model.IterNext(iter)
	}
}
