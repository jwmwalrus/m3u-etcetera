package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/gear-pieces/idler"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/discover"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	rtc "github.com/jwmwalrus/rtcycler"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// Track defines a track row.
type Track struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Location    string `json:"location" gorm:"uniqueIndex:unique_idx_track_location,not null"`
	Format      string `json:"format"`
	Type        string `json:"type"`
	Title       string `json:"title" gorm:"index:idx_track_title"`
	Album       string `json:"album" gorm:"index:idx_track_album"`
	Artist      string `json:"artist" gorm:"index:idx_track_artist"`
	Albumartist string `json:"albumartist" gorm:"index:idx_track_album_artist"`
	Composer    string `json:"composer" gorm:"index:idx_track_composer"`
	Genre       string `json:"genre" gorm:"index:idx_track_genre"`

	Comment     string `json:"comment"`
	Lyrics      string `json:"lyrics"`
	Cover       string `json:"cover"`
	Year        int    `json:"year" gorm:"index:idx_track_year"`
	Tracknumber int    `json:"tracknumber"`
	Tracktotal  int    `json:"tracktotal"`
	Discnumber  int    `json:"discnumber"`
	Disctotal   int    `json:"disctotal"`
	Date        int64  `json:"date" gorm:"index:idx_track_date"`
	Duration    int64  `json:"duration"`

	Rating       int        `json:"rating" gorm:"index:idx_track_rating"`
	Playcount    int        `json:"playcount,"`
	Remote       bool       `json:"remote"` // if track is remote but not in a remote collection
	Lastplayed   int64      `json:"lastplayed"`
	Tags         string     `json:"tags"`
	CreatedAt    int64      `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt    int64      `json:"updatedAt" gorm:"autoUpdateTime:nano"`
	CollectionID int64      `json:"collectionId" gorm:"index:idx_track_collection_id,not null"`
	Collection   Collection `json:"collection" gorm:"foreignKey:CollectionID"`
}

func (t *Track) Create() error {
	return t.CreateTx(db)
}

func (t *Track) CreateTx(tx *gorm.DB) error {
	return tx.Create(t).Error
}

func (t *Track) Delete() error {
	return t.DeleteTx(db)
}

func (t *Track) DeleteTx(tx *gorm.DB) error {
	return tx.Delete(t).Error
}

func (t *Track) Read(id int64) error {
	return t.ReadTx(db, id)
}

func (t *Track) ReadTx(tx *gorm.DB, id int64) error {
	return tx.First(t, id).Error
}

func (t *Track) Save() error {
	return t.SaveTx(db)
}

func (t *Track) SaveTx(tx *gorm.DB) error {
	return tx.Save(t).Error
}

func (t *Track) ToProtobuf() proto.Message {
	bv, err := json.Marshal(t)
	if err != nil {
		slog.Error("Failed to marshal track", "error", err)
		return &m3uetcpb.Track{}
	}

	out := &m3uetcpb.Track{}
	err = jsonUnmarshaler.Unmarshal(bv, out)
	onerror.Log(err)

	// Unmatched
	if !t.Remote {
		path, err := urlstr.URLToPath(t.Location)
		if err == nil {
			if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
				out.Dangling = true
			}
		}
	}

	return out
}

// AfterCreate is a GORM hook.
func (t *Track) AfterCreate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToCollectionStoreEvent,
			subscription.Event{
				Idx:  int(CollectionEventItemAdded),
				Data: t,
			},
		)
	}()
	return nil
}

// AfterUpdate is a GORM hook.
func (t *Track) AfterUpdate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToCollectionStoreEvent,
			subscription.Event{
				Idx:  int(CollectionEventItemChanged),
				Data: t,
			},
		)
	}()
	return nil
}

// AfterDelete is a GORM hook.
func (t *Track) AfterDelete(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToCollectionStoreEvent,
			subscription.Event{
				Idx:  int(CollectionEventItemRemoved),
				Data: t,
			},
		)
	}()
	return nil
}

func (t *Track) createTransient(tx *gorm.DB, raw map[string]interface{}) (err error) {
	c, err := TransientCollection.Get()
	if err != nil {
		return
	}
	t.CollectionID = c.ID
	if err = t.updateTags(); err != nil {
		return
	}

	t.fillMissingTags(raw)
	err = t.SaveTx(tx)
	return
}

func (t *Track) discoverDuration() {
	slog.Debug("Discovering duration", "location", t.Location)
	info, err := discover.Execute(t.Location)
	if err != nil {
		slog.Error("Failed to execute `discover`", "error", err)
		return
	}

	t.Duration = info.Duration
	slog.Debug("discovered duration", "duration", time.Duration(t.Duration)*time.Nanosecond)
}

func (t *Track) fillMissingTags(raw map[string]interface{}) {
	if raw == nil {
		return
	}

	const unknownTxt = "[Unknown]"
	var title, album, artist, albumArtist, genre string
	var full, partial time.Time
	for k, v := range raw {
		str, _ := v.(string)
		switch k {
		case "TIT2":
			title = v.(string)
		case "TALB":
			album = v.(string)
		case "TPE1":
			artist = v.(string)
		case "TPE2":
			albumArtist = v.(string)
		case "TCON":
			genre = v.(string)
		case "TLEN":
			msec, err := strconv.ParseInt(str, 10, 64)
			if err == nil && t.Duration == 0 {
				t.Duration = msec * 1e6
			}
		case "TDAT", "TRDA", "TDRL", "TDRC", "TDOR", "TYER":
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

	unesc := ""
	unesc, _ = urlstr.URLToPath(t.Location)
	if strings.TrimSpace(t.Title) == "" {
		if title != "" {
			t.Title = title
		} else {
			t.Title = unknownTxt
			if unesc != "" {
				t.Title += " (" + filepath.Base(unesc) + ")"
			}
		}
	}
	if strings.TrimSpace(t.Album) == "" {
		if album != "" {
			t.Album = album
		} else {
			t.Album = unknownTxt
			if unesc != "" {
				t.Album += " (" + filepath.Dir(unesc) + ")"
			}
		}
	}
	if strings.TrimSpace(t.Artist) == "" &&
		strings.TrimSpace(t.Albumartist) == "" {
		if artist != "" {
			t.Artist = artist
		} else if albumArtist != "" {
			t.Artist = albumArtist
		} else {
			t.Artist = unknownTxt
		}
	}
	if strings.TrimSpace(t.Genre) == "" {
		if genre != "" {
			t.Genre = genre
		} else {
			t.Genre = unknownTxt
		}
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
		logw := slog.With("file", fn)

		file := filepath.Join(base.CoversDir(), fn)
		err := os.WriteFile(file, p.Data, 0644)
		if err != nil {
			logw.Error("Failed to write picture info file", "error", err)
			return
		}
		logw.Debug("Picture info saved")
		t.Cover = fn
	}
}

func (t *Track) updateTags() (err error) {
	idler.GetBusy(idler.StatusFileOperations)
	defer idler.GetFree(idler.StatusFileOperations)

	var path string
	path, err = urlstr.URLToPath(t.Location)
	if err != nil {
		slog.With(
			"location", t.Location,
			"error", err,
		).Error("Failed to convert URL to path")
		return
	}

	f, err := os.Open(path)
	if err != nil {
		slog.With(
			"path", path,
			"error", err,
		).Error("Failed to open path")
		return
	}
	defer f.Close()

	m, err2 := tag.ReadFrom(f)
	if err2 != nil {
		slog.With(
			"location", t.Location,
			"error", err2,
		).Error("Failed to read tags from file")
	}

	raw := map[string]interface{}{}
	if m != nil {
		t.Format = string(m.Format())
		t.Type = string(m.FileType())
		t.Title = m.Title()
		t.Album = m.Album()
		t.Artist = m.Artist()
		t.Albumartist = m.AlbumArtist()
		t.Composer = m.Composer()
		t.Genre = m.Genre()
		t.Year = m.Year()
		t.Tracknumber, t.Tracktotal = m.Track()
		t.Discnumber, t.Disctotal = m.Disc()

		dir := filepath.Dir(t.Location)
		hasher := md5.New()
		hasher.Write([]byte(dir))
		t.savePicture(m.Picture(), hex.EncodeToString(hasher.Sum(nil)))

		raw = m.Raw()
	}

	t.fillMissingTags(raw)

	if t.Duration == 0 {
		t.discoverDuration()
	}

	return
}

// DeleteDanglingTrack removes a (presumably) non-existent track from collection.
func DeleteDanglingTrack(t *Track, c *Collection, withRemote bool) (err error) {
	if !withRemote && (c.Remote || t.Remote) {
		return
	}

	pts := []PlaylistTrack{}
	err = db.Where("track_id = ?", t.ID).Find(&pts).Error
	if err != nil {
		return
	}

	for i := range pts {
		pl := Playlist{}
		err = db.Joins("Playbar").First(&pl, pts[i].PlaylistID).Error
		if err != nil {
			return
		}
		pl.Playbar.DeleteFromPlaylist(&pl, pts[i].Position)
	}
	err = t.Delete()
	return
}

// DeleteDanglingTrackByID removes a (presumably) non-existent track from collection.
func DeleteDanglingTrackByID(id int64, withRemote bool) error {
	t := &Track{}
	err := t.Read(id)
	if err != nil {
		return err
	}

	c := &Collection{}
	err = db.First(c, t.CollectionID).Error
	if err != nil {
		return err
	}

	return DeleteDanglingTrack(t, c, withRemote)
}

// DeleteLocalTrackIfDangling deletes the track identified by id if the given
// location does not exist.
func DeleteLocalTrackIfDangling(id int64, location string) {
	path, err := urlstr.URLToPath(location)
	if err != nil {
		slog.With(
			"location", location,
			"error", err,
		).Warn("Failed to convert URL to path")
		return
	}

	if _, err := os.Stat(path); err == nil {
		return
	}

	onerror.Warn(DeleteDanglingTrackByID(id, false))
}

// DeleteTrackIfTransient removes a track only if it belongs
// to the transient collection and is not being used anywhere else.
func DeleteTrackIfTransient(id int64) {
	t := Track{}
	err := t.Read(id)
	if err != nil {
		slog.With(
			"id", id,
			"error", err,
		).Error("Failed to read track")
		return
	}

	trc, err := TransientCollection.Get()
	if err != nil {
		slog.Error("Failed to get transient collection", "error", err)
		return
	}

	if t.CollectionID == trc.ID {
		pts := []PlaylistTrack{}
		db.Where("track_id = ?", t.ID).Find(&pts)
		if len(pts) == 0 {
			t.Delete()
		}
		return
	}
}

// FindTracksIn returns the tracks for the given IDs.
func FindTracksIn(ids []int64) (ts []*Track, notFound []int64) {
	ts = []*Track{}
	if len(ids) < 1 {
		return
	}

	list := []Track{}
	err := db.Where("id in ?", ids).Find(&list).Error
	if err != nil {
		slog.With(
			"ids", ids,
			"error", err,
		).Error("Failed to find tracks in database")
		return
	}

	actual := []int64{}
	for i := range list {
		actual = append(actual, list[i].ID)
		ts = append(ts, &list[i])
	}

	for _, id := range ids {
		if !slices.Contains(actual, id) {
			notFound = append(notFound, id)
		}
	}
	return
}

// ReadTagsForLocation returns a track containing the tags read
// for the given location.
func ReadTagsForLocation(location string) (t *Track, err error) {
	t = &Track{Location: location}
	err = t.updateTags()
	return
}
