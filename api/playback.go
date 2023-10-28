package api

import (
	"context"
	"log/slog"
	"slices"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PlaybackSvc implements the m3uetcpb.PlaybackSvcServer interface.
type PlaybackSvc struct {
	PbEvents playback.IEvents
	m3uetcpb.UnimplementedPlaybackSvcServer
}

func (svc *PlaybackSvc) GetPlayback(_ context.Context,
	_ *m3uetcpb.Empty) (*m3uetcpb.GetPlaybackResponse, error) {

	res := &m3uetcpb.GetPlaybackResponse{
		IsStreaming: svc.PbEvents.IsStreaming(),
		IsPlaying:   svc.PbEvents.IsPlaying(),
		IsPaused:    svc.PbEvents.IsPaused(),
		IsStopped:   svc.PbEvents.IsStopped(),
		IsReady:     svc.PbEvents.IsReady(),
	}
	pb, t := svc.PbEvents.GetPlayback()
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

func (*PlaybackSvc) GetPlaybackList(_ context.Context,
	_ *m3uetcpb.Empty) (*m3uetcpb.GetPlaybackListResponse, error) {

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

func (svc *PlaybackSvc) ExecutePlaybackAction(_ context.Context,
	req *m3uetcpb.ExecutePlaybackActionRequest) (*m3uetcpb.Empty, error) {

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
				slog.Warn("Ignoring given location", "location", v)
			}
			for _, v := range req.Ids {
				slog.Warn("Ignoring given ID", "ID", v)
			}
		}

		switch req.Action {
		case m3uetcpb.PlaybackAction_PB_SEEK:
			svc.PbEvents.SeekInStream(req.Seek)
		case m3uetcpb.PlaybackAction_PB_NEXT:
			svc.PbEvents.NextStream()
		case m3uetcpb.PlaybackAction_PB_PAUSE:
			svc.PbEvents.PauseStream(false)
		case m3uetcpb.PlaybackAction_PB_PLAY:
			if len(req.Locations) > 0 || len(req.Ids) > 0 {
				if req.Force {
					svc.PbEvents.PlayStreams(req.Force, req.Locations, req.Ids)
				} else {
					q, _ := models.
						PerspectiveIndex(req.Perspective).
						GetPerspectiveQueue()

					q.Add(req.Locations, req.Ids)
				}
			} else {
				svc.PbEvents.PauseStream(true)
			}
		case m3uetcpb.PlaybackAction_PB_PREVIOUS:
			svc.PbEvents.PreviousStream()
		case m3uetcpb.PlaybackAction_PB_STOP:
			svc.PbEvents.StopAll()
		default:
			return
		}
	}()

	return &m3uetcpb.Empty{}, nil
}

func (svc *PlaybackSvc) SubscribeToPlayback(_ *m3uetcpb.Empty,
	stream m3uetcpb.PlaybackSvc_SubscribeToPlaybackServer) error {

	sub, id := subscription.Subscribe(subscription.ToPlaybackEvent)
	defer func() { sub.Unsubscribe() }()

	go func() {
		sub.Event <- subscription.Event{Data: struct{}{}}
	}()

sLoop:
	for {
		select {
		case e := <-sub.Event:
			if sub.MustUnsubscribe(e) {
				break sLoop
			}

			res := &m3uetcpb.SubscribeToPlaybackResponse{
				SubscriptionId: id,
				IsStreaming:    svc.PbEvents.IsStreaming(),
				IsPlaying:      svc.PbEvents.IsPlaying(),
				IsPaused:       svc.PbEvents.IsPaused(),
				IsStopped:      svc.PbEvents.IsStopped(),
				IsReady:        svc.PbEvents.IsReady(),
			}

			pb, t := svc.PbEvents.GetPlayback()
			if pb != nil {
				res.Playback = pb.ToProtobuf().(*m3uetcpb.Playback)
				res.Track = &m3uetcpb.Track{}
				if t != nil {
					res.Track = t.ToProtobuf().(*m3uetcpb.Track)
				}
				err := stream.Send(res)
				if err != nil {
					slog.Warn("Failed to send stream", "error", err)
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

func (*PlaybackSvc) UnsubscribeFromPlayback(_ context.Context,
	req *m3uetcpb.UnsubscribeFromPlaybackRequest) (*m3uetcpb.Empty, error) {

	if req.SubscriptionId == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"A non-empty subscription ID is required")
	}
	subscription.Broadcast(
		subscription.ToNone,
		subscription.Event{Data: req.SubscriptionId},
	)

	return &m3uetcpb.Empty{}, nil
}
