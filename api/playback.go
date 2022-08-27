package api

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PlaybackSvc defines the playback service
type PlaybackSvc struct {
	m3uetcpb.UnimplementedPlaybackSvcServer
}

// GetPlayback implements m3uetcpb.PlaybackSvcServer
func (*PlaybackSvc) GetPlayback(_ context.Context,
	_ *empty.Empty) (*m3uetcpb.GetPlaybackResponse, error) {

	res := &m3uetcpb.GetPlaybackResponse{
		IsStreaming: playback.IsStreaming(),
		IsPlaying:   playback.IsPlaying(),
		IsPaused:    playback.IsPaused(),
		IsStopped:   playback.IsStopped(),
		IsReady:     playback.IsReady(),
	}
	pb, t := playback.GetPlayback()
	if pb != nil {
		res.Playback = pb.ToProtobuf().(*m3uetcpb.Playback)
		res.Track = &m3uetcpb.Track{}
		if t != nil {
			res.Track = t.ToProtobuf().(*m3uetcpb.Track)
		}
		return res, nil
	}

	return res, nil
}

// GetPlaybackList implements m3uetcpb.PlaybackSvcServer
func (*PlaybackSvc) GetPlaybackList(_ context.Context,
	_ *empty.Empty) (*m3uetcpb.GetPlaybackListResponse, error) {

	res := &m3uetcpb.GetPlaybackListResponse{}
	pbs := models.GetAllPlayback()
	for i := range pbs {
		res.PlaybackEntries = append(
			res.PlaybackEntries,
			pbs[i].ToProtobuf().(*m3uetcpb.Playback),
		)
	}

	return res, nil
}

// ExecutePlaybackAction implements m3uetcpb.PlaybackSvcServer
func (*PlaybackSvc) ExecutePlaybackAction(_ context.Context,
	req *m3uetcpb.ExecutePlaybackActionRequest) (*empty.Empty, error) {

	if req.Action == m3uetcpb.PlaybackAction_PB_PLAY {
		if len(req.Locations) > 0 {
			unsup := base.CheckUnsupportedFiles(req.Locations)
			if len(unsup) > 0 {
				return nil, status.Errorf(codes.InvalidArgument,
					"Unsupported locations were provided: %+q", unsup)
			}
		}
		if len(req.Ids) > 0 {
			_, notFound := models.FindTracksIn(req.Ids)
			if len(notFound) > 0 {
				return nil, status.Errorf(codes.InvalidArgument,
					"Non-existing track IDs were provided: %+v", notFound)
			}
		}
	}

	go func() {
		if !slices.Contains(
			[]m3uetcpb.PlaybackAction{
				m3uetcpb.PlaybackAction_PB_PLAY,
				m3uetcpb.PlaybackAction_PB_SEEK,
			},
			req.Action,
		) &&
			(len(req.Locations) > 0 || len(req.Ids) > 0) {

			for _, v := range req.Locations {
				log.Warnf("Ignoring given location: %v", v)
			}
			for _, v := range req.Ids {
				log.Warnf("Ignoring given ID: %v", v)
			}
		}

		switch req.Action {
		case m3uetcpb.PlaybackAction_PB_SEEK:
			playback.SeekInStream(req.Seek)
		case m3uetcpb.PlaybackAction_PB_NEXT:
			playback.NextStream()
		case m3uetcpb.PlaybackAction_PB_PAUSE:
			playback.PauseStream(false)
		case m3uetcpb.PlaybackAction_PB_PLAY:
			if len(req.Locations) > 0 || len(req.Ids) > 0 {
				if req.Force {
					playback.PlayStreams(req.Force, req.Locations, req.Ids)
				} else {
					q, _ := models.
						PerspectiveIndex(req.Perspective).
						GetPerspectiveQueue()

					q.Add(req.Locations, req.Ids)
				}
			} else {
				playback.PauseStream(true)
			}
		case m3uetcpb.PlaybackAction_PB_PREVIOUS:
			playback.PreviousStream()
		case m3uetcpb.PlaybackAction_PB_STOP:
			playback.StopAll()
		default:
			return
		}
	}()

	return &empty.Empty{}, nil
}

// SubscribeToPlayback implements m3uetcpb.PlaybackSvcServer
func (*PlaybackSvc) SubscribeToPlayback(_ *empty.Empty,
	stream m3uetcpb.PlaybackSvc_SubscribeToPlaybackServer) error {

	s, id := subscription.Subscribe(subscription.ToPlaybackEvent)
	defer func() { s.Unsubscribe() }()

	go func() {
		s.Event <- subscription.Event{Data: struct{}{}}
	}()

sLoop:
	for {
		select {
		case e := <-s.Event:
			if s.MustUnsubscribe(e) {
				break sLoop
			}

			res := &m3uetcpb.SubscribeToPlaybackResponse{
				SubscriptionId: id,
				IsStreaming:    playback.IsStreaming(),
				IsPlaying:      playback.IsPlaying(),
				IsPaused:       playback.IsPaused(),
				IsStopped:      playback.IsStopped(),
				IsReady:        playback.IsReady(),
			}

			pb, t := playback.GetPlayback()
			if pb != nil {
				res.Playback = pb.ToProtobuf().(*m3uetcpb.Playback)
				res.Track = &m3uetcpb.Track{}
				if t != nil {
					res.Track = t.ToProtobuf().(*m3uetcpb.Track)
				}
				err := stream.Send(res)
				if err != nil {
					log.Warn(err)
					return status.Errorf(codes.Internal,
						"Error sending playback: %v", err)
				}
				continue sLoop
			}
			err := stream.Send(res)
			if err != nil {
				return status.Errorf(codes.Internal,
					"Error sending playback: %v", err)
			}
		}
	}
	return nil
}

// UnsubscribeFromPlayback implements m3uetcpb.PlaybackSvcServer
func (*PlaybackSvc) UnsubscribeFromPlayback(_ context.Context,
	req *m3uetcpb.UnsubscribeFromPlaybackRequest) (*empty.Empty, error) {

	if req.SubscriptionId == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"A non-empty subscription ID is required")
	}
	subscription.Broadcast(
		subscription.ToNone,
		subscription.Event{Data: req.SubscriptionId},
	)

	return &empty.Empty{}, nil
}
