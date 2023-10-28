package api

import (
	"context"
	"time"

	"github.com/jwmwalrus/gear-pieces/idler"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
)

// RootSvc implements the m3uetcpb.RootSvcServer interface.
type RootSvc struct {
	m3uetcpb.UnimplementedRootSvcServer
}

func (*RootSvc) Off(_ context.Context,
	req *m3uetcpb.OffRequest) (*m3uetcpb.OffResponse, error) {

	idler.DoTerminate(req.Force)
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go idler.Idle(ctx)
		time.Sleep(5 * time.Second)
	}()
	return &m3uetcpb.OffResponse{GoingOff: true}, nil
}

func (*RootSvc) Status(_ context.Context,
	_ *m3uetcpb.Empty) (*m3uetcpb.StatusResponse, error) {
	return &m3uetcpb.StatusResponse{Healthy: true}, nil
}
