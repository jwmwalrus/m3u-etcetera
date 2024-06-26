package models

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/jwmwalrus/bnp/chars"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/pointers"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/impexp"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	rtc "github.com/jwmwalrus/rtcycler"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	// MaxOpenTransientPlaylists -.
	MaxOpenTransientPlaylists = 2047
)

// Playlist defines a playlist.
type Playlist struct {
	Model
	Name            string        `json:"name" gorm:"uniqueIndex:unique_idx_playlist_name,not null"`
	Description     string        `json:"description"`
	Open            bool          `json:"open"`
	Active          bool          `json:"active"`
	Transient       bool          `json:"transient"`
	Bucket          bool          `json:"bucket"`
	QueryID         int64         `json:"queryId"`
	PlaylistGroupID int64         `json:"playlistGroupId" gorm:"index:idx_playlist_playlist_group_id,not null"`
	PlaybarID       int64         `json:"playbarId" gorm:"index:idx_playlist_playbar_id,not null"`
	PlaylistGroup   PlaylistGroup `json:"playlistGroup" gorm:"foreignKey:PlaylistGroupID"`
	Playbar         Playbar       `json:"playbar" gorm:"foreignKey:PlaybarID"`
}

func (pl *Playlist) Create() (err error) {
	return pl.CreateTx(db)
}

func (pl *Playlist) CreateTx(tx *gorm.DB) (err error) {
	err = tx.Create(pl).Error
	return
}

func (pl *Playlist) Delete() (err error) {
	return pl.DeleteTx(db)
}

func (pl *Playlist) DeleteTx(tx *gorm.DB) (err error) {
	logw := slog.With("pl", *pl)
	logw.Info("Deleting playlist")

	pl.DeleteDynamicTracks(tx)

	pqys := []PlaylistQuery{}
	if err = tx.Where("playlist_id = ?", pl.ID).Find(&pqys).Error; err != nil {
		logw.Error("Failed to find playlist queries in database", "error", err)
		return
	}

	pts := []PlaylistTrack{}
	if err = tx.Where("playlist_id = ?", pl.ID).Find(&pts).Error; err != nil {
		logw.Error("Failed to find playlist tracks in database", "error", err)
		return
	}

	for i := range pqys {
		if err := pqys[i].DeleteTx(tx); err != nil {
			logw.With(
				"pqy", pqys[i],
				"error", err,
			).Warn("Failed to delete playlist query")
		}
	}

	for i := range pts {
		onerror.NewRecorder(logw).Log(pts[i].DeleteTx(tx))
	}

	pl.Transient = true
	rl, _ := chars.GetRandomLetters(8)
	pl.Name += fmt.Sprintf(" (deleted %v)", rl)
	err = tx.Save(pl).Error
	return
}

func (pl *Playlist) Read(id int64) error {
	return pl.ReadTx(db, id)
}

func (pl *Playlist) ReadTx(tx *gorm.DB, id int64) error {
	return tx.Joins("Playbar").
		First(pl, id).
		Error
}

func (pl *Playlist) Save() error {
	return db.Save(pl).Error
}

func (pl *Playlist) SaveTx(tx *gorm.DB) error {
	return tx.Save(pl).Error
}

func (pl *Playlist) ToProtobuf() proto.Message {
	var bar Playbar
	bar.Read(pl.PlaybarID)

	return &m3uetcpb.Playlist{
		Id:              pl.ID,
		Name:            pl.Name,
		Description:     pl.Description,
		Open:            pl.Open,
		Active:          pl.Active,
		Transient:       pl.Transient,
		Bucket:          pl.Bucket,
		QueryId:         pl.QueryID,
		PlaylistGroupId: pl.PlaylistGroupID,
		Duration:        pl.Duration(), Perspective: m3uetcpb.Perspective(bar.getPerspectiveIndex()),
		CreatedAt: timestamppb.New(time.Unix(0, pl.CreatedAt)),
		UpdatedAt: timestamppb.New(time.Unix(0, pl.UpdatedAt)),
	}
}

// AfterCreate is a GORM hook.
func (pl *Playlist) AfterCreate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToPlaybarStoreEvent,
			subscription.Event{
				Idx:  int(PlaybarEventItemAdded),
				Data: pl,
			},
		)
	}()
	return nil
}

// AfterUpdate is a GORM hook.
func (pl *Playlist) AfterUpdate(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
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

// AfterDelete is a GORM hook.
func (pl *Playlist) AfterDelete(tx *gorm.DB) error {
	go func() {
		if rtc.FlagTestMode() {
			return
		}
		subscription.Broadcast(
			subscription.ToPlaybarStoreEvent,
			subscription.Event{
				Idx:  int(PlaybarEventItemRemoved),
				Data: pl,
			},
		)
	}()
	return nil
}

// Count returns the number of tracks in a playlist.
func (pl *Playlist) Count() (count int64) {
	err := db.
		Model(&PlaylistTrack{}).
		Where("playlist_id = ?", pl.ID).Count(&count).
		Error

	onerror.Warn(err)
	return
}

// DeleteDynamicTracks removes a dynamic track from the database.
func (pl *Playlist) DeleteDynamicTracks(tx *gorm.DB) {
	pts := []PlaylistTrack{}
	err := tx.Where("dynamic = 1 AND playlist_id = ?", pl.ID).
		Find(&pts).
		Error
	if err != nil {
		slog.Error("Failed to find dynamic playlist tracks in database", "error", err)
		return
	}

	for i := range pts {
		pts[i].DeleteTx(tx)
	}
}

// Duration returns the duration of the playlist.
func (pl *Playlist) Duration() int64 {
	var d sql.NullInt64
	err := db.Raw("SELECT sum(t.duration) FROM track t JOIN playlist_track pt ON pt.track_id = t.id WHERE pt.playlist_id = ?", pl.ID).
		Row().
		Scan(&d)
	onerror.Log(err)
	return d.Int64
}

// Export exports a playlist with the given format to the given location.
func (pl *Playlist) Export(format impexp.PlaylistType, location string) (err error) {
	path, err := urlstr.URLToPath(location)
	if err != nil {
		return
	}

	var f *os.File
	if f, err = os.Create(path); err != nil {
		return
	}
	defer f.Close()

	props := impexp.PlaylistProps{impexp.NamePropKey: pl.Name}
	m3u, err := impexp.New(format, props)
	if err != nil {
		return
	}

	pts, _ := pl.GetTracks(0)

	tis := []impexp.TrackInfo{}
	for _, pt := range pts {
		var un string
		if un, err = urlstr.URLToPath(pt.Track.Location); err != nil {
			return
		}
		t := impexp.TrackInfo{
			Location:    un,
			Title:       pt.Track.Title,
			ArtistTitle: pt.Track.Artist + " - " + pt.Track.Title,
			Album:       pt.Track.Album,
			Artist:      pt.Track.Artist,
			Albumartist: pt.Track.Albumartist,
			Genre:       pt.Track.Genre,
			Duration:    pt.Track.Duration,
			Year:        pt.Track.Year,
		}
		tis = append(tis, t)
	}

	m3u.Add(tis)

	if _, err = m3u.Format(f); err != nil {
		return
	}
	return
}

// GetQueries returns all queries bound by the given playlist.
func (pl *Playlist) GetQueries() []*PlaylistQuery {
	pqs := []PlaylistQuery{}
	err := db.Joins("Query").
		Where("playlist_id = ?", pl.ID).
		Find(&pqs).
		Error
	if err != nil {
		slog.Error("Failed to find playlist queries in database", "error", err)
		return []*PlaylistQuery{}
	}

	return pointers.FromSlice(pqs)
}

// GetTrackAfter returns the next playing track, if any, after the given position.
// Alternatively, return the previous one instead.
func (pl *Playlist) GetTrackAfter(curr PlaylistTrack,
	previous bool) (pt *PlaylistTrack, err error) {

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

// GetTrackAt returns the track at the given position.
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

// GetTracks returns all tracks in the playlist.
func (pl *Playlist) GetTracks(limit int) ([]*PlaylistTrack, []*Track) {
	pts := []PlaylistTrack{}

	tx := db.Joins("Track").
		Where("playlist_id = ?", pl.ID).
		Order("position ASC")

	if limit > 0 {
		tx.Limit(limit)
	}

	if err := tx.Find(&pts).Error; err != nil {
		slog.Error("Failed to find playlist tracks in database", "error", err)
		return []*PlaylistTrack{}, []*Track{}
	}

	ts := []Track{}
	for i := range pts {
		ts = append(ts, pts[i].Track)
	}
	return pointers.FromSlice(pts), pointers.FromSlice(ts)
}

func (pl *Playlist) createTracks(trackIds []int64,
	locations []string) (pts []PlaylistTrack, err error) {

	for _, id := range trackIds {
		t := Track{}
		if err = t.Read(id); err != nil {
			slog.With(
				"id", id,
				"error", err,
			).Error("Failed to read track")
			return
		}
		pts = append(pts, PlaylistTrack{PlaylistID: pl.ID, TrackID: id})
	}

	tx := db.Session(&gorm.Session{SkipHooks: true})
	for _, loc := range locations {
		t := Track{}
		err = tx.Where("location = ?", loc).First(&t).Error
		if err != nil {
			t.Location = loc
			if err = t.createTransient(tx, nil); err != nil {
				return
			}
		}
		pts = append(pts, PlaylistTrack{PlaylistID: pl.ID, TrackID: t.ID})
	}
	return
}

// FindPlaylistsIn returns the playlists for the given IDs.
func FindPlaylistsIn(ids []int64) (pls []*Playlist, notFound []int64) {
	pls = []*Playlist{}
	if len(ids) < 1 {
		return
	}

	s := []Playlist{}
	err := db.Where("id in ?", ids).
		Find(&s).
		Error
	if err != nil {
		slog.With(
			"ids", ids,
			"error", err,
		).Error("Failed to find playlists in database")
		return
	}

	actual := []int64{}
	for i := range s {
		actual = append(actual, s[i].ID)
		pls = append(pls, &s[i])
	}

	for _, id := range ids {
		if !slices.Contains(actual, id) {
			notFound = append(notFound, id)
		}
	}
	return
}

// GetTransientNameForPlaylist returns the next string.
func GetTransientNameForPlaylist(queryID int64) string {
	pls := []Playlist{}
	err := db.Find(&pls).Error
	if err != nil {
		slog.Warn("Failed to find all playlists in database", "error", err)
		return "Playlist unknown"
	}
	names := []string{}
	for i := range pls {
		names = append(names, pls[i].Name)
	}

	prefix := "Playlist "
	if queryID > 0 {
		qy := Query{}
		err = qy.Read(queryID)
		onerror.Log(err)
		if err == nil {
			descr := QueryIndex(qy.Idx).Description()
			if descr != "" {
				prefix = descr + " "
			}
		}
	}

	for i := 1; i <= MaxOpenTransientPlaylists; i++ {
		name := prefix + strconv.Itoa(i)
		if slices.Contains(names, name) {
			continue
		}
		return name
	}

	return ""
}

func broadcastOpenPlaylist(id int64) {
	if rtc.FlagTestMode() {
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
