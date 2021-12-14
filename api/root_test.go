package api

import (
	"context"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
)

func TestRootOff(t *testing.T) {
	c := RootSvc{}

	req := &m3uetcpb.Empty{}
	res, err := c.Off(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !res.GetGoingOff() || res.GetReason() != "" {
		t.Fail()
	}
}

func TestRootStatus(t *testing.T) {
	c := RootSvc{}

	req := &m3uetcpb.Empty{}
	res, err := c.Status(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !res.GetAlive() {
		t.Fail()
	}
}
