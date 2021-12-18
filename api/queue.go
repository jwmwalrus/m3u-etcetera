package api

import (
	"context"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
)

// QueueSvc defines the queue server
type QueueSvc struct {
	m3uetcpb.UnimplementedQueueSvcServer
}

// GetQueue implements m3uetcpb.QueueSvcServer
func (*QueueSvc) GetQueue(_ context.Context, req *m3uetcpb.GetQueueRequest) (*m3uetcpb.GetQueueResponse, error) {

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
	}

	res.Tracks = tList
	return res, nil
}

// ExecuteQueueAction implements m3uetcpb.QueueSvcServer
func (*QueueSvc) ExecuteQueueAction(_ context.Context, req *m3uetcpb.ExecuteQueueActionRequest) (*m3uetcpb.Empty, error) {

	q, _ := models.PerspectiveIndex(req.Perspective).GetPerspectiveQueue()

	go func() {
		switch req.Action {
		case m3uetcpb.QueueAction_Q_APPEND:
			q.Add(req.Locations, req.Ids)
		case m3uetcpb.QueueAction_Q_INSERT:
			q.InsertAt(int(req.Position), req.Locations, req.Ids)
		case m3uetcpb.QueueAction_Q_PREPPEND:
			q.InsertAt(0, req.Locations, req.Ids)
		case m3uetcpb.QueueAction_Q_DELETE:
			q.DeleteAt(int(req.Position))
		case m3uetcpb.QueueAction_Q_CLEAR:
			q.Clear()
		case m3uetcpb.QueueAction_Q_MOVE:
			q.MoveTo(int(req.FromPosition), int(req.Position))
		default:
		}
	}()

	return &m3uetcpb.Empty{}, nil
}
