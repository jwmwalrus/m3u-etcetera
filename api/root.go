package api

import (
	"context"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
)

// Root defines the root service
type Root struct {
	m3uetcpb.UnimplementedRootServer
}

// Off initiates the process to unload the server
func (r *Root) Off(_ context.Context, req *m3uetcpb.Empty) (*m3uetcpb.OffResponse, error) {
	go func() {
		time.Sleep(5 * time.Second)
		base.Idle(true)
	}()
	return &m3uetcpb.OffResponse{GoingOff: true}, nil
}

// Status returns the status of the server
func (r *Root) Status(_ context.Context, req *m3uetcpb.Empty) (*m3uetcpb.StatusResponse, error) {
	return &m3uetcpb.StatusResponse{Alive: true}, nil
}
