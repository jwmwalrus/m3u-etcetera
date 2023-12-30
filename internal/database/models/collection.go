package models

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/pointers"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/gear-pieces/idler"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	rtc "github.com/jwmwalrus/rtcycler"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// CollectionIndex defines indexes for collections.
type CollectionIndex int

const (
	// DefaultCollection for the default collection.
	DefaultCollection CollectionIndex = iota + 1

	// TransientCollection for the transient collection.
	TransientCollection
)

func (idx CollectionIndex) String() string {
	return [...]string{"", "\t", "\t\t"}[idx]
}

// Get returns the collection associated to the index.
func (idx CollectionIndex) Get() (c *Collection, err error) {
	c = &Collection{}
	err = db.Where("idx = ?", int(idx)).First(c).Error
	return
}

// CollectionEvent defines a collection event.
type CollectionEvent int

// CollectionEvent enum.
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

// Collection defines a collection row.
type Collection struct {
	Model
	Idx            int         `json:"idx" gorm:"not null"`
	Name           string      `json:"name" gorm:"uniqueIndex:unique_idx_collection_name,not null"`
	Description    string      `json:"description"`
	Location       string      `json:"location" gorm:"uniqueIndex:unique_idx_collection_location,not null"`
	Remotelocation string      `json:"remoteLocation"`
	Hidden         bool        `json:"hidden"`
	Disabled       bool        `json:"disabled"`
	Remote         bool        `json:"remote"`
	Scanned        int         `json:"scanned"`
	Tracks         int64       `json:"tracks" gorm:"-"`
	PerspectiveID  int64       `json:"perspectiveId" gorm:"index:idx_collection_perspective_id,not null"`
	Perspective    Perspective `json:"-" gorm:"foreignKey:PerspectiveID"`
}

func (c *Collection) Create() error {
	return c.CreateTx(db)
}

func (c *Collection) CreateTx(tx *gorm.DB) error {
	return tx.Create(c).Error
}

func (c *Collection) Delete() (err error) {
	return c.DeleteTx(db)
}

func (c *Collection) DeleteTx(tx *gorm.DB) (err error) {
	slog.Info("Deleting collection")

	storageGuard <- struct{}{}
	defer func() { <-storageGuard }()

	c.Disabled = true
	if err = c.Save(); err != nil {
		slog.Error("Failed to save collection", "error", err)
		return
	}

	s := []Track{}
	if err = tx.Where("collection_id = ?", c.ID).Find(&s).Error; err != nil {
		return
	}

	nTrack := len(s)
	doNotDelete := 0
	for i := 0; i < nTrack; i++ {
		err := DeleteDanglingTrack(&s[i], c, true)
		if err != nil {
			slog.With(
				"collection", c.Name,
				"track", s[i],
				"error", err,
			).Warn("Failed to delete dangling track")
			doNotDelete++
			continue
		}
		if i%100 == 0 {
			c.Scanned = int((float32(nTrack-i) / float32(nTrack)) * 100)
			tx.Save(c)
		}
	}

	if doNotDelete > 0 {
		slog.With(
			"collection_id", c.ID,
			"tracks-still-in-colletion", doNotDelete,
		).Warn("Collection could not be deleted")
		return
	}

	// delete collection
	err = tx.Delete(c).Error
	return
}

func (c *Collection) Read(id int64) error {
	return c.ReadTx(db, id)
}

func (c *Collection) ReadTx(tx *gorm.DB, id int64) error {
	return tx.Joins("Perspective").
		First(c, id).
		Error
}

func (c *Collection) Save() error {
	return c.SaveTx(db)
}

func (c *Collection) SaveTx(tx *gorm.DB) error {
	return tx.Save(c).Error
}

func (c *Collection) ToProtobuf() proto.Message {
	return &m3uetcpb.Collection{
		Id:             c.ID,
		Name:           c.Name,
		Description:    c.Description,
		Location:       c.Location,
		RemoteLocation: c.Remotelocation,
		Disabled:       c.Disabled,
		Remote:         c.Remote,
		Scanned:        int32(c.Scanned),
		Tracks:         c.Tracks,
		Perspective:    m3uetcpb.Perspective(c.Perspective.Idx),
		CreatedAt:      timestamppb.New(time.Unix(0, c.CreatedAt)),
		UpdatedAt:      timestamppb.New(time.Unix(0, c.UpdatedAt)),
	}
}

// AfterCreate is a GORM hook.
func (c *Collection) AfterCreate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
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

// AfterUpdate is a GORM hook.
func (c *Collection) AfterUpdate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
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

// AfterDelete is a GORM hook.
func (c *Collection) AfterDelete(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
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

// CountTracks counts tracks that belong to the collection.
func (c *Collection) CountTracks() {
	if c.Scanned != 100 {
		return
	}

	slog.Info("Counting tracks in collection", "collection", *c)

	var tracks int64
	err := db.Model(&Track{}).
		Where("collection_id = ?", c.ID).
		Count(&tracks).
		Error
	if err != nil {
		return
	}
	c.Tracks = tracks
}

// Scan adds tracks to collection.
func (c *Collection) Scan(withTags bool) {
	logw := slog.With(
		"c", *c,
		"with_tags", withTags,
	)

	if c.Disabled {
		logw.Info("Cannot scan collection while disabled")
		return
	}

	logw.Info("Scanning collection")

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
		logw.Error("Failed to convert URL to path", "error", err)
		return
	}

	if _, err = os.Stat(rootDir); errors.Is(err, os.ErrNotExist) {
		logw.Warn("Failed to stat collection's location", "error", err)
		return
	}

	idler.GetBusy(idler.StatusDbOperations)
	defer func() { idler.GetFree(idler.StatusDbOperations) }()

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

	onerrorw := onerror.NewRecorder(logw)
	onerrorw.Warn(err)

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
				logw.With(
					"path", path,
					"extension", filepath.Ext(path),
				).
					Info("Unsupported file")

				unsupp++
			}
			return nil
		}

		iTrack++

		if _, err = c.addTrackFromPath(tx, path, withTags); err != nil {
			logw.Warn("Failed to add track from path", "error", err)
			scanErr++
			return nil
		}

		if iTrack%100 == 0 {
			c.Scanned = int((float32(iTrack) / float32(nTrack)) * 100)
			onerrorw.Log(c.Save())
		}

		return nil
	})
	onerrorw.Warn(err)
	logw = logw.With(
		"tracks-expected", nTrack,
		"tracks-found", iTrack,
		"unsupported-tracks", unsupp,
		"scanning-errors-count", scanErr,
	)
	logw.Info("ScanCollection Summary")

	c.Scanned = 100
	onerror.NewRecorder(logw).Log(c.Save())
}

// Verify removes tracks that do not exist in the collection anymore.
func (c *Collection) Verify() {
	logw := slog.With("c", *c)
	logw.Info("Verifying collection")

	if c.Disabled {
		logw.Info("Cannot verify collection while disabled")
		return
	}

	s := []Track{}
	if err := db.Where("collection_id = ?", c.ID).Find(&s).Error; err != nil {
		return
	}

	for i := range s {
		if urlstr.URLExists(s[i].Location) {
			continue
		}

		DeleteDanglingTrack(&s[i], c, true)
	}
}

func (c *Collection) addTrackFromLocation(tx *gorm.DB, location string,
	withTags bool) (t *Track, err error) {

	logw := slog.With("location", location)

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
			if newt.CollectionID != c.ID {
				err = fmt.Errorf("Track already belongs to another collection")
				return
			}
			logw.Info("Track already in a collection", "collection", c.Name)
		} else {
			logw.Info("Reusing transient track")
			newt.CollectionID = c.ID
		}
	}

	t = newt

	if withTags || doTag {
		err2 := t.updateTags()
		if err2 != nil {
			logw.Warn("Failed to update tags", "error", err2)
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

// GetAllCollections returns all valid collections.
func GetAllCollections() []*Collection {
	s := []Collection{}
	err := db.Where("hidden = 0").Find(&s).Error
	if err != nil {
		slog.Error("Failed to find all collections in database", "error", err)
		return []*Collection{}
	}
	return pointers.FromSlice(s)
}

// GetCollectionStore returns all valid collection tracks.
func GetCollectionStore() ([]*Collection, []*Track) {
	cs := []Collection{}
	ts := []Track{}

	db.Joins(
		"JOIN collection ON track.collection_id = collection.id AND collection.hidden = 0 AND collection.disabled = 0",
	).
		Find(&ts)

	db.Joins("Perspective").
		// Joins("JOIN perspective ON collection.perspective_id = perspective.id").
		Where("hidden = 0").
		Find(&cs)

	for i := range cs {
		cs[i].CountTracks()
	}

	return pointers.FromSlice(cs), pointers.FromSlice(ts)
}
