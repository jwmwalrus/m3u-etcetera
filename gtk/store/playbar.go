package store

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

type playlistModel struct {
	id    int64
	model *gtk.ListStore
}

var (
	// BData playbar store
	BData struct {
		subscriptionID              string
		Mu                          sync.Mutex
		ActiveID                    int64
		OpenPlaylist                []*m3uetcpb.Playlist
		OpenPlaylistTrack           []*m3uetcpb.PlaylistTrack
		OpenTrack                   []*m3uetcpb.Track
		PlaylistGroup               []*m3uetcpb.PlaylistGroup
		Playlist                    []*m3uetcpb.Playlist
		PlaylistReplacementID       int64
		PlaylistTrackReplacementIDs []int64
	}

	playlistTrackToTrack    map[int64]*m3uetcpb.Track
	playlistToPlaylistGroup map[int64]*m3uetcpb.PlaylistGroup
	playlists               = []*playlistModel{}
	updatePlaybarView       func()
	barTree                 playbarTree
	playlistGroupsModel     *gtk.ListStore

	// PerspectiveToPlaylists -
	PerspectiveToPlaylists map[m3uetcpb.Perspective][]*m3uetcpb.Playlist
)

// CreatePlaylistModel creates a playlist model
func CreatePlaylistModel(id int64) (model *gtk.ListStore, err error) {
	log.Info("Creating a playlist model")

	model = GetPlaylistModel(id)
	if model != nil {
		return
	}

	model, err = gtk.ListStoreNew(TColumns.getTypes()...)
	if err != nil {
		return
	}

	playlists = append(playlists, &playlistModel{id, model})
	return
}

// CreatePlaylistGroupsModel creates a playlist model
func CreatePlaylistGroupsModel() (model *gtk.ListStore, err error) {
	log.Info("Creating a playlist model")

	playlistGroupsModel, err = gtk.ListStoreNew(PGColumns.getTypes()...)
	if err != nil {
		return
	}

	model = playlistGroupsModel
	return
}

// CreatePlaylistsTreeModel creates a playlist model
func CreatePlaylistsTreeModel(p m3uetcpb.Perspective) (model *gtk.TreeStore, err error) {
	log.Info("Creating playlists model")

	model, err = gtk.TreeStoreNew(PLTreeColumn.getTypes()...)
	if err != nil {
		return
	}

	v := barTree.pplt[p]
	v.model = model
	barTree.pplt[p] = v
	return
}

// DestroyPlaylistModel destroy a playlist model
func DestroyPlaylistModel(id int64) (err error) {
	log.Info("Destroying a playlist model")

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

// FilterPlaylistTreeBy filters the playlist tree by the given value
func FilterPlaylistTreeBy(p m3uetcpb.Perspective, val string) {
	v := barTree.pplt[p]
	v.filterVal = val
	barTree.pplt[p] = v
	barTree.update()
}

// GetOpenPlaylist returns the open playlist for the given id.
func GetOpenPlaylist(id int64) (pl *m3uetcpb.Playlist) {
	log.WithField("id", id).
		Info("Returning open playlist")

	BData.Mu.Lock()
	defer BData.Mu.Unlock()
	for _, pl := range BData.OpenPlaylist {
		if pl.Id == id {
			return pl
		}
	}
	return nil
}

// GetPlaylist returns the playlist for the given id.
func GetPlaylist(id int64) *m3uetcpb.Playlist {
	log.WithField("id", id).
		Info("Returning playlist")

	BData.Mu.Lock()
	defer BData.Mu.Unlock()
	for _, pl := range BData.Playlist {
		if pl.Id == id {
			return pl
		}
	}
	return nil
}

// GetPlaylistGroup returns the playlist group for the given id.
func GetPlaylistGroup(id int64) *m3uetcpb.PlaylistGroup {
	log.WithField("id", id).
		Info("Returning playlist group")

	BData.Mu.Lock()
	defer BData.Mu.Unlock()
	for _, pg := range BData.PlaylistGroup {
		if pg.Id == id {
			return pg
		}
	}
	return nil
}

// GetPlaylistModel returns the playlist model for the given ID
func GetPlaylistModel(id int64) *gtk.ListStore {
	log.Info("Returning playlist model")

	for _, pl := range playlists {
		if pl.id == id {
			return pl.model
		}
	}
	return nil
}

// GetPlaylistsTreeModel returns the current playlist tree model
func GetPlaylistsTreeModel(p m3uetcpb.Perspective) *gtk.TreeStore {
	v := barTree.pplt[p]
	return v.model
}

// PlaylistAlreadyExists returns true if a playlist with the given
// name already exists
func PlaylistAlreadyExists(name string) bool {
	BData.Mu.Lock()
	defer BData.Mu.Unlock()

	for _, pl := range BData.Playlist {
		if strings.ToLower(pl.Name) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

// PlaylistGroupAlreadyExists returns true if a playlist group with the
//given name already exists
func PlaylistGroupAlreadyExists(name string) bool {
	BData.Mu.Lock()
	defer BData.Mu.Unlock()

	for _, pg := range BData.PlaylistGroup {
		if strings.ToLower(pg.Name) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

// SetUpdatePlaybarViewFn sets the update-playbar-view function
func SetUpdatePlaybarViewFn(fn func()) {
	updatePlaybarView = fn
}

func updatePlaybarMaps() {
	BData.Mu.Lock()
	defer BData.Mu.Unlock()

	sort.SliceStable(BData.OpenPlaylistTrack, func(i, j int) bool {
		if BData.OpenPlaylistTrack[i].PlaylistId != BData.OpenPlaylistTrack[j].PlaylistId {
			return BData.OpenPlaylistTrack[i].PlaylistId < BData.OpenPlaylistTrack[j].PlaylistId
		}
		return BData.OpenPlaylistTrack[i].Position < BData.OpenPlaylistTrack[j].Position
	})

	playlistTrackToTrack = map[int64]*m3uetcpb.Track{}
	for _, pt := range BData.OpenPlaylistTrack {
		t := &m3uetcpb.Track{}
		for i := range BData.OpenTrack {
			if pt.TrackId == BData.OpenTrack[i].Id {
				t = BData.OpenTrack[i]
				break
			}
		}
		playlistTrackToTrack[pt.Id] = t
	}

	playlistToPlaylistGroup = map[int64]*m3uetcpb.PlaylistGroup{}
	for _, pl := range BData.Playlist {
		pg := &m3uetcpb.PlaylistGroup{}
		for i := range BData.PlaylistGroup {
			if pl.PlaylistGroupId == BData.PlaylistGroup[i].Id {
				pg = BData.PlaylistGroup[i]
				break
			}
		}
		playlistToPlaylistGroup[pl.Id] = pg
	}

	PerspectiveToPlaylists = map[m3uetcpb.Perspective][]*m3uetcpb.Playlist{}
	for i := range BData.OpenPlaylist {
		var list []*m3uetcpb.Playlist
		list, ok := PerspectiveToPlaylists[BData.OpenPlaylist[i].Perspective]
		if !ok {
			list = []*m3uetcpb.Playlist{}
		}
		list = append(list, BData.OpenPlaylist[i])
		PerspectiveToPlaylists[BData.OpenPlaylist[i].Perspective] = list
	}
}

func updatePlaybarModel() bool {
	if barTree.initialMode || barTree.receivingOpenItems {
		return false
	}

	log.Info("Updating playbar model")

	updatePlaybarMaps()

	pbdata.mu.Lock()
	playbackTrackID := pbdata.trackID
	pbdata.mu.Unlock()

	BData.Mu.Lock()
	for _, pl := range BData.OpenPlaylist {
		model := GetPlaylistModel(pl.Id)
		if model == nil {
			var err error
			model, err = CreatePlaylistModel(pl.Id)
			if err != nil {
				log.Error(err)
				return false
			}
		}

		if model.GetNColumns() == 0 {
			return false
		}

		_, ok := model.GetIterFirst()
		if ok {
			model.Clear()
		}

		if len(BData.OpenPlaylistTrack) > 0 {
			var iter *gtk.TreeIter
			nTracks := 0
			for _, pt := range BData.OpenPlaylistTrack {
				if pt.PlaylistId != pl.Id {
					continue
				}
				nTracks++
			}
			for _, pt := range BData.OpenPlaylistTrack {
				if pt.PlaylistId != pl.Id {
					continue
				}

				t, ok := playlistTrackToTrack[pt.Id]
				if !ok {
					continue
				}
				iter = model.Append()

				weight := 400
				if pt.PlaylistId == BData.ActiveID &&
					playbackTrackID > 0 &&
					playbackTrackID == pt.TrackId {
					weight = 700
				}

				dur := time.Duration(t.Duration) * time.Nanosecond
				err := model.Set(
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
						int(TColDiscnumber),
						int(TColDisctotal),
						int(TColLyrics),
						int(TColComment),
						int(TColPlaycount),

						int(TColRating),
						int(TColDuration),
						int(TColRemote),
						int(TColLastplayed),
						int(TColPosition),
						int(TColLastPosition),
						int(TColDynamic),
						int(TColFontWeight),
					},
					[]interface{}{
						t.Id,
						t.CollectionId,
						t.Location,
						t.Format,
						t.Type,
						t.Title,
						t.Album,
						t.Artist,
						t.Albumartist,
						t.Composer,
						t.Genre,

						int(t.Year),
						int(t.Tracknumber),
						int(t.Tracktotal),
						int(t.Discnumber),
						int(t.Disctotal),
						t.Lyrics,
						t.Comment,
						int(t.Playcount),

						int(t.Rating),
						fmt.Sprint(dur.Truncate(time.Second)),
						t.Remote,
						t.Lastplayed,
						int(pt.Position),
						nTracks,
						pt.Dynamic,
						weight,
					},
				)
				onerror.Log(err)
			}
		}
	}
	BData.Mu.Unlock()

	updatePlaybarView()

	return false
}

func updatePlaylistGroupModel() bool {
	if barTree.initialMode || barTree.receivingOpenItems {
		return false
	}

	log.Info("Updating playlit-group model")

	model := playlistGroupsModel
	if model == nil {
		return false
	}

	if model.GetNColumns() == 0 {
		return false
	}

	_, ok := model.GetIterFirst()
	if ok {
		model.Clear()
	}

	BData.Mu.Lock()
	var iter *gtk.TreeIter
	for _, pg := range BData.PlaylistGroup {
		iter = model.Append()
		err := model.Set(
			iter,
			[]int{
				int(PGColPlaylistGroupID),
				int(PGColName),
				int(PGColDescription),
				int(PGColPerspective),
			},
			[]interface{}{
				pg.Id,
				pg.Name,
				pg.Description,
				pg.Perspective.String(),
			},
		)
		onerror.Log(err)
	}
	BData.Mu.Unlock()

	return false
}
