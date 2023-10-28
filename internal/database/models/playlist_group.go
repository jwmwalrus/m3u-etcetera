package models

import (
	"encoding/json"
	"log/slog"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	rtc "github.com/jwmwalrus/rtcycler"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// PlaylistGroupIndex -.
type PlaylistGroupIndex int

// Defines the default playlist groups.
const (
	MusicPlaylistGroup PlaylistGroupIndex = iota + 1
	RadioPlaylistGroup
	PodcastsPlaylistGroup
	AudiobooksPlaylistGroup
)

func (idx PlaylistGroupIndex) String() string {
	return [...]string{"", "\t", "\t\t", "\t\t\t", "\t\t\t\t"}[idx]
}

// Get returns the playlist group for the given index.
func (idx PlaylistGroupIndex) Get() (plg *PlaylistGroup, err error) {
	plg = &PlaylistGroup{}
	err = db.Where("idx = ?", int(idx)).First(plg).Error
	return
}

// PlaylistGroup defines a playlist group.
type PlaylistGroup struct {
	ID            int64       `json:"id" gorm:"primaryKey"`
	Idx           int         `json:"idx" gorm:"not null"`
	Name          string      `json:"name" gorm:"uniqueIndex:unique_idx_playlist_group_name,not null"`
	Description   string      `json:"description"`
	Hidden        bool        `json:"hidden"`
	CreatedAt     int64       `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt     int64       `json:"updatedAt" gorm:"autoUpdateTime:nano"`
	PerspectiveID int64       `json:"perspectiveId" gorm:"index:idx_playlist_group_perspective_id,not null"`
	Perspective   Perspective `json:"-" gorm:"foreignKey:PerspectiveID"`
}

func (pg *PlaylistGroup) Create() error {
	return pg.CreateTx(db)
}

func (pg *PlaylistGroup) CreateTx(tx *gorm.DB) error {
	return tx.Create(pg).Error
}

func (pg *PlaylistGroup) Delete() error {
	return pg.DeleteTx(db)
}

func (pg *PlaylistGroup) DeleteTx(tx *gorm.DB) error {
	pls := []Playlist{}
	err := tx.Where("playlist_group_id = ?", pg.ID).Find(&pls).Error
	if err != nil {
		return err
	}
	if len(pls) > 0 {
		pgd := PlaylistGroup{}
		err = tx.Where("hidden = 1 and idx > 0 and perspective_id = ?", pg.PerspectiveID).Find(&pgd).Error
		if err != nil {
			return err
		}
		for i := range pls {
			pls[i].PlaylistGroupID = pgd.ID
		}
		err = tx.Save(&pls).Error
		if err != nil {
			return err
		}
	}

	return tx.Delete(pg).Error
}

func (pg *PlaylistGroup) Read(id int64) error {
	return pg.ReadTx(db, id)
}

func (pg *PlaylistGroup) ReadTx(tx *gorm.DB, id int64) error {
	return tx.Joins("Perspective").
		First(pg, id).
		Error
}

func (pg *PlaylistGroup) Save() error {
	return pg.SaveTx(db)
}

func (pg *PlaylistGroup) SaveTx(tx *gorm.DB) error {
	return tx.Save(pg).Error
}

func (pg *PlaylistGroup) ToProtobuf() proto.Message {
	bv, err := json.Marshal(pg)
	if err != nil {
		slog.Error("Failed to marshal playlist group", "error", err)
		return &m3uetcpb.PlaylistGroup{}
	}

	out := &m3uetcpb.PlaylistGroup{}
	err = jsonUnmarshaler.Unmarshal(bv, out)
	onerror.Log(err)

	// Unmatched
	out.Perspective = m3uetcpb.Perspective(pg.Perspective.Idx)

	return out
}

// AfterCreate is a GORM hook.
func (pg *PlaylistGroup) AfterCreate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToPlaybarStoreEvent,
			subscription.Event{
				Idx:  int(PlaybarEventItemAdded),
				Data: pg,
			},
		)
	}()
	return nil
}

// AfterUpdate is a GORM hook.
func (pg *PlaylistGroup) AfterUpdate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToPlaybarStoreEvent,
			subscription.Event{
				Idx:  int(PlaybarEventItemChanged),
				Data: pg,
			},
		)
	}()
	return nil
}

// AfterDelete is a GORM hook.
func (pg *PlaylistGroup) AfterDelete(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToPlaybarStoreEvent,
			subscription.Event{
				Idx:  int(PlaybarEventItemRemoved),
				Data: pg,
			},
		)
	}()
	return nil
}

// ReadDefaultForPerspective returns the default playlist group for
// the given perspective.
func (pg *PlaylistGroup) ReadDefaultForPerspective(id int64) error {
	return db.Joins("Perspective").
		Where("perspective_id = ? and playlist_group.idx > 0", id).
		First(pg, id).
		Error
}
