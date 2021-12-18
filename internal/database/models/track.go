package models

import (
	"encoding/json"
	"os"

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
	Year        int    `json:"year" gorm:"index:idx_track_year"`
	Tracknumber int    `json:"trackNumber"`
	Tracktotal  int    `json:"trackTotal"`
	Discnumber  int    `json:"discNumber"`
	Disctotal   int    `json:"discTotal"`
	Lyrics      string `json:"lyrics"`
	Comment     string `json:"comment"`
	Tags        string `json:"tags"`
	Sum         string `json:"sum" gorm:"index:idx_track_sum"`
	Playcount   int    `json:"playCount"`
	Rating      int    `json:"rating" gorm:"index:idx_track_rating"`
	Duration    int64  `json:"duration"`
	Remote      bool   `json:"remote"` // if track is remote
	Lastplayed  int64  `json:"lastPlayed"`
	CreatedAt   int64  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   int64  `json:"updatedAt" gorm:"autoUpdateTime"`
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
	out.CreatedAt = t.CreatedAt
	out.UpdatedAt = t.UpdatedAt
	return out
}

// FindBy implements the DataFinder interface
func (t *Track) FindBy(query interface{}) (err error) {
	err = db.Where(query).First(t).Error
	return
}

func (t *Track) updateTags() (err error) {
	base.GetBusy(base.IdleStatusFileOperations)
	defer base.GetFree(base.IdleStatusFileOperations)

	var path string
	path, err = urlstr.URLToPath(t.Location)

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
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
	t.Year = m.Year()
	t.Tracknumber, t.Tracktotal = m.Track()
	t.Discnumber, t.Disctotal = m.Disc()
	t.Lyrics = m.Lyrics()
	t.Comment = m.Comment()

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

	err = db.Save(t).Error
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

// CheckNotFoundTracks returns the not-found track IDs from a given list
func CheckNotFoundTracks(ids []int64) (notFound []int64) {
	ts := FindTracksIn(ids)
	actual := []int64{}
	for i := range ts {
		actual = append(actual, ts[i].ID)
	}

	for _, id := range ids {
		if !slice.Contains(actual, id) {
			notFound = append(notFound, id)
		}
	}
	return
}

// FindTracksIn returns the tracks for the given IDs
func FindTracksIn(ids []int64) (ts []*Track) {
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

	ts = []*Track{}
	for i := range list {
		ts = append(ts, &list[i])
	}
	return
}

func appendToTrackList(list []*Track, ts []*Track) {
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
}
