package models

import (
	"encoding/json"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// PlaylistTrack defines a track in a playlist
type PlaylistTrack struct {
	ID         int64    `json:"id" gorm:"primaryKey"`
	Position   int      `json:"position"`
	Dynamic    bool     `json:"dynamic"` // playlist is populated dynamically
	CreatedAt  int64    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  int64    `json:"updatedAt" gorm:"autoUpdateTime"`
	PlaylistID int64    `json:"playlistId" gorm:"index:idx_playlist_track_playlist_id,not null"`
	TrackID    int64    `json:"trackId" gorm:"index:idx_playlist_track_track_id,not null"`
	Playlist   Playlist `json:"playlist" gorm:"foreignKey:PlaylistID"`
	Track      Track    `json:"track" gorm:"foreignKey:TrackID"`
}

// Create implements the DataCreator interface
func (pt *PlaylistTrack) Create() error {
	return db.Create(pt).Error
}

// Delete implements the DataDeleter interface
func (pt *PlaylistTrack) Delete() (err error) {
	return pt.DeleteTx(db)
}

// DeleteTx implements the DataDeleterTx interface
func (pt *PlaylistTrack) DeleteTx(tx *gorm.DB) (err error) {
	defer DeleteTrackIfTransient(pt.TrackID)

	err = tx.Delete(pt).Error
	return
}

// Read implements the DataReader interface
func (pt *PlaylistTrack) Read(id int64) (err error) {
	return db.Joins("Playlist").
		Joins("Track").
		First(pt, id).
		Error
}

// ToProtobuf implments ProtoOut interface
func (pt *PlaylistTrack) ToProtobuf() proto.Message {
	bv, err := json.Marshal(pt)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.PlaylistTrack{}
	}

	out := &m3uetcpb.PlaylistTrack{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)

	// Unmatched
	out.PlaylistId = pt.PlaylistID
	out.TrackId = pt.TrackID
	out.CreatedAt = pt.CreatedAt
	out.UpdatedAt = pt.UpdatedAt
	return out
}

// AfterCreate is a GORM hook
func (pt *PlaylistTrack) AfterCreate(tx *gorm.DB) error {
	go broadcastOpenPlaylist(pt.PlaylistID)
	return nil
}

// AfterUpdate is a GORM hook
func (pt *PlaylistTrack) AfterUpdate(tx *gorm.DB) error {
	go broadcastOpenPlaylist(pt.PlaylistID)
	return nil
}

// AfterDelete is a GORM hook
func (pt *PlaylistTrack) AfterDelete(tx *gorm.DB) error {
	go broadcastOpenPlaylist(pt.PlaylistID)
	return nil
}

// GetPosition implements the Poser interface
func (pt PlaylistTrack) GetPosition() int {
	return pt.Position
}

// SetPosition implements the Poser interface
func (pt PlaylistTrack) SetPosition(pos int) {
	pt.Position = pos
}

// GetIgnore implements the Poser interface
func (pt PlaylistTrack) GetIgnore() bool {
	return false
}

// SetIgnore implements the Poser interface
func (pt PlaylistTrack) SetIgnore(_ bool) {
	return
}

// GetActivePlaylistTrack deletes a playlist
func GetActivePlaylistTrack() (pt *PlaylistTrack) {
	pb := Playback{}
	err := pb.GetNextToPlay()
	if err != nil || pb.ID == 0 || pb.TrackID == 0 {
		return
	}

	active := &PlaylistTrack{}
	err = db.
		Joins("JOIN playlist ON playlist_track.playlist_id = playlist.id").
		Where("playlist.active = 1 AND playlist_track.track_id = ?", pb.TrackID).
		First(active).
		Error
	if err != nil {
		return
	}

	pt = &PlaylistTrack{}
	err = db.Joins("Playlist").
		Joins("Track").
		First(pt, active.ID).
		Error
	onerror.Log(err)
	return
}
