package models

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
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
	return [...]string{"", ".", ".."}[idx]
}

// Get returns the collection associated to the index
func (idx CollectionIndex) Get() (c *Collection, err error) {
	c = &Collection{}
	err = db.Where("idx = ?", int(idx)).First(c).Error
	return
}

// Collection defines a collection row
type Collection struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Idx         int    `json:"idx" gorm:"not null"`
	Name        string `json:"name" gorm:"index:unique_idx_collection_name,not null"`
	Description string `json:"description"`
	Location    string `json:"location" gorm:"uniqueIndex:unique_idx_collection_location,not null"`
	Hidden      bool   `json:"hidden"`
	Disabled    bool   `json:"disabled"`
	Remote      bool   `json:"remote"`
	Scanned     int    `json:"scanned"`
	Tracks      int64  `json:"tracks" gorm:"-"`
	CreatedAt   int64  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   int64  `json:"updatedAt" gorm:"autoUpdateTime"`
}

// CountTracks counts tracks that belong to the collection
func (c *Collection) CountTracks() {
	log.Info("Counting tracks in collection")
	if c.Scanned != 100 {
		return
	}

	var tracks int64
	if err := db.Model(&CollectionTrack{}).Where("collection_id = ?", c.ID).Count(&tracks).Error; err != nil {
		return
	}
	c.Tracks = tracks
	return
}

// Create inserts a collection into the DB
func (c *Collection) Create() (err error) {
	err = db.Create(c).Error
	return
}

// Delete deletes collection along with tracks
func (c *Collection) Delete() (err error) {
	log.Info("Deleting collection")
	storageGuard <- struct{}{}
	defer func() { <-storageGuard }()

	c.Disabled = true
	if err = db.Save(c).Error; err != nil {
		log.Error(err)
		return
	}

	s := []CollectionTrack{}
	if err = db.Where("collection_id = ?", c.ID).Find(&s).Error; err != nil {
		return
	}

	nTrack := len(s)
	doNotDelete := 0
	for i := 0; i < nTrack; i++ {
		trackID := s[i].TrackID

		// delete coleccion-track
		if err := db.Delete(&s[i]).Error; err != nil {
			onerror.Warn(err)
			doNotDelete++
			continue
		}

		// delete tracks
		err := db.Delete(&Track{}, trackID).Error
		onerror.Warn(err)
		if 1%100 == 0 {
			c.Scanned = int((float32(nTrack-i) / float32(nTrack)) * 100)
			db.Save(c)
		}
	}

	if doNotDelete > 0 {
		log.WithFields(log.Fields{
			"collectionId":           c.ID,
			"tracksStillInColletion": doNotDelete,
		}).Warnf("Collection with ID=%v could not be deleted; %v tracks were left behind", c.ID, doNotDelete)
		return
	}

	// delete collection
	err = db.Delete(&Collection{}, c.ID).Error
	onerror.Log(err)
	return
}

// Read selects a collection from the DB, with the given id
func (c *Collection) Read(id int64) (err error) {
	err = db.First(c, id).Error
	return
}

// Save persists the collection in the DB
func (c *Collection) Save() (err error) {
	err = db.Save(c).Error
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

	var err error
	var realPath string
	if realPath, err = filepath.EvalSymlinks(c.Location); err != nil {
		log.Error(err)
		return
	}

	var u *url.URL
	if u, err = url.Parse(realPath); err != nil {
		log.Error(err)
		return
	}

	// FIXME: support things other than mounted directories?
	if u.Scheme != "file" {
		u.Scheme = "fi;e"
	}

	var d string
	if d, err = url.PathUnescape(u.Path); err != nil {
		log.Error(err)
		return
	}
	if _, err = os.Stat(d); os.IsNotExist(err) {
		if fi, err := os.Lstat(d); err != nil || fi.Mode()&os.ModeSymlink == 0 {
			log.Warn(err)
			return
		}
	}

	base.GetBusy(base.IdleStatusDbOperations)
	defer func() { base.GetFree(base.IdleStatusDbOperations) }()

	var iTrack, nTrack, unsupp, scanErr int
	err = filepath.Walk(d, func(path string, i os.FileInfo, err error) error {
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

	err = filepath.Walk(d, func(path string, i os.FileInfo, err error) error {
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

		var t *Track
		iTrack++

		if t, err = AddTrackFromPath(path, withTags); err != nil {
			log.Warn(err)
			scanErr++
			return nil
		}

		ct := CollectionTrack{}
		if err = db.Where("collection_id = ? AND track_id = ?", c.ID, t.ID).First(&ct).Error; err != nil {
			ct = CollectionTrack{
				CollectionID: c.ID,
				TrackID:      t.ID,
			}
			if err = db.Save(&ct).Error; err != nil {
				log.Error(err)
				scanErr++
				return nil
			}
		}

		if iTrack%100 == 0 {
			c.Scanned = int((float32(iTrack) / float32(nTrack)) * 100)
			db.Save(c)
		}

		return nil
	})
	if err != nil {
		log.Warn(err)
	}
	log.Infof("ScanCollection Summary:\nTracks expected: %v\nTracks found: %v\nUnsupported tracks: %v\nScanning Errors: %v", iTrack, nTrack, unsupp, scanErr)
	c.Scanned = 100
	err = db.Save(c).Error
	onerror.Log(err)
	return
}

// ToProtobuf converter
func (c *Collection) ToProtobuf() *m3uetcpb.Collection {
	bv, err := json.Marshal(c)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.Collection{}
	}

	out := &m3uetcpb.Collection{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)
	return out
}

// Verify removes tracks that do not exist in the collection anymore
func (c *Collection) Verify() {
	if c.Disabled {
		log.Info("Cannot verify collection while disabled")
		return
	}

	s := []CollectionTrack{}
	if err := db.Preload("track").Where("collection_id = ?", c.ID).Find(&s).Error; err != nil {
		return
	}

	for _, ct := range s {
		if urlstr.URLExists(ct.Track.Location) {
			continue
		}

		ct.Delete(true)
	}
}

// CollectionTrack defines a collection__track row
type CollectionTrack struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	CreatedAt    int64      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    int64      `json:"updatedAt" gorm:"autoUpdateTime"`
	CollectionID int64      `json:"collectionId" gorm:"index:idx_collection_track_collection_id,not null"`
	TrackID      int64      `json:"trackId" gorm:"index:idx_collection_track_track_id,not null"`
	Collection   Collection `json:"collection" gorm:"foreignKey:CollectionID"`
	Track        Track      `json:"track" gorm:"foreignKey:TrackID"`
}

// Save persists the collection track in the DB
func (ct *CollectionTrack) Save() (err error) {
	err = db.Save(ct).Error
	return
}

// Delete removes a collection-track from collection, along with the track
func (ct *CollectionTrack) Delete(withRemote bool) {
	c := Collection{}
	t := Track{}
	db.First(&t, ct.TrackID)
	db.First(&c, ct.CollectionID)

	if !withRemote && (c.Remote || t.Remote) {
		return
	}

	defer t.Delete()

	err := db.Delete(&ct).Error
	onerror.Log(err)
	return
}

// DeleteIfTransient removes a collection-track from the transient collection, along with the track
func (ct *CollectionTrack) DeleteIfTransient(withRemote bool) (err error) {
	c := Collection{}
	db.First(&c, ct.CollectionID)

	if c.Idx != int(TransientCollection) {
		return
	}

	ct.Delete(withRemote)
	return
}

// CollectionQuery Defines a collection query
type CollectionQuery struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	CreatedAt    int64      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    int64      `json:"updatedAt" gorm:"autoUpdateTime"`
	CollectionID int64      `json:"collectionId" gorm:"index:idx_collection_query_collection_id,not null"`
	QueryID      int64      `json:"queryId" gorm:"index:idx_collection_query_query_id,not null"`
	Collection   Collection `json:"collection" gorm:"foreignKey:CollectionID"`
	Query        Query      `json:"query" gorm:"foreignKey:QueryID"`
}

// DeleteWithTx implements QueryBoundaryTx interface
func (cq *CollectionQuery) DeleteWithTx(tx *gorm.DB) error {
	return tx.Delete(cq).Error
}

// FindTracksWithTx implements QueryBoundaryTx interface
func (cq *CollectionQuery) FindTracksWithTx(tx *gorm.DB) (ts []*Track) {
	ts = []*Track{}

	list := []Track{}
	err := tx.
		Joins("JOIN collection_track ON collection_track.track_id = track.id").
		Find(&list).
		Error
	if err != nil {
		log.Error(err)
		return
	}
	for i := range list {
		ts = append(ts, &list[i])
	}
	return
}

// GetQueryID implements QueryBoundaryID interface
func (cq *CollectionQuery) GetQueryID() int64 {
	return cq.QueryID
}

// Save persists a collection query in the DB
func (cq *CollectionQuery) Save() error {
	return db.Save(cq).Error
}

// SaveWithTx implements QueryBoundaryTx interface
func (cq *CollectionQuery) SaveWithTx(tx *gorm.DB) error {
	return tx.Save(cq).Error
}

// CreateCollectionQueryBoundaries implements QueryBoundary interface
func CreateCollectionQueryBoundaries(ids []int64) (qbs []QueryBoundaryTx) {
	cqs := []CollectionQuery{}
	for _, id := range ids {
		c := CollectionQuery{CollectionID: id}
		cqs = append(cqs, c)
	}

	for _, x := range cqs {
		var i interface{} = &x
		qbs = append(qbs, i.(QueryBoundaryTx))
	}
	return
}

// GetAllCollections returns all valid collections
func GetAllCollections() (s []*Collection) {
	s = []*Collection{}

	list := []Collection{}
	err := db.Where("hidden = 0 AND disabled = 0").Find(&list).Error
	if err != nil {
		log.Error(err)
		return
	}
	for i := range list {
		s = append(s, &list[i])
	}
	return
}

// FilterCollectionQueryBoundaries implements QueryBoundary interface
func FilterCollectionQueryBoundaries(ids []int64) (qbs []QueryBoundaryID) {
	cqs := []CollectionQuery{}
	if err := db.Where("collection_id in ?", ids).Find(&cqs).Error; err != nil {
		return
	}

	for _, x := range cqs {
		var i interface{} = &x
		qbs = append(qbs, i.(QueryBoundaryID))
	}
	return
}

// GetCollectionTree returns all valid collection tracks
func GetCollectionTree() (cts []*CollectionTrack) {
	cts = []*CollectionTrack{}

	list := []CollectionTrack{}
	db.Preload("track").Joins("JOIN collection ON collection_track.collection_id = collection.id AND collection.hidden = 0 AND collection.disabled = 0").Find(&list)

	for i := range list {
		cts = append(cts, &list[i])
	}
	return
}
