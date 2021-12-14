package api

import (
	"context"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
)

type QueueSvc struct {
	m3uetcpb.UnimplementedQueueSvcServer
}

func (*QueueSvc) GetQueue(_ context.Context, req *m3uetcpb.GetQueueRequest) (*m3uetcpb.GetQueueResponse, error) {

	res := &m3uetcpb.GetQueueResponse{}
	s := models.GetAllQueueTracks(
		models.PerspectiveIndex(req.Perspective),
		int(req.Limit),
	)
	list := []*m3uetcpb.QueueTrack{}
	for _, qt := range s {
		list = append(list, qt.ToProtobuf())
	}

	res.QueueTracks = list
	return res, nil
}

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
		default:
		}
	}()

	return &m3uetcpb.Empty{}, nil
}
