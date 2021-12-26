package api

import (
	"context"
	"time"

	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// PlaybackSvc defines the playback service
type PlaybackSvc struct {
	m3uetcpb.UnimplementedPlaybackSvcServer
}

// GetPlayback implements m3uetcpb.PlaybackSvcServer
func (*PlaybackSvc) GetPlayback(_ context.Context, _ *m3uetcpb.Empty) (*m3uetcpb.GetPlaybackResponse, error) {

	pb := playback.GetPlayback()
	if pb != nil {
		res := &m3uetcpb.GetPlaybackResponse{Playing: true}
		res.Playback = pb.ToProtobuf().(*m3uetcpb.Playback)
		res.Track = &m3uetcpb.Track{}
		if pb.TrackID > 0 {
			if t, err := pb.GetTrack(); err == nil {
				res.Track = t.ToProtobuf().(*m3uetcpb.Track)
			}
		} else {
			if t, err := models.ReadTagsForLocation(pb.Location); err == nil {
				res.Track = t.ToProtobuf().(*m3uetcpb.Track)
			}
		}
		return res, nil
	}

	return &m3uetcpb.GetPlaybackResponse{Playing: false}, nil
}

// ExecutePlaybackAction implements m3uetcpb.PlaybackSvcServer
func (*PlaybackSvc) ExecutePlaybackAction(_ context.Context, req *m3uetcpb.ExecutePlaybackActionRequest) (*m3uetcpb.Empty, error) {
	if req.Action == m3uetcpb.PlaybackAction_PB_PLAY {
		if len(req.Locations) > 0 {
			unsup := models.CheckUnsupportedFiles(req.Locations)
			if len(unsup) > 0 {
				return nil, grpc.Errorf(codes.InvalidArgument, "Unsupported locations were provided: %+q", unsup)
			}
		}
		if len(req.Ids) > 0 {
			_, notFound := models.FindTracksIn(req.Ids)
			if len(notFound) > 0 {
				return nil, grpc.Errorf(codes.InvalidArgument, "Non-existing track IDs were provided: %+v", notFound)
			}
		}
	}

	go func() {
		go func() {
			if !slice.Contains([]m3uetcpb.PlaybackAction{m3uetcpb.PlaybackAction_PB_PLAY, m3uetcpb.PlaybackAction_PB_SEEK}, req.Action) &&
				(len(req.Locations) > 0 || len(req.Ids) > 0) {
				for _, v := range req.Locations {
					log.Warnf("Ignoring given location: %v", v)
				}
				for _, v := range req.Ids {
					log.Warnf("Ignoring given ID: %v", v)
				}
			}
		}()

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
					q, _ := models.PerspectiveIndex(req.Perspective).GetPerspectiveQueue()
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
		}
	}()

	return &m3uetcpb.Empty{}, nil
}

// SubscribeToPlayback implements m3uetcpb.PlaybackSvcServer
func (*PlaybackSvc) SubscribeToPlayback(_ *m3uetcpb.Empty, stream m3uetcpb.PlaybackSvc_SubscribeToPlaybackServer) error {

	base.GetBusy(base.IdleStatusSubscription)
	defer func() { base.GetFree(base.IdleStatusSubscription) }()

	subs := playback.Subscribe(playback.SubscribedToPlayback)
	defer func() { subs.Unsubscribe() }()

	go func() { time.Sleep(5 * time.Second); subs.Data <- playback.GetPlayback() }()

subsLoop:
	for {
		select {
		case i := <-subs.Data:
			pb, ok := i.(*models.Playback)
			if !ok {
				pb = nil
			}
			if pb != nil {
				res := &m3uetcpb.SubscribeToPlaybackResponse{Playing: true}
				res.Playback = pb.ToProtobuf().(*m3uetcpb.Playback)
				res.Track = &m3uetcpb.Track{}
				if pb.TrackID > 0 {
					if t, err := pb.GetTrack(); err == nil {
						res.Track = t.ToProtobuf().(*m3uetcpb.Track)
					}
				} else {
					if t, err := models.ReadTagsForLocation(pb.Location); err == nil {
						res.Track = t.ToProtobuf().(*m3uetcpb.Track)
					}
				}
				err := stream.Send(res)
				if err != nil {
					log.Warn(err)
					break subsLoop
				}
				continue subsLoop
			}
			res := &m3uetcpb.SubscribeToPlaybackResponse{Playing: false}
			err := stream.Send(res)
			if err != nil {
				log.Warn(err)
				break subsLoop
			}
		}
	}
	return nil
}
