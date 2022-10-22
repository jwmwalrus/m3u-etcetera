package dialer

import (
	"context"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

// ExecuteQueueAction sends an ExecuteQueueAction request
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
		log.Error(s.Message())
		return
	}
	return
}

func subscribeToQueueStore() {
	log.Info("Subscribing to queue store")

	defer wgqueue.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	stream, err := cl.SubscribeToQueueStore(
		context.Background(),
		&m3uetcpb.Empty{},
	)
	if err != nil {
		log.Errorf("Error subscribing to queue store: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
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
	log.Info("Unsuubscribing from queue store")

	id := store.QData.GetSubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
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
