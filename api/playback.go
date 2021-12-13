package api

import (
	"context"

	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
	log "github.com/sirupsen/logrus"
)

type PlaybackSvc struct {
	m3uetcpb.UnimplementedPlaybackSvcServer
}

func (*PlaybackSvc) Get(_ context.Context, _ *m3uetcpb.Empty) (*m3uetcpb.GetResponse, error) {

	pb := playback.GetPlayback()
	if pb != nil {
		if pb.TrackID > 0 {
			if t, err := pb.GetTrack(); err == nil {
				res := &m3uetcpb.GetResponse_Track{Track: t.ToProtobuf()}
				return &m3uetcpb.GetResponse{Playing: res}, nil
			}
		}
		res := &m3uetcpb.GetResponse_Playback{Playback: pb.ToProtobuf()}
		return &m3uetcpb.GetResponse{Playing: res}, nil
	}

	res := &m3uetcpb.GetResponse_Empty{Empty: &m3uetcpb.Empty{}}
	return &m3uetcpb.GetResponse{Playing: res}, nil
}

func (*PlaybackSvc) ExecuteAction(_ context.Context, req *m3uetcpb.ExecuteActionRequest) (*m3uetcpb.Empty, error) {

	go func() {
		go func() {
			if !slice.Contains([]m3uetcpb.PlaybackAction{m3uetcpb.PlaybackAction_PLAY, m3uetcpb.PlaybackAction_SEEK}, req.Action) &&
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
		case m3uetcpb.PlaybackAction_SEEK:
			playback.SeekInStream(req.Seek)
		case m3uetcpb.PlaybackAction_NEXT:
			playback.NextStream()
		case m3uetcpb.PlaybackAction_PAUSE:
			playback.PauseStream(false)
		case m3uetcpb.PlaybackAction_PLAY:
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
		case m3uetcpb.PlaybackAction_PREVIOUS:
			playback.PreviousStream()
		case m3uetcpb.PlaybackAction_STOP:
			playback.StopAll()
		default:
		}
	}()

	return &m3uetcpb.Empty{}, nil
}
