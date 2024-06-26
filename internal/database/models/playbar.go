package models

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jwmwalrus/bnp/chars"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/pointers"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/impexp"
	"github.com/jwmwalrus/m3u-etcetera/pkg/poser"
	"gorm.io/gorm"
)

// GetPlaybar returns the playbar associated to the given perspective.
func (idx PerspectiveIndex) GetPlaybar() (bar *Playbar, err error) {
	bar = &Playbar{}
	err = db.
		Joins(
			"JOIN perspective ON playbar.perspective_id = perspective.id AND perspective.idx = ?",
			int(idx),
		).
		First(bar).
		Error
	return

}

// PlaybarEvent defines a collection event.
type PlaybarEvent int

// PlaybarEvent enum.
const (
	PlaybarEventNone PlaybarEvent = iota
	PlaybarEventInitial
	_
	_
	PlaybarEventItemAdded
	PlaybarEventItemChanged
	PlaybarEventItemRemoved
	PlaybarEventOpenItems
	_
	_
)

func (ce PlaybarEvent) String() string {
	return []string{
		"none",
		"initial",
		"initial-item",
		"initial-done",
		"item-added",
		"item-changed",
		"item-removed",
		"open-items",
		"open-items-item",
		"open-items-done",
	}[ce]
}

// Playbar defines the playlist bar for each perspective.
type Playbar struct {
	Model
	PerspectiveID int64       `json:"perspectiveId" gorm:"uniqueIndex:unique_idx_playbar_perspective_id,not null"`
	Perspective   Perspective `json:"perspective" gorm:"foreignKey:PerspectiveID"`
}

func (b *Playbar) Read(id int64) error {
	return b.ReadTx(db, id)
}

func (b *Playbar) ReadTx(tx *gorm.DB, id int64) error {
	return tx.Joins("Perspective").First(b, id).Error
}

// ActivateEntry activates the given entry in a playbar.
func (b *Playbar) ActivateEntry(pl *Playlist) {
	logw := slog.With("pl", *pl)
	logw.Info("Activating in playbar")

	pl.Open = true
	if err := pl.Save(); err != nil {
		logw.Error("Failed to save playlist", "error", err)
		return
	}

	pls := []Playlist{}
	err := db.Where("playbar_id = ? and open=1", b.ID).Find(&pls).Error
	if err != nil {
		logw.Error("Failed to find playlists in database", "error", err)
		return
	}

	for i := range pls {
		if pls[i].ID == pl.ID {
			pls[i].Active = true
			continue
		}
		pls[i].Active = false
	}
	onerror.NewRecorder(logw).Log(db.Save(&pls).Error)
}

// AppendToPlaylist -.
func (b *Playbar) AppendToPlaylist(pl *Playlist, trackIds []int64,
	locations []string) {

	logw := slog.With(
		"pl", *pl,
		"track_ids", trackIds,
		"locations", locations,
	)
	logw.Info("Appending tracks/locations to playlist")

	tx := db.Session(&gorm.Session{SkipHooks: true})

	pts := []PlaylistTrack{}
	err := tx.Where("playlist_id = ?", pl.ID).Order("position ASC").
		Find(&pts).
		Error
	if err != nil {
		logw.Error("Failed to find playlist tracks in database", "error", err)
		return
	}

	s, err := pl.createTracks(trackIds, locations)
	if err != nil {
		logw.Error("Failed to create tracks", "error", err)
		return
	}

	list := poser.AppendTo(pointers.FromSlice(pts), pointers.FromSlice(s)...)
	pts = pointers.ToValues(list)
	err = tx.Save(&pts).Error
	onerror.NewRecorder(logw).Log(err)

	if pl.Open {
		broadcastOpenPlaylist(pl.ID)
	}
}

// ClearPlaylist -.
func (b *Playbar) ClearPlaylist(pl *Playlist) {
	logw := slog.With("pl", *pl)
	logw.Info("Clearing tracks/locations in playlist")

	pts := []PlaylistTrack{}
	if err := db.Where("playlist_id = ?", pl.ID).Find(&pts).Error; err != nil {
		logw.Error("Failed to find playlist tracks in database", "error", err)
		return
	}

	if len(pts) > 0 {
		err := db.Session(&gorm.Session{SkipHooks: true}).
			Where("id > 0").
			Delete(&pts).
			Error
		onerror.NewRecorder(logw).Log(err)

		broadcastOpenPlaylist(pl.ID)
	}
}

// CloseEntry closes the given playbar entry.
func (b *Playbar) CloseEntry(pl *Playlist) {
	if !pl.Transient {
		pl.DeleteDynamicTracks(db)
	}

	pl.Open = false
	pl.Active = false
	onerror.Log(pl.Save())
}

// CreateEntry creates a playlist.
func (b *Playbar) CreateEntry(name, description string, queryID int64) (pl *Playlist, err error) {
	pg := &PlaylistGroup{}
	err = pg.ReadDefaultForPerspective(b.PerspectiveID)
	if err != nil {
		return
	}

	var plname string

	isTransient := false
	if name != "" {
		plname = name
	} else {
		isTransient = true
		plname = GetTransientNameForPlaylist(queryID)
	}

	pl = &Playlist{
		Name:            plname,
		Description:     description,
		PlaybarID:       b.ID,
		PlaylistGroupID: pg.ID,
		Transient:       isTransient,
		QueryID:         queryID,
	}
	err = pl.Create()
	if err != nil {
		return
	}

	if isTransient {
		pl.Open = true
		err = pl.Save()
	}
	return
}

// CreateGroup creates a playlist.
func (b *Playbar) CreateGroup(name, description string) (pg *PlaylistGroup, err error) {
	pg = &PlaylistGroup{
		Name:          name,
		Description:   description,
		PerspectiveID: b.PerspectiveID,
	}
	err = pg.Create()
	return
}

// DeactivateEntry -.
func (b *Playbar) DeactivateEntry(pl *Playlist) {
	pl.Active = false
	onerror.Log(pl.Save())
}

// DeleteFromPlaylist -.
func (b *Playbar) DeleteFromPlaylist(pl *Playlist, position int) {
	logw := slog.With(
		"pl", *pl,
		"position", position,
	)
	logw.Info("Deleting position in playlist")

	pts := []PlaylistTrack{}
	err := db.Where("playlist_id = ?", pl.ID).Order("position ASC").
		Find(&pts).
		Error
	if err != nil {
		logw.Error("Failed to find playlist tracks in database", "error", err)
		return
	}
	list, pt := poser.DeleteAt(pointers.FromSlice(pts), position)
	pts = pointers.ToValues(list)

	if pt != nil && pt.ID > 0 {
		if err := pt.Delete(); err != nil {
			logw.Error("Failed to delete playlist track", "error", err)
			return
		}
	}
	onerror.NewRecorder(logw).Log(db.Save(&pts).Error)
}

// DestroyEntry deletes a playlist.
func (b *Playbar) DestroyEntry(pl *Playlist) error {
	return pl.Delete()
}

// DestroyGroup deletes a playlist.
func (b *Playbar) DestroyGroup(pg *PlaylistGroup) error {
	return pg.Delete()
}

// GetAllEntries -.
func (b *Playbar) GetAllEntries(limit int) []*Playlist {
	pls := []Playlist{}

	tx := db.Joins("Playbar").
		Where("playbar_id = ?", b.ID)
	if limit > 0 {
		tx.Limit(limit)
	}
	err := tx.Find(&pls).Error
	if err != nil {
		return []*Playlist{}
	}

	return pointers.FromSlice(pls)
}

// GetAllGroups -.
func (b *Playbar) GetAllGroups(limit int) []*PlaylistGroup {
	pgs := []PlaylistGroup{}

	tx := db.Joins("Perspective").
		Where("hidden = 0 and perspective_id = ?", b.PerspectiveID)
	if limit > 0 {
		tx.Limit(limit)
	}
	err := tx.Find(&pgs).Error
	if err != nil {
		slog.Error("Failed to find all playlist groups in database", "error", err)
		return []*PlaylistGroup{}
	}

	return pointers.FromSlice(pgs)
}

// GetAllOpenEntries -.
func (b *Playbar) GetAllOpenEntries() []*Playlist {
	pls := []Playlist{}
	err := db.Joins("Playbar").
		Where("open = 1 and playbar_id = ?", b.ID).
		Find(&pls).
		Error
	if err != nil {
		return []*Playlist{}
	}

	return pointers.FromSlice(pls)
}

// ImportPlaylist creates a playlist from the given location, if supported.
func (b *Playbar) ImportPlaylist(location string, asTransient bool) (pl *Playlist, msgs []string, err error) {
	def, err := importPlaylist(location)
	if err != nil {
		return
	}

	pl, msgs, err = b.newPlaylistFromImported(def)
	if err != nil {
		return
	}

	name := def.Name()
	description := ""

	if name == "" {
		var path string
		path, err = urlstr.URLToPath(location)
		if err != nil {
			return
		}

		if !asTransient {
			rl, _ := chars.GetRandomLetters(8)
			name = strings.Join([]string{
				strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
				rl,
				def.Type(),
			}, "-")
			description = "Imported from " + filepath.Base(path) +
				" on " + time.Now().Format(time.RFC3339)
		} else {
			name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		}
	}

	err = db.Where("name = ?", name).First(&Playlist{}).Error
	if err == nil {
		var found bool
		for i := 1; i <= MaxOpenTransientPlaylists; i++ {
			tmpname := name + " " + strconv.Itoa(i)
			err := db.Where("name = ?", tmpname).First(&Playlist{}).Error
			if err != nil {
				name = tmpname
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("Unable to find new name for playlist location")
			return
		}
		err = nil
	}

	pl.Name = name
	pl.Description = description
	pl.Transient = asTransient
	pl.Open = asTransient
	err = pl.Save()

	return
}

// InsertIntoPlaylist -.
func (b *Playbar) InsertIntoPlaylist(pl *Playlist, position int,
	trackIds []int64, locations []string) {

	logw := slog.With(
		"pl", *pl,
		"position", position,
		"track_ids", trackIds,
		"locations", locations,
	)
	logw.Info("Inserting tracks/locations into playlist")

	pts := []PlaylistTrack{}
	err := db.Where("playlist_id = ?", pl.ID).Order("position ASC").
		Find(&pts).
		Error
	if err != nil {
		logw.Error("Failed to find playlist tracks in database", "error", err)
		return
	}

	s, err := pl.createTracks(trackIds, locations)
	if err != nil {
		logw.Error("Failed to create tracks", "error", err)
		return
	}

	list := poser.InsertInto(pointers.FromSlice(pts), position, pointers.FromSlice(s)...)
	pts = pointers.ToValues(list)

	err = db.Session(&gorm.Session{SkipHooks: true}).
		Save(&pts).
		Error
	onerror.NewRecorder(logw).Log(err)

	broadcastOpenPlaylist(pl.ID)
}

// MergePlaylists -.
func (b *Playbar) MergePlaylists(pl1, pl2 *Playlist) (err error) {
	pts1, _ := pl1.GetTracks(0)
	pts2, _ := pl2.GetTracks(0)

	if len(pts2) == 0 {
		return
	}

	s := []PlaylistTrack{}

	for _, x := range pts2 {
		pt := *x
		pt.PlaylistID = pl1.ID
		s = append(s, pt)
	}

	pts1 = poser.AppendTo(pts1, pointers.FromSlice(s)...)
	s = pointers.ToValues(pts1)
	err = db.Save(&s).Error
	if err == nil {
		isTransient := pl2.Transient
		b.CloseEntry(pl2)
		if !isTransient {
			err = b.DestroyEntry(pl2)
		}
	}
	return
}

// MovePlaylistTrack -.
func (b *Playbar) MovePlaylistTrack(pl *Playlist, to, from int) {
	if from == to || from < 1 {
		return
	}

	logw := slog.With(
		"pl", *pl,
		"to", to,
		"from", from,
	)
	logw.Info("Moving track in playlist")

	pts := []PlaylistTrack{}
	err := db.Where("playlist_id = ?", pl.ID).Order("position ASC").
		Find(&pts).
		Error
	if err != nil {
		logw.Error("Failed to find playlist tracks in database", "error", err)
		return
	}

	list := poser.MoveTo(pointers.FromSlice(pts), to, from)
	moved := pointers.ToValues(list)
	onerror.NewRecorder(logw).Log(db.Save(&moved).Error)
}

// OpenEntry opens the given playbar entry.
func (b *Playbar) OpenEntry(pl *Playlist) {
	if pl.Transient && !pl.Open {
		slog.Warn("Ignoring attempt to reopen transient playlist marked for deletiion")
		return
	}
	pl.Open = true
	onerror.Log(pl.Save())
}

// PrependToPlaylist -.
func (b *Playbar) PrependToPlaylist(pl *Playlist, trackIds []int64,
	locations []string) {
	b.InsertIntoPlaylist(pl, 0, trackIds, locations)
}

// QueryInPlaylist -.
func (bar *Playbar) QueryInPlaylist(qy *Query, qybs []QueryBoundaryTx, pl *Playlist) {
	logw := slog.With(
		"pl", *pl,
		"qy", *qy,
		"len(qybs)", len(qybs),
	)
	logw.Info("Appending query result tracks to playlist")

	var lpf []int64
	var ts []*Track

	switch QueryIndex(qy.Idx) {
	case HistoryQuery:
		if pl.QueryID == qy.ID {
			ts, lpf = findHistoryTracks()
		} else {
			ts = findUniqueHistoryTracks()
			lpf = make([]int64, len(ts))
		}
	case TopTracksQuery:
		ts = findTopTracks()
		lpf = make([]int64, len(ts))
	default:
		ts = qy.FindTracks(qybs)
		lpf = make([]int64, len(ts))
	}

	tx := db.Session(&gorm.Session{SkipHooks: true})

	var pts, s []PlaylistTrack

	err := tx.Where("playlist_id = ?", pl.ID).Find(&pts).Error
	if err != nil {
		logw.Error("Failed to find playlist tracks in database", "error", err)
		return
	}

	for i, t := range ts {
		s = append(s, PlaylistTrack{PlaylistID: pl.ID, TrackID: t.ID, Lastplayedfor: lpf[i]})
	}

	list := poser.AppendTo(pointers.FromSlice(pts), pointers.FromSlice(s)...)
	pts = pointers.ToValues(list)
	err = tx.Save(&pts).Error
	onerror.NewRecorder(logw).Log(err)

	if pl.Open {
		broadcastOpenPlaylist(pl.ID)
	}
}

// UpdateEntry updates a playlist.
func (b *Playbar) UpdateEntry(pl *Playlist, name, descr string, groupID int64,
	resetDescr bool, bucket int) (err error) {

	isTransient := pl.Transient
	queryID := pl.QueryID

	newName := pl.Name
	if name != "" {
		newName = name
		isTransient = false
		queryID = 0
	}

	newDescr := pl.Description
	if descr != "" {
		newDescr = descr
	} else if resetDescr {
		newDescr = ""
	}

	newGroupID := pl.PlaylistGroupID
	if groupID > 0 {
		upg := PlaylistGroup{}
		err = upg.Read(groupID)
		if err != nil {
			return
		}
		if upg.PerspectiveID != b.PerspectiveID {
			err = fmt.Errorf("The provided group ID does not match perspective")
			return
		}
		newGroupID = groupID
	} else if groupID < 0 {
		dpg := &PlaylistGroup{}
		err = dpg.ReadDefaultForPerspective(b.PerspectiveID)
		if err == nil {
			if dpg.PerspectiveID != b.PerspectiveID {
				err = fmt.Errorf("The provided group ID does not match perspective")
				return
			}
			newGroupID = dpg.ID
		}
	} else {
		newGroupID = pl.PlaylistGroupID
	}

	pl.Name = newName
	pl.Description = newDescr
	pl.Transient = isTransient
	pl.QueryID = queryID
	pl.PlaylistGroupID = newGroupID

	switch bucket {
	case 1:
		pl.Bucket = true
	case 2:
		pl.Bucket = false
	}
	err = pl.Save()
	return
}

// UpdateGroup updates a playlist.
func (b *Playbar) UpdateGroup(pg *PlaylistGroup, name, descr string,
	resetDescr bool) (err error) {

	newName := pg.Name
	if name != "" {
		newName = name
	}

	newDescr := pg.Description
	if descr != "" {
		newDescr = descr
	} else if resetDescr {
		newDescr = ""
	}

	pg.Name = newName
	pg.Description = newDescr
	err = pg.Save()
	return
}

func (b *Playbar) getPerspectiveIndex() (idx PerspectiveIndex) {
	p := Perspective{}
	if err := p.Read(b.PerspectiveID); err != nil {
		slog.With(
			"perspective_id", b.PerspectiveID,
			"error", err,
		).Error("Failed to read perspective")
		return
	}
	idx = PerspectiveIndex(p.Idx)
	return
}

func (b *Playbar) newPlaylistFromImported(def impexp.Playlist) (pl *Playlist, msgs []string, err error) {
	pg := PlaylistGroup{}
	err = pg.ReadDefaultForPerspective(b.PerspectiveID)
	if err != nil {
		return
	}

	pl = &Playlist{
		Name:            GetTransientNameForPlaylist(0),
		PlaylistGroupID: pg.ID,
		PlaybarID:       b.ID,
	}

	tx := db.Session(&gorm.Session{SkipHooks: true})

	err = pl.CreateTx(tx)
	if err != nil {
		return
	}

	pts := []PlaylistTrack{}
	for _, dt := range def.Tracks() {
		t := Track{}
		err2 := db.Where("location = ?", dt.Location).First(&t).Error
		if err2 != nil {
			t = Track{Location: dt.Location}
			err2 := t.createTransient(tx, dt.ToRaw())
			if err2 != nil {
				msgs = append(
					msgs,
					fmt.Sprintf(
						"Error creating track at `%v`: %v",
						dt.Location,
						err2,
					),
				)
				continue
			}
		}

		pts = append(pts, PlaylistTrack{PlaylistID: pl.ID, TrackID: t.ID})
	}

	pts = reasignPlaylistTrackPositions(pts)
	for i := range pts {
		err2 := tx.Save(&pts[i]).Error
		if err2 != nil {
			t := Track{}
			t.Read(pts[i].TrackID)
			msgs = append(
				msgs,
				fmt.Sprintf(
					"Error saving playlist track at `%v`: %v",
					t.Location,
					err2,
				),
			)
		}
	}

	return
}

// DeactivatePlaybars deactivates all playbars.
func DeactivatePlaybars() {
	pls := []Playlist{}
	err := db.Where("active = 1").
		Find(&pls).
		Error
	if err != nil {
		slog.Error("Failed to find active playlists in database", "error", err)
		return
	}

	if len(pls) == 0 {
		return
	}

	for i := range pls {
		pls[i].Active = false
	}

	onerror.Log(db.Save(&pls).Error)
}

// GetActiveEntry returns the active playlist.
func GetActiveEntry() *Playlist {
	pl := Playlist{}
	db.Where("active = 1").First(&pl)
	return &pl
}

// GetOpenEntries returns the list of open entries.
func GetOpenEntries() (pls []*Playlist, pts []*PlaylistTrack, ts []*Track) {
	pls = []*Playlist{}
	pts = []*PlaylistTrack{}
	ts = []*Track{}

	plsaux := []Playlist{}
	err := db.Where("open = 1").Find(&plsaux).Error
	if err != nil {
		slog.Error("Failed to find all open playlists in database", "error", err)
		return
	}

	ptsaux := []PlaylistTrack{}
	err = db.Preload("Track").
		Joins("JOIN playlist ON playlist_track.playlist_id = playlist.id").
		Where("playlist.open = 1").
		Order("playlist_id").
		Find(&ptsaux).
		Error
	if err != nil {
		slog.Error("Failed to find all open playlists in database", "error", err)
		return
	}
	tsaux := []Track{}
	for i := range ptsaux {
		tsaux = append(tsaux, ptsaux[i].Track)
	}

	pls = pointers.FromSlice(plsaux)
	pts = pointers.FromSlice(ptsaux)
	ts = pointers.FromSlice(tsaux)

	return
}

// GetPlaybarStore returns the initial status of the playbar store.
func GetPlaybarStore() (
	pgs []*PlaylistGroup,
	pls []*Playlist,
	opls []*Playlist,
	opts []*PlaylistTrack,
	ots []*Track,
) {
	pgs = []*PlaylistGroup{}
	pls = []*Playlist{}
	opls = []*Playlist{}
	opts = []*PlaylistTrack{}
	ots = []*Track{}

	pgsaux := []PlaylistGroup{}
	err := db.Joins("Perspective").
		Where("hidden = 0").
		Order("perspective_id").
		Find(&pgsaux).
		Error
	if err != nil {
		slog.Error("Failed to find all public playlists groups in database", "error", err)
		return
	}

	plsaux := []Playlist{}
	err = db.Joins("Playbar").
		Order("playbar_id").
		Find(&plsaux).
		Error
	if err != nil {
		slog.Error("Failed to find all playlists in database", "error", err)
		return
	}

	pls = pointers.FromSlice(plsaux)
	pgs = pointers.FromSlice(pgsaux)
	opls, opts, ots = GetOpenEntries()

	return
}

func importPlaylist(location string) (impexp.Playlist, error) {
	path, err := urlstr.URLToPath(location)
	if err != nil {
		return nil, err
	}

	if !base.IsSupportedPlaylist(path) {
		err = fmt.Errorf("Unsupported playlist format: %v", location)
		return nil, err
	}

	def, err := impexp.NewFromPath(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	err = def.Parse(f)
	if err != nil {
		return nil, err
	}

	return def, nil
}
