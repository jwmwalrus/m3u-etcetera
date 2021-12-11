package api

import (
	"context"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/api/pb"
	"google.golang.org/grpc"
)

func getClient(t *testing.T) (c pb.RootClient, cc *grpc.ClientConn) {
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		t.Fatal(err)
	}

	c = pb.NewRootClient(cc)
	return
}

func TestRootOff(t *testing.T) {
	c := Root{}

	req := &pb.Empty{}
	res, err := c.Off(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !res.GetGoingOff() || res.GetReason() != "" {
		t.Fail()
	}
}

func TestRootStatus(t *testing.T) {
	c := Root{}

	req := &pb.Empty{}
	res, err := c.Status(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !res.GetAlive() {
		t.Fail()
	}
}
