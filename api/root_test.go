package api

import (
	"context"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/stretchr/testify/assert"
)

func TestRootOff(t *testing.T) {
	c := RootSvc{}

	req := &m3uetcpb.OffRequest{}
	res, err := c.Off(context.Background(), req)
	assert.Equal(t, err != nil, false)
	assert.Equal(t, res.GetGoingOff() || res.GetReason() == "", true)
}

func TestRootStatus(t *testing.T) {
	c := RootSvc{}

	req := &m3uetcpb.Empty{}
	res, err := c.Status(context.Background(), req)
	assert.Equal(t, err != nil, false)
	assert.Equal(t, res.GetAlive(), true)
}
