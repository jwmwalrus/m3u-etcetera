package dialer

import (
	"context"
	"log/slog"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
)

// SetActivePerspective sets a new active perspective.
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
	slog.Info("Subscribing to perspective")

	defer wgperspective.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPerspectiveSvcClient(cc)
	stream, err := cl.SubscribeToPerspective(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		slog.Error("Failed to subscribe to perspective", "error", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			slog.Info("Subscription closed by server", "error", err)
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
	slog.Info("Unsubscribing from perspective")

	id := store.PerspData.GetSubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
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
