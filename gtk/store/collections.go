package store

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	cTree = &collectionTree{}

	collectionsModel *gtk.ListStore
	cProgress        *gtk.ProgressBar

	// CData collection store
	CData struct {
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
func (co *CollectionsOptions) SetDefaults() {}

// AddCollection adds a collection
func AddCollection(req *m3uetcpb.AddCollectionRequest) (res *m3uetcpb.AddCollectionResponse, err error) {
	log.Info("Adding collection")

	cc, err := getClientConn()
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

// CollectionAlreadyExists returns true if the location and the name are not
// already in use by another collection
func CollectionAlreadyExists(location, name string) bool {
	CData.Mu.Lock()
	defer CData.Mu.Unlock()

	for _, c := range CData.Collection {
		if c.Location == location ||
			c.Name == name {
			return true
		}
	}
	return false
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

	var errp error
	cProgress, errp = builder.GetProgressBar("collections_scanning_progress")
	if errp != nil {
		cProgress = nil
		log.Error(errp)
	}
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

		glib.IdleAdd(updateCollectionsModel)
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
	onerror.Log(err)
}

func updateCollectionsModel() bool {
	if cTree.initialMode {
		return false
	}

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

	CData.Mu.Lock()
	if len(CData.Collection) > 0 {
		var iter *gtk.TreeIter
		for _, c := range CData.Collection {
			iter = model.Append()
			tracks := strconv.FormatInt(c.Tracks, 10)
			if c.Scanned != 100 {
				tracks = strconv.Itoa(int(c.Scanned)) + "%"
			}
			persp := m3uetcpb.Perspective_name[int32(c.Perspective)]
			err := model.Set(
				iter,
				[]int{
					int(CColCollectionID),
					int(CColName),
					int(CColDescription),
					int(CColLocation),
					int(CColPerspective),
					int(CColDisabled),
					int(CColRemote),
					int(CColTracks),
					int(CColTracksView),
					int(CColRescan),
				},
				[]interface{}{
					c.Id,
					c.Name,
					c.Description,
					c.Location,
					persp,
					c.Disabled,
					c.Remote,
					c.Tracks,
					tracks,
					false,
				},
			)
			onerror.Log(err)
		}
	}
	CData.Mu.Unlock()
	return false
}

func updateScanningProgress() bool {
	if cProgress == nil {
		log.Error("Collections progress bar is unavailable")
		return false
	}

	CData.Mu.Lock()
	defer CData.Mu.Unlock()

	scanned := 0
	for _, c := range CData.Collection {
		if c.Scanned != 100 && int(c.Scanned) > scanned {
			scanned = int(c.Scanned)
		}
	}
	if scanned > 0 {
		cProgress.SetVisible(true)
		cProgress.SetFraction(float64(scanned) / float64(100))
		cProgress.SetText(fmt.Sprintf("Scanning: %v%%", scanned))
		return false
	}
	cProgress.SetVisible(false)
	return false
}
