package models

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/dhowden/tag"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"google.golang.org/protobuf/proto"

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
	Year        int    `json:"year" gorm:"index:idx_track_year"`
	Tracknumber int    `json:"trackNumber"`
	Tracktotal  int    `json:"trackTotal"`
	Discnumber  int    `json:"discNumber"`
	Disctotal   int    `json:"discTotal"`
	Date        int64  `json:"date" gorm:"index:idx_track_date"`
	Duration    int64  `json:"duration"`

	Rating     int    `json:"rating" gorm:"index:idx_track_rating"`
	Playcount  int    `json:"playCount"`
	Remote     bool   `json:"remote"` // if track is remote but not in a remote collection
	Lastplayed int64  `json:"lastPlayed"`
	Tags       string `json:"tags"`
	CreatedAt  int64  `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt  int64  `json:"updatedAt" gorm:"autoUpdateTime:nano"`
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

func (t *Track) savePicture(p *tag.Picture) {
	// TODO
	return
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
	t.savePicture(m.Picture())

	return
}

// AddTrackFromLocation adds a track, given its location
func AddTrackFromLocation(location string, withTags bool) (t *Track, err error) {
	doTag := false
	t = &Track{}
	if err := db.Where("location = ?", location).First(t).Error; err != nil {
		t = &Track{
			Location: location,
		}
		doTag = true
	}

	if withTags || doTag {
		err = t.updateTags()
		onerror.Log(err)
	}

	err = t.Save()
	return
}

// AddTrackFromPath adds a track, given its location
func AddTrackFromPath(path string, withTags bool) (t *Track, err error) {
	var u string
	if u, err = urlstr.PathToURL(path); err != nil {
		return
	}

	t, err = AddTrackFromLocation(u, withTags)
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
