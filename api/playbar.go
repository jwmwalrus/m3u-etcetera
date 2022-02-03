package api

import (
	"context"
	"time"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/m3u-etcetera/pkg/impexp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// PlaybarSvc defines the playbar service
type PlaybarSvc struct {
	m3uetcpb.UnimplementedPlaybarSvcServer
}

// GetPlaybar implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) GetPlaybar(_ context.Context, req *m3uetcpb.GetPlaybarRequest) (*m3uetcpb.GetPlaybarResponse, error) {

	bar, err := models.PerspectiveIndex(req.Perspective).GetPlaybar()
	if err != nil {
		return nil,
			grpc.Errorf(codes.Internal,
				"Error while obtaining perspective playbar: %v",
				err)
	}

	pls := bar.GetAllOpenEntries()

	res := &m3uetcpb.GetPlaybarResponse{}

	list := []*m3uetcpb.Playlist{}
	for _, pl := range pls {
		out := pl.ToProtobuf().(*m3uetcpb.Playlist)
		list = append(list, out)
	}
	res.Playlists = list

	return res, nil
}

// GetPlaylist implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) GetPlaylist(_ context.Context, req *m3uetcpb.GetPlaylistRequest) (*m3uetcpb.GetPlaylistResponse, error) {
	if req.Id < 1 {
		return nil,
			grpc.Errorf(codes.InvalidArgument, "Playlist ID must be greater than zero")
	}

	pl := models.Playlist{}
	err := pl.Read(req.Id)
	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Playlist not found: %v", err)
	}

	res := &m3uetcpb.GetPlaylistResponse{
		Playlist: pl.ToProtobuf().(*m3uetcpb.Playlist),
	}

	pts := []*m3uetcpb.PlaylistTrack{}
	ts := []*m3uetcpb.Track{}
	ptlist, tlist := pl.GetTracks(int(req.Limit))
	for i := range ptlist {
		pts = append(pts, ptlist[i].ToProtobuf().(*m3uetcpb.PlaylistTrack))
		ts = append(ts, tlist[i].ToProtobuf().(*m3uetcpb.Track))
	}
	res.PlaylistTracks = pts
	res.Tracks = ts

	return res, nil
}

// GetAllPlaylists implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) GetAllPlaylists(_ context.Context, req *m3uetcpb.GetAllPlaylistsRequest) (*m3uetcpb.GetAllPlaylistsResponse, error) {
	bar, err := models.PerspectiveIndex(req.Perspective).GetPlaybar()
	if err != nil {
		return nil,
			grpc.Errorf(codes.Internal,
				"Error while obtaining perspective playbar: %v", err)
	}

	res := &m3uetcpb.GetAllPlaylistsResponse{}

	pls := bar.GetAllEntries(int(req.Limit))

	s := []*m3uetcpb.Playlist{}
	for i := range pls {
		s = append(s, pls[i].ToProtobuf().(*m3uetcpb.Playlist))
	}
	res.Playlists = s

	return res, nil
}

// GetPlaylistGroup implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) GetPlaylistGroup(_ context.Context, req *m3uetcpb.GetPlaylistGroupRequest) (*m3uetcpb.GetPlaylistGroupResponse, error) {
	if req.Id < 1 {
		return nil,
			grpc.Errorf(codes.InvalidArgument,
				"Playlist Group ID must be greater than zero")
	}
	pg := models.PlaylistGroup{}
	if err := pg.Read(req.Id); err != nil {
		return nil,
			grpc.Errorf(codes.NotFound, "Playlist Group not found: %v", err)
	}
	res := &m3uetcpb.GetPlaylistGroupResponse{
		PlaylistGroup: pg.ToProtobuf().(*m3uetcpb.PlaylistGroup),
	}
	return res, nil
}

// GetAllPlaylistGroups implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) GetAllPlaylistGroups(_ context.Context, req *m3uetcpb.GetAllPlaylistGroupsRequest) (*m3uetcpb.GetAllPlaylistGroupsResponse, error) {
	bar, err := models.PerspectiveIndex(req.Perspective).GetPlaybar()
	if err != nil {
		return nil,
			grpc.Errorf(codes.Internal,
				"Error while obtaining perspective playbar: %v", err)
	}

	res := &m3uetcpb.GetAllPlaylistGroupsResponse{}

	pgs := bar.GetAllGroups(int(req.Limit))

	s := []*m3uetcpb.PlaylistGroup{}
	for i := range pgs {
		s = append(s, pgs[i].ToProtobuf().(*m3uetcpb.PlaylistGroup))
	}
	res.PlaylistGroups = s

	return res, nil
}

// ExecutePlaybarAction implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) ExecutePlaybarAction(_ context.Context, req *m3uetcpb.ExecutePlaybarActionRequest) (*m3uetcpb.Empty, error) {
	if len(req.Ids) < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "At least one playlist ID is required")
	}

	if (req.Action == m3uetcpb.PlaybarAction_BAR_ACTIVATE ||
		req.Action == m3uetcpb.PlaybarAction_BAR_DEACTIVATE) &&
		len(req.Ids) != 1 {
		return nil,
			grpc.Errorf(codes.InvalidArgument,
				"Only one playlist ID must be provided to activate|deactivate")
	}

	pls, notfound := models.FindPlaylistsIn(req.Ids)
	if len(notfound) > 0 {
		return nil,
			grpc.Errorf(codes.InvalidArgument,
				"Some of the provided playlist IDs do not exist: %v", notfound)
	}

	go func() {
		switch req.Action {
		case m3uetcpb.PlaybarAction_BAR_OPEN:
			for i := range pls {
				bar := models.Playbar{}
				if err := bar.Read(pls[i].PlaybarID); err != nil {
					log.Error(err)
					continue
				}
				bar.OpenEntry(pls[i])
			}
		case m3uetcpb.PlaybarAction_BAR_ACTIVATE:
			bar := models.Playbar{}
			if err := bar.Read(pls[0].PlaybarID); err != nil {
				log.Error(err)
				return
			}
			playback.TryPlayingFromBar(pls[0], int(req.Position))
		case m3uetcpb.PlaybarAction_BAR_DEACTIVATE:
			bar := models.Playbar{}
			if err := bar.Read(pls[0].PlaybarID); err != nil {
				log.Error(err)
				return
			}
			playback.QuitPlayingFromBar(pls[0])
		case m3uetcpb.PlaybarAction_BAR_CLOSE:
			for i := range pls {
				bar := models.Playbar{}
				if err := bar.Read(pls[i].PlaybarID); err != nil {
					log.Error(err)
					continue
				}
				bar.CloseEntry(pls[i])
			}
		}
	}()

	return &m3uetcpb.Empty{}, nil
}

// ExecutePlaylistAction implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) ExecutePlaylistAction(_ context.Context, req *m3uetcpb.ExecutePlaylistActionRequest) (*m3uetcpb.ExecutePlaylistActionResponse, error) {

	if (req.Action == m3uetcpb.PlaylistAction_PL_UPDATE ||
		req.Action == m3uetcpb.PlaylistAction_PL_DESTROY) &&
		req.Id < 1 {
		return nil,
			grpc.Errorf(codes.InvalidArgument, "Invalid playlist ID: %v", req.Id)
	}

	pl := &models.Playlist{}
	if req.Action != m3uetcpb.PlaylistAction_PL_CREATE {
		if err := pl.Read(req.Id); err != nil {
			return nil,
				grpc.Errorf(codes.NotFound,
					"Playlist with ID=%v does not exist: %v", req.Id, err)
		}
	}

	if req.Action == m3uetcpb.PlaylistAction_PL_DESTROY && pl.Open {
		return nil,
			grpc.Errorf(codes.InvalidArgument,
				"A playlist cannot be deleted while open")
	}

	switch req.Action {
	case m3uetcpb.PlaylistAction_PL_CREATE:
		bar, err := models.PerspectiveIndex(req.Perspective).GetPlaybar()
		if err != nil {
			return nil,
				grpc.Errorf(codes.Internal,
					"Error obtaining perspective:", err)
		}

		pl, err = bar.CreateEntry(req.Name, req.Description)
		if err != nil {
			return nil,
				grpc.Errorf(codes.Internal,
					"Error creating playlist:", err)
		}
	case m3uetcpb.PlaylistAction_PL_UPDATE:
		bar := pl.Playbar
		err := bar.UpdateEntry(
			pl,
			req.Name,
			req.Description,
			req.PlaylistGroupId,
			req.ResetDescription,
		)
		if err != nil {
			return nil,
				grpc.Errorf(codes.Internal,
					"Error updating playlist:", err)
		}
	case m3uetcpb.PlaylistAction_PL_DESTROY:
		bar := pl.Playbar
		err := bar.DestroyEntry(pl)
		if err != nil {
			return nil,
				grpc.Errorf(codes.Internal,
					"Error deleting playlist:", err)
		}
	}

	return &m3uetcpb.ExecutePlaylistActionResponse{Id: pl.ID}, nil
}

// ExecutePlaylistGroupAction implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) ExecutePlaylistGroupAction(_ context.Context, req *m3uetcpb.ExecutePlaylistGroupActionRequest) (*m3uetcpb.ExecutePlaylistGroupActionResponse, error) {

	if req.Action != m3uetcpb.PlaylistGroupAction_PG_CREATE &&
		req.Id < 1 {
		return nil,
			grpc.Errorf(codes.InvalidArgument, "Invalid playlist group ID")
	}

	pg := &models.PlaylistGroup{}
	bar := &models.Playbar{}
	if req.Action != m3uetcpb.PlaylistGroupAction_PG_CREATE {
		var err error
		if err = pg.Read(req.Id); err != nil {
			return nil,
				grpc.Errorf(codes.NotFound,
					"Playlist Group with ID=%v does not exist: %v",
					req.Id, err)
		}
		bar, err = models.PerspectiveIndex(pg.Perspective.Idx).GetPlaybar()
		if err != nil {
			return nil,
				grpc.Errorf(codes.Internal, "Error obtaining perspective:", err)
		}
	}

	switch req.Action {
	case m3uetcpb.PlaylistGroupAction_PG_CREATE:
		bar, err := models.PerspectiveIndex(req.Perspective).GetPlaybar()
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "Error obtaining perspective:", err)
		}
		pg, err = bar.CreateGroup(req.Name, req.Description)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "Error creating playlist group:", err)
		}
	case m3uetcpb.PlaylistGroupAction_PG_UPDATE:
		err := bar.UpdateGroup(pg, req.Name, req.Description, req.ResetDescription)
		if err != nil {
			return nil,
				grpc.Errorf(codes.Internal, "Error updating playlist group:", err)
		}
	case m3uetcpb.PlaylistGroupAction_PG_DESTROY:
		err := bar.DestroyGroup(pg)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "Error deleting playlist group:", err)
		}
	}
	return &m3uetcpb.ExecutePlaylistGroupActionResponse{Id: pg.ID}, nil
}

// ExecutePlaylistTrackAction implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) ExecutePlaylistTrackAction(_ context.Context, req *m3uetcpb.ExecutePlaylistTrackActionRequest) (*m3uetcpb.Empty, error) {
	if req.PlaylistId < 1 {
		return nil,
			grpc.Errorf(codes.InvalidArgument,
				"Invalid playlist ID: %v", req.PlaylistId)
	}

	pl := &models.Playlist{}
	if err := pl.Read(req.PlaylistId); err != nil {
		return nil,
			grpc.Errorf(codes.NotFound,
				"Playlist with ID=%v does not exist: %v", req.PlaylistId, err)
	}

	if !pl.Open {
		return nil,
			grpc.Errorf(codes.InvalidArgument,
				"Playlist must be open to operate on it")
	}

	go func() {
		bar := pl.Playbar
		switch req.Action {
		case m3uetcpb.PlaylistTrackAction_PT_APPEND:
			bar.AppendToPlaylist(pl, req.TrackIds, req.Locations)
		case m3uetcpb.PlaylistTrackAction_PT_PREPEND:
			bar.PrependToPlaylist(pl, req.TrackIds, req.Locations)
		case m3uetcpb.PlaylistTrackAction_PT_INSERT:
			bar.InsertIntoPlaylist(pl, int(req.Position), req.TrackIds, req.Locations)
		case m3uetcpb.PlaylistTrackAction_PT_DELETE:
			bar.DeleteFromPlaylist(pl, int(req.Position))
		case m3uetcpb.PlaylistTrackAction_PT_CLEAR:
			bar.ClearPlaylist(pl)
		case m3uetcpb.PlaylistTrackAction_PT_MOVE:
			bar.MovePlaylistTrack(pl, int(req.Position), int(req.FromPosition))
		}
	}()

	return &m3uetcpb.Empty{}, nil
}

// ImportPlaylists implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) ImportPlaylists(req *m3uetcpb.ImportPlaylistsRequest, stream m3uetcpb.PlaybarSvc_ImportPlaylistsServer) error {
	bar, err := models.PerspectiveIndex(req.Perspective).GetPlaybar()
	if err != nil {
		return grpc.Errorf(codes.Internal,
			"Error while obtaining perspective playbar: %v", err)
	}

	if len(req.Locations) < 1 {
		return grpc.Errorf(codes.InvalidArgument,
			"At least one playlist location is required")
	}

	for _, l := range req.Locations {
		pl, msgs, err := bar.ImportPlaylist(l)
		if err != nil {
			un, err2 := urlstr.URLToPath(l)
			if err2 != nil {
				un = l
			}
			return grpc.Errorf(codes.InvalidArgument,
				"Error importing playlist at `%v`: %v", un, err)
		}
		req := &m3uetcpb.ImportPlaylistsResponse{Id: pl.ID, ImportErrors: msgs}
		err = stream.Send(req)
		if err != nil {
			return grpc.Errorf(codes.Internal,
				"Error sending stream: %v", err)
		}
	}

	return nil
}

// ExportPlaylist implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) ExportPlaylist(_ context.Context, req *m3uetcpb.ExportPlaylistRequest) (*m3uetcpb.Empty, error) {
	if req.Id < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument,
			"A valid playlist ID is required")
	}

	pl := models.Playlist{}
	if err := pl.Read(req.Id); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument,
			"Playlist does not exist: %v", err)
	}

	if req.Location == "" {
		return nil, grpc.Errorf(codes.InvalidArgument,
			"The target location for the playlist is required")
	}

	var format impexp.PlaylistType
	switch req.Format {
	case m3uetcpb.PlaylistExportFormat_PLEF_M3U:
		format = impexp.M3UPlaylist
	case m3uetcpb.PlaylistExportFormat_PLEF_PLS:
		format = impexp.PLSPlaylist
	default:
		return nil, grpc.Errorf(codes.InvalidArgument,
			"Unsupported export format: %v", req.Format)
	}

	if err := pl.Export(format, req.Location); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Error exporting playlist: %v", err)
	}
	return &m3uetcpb.Empty{}, nil
}

// SubscribeToPlaybarStore implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) SubscribeToPlaybarStore(_ *m3uetcpb.Empty, stream m3uetcpb.PlaybarSvc_SubscribeToPlaybarStoreServer) error {

	s, id := subscription.Subscribe(subscription.ToPlaybarStoreEvent)
	defer func() { s.Unsubscribe() }()

	go func() {
		time.Sleep(2 * time.Second)
		s.Event <- subscription.Event{Idx: int(models.PlaybarEventInitial)}
	}()

	sendPlaylistGroup := func(e m3uetcpb.PlaybarEvent, c models.ProtoOut) error {
		res := &m3uetcpb.SubscribeToPlaybarStoreResponse{
			SubscriptionId:   id,
			Event:            e,
			ActivePlaylistId: models.GetActiveEntry().ID,
			Item: &m3uetcpb.SubscribeToPlaybarStoreResponse_PlaylistGroup{
				PlaylistGroup: c.ToProtobuf().(*m3uetcpb.PlaylistGroup),
			},
		}
		return stream.Send(res)
	}

	sendPlaylist := func(e m3uetcpb.PlaybarEvent, c models.ProtoOut) error {
		res := &m3uetcpb.SubscribeToPlaybarStoreResponse{
			SubscriptionId:   id,
			Event:            e,
			ActivePlaylistId: models.GetActiveEntry().ID,
			Item: &m3uetcpb.SubscribeToPlaybarStoreResponse_Playlist{
				Playlist: c.ToProtobuf().(*m3uetcpb.Playlist),
			},
		}
		return stream.Send(res)
	}

	sendOpenPlaylist := func(e m3uetcpb.PlaybarEvent, c models.ProtoOut) error {
		res := &m3uetcpb.SubscribeToPlaybarStoreResponse{
			SubscriptionId:   id,
			Event:            e,
			ActivePlaylistId: models.GetActiveEntry().ID,
			Item: &m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylist{
				OpenPlaylist: c.ToProtobuf().(*m3uetcpb.Playlist),
			},
		}
		return stream.Send(res)
	}

	sendOpenPlaylistTrack := func(e m3uetcpb.PlaybarEvent, t models.ProtoOut) error {
		res := &m3uetcpb.SubscribeToPlaybarStoreResponse{
			SubscriptionId:   id,
			Event:            e,
			ActivePlaylistId: models.GetActiveEntry().ID,
			Item: &m3uetcpb.SubscribeToPlaybarStoreResponse_OpenPlaylistTrack{
				OpenPlaylistTrack: t.ToProtobuf().(*m3uetcpb.PlaylistTrack),
			},
		}
		return stream.Send(res)
	}

	sendOpenTrack := func(e m3uetcpb.PlaybarEvent, t models.ProtoOut) error {
		res := &m3uetcpb.SubscribeToPlaybarStoreResponse{
			SubscriptionId:   id,
			Event:            e,
			ActivePlaylistId: models.GetActiveEntry().ID,
			Item: &m3uetcpb.SubscribeToPlaybarStoreResponse_OpenTrack{
				OpenTrack: t.ToProtobuf().(*m3uetcpb.Track),
			},
		}
		return stream.Send(res)
	}

sLoop:
	for {

		select {
		case e := <-s.Event:
			if s.MustUnsubscribe(e) {
				break sLoop
			}

			if models.PlaybarEvent(e.Idx) == models.PlaybarEventInitial {
				err := stream.Send(
					&m3uetcpb.SubscribeToPlaybarStoreResponse{
						SubscriptionId: id,
						Event:          m3uetcpb.PlaybarEvent_BE_INITIAL,
					},
				)
				if err != nil {
					return err
				}

				pgs, pls, opls, opts, ots := models.GetPlaybarStore()

				for i := range pgs {
					err := sendPlaylistGroup(
						m3uetcpb.PlaybarEvent_BE_INITIAL_ITEM,
						pgs[i],
					)
					if err != nil {
						return err
					}
				}

				for i := range pls {
					err := sendPlaylist(
						m3uetcpb.PlaybarEvent_BE_INITIAL_ITEM,
						pls[i],
					)
					if err != nil {
						return err
					}
				}

				for i := range opls {
					err := sendOpenPlaylist(
						m3uetcpb.PlaybarEvent_BE_INITIAL_ITEM,
						opls[i],
					)
					if err != nil {
						return err
					}
				}

				for i := range opts {
					err := sendOpenPlaylistTrack(
						m3uetcpb.PlaybarEvent_BE_INITIAL_ITEM,
						opts[i],
					)
					if err != nil {
						return err
					}
				}

				for i := range ots {
					err := sendOpenTrack(
						m3uetcpb.PlaybarEvent_BE_INITIAL_ITEM,
						ots[i],
					)
					if err != nil {
						return err
					}
				}

				err = stream.Send(
					&m3uetcpb.SubscribeToPlaybarStoreResponse{
						SubscriptionId: id,
						Event:          m3uetcpb.PlaybarEvent_BE_INITIAL_DONE,
					},
				)
				if err != nil {
					return err
				}
				continue sLoop
			}

			if models.PlaybarEvent(e.Idx) == models.PlaybarEventOpenItems {
				pl := &models.Playlist{}
				err := pl.Read(e.Data.(int64))
				if err != nil {
					return err
				}

				err = stream.Send(
					&m3uetcpb.SubscribeToPlaybarStoreResponse{
						SubscriptionId: id,
						Event:          m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS,
					},
				)
				if err != nil {
					return err
				}

				err = sendOpenPlaylist(
					m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS_ITEM,
					pl,
				)
				if err != nil {
					return err
				}

				if pl.Open {
					opts, ots := pl.GetTracks(0)

					for i := range opts {
						err := sendOpenPlaylistTrack(
							m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS_ITEM,
							opts[i],
						)
						if err != nil {
							return err
						}
					}

					for i := range ots {
						err := sendOpenTrack(
							m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS_ITEM,
							ots[i],
						)
						if err != nil {
							return err
						}
					}
				}

				err = stream.Send(
					&m3uetcpb.SubscribeToPlaybarStoreResponse{
						SubscriptionId: id,
						Event:          m3uetcpb.PlaybarEvent_BE_OPEN_ITEMS_DONE,
					},
				)
				if err != nil {
					return err
				}
				continue sLoop
			}

			var eout m3uetcpb.PlaybarEvent
			var fn func(m3uetcpb.PlaybarEvent, models.ProtoOut) error

			switch models.PlaybarEvent(e.Idx) {
			case models.PlaybarEventItemAdded:
				eout = m3uetcpb.PlaybarEvent_BE_ITEM_ADDED
			case models.PlaybarEventItemChanged:
				eout = m3uetcpb.PlaybarEvent_BE_ITEM_CHANGED
			case models.PlaybarEventItemRemoved:
				eout = m3uetcpb.PlaybarEvent_BE_ITEM_REMOVED
			default:
				log.Errorf("Ignoring unsupported playbar event: %v", e.Idx)
				continue sLoop

			}

			switch e.Data.(type) {
			case *models.PlaylistGroup:
				fn = sendPlaylistGroup
			case *models.Playlist:
				fn = sendPlaylist
			default:
				log.Errorf("Ignoring unsupported data for %v: %+v", e.Idx, e.Data)
				continue sLoop
			}

			if err := fn(eout, e.Data.(models.ProtoOut)); err != nil {
				return err
			}
		}
	}
	return nil
}

// UnsubscribeFromPlaybarStore implements m3uetcpb.PlaybarSvcServer
func (*PlaybarSvc) UnsubscribeFromPlaybarStore(_ context.Context, req *m3uetcpb.UnsubscribeFromPlaybarStoreRequest) (*m3uetcpb.Empty, error) {
	if req.SubscriptionId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "A non-empty subscription ID is required")
	}
	subscription.Broadcast(
		subscription.ToNone,
		subscription.Event{Data: req.SubscriptionId},
	)

	return &m3uetcpb.Empty{}, nil
}
