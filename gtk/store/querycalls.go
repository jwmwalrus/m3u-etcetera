package store

import (
	"context"
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
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

	defer wgquery.Done()

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
