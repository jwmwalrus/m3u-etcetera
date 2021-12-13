package models

import (
	"encoding/json"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"

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

func (t *Track) Create() (err error) {
	err = db.Create(t).Error
	return
}

func (t *Track) Delete() {
	err := db.Delete(&Track{}, t.ID).Error
	onerror.Log(err)
}

func (t *Track) FindBy(query interface{}) (err error) {
	err = db.Where(query).First(t).Error
	return
}

func (t *Track) Read(id int64) (err error) {
	err = db.First(t, id).Error
	return
}

func (t *Track) Save() (err error) {
	err = db.Save(t).Error
	return
}

func (pb *Track) ToProtobuf() *m3uetcpb.Track {
	bv, err := json.Marshal(pb)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.Track{}
	}

	out := &m3uetcpb.Track{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)
	return out
}
