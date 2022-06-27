package api

import (
	"context"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueueSvc defines the queue server
type QueueSvc struct {
	m3uetcpb.UnimplementedQueueSvcServer
}

// GetQueue implements m3uetcpb.QueueSvcServer
func (*QueueSvc) GetQueue(_ context.Context,
	req *m3uetcpb.GetQueueRequest) (*m3uetcpb.GetQueueResponse, error) {

	res := &m3uetcpb.GetQueueResponse{}
	qs, ts := models.GetAllQueueTracks(
		models.PerspectiveIndex(req.Perspective),
		int(req.Limit),
	)

	qtList := []*m3uetcpb.QueueTrack{}
	for _, qt := range qs {
		out := qt.ToProtobuf().(*m3uetcpb.QueueTrack)
		qtList = append(qtList, out)
	}
	res.QueueTracks = qtList

	tList := []*m3uetcpb.Track{}
	for _, t := range ts {
		out := t.ToProtobuf().(*m3uetcpb.Track)
		tList = append(tList, out)
		res.Duration += t.Duration
	}

	res.Tracks = tList
	return res, nil
}

// ExecuteQueueAction implements m3uetcpb.QueueSvcServer
func (*QueueSvc) ExecuteQueueAction(_ context.Context,
	req *m3uetcpb.ExecuteQueueActionRequest) (*m3uetcpb.Empty, error) {

	if slices.Contains(
		[]m3uetcpb.QueueAction{
			m3uetcpb.QueueAction_Q_APPEND,
			m3uetcpb.QueueAction_Q_PREPEND,
		},
		req.Action,
	) {
		if len(req.Locations) > 0 || len(req.Ids) > 0 {
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

	q, _ := models.PerspectiveIndex(req.Perspective).GetPerspectiveQueue()

	go func() {
		switch req.Action {
		case m3uetcpb.QueueAction_Q_APPEND:
			q.Add(req.Locations, req.Ids)
		case m3uetcpb.QueueAction_Q_INSERT:
			q.InsertAt(int(req.Position), req.Locations, req.Ids)
		case m3uetcpb.QueueAction_Q_PREPEND:
			q.InsertAt(0, req.Locations, req.Ids)
		case m3uetcpb.QueueAction_Q_DELETE:
			q.DeleteAt(int(req.Position))
		case m3uetcpb.QueueAction_Q_CLEAR:
			q.Clear()
		case m3uetcpb.QueueAction_Q_MOVE:
			q.MoveTo(int(req.Position), int(req.FromPosition))
		default:
		}
	}()

	return &m3uetcpb.Empty{}, nil
}

// SubscribeToQueueStore implements m3uetcpb.QueueSvcServer
func (*QueueSvc) SubscribeToQueueStore(_ *m3uetcpb.Empty,
	stream m3uetcpb.QueueSvc_SubscribeToQueueStoreServer) error {

	s, id := subscription.Subscribe(subscription.ToQueueStoreEvent)
	defer func() { s.Unsubscribe() }()

	go func() {
		time.Sleep(2 * time.Second)
		s.Event <- subscription.Event{}
	}()

sLoop:
	for {
		select {
		case e := <-s.Event:
			if s.MustUnsubscribe(e) {
				break sLoop
			}

			res := &m3uetcpb.SubscribeToQueueStoreResponse{SubscriptionId: id}
			qs, ts, pds := models.GetQueueStore()

			qtList := []*m3uetcpb.QueueTrack{}
			for _, qt := range qs {
				out := qt.ToProtobuf().(*m3uetcpb.QueueTrack)
				qtList = append(qtList, out)
			}
			res.QueueTracks = qtList

			tList := []*m3uetcpb.Track{}
			for _, t := range ts {
				out := t.ToProtobuf().(*m3uetcpb.Track)
				tList = append(tList, out)
			}
			res.Tracks = tList

			pdList := []*m3uetcpb.PerspectiveDigest{}
			for _, pd := range pds {
				out := pd.ToProtobuf().(*m3uetcpb.PerspectiveDigest)
				pdList = append(pdList, out)
			}
			res.Digest = pdList

			if err := stream.Send(res); err != nil {
				return status.Errorf(codes.Internal,
					"Error sending queue event: %v",
					err)
			}
		}
	}

	return nil
}

// UnsubscribeFromQueueStore implements m3uetcpb.QueueSvcServer
func (*QueueSvc) UnsubscribeFromQueueStore(_ context.Context,
	req *m3uetcpb.UnsubscribeFromQueueStoreRequest) (*m3uetcpb.Empty, error) {

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
