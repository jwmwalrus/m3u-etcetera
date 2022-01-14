package store

import (
	"context"
	"strconv"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	cTree = &collectionTree{}

	collectionsModel *gtk.ListStore

	// CStore collection store
	CStore struct {
		subscriptionID string
		Mu             sync.Mutex
		Collection     []*m3uetcpb.Collection
		Track          []*m3uetcpb.Track
	}
)

// CollectionsOptions for the collections dialog
type CollectionsOptions struct {
	Discover   bool
	UpdateTags bool
}

// SetDefaults -
func (co *CollectionsOptions) SetDefaults() {
	co = &CollectionsOptions{}
}

// ApplyCollectionChanges applies collection changes
func ApplyCollectionChanges(o ...CollectionsOptions) {
	log.WithField("collectionOptions", o).
		Info("Applying collection changes")

	cc, err := getClientConn()
	if err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewCollectionSvcClient(cc)

	model := collectionsModel

	var toScan []int64
	var opts CollectionsOptions
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
		for _, c := range CStore.Collection {
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

// CreateCollectionsModel creates a collection model
func CreateCollectionsModel() (model *gtk.ListStore, err error) {
	log.Info("Creating collections model")

	collectionsModel, err = gtk.ListStoreNew(CColumns.getTypes()...)
	if err != nil {
		return
	}

	model = collectionsModel
	return
}

// CreateCollectionTreeModel creates a collection model
func CreateCollectionTreeModel(h collectionTreeHierarchy) (model *gtk.TreeStore, err error) {
	log.WithField("hierarchy", h).
		Info("Creating collections model")

	cTree.model, err = gtk.TreeStoreNew(CTreeColumn.getTypes()...)
	if err != nil {
		return
	}

	cTree.hierarchy = h
	model = cTree.model
	return
}

// FilterCollectionsBy filters the collections by the given value
func FilterCollectionsBy(val string) {
	cTree.filterVal = val
	cTree.rebuild()
}

// GetCollectionsModel returns the current collections model
func GetCollectionsModel() *gtk.ListStore {
	return collectionsModel
}

// GetCollectionTreeModel returns the current collections model
func GetCollectionTreeModel() *gtk.TreeStore {
	return cTree.model
}

func subscribeToCollectionStore() {
	log.Info("Subscribing to collection store")

	defer wgcollections.Done()

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

	appendItem := func(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
		switch res.Item.(type) {
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
			CStore.Collection = append(
				CStore.Collection,
				res.GetCollection(),
			)
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
			CStore.Track = append(CStore.Track, res.GetTrack())
		default:
		}
	}

	changeItem := func(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
		switch res.Item.(type) {
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
			c := res.GetCollection()
			for i := range CStore.Collection {
				if CStore.Collection[i].Id == c.Id {
					CStore.Collection[i] = c
					break
				}
			}
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
			t := res.GetTrack()
			for i := range CStore.Track {
				if CStore.Track[i].Id == t.Id {
					CStore.Track[i] = t
					break
				}
			}
		}
	}

	removeItem := func(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
		switch res.Item.(type) {
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
			c := res.GetCollection()
			n := len(CStore.Collection)
			for i := range CStore.Collection {
				if CStore.Collection[i].Id == c.Id {
					CStore.Collection[i] = CStore.Collection[n-1]
					CStore.Collection = CStore.Collection[:n-1]
					break
				}
			}
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
			t := res.GetTrack()
			n := len(CStore.Track)
			for i := range CStore.Track {
				if CStore.Track[i].Id == t.Id {
					CStore.Track[i] = CStore.Track[n-1]
					CStore.Track = CStore.Track[:n-1]
					break
				}
			}
		}
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		CStore.Mu.Lock()

		if CStore.subscriptionID == "" {
			CStore.subscriptionID = res.SubscriptionId
		}

		cTree.lastEvent = res.Event
		switch res.Event {
		case m3uetcpb.CollectionEvent_CE_INITIAL:
			cTree.initialMode = true
			CStore.Collection = []*m3uetcpb.Collection{}
			CStore.Track = []*m3uetcpb.Track{}
		case m3uetcpb.CollectionEvent_CE_INITIAL_ITEM:
			appendItem(res)
		case m3uetcpb.CollectionEvent_CE_INITIAL_DONE:
			cTree.initialMode = false
			// pass
		case m3uetcpb.CollectionEvent_CE_ITEM_ADDED:
			appendItem(res)
		case m3uetcpb.CollectionEvent_CE_ITEM_CHANGED:
			changeItem(res)
		case m3uetcpb.CollectionEvent_CE_ITEM_REMOVED:
			removeItem(res)
		case m3uetcpb.CollectionEvent_CE_SCANNING:
			cTree.scanningMode = true
		case m3uetcpb.CollectionEvent_CE_SCANNING_DONE:
			cTree.scanningMode = false
		}

		CStore.Mu.Unlock()
		if !cTree.initialMode && !cTree.scanningMode {
			glib.IdleAdd(cTree.update)
			glib.IdleAdd(updateCollectionsModel)
		}
		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromCollectionStore() {
	log.Info("Unsubscribing from collection store")

	CStore.Mu.Lock()
	id := CStore.subscriptionID
	CStore.Mu.Unlock()

	var cc *grpc.ClientConn
	var err error
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		log.Error(err)
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
	if err != nil {
		log.Error(err)
		return
	}
}

func updateCollectionsModel() bool {
	log.Info("Updating collection model")

	model := collectionsModel
	if model == nil {
		return false
	}

	if model.GetNColumns() == 0 {
		return false
	}

	_, ok := model.GetIterFirst()
	if ok {
		model.Clear()
	}

	CStore.Mu.Lock()
	if len(CStore.Collection) > 0 {
		var iter *gtk.TreeIter
		for _, c := range CStore.Collection {
			iter = model.Append()
			tracks := strconv.FormatInt(c.Tracks, 10)
			if c.Scanned != 100 {
				tracks = strconv.Itoa(int(c.Scanned)) + "%"
			}
			err := model.Set(
				iter,
				[]int{
					int(CColCollectionID),
					int(CColName),
					int(CColDescription),
					int(CColLocation),
					int(CColHidden),
					int(CColDisabled),
					int(CColRemote),
					int(CColScanned),
					int(CColTracks),
					int(CColTracksView),
					int(CColRescan),
				},
				[]interface{}{
					c.Id,
					c.Name,
					c.Description,
					c.Location,
					c.Hidden,
					c.Disabled,
					c.Remote,
					c.Scanned,
					c.Tracks,
					tracks,
					false,
				},
			)
			onerror.Log(err)
		}
	}
	CStore.Mu.Unlock()
	return false
}
