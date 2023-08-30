package dialer

import (
	"context"
	"log/slog"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
)

// ExecutePlaybackAction -.
func ExecutePlaybackAction(req *m3uetcpb.ExecutePlaybackActionRequest) (err error) {
	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	return
}

func subscribeToPlayback() {
	slog.Info("Subscribing to playback")

	defer wgplayback.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	stream, err := cl.SubscribeToPlayback(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		slog.Error("Failed to subscribe to playback", "error", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			slog.Info("Subscription closed by server", "error", err)
			break
		}

		store.PbData.ProcessSubscriptionResponse(res)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromPlayback() {
	slog.Info("Unsubscribing from playback")

	id := store.PbData.GetSubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
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
