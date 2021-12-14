package api

import (
	"context"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CollectionSvc defines the collection service
type CollectionSvc struct {
	m3uetcpb.UnimplementedCollectionSvcServer
}

// GetCollection implements m3uetcpb.CollectionSvcServer
func (*CollectionSvc) GetCollection(_ context.Context, req *m3uetcpb.GetCollectionRequest) (*m3uetcpb.GetCollectionResponse, error) {
	if req.Id < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Collection ID must be greater than zero")
	}

	coll := models.Collection{}
	if err := coll.Read(req.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Collection not found: %v", err)
	}

	coll.CountTracks()

	return &m3uetcpb.GetCollectionResponse{Collection: coll.ToProtobuf()}, nil
}

// GetAllCollections implements m3uetcpb.CollectionSvcServer
func (*CollectionSvc) GetAllCollections(_ context.Context, _ *m3uetcpb.Empty) (*m3uetcpb.GetAllCollectionsResponse, error) {
	s := models.GetAllCollections()

	all := []*m3uetcpb.Collection{}
	for _, c := range s {
		c.CountTracks()
		aux := c.ToProtobuf()
		all = append(all, aux)
	}
	return &m3uetcpb.GetAllCollectionsResponse{Collections: all}, nil
}

// AddCollection implements m3uetcpb.CollectionSvcServer
func (*CollectionSvc) AddCollection(_ context.Context, req *m3uetcpb.AddCollectionRequest) (*m3uetcpb.AddCollectionResponse, error) {
	if req.Location == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Collection location must not be empty")
	}

	name := "A collection"
	if req.Name != "" {
		name = req.Name
	}

	coll := models.Collection{
		Name:     name,
		Location: req.Location,
		Disabled: req.Disabled,
		Remote:   req.Remote,
	}

	if err := coll.Create(); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Error creating collection: %v", err)
	}

	go func() {
		if !coll.Remote {
			coll.Scan(false)
		}
	}()

	return &m3uetcpb.AddCollectionResponse{Id: coll.ID}, nil
}

// RemoveCollection implements m3uetcpb.CollectionSvcServer
func (*CollectionSvc) RemoveCollection(_ context.Context, req *m3uetcpb.RemoveCollectionRequest) (*m3uetcpb.Empty, error) {
	if req.Id < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Collection ID must be greater than zero")
	}

	coll := models.Collection{}
	if err := coll.Read(req.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Collection not found: %v", err)
	}

	go func() {
		err := coll.Delete()
		onerror.Log(err)
	}()

	return &m3uetcpb.Empty{}, nil
}

// UpdateCollection implements m3uetcpb.CollectionSvcServer
func (*CollectionSvc) UpdateCollection(_ context.Context, req *m3uetcpb.UpdateCollectionRequest) (*m3uetcpb.Empty, error) {
	if req.Id < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Collection ID must be greater than zero")
	}

	coll := models.Collection{}
	if err := coll.Read(req.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Collection not found: %v", err)
	}

	if req.NewName != "" {
		coll.Name = req.NewName
	}

	if req.NewDescription != "" {
		coll.Description = req.NewDescription
	}

	if req.Enable {
		coll.Disabled = false
	}

	if req.Disable {
		coll.Disabled = true
	}

	if req.MakeLocal {
		coll.Remote = false
	}

	if req.MakeRemote {
		coll.Remote = true
	}

	if err := coll.Save(); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Error updating collection: %v", err)
	}

	return &m3uetcpb.Empty{}, nil
}

// ScanCollection implements m3uetcpb.CollectionSvcServer
func (*CollectionSvc) ScanCollection(_ context.Context, req *m3uetcpb.ScanCollectionRequest) (*m3uetcpb.Empty, error) {
	if req.Id < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Collection ID must be greater than zero")
	}

	coll := models.Collection{}
	if err := coll.Read(req.Id); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Collection not found: %v", err)
	}

	go func() {
		coll.Verify()
		coll.Scan(req.UpdateTags)
	}()

	return &m3uetcpb.Empty{}, nil
}

// DiscoverCollections implements m3uetcpb.CollectionSvcServer
func (*CollectionSvc) DiscoverCollections(_ context.Context, _ *m3uetcpb.Empty) (*m3uetcpb.Empty, error) {
	s := models.GetAllCollections()

	go func() {
		for _, coll := range s {
			coll.Scan(false)
		}
	}()

	return &m3uetcpb.Empty{}, nil
}
