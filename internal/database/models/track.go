package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dhowden/tag"
	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/onerror"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

// Track defines a track row
type Track struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Location    string `json:"location" gorm:"uniqueIndex:unique_idx_track_location,not null"`
	Format      string `json:"format"`
	Type        string `json:"type"`
	Title       string `json:"title" gorm:"index:idx_track_title"`
	Album       string `json:"album" gorm:"index:idx_track_album"`
	Artist      string `json:"artist" gorm:"index:idx_track_artist"`
	Albumartist string `json:"albumArtist" gorm:"index:idx_track_album_artist"`
	Composer    string `json:"composer" gorm:"index:idx_track_composer"`
	Genre       string `json:"genre" gorm:"index:idx_track_genre"`

	Comment     string `json:"comment"`
	Lyrics      string `json:"lyrics"`
	Cover       string `json:"cover"`
	Year        int    `json:"year" gorm:"index:idx_track_year"`
	Tracknumber int    `json:"trackNumber"`
	Tracktotal  int    `json:"trackTotal"`
	Discnumber  int    `json:"discNumber"`
	Disctotal   int    `json:"discTotal"`
	Date        int64  `json:"date" gorm:"index:idx_track_date"`
	Duration    int64  `json:"duration"`

	Rating       int        `json:"rating" gorm:"index:idx_track_rating"`
	Playcount    int        `json:"playCount"`
	Remote       bool       `json:"remote"` // if track is remote but not in a remote collection
	Lastplayed   int64      `json:"lastPlayed"`
	Tags         string     `json:"tags"`
	CollectionID int64      `json:"collectionId" gorm:"index:idx_track_collection_id,not null"`
	Collection   Collection `json:"collection" gorm:"foreignKey:CollectionID"`
	CreatedAt    int64      `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt    int64      `json:"updatedAt" gorm:"autoUpdateTime:nano"`
}

// Create implements the DataCreator interface
func (t *Track) Create() (err error) {
	err = db.Create(t).Error
	return
}

// Delete implements the DataDeleter interface
func (t *Track) Delete() {
	err := db.Delete(&Track{}, t.ID).Error
	onerror.Log(err)
}

// Read implements the DataReader interface
func (t *Track) Read(id int64) (err error) {
	err = db.First(t, id).Error
	return
}

// Save implements the DataUpdater interface
func (t *Track) Save() (err error) {
	err = db.Save(t).Error
	return
}

// ToProtobuf implements the ProtoOut interface
func (t *Track) ToProtobuf() proto.Message {
	bv, err := json.Marshal(t)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.Track{}
	}

	out := &m3uetcpb.Track{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)

	// Unmatched
	out.Albumartist = t.Albumartist
	out.Tracknumber = int32(t.Tracknumber)
	out.Tracktotal = int32(t.Tracktotal)
	out.Discnumber = int32(t.Discnumber)
	out.Disctotal = int32(t.Disctotal)
	out.Playcount = int32(t.Playcount)
	out.Lastplayed = t.Lastplayed
	out.CreatedAt = t.CreatedAt
	out.UpdatedAt = t.UpdatedAt
	return out
}

// FindBy implements the DataFinder interface
func (t *Track) FindBy(query interface{}) (err error) {
	err = db.Where(query).First(t).Error
	return
}

// AfterCreate is a GORM hook
func (t *Track) AfterCreate(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode && globalCollectionEvent != CollectionEventInitial {
			subscription.Broadcast(
				subscription.ToCollectionStoreEvent,
				subscription.Event{
					Idx:  int(CollectionEventItemAdded),
					Data: t,
				},
			)
		}
	}()
	return nil
}

// AfterSave is a GORM hook
func (t *Track) AfterSave(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode && globalCollectionEvent != CollectionEventInitial {
			subscription.Broadcast(
				subscription.ToCollectionStoreEvent,
				subscription.Event{
					Idx:  int(CollectionEventItemChanged),
					Data: t,
				},
			)
		}
	}()
	return nil
}

// AfterDelete is a GORM hook
func (t *Track) AfterDelete(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode && globalCollectionEvent != CollectionEventInitial {
			subscription.Broadcast(
				subscription.ToCollectionStoreEvent,
				subscription.Event{
					Idx:  int(CollectionEventItemRemoved),
					Data: t,
				},
			)
		}
	}()
	return nil
}

// DeleteWithRemote removes a collection-track from collection, along with the track
func (t *Track) DeleteWithRemote(withRemote bool) {
	c := Collection{}
	db.First(&c, t.CollectionID)

	if !withRemote && (c.Remote || t.Remote) {
		return
	}

	defer t.Delete()
	return
}

func (t *Track) fillMissingTags(raw map[string]interface{}) {
	var full, partial time.Time
	for k, v := range raw {
		str, _ := v.(string)
		switch k {
		case "TDAT", "TRDA", "TDRL", "TDRC", "TDOR":
			var dt time.Time
			dt, err := time.Parse("2006-01-02", str)
			if err != nil && str != "0000" {
				dt, err = time.Parse("2006-01-02", str+"-01-01")
				if err != nil {
					dt, err = time.Parse("2006-01-02", str+"-01")
					if err == nil && partial.IsZero() {
						partial = dt
					}
				} else if partial.IsZero() {
					partial = dt
				}
			} else if full.IsZero() {
				full = dt
			}
		case "TLEN":
			msec, err := strconv.ParseInt(str, 10, 64)
			if err == nil && t.Duration == 0 {
				t.Duration = msec * 1e6
			}
		default:
		}
	}
	var year int
	var unano int64
	if !full.IsZero() {
		year = full.Year()
		unano = full.UnixNano()
	} else if !partial.IsZero() {
		year = partial.Year()
		unano = partial.UnixNano()
	}
	if year != 0 && t.Year == 0 {
		t.Year = year
	}
	if unano != 0 && t.Date == 0 {
		t.Date = unano
	}
}

func (t *Track) savePicture(p *tag.Picture, sum string) {
	if p == nil || base.Conf.Server.Collection.Scanning.SkipCover {
		return
	}

	if p.Ext != "" && p.MIMEType != "" {
		fn := t.Cover
		if fn == "" {
			fn = sum + "." + p.Ext
		}
		file := filepath.Join(base.CoversDir, fn)
		err := ioutil.WriteFile(file, p.Data, 0644)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debugf("Picture info saved as %v", fn)
		t.Cover = fn
	}
}
func (t *Track) updateTags() (err error) {
	base.GetBusy(base.IdleStatusFileOperations)
	defer base.GetFree(base.IdleStatusFileOperations)

	var path string
	path, err = urlstr.URLToPath(t.Location)
	if err != nil {
		log.Error(err)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		log.Error(err)
		return err
	}

	t.Format = string(m.Format())
	t.Type = string(m.FileType())
	t.Title = m.Title()
	t.Album = m.Album()
	t.Artist = m.Artist()
	t.Albumartist = m.AlbumArtist()
	t.Composer = m.Composer()
	t.Genre = m.Genre()
	// t.Comment = m.Comment()
	// t.Lyrics = m.Lyrics()
	t.Year = m.Year()
	t.Tracknumber, t.Tracktotal = m.Track()
	t.Discnumber, t.Disctotal = m.Disc()
	t.fillMissingTags(m.Raw())

	dir := filepath.Dir(t.Location)
	hasher := md5.New()
	hasher.Write([]byte(dir))
	t.savePicture(m.Picture(), hex.EncodeToString(hasher.Sum(nil)))

	return
}

// FindTracksIn returns the tracks for the given IDs
func FindTracksIn(ids []int64) (ts []*Track, notFound []int64) {
	ts = []*Track{}
	if len(ids) < 1 {
		return
	}

	list := []Track{}
	err := db.Where("id in ?", ids).Find(&list).Error
	if err != nil {
		log.Error(err)
		return
	}

	actual := []int64{}
	for i := range list {
		actual = append(actual, list[i].ID)
		ts = append(ts, &list[i])
	}

	for _, id := range ids {
		if !slice.Contains(actual, id) {
			notFound = append(notFound, id)
		}
	}
	return
}

// ReadTagsForLocation returns a track containing the tags read
// for the given location
func ReadTagsForLocation(location string) (t *Track, err error) {
	t = &Track{Location: location}
	err = t.updateTags()
	return
}

func appendToTrackList(list []*Track, ts []*Track) []*Track {
	for i := range ts {
		found := false
		for j := range list {
			if list[j].ID == ts[i].ID {
				found = true
				break
			}
		}
		if !found {
			list = append(list, ts[i])
		}
	}
	return list
}
