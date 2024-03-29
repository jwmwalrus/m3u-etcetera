package api

import (
	"context"
	"log/slog"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/m3u-etcetera/pkg/qparams"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QuerySvc implemets the m3uetcpb.QuerySvcServer interface.
type QuerySvc struct {
	m3uetcpb.UnimplementedQuerySvcServer
}

func (*QuerySvc) GetQuery(_ context.Context,
	req *m3uetcpb.GetQueryRequest) (*m3uetcpb.GetQueryResponse, error) {

	if req.Id < 1 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Query id must be greater than zero")
	}

	qy := models.Query{}
	if err := qy.Read(req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}

	out := qy.ToProtobuf().(*m3uetcpb.Query)
	return &m3uetcpb.GetQueryResponse{Query: out}, nil
}

func (*QuerySvc) GetQueries(_ context.Context,
	req *m3uetcpb.GetQueriesRequest) (*m3uetcpb.GetQueriesResponse, error) {

	qybs := models.FilterCollectionQueryBoundaries(req.CollectionIds)
	qys := models.GetAllQueries(int(req.Limit), qybs...)

	out := []*m3uetcpb.Query{}
	for _, x := range qys {
		aux := x.ToProtobuf().(*m3uetcpb.Query)
		out = append(out, aux)
	}

	return &m3uetcpb.GetQueriesResponse{Queries: out}, nil
}

func (*QuerySvc) AddQuery(_ context.Context,
	req *m3uetcpb.AddQueryRequest) (*m3uetcpb.AddQueryResponse, error) {

	if req.Query.Params != "" {

		if _, err := qparams.ParseParams(req.Query.Params); err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"Error parsing query params: %v", err)
		}
	}

	qy := models.FromProtobuf(req.Query)

	qybs := models.CollectionsToBoundaries(
		models.CreateCollectionQueries(req.Query.CollectionIds),
	)
	if err := qy.SaveBound(qybs); err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"Error saving query: %v", err)
	}

	return &m3uetcpb.AddQueryResponse{Id: qy.ID}, nil
}

func (*QuerySvc) UpdateQuery(_ context.Context,
	req *m3uetcpb.UpdateQueryRequest) (*m3uetcpb.Empty, error) {

	qy := models.Query{}
	if err := qy.Read(req.Query.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}
	if qy.IsReadOnly() {
		return nil, status.Errorf(codes.InvalidArgument, "Query is read-only")
	}

	qy.FromProtobuf(req.Query)

	if err := models.DeleteCollectionQueries(qy.ID); err != nil {
		return nil, status.Errorf(codes.Internal,
			"Error replacing collection boundaries: %v", err)
	}

	qybs := models.CollectionsToBoundaries(
		models.CreateCollectionQueries(req.Query.CollectionIds),
	)
	if err := qy.SaveBound(qybs); err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"Error saving query: %v", err)
	}

	return &m3uetcpb.Empty{}, nil
}

func (*QuerySvc) RemoveQuery(_ context.Context,
	req *m3uetcpb.RemoveQueryRequest) (*m3uetcpb.Empty, error) {
	qy := models.Query{}
	if err := qy.Read(req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}
	if qy.IsReadOnly() {
		return nil, status.Errorf(codes.InvalidArgument, "Query is read-only")
	}

	if err := qy.Delete(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	return &m3uetcpb.Empty{}, nil
}

func (*QuerySvc) QueryBy(_ context.Context,
	req *m3uetcpb.QueryByRequest) (*m3uetcpb.QueryByResponse, error) {
	qy := models.FromProtobuf(req.Query)

	qybs := models.CollectionsToBoundaries(
		models.GetApplicableCollectionQueries(qy, req.Query.CollectionIds...),
	)
	ts := qy.FindTracks(qybs)

	if qy.Name != "" {
		go func() {
			if len(req.Query.CollectionIds) > 0 {
				onerror.Log(qy.SaveBound(qybs))
				return
			}
			onerror.Log(qy.Save())
		}()
	}

	out := []*m3uetcpb.Track{}
	for _, x := range ts {
		aux := x.ToProtobuf().(*m3uetcpb.Track)
		out = append(out, aux)
	}

	return &m3uetcpb.QueryByResponse{Tracks: out}, nil
}

func (*QuerySvc) QueryInPlaylist(_ context.Context,
	req *m3uetcpb.QueryInPlaylistRequest) (*m3uetcpb.QueryInPlaylistResponse, error) {

	if req.Id < 1 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Query id must be greater than zero")
	}

	qy := &models.Query{}
	if err := qy.Read(req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}

	qybs := models.CollectionsToBoundaries(
		models.GetApplicableCollectionQueries(qy),
	)

	var bar *models.Playbar
	pl := &models.Playlist{}
	if req.PlaylistId > 0 {
		pl.Read(req.PlaylistId)
		bar = &pl.Playbar
	} else {
		var err error
		bar, err = models.DefaultPerspective.GetPlaybar()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
		pl, err = bar.CreateEntry("", "", req.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
	}

	go bar.QueryInPlaylist(qy, qybs, pl)

	return &m3uetcpb.QueryInPlaylistResponse{PlaylistId: pl.ID}, nil
}

func (*QuerySvc) QueryInQueue(_ context.Context,
	req *m3uetcpb.QueryInQueueRequest) (*m3uetcpb.Empty, error) {

	if req.Id < 1 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Query id must be greater than zero")
	}

	qy := &models.Query{}
	if err := qy.Read(req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}

	qybs := models.CollectionsToBoundaries(
		models.GetApplicableCollectionQueries(qy),
	)

	q, err := models.PerspectiveIndex(req.Perspective).GetPerspectiveQueue()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	go q.QueryIn(qy, qybs)

	return &m3uetcpb.Empty{}, nil
}

func (*QuerySvc) SubscribeToQueryStore(_ *m3uetcpb.Empty,
	stream m3uetcpb.QuerySvc_SubscribeToQueryStoreServer) error {

	sub, id := subscription.Subscribe(subscription.ToQueryStoreEvent)
	defer func() { sub.Unsubscribe() }()

	go func() {
		sub.Event <- subscription.Event{Idx: int(models.QueryEventInitial)}
	}()

	sendQuery := func(e m3uetcpb.QueryEvent, qy models.ProtoOut) error {
		res := &m3uetcpb.SubscribeToQueryStoreResponse{
			SubscriptionId: id,
			Event:          e,
			Query:          qy.ToProtobuf().(*m3uetcpb.Query),
		}
		return stream.Send(res)
	}

sLoop:
	for {

		select {
		case e := <-sub.Event:
			if sub.MustUnsubscribe(e) {
				break sLoop
			}

			if models.QueryEvent(e.Idx) == models.QueryEventInitial {
				res := &m3uetcpb.SubscribeToQueryStoreResponse{
					SubscriptionId: id,
					Event:          m3uetcpb.QueryEvent_QYE_INITIAL,
				}
				err := stream.Send(res)
				if err != nil {
					return status.Errorf(codes.Internal,
						"Error sending query event (%v): %v",
						m3uetcpb.QueryEvent_QYE_INITIAL, err)
				}

				qys := models.GetAllQueries(0)
				for i := range qys {
					err := sendQuery(
						m3uetcpb.QueryEvent_QYE_INITIAL_ITEM,
						qys[i],
					)
					if err != nil {
						return status.Errorf(codes.Internal,
							"Error sending query event (%v): %v",
							m3uetcpb.QueryEvent_QYE_INITIAL_ITEM, err)
					}
				}

				res = &m3uetcpb.SubscribeToQueryStoreResponse{
					SubscriptionId: id,
					Event:          m3uetcpb.QueryEvent_QYE_INITIAL_DONE,
				}
				onerror.Log(stream.Send(res))
				continue sLoop
			}

			var eout m3uetcpb.QueryEvent

			switch models.QueryEvent(e.Idx) {
			case models.QueryEventItemAdded:
				eout = m3uetcpb.QueryEvent_QYE_ITEM_ADDED
			case models.QueryEventItemChanged:
				eout = m3uetcpb.QueryEvent_QYE_ITEM_CHANGED
			case models.QueryEventItemRemoved:
				eout = m3uetcpb.QueryEvent_QYE_ITEM_REMOVED
			default:
				slog.Error("Ignoring unsupported query event", "event", e.Idx)
				continue sLoop

			}

			if err := sendQuery(eout, e.Data.(models.ProtoOut)); err != nil {
				return status.Errorf(codes.Internal,
					"Error sending query event (%v): %v",
					eout, err)
			}
		}
	}
	return nil
}

func (*QuerySvc) UnsubscribeFromQueryStore(_ context.Context,
	req *m3uetcpb.UnsubscribeFromQueryStoreRequest) (*m3uetcpb.Empty, error) {
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
