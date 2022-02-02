package store

import (
	"context"
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
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
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

// AddQuery adds the query defined by the request
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
		err = fmt.Errorf(s.Message())
		return
	}
	return
}

// ApplyQuery apply the query defined by the request and add the results
// to the given target
func ApplyQuery(req *m3uetcpb.ApplyQueryRequest, targetID int64) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	res, err := cl.ApplyQuery(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	var ids []int64
	for i := range res.Tracks {
		ids = append(ids, res.Tracks[i].Id)
	}

	if targetID > 0 {
		reqpl := &m3uetcpb.ExecutePlaylistTrackActionRequest{
			PlaylistId: targetID,
			Action:     m3uetcpb.PlaylistTrackAction_PT_APPEND,
			TrackIds:   ids,
		}

		err = ExecutePlaylistTrackAction(reqpl)
		if err != nil {
			s := status.Convert(err)
			err = fmt.Errorf(s.Message())
		}
		return
	}

	reqpb := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_PLAY,
		Ids:    ids,
	}

	err = ExecutePlaybackAction(reqpb)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
	}
	return
}

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
	log.Info("Creating queries model")

	queryTree.model, err = gtk.TreeStoreNew(QYTreeColumn.getTypes()...)
	if err != nil {
		return
	}

	model = queryTree.model
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

// GetQueryTreeModel returns the current collections model
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

// QueryBy performs the query defined by the request and displays
// the results
func QueryBy(req *m3uetcpb.QueryByRequest) (count int, err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	res, err := cl.QueryBy(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	QYData.Mu.Lock()
	QYData.tracks = res.Tracks
	count = len(res.Tracks)
	QYData.Mu.Unlock()

	glib.IdleAdd(updateQueryResults)
	return
}

// RemoveQuery removes the query defined by the request
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
		err = fmt.Errorf(s.Message())
		return
	}
	return
}

// UpdateQuery updates the query defined by the request
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
		err = fmt.Errorf(s.Message())
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
		for _, qy := range QYData.Query {
			if qy.Id == res.Query.Id {
				return
			}
		}
		QYData.Query = append(
			QYData.Query,
			res.Query,
		)
	}

	changeItem := func(res *m3uetcpb.SubscribeToQueryStoreResponse) {
		qy := res.Query
		for i := range QYData.Query {
			if QYData.Query[i].Id == qy.Id {
				QYData.Query[i] = qy
				break
			}
		}
	}

	removeItem := func(res *m3uetcpb.SubscribeToQueryStoreResponse) {
		n := len(QYData.Query)
		for i := range QYData.Query {
			if QYData.Query[i].Id == res.Query.Id {
				QYData.Query[i] = QYData.Query[n-1]
				QYData.Query = QYData.Query[:n-1]
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

		QYData.Mu.Lock()

		if QYData.subscriptionID == "" {
			QYData.subscriptionID = res.SubscriptionId
		}

		switch res.Event {
		case m3uetcpb.QueryEvent_QYE_INITIAL:
			queryTree.initialMode = true
			QYData.Query = []*m3uetcpb.Query{}
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
		QYData.Mu.Unlock()

		if !queryTree.initialMode {
			glib.IdleAdd(queryTree.update)
		}
		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromQueryStore() {
	log.Info("Unsubscribing from query store")

	QYData.Mu.Lock()
	id := QYData.subscriptionID
	QYData.Mu.Unlock()

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
	onerror.Log(err)
}

func (qyt *queryTreeModel) update() bool {
	log.Info("Updating queries model")

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
