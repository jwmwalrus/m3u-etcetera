package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// CollectionIndex defines indexes for collections
type CollectionIndex int

const (
	// DefaultCollection for the default collection
	DefaultCollection CollectionIndex = iota + 1

	// TransientCollection for the transient collection
	TransientCollection
)

func (idx CollectionIndex) String() string {
	return [...]string{"", "\t", "\t\t"}[idx]
}

// Get returns the collection associated to the index
func (idx CollectionIndex) Get() (c *Collection, err error) {
	c = &Collection{}
	err = db.Where("idx = ?", int(idx)).First(c).Error
	return
}

// CollectionEvent defines a collection event
type CollectionEvent int

// CollectionEvent enum
const (
	CollectionEventNone CollectionEvent = iota
	CollectionEventInitial
	_
	_
	CollectionEventItemAdded
	CollectionEventItemChanged
	CollectionEventItemRemoved
	CollectionEventScanning
	CollectionEventScanningDone
)

func (ce CollectionEvent) String() string {
	return []string{
		"none",
		"initial",
		"initial-item",
		"initial-done",
		"item-added",
		"item-changed",
		"item-removed",
		"scanning",
		"scanning-done",
	}[ce]
}

// Collection defines a collection row
type Collection struct {
	ID             int64       `json:"id" gorm:"primaryKey"`
	Idx            int         `json:"idx" gorm:"not null"`
	Name           string      `json:"name" gorm:"index:unique_idx_collection_name,not null"`
	Description    string      `json:"description"`
	Location       string      `json:"location" gorm:"uniqueIndex:unique_idx_collection_location,not null"`
	Remotelocation string      `json:"remoteLocation"`
	Hidden         bool        `json:"hidden"`
	Disabled       bool        `json:"disabled"`
	Remote         bool        `json:"remote"`
	Scanned        int         `json:"scanned"`
	Tracks         int64       `json:"tracks" gorm:"-"`
	PerspectiveID  int64       `json:"perspectiveId" gorm:"index:idx_collection_perspective_id,not null"`
	Perspective    Perspective `json:"perspective" gorm:"foreignKey:PerspectiveID"`
	CreatedAt      int64       `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt      int64       `json:"updatedAt" gorm:"autoUpdateTime:nano"`
}

// Create implements DataCreator interface
func (c *Collection) Create() error {
	return db.Create(c).Error
}

// Delete implements DataDeleter interface
func (c *Collection) Delete() (err error) {
	log.Info("Deleting collection")

	// TODO: broadcast scanning event

	storageGuard <- struct{}{}
	defer func() { <-storageGuard }()

	c.Disabled = true
	if err = c.Save(); err != nil {
		log.Error(err)
		return
	}

	s := []Track{}
	if err = db.Where("collection_id = ?", c.ID).Find(&s).Error; err != nil {
		return
	}

	nTrack := len(s)
	doNotDelete := 0
	for i := 0; i < nTrack; i++ {
		// delete track
		if err := db.Delete(&s[i]).Error; err != nil {
			onerror.Warn(err)
			doNotDelete++
			continue
		}
		if 1%100 == 0 {
			c.Scanned = int((float32(nTrack-i) / float32(nTrack)) * 100)
			db.Save(c)
		}
	}

	if doNotDelete > 0 {
		log.WithFields(log.Fields{
			"collectionId":           c.ID,
			"tracksStillInColletion": doNotDelete,
		}).Warnf(
			"Collection with ID=%v could not be deleted; %v tracks were left behind",
			c.ID,
			doNotDelete,
		)
		return
	}

	// delete collection
	err = c.Delete()
	return
}

// Read implements DataReader interface
func (c *Collection) Read(id int64) error {
	return db.Joins("Perspective").
		First(c, id).
		Error
}

// Save implements DataUpdater interface
func (c *Collection) Save() error {
	return db.Save(c).Error
}

// ToProtobuf implements ProtoOut interface
func (c *Collection) ToProtobuf() proto.Message {
	bv, err := json.Marshal(c)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.Collection{}
	}

	out := &m3uetcpb.Collection{}
	json.Unmarshal(bv, out)

	// Unmatched
	out.RemoteLocation = c.Remotelocation
	out.Perspective = m3uetcpb.Perspective(c.Perspective.Idx)
	out.CreatedAt = c.CreatedAt
	out.UpdatedAt = c.UpdatedAt
	return out
}

// AfterCreate is a GORM hook
func (c *Collection) AfterCreate(tx *gorm.DB) error {
	go func() {
		if base.FlagTestingMode {
			return
		}
		subscription.Broadcast(
			subscription.ToCollectionStoreEvent,
			subscription.Event{
				Idx:  int(CollectionEventItemAdded),
				Data: c,
			},
		)
	}()
	return nil
}

// AfterUpdate is a GORM hook
func (c *Collection) AfterUpdate(tx *gorm.DB) error {
	go func() {
		if base.FlagTestingMode {
			return
		}
		subscription.Broadcast(
			subscription.ToCollectionStoreEvent,
			subscription.Event{
				Idx:  int(CollectionEventItemChanged),
				Data: c,
			},
		)
	}()
	return nil
}

// AfterDelete is a GORM hook
func (c *Collection) AfterDelete(tx *gorm.DB) error {
	go func() {
		if base.FlagTestingMode {
			return
		}
		subscription.Broadcast(
			subscription.ToCollectionStoreEvent,
			subscription.Event{
				Idx:  int(CollectionEventItemRemoved),
				Data: c,
			},
		)
	}()
	return nil
}

// CountTracks counts tracks that belong to the collection
func (c *Collection) CountTracks() {
	if c.Scanned != 100 {
		return
	}

	log.Info("Counting tracks in collection")

	var tracks int64
	err := db.Model(&Track{}).
		Where("collection_id = ?", c.ID).
		Count(&tracks).
		Error
	if err != nil {
		return
	}
	c.Tracks = tracks
	return
}

// Scan adds tracks to collection
func (c *Collection) Scan(withTags bool) {
	if c.Disabled {
		log.Info("Cannot scan collection while disabled")
		return
	}

	log.Info("Scanning collection")

	storageGuard <- struct{}{}
	defer func() { <-storageGuard }()

	subscription.Broadcast(
		subscription.ToCollectionStoreEvent,
		subscription.Event{Idx: int(CollectionEventScanning)})
	defer func() {
		subscription.Broadcast(
			subscription.ToCollectionStoreEvent,
			subscription.Event{Idx: int(CollectionEventScanningDone)})
	}()

	rootDir, err := urlstr.URLToPath(c.Location)
	if err != nil {
		log.Error(err)
		return
	}

	if _, err = os.Stat(rootDir); os.IsNotExist(err) {
		log.Warn(err)
		return
	}

	base.GetBusy(base.IdleStatusDbOperations)
	defer func() { base.GetFree(base.IdleStatusDbOperations) }()

	var iTrack, nTrack, unsupp, scanErr int
	err = filepath.Walk(rootDir, func(path string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if i.IsDir() {
			return nil
		}
		nTrack++
		return nil
	})

	if nTrack == 0 {
		return
	}

	tx := db.Session(&gorm.Session{SkipHooks: true})
	err = filepath.Walk(rootDir, func(path string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if i.IsDir() {
			return nil
		}

		if !base.IsSupportedFile(path) {
			if !base.IsIgnoredFile(path) {
				log.WithFields(log.Fields{
					"path":      path,
					"extension": filepath.Ext(path),
				}).
					Info("Unsupported file:")

				unsupp++
			}
			return nil
		}

		iTrack++

		if _, err = c.addTrackFromPath(tx, path, withTags); err != nil {
			log.Warn(err)
			scanErr++
			return nil
		}

		if iTrack%100 == 0 {
			c.Scanned = int((float32(iTrack) / float32(nTrack)) * 100)
			onerror.Log(c.Save())
		}

		return nil
	})
	if err != nil {
		log.Warn(err)
	}
	log.WithFields(log.Fields{
		"tracksExpected":      nTrack,
		"tracksFound":         iTrack,
		"unsupportedTracks":   unsupp,
		"scanningErrorsCount": scanErr,
	}).
		Info("ScanCollection Summary")

	c.Scanned = 100
	onerror.Log(c.Save())
	return
}

// Verify removes tracks that do not exist in the collection anymore
func (c *Collection) Verify() {
	log.WithField("c", *c).
		Info("Verifying collection")

	if c.Disabled {
		log.Info("Cannot verify collection while disabled")
		return
	}

	s := []Track{}
	if err := db.Where("collection_id = ?", c.ID).Find(&s).Error; err != nil {
		return
	}

	for _, t := range s {
		if urlstr.URLExists(t.Location) {
			continue
		}

		DeleteDanglingTrack(t.ID, true)
	}
}

func (c *Collection) addTrackFromLocation(tx *gorm.DB, location string,
	withTags bool) (t *Track, err error) {

	doTag := false
	newt := &Track{}
	err2 := tx.Where("location = ?", location).First(newt).Error
	if err2 != nil {
		newt = &Track{
			Location:     location,
			CollectionID: c.ID,
		}
		doTag = true
	} else {
		trColl, err2 := TransientCollection.Get()
		if err2 != nil {
			err = err2
			return
		}

		if newt.CollectionID != trColl.ID {
			err = fmt.Errorf("Track already belongs to another collection")
			return
		}
		log.WithField("location", newt.Location).Info("Reusing transient track")
		newt.CollectionID = c.ID
	}

	t = newt

	if withTags || doTag {
		err2 := t.updateTags()
		if err2 != nil {
			log.WithField("location", t.Location)
		}
	}

	err = t.SaveTx(tx)
	return
}

func (c *Collection) addTrackFromPath(tx *gorm.DB, path string,
	withTags bool) (t *Track, err error) {

	var u string
	if u, err = urlstr.PathToURL(path); err != nil {
		return
	}

	t, err = c.addTrackFromLocation(tx, u, withTags)
	return
}

// GetAllCollections returns all valid collections
func GetAllCollections() (s []*Collection) {
	s = []*Collection{}

	list := []Collection{}
	err := db.Where("hidden = 0").Find(&list).Error
	if err != nil {
		log.Error(err)
		return
	}
	for i := range list {
		s = append(s, &list[i])
	}
	return
}

// GetCollectionStore returns all valid collection tracks
func GetCollectionStore() (cs []*Collection, ts []*Track) {
	cs = []*Collection{}
	ts = []*Track{}

	cList := []Collection{}
	tList := []Track{}

	db.Joins(
		"JOIN collection ON track.collection_id = collection.id AND collection.hidden = 0 AND collection.disabled = 0",
	).
		Find(&tList)

	db.Joins("Perspective").
		// Joins("JOIN perspective ON collection.perspective_id = perspective.id").
		Where("hidden = 0").
		Find(&cList)

	for i := range tList {
		ts = append(ts, &tList[i])
	}

	for i := range cList {
		cList[i].CountTracks()
		cs = append(cs, &cList[i])
	}
	return
}
