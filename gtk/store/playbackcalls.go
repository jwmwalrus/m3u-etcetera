package store

import (
	"context"

	"github.com/gotk3/gotk3/glib"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

// ExecutePlaybackAction -
func ExecutePlaybackAction(req *m3uetcpb.ExecutePlaybackActionRequest) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	return
}

func subscribeToPlayback() {
	log.Info("Subscribing to playback")

	defer wgplayback.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	stream, err := cl.SubscribeToPlayback(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Errorf("Error subscribing to playback: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		pbdata.mu.Lock()
		pbdata.res = res
		pbdata.mu.Unlock()

		glib.IdleAdd(pbdata.updatePlayback)
		glib.IdleAdd(pbdata.setCover)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromPlayback() {
	log.Info("Unsubscribing from playback")

	pbdata.mu.Lock()
	id := pbdata.res.SubscriptionId
	pbdata.mu.Unlock()

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.UnsubscribeFromPlayback(
		context.Background(),
		&m3uetcpb.UnsubscribeFromPlaybackRequest{
			SubscriptionId: id,
		},
	)
	onerror.Log(err)
}
