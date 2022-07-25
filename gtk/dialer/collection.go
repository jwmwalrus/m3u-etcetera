package dialer

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

// AddCollection adds a collection
func AddCollection(req *m3uetcpb.AddCollectionRequest) (
	res *m3uetcpb.AddCollectionResponse, err error) {

	log.Info("Adding collection")

	cc, err := getClientConn1()
	if err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewCollectionSvcClient(cc)
	res, err = cl.AddCollection(context.Background(), req)
	return
}

// ApplyCollectionChanges applies collection changes
func ApplyCollectionChanges(o ...store.CollectionOptions) {
	log.WithField("collectionOptions", o).
		Info("Applying collection changes")

	requests, err := store.CData.GetUpdateCollectionRequests()
	if err != nil {
		log.Error(err)
	}

	cc, err := getClientConn1()
	if err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewCollectionSvcClient(cc)

	for _, req := range requests {
		_, err := cl.UpdateCollection(context.Background(), req)
		onerror.Log(err)
	}

	applyCollectionActionsChanges(o...)
}

func applyCollectionActionsChanges(o ...store.CollectionOptions) {
	log.WithField("collectionOptions", o).
		Info("Applying collection actions changes")

	toScan, toRemove := store.CData.GetCollectionActionsChanges()

	cc, err := getClientConn1()
	if err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewCollectionSvcClient(cc)

	var opts store.CollectionOptions
	if len(o) > 0 {
		opts = o[0]
	}

	for _, id := range toRemove {
		req := &m3uetcpb.RemoveCollectionRequest{Id: id}
		_, err := cl.RemoveCollection(context.Background(), req)
		onerror.Log(err)
	}

	if opts.Discover {
		_, err := cl.DiscoverCollections(
			context.Background(),
			&empty.Empty{},
		)
		onerror.Log(err)
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
			onerror.Log(err)
		}
	}
}

func subscribeToCollectionStore() {
	log.Info("Subscribing to collection store")

	defer wgcollection.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewCollectionSvcClient(cc)
	stream, err := cl.SubscribeToCollectionStore(
		context.Background(),
		&empty.Empty{},
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

		store.CData.ProcessSubscriptionResponse(res)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromCollectionStore() {
	log.Info("Unsubscribing from collection store")

	id := store.CData.GetSubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
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
