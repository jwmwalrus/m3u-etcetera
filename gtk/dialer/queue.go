package dialer

import (
	"context"
	"log/slog"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"google.golang.org/grpc/status"
)

// ExecuteQueueAction sends an ExecuteQueueAction request.
func ExecuteQueueAction(req *m3uetcpb.ExecuteQueueActionRequest) (err error) {
	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.ExecuteQueueAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		slog.Error(s.Message())
		return
	}
	return
}

func subscribeToQueueStore() {
	slog.Info("Subscribing to queue store")

	defer wgqueue.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	stream, err := cl.SubscribeToQueueStore(
		context.Background(),
		&m3uetcpb.Empty{},
	)
	if err != nil {
		slog.Error("Failed to subscribe to queue store", "error", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			slog.Info("Subscription closed by server", "error", err)
			break
		}

		store.QData.ProcessSubscriptionResponse(res)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromQueueStore() {
	slog.Info("Unsuubscribing from queue store")

	id := store.QData.GetSubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.UnsubscribeFromQueueStore(
		context.Background(),
		&m3uetcpb.UnsubscribeFromQueueStoreRequest{
			SubscriptionId: id,
		},
	)
	onerror.Log(err)
}
