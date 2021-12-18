package api

import (
	"context"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/pkg/qparams"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// QuerySvc defines the query service
type QuerySvc struct {
	m3uetcpb.UnimplementedQuerySvcServer
}

// GetQuery implements m3uetcpb.QuerySvcServer
func (*QuerySvc) GetQuery(_ context.Context, req *m3uetcpb.GetQueryRequest) (*m3uetcpb.GetQueryResponse, error) {

	if req.Id < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Query id must be greater than zero")
	}

	q := models.Query{}
	if err := q.Read(req.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", err)
	}

	out := q.ToProtobuf().(*m3uetcpb.Query)
	return &m3uetcpb.GetQueryResponse{Query: out}, nil
}

// GetQueries implements m3uetcpb.QuerySvcServer
func (*QuerySvc) GetQueries(_ context.Context, req *m3uetcpb.GetQueriesRequest) (*m3uetcpb.GetQueriesResponse, error) {

	qbs := models.FilterCollectionQueryBoundaries(req.CollectionIds)
	qs := models.GetAllQueries(int(req.Limit), qbs...)

	out := []*m3uetcpb.Query{}
	for _, x := range qs {
		aux := x.ToProtobuf().(*m3uetcpb.Query)
		out = append(out, aux)
	}

	return &m3uetcpb.GetQueriesResponse{Queries: out}, nil
}

// AddQuery implements m3uetcpb.QuerySvcServer
func (*QuerySvc) AddQuery(_ context.Context, req *m3uetcpb.AddQueryRequest) (*m3uetcpb.AddQueryResponse, error) {
	if req.Query.Params != "" {

		if _, err := qparams.ParseParams(req.Query.Params); err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Error parsing query params: %v", err)
		}
	}

	q := models.FromProtobuf(req.Query)

	qbs := models.CreateCollectionQueryBoundaries(req.Query.CollectionIds)
	if err := q.SaveBound(qbs); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Error saving query: %v", err)
	}

	return &m3uetcpb.AddQueryResponse{Id: q.ID}, nil
}

// UpdateQuery implements m3uetcpb.QuerySvcServer
func (*QuerySvc) UpdateQuery(_ context.Context, req *m3uetcpb.UpdateQueryRequest) (*m3uetcpb.Empty, error) {
	q := models.Query{}
	if err := q.Read(req.Query.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", err)
	}
	q.FromProtobuf(req.Query)

	models.RemoveCollections(q.GetCollections())
	qbs := models.CreateCollectionQueryBoundaries(req.Query.CollectionIds)
	if err := q.SaveBound(qbs); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Error saving query: %v", err)
	}

	return &m3uetcpb.Empty{}, nil
}

// RemoveQuery implements m3uetcpb.QuerySvcServer
func (*QuerySvc) RemoveQuery(_ context.Context, req *m3uetcpb.RemoveQueryRequest) (*m3uetcpb.Empty, error) {
	q := models.Query{}
	if err := q.Read(req.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", err)
	}

	if err := q.Delete(); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "%v", err)
	}

	return &m3uetcpb.Empty{}, nil
}

// ApplyQuery implements m3uetcpb.QuerySvcServer
func (*QuerySvc) ApplyQuery(_ context.Context, req *m3uetcpb.ApplyQueryRequest) (*m3uetcpb.ApplyQueryResponse, error) {

	if req.Id < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Query id must be greater than zero")
	}

	q := models.Query{}
	if err := q.Read(req.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", err)
	}

	qbs := models.CollectionsToBoundaries(q.GetCollections())
	ts := q.FindTracks(qbs...)

	out := []*m3uetcpb.Track{}
	for _, x := range ts {
		aux := x.ToProtobuf().(*m3uetcpb.Track)
		out = append(out, aux)
	}

	return &m3uetcpb.ApplyQueryResponse{Tracks: out}, nil
}

// QueryBy implements m3uetcpb.QuerySvcServer
func (*QuerySvc) QueryBy(_ context.Context, req *m3uetcpb.QueryByRequest) (*m3uetcpb.QueryByResponse, error) {
	q := models.FromProtobuf(req.Query)

	ts := q.FindTracks()

	if q.Name != "" {
		go func() {
			err := q.Save()
			onerror.Log(err)
		}()
	}

	out := []*m3uetcpb.Track{}
	for _, x := range ts {
		aux := x.ToProtobuf().(*m3uetcpb.Track)
		out = append(out, aux)
	}

	return &m3uetcpb.QueryByResponse{Tracks: out}, nil
}
