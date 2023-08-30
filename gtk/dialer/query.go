package dialer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
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

// QueryBy performs the query defined by the request and displays
// the results.
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

// QueryInPlaylist apply the query defined by the request and add the results
// to the given target.
func QueryInPlaylist(req *m3uetcpb.QueryInPlaylistRequest) (playlistID int64, err error) {
	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	res, err := cl.QueryInPlaylist(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}

	playlistID = res.PlaylistId
	return
}

// QueryInQueue apply the query defined by the request and add the results
// to the given target.
func QueryInQueue(req *m3uetcpb.QueryInQueueRequest) (err error) {
	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	_, err = cl.QueryInQueue(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return
	}
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
	slog.Info("Subscribing to query store")

	defer wgquery.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to get client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQuerySvcClient(cc)
	stream, err := cl.SubscribeToQueryStore(
		context.Background(),
		&m3uetcpb.Empty{},
	)
	if err != nil {
		slog.Error("Failed to subscribe to collection store", "error", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			slog.Info("Subscription closed by server", "error", err)
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
	slog.Info("Unsubscribing from query store")

	id := store.QYData.GetSubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to get client connection", "error", err)
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
