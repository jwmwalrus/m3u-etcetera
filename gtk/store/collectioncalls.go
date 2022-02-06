package store

import (
	"context"

	"github.com/gotk3/gotk3/glib"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

// AddCollection adds a collection
func AddCollection(req *m3uetcpb.AddCollectionRequest) (res *m3uetcpb.AddCollectionResponse, err error) {
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
func ApplyCollectionChanges(o ...CollectionOptions) {
	log.WithField("collectionOptions", o).
		Info("Applying collection changes")

	cc, err := getClientConn1()
	if err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewCollectionSvcClient(cc)

	model := collectionModel

	var toScan []int64
	var opts CollectionOptions
	if len(o) > 0 {
		opts = o[0]
	}

	iter, ok := model.GetIterFirst()
	for ok {
		row, err := GetListStoreModelValues(model, iter, []ModelColumn{CColCollectionID, CColName, CColDescription, CColRemoteLocation, CColDisabled, CColRemote, CColRescan})
		if err != nil {
			log.Error(err)
			return
		}
		id := row[CColCollectionID].(int64)
		for _, c := range CData.Collection {
			if id != c.Id {
				continue
			}

			req := &m3uetcpb.UpdateCollectionRequest{Id: c.Id}

			newName := row[CColName].(string)
			if newName != c.Name && newName != "" {
				req.NewName = newName
			}

			newDescription := row[CColDescription].(string)
			if newDescription != c.Description {
				if newDescription != "" {
					req.NewDescription = newDescription
				} else {
					req.ResetDescription = true
				}
			}

			remoteLocation := row[CColRemoteLocation].(string)
			if remoteLocation != c.RemoteLocation {
				if remoteLocation != "" {
					req.NewRemoteLocation = remoteLocation
				} else {
					req.ResetRemoteLocation = true
				}
			}

			disabled := row[CColDisabled].(bool)
			if disabled {
				req.Disable = true
			} else {
				req.Enable = true
			}

			remote := row[CColRemote].(bool)
			if remote {
				req.MakeRemote = true
			} else {
				req.MakeLocal = true
			}

			rescan := row[CColRescan].(bool)
			if rescan {
				toScan = append(toScan, id)
			}

			_, err := cl.UpdateCollection(context.Background(), req)
			onerror.Log(err)
			break
		}
		ok = model.IterNext(iter)
	}

	if opts.Discover {
		_, err := cl.DiscoverCollections(context.Background(), &m3uetcpb.Empty{})
		onerror.Log(err)
	} else {
		for _, id := range toScan {
			req := &m3uetcpb.ScanCollectionRequest{Id: id, UpdateTags: opts.UpdateTags}
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
	stream, err := cl.SubscribeToCollectionStore(context.Background(), &m3uetcpb.Empty{})
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

		CData.Mu.Lock()

		if CData.subscriptionID == "" {
			CData.subscriptionID = res.SubscriptionId
		}

		cTree.lastEvent = res.Event
		switch res.Event {
		case m3uetcpb.CollectionEvent_CE_INITIAL:
			cTree.initialMode = true
			CData.Collection = []*m3uetcpb.Collection{}
			CData.Track = []*m3uetcpb.Track{}
		case m3uetcpb.CollectionEvent_CE_INITIAL_ITEM:
			appendCDataItem(res)
		case m3uetcpb.CollectionEvent_CE_INITIAL_DONE:
			cTree.initialMode = false
		case m3uetcpb.CollectionEvent_CE_ITEM_ADDED:
			appendCDataItem(res)
		case m3uetcpb.CollectionEvent_CE_ITEM_CHANGED:
			changeCDataItem(res)
		case m3uetcpb.CollectionEvent_CE_ITEM_REMOVED:
			removeCDataItem(res)
		case m3uetcpb.CollectionEvent_CE_SCANNING:
			cTree.scanningMode = true
		case m3uetcpb.CollectionEvent_CE_SCANNING_DONE:
			cTree.scanningMode = false
		}

		CData.Mu.Unlock()

		glib.IdleAdd(updateCollectionModel)
		if !cTree.initialMode && !cTree.scanningMode {
			glib.IdleAdd(cTree.update)
		} else {
			glib.IdleAdd(updateScanningProgress)
		}

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func appendCDataItem(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
		item := res.GetCollection()
		for _, c := range CData.Collection {
			if c.Id == item.Id {
				return
			}
		}
		CData.Collection = append(
			CData.Collection,
			item,
		)
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
		item := res.GetTrack()
		for _, t := range CData.Track {
			if t.Id == item.Id {
				return
			}
		}
		CData.Track = append(CData.Track, item)
	default:
	}
}

func changeCDataItem(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
		c := res.GetCollection()
		for i := range CData.Collection {
			if CData.Collection[i].Id == c.Id {
				CData.Collection[i] = c
				break
			}
		}
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
		t := res.GetTrack()
		for i := range CData.Track {
			if CData.Track[i].Id == t.Id {
				CData.Track[i] = t
				break
			}
		}
	}
}

func removeCDataItem(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
		c := res.GetCollection()
		n := len(CData.Collection)
		for i := range CData.Collection {
			if CData.Collection[i].Id == c.Id {
				CData.Collection[i] = CData.Collection[n-1]
				CData.Collection = CData.Collection[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
		t := res.GetTrack()
		n := len(CData.Track)
		for i := range CData.Track {
			if CData.Track[i].Id == t.Id {
				CData.Track[i] = CData.Track[n-1]
				CData.Track = CData.Track[:n-1]
				break
			}
		}
	}
}

func unsubscribeFromCollectionStore() {
	log.Info("Unsubscribing from collection store")

	CData.Mu.Lock()
	id := CData.subscriptionID
	CData.Mu.Unlock()

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
