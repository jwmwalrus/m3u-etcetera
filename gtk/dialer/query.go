package dialer

import (
	"context"
	"fmt"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

// AddQuery adds the query defined by the request.
func AddQuery(req *m3uetcpb.AddQueryRequest) (err error) {
	cc, err := getClientConn1()
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
	cc, err := getClientConn1()
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
	cc, err := getClientConn1()
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

	count = store.QYData.UpdateQueryByResults(res)
	return
}

// RemoveQuery removes the query defined by the request.
func RemoveQuery(req *m3uetcpb.RemoveQueryRequest) (err error) {
	cc, err := getClientConn1()
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

// UpdateQuery updates the query defined by the request.
func UpdateQuery(req *m3uetcpb.UpdateQueryRequest) (err error) {
	cc, err := getClientConn1()
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
	stream, err := cl.SubscribeToQueryStore(
		context.Background(),
		&m3uetcpb.Empty{},
	)
	if err != nil {
		log.Errorf("Error subscribing to collection store: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		store.QYData.ProcessSubscriptionResponse(res)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromQueryStore() {
	log.Info("Unsubscribing from query store")

	id := store.QYData.GetSubscriptionID()

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
