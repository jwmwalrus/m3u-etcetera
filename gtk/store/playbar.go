package store

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
)

type playlistModelRow struct {
	trackID int64
	path    string
}

type playlistModel struct {
	id    int64
	model *gtk.ListStore
	rows  map[int]playlistModelRow
}
type playbarData struct {
	subscriptionID              string
	activeID                    int64
	openPlaylist                []*m3uetcpb.Playlist
	openPlaylistTrack           []*m3uetcpb.PlaylistTrack
	openTrack                   []*m3uetcpb.Track
	playlistGroup               []*m3uetcpb.PlaylistGroup
	playlist                    []*m3uetcpb.Playlist
	playlistReplacementID       int64
	playlistTrackReplacementIDs []int64

	mu sync.RWMutex
}

var (
	// BData playbar store.
	BData = &playbarData{}

	// PerspectiveToPlaylists -.
	PerspectiveToPlaylists map[m3uetcpb.Perspective][]*m3uetcpb.Playlist

	playlistTrackToTrack    map[int64]*m3uetcpb.Track
	playlistToPlaylistGroup map[int64]*m3uetcpb.PlaylistGroup
	playlists               = []*playlistModel{}
	updatePlaybarView       func()
	playlistGroupsModel     *gtk.ListStore

	barTree          playbarTree
	perspectivesList []m3uetcpb.Perspective
)

func init() {
	perspectivesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_RADIO,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}

	barTree.pplt = map[m3uetcpb.Perspective]playlistTree{
		m3uetcpb.Perspective_MUSIC:      {},
		m3uetcpb.Perspective_RADIO:      {},
		m3uetcpb.Perspective_PODCASTS:   {},
		m3uetcpb.Perspective_AUDIOBOOKS: {},
	}
}

func (bd *playbarData) ActiveID() int64 {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	return bd.activeID
}

// GetOpenPlaylist returns the open playlist for the given id.
func (bd *playbarData) GetOpenPlaylist(id int64) (pl *m3uetcpb.Playlist) {
	logw := slog.With("id", id)
	logw.Info("Returning open playlist")

	bd.mu.RLock()
	defer bd.mu.RUnlock()

	for _, pl := range bd.openPlaylist {
		if pl.Id == id {
			return pl
		}
	}

	logw.Debug("No open playlist exists for the given ID")
	return nil
}

// GetOpenPlaylists returns the list of open playlists.
func (bd *playbarData) GetOpenPlaylists() []*m3uetcpb.Playlist {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	return bd.openPlaylist
}

// GetPlaylist returns the playlist for the given id.
func (bd *playbarData) GetPlaylist(id int64) *m3uetcpb.Playlist {
	logw := slog.With("id", id)
	logw.Info("Returning playlist")

	bd.mu.RLock()
	defer bd.mu.RUnlock()

	for _, pl := range bd.playlist {
		if pl.Id == id {
			return pl
		}
	}

	logw.Debug("No playlist exists for the given ID")
	return nil
}

// GetPlaylistGroup returns the playlist group for the given id.
func (bd *playbarData) GetPlaylistGroup(id int64) *m3uetcpb.PlaylistGroup {
	logw := slog.With("id", id)
	logw.Info("Returning playlist group")

	bd.mu.RLock()
	defer bd.mu.RUnlock()

	for _, pg := range bd.playlistGroup {
		if pg.Id == id {
			return pg
		}
	}

	logw.Debug("No playlist group exists for the given ID")
	return nil
}

// PlaylistGroupActionsChanges returns the list of playlist group actions
// to be applied.
func (bd *playbarData) PlaylistGroupActionsChanges() (toRemove []int64) {
	slog.Info("Returning playlist group actions changes")

	model := playlistGroupsModel

	bd.mu.Lock()
	defer bd.mu.Unlock()

	iter, ok := model.IterFirst()
	for ok {
		row, err := GetTreeModelValues(
			&model.TreeModel,
			iter,
			[]ModelColumn{
				PGColPlaylistGroupID,
				PGColActionRemove,
			},
		)
		if err != nil {
			slog.Error("Failed to get tree-model values", "error", err)
			return
		}
		id := row[PGColPlaylistGroupID].(int64)
		for _, pg := range bd.playlistGroup {
			if id != pg.Id {
				continue
			}

			remove := row[PGColActionRemove].(bool)
			if remove {
				toRemove = append(toRemove, id)
			}
		}
		ok = model.IterNext(iter)
	}

	return
}

func (bd *playbarData) PlaylistGroupNames() map[int64]string {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	out := make(map[int64]string)
	for _, pg := range bd.playlistGroup {
		out[pg.Id] = pg.Name
	}
	return out
}

func (bd *playbarData) PlaylistTracksCount(id int64) int64 {
	logw := slog.With("id", id)
	logw.Info("Returning playlist tracks count")

	bd.mu.RLock()
	defer bd.mu.RUnlock()

	var count int64
	for _, opt := range bd.openPlaylistTrack {
		if opt.PlaylistId == id {
			count++
		}
	}

	logw.Debug("Playlist tracks count", "count", count)
	return count
}

func (bd *playbarData) SubscriptionID() string {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	return bd.subscriptionID
}

func (bd *playbarData) GetUpdatePlaylistGroupRequests() (
	[]*m3uetcpb.ExecutePlaylistGroupActionRequest, error) {

	slog.Info("Returning update-playlist-group requests")

	model := playlistGroupsModel

	requests := []*m3uetcpb.ExecutePlaylistGroupActionRequest{}

	iter, ok := model.IterFirst()
	for ok {
		row, err := GetTreeModelValues(
			&model.TreeModel,
			iter,
			[]ModelColumn{PGColPlaylistGroupID, PGColName, PGColDescription},
		)
		if err != nil {
			return nil, err
		}
		id := row[PGColPlaylistGroupID].(int64)
		for _, pg := range bd.playlistGroup {
			if id != pg.Id {
				continue
			}

			req := &m3uetcpb.ExecutePlaylistGroupActionRequest{
				Action: m3uetcpb.PlaylistGroupAction_PG_UPDATE,
				Id:     pg.Id,
			}

			newName := row[PGColName].(string)
			if newName != pg.Name && newName != "" {
				req.Name = newName
			}

			newDescription := row[PGColDescription].(string)
			if newDescription != pg.Description {
				if newDescription != "" {
					req.Description = newDescription
				} else {
					req.ResetDescription = true
				}
			}

			requests = append(requests, req)
			break
		}
		ok = model.IterNext(iter)
	}

	return requests, nil
}

func (bd *playbarData) HasLastPlayedFor(id int64) bool {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	var lpf int64
	for _, opt := range bd.openPlaylistTrack {
		if opt.PlaylistId == id {
			lpf += opt.Lastplayedfor
		}
	}

	return lpf > 0
}

// PlaylistAlreadyExists returns true if a playlist with the given
// name already exists.
func (bd *playbarData) PlaylistAlreadyExists(name string) bool {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	for _, pl := range bd.playlist {
		if strings.EqualFold(pl.Name, name) {
			return true
		}
	}

	return false
}

// PlaylistIsOpen returns true if the playlist with the given
// id is already open.
func (bd *playbarData) PlaylistIsOpen(id int64) bool {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	for _, op := range bd.openPlaylist {
		if op.Id == id {
			return true
		}
	}

	return false
}

// PlaylistGroupAlreadyExists returns true if a playlist group with the
// given name already exists.
func (bd *playbarData) PlaylistGroupAlreadyExists(name string) bool {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	for _, pg := range bd.playlistGroup {
		if strings.EqualFold(pg.Name, name) {
			return true
		}
	}

	return false
}

func (bd *playbarData) ProcessSubscriptionResponse(
	res *m3uetcpb.SubscribeToPlaybarStoreResponse) {

	bd.mu.Lock()
	defer bd.mu.Unlock()

	if bd.subscriptionID == "" {
		bd.subscriptionID = res.SubscriptionId
	}

	switch res.Event {
	case m3uetcpb.PlaybarEvent_BE_INITIAL:
		barTree.setInitialMode(true)
		bd.activeID = 0
		bd.openPlaylist = []*m3uetcpb.Playlist{}
		bd.openPlaylistTrack = []*m3uetcpb.PlaylistTrack{}
		bd.openTrack = []*m3uetcpb.Track{}
		bd.playlistGroup = []*m3uetcpb.PlaylistGroup{}
		bd.playlist = []*m3uetcpb.Playlist{}

		bd.playlistReplacementID = 0
		bd.playlistTrackReplacementIDs = []int64{}
	case m3uetcpb.PlaybarEvent_BE_INITIAL_ITEM:
		bd.activeID = res.ActivePlaylistId
		bd.appendBDataItem(res)
	case m3uetcpb.PlaybarEvent_BE_INITIAL_DONE:
		barTree.setInitialMode(false)
	case m3uetcpb.PlaybarEvent_BE_ITEM_ADDED:
		bd.activeID = res.ActivePlaylistId
		bd.appendBDataItem(res)
	case m3uetcpb.PlaybarEvent_BE_ITEM_CHANGED:
		bd.activeID = res.ActivePlaylistId
		bd.appendBDataItem(res)
	case m3uetcpb.PlaybarEvent_BE_ITEM_REMOVED:
		bd.activeID = res.ActivePlaylistId
		bd.removeBDataItem(res)
	case m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS:
		barTree.setReceivingOpenItems(true)
	case m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS_ITEM:
		bd.appendBDataItem(res)
		bd.trackBDataItemReplacements(res)
	case m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS_DONE:
		bd.processBDataItemReplacements()
		barTree.setReceivingOpenItems(false)
	}

	if barTree.canBeUpdated() {
		glib.IdleAdd(bd.updatePlaybarModel)
		glib.IdleAdd(barTree.update)
		glib.IdleAdd(bd.updatePlaylistGroupModel)
	}
}

func (bd *playbarData) appendBDataItem(res *m3uetcpb.SubscribeToPlaybarStoreResponse) {
	// NOTE: bd.mu lock is already set

	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylist:
		item := res.GetOpenPlaylist()
		for i := range bd.openPlaylist {
			if bd.openPlaylist[i].Id == item.Id {
				bd.openPlaylist[i] = item
				return
			}
		}
		bd.openPlaylist = append(
			bd.openPlaylist,
			item,
		)
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylistTrack:
		item := res.GetOpenPlaylistTrack()
		for i := range bd.openPlaylistTrack {
			if bd.openPlaylistTrack[i].Id == item.Id {
				bd.openPlaylistTrack[i] = item
				return
			}
		}
		bd.openPlaylistTrack = append(
			bd.openPlaylistTrack,
			item,
		)
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenTrack:
		item := res.GetOpenTrack()
		for i := range bd.openTrack {
			if bd.openTrack[i].Id == item.Id {
				bd.openTrack[i] = item
				return
			}
		}
		bd.openTrack = append(bd.openTrack, item)
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_PlaylistGroup:
		item := res.GetPlaylistGroup()
		for i := range bd.playlistGroup {
			if bd.playlistGroup[i].Id == item.Id {
				bd.playlistGroup[i] = item
				return
			}
		}
		bd.playlistGroup = append(
			bd.playlistGroup,
			item,
		)
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_Playlist:
		item := res.GetPlaylist()
		for i := range bd.playlist {
			if bd.playlist[i].Id == item.Id {
				bd.playlist[i] = item
				return
			}
		}
		bd.playlist = append(bd.playlist, item)
	default:
	}
}

func (bd *playbarData) processBDataItemReplacements() {
	// NOTE: bd.mu lock is already set

	defer func() {
		bd.playlistReplacementID = 0
		bd.playlistTrackReplacementIDs = []int64{}
	}()

	var pl *m3uetcpb.Playlist
	for _, opl := range bd.openPlaylist {
		if opl.Id == bd.playlistReplacementID {
			pl = opl
			break
		}
	}

	if pl.Open {
		newpts := []*m3uetcpb.PlaylistTrack{}
		newts := []*m3uetcpb.Track{}

		for i := range bd.openPlaylistTrack {
			if bd.openPlaylistTrack[i].PlaylistId != pl.Id ||
				slices.Contains(
					bd.playlistTrackReplacementIDs,
					bd.openPlaylistTrack[i].Id,
				) {

				newpts = append(newpts, bd.openPlaylistTrack[i])
				for j := range bd.openTrack {
					if bd.openPlaylistTrack[i].TrackId ==
						bd.openTrack[j].Id {
						newts = append(newts, bd.openTrack[j])
						break
					}
				}
			}
		}

		bd.openPlaylistTrack = newpts
		bd.openTrack = newts
	} else {
		newpts := []*m3uetcpb.PlaylistTrack{}
		newts := []*m3uetcpb.Track{}

		for i := range bd.openPlaylistTrack {
			if bd.openPlaylistTrack[i].PlaylistId == pl.Id {
				continue
			}

			newpts = append(newpts, bd.openPlaylistTrack[i])
			for j := range bd.openTrack {
				if bd.openPlaylistTrack[i].TrackId ==
					bd.openTrack[j].Id {
					newts = append(newts, bd.openTrack[j])
					break
				}
			}
		}
		bd.openPlaylistTrack = newpts
		bd.openTrack = newts

		for i := range bd.openPlaylist {
			if pl.Id == bd.openPlaylist[i].Id {
				n := len(bd.openPlaylist)
				bd.openPlaylist[i] = bd.openPlaylist[n-1]
				bd.openPlaylist = bd.openPlaylist[:n-1]
				break
			}
		}
	}
}

func (bd *playbarData) removeBDataItem(res *m3uetcpb.SubscribeToPlaybarStoreResponse) {
	// NOTE: bd.mu lock is already set

	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylist:
		item := res.GetOpenPlaylist()

		n := len(bd.openPlaylist)
		for i := range bd.openPlaylist {
			if bd.openPlaylist[i].Id == item.Id {
				bd.openPlaylist[i] = bd.openPlaylist[n-1]
				bd.openPlaylist = bd.openPlaylist[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylistTrack:
		item := res.GetOpenPlaylistTrack()

		n := len(bd.openPlaylistTrack)
		for i := range bd.openPlaylistTrack {
			if bd.openPlaylistTrack[i].Id == item.Id {
				bd.openPlaylistTrack[i] = bd.openPlaylistTrack[n-1]
				bd.openPlaylistTrack = bd.openPlaylistTrack[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenTrack:
		item := res.GetOpenTrack()

		n := len(bd.openTrack)
		for i := range bd.openTrack {
			if bd.openTrack[i].Id == item.Id {
				bd.openTrack[i] = bd.openTrack[n-1]
				bd.openTrack = bd.openTrack[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_PlaylistGroup:
		item := res.GetPlaylistGroup()

		n := len(bd.playlistGroup)
		for i := range bd.playlistGroup {
			if bd.playlistGroup[i].Id == item.Id {
				bd.playlistGroup[i] = bd.playlistGroup[n-1]
				bd.playlistGroup = bd.playlistGroup[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_Playlist:
		item := res.GetPlaylist()

		n := len(bd.playlist)
		for i := range bd.playlist {
			if bd.playlist[i].Id == item.Id {
				bd.playlist[i] = bd.playlist[n-1]
				bd.playlist = bd.playlist[:n-1]
				break
			}
		}
	}
}

func (bd *playbarData) trackBDataItemReplacements(res *m3uetcpb.SubscribeToPlaybarStoreResponse) {
	// NOTE: bd.mu lock is already set

	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylist:
		bd.playlistReplacementID = res.GetOpenPlaylist().Id
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylistTrack:
		bd.playlistTrackReplacementIDs = append(
			bd.playlistTrackReplacementIDs,
			res.GetOpenPlaylistTrack().Id,
		)
	default:
	}
}

func (bd *playbarData) updatePlaybarMaps() {
	bd.mu.Lock()
	defer bd.mu.Unlock()

	slog.Info("Updating playbar maps")

	sort.SliceStable(bd.openPlaylistTrack, func(i, j int) bool {
		if bd.openPlaylistTrack[i].PlaylistId !=
			bd.openPlaylistTrack[j].PlaylistId {
			return bd.openPlaylistTrack[i].PlaylistId < bd.openPlaylistTrack[j].PlaylistId
		}
		return bd.openPlaylistTrack[i].Position < bd.openPlaylistTrack[j].Position
	})

	playlistTrackToTrack = map[int64]*m3uetcpb.Track{}
	for _, pt := range bd.openPlaylistTrack {
		t := &m3uetcpb.Track{}
		for i := range bd.openTrack {
			if pt.TrackId == bd.openTrack[i].Id {
				t = bd.openTrack[i]
				break
			}
		}
		playlistTrackToTrack[pt.Id] = t
	}

	playlistToPlaylistGroup = map[int64]*m3uetcpb.PlaylistGroup{}
	for _, pl := range bd.playlist {
		pg := &m3uetcpb.PlaylistGroup{}
		for i := range bd.playlistGroup {
			if pl.PlaylistGroupId == bd.playlistGroup[i].Id {
				pg = bd.playlistGroup[i]
				break
			}
		}
		playlistToPlaylistGroup[pl.Id] = pg
	}

	PerspectiveToPlaylists = map[m3uetcpb.Perspective][]*m3uetcpb.Playlist{}
	for i := range bd.openPlaylist {
		var list []*m3uetcpb.Playlist
		list, ok := PerspectiveToPlaylists[bd.openPlaylist[i].Perspective]
		if !ok {
			list = []*m3uetcpb.Playlist{}
		}
		list = append(list, bd.openPlaylist[i])
		PerspectiveToPlaylists[bd.openPlaylist[i].Perspective] = list
	}
}

func (bd *playbarData) updatePlaybarModel() bool {
	if !barTree.canBeUpdated() {
		return false
	}

	slog.Info("Updating playbar model")

	bd.updatePlaybarMaps()

	playbackTrackID := PbData.getTrackID()

	lastPlayedForOverDuration := func(lpf, d int64, dur time.Duration) string {
		diff := math.Abs(float64(lpf)/1e9 - float64(d)/1e9)
		if diff <= 3 {
			return fmt.Sprintf("%v", dur.Truncate(time.Second))
		}
		lpfdur := time.Duration(lpf) * time.Nanosecond
		return fmt.Sprintf("%v / %v", lpfdur.Truncate(time.Second), dur.Truncate(time.Second))
	}

	slog.Debug("Updating open playlists")

	bd.mu.Lock()
	for _, pl := range bd.openPlaylist {
		logw := slog.With("open-pl-id", pl.Id)
		logw.Debug("Updating open playlist")

		model, rows := getPlaylistModel(pl.Id)
		if model == nil {
			logw.Debug("Creating playlist")

			var err error
			model, rows, err = createPlaylistModel(pl.Id)
			if err != nil {
				logw.Error("Failed to create playlist model", "error", err)
				return false
			}
		}

		if model.NColumns() == 0 {
			return false
		}

		logw.Debug("Counting playlist tracks")
		nTracks := 0
		for _, pt := range bd.openPlaylistTrack {
			if pt.PlaylistId != pl.Id {
				continue
			}
			nTracks++
		}

		logw.Debug("Playlist tracks count", "ntracks", nTracks)
		if nTracks == 0 {
			_, ok := model.IterFirst()
			if ok {
				model.Clear()
			}
			continue
		}

		logw.Debug("Updating playlist tracks")
	ptLoop:
		for _, pt := range bd.openPlaylistTrack {
			if pt.PlaylistId != pl.Id {
				continue
			}

			logw2 := logw.With(
				"pt-id", pt.Id,
				"track-id", pt.TrackId,
			)
			logw2.Debug("Updating playlist track")

			weight := 400
			if pt.PlaylistId == bd.activeID &&
				playbackTrackID > 0 &&
				playbackTrackID == pt.TrackId {
				weight = 700
			}

			r, ok := rows[int(pt.Position)]
			if ok {
				if r.trackID == pt.TrackId {
					// change weight
					iter, ok := model.IterFromString(r.path)
					if !ok {
						logw2.Error("Failed to change weight")
						continue
					}
					model.SetValue(
						iter,
						int(TColFontWeight),
						glib.NewValue(weight),
					)
					continue
				}

				// remove
				iter, ok := model.IterFromString(r.path)
				if !ok {
					logw2.Error("Failed to remove row")
					continue
				}
				if !model.Remove(iter) {
					logw2.Error("Failed to remove playlist row", "path", r.path)
				}
				delete(rows, int(pt.Position))
			}

			// add

			t, ok := playlistTrackToTrack[pt.Id]
			if !ok {
				continue
			}

			var iter *gtk.TreeIter
			pos := int(pt.Position)
			if pos == 1 {
				iter = model.Prepend()
			} else {
				var prev *gtk.TreeIter
				for j := pos - 1; j >= 1; j-- {
					r, ok := rows[j]
					if ok {
						prev, ok = model.IterFromString(r.path)
						if !ok {
							logw2.Error("Failed to get iter from string")
							continue ptLoop
						}
						break
					}
				}

				if prev != nil {
					iter = model.InsertAfter(prev)
				} else {
					iter = model.Prepend()
				}
			}

			dur := time.Duration(t.Duration) * time.Nanosecond
			model.Set(
				iter,
				[]int{
					int(TColTrackID),
					int(TColCollectionID),

					int(TColLocation),
					int(TColFormat),
					int(TColType),
					int(TColTitle),
					int(TColAlbum),
					int(TColArtist),
					int(TColAlbumartist),
					int(TColComposer),
					int(TColGenre),

					int(TColYear),
					int(TColTracknumber),
					int(TColTracktotal),
					int(TColTrackNumberOverTotal),
					int(TColDiscnumber),
					int(TColDisctotal),
					int(TColDiscNumberOverTotal),
					int(TColLyrics),
					int(TColComment),
					int(TColPlaycount),

					int(TColRating),
					int(TColDuration),
					int(TColPlayedOverDuration),
					int(TColRemote),
					int(TColLastplayed),
					int(TColPosition),
					int(TColDynamic),

					int(TColLastPosition),
					int(TColFontWeight),
				},
				[]glib.Value{
					*glib.NewValue(t.Id),
					*glib.NewValue(t.CollectionId),

					*glib.NewValue(t.Location),
					*glib.NewValue(t.Format),
					*glib.NewValue(t.Type),
					*glib.NewValue(t.Title),
					*glib.NewValue(t.Album),
					*glib.NewValue(t.Artist),
					*glib.NewValue(t.Albumartist),
					*glib.NewValue(t.Composer),
					*glib.NewValue(t.Genre),

					*glib.NewValue(int(t.Year)),
					*glib.NewValue(int(t.Tracknumber)),
					*glib.NewValue(int(t.Tracktotal)),
					*glib.NewValue(fmt.Sprintf("%d / %d", t.Tracknumber, t.Tracktotal)),
					*glib.NewValue(int(t.Discnumber)),
					*glib.NewValue(int(t.Disctotal)),
					*glib.NewValue(fmt.Sprintf("%d / %d", t.Discnumber, t.Disctotal)),
					*glib.NewValue(t.Lyrics),
					*glib.NewValue(t.Comment),
					*glib.NewValue(int(t.Playcount)),
					*glib.NewValue(int(t.Rating)),
					*glib.NewValue(fmt.Sprint(dur.Truncate(time.Second))),
					*glib.NewValue(lastPlayedForOverDuration(pt.Lastplayedfor, t.Duration, dur)),
					*glib.NewValue(t.Remote),
					*glib.NewValue(time.Unix(0, t.Lastplayed).Format(lastPlayedLayout)),
					*glib.NewValue(int(pt.Position)),
					*glib.NewValue(pt.Dynamic),

					*glib.NewValue(nTracks),
					*glib.NewValue(weight),
				},
			)

			var ps string
			path := model.Path(iter)
			if path == nil {
				logw2.Error("Failed to get iter-path")
			} else {
				ps = path.String()
			}
			rows[pos] = playlistModelRow{pt.TrackId, ps}
		}

		newRows := map[int]playlistModelRow{}
		for p, r := range rows {
			if p > nTracks {
				iter, ok := model.IterFromString(r.path)
				if !ok {
					logw.Error("Failed to get iter from string")
					continue
				}
				if !model.Remove(iter) {
					logw.Error("Failed to remove residual row")
				}
				continue
			}
			newRows[p] = r
		}
		setPlaylistModelRows(pl.Id, newRows)
	}
	bd.mu.Unlock()

	updatePlaybarView()

	return false
}

func (bd *playbarData) updatePlaylistGroupModel() bool {
	if !barTree.canBeUpdated() {
		return false
	}

	slog.Info("Updating playlist-group model")

	model := playlistGroupsModel
	if model == nil {
		return false
	}

	if model.NColumns() == 0 {
		return false
	}

	_, ok := model.IterFirst()
	if ok {
		model.Clear()
	}

	bd.mu.Lock()
	defer bd.mu.Unlock()

	var iter *gtk.TreeIter
	for _, pg := range bd.playlistGroup {
		iter = model.Append()
		model.Set(
			iter,
			[]int{
				int(PGColPlaylistGroupID),
				int(PGColName),
				int(PGColDescription),
				int(PGColPerspective),
				int(PGColActionRemove),
			},
			[]glib.Value{
				*glib.NewValue(pg.Id),
				*glib.NewValue(pg.Name),
				*glib.NewValue(pg.Description),
				*glib.NewValue(pg.Perspective.String()),
				*glib.NewValue(false),
			},
		)
	}

	return false
}

// CreatePlaylistModel creates a playlist model.
func CreatePlaylistModel(id int64) (model *gtk.ListStore, err error) {
	slog.Info("Creating playlist model", "id", id)

	model, _, err = createPlaylistModel(id)
	return
}

// CreatePlaylistGroupsModel creates a playlist model.
func CreatePlaylistGroupsModel() (model *gtk.ListStore, err error) {
	slog.Info("Creating playlist group model")

	playlistGroupsModel = gtk.NewListStore(PGColumns.getTypes())
	if playlistGroupsModel == nil {
		err = fmt.Errorf("failed to create list-store")
		return
	}

	model = playlistGroupsModel
	return
}

// CreatePlaylistsTreeModel creates a playlist model.
func CreatePlaylistsTreeModel(p m3uetcpb.Perspective) (
	model *gtk.TreeStore, err error) {

	slog.Info("Creating playlists tree model", "perspective", p)

	model = gtk.NewTreeStore(PLTreeColumn.getTypes())
	if model == nil {
		err = fmt.Errorf("failed to create tree-store")
		return
	}

	barTree.setPlaylistTree(p, playlistTree{model: model})
	return
}

// DestroyPlaylistModel destroy a playlist model.
func DestroyPlaylistModel(id int64) (err error) {
	slog.Info("Destroying playlist model", "id", id)

	n := len(playlists)
	for i := range playlists {
		if playlists[i].id == id {
			playlists[i] = playlists[n-1]
			playlists = playlists[:n-1]
			return
		}
	}

	err = fmt.Errorf("Playlist model is not in store")
	return
}

// FilterPlaylistTreeBy filters the playlist tree by the given value.
func FilterPlaylistTreeBy(p m3uetcpb.Perspective, val string) {
	v := barTree.getPlaylistTree(p)
	v.filterVal = val
	barTree.setPlaylistTree(p, v)
	barTree.update()
}

// GetPlaylistModel returns the playlist model for the given ID.
func GetPlaylistModel(id int64) *gtk.ListStore {
	model, _ := getPlaylistModel(id)
	return model
}

// GetPlaylistsTreeModel returns the current playlist tree model.
func GetPlaylistsTreeModel(p m3uetcpb.Perspective) *gtk.TreeStore {
	v := barTree.pplt[p]
	return v.model
}

// SetUpdatePlaybarViewFn sets the update-playbar-view function.
func SetUpdatePlaybarViewFn(fn func()) {
	updatePlaybarView = fn
}

func createPlaylistModel(id int64) (model *gtk.ListStore,
	rows map[int]playlistModelRow, err error) {

	logw := slog.With("id", id)
	logw.Info("Creating a playlist model")

	model, rows = getPlaylistModel(id)
	if model != nil {
		logw.Debug("Playlist model already exists")
		return
	}

	model = gtk.NewListStore(TColumns.getTypes())
	if model == nil {
		err = fmt.Errorf("failed to create list-store")
		return
	}

	rows = map[int]playlistModelRow{}

	playlists = append(playlists, &playlistModel{id, model, rows})
	return
}

func getPlaylistModel(id int64) (*gtk.ListStore, map[int]playlistModelRow) {
	logw := slog.With("id", id)
	logw.Debug("Returning playlist model")

	for _, pl := range playlists {
		if pl.id == id {
			return pl.model, pl.rows
		}
	}

	logw.Debug("No playlist model exists for the given ID")
	return nil, nil
}

func setPlaylistModelRows(id int64, rows map[int]playlistModelRow) {
	for _, pl := range playlists {
		if pl.id == id {
			pl.rows = rows
		}
	}
}
