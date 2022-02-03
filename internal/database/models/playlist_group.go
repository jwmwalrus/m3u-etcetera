package models

import (
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// PlaylistGroupIndex -
type PlaylistGroupIndex int

// Defines the default playlist groups
const (
	DefaultPlaylistGroup PlaylistGroupIndex = iota + 1
	TransientPlaylistGroup
)

func (idx PlaylistGroupIndex) String() string {
	return [...]string{"", "\t", "\t\t"}[idx]
}

// Get returns the playlist group for the given index
func (idx PlaylistGroupIndex) Get() (plg *PlaylistGroup, err error) {
	plg = &PlaylistGroup{}
	err = db.Where("idx = ?", int(idx)).First(plg).Error
	return
}

// PlaylistGroup defines a playlist group
type PlaylistGroup struct {
	ID            int64       `json:"id" gorm:"primaryKey"`
	Idx           int         `json:"idx" gorm:"not null"`
	Name          string      `json:"name" gorm:"uniqueIndex:unique_idx_playlist_group_name,not null"`
	Description   string      `json:"description"`
	Hidden        bool        `json:"hidden"`
	CreatedAt     int64       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     int64       `json:"updatedAt" gorm:"autoUpdateTime"`
	PerspectiveID int64       `json:"perspectiveId" gorm:"index:idx_playlist_group_perspective_id,not null"`
	Perspective   Perspective `json:"perspective" gorm:"foreignKey:PerspectiveID"`
}

// Create implements the DataCreator interface
func (pg *PlaylistGroup) Create() (err error) {
	err = db.Create(pg).Error
	return
}

// Delete implements the DataDeleter interface
func (pg *PlaylistGroup) Delete() error {
	return db.Delete(pg).Error
}

// Read implements the DataReader interface
func (pg *PlaylistGroup) Read(id int64) (err error) {
	return db.Joins("Perspective").
		First(pg, id).
		Error
}

// Save implements the DataUpdater interface
func (pg *PlaylistGroup) Save() error {
	return db.Save(pg).Error
}

// ToProtobuf implments ProtoOut interface
func (pg *PlaylistGroup) ToProtobuf() proto.Message {
	out := &m3uetcpb.PlaylistGroup{}

	out.Id = pg.ID
	out.Name = pg.Name
	out.Description = pg.Description
	out.Perspective = m3uetcpb.Perspective(pg.Perspective.Idx)
	out.CreatedAt = pg.CreatedAt
	out.UpdatedAt = pg.UpdatedAt
	return out
}

// AfterCreate is a GORM hook
func (pg *PlaylistGroup) AfterCreate(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode {
			subscription.Broadcast(
				subscription.ToPlaybarStoreEvent,
				subscription.Event{
					Idx:  int(PlaybarEventItemAdded),
					Data: pg,
				},
			)
		}
	}()
	return nil
}

// AfterUpdate is a GORM hook
func (pg *PlaylistGroup) AfterUpdate(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode {
			subscription.Broadcast(
				subscription.ToPlaybarStoreEvent,
				subscription.Event{
					Idx:  int(PlaybarEventItemChanged),
					Data: pg,
				},
			)
		}
	}()
	return nil
}

// AfterDelete is a GORM hook
func (pg *PlaylistGroup) AfterDelete(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode {
			subscription.Broadcast(
				subscription.ToPlaybarStoreEvent,
				subscription.Event{
					Idx:  int(PlaybarEventItemRemoved),
					Data: pg,
				},
			)
		}
	}()
	return nil
}