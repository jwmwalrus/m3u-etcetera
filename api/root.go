package api

import (
	"context"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/pb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
)

// Root defines the root service
type Root struct {
	pb.UnimplementedRootServer
}

// Off initiates the process to unload the server
func (r *Root) Off(_ context.Context, req *pb.Empty) (*pb.OffResponse, error) {
	go func() {
		time.Sleep(5 * time.Second)
		base.Idle(true)
	}()
	return &pb.OffResponse{GoingOff: true}, nil
}

// Status returns the status of the server
func (r *Root) Status(_ context.Context, req *pb.Empty) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{Alive: true}, nil
}
