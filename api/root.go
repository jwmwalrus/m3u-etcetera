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
func (*RootSvc) Off(_ context.Context,
	req *m3uetcpb.OffRequest) (*m3uetcpb.OffResponse, error) {

	base.DoTerminate(req.Force)
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go base.Idle(ctx)
		time.Sleep(5 * time.Second)
	}()
	return &m3uetcpb.OffResponse{GoingOff: true}, nil
}

// Status implements RootSvcServer
// Returns the status of the server
func (*RootSvc) Status(_ context.Context,
	_ *m3uetcpb.Empty) (*m3uetcpb.StatusResponse, error) {
	return &m3uetcpb.StatusResponse{Alive: true}, nil
}
