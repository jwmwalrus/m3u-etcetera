package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
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

const (
	// MaxOpenTransientPlaylists -
	MaxOpenTransientPlaylists = 2047
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

// Playlist defines a playlist
type Playlist struct {
	ID              int64         `json:"id" gorm:"primaryKey"`
	Name            string        `json:"name" gorm:"uniqueIndex:unique_idx_playlist_name,not null"`
	Description     string        `json:"description"`
	Open            bool          `json:"open"`
	Active          bool          `json:"active"`
	CreatedAt       int64         `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       int64         `json:"updatedAt" gorm:"autoUpdateTime"`
	PlaylistGroupID int64         `json:"playlistGroupId" gorm:"index:idx_playlist_playlist_group_id,not null"`
	PlaybarID       int64         `json:"playbarId" gorm:"index:idx_playlist_playbar_id,not null"`
	PlaylistGroup   PlaylistGroup `json:"playlistgroup" gorm:"foreignKey:PlaylistGroupID"`
	Playbar         Playbar       `json:"playbar" gorm:"foreignKey:PlaybarID"`
}

// Create implements the DataCreator interface
func (pl *Playlist) Create() (err error) {
	err = db.Create(pl).Error
	return
}

// Delete implements the DataDeleter interface
func (pl *Playlist) Delete() (err error) {
	log.WithField("pl", pl).
		Info("Deleting playlist")

	pl.DeleteDynamicTracks()

	pqys := []PlaylistQuery{}
	if err = db.Where("playlist_id = ?", pl.ID).Find(&pqys).Error; err != nil {
		log.Error(err)
		return
	}

	pts := []PlaylistTrack{}
	if err = db.Where("playlist_id = ?", pl.ID).Find(&pts).Error; err != nil {
		log.Error(err)
		return
	}

	for i := range pqys {
		if err := pqys[i].Delete(); err != nil {
			log.Warn(err)
		}
	}

	for i := range pts {
		onerror.Log(pts[i].Delete())
	}

	err = db.Delete(pl).Error
	return
}

// Read implements the DataReader interface
func (pl *Playlist) Read(id int64) (err error) {
	return db.Joins("Playbar").
		First(pl, id).
		Error
}

// Save implements the DataUpdater interface
func (pl *Playlist) Save() (err error) {
	err = db.Save(pl).Error
	return
}

// ToProtobuf implments ProtoOut interface
func (pl *Playlist) ToProtobuf() proto.Message {
	bv, err := json.Marshal(pl)
	if err != nil {
		log.Error(err)
		return &m3uetcpb.Playlist{}
	}

	out := &m3uetcpb.Playlist{}
	err = json.Unmarshal(bv, out)
	onerror.Log(err)

	// Unmatched
	out.PlaylistGroupId = pl.PlaylistGroupID
	bar := Playbar{}
	bar.Read(pl.PlaybarID)
	if bar.ID == 0 || bar.PerspectiveID == 0 {
		fmt.Println(fmt.Sprintf("SHAMEFUL PLAYLIST: %+v", pl))
	}
	out.Perspective = m3uetcpb.Perspective(bar.getPerspectiveIndex())
	out.Transient = pl.IsTransient()
	out.CreatedAt = pl.CreatedAt
	out.UpdatedAt = pl.UpdatedAt
	return out
}

// AfterCreate is a GORM hook
func (pl *Playlist) AfterCreate(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode {
			subscription.Broadcast(
				subscription.ToPlaybarStoreEvent,
				subscription.Event{
					Idx:  int(PlaybarEventItemAdded),
					Data: pl,
				},
			)
		}
	}()
	return nil
}

// AfterUpdate is a GORM hook
func (pl *Playlist) AfterUpdate(tx *gorm.DB) error {
	go func() {
		if base.FlagTestingMode {
			return
		}
		subscription.Broadcast(
			subscription.ToPlaybarStoreEvent,
			subscription.Event{
				Idx:  int(PlaybarEventItemChanged),
				Data: pl,
			},
		)
		broadcastOpenPlaylist(pl.ID)
	}()
	return nil
}

// AfterDelete is a GORM hook
func (pl *Playlist) AfterDelete(tx *gorm.DB) error {
	go func() {
		if !base.FlagTestingMode {
			subscription.Broadcast(
				subscription.ToPlaybarStoreEvent,
				subscription.Event{
					Idx:  int(PlaybarEventItemRemoved),
					Data: pl,
				},
			)
		}
	}()
	return nil
}

// Count returns the number of tracks in a playlist
func (pl *Playlist) Count() (count int64) {
	err := db.
		Model(&PlaylistTrack{}).
		Where("playlist_id = ?", pl.ID).Count(&count).
		Error

	onerror.Warn(err)
	return
}

// DeleteDelayed deletes a playlist after 5 seconds
func (pl *Playlist) DeleteDelayed() {
	time.Sleep(5 * time.Second)
	onerror.Log(pl.Delete())
}

// DeleteDynamicTracks removes a dynamic track from the database
func (pl *Playlist) DeleteDynamicTracks() {
	pts := []PlaylistTrack{}
	err := db.Where("dynamic = 1 AND playlist_id = ?", pl.ID).
		Find(&pts).
		Error
	if err != nil {
		log.Error(err)
		return
	}

	for i := range pts {
		pts[i].Delete()
	}
}

// GetQueries returns all queries bound by the given playlist
func (pl *Playlist) GetQueries() (pqs []*PlaylistQuery) {
	pqs = []*PlaylistQuery{}
	s := []PlaylistQuery{}
	err := db.Joins("Query").
		Where("playlist_id = ?", pl.ID).
		Find(&s).
		Error
	if err != nil {
		log.Error(err)
		return
	}

	for i := range s {
		pqs = append(pqs, &s[i])
	}
	return
}

// GetTrackAfter returns the next playing track, if any, after the given position.
// Alternatively, return the previous one instead
func (pl *Playlist) GetTrackAfter(curr PlaylistTrack, previous bool) (pt *PlaylistTrack, err error) {
	// Position might have changed, so, reread
	if err = curr.Read(curr.ID); err != nil {
		return
	}

	after := &PlaylistTrack{}

	if previous {
		db.Where("playlist_id = ? AND position < ?", pl.ID, curr.Position).
			Order("position DESC").
			Limit(1).
			Find(after)
	} else {
		db.Where("playlist_id = ? AND position > ?", pl.ID, curr.Position).
			Order("position ASC").
			Limit(1).
			Find(after)
	}

	if after.ID == 0 {
		err = fmt.Errorf("There is no track after")
		return
	}

	pt = &PlaylistTrack{}
	err = db.Joins("Playlist").
		Joins("Track").
		First(pt, after.ID).
		Error
	return
}

// GetTrackAt returns the track at the given position
func (pl *Playlist) GetTrackAt(position int) (pt *PlaylistTrack, err error) {
	at := &PlaylistTrack{}
	err = db.Where("playlist_id = ? AND position = ?", pl.ID, position).
		First(at).
		Error
	if err != nil {
		return
	}

	pt = &PlaylistTrack{}
	err = db.Joins("Playlist").
		Joins("Track").
		First(pt, at.ID).
		Error
	return
}

// GetTracks returns all tracks in the playlist
func (pl *Playlist) GetTracks(limit int) (pts []*PlaylistTrack, ts []*Track) {
	pts = []*PlaylistTrack{}
	ts = []*Track{}

	s := []PlaylistTrack{}
	tx := db.Joins("Track").
		Where("playlist_id = ?", pl.ID).
		Order("position ASC")

	if limit > 0 {
		tx.Limit(limit)
	}

	if err := tx.Find(&s).Error; err != nil {
		log.Error(err)
		return
	}

	for i := range s {
		pts = append(pts, &s[i])
		ts = append(ts, &s[i].Track)
	}
	return
}

// IsTransient checks if playlist is transient
func (pl *Playlist) IsTransient() bool {
	plg := PlaylistGroup{}
	db.First(&plg, pl.PlaylistGroupID)
	return plg.Idx == int(TransientPlaylistGroup)
}

func (pl *Playlist) createTracks(trackIds []int64, locations []string) (pts []PlaylistTrack, err error) {
	for _, id := range trackIds {
		t := Track{}
		if err = t.Read(id); err != nil {
			log.Error(err)
			return
		}
		pts = append(pts, PlaylistTrack{PlaylistID: pl.ID, TrackID: id})
	}

	for _, loc := range locations {
		t := Track{}
		err = db.Where("location = ?", loc).First(&t).Error
		if err != nil {
			t.Location = loc
			if err = t.createTransient(); err != nil {
				return
			}
		}
		pts = append(pts, PlaylistTrack{PlaylistID: pl.ID, TrackID: t.ID})
	}
	return
}

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
	defer DeleteTrackIfTransient(pt.TrackID)

	err = db.Delete(pt).Error
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

// PlaylistQuery Defines a playlist query
type PlaylistQuery struct {
	ID         int64    `json:"id" gorm:"primaryKey"`
	CreatedAt  int64    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  int64    `json:"updatedAt" gorm:"autoUpdateTime"`
	PlaylistID int64    `json:"playlistId" gorm:"index:idx_playlist_query_playlist_id,not null"`
	QueryID    int64    `json:"queryId" gorm:"index:idx_playlist_query_query_id,not null"`
	Playlist   Playlist `json:"playlist" gorm:"foreignKey:PlaylistID"`
	Query      Query    `json:"query" gorm:"foreignKey:QueryID"`
}

// Delete implements the DataDeleter interface
func (pqy *PlaylistQuery) Delete() error {
	return db.Delete(pqy).Error
}

// FindPlaylistsIn returns the tracks for the given IDs
func FindPlaylistsIn(ids []int64) (pls []*Playlist, notFound []int64) {
	pls = []*Playlist{}
	if len(ids) < 1 {
		return
	}

	list := []Playlist{}
	err := db.Where("id in ?", ids).
		Find(&list).
		Error
	if err != nil {
		log.Error(err)
		return
	}

	actual := []int64{}
	for i := range list {
		actual = append(actual, list[i].ID)
		pls = append(pls, &list[i])
	}

	for _, id := range ids {
		if !slice.Contains(actual, id) {
			notFound = append(notFound, id)
		}
	}
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

// GetTransientNameForPlaylist returns the next string
func GetTransientNameForPlaylist() string {
	pls := []Playlist{}
	err := db.Find(&pls).Error
	if err != nil {
		log.Warn(err)
		return "Playlist unknown"
	}
	names := []string{}
	for i := range pls {
		names = append(names, pls[i].Name)
	}

	for i := 1; i <= MaxOpenTransientPlaylists; i++ {
		name := "Playlist " + strconv.Itoa(i)
		if exists, _ := slice.HasValue(names, name); exists {
			continue
		}
		return name
	}

	return ""
}

func broadcastOpenPlaylist(id int64) {
	if base.FlagTestingMode {
		return
	}

	subscription.Broadcast(
		subscription.ToPlaybarStoreEvent,
		subscription.Event{
			Idx:  int(PlaybarEventOpenItems),
			Data: id,
		},
	)
}

func reasignPlaylistTrackPositions(s []PlaylistTrack) []PlaylistTrack {
	pos := 0
	for i := range s {
		pos++
		s[i].Position = pos
	}
	return s
}
