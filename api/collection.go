package api

import (
	"context"
	"log/slog"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CollectionSvc implements the m3uetcpb.CollectionSvcServer interface.
type CollectionSvc struct {
	m3uetcpb.UnimplementedCollectionSvcServer
}

func (*CollectionSvc) GetCollection(_ context.Context,
	req *m3uetcpb.GetCollectionRequest) (*m3uetcpb.GetCollectionResponse, error) {

	if req.Id < 1 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Collection ID must be greater than zero")
	}

	coll := models.Collection{}
	if err := coll.Read(req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "Collection not found: %v", err)
	}

	coll.CountTracks()

	out := coll.ToProtobuf().(*m3uetcpb.Collection)
	return &m3uetcpb.GetCollectionResponse{
			Collection: out,
		},
		nil
}

func (*CollectionSvc) GetAllCollections(_ context.Context,
	_ *m3uetcpb.Empty) (*m3uetcpb.GetAllCollectionsResponse, error) {

	s := models.GetAllCollections()

	all := []*m3uetcpb.Collection{}
	for _, c := range s {
		c.CountTracks()
		out := c.ToProtobuf().(*m3uetcpb.Collection)
		all = append(all, out)
	}
	return &m3uetcpb.GetAllCollectionsResponse{Collections: all}, nil
}

func (*CollectionSvc) AddCollection(_ context.Context,
	req *m3uetcpb.AddCollectionRequest) (*m3uetcpb.AddCollectionResponse, error) {

	if req.Location == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"Collection location must not be empty")
	}

	name := "A collection"
	if req.Name != "" {
		name = req.Name
	}

	perspID := models.PerspectiveIndex(req.Perspective).Get().ID
	coll := models.Collection{
		Name:          name,
		Location:      req.Location,
		Disabled:      req.Disabled,
		Remote:        req.Remote,
		PerspectiveID: perspID,
	}

	if err := coll.Create(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"Error creating collection: %v", err)
	}

	go func() {
		if !coll.Remote {
			coll.Scan(false)
		}
	}()

	return &m3uetcpb.AddCollectionResponse{Id: coll.ID}, nil
}

func (*CollectionSvc) RemoveCollection(_ context.Context,
	req *m3uetcpb.RemoveCollectionRequest) (*m3uetcpb.Empty, error) {
	if req.Id < 1 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Collection ID must be greater than zero")
	}

	coll := models.Collection{}
	if err := coll.Read(req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "Collection not found: %v", err)
	}

	go func() {
		onerror.Log(coll.Delete())
	}()

	return &m3uetcpb.Empty{}, nil
}

func (*CollectionSvc) UpdateCollection(_ context.Context,
	req *m3uetcpb.UpdateCollectionRequest) (*m3uetcpb.Empty, error) {

	if req.Id < 1 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Collection ID must be greater than zero")
	}

	coll := models.Collection{}
	if err := coll.Read(req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "Collection not found: %v", err)
	}

	if req.NewName != "" {
		coll.Name = req.NewName
	}

	if req.NewDescription != "" {
		coll.Description = req.NewDescription
	}

	if req.ResetDescription {
		coll.Description = ""
	}

	if req.NewRemoteLocation != "" {
		coll.Remotelocation = req.NewRemoteLocation
	}

	if req.ResetRemoteLocation {
		coll.Remotelocation = ""
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
		return nil, status.Errorf(codes.InvalidArgument,
			"Error updating collection: %v", err)
	}

	return &m3uetcpb.Empty{}, nil
}

func (*CollectionSvc) ScanCollection(_ context.Context,
	req *m3uetcpb.ScanCollectionRequest) (*m3uetcpb.Empty, error) {

	if req.Id < 1 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Collection ID must be greater than zero")
	}

	coll := models.Collection{}
	if err := coll.Read(req.Id); err != nil {
		return nil, status.Errorf(codes.NotFound,
			"Collection not found: %v", err)
	}

	go func() {
		coll.Verify()
		coll.Scan(req.UpdateTags)
	}()

	return &m3uetcpb.Empty{}, nil
}

func (*CollectionSvc) DiscoverCollections(_ context.Context,
	_ *m3uetcpb.Empty) (*m3uetcpb.Empty, error) {

	s := models.GetAllCollections()

	go func() {
		for _, coll := range s {
			coll.Verify()
			coll.Scan(false)
		}
	}()

	return &m3uetcpb.Empty{}, nil
}

func (*CollectionSvc) SubscribeToCollectionStore(_ *m3uetcpb.Empty,
	stream m3uetcpb.CollectionSvc_SubscribeToCollectionStoreServer) error {

	sub, id := subscription.Subscribe(subscription.ToCollectionStoreEvent)
	defer func() { sub.Unsubscribe() }()

	go func() {
		sub.Event <- subscription.Event{Idx: int(models.CollectionEventInitial)}
	}()

	sendCollection := func(e m3uetcpb.CollectionEvent, c models.ProtoOut) error {
		res := &m3uetcpb.SubscribeToCollectionStoreResponse{
			SubscriptionId: id,
			Event:          e,
			Item: &m3uetcpb.SubscribeToCollectionStoreResponse_Collection{
				Collection: c.ToProtobuf().(*m3uetcpb.Collection),
			},
		}
		return stream.Send(res)
	}

	sendTrack := func(e m3uetcpb.CollectionEvent, t models.ProtoOut) error {
		res := &m3uetcpb.SubscribeToCollectionStoreResponse{
			SubscriptionId: id,
			Event:          e,
			Item: &m3uetcpb.SubscribeToCollectionStoreResponse_Track{
				Track: t.ToProtobuf().(*m3uetcpb.Track),
			},
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

			if models.CollectionEvent(e.Idx) == models.CollectionEventInitial ||
				models.CollectionEvent(e.Idx) == models.CollectionEventScanningDone {
				if models.CollectionEvent(e.Idx) == models.CollectionEventScanningDone {
					res := &m3uetcpb.SubscribeToCollectionStoreResponse{
						SubscriptionId: id,
						Event:          m3uetcpb.CollectionEvent_CE_SCANNING_DONE,
					}
					err := stream.Send(res)
					if err != nil {
						return status.Errorf(codes.Internal,
							"Error sending event (%v): %v",
							m3uetcpb.CollectionEvent_CE_SCANNING_DONE, err)
					}
				}

				res := &m3uetcpb.SubscribeToCollectionStoreResponse{
					SubscriptionId: id,
					Event:          m3uetcpb.CollectionEvent_CE_INITIAL,
				}
				err := stream.Send(res)
				if err != nil {
					return status.Errorf(codes.Internal,
						"Error sending event (%v): %v",
						m3uetcpb.CollectionEvent_CE_INITIAL, err)
				}

				cs, ts := models.GetCollectionStore()
				for i := range cs {
					err := sendCollection(
						m3uetcpb.CollectionEvent_CE_INITIAL_ITEM,
						cs[i],
					)
					if err != nil {
						return status.Errorf(codes.Internal,
							"Error sending event (%v): %v",
							m3uetcpb.CollectionEvent_CE_INITIAL_ITEM, err)
					}
				}

				for i := range ts {
					err := sendTrack(
						m3uetcpb.CollectionEvent_CE_INITIAL_ITEM,
						ts[i],
					)
					if err != nil {
						return status.Errorf(codes.Internal,
							"Error sending event (%v): %v",
							m3uetcpb.CollectionEvent_CE_INITIAL_ITEM, err)
					}
				}

				res = &m3uetcpb.SubscribeToCollectionStoreResponse{
					SubscriptionId: id,
					Event:          m3uetcpb.CollectionEvent_CE_INITIAL_DONE,
				}
				onerror.Log(stream.Send(res))
				continue sLoop
			}

			if models.CollectionEvent(e.Idx) == models.CollectionEventScanning {
				res := &m3uetcpb.SubscribeToCollectionStoreResponse{
					SubscriptionId: id,
					Event:          m3uetcpb.CollectionEvent_CE_SCANNING,
				}
				err := stream.Send(res)
				if err != nil {
					return status.Errorf(codes.Internal,
						"Error sending event (%v): %v",
						m3uetcpb.CollectionEvent_CE_SCANNING, err)
				}
				continue sLoop
			}

			var eout m3uetcpb.CollectionEvent
			var fn func(m3uetcpb.CollectionEvent, models.ProtoOut) error

			switch models.CollectionEvent(e.Idx) {
			case models.CollectionEventItemAdded:
				eout = m3uetcpb.CollectionEvent_CE_ITEM_ADDED
			case models.CollectionEventItemChanged:
				eout = m3uetcpb.CollectionEvent_CE_ITEM_CHANGED
			case models.CollectionEventItemRemoved:
				eout = m3uetcpb.CollectionEvent_CE_ITEM_REMOVED
			default:
				slog.Error("Ignoring unsupported collection event", "event", e.Idx)
				continue sLoop

			}

			switch e.Data.(type) {
			case *models.Collection:
				fn = sendCollection
			case *models.Track:
				fn = sendTrack
			default:
				slog.With(
					"event", e.Idx,
					"data", e.Data,
				).Error("Ignoring unsupported data")
				continue sLoop
			}

			if err := fn(eout, e.Data.(models.ProtoOut)); err != nil {
				return status.Errorf(codes.Internal,
					"Error sending event (%v): %v",
					eout, err)
			}
		}
	}
	return nil

}

func (*CollectionSvc) UnsubscribeFromCollectionStore(_ context.Context,
	req *m3uetcpb.UnsubscribeFromCollectionStoreRequest) (*m3uetcpb.Empty, error) {
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
