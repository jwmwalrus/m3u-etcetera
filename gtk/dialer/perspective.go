package dialer

import (
	"context"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

// SetActivePerspective
func SetActivePerspective(req *m3uetcpb.SetActivePerspectiveRequest) (err error) {
	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPerspectiveSvcClient(cc)
	_, err = cl.SetActivePerspective(context.Background(), req)
	return
}

func subscribeToPerspective() {
	log.Info("Subscribing to perspective")

	defer wgperspective.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPerspectiveSvcClient(cc)
	stream, err := cl.SubscribeToPerspective(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Errorf("Error subscribing to perspective: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		store.PerspData.ProcessSubscriptionResponse(res)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromPerspective() {
	log.Info("Unsubscribing from perspective")

	id := store.PerspData.GetSubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPerspectiveSvcClient(cc)
	_, err = cl.UnsubscribeFromPerspective(
		context.Background(),
		&m3uetcpb.UnsubscribeFromPerspectiveRequest{
			SubscriptionId: id,
		},
	)
	onerror.Log(err)
}
