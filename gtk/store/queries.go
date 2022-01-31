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
		req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
			PlaylistId: targetID,
			Action:     m3uetcpb.PlaylistTrackAction_PT_APPEND,
			TrackIds:   ids,
		}

		err = ExecutePlaylistTrackAction(req)
		if err != nil {
			s := status.Convert(err)
			err = fmt.Errorf(s.Message())
			return
		}
	} else {
		req := &m3uetcpb.ExecutePlaybackActionRequest{
			Action: m3uetcpb.PlaybackAction_PB_PLAY,
			Ids:    ids,
		}

		err = ExecutePlaybackAction(req)
		if err != nil {
			s := status.Convert(err)
			err = fmt.Errorf(s.Message())
			return
		}
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

		if len(qy.CollectionIds) > 0 {
			// TODO: bound by collections
			continue
		}
		all.unbounded = append(all.unbounded, queryType{qy.Id, qy.Name, getKeywords(qy)})
		all.ids = append(all.ids, qy.Id)
	}
	QYData.Mu.Unlock()

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

	QYData.Mu.Lock()
	var iter *gtk.TreeIter
	for i, t := range QYData.tracks {
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
	QYData.Mu.Unlock()
	return false
}
