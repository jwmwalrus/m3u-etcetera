package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jwmwalrus/bnp/stringing"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/pkg/impexp"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetPlaybar returns the playbar associated to the given perspective
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

// PlaybarEvent defines a collection event
type PlaybarEvent int

// PlaybarEvent enum
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

// Playbar defines the playlist bar for each perspective
type Playbar struct {
	ID            int64       `json:"id" gorm:"primaryKey"`
	CreatedAt     int64       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     int64       `json:"updatedAt" gorm:"autoUpdateTime"`
	PerspectiveID int64       `json:"perspectiveId" gorm:"uniqueIndex:unique_idx_playbar_perspective_id,not null"`
	Perspective   Perspective `json:"perspective" gorm:"foreignKey:PerspectiveID"`
}

// Read implements the DataReader interface
func (b *Playbar) Read(id int64) error {
	return db.First(b, id).Error
}

// ActivateEntry activates the given entry in a playbar
func (b *Playbar) ActivateEntry(pl *Playlist) {
	log.WithField("pl", pl).
		Info("Activating in playbar")

	pl.Open = true
	if err := pl.Save(); err != nil {
		log.Error(err)
		return
	}

	pls := []Playlist{}
	err := db.Where("playbar_id = ? and open=1", b.ID).Find(&pls).Error
	if err != nil {
		log.Error(err)
		return
	}

	for i := range pls {
		if pls[i].ID == pl.ID {
			pls[i].Active = true
			continue
		}
		pls[i].Active = false
	}
	onerror.Log(db.Save(&pls).Error)
}

// AppendToPlaylist -
func (b *Playbar) AppendToPlaylist(pl *Playlist, trackIds []int64,
	locations []string) {

	log.WithFields(log.Fields{
		"pl":        *pl,
		"trackIds":  trackIds,
		"locations": locations,
	}).
		Info("Appending tracks/locations to playlist")

	pts := []PlaylistTrack{}
	if err := db.Where("playlist_id = ?", pl.ID).Find(&pts).Error; err != nil {
		log.Error(err)
		return
	}

	s, err := pl.createTracks(trackIds, locations)
	if err != nil {
		log.Error(err)
		return
	}

	pts = append(pts, s...)
	pts = reasignPlaylistTrackPositions(pts)

	err = db.Session(&gorm.Session{SkipHooks: true}).
		Save(&pts).
		Error
	onerror.Log(err)

	broadcastOpenPlaylist(pl.ID)
}

// ClearPlaylist -
func (b *Playbar) ClearPlaylist(pl *Playlist) {
	log.WithFields(log.Fields{
		"pl": *pl,
	}).
		Info("Clearing tracks/locations in playlist")

	pts := []PlaylistTrack{}
	if err := db.Where("playlist_id = ?", pl.ID).Find(&pts).Error; err != nil {
		log.Error(err)
		return
	}

	if len(pts) > 0 {
		err := db.Session(&gorm.Session{SkipHooks: true}).
			Where("id > 0").
			Delete(&pts).
			Error
		onerror.Log(err)

		broadcastOpenPlaylist(pl.ID)
	}
}

// CloseEntry closes the given playbar entry
func (b *Playbar) CloseEntry(pl *Playlist) {
	if !pl.Transient {
		pl.DeleteDynamicTracks(db)
	}

	pl.Open = false
	pl.Active = false
	onerror.Log(pl.Save())
}

// CreateEntry creates a playlist
func (b *Playbar) CreateEntry(name, description string) (pl *Playlist, err error) {
	pg := &PlaylistGroup{}
	err = pg.ReadDefaultForPerspective(b.PerspectiveID)
	if err != nil {
		return
	}

	var plname string
	pl = &Playlist{}

	isTransient := false
	if name != "" {
		plname = name
	} else {
		isTransient = true
		plname = GetTransientNameForPlaylist()
	}

	pl = &Playlist{
		Name:            plname,
		Description:     description,
		PlaybarID:       b.ID,
		PlaylistGroupID: pg.ID,
		Transient:       isTransient,
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

// CreateGroup creates a playlist
func (b *Playbar) CreateGroup(name, description string) (pg *PlaylistGroup, err error) {
	pg = &PlaylistGroup{
		Name:          name,
		Description:   description,
		PerspectiveID: b.PerspectiveID,
	}
	err = pg.Create()
	return
}

// DeactivateEntry -
func (b *Playbar) DeactivateEntry(pl *Playlist) {
	pl.Active = false
	onerror.Log(pl.Save())
}

// DeleteFromPlaylist -
func (b *Playbar) DeleteFromPlaylist(pl *Playlist, position int) {
	log.WithFields(log.Fields{
		"pl":       *pl,
		"position": position,
	}).
		Info("Deleting position in playlist")

	pts := []PlaylistTrack{}
	if err := db.Where("playlist_id = ?", pl.ID).Find(&pts).Error; err != nil {
		log.Error(err)
		return
	}
	s := []PlaylistTrack{}
	pt := PlaylistTrack{}
	for i := range pts {
		if pts[i].Position == position {
			pt = pts[i]
			continue
		}
		s = append(s, pts[i])
	}

	if pt.ID > 0 {
		if err := pt.Delete(); err != nil {
			log.Error(err)
			return
		}
	}
	s = reasignPlaylistTrackPositions(s)
	onerror.Log(db.Save(&s).Error)
}

// DestroyEntry deletes a playlist
func (b *Playbar) DestroyEntry(pl *Playlist) error {
	if pl.Transient {
		tx := db.Session(&gorm.Session{SkipHooks: true})
		pld := Playlist{}
		if errd := tx.First(&pld).Error; errd == nil {
			pld.DeleteTx(tx)
		}
		return nil
	}
	return pl.Delete()
}

// DestroyGroup deletes a playlist
func (b *Playbar) DestroyGroup(pg *PlaylistGroup) error {
	return pg.Delete()
}

// GetAllEntries -
func (b *Playbar) GetAllEntries(limit int) (pls []*Playlist) {
	pls = []*Playlist{}

	s := []Playlist{}
	tx := db.Joins("Playbar").
		Where("playbar_id = ?", b.ID)
	if limit > 0 {
		tx.Limit(limit)
	}
	err := tx.Find(&s).Error
	if err != nil {
		return
	}

	for i := range s {
		pls = append(pls, &s[i])
	}
	return
}

// GetAllGroups -
func (b *Playbar) GetAllGroups(limit int) (pgs []*PlaylistGroup) {
	pgs = []*PlaylistGroup{}

	s := []PlaylistGroup{}
	tx := db.Joins("Perspective").
		Where("hidden = 0 and perspective_id = ?", b.PerspectiveID)
	if limit > 0 {
		tx.Limit(limit)
	}
	err := tx.Find(&s).Error
	if err != nil {
		log.Error(err)
		return
	}

	for i := range s {
		pgs = append(pgs, &s[i])
	}
	return
}

// GetAllOpenEntries -
func (b *Playbar) GetAllOpenEntries() (pls []*Playlist) {
	pls = []*Playlist{}

	s := []Playlist{}
	err := db.Joins("Playbar").
		Where("open = 1 and playbar_id = ?", b.ID).
		Find(&s).
		Error
	if err != nil {
		return
	}

	for i := range s {
		pls = append(pls, &s[i])
	}
	return
}

// ImportPlaylist creates a playlist from the given location, if supported
func (b *Playbar) ImportPlaylist(location string) (pl *Playlist, msgs []string, err error) {
	path, err := urlstr.URLToPath(location)
	if err != nil {
		return
	}

	if !base.IsSupportedPlaylist(path) {
		err = fmt.Errorf("Unsupported playlist format: %v", location)
	}

	def, err := impexp.NewFromPath(path)
	if err != nil {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}

	err = def.Parse(f)
	if err != nil {
		return
	}

	name := def.Name()

	if name == "" {
		name = strings.Join([]string{
			strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
			stringing.GetRandomString(8),
			def.Type(),
		}, "-")
	}

	pg := PlaylistGroup{}
	err = pg.ReadDefaultForPerspective(b.PerspectiveID)
	if err != nil {
		return
	}

	pl = &Playlist{
		Name: name,
		Description: "Imported from " + filepath.Base(path) +
			" on " + time.Now().Format(time.RFC3339),
		PlaylistGroupID: pg.ID,
		PlaybarID:       b.ID,
	}

	err = pl.Create()
	if err != nil {
		return
	}

	tx := db.Session(&gorm.Session{SkipHooks: true})
	pts := []PlaylistTrack{}
	for _, dt := range def.Tracks() {
		t := Track{}
		err2 := db.Where("location = ?", dt.Location).First(&t).Error
		if err2 != nil {
			t = Track{Location: dt.Location}
			err2 := t.createTransientWithRaw(tx, dt.ToRaw())
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

// InsertIntoPlaylist -
func (b *Playbar) InsertIntoPlaylist(pl *Playlist, position int,
	trackIds []int64, locations []string) {

	log.WithFields(log.Fields{
		"pl":        *pl,
		"position":  position,
		"trackIds":  trackIds,
		"locations": locations,
	}).
		Info("Inserting tracks/locations into playlist")

	pts := []PlaylistTrack{}
	if err := db.Where("playlist_id = ?", pl.ID).Find(&pts).Error; err != nil {
		log.Error(err)
		return
	}

	if position <= 1 {
		s, err := pl.createTracks(trackIds, locations)
		if err != nil {
			log.Error(err)
			return
		}
		aux := pts
		pts = s
		pts = append(pts, aux...)
	} else if position > 1 && position <= len(pts) {
		s, err := pl.createTracks(trackIds, locations)
		if err != nil {
			log.Error(err)
			return
		}
		aux := pts
		piv := position - 1
		pts = aux[:piv]
		pts = append(pts, s...)
		pts = append(pts, aux[piv:]...)

	} else {
		b.AppendToPlaylist(pl, trackIds, locations)
		return
	}

	pts = reasignPlaylistTrackPositions(pts)

	err := db.Session(&gorm.Session{SkipHooks: true}).
		Save(&pts).
		Error
	onerror.Log(err)

	broadcastOpenPlaylist(pl.ID)
}

// MovePlaylistTrack -
func (b *Playbar) MovePlaylistTrack(pl *Playlist, to, from int) {
	if from == to || from < 1 {
		return
	}

	log.WithFields(log.Fields{
		"pl":   *pl,
		"to":   to,
		"from": from,
	}).
		Info("Moving track in playlist")

	pts := []PlaylistTrack{}
	if err := db.Where("playlist_id = ?", pl.ID).Find(&pts).Error; err != nil {
		log.Error(err)
		return
	}

	var moved, afterPiv []PlaylistTrack
	var piv *PlaylistTrack
	for i := range pts {
		if pts[i].Position == from {
			piv = &pts[i]
		} else if pts[i].Position < to {
			moved = append(moved, pts[i])
		} else if pts[i].Position > to {
			afterPiv = append(afterPiv, pts[i])
		} else if pts[i].Position == to {
			if from < to {
				moved = append(moved, pts[i])
			} else {
				afterPiv = append(afterPiv, pts[i])
			}
		}
	}

	if piv != nil {
		moved = append(moved, *piv)
	}
	moved = append(moved, afterPiv...)

	moved = reasignPlaylistTrackPositions(moved)
	onerror.Log(db.Save(&moved).Error)
}

// OpenEntry opens the given playbar entry
func (b *Playbar) OpenEntry(pl *Playlist) {
	if pl.Transient && !pl.Open {
		log.Warn("Ignoring attempt to reopen transient playlist marked for deletiion")
		return
	}
	pl.Open = true
	onerror.Log(pl.Save())
}

// PrependToPlaylist -
func (b *Playbar) PrependToPlaylist(pl *Playlist, trackIds []int64,
	locations []string) {
	b.InsertIntoPlaylist(pl, 0, trackIds, locations)
}

// UpdateEntry updates a playlist
func (b *Playbar) UpdateEntry(pl *Playlist, name, descr string, groupID int64,
	resetDescr bool) (err error) {

	isTransient := pl.Transient

	newName := pl.Name
	if name != "" {
		newName = name
		isTransient = false
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
		var dpg *PlaylistGroup
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
	pl.PlaylistGroupID = newGroupID
	err = pl.Save()
	return
}

// UpdateGroup updates a playlist
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
		log.Error(err)
		return
	}
	idx = PerspectiveIndex(p.Idx)
	return
}

// DeactivatePlaybars deactivates all playbars
func DeactivatePlaybars() {
	pls := []Playlist{}
	err := db.Where("active = 1").
		Find(&pls).
		Error
	if err != nil {
		log.Error(err)
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

// GetActiveEntry returns the active playlist
func GetActiveEntry() *Playlist {
	pl := Playlist{}
	db.Where("active = 1").First(&pl)
	return &pl
}

// GetOpenEntries returns the list of open entries
func GetOpenEntries() (pls []*Playlist, pts []*PlaylistTrack, ts []*Track) {
	pls = []*Playlist{}
	pts = []*PlaylistTrack{}
	ts = []*Track{}

	pllist := []Playlist{}
	err := db.Where("open = 1").Find(&pllist).Error
	if err != nil {
		log.Error(err)
		return
	}
	for i := range pllist {
		pls = append(pls, &pllist[i])
	}

	ptlist := []PlaylistTrack{}
	err = db.Preload("Track").
		Joins("JOIN playlist ON playlist_track.playlist_id = playlist.id").
		Where("playlist.open = 1").
		Order("playlist_id").
		Find(&ptlist).
		Error
	if err != nil {
		log.Error(err)
		return
	}
	for i := range ptlist {
		pts = append(pts, &ptlist[i])
		ts = append(ts, &ptlist[i].Track)
	}
	return
}

// GetPlaybarStore returns the initial status of the playbar store
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

	pglist := []PlaylistGroup{}
	err := db.Joins("Perspective").
		Where("hidden = 0").
		Order("perspective_id").
		Find(&pglist).
		Error
	if err != nil {
		log.Error(err)
		return
	}

	pllist := []Playlist{}
	err = db.Joins("Playbar").
		Order("playbar_id").
		Find(&pllist).
		Error
	if err != nil {
		log.Error(err)
		return
	}

	for i := range pglist {
		pgs = append(pgs, &pglist[i])
	}

	for i := range pllist {
		pls = append(pls, &pllist[i])
	}

	opls, opts, ots = GetOpenEntries()

	return
}
