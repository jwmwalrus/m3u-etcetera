package store

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/pkg/qparams"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

var (
	queryTreeModel    *gtk.TreeStore
	queryResultsModel *gtk.ListStore
	queriesFilterVal  string

	// QYStore query store
	QYStore struct {
		subscriptionID string
		Mu             sync.Mutex
		Query          []*m3uetcpb.Query
		tracks         []*m3uetcpb.Track
	}
)

func AddQuery(req *m3uetcpb.AddQueryRequest) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	_, err = cl.AddQuery(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = errors.New(s.Message())
		return
	}
	return
}

func ApplyQuery(req *m3uetcpb.ApplyQueryRequest) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	res, err := cl.ApplyQuery(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = errors.New(s.Message())
		return
	}

	var ids []int64
	for i := range res.Tracks {
		ids = append(ids, res.Tracks[i].Id)
	}

	req2 := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_PLAY,
		Ids:    ids,
	}

	err = ExecutePlaybackAction(req2)
	if err != nil {
		s := status.Convert(err)
		err = errors.New(s.Message())
		return
	}
	return
}

func ClearQueryResults() {
	model := queryResultsModel

	_, ok := model.GetIterFirst()
	if ok {
		model.Clear()
	}
}

// CreateQueryTreeModel creates a query model
func CreateQueryTreeModel() (model *gtk.TreeStore, err error) {
	log.Info("Creating queries model")

	queryTreeModel, err = gtk.TreeStoreNew(QYTreeColumn.getTypes()...)
	if err != nil {
		return
	}

	model = queryTreeModel
	return
}

// CreateQueryResultsModel creates a query model
func CreateQueryResultsModel() (model *gtk.ListStore, err error) {
	log.Info("Creating queries model")

	queryResultsModel, err = gtk.ListStoreNew(TColumns.getTypes()...)
	if err != nil {
		return
	}

	model = queryResultsModel
	return
}

// FilterQueriesBy filters the collections by the given value
func FilterQueriesBy(val string) {
	queriesFilterVal = val
	updateQueryTree()
}

func GetQuery(id int64) (qy *m3uetcpb.Query) {
	QYStore.Mu.Lock()
	for _, v := range QYStore.Query {
		if v.Id == id {
			qy = v
			break
		}
	}
	QYStore.Mu.Unlock()
	return
}

// GetQueryTreeModel returns the current collections model
func GetQueryTreeModel() *gtk.TreeStore {
	return queryTreeModel
}

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

func QueryBy(req *m3uetcpb.QueryByRequest) (err error, count int) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	res, err := cl.QueryBy(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = errors.New(s.Message())
		return
	}

	QYStore.Mu.Lock()
	QYStore.tracks = res.Tracks
	count = len(res.Tracks)
	QYStore.Mu.Unlock()
	glib.IdleAdd(updateQueryResults)
	return
}

func RemoveQuery(req *m3uetcpb.RemoveQueryRequest) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	_, err = cl.RemoveQuery(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = errors.New(s.Message())
		return
	}
	return
}

func UpdateQuery(req *m3uetcpb.UpdateQueryRequest) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	_, err = cl.UpdateQuery(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = errors.New(s.Message())
		return
	}
	return
}

func subscribeToQueryStore() {
	log.Info("Subscribing to query store")

	defer wgqueries.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	stream, err := cl.SubscribeToQueryStore(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Errorf("Error subscribing to collection store: %v", err)
		return
	}

	appendItem := func(res *m3uetcpb.SubscribeToQueryStoreResponse) {
		QYStore.Query = append(
			QYStore.Query,
			res.Query,
		)
	}

	changeItem := func(res *m3uetcpb.SubscribeToQueryStoreResponse) {
		qy := res.Query
		for i := range QYStore.Query {
			if QYStore.Query[i].Id == qy.Id {
				QYStore.Query[i] = qy
				break
			}
		}
	}

	removeItem := func(res *m3uetcpb.SubscribeToQueryStoreResponse) {
		n := len(QYStore.Query)
		for i := range QYStore.Query {
			if QYStore.Query[i].Id == res.Query.Id {
				QYStore.Query[i] = QYStore.Query[n-1]
				QYStore.Query = QYStore.Query[:n-1]
				break
			}
		}
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		QYStore.Mu.Lock()

		if QYStore.subscriptionID == "" {
			QYStore.subscriptionID = res.SubscriptionId
		}

		switch res.Event {
		case m3uetcpb.QueryEvent_QYE_INITIAL:
			QYStore.Query = []*m3uetcpb.Query{}
		case m3uetcpb.QueryEvent_QYE_INITIAL_ITEM:
			appendItem(res)
		case m3uetcpb.QueryEvent_QYE_INITIAL_DONE:
			// pass
		case m3uetcpb.QueryEvent_QYE_ITEM_ADDED:
			appendItem(res)
		case m3uetcpb.QueryEvent_QYE_ITEM_CHANGED:
			changeItem(res)
		case m3uetcpb.QueryEvent_QYE_ITEM_REMOVED:
			removeItem(res)
		}
		QYStore.Mu.Unlock()

		if res.Event != m3uetcpb.QueryEvent_QYE_INITIAL &&
			res.Event != m3uetcpb.QueryEvent_QYE_INITIAL_ITEM {
			glib.IdleAdd(updateQueryTree)
		}
		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromQueryStore() {
	log.Info("Unsubscribing from query store")

	QYStore.Mu.Lock()
	id := QYStore.subscriptionID
	QYStore.Mu.Unlock()

	cc, err := getClientConn()
	if err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	_, err = cl.UnsubscribeFromQueryStore(
		context.Background(),
		&m3uetcpb.UnsubscribeFromQueryStoreRequest{
			SubscriptionId: id,
		},
	)
	if err != nil {
		log.Error(err)
		return
	}
}

func updateQueryTree() bool {
	log.Info("Updating queries model")

	model := queryTreeModel
	if model == nil {
		return false
	}

	if model.GetNColumns() == 0 {
		return false
	}

	type queryType struct {
		id   int64
		name string
		kw   string
	}

	type kind int

	const (
		kindCollection kind = iota
	)

	type boundaryType struct {
		name       string
		boundaryID int64
		ids        []int64
		query      []queryType
	}

	type boundaryKind struct {
		ids      []int64
		bMap     map[int64]int
		boundary []boundaryType
	}

	type byType struct {
		ids       []int64
		kind      []boundaryKind
		unbounded []queryType
	}

	_, ok := model.GetIterFirst()
	if ok {
		model.Clear()
	}

	all := byType{kind: []boundaryKind{{}}}

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

	QYStore.Mu.Lock()
	for _, qy := range QYStore.Query {
		if queriesFilterVal != "" {
			kw := getKeywords(qy)
			match := false
			for _, s := range strings.Split(queriesFilterVal, " ") {
				match = match || strings.Contains(kw, s)
				if match {
					break
				}
			}
			if !match {
				continue
			}
		}

		if len(qy.CollectionIds) > 0 {
			// TODO: bound by collections
			continue
		}
		all.unbounded = append(all.unbounded, queryType{qy.Id, qy.Name, getKeywords(qy)})
		all.ids = append(all.ids, qy.Id)
	}
	QYStore.Mu.Unlock()

	sort.SliceStable(all.unbounded, func(i, j int) bool {
		return all.unbounded[i].name < all.unbounded[j].name
	})

	for _, ub := range all.unbounded {
		iter := model.Append(nil)
		model.SetValue(iter, int(QYColTree), ub.name)
		model.SetValue(iter, int(QYColTreeIDList), strconv.FormatInt(ub.id, 10))
		model.SetValue(iter, int(QYColTreeKeywords), ub.kw)
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

	QYStore.Mu.Lock()
	var iter *gtk.TreeIter
	for i, t := range QYStore.tracks {
		iter = model.Append()
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
				fmt.Sprint(time.Duration(t.Duration) * time.Nanosecond),
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
	QYStore.Mu.Unlock()
	return false
}
