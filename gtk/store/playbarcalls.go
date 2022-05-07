package store

import (
	"context"
	"fmt"
	"io"

	"github.com/gotk3/gotk3/glib"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/status"
)

// ApplyPlaylistGroupChanges applies collection changes
func ApplyPlaylistGroupChanges() {
	log.Info("Applying playlist-group changes")

	cc, err := getClientConn1()
	if err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewPlaybarSvcClient(cc)

	model := playlistGroupsModel

	iter, ok := model.GetIterFirst()
	for ok {
		row, err := GetListStoreModelValues(
			model,
			iter,
			[]ModelColumn{PGColPlaylistGroupID, PGColName, PGColDescription},
		)
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

// ExecutePlaybarAction -
func ExecutePlaybarAction(req *m3uetcpb.ExecutePlaybarActionRequest) (err error) {
	cc, err := getClientConn1()
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
func ExecutePlaylistAction(req *m3uetcpb.ExecutePlaylistActionRequest) (
	*m3uetcpb.ExecutePlaylistActionResponse, error) {

	cc, err := getClientConn1()
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
func ExecutePlaylistGroupAction(req *m3uetcpb.ExecutePlaylistGroupActionRequest) (
	*m3uetcpb.ExecutePlaylistGroupActionResponse, error) {

	cc, err := getClientConn1()
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
	cc, err := getClientConn1()
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

// ImportPlaylists -
func ImportPlaylists(req *m3uetcpb.ImportPlaylistsRequest) (
	msgList []string, err error) {

	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybarSvcClient(cc)
	stream, err := cl.ImportPlaylists(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		log.Error(s.Message())
		return
	}

	for {
		var res *m3uetcpb.ImportPlaylistsResponse
		res, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		for _, msg := range res.ImportErrors {
			msgList = append(msgList, msg)
		}
	}
	return
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
				slices.Contains(
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

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
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
