package store

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
	barTree                 playbarsTree
	playlistGroupsModel     *gtk.ListStore

	// PerspectiveToPlaylists -
	PerspectiveToPlaylists map[m3uetcpb.Perspective][]*m3uetcpb.Playlist
)

// ApplyPlaylistGroupChanges applies collection changes
func ApplyPlaylistGroupChanges() {
	log.Info("Applying playlist-group changes")

	cc, err := getClientConn()
	if err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)

	model := playlistGroupsModel

	iter, ok := model.GetIterFirst()
	for ok {
		row, err := GetListStoreModelValues(model, iter, []ModelColumn{PGColPlaylistGroupID, PGColName, PGColDescription})
		if err != nil {
			log.Error(err)
			return
		}
		id := row[PGColPlaylistGroupID].(int64)
		for _, pg := range BData.PlaylistGroup {
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

			_, err := cl.ExecutePlaylistGroupAction(context.Background(), req)
			onerror.Log(err)
			break
		}
		ok = model.IterNext(iter)
	}
}

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

// ExecutePlaybarAction -
func ExecutePlaybarAction(req *m3uetcpb.ExecutePlaybarActionRequest) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	_, err = cl.ExecutePlaybarAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		log.Error(s.Message())
		return
	}
	return
}

// ExecutePlaylistAction -
func ExecutePlaylistAction(req *m3uetcpb.ExecutePlaylistActionRequest) (*m3uetcpb.ExecutePlaylistActionResponse, error) {
	cc, err := GetClientConn()
	if err != nil {
		return nil, err
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	res, err := cl.ExecutePlaylistAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return nil, err
	}
	return res, nil
}

// ExecutePlaylistGroupAction -
func ExecutePlaylistGroupAction(req *m3uetcpb.ExecutePlaylistGroupActionRequest) (*m3uetcpb.ExecutePlaylistGroupActionResponse, error) {
	cc, err := GetClientConn()
	if err != nil {
		return nil, err
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	res, err := cl.ExecutePlaylistGroupAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		err = fmt.Errorf(s.Message())
		return nil, err
	}
	return res, nil
}

// ExecutePlaylistTrackAction -
func ExecutePlaylistTrackAction(req *m3uetcpb.ExecutePlaylistTrackActionRequest) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	_, err = cl.ExecutePlaylistTrackAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		log.Error(s.Message())
		return
	}
	return
}

// FilterPlaylistTreeBy filters the collections by the given value
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

// GetPlaylistsTreeModel returns the current collections model
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

func appendBDataItem(res *m3uetcpb.SubscribeToPlaybarStoreResponse) {
	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylist:
		item := res.GetOpenPlaylist()
		for i := range BData.OpenPlaylist {
			if BData.OpenPlaylist[i].Id == item.Id {
				BData.OpenPlaylist[i] = item
				return
			}
		}
		BData.OpenPlaylist = append(
			BData.OpenPlaylist,
			item,
		)
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylistTrack:
		item := res.GetOpenPlaylistTrack()
		for i := range BData.OpenPlaylistTrack {
			if BData.OpenPlaylistTrack[i].Id == item.Id {
				BData.OpenPlaylistTrack[i] = item
				return
			}
		}
		BData.OpenPlaylistTrack = append(
			BData.OpenPlaylistTrack,
			item,
		)
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenTrack:
		item := res.GetOpenTrack()
		for i := range BData.OpenTrack {
			if BData.OpenTrack[i].Id == item.Id {
				BData.OpenTrack[i] = item
				return
			}
		}
		BData.OpenTrack = append(BData.OpenTrack, item)
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_PlaylistGroup:
		item := res.GetPlaylistGroup()
		for i := range BData.PlaylistGroup {
			if BData.PlaylistGroup[i].Id == item.Id {
				BData.PlaylistGroup[i] = item
				return
			}
		}
		BData.PlaylistGroup = append(
			BData.PlaylistGroup,
			item,
		)
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_Playlist:
		item := res.GetPlaylist()
		for i := range BData.Playlist {
			if BData.Playlist[i].Id == item.Id {
				BData.Playlist[i] = item
				return
			}
		}
		BData.Playlist = append(BData.Playlist, item)
	default:
	}
}

func processBDataItemReplacements() {
	defer func() {
		BData.PlaylistReplacementID = 0
		BData.PlaylistTrackReplacementIDs = []int64{}
	}()

	var pl *m3uetcpb.Playlist
	for _, opl := range BData.OpenPlaylist {
		if opl.Id == BData.PlaylistReplacementID {
			pl = opl
			break
		}
	}

	if pl.Open {
		newpts := []*m3uetcpb.PlaylistTrack{}
		newts := []*m3uetcpb.Track{}

		for i := range BData.OpenPlaylistTrack {
			if BData.OpenPlaylistTrack[i].PlaylistId != pl.Id ||
				slice.Contains(
					BData.PlaylistTrackReplacementIDs,
					BData.OpenPlaylistTrack[i].Id,
				) {

				newpts = append(newpts, BData.OpenPlaylistTrack[i])
				for j := range BData.OpenTrack {
					if BData.OpenPlaylistTrack[i].TrackId ==
						BData.OpenTrack[j].Id {
						newts = append(newts, BData.OpenTrack[j])
						break
					}
				}
			}
		}

		BData.OpenPlaylistTrack = newpts
		BData.OpenTrack = newts
	} else {
		newpts := []*m3uetcpb.PlaylistTrack{}
		newts := []*m3uetcpb.Track{}

		for i := range BData.OpenPlaylistTrack {
			if BData.OpenPlaylistTrack[i].PlaylistId == pl.Id {
				continue
			}

			newpts = append(newpts, BData.OpenPlaylistTrack[i])
			for j := range BData.OpenTrack {
				if BData.OpenPlaylistTrack[i].TrackId ==
					BData.OpenTrack[j].Id {
					newts = append(newts, BData.OpenTrack[j])
					break
				}
			}
		}
		BData.OpenPlaylistTrack = newpts
		BData.OpenTrack = newts

		for i := range BData.OpenPlaylist {
			if pl.Id == BData.OpenPlaylist[i].Id {
				n := len(BData.OpenPlaylist)
				BData.OpenPlaylist[i] = BData.OpenPlaylist[n-1]
				BData.OpenPlaylist = BData.OpenPlaylist[:n-1]
				break
			}
		}
	}
}

func removeBDataItem(res *m3uetcpb.SubscribeToPlaybarStoreResponse) {
	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylist:
		item := res.GetOpenPlaylist()
		n := len(BData.OpenPlaylist)
		for i := range BData.OpenPlaylist {
			if BData.OpenPlaylist[i].Id == item.Id {
				BData.OpenPlaylist[i] = BData.OpenPlaylist[n-1]
				BData.OpenPlaylist = BData.OpenPlaylist[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylistTrack:
		item := res.GetOpenPlaylistTrack()
		n := len(BData.OpenPlaylistTrack)
		for i := range BData.OpenPlaylistTrack {
			if BData.OpenPlaylistTrack[i].Id == item.Id {
				BData.OpenPlaylistTrack[i] = BData.OpenPlaylistTrack[n-1]
				BData.OpenPlaylistTrack = BData.OpenPlaylistTrack[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenTrack:
		item := res.GetOpenTrack()
		n := len(BData.OpenTrack)
		for i := range BData.OpenTrack {
			if BData.OpenTrack[i].Id == item.Id {
				BData.OpenTrack[i] = BData.OpenTrack[n-1]
				BData.OpenTrack = BData.OpenTrack[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_PlaylistGroup:
		item := res.GetPlaylistGroup()
		n := len(BData.PlaylistGroup)
		for i := range BData.PlaylistGroup {
			if BData.PlaylistGroup[i].Id == item.Id {
				BData.PlaylistGroup[i] = BData.PlaylistGroup[n-1]
				BData.PlaylistGroup = BData.PlaylistGroup[:n-1]
				break
			}
		}
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_Playlist:
		item := res.GetPlaylist()
		n := len(BData.Playlist)
		for i := range BData.Playlist {
			if BData.Playlist[i].Id == item.Id {
				BData.Playlist[i] = BData.Playlist[n-1]
				BData.Playlist = BData.Playlist[:n-1]
				break
			}
		}
	}
}

func subscribeToPlaybarStore() {
	log.Info("Subscribing to playbar store")

	defer wgplaybar.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	stream, err := cl.SubscribeToPlaybarStore(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Errorf("Error subscribing to playbar store: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		BData.Mu.Lock()

		if BData.subscriptionID == "" {
			BData.subscriptionID = res.SubscriptionId
		}

		switch res.Event {
		case m3uetcpb.PlaybarEvent_BE_INITIAL:
			barTree.initialMode = true
			BData.ActiveID = 0
			BData.OpenPlaylist = []*m3uetcpb.Playlist{}
			BData.OpenPlaylistTrack = []*m3uetcpb.PlaylistTrack{}
			BData.OpenTrack = []*m3uetcpb.Track{}
			BData.PlaylistGroup = []*m3uetcpb.PlaylistGroup{}
			BData.Playlist = []*m3uetcpb.Playlist{}

			BData.PlaylistReplacementID = 0
			BData.PlaylistTrackReplacementIDs = []int64{}
		case m3uetcpb.PlaybarEvent_BE_INITIAL_ITEM:
			BData.ActiveID = res.ActivePlaylistId
			appendBDataItem(res)
		case m3uetcpb.PlaybarEvent_BE_INITIAL_DONE:
			barTree.initialMode = false
		case m3uetcpb.PlaybarEvent_BE_ITEM_ADDED:
			BData.ActiveID = res.ActivePlaylistId
			appendBDataItem(res)
		case m3uetcpb.PlaybarEvent_BE_ITEM_CHANGED:
			BData.ActiveID = res.ActivePlaylistId
			appendBDataItem(res)
		case m3uetcpb.PlaybarEvent_BE_ITEM_REMOVED:
			BData.ActiveID = res.ActivePlaylistId
			removeBDataItem(res)
		case m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS:
			barTree.receivingOpenItems = true
		case m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS_ITEM:
			appendBDataItem(res)
			trackBDataItemReplacements(res)
		case m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS_DONE:
			processBDataItemReplacements()
			barTree.receivingOpenItems = false
		}

		BData.Mu.Unlock()

		if !barTree.initialMode && !barTree.receivingOpenItems {
			glib.IdleAdd(updatePlaybarModel)
			glib.IdleAdd(barTree.update)
			glib.IdleAdd(updatePlaylistGroupModel)
		}

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}

}

func trackBDataItemReplacements(res *m3uetcpb.SubscribeToPlaybarStoreResponse) {
	switch res.Item.(type) {
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylist:
		BData.PlaylistReplacementID = res.GetOpenPlaylist().Id
	case *m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylistTrack:
		BData.PlaylistTrackReplacementIDs = append(
			BData.PlaylistTrackReplacementIDs,
			res.GetOpenPlaylistTrack().Id,
		)
	default:
	}
}

func unsubscribeFromPlaybarStore() {
	log.Info("Unsubscribing from playbar store")

	BData.Mu.Lock()
	id := BData.subscriptionID
	BData.Mu.Unlock()

	var cc *grpc.ClientConn
	var err error
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	_, err = cl.UnsubscribeFromPlaybarStore(
		context.Background(),
		&m3uetcpb.UnsubscribeFromPlaybarStoreRequest{
			SubscriptionId: id,
		},
	)
	onerror.Log(err)
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
