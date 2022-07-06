package store

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

// CollectionOptions for the collections dialog
type CollectionOptions struct {
	Discover   bool
	UpdateTags bool
}

// SetDefaults -
func (co *CollectionOptions) SetDefaults() {}

type collectionData struct {
	subscriptionID string
	collection     []*m3uetcpb.Collection
	track          []*m3uetcpb.Track

	mu sync.Mutex
}

var (
	// CData collection store
	CData = &collectionData{}

	cTree = &collectionTree{}

	collectionModel *gtk.ListStore
	cProgress       *gtk.ProgressBar

	collectionNameMap          map[int64]string
	collectionTreeHierarchyMap map[string]collectionTreeHierarchy
)

// CollectionAlreadyExists returns true if the location and the name are not
// already in use by another collection
func (cd *collectionData) CollectionAlreadyExists(location, name string) bool {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	for _, c := range cd.collection {
		if c.Location == location ||
			strings.EqualFold(c.Name, name) {
			return true
		}
	}
	return false
}

func (cd *collectionData) GetCollectionActionsChanges() (toScan, toRemove []int64) {
	model := collectionModel

	cd.mu.Lock()
	defer cd.mu.Unlock()

	iter, ok := model.GetIterFirst()
	for ok {
		row, err := GetListStoreModelValues(
			model,
			iter,
			[]ModelColumn{
				CColCollectionID,
				CColName,
				CColActionRescan,
				CColActionRemove,
			},
		)
		if err != nil {
			log.Error(err)
			return
		}
		id := row[CColCollectionID].(int64)
		for _, c := range cd.collection {
			if id != c.Id {
				continue
			}

			rescan := row[CColActionRescan].(bool)
			if rescan {
				toScan = append(toScan, id)
			}

			remove := row[CColActionRemove].(bool)
			if remove {
				toRemove = append(toRemove, id)
			}
		}
		ok = model.IterNext(iter)
	}

	return
}

func (cd *collectionData) GetSubscriptionID() string {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	return cd.subscriptionID
}

func (cd *collectionData) GetTracksTotalCount() int64 {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	var total int64
	for _, c := range cd.collection {
		total += c.Tracks
	}
	return total
}

func (cd *collectionData) GetUpdateCollectionRequests() ([]*m3uetcpb.UpdateCollectionRequest, error) {
	requests := []*m3uetcpb.UpdateCollectionRequest{}

	model := GetCollectionModel()

	cd.mu.Lock()
	defer cd.mu.Unlock()

	iter, ok := model.GetIterFirst()
	for ok {
		row, err := GetListStoreModelValues(
			model,
			iter,
			[]ModelColumn{
				CColCollectionID,
				CColName,
				CColDescription,
				CColRemoteLocation,
				CColDisabled,
				CColRemote,
			},
		)
		if err != nil {
			return nil, err
		}
		id := row[CColCollectionID].(int64)
		for _, c := range cd.collection {
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

			requests = append(requests, req)
			break
		}
		ok = model.IterNext(iter)
	}

	return requests, nil
}

func (cd *collectionData) ProcessSubscriptionResponse(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	if cd.subscriptionID == "" {
		cd.subscriptionID = res.SubscriptionId
	}

	cTree.lastEvent = res.Event
	switch res.Event {
	case m3uetcpb.CollectionEvent_CE_INITIAL:
		cTree.initialMode = true
		cd.collection = []*m3uetcpb.Collection{}
		cd.track = []*m3uetcpb.Track{}
	case m3uetcpb.CollectionEvent_CE_INITIAL_ITEM:
		cd.appendCDataItem(res)
	case m3uetcpb.CollectionEvent_CE_INITIAL_DONE:
		cTree.initialMode = false
	case m3uetcpb.CollectionEvent_CE_ITEM_ADDED:
		cd.appendCDataItem(res)
	case m3uetcpb.CollectionEvent_CE_ITEM_CHANGED:
		cd.changeCDataItem(res)
	case m3uetcpb.CollectionEvent_CE_ITEM_REMOVED:
		cd.removeCDataItem(res)
	case m3uetcpb.CollectionEvent_CE_SCANNING:
		cTree.scanningMode = true
	case m3uetcpb.CollectionEvent_CE_SCANNING_DONE:
		cTree.scanningMode = false
	}

	glib.IdleAdd(cd.updateCollectionModel)
	if !cTree.initialMode && !cTree.scanningMode {
		glib.IdleAdd(cTree.update)
	} else {
		glib.IdleAdd(cd.updateScanningProgress)
	}
}

func (cd *collectionData) SwitchHierarchyTo(id string, grouped bool) {
	if cTree.initialMode || cTree.scanningMode {
		return
	}

	cTree.hierarchy = collectionTreeHierarchyMap[id]
	cTree.groupByCollection = grouped

	glib.IdleAdd(cTree.update)
}

func (cd *collectionData) appendCDataItem(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
		item := res.GetCollection()
		for _, c := range cd.collection {
			if c.Id == item.Id {
				return
			}
		}
		cd.collection = append(
			cd.collection,
			item,
		)
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
		item := res.GetTrack()
		for _, t := range cd.track {
			if t.Id == item.Id {
				return
			}
		}
		cd.track = append(cd.track, item)
	default:
	}
}

func (cd *collectionData) changeCDataItem(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
		c := res.GetCollection()
		for i := range cd.collection {
			if cd.collection[i].Id == c.Id {
				cd.collection[i] = c
				break
			}
		}
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
		t := res.GetTrack()
		for i := range cd.track {
			if cd.track[i].Id == t.Id {
				cd.track[i] = t
				break
			}
		}
	}
}

func (cd *collectionData) removeCDataItem(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
		c := res.GetCollection()
		n := len(cd.collection)
		for i := range cd.collection {
			if cd.collection[i].Id == c.Id {
				cd.collection[i] = cd.collection[n-1]
				cd.collection = cd.collection[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
		t := res.GetTrack()
		n := len(cd.track)
		for i := range cd.track {
			if cd.track[i].Id == t.Id {
				cd.track[i] = cd.track[n-1]
				cd.track = cd.track[:n-1]
				break
			}
		}
	}
}

func (cd *collectionData) updateCollectionModel() bool {
	if cTree.initialMode {
		return false
	}

	log.Info("Updating collection model")

	cd.mu.Lock()
	defer cd.mu.Unlock()

	model := collectionModel
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

	if len(cd.collection) > 0 {
		var iter *gtk.TreeIter
		for _, c := range cd.collection {
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
					int(CColActionRescan),
					int(CColActionRemove),
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
					false,
				},
			)
			onerror.Log(err)
		}
	}

	return false
}

func (cd *collectionData) updateCollectionNamesMap() {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	collectionNameMap = make(map[int64]string)
	for _, c := range cd.collection {
		collectionNameMap[c.Id] = c.Name
	}
}

func (cd *collectionData) updateScanningProgress() bool {
	if cProgress == nil {
		log.Error("Collections progress bar is unavailable")
		return false
	}

	cd.mu.Lock()
	defer cd.mu.Unlock()

	scanned := 0
	for _, c := range cd.collection {
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

// CreateCollectionModel creates a collection model
func CreateCollectionModel() (model *gtk.ListStore, err error) {
	log.Info("Creating collection model")

	collectionModel, err = gtk.ListStoreNew(CColumns.getTypes()...)
	if err != nil {
		return
	}

	model = collectionModel
	return
}

// CreateCollectionTreeModel creates a collection model
func CreateCollectionTreeModel(h collectionTreeHierarchy) (
	model *gtk.TreeStore, err error) {

	log.WithField("hierarchy", h).
		Info("Creating collection tree model")

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

// FilterCollectionTreeBy filters the collection tree by the given value
func FilterCollectionTreeBy(val string) {
	cTree.filterVal = val
	cTree.rebuild()
}

// GetCollectionModel returns the current collection model
func GetCollectionModel() *gtk.ListStore {
	return collectionModel
}

// GetCollectionTreeModel returns the current collection tree model
func GetCollectionTreeModel() *gtk.TreeStore {
	return cTree.model
}

func init() {
	collectionNameMap = make(map[int64]string)

	hl := []collectionTreeHierarchy{
		ArtistYearAlbumTree,
		ArtistAlbumTree,
		AlbumTree,
		GenreArtistAlbumTree,
		YearArtistAlbumTree,
	}
	collectionTreeHierarchyMap = make(map[string]collectionTreeHierarchy)
	for _, h := range hl {
		collectionTreeHierarchyMap[h.String()] = h
	}
}
