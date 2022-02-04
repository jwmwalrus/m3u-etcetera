package store

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

var (
	cTree = &collectionTree{}

	collectionModel *gtk.ListStore
	cProgress       *gtk.ProgressBar

	// CData collection store
	CData struct {
		subscriptionID string
		Mu             sync.Mutex
		Collection     []*m3uetcpb.Collection
		Track          []*m3uetcpb.Track
	}
)

// CollectionOptions for the collections dialog
type CollectionOptions struct {
	Discover   bool
	UpdateTags bool
}

// SetDefaults -
func (co *CollectionOptions) SetDefaults() {}

// CollectionAlreadyExists returns true if the location and the name are not
// already in use by another collection
func CollectionAlreadyExists(location, name string) bool {
	CData.Mu.Lock()
	defer CData.Mu.Unlock()

	for _, c := range CData.Collection {
		if c.Location == location ||
			strings.ToLower(c.Name) == strings.ToLower(name) {
			return true
		}
	}
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
func CreateCollectionTreeModel(h collectionTreeHierarchy) (model *gtk.TreeStore, err error) {
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

func updateCollectionModel() bool {
	if cTree.initialMode {
		return false
	}

	log.Info("Updating collection model")

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
