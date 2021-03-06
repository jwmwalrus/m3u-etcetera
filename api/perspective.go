package api

import (
	"context"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PerspectiveSvc defines the perspective service
type PerspectiveSvc struct {
	m3uetcpb.UnimplementedPerspectiveSvcServer
}

// GetActivePerspective implements m3uetcpb.PerspectiveSvcServer
func (p *PerspectiveSvc) GetActivePerspective(_ context.Context,
	_ *m3uetcpb.Empty) (*m3uetcpb.GetActivePerspectiveResponse, error) {
	res := &m3uetcpb.GetActivePerspectiveResponse{
		Perspective: m3uetcpb.Perspective(models.GetActivePerspectiveIndex()),
	}

	return res, nil
}

// SetActivePerspective implements m3uetcpb.PerspectiveSvcServer
func (p *PerspectiveSvc) SetActivePerspective(_ context.Context,
	req *m3uetcpb.SetActivePerspectiveRequest) (*m3uetcpb.Empty, error) {

	persp := models.PerspectiveIndex(req.Perspective)
	err := persp.Activate()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"Error activating the %v perspective: %v", persp, err)
	}

	return &m3uetcpb.Empty{}, nil
}

// SubscribeToPerspective implements m3uetcpb.PerspectiveSvcServer
func (p *PerspectiveSvc) SubscribeToPerspective(_ *m3uetcpb.Empty,
	stream m3uetcpb.PerspectiveSvc_SubscribeToPerspectiveServer) error {

	s, id := subscription.Subscribe(subscription.ToPerspectiveEvent)
	defer func() { s.Unsubscribe() }()

	go func() {
		time.Sleep(2 * time.Second)
		s.Event <- subscription.Event{Data: struct{}{}}
	}()

sLoop:
	for {
		select {
		case e := <-s.Event:
			if s.MustUnsubscribe(e) {
				break sLoop
			}

			res := &m3uetcpb.SubscribeToPerspectiveResponse{
				SubscriptionId: id,
				ActivePerspective: m3uetcpb.Perspective(
					models.GetActivePerspectiveIndex(),
				),
			}

			err := stream.Send(res)
			if err != nil {
				return status.Errorf(codes.Internal,
					"Error sending perspective: %v", err)
			}
		}
	}
	return nil
}

// UnsubscribeFromPerspective implements m3uetcpb.PerspectiveSvcServer
func (p *PerspectiveSvc) UnsubscribeFromPerspective(_ context.Context,
	req *m3uetcpb.UnsubscribeFromPerspectiveRequest) (*m3uetcpb.Empty, error) {

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
