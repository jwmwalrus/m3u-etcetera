package store

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
)

// CollectionOptions for the collections dialog.
type CollectionOptions struct {
	Discover   bool
	UpdateTags bool
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

// SetDefaults -.
func (co *CollectionOptions) SetDefaults() {}

type collectionData struct {
	subscriptionID string
	collection     []*m3uetcpb.Collection
	track          []*m3uetcpb.Track

	mu sync.RWMutex
}

var (
	// CData collection store.
	CData = &collectionData{}

	cTree = &collectionTree{}

	collectionModel *gtk.ListStore
	cProgress       *gtk.ProgressBar

	collectionNameMap          map[int64]string
	collectionTreeHierarchyMap map[string]collectionTreeHierarchy
)

// CollectionAlreadyExists returns true if the location and the name are not
// already in use by another collection.
func (cd *collectionData) CollectionAlreadyExists(location, name string) bool {
	cd.mu.RLock()
	defer cd.mu.RUnlock()

	for _, c := range cd.collection {
		if c.Location == location ||
			strings.EqualFold(c.Name, name) {
			return true
		}
	}
	return false
}

func (cd *collectionData) CollectionActionsChanges() (toScan, toRemove []int64) {
	model := collectionModel

	cd.mu.RLock()
	defer cd.mu.RUnlock()

	iter, ok := model.IterFirst()
	for ok {
		row, err := GetTreeModelValues(
			&model.TreeModel,
			iter,
			[]ModelColumn{
				CColCollectionID,
				CColName,
				CColActionRescan,
				CColActionRemove,
			},
		)
		if err != nil {
			slog.Error("Failed to get tree-model values", "error", err)
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

func (cd *collectionData) SubscriptionID() string {
	cd.mu.RLock()
	defer cd.mu.RUnlock()

	return cd.subscriptionID
}

func (cd *collectionData) TracksTotalCount() int64 {
	cd.mu.RLock()
	defer cd.mu.RUnlock()

	var total int64
	for _, c := range cd.collection {
		total += c.Tracks
	}
	return total
}

func (cd *collectionData) FindUpdateCollectionRequests() ([]*m3uetcpb.UpdateCollectionRequest, error) {
	requests := []*m3uetcpb.UpdateCollectionRequest{}

	model := GetCollectionModel()

	cd.mu.RLock()
	defer cd.mu.RUnlock()

	iter, ok := model.IterFirst()
	for ok {
		row, err := GetTreeModelValues(
			&model.TreeModel,
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

	cTree.setLastEvent(res.Event)
	switch res.Event {
	case m3uetcpb.CollectionEvent_CE_INITIAL:
		cTree.setInitialMode(true)
		cd.collection = []*m3uetcpb.Collection{}
		cd.track = []*m3uetcpb.Track{}
	case m3uetcpb.CollectionEvent_CE_INITIAL_ITEM:
		cd.appendCDataItem(res)
	case m3uetcpb.CollectionEvent_CE_INITIAL_DONE:
		cTree.setInitialMode(false)
	case m3uetcpb.CollectionEvent_CE_ITEM_ADDED:
		cd.appendCDataItem(res)
	case m3uetcpb.CollectionEvent_CE_ITEM_CHANGED:
		cd.changeCDataItem(res)
	case m3uetcpb.CollectionEvent_CE_ITEM_REMOVED:
		cd.removeCDataItem(res)
	case m3uetcpb.CollectionEvent_CE_SCANNING:
		cTree.setScanningMode(true)
	case m3uetcpb.CollectionEvent_CE_SCANNING_DONE:
		cTree.setScanningMode(false)
	}

	glib.IdleAdd(cd.updateCollectionModel)
	if cTree.canBeUpdated() {
		glib.IdleAdd(cTree.update)
	} else {
		glib.IdleAdd(cd.updateScanningProgress)
	}
}

func (cd *collectionData) SwitchHierarchyTo(id string, grouped bool) {
	if !cTree.canBeUpdated() {
		return
	}

	cTree.
		setHierarchy(collectionTreeHierarchyMap[id]).
		setGroupByCollection(grouped)

	glib.IdleAdd(cTree.update)
}

func (cd *collectionData) appendCDataItem(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
	// NOTE: cd.mu lock is already set
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
	// NOTE: cd.mu lock is already set
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
	// NOTE: cd.mu lock is already set
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
	if cTree.isInInitialMode() {
		return false
	}

	slog.Info("Updating collection model")

	cd.mu.RLock()
	defer cd.mu.RUnlock()

	model := collectionModel
	if model == nil {
		return false
	}

	if model.NColumns() == 0 {
		return false
	}

	_, ok := model.IterFirst()
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
			model.Set(
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
				[]glib.Value{
					*glib.NewValue(c.Id),
					*glib.NewValue(c.Name),
					*glib.NewValue(c.Description),
					*glib.NewValue(c.Location),
					*glib.NewValue(persp),
					*glib.NewValue(c.Disabled),
					*glib.NewValue(c.Remote),
					*glib.NewValue(c.Tracks),
					*glib.NewValue(tracks),
					*glib.NewValue(false),
					*glib.NewValue(false),
				},
			)
		}
	}

	return false
}

func (cd *collectionData) updateCollectionNamesMap() {
	cd.mu.RLock()
	defer cd.mu.RUnlock()

	collectionNameMap = make(map[int64]string)
	for _, c := range cd.collection {
		collectionNameMap[c.Id] = c.Name
	}
}

func (cd *collectionData) updateScanningProgress() bool {
	if cProgress == nil {
		slog.Error("Collections progress bar is unavailable")
		return false
	}

	cd.mu.RLock()
	defer cd.mu.RUnlock()

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

// CreateCollectionModel creates a collection model.
func CreateCollectionModel() (model *gtk.ListStore, err error) {
	slog.Info("Creating collection model")

	collectionModel = gtk.NewListStore(CColumns.getTypes())
	if collectionModel == nil {
		err = fmt.Errorf("failed to create list-store")
		return
	}

	model = collectionModel
	return
}

// CreateCollectionTreeModel creates a collection model.
func CreateCollectionTreeModel(h collectionTreeHierarchy) (
	model *gtk.TreeStore, err error) {

	logw := slog.With("hierarchy", h)
	logw.Info("Creating collection tree model")

	model = gtk.NewTreeStore(CTreeColumn.getTypes())
	if model == nil {
		err = fmt.Errorf("fsiled to create tree-store")
		return
	}

	cTree.
		setHierarchy(h).
		setModel(model)

	var errp error
	cProgress, errp = builder.GetProgressBar("collections_scanning_progress")
	if errp != nil {
		cProgress = nil
		logw.Error("Failed to get progress bar from builder", "error", errp)
	}
	return
}

// FilterCollectionTreeBy filters the collection tree by the given value.
func FilterCollectionTreeBy(val string) {
	cTree.
		setFilterVal(val).
		rebuild()
}

// GetCollectionModel returns the current collection model.
func GetCollectionModel() *gtk.ListStore {
	return collectionModel
}

// GetCollectionTreeModel returns the current collection tree model.
func GetCollectionTreeModel() *gtk.TreeStore {
	return cTree.getModel()
}
