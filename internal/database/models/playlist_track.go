package models

import (
	"fmt"
	"time"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// PlaylistTrack defines a track in a playlist.
type PlaylistTrack struct {
	Model
	Position      int      `json:"position"`
	Dynamic       bool     `json:"dynamic"` // playlist is populated dynamically
	Lastplayedfor int64    `json:"lastplayedfor"`
	PlaylistID    int64    `json:"playlistId" gorm:"index:idx_playlist_track_playlist_id,not null"`
	TrackID       int64    `json:"trackId" gorm:"index:idx_playlist_track_track_id,not null"`
	Playlist      Playlist `json:"playlist" gorm:"foreignKey:PlaylistID"`
	Track         Track    `json:"track" gorm:"foreignKey:TrackID"`
}

func (pt *PlaylistTrack) Create() error {
	return pt.CreateTx(db)
}

func (pt *PlaylistTrack) CreateTx(tx *gorm.DB) error {
	return tx.Create(pt).Error
}

func (pt *PlaylistTrack) Delete() error {
	return pt.DeleteTx(db)
}

func (pt *PlaylistTrack) DeleteTx(tx *gorm.DB) error {
	defer DeleteTrackIfTransient(pt.TrackID)

	return tx.Delete(pt).Error
}

func (pt *PlaylistTrack) Read(id int64) error {
	return pt.ReadTx(db, id)
}

func (pt *PlaylistTrack) ReadTx(tx *gorm.DB, id int64) error {
	return tx.Joins("Playlist").
		Joins("Track").
		First(pt, id).
		Error
}

func (pt *PlaylistTrack) Save() error {
	return pt.SaveTx(db)
}

func (pt *PlaylistTrack) SaveTx(tx *gorm.DB) error {
	return tx.Save(pt).Error
}

func (pt *PlaylistTrack) ToProtobuf() proto.Message {
	return &m3uetcpb.PlaylistTrack{
		Id:            pt.ID,
		Position:      int32(pt.Position),
		Dynamic:       pt.Dynamic,
		Lastplayedfor: pt.Lastplayedfor,
		PlaylistId:    pt.PlaylistID,
		TrackId:       pt.TrackID,
		CreatedAt:     timestamppb.New(time.Unix(0, pt.CreatedAt)),
		UpdatedAt:     timestamppb.New(time.Unix(0, pt.UpdatedAt)),
	}
}

// AfterCreate is a GORM hook.
func (pt *PlaylistTrack) AfterCreate(tx *gorm.DB) error {
	go broadcastOpenPlaylist(pt.PlaylistID)
	return nil
}

// AfterUpdate is a GORM hook.
func (pt *PlaylistTrack) AfterUpdate(tx *gorm.DB) error {
	go broadcastOpenPlaylist(pt.PlaylistID)
	return nil
}

// AfterDelete is a GORM hook.
func (pt *PlaylistTrack) AfterDelete(tx *gorm.DB) error {
	go broadcastOpenPlaylist(pt.PlaylistID)
	return nil
}

// GetPosition implements the Poser interface.
func (pt *PlaylistTrack) GetPosition() int {
	return pt.Position
}

// SetPosition implements the Poser interface.
func (pt *PlaylistTrack) SetPosition(pos int) {
	pt.Position = pos
}

// GetIgnore implements the Poser interface.
func (pt *PlaylistTrack) GetIgnore() bool {
	return false
}

// SetIgnore implements the Poser interface.
func (pt *PlaylistTrack) SetIgnore(_ bool) {}

// GetTrackAfter returns the track after the current one.
func (pt *PlaylistTrack) GetTrackAfter(goingBack bool) (*PlaylistTrack, error) {
	pl := pt.Playlist
	if pl.ID == 0 {
		return nil, fmt.Errorf("There is no list to play from")
	}

	return pl.GetTrackAfter(*pt, goingBack)
}

// GetActivePlaylistTrack deletes a playlist.
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
