package dialer

import (
	"context"
	"log/slog"
	"slices"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
)

// AddCollection adds a collection.
func AddCollection(req *m3uetcpb.AddCollectionRequest) (
	res *m3uetcpb.AddCollectionResponse, err error) {

	slog.Info("Adding collection")

	cc, err := getClientConn1()
	if err != nil {
		slog.Error("Failed to get client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewCollectionSvcClient(cc)
	res, err = cl.AddCollection(context.Background(), req)
	return
}

// ApplyCollectionChanges applies collection changes.
func ApplyCollectionChanges(o ...store.CollectionOptions) {
	logw := slog.With("collection_options", o)
	logw.Info("Applying collection changes")

	requests, err := store.CData.FindUpdateCollectionRequests()
	if err != nil {
		logw.Error("Failed to get updated collection requests", "error", err)
	}

	cc, err := getClientConn1()
	if err != nil {
		logw.Error("Failed to get client connection", "error", err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewCollectionSvcClient(cc)

	for _, req := range requests {
		_, err := cl.UpdateCollection(context.Background(), req)
		onerror.NewRecorder(logw).Log(err)
	}

	applyCollectionActionsChanges(o...)
}

func applyCollectionActionsChanges(o ...store.CollectionOptions) {
	logw := slog.With("collection_options", o)
	logw.Info("Applying collection actions changes")

	toScan, toRemove := store.CData.CollectionActionsChanges()

	cc, err := getClientConn1()
	if err != nil {
		logw.Error("Failed to get client connection", "error", err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewCollectionSvcClient(cc)

	var opts store.CollectionOptions
	if len(o) > 0 {
		opts = o[0]
	}

	onerrorw := onerror.NewRecorder(logw)
	for _, id := range toRemove {
		req := &m3uetcpb.RemoveCollectionRequest{Id: id}
		_, err := cl.RemoveCollection(context.Background(), req)
		onerrorw.Log(err)
	}

	if opts.Discover {
		_, err := cl.DiscoverCollections(
			context.Background(),
			&m3uetcpb.Empty{},
		)
		onerrorw.Log(err)
	} else {
		for _, id := range toScan {
			if slices.Contains(toRemove, id) {
				continue
			}
			req := &m3uetcpb.ScanCollectionRequest{
				Id:         id,
				UpdateTags: opts.UpdateTags,
			}
			_, err := cl.ScanCollection(context.Background(), req)
			onerrorw.Log(err)
		}
	}
}

func subscribeToCollectionStore() {
	slog.Info("Subscribing to collection store")

	defer wgcollection.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewCollectionSvcClient(cc)
	stream, err := cl.SubscribeToCollectionStore(
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

		store.CData.ProcessSubscriptionResponse(res)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromCollectionStore() {
	slog.Info("Unsubscribing from collection store")

	id := store.CData.SubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		slog.Error("Failed to obtain client connection", "error", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewCollectionSvcClient(cc)
	_, err = cl.UnsubscribeFromCollectionStore(
		context.Background(),
		&m3uetcpb.UnsubscribeFromCollectionStoreRequest{
			SubscriptionId: id,
		},
	)
	onerror.Log(err)
}
