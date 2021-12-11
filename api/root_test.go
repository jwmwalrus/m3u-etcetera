package api

import (
	"context"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/api/pb"
)

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
