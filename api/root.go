package api

import (
	"context"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
)

// RootSvc defines the root service
type RootSvc struct {
	m3uetcpb.UnimplementedRootSvcServer
}

// Off implements RootSvcServer
// Initiates the process to unload the server
func (*RootSvc) Off(_ context.Context, _ *m3uetcpb.Empty) (*m3uetcpb.OffResponse, error) {
	go func() {
		time.Sleep(5 * time.Second)
		base.Idle(true)
	}()
	return &m3uetcpb.OffResponse{GoingOff: true}, nil
}

// Status implements RootSvcServer
// Returns the status of the server
func (*RootSvc) Status(_ context.Context, _ *m3uetcpb.Empty) (*m3uetcpb.StatusResponse, error) {
	return &m3uetcpb.StatusResponse{Alive: true}, nil
}
