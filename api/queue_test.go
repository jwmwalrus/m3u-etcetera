package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/stretchr/testify/assert"
)

func TestGetQueue(t *testing.T) {
	table := []testCase{
		{
			"Get with empty queue",
			"api/queue/get-empty",
			&m3uetcpb.GetQueueRequest{},
			&m3uetcpb.GetQueueResponse{},
			nil,
		},
		{
			"Get with location",
			"api/queue/get-location",
			&m3uetcpb.GetQueueRequest{},
			&m3uetcpb.GetQueueResponse{
				QueueTracks: []*m3uetcpb.QueueTrack{{
					Id:       1,
					Position: 1,
					Location: "./data/testing/audio1/track01.ogg",
				}},
			},
			nil,
		},
		{
			"Get with track",
			"api/queue/get-track",
			&m3uetcpb.GetQueueRequest{},
			&m3uetcpb.GetQueueResponse{
				QueueTracks: []*m3uetcpb.QueueTrack{{
					Id:       1,
					Position: 1,
					Played:   false,
					Location: "./data/testing/audio1/track01.ogg",
					TrackId:  1,
				}},
				Tracks: []*m3uetcpb.Track{{
					Id:       1,
					Location: "./data/testing/audio1/track01.ogg",
					Title:    "track",
					Album:    "tracks",
					Artist:   "tracker",
				}},
			},
			nil,
		},
	}

	svc := QueueSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetQueueResponse)

			res, err := svc.GetQueue(context.Background(), tc.req.(*m3uetcpb.GetQueueRequest))

			assert.Equal(t, err != nil, tc.err)
			assert.Equal(t, len(res.QueueTracks), len(exp.QueueTracks))
			if len(exp.QueueTracks) > 0 {
				assert.Equal(t, res.QueueTracks[0].Id, exp.QueueTracks[0].Id)
				assert.Equal(t, res.QueueTracks[0].Location, exp.QueueTracks[0].Location)
				assert.Equal(t, res.QueueTracks[0].Played, exp.QueueTracks[0].Played)
				assert.Equal(t, res.QueueTracks[0].Position, exp.QueueTracks[0].Position)
			}

			assert.Equal(t, len(res.Tracks), len(exp.Tracks))
			if len(exp.Tracks) > 0 {
				assert.Equal(t, res.Tracks[0].Id, exp.Tracks[0].Id)
				assert.Equal(t, res.Tracks[0].Location, exp.Tracks[0].Location)
				assert.Equal(t, res.Tracks[0].Title, exp.Tracks[0].Title)
				assert.Equal(t, res.Tracks[0].Album, exp.Tracks[0].Album)
				assert.Equal(t, res.Tracks[0].Artist, exp.Tracks[0].Artist)
			}
		})
	}
	return
}

func TestExecuteQueueAction(t *testing.T) {
	table := []testCase{
		{
			"Execute queue invalid location",
			"api/queue/exec-invalid-loc",
			&m3uetcpb.ExecuteQueueActionRequest{
				Action:    m3uetcpb.QueueAction_Q_APPEND,
				Locations: []string{"2"},
			},
			&m3uetcpb.Empty{},
			fmt.Errorf("error"),
		},
		{
			"Execute queue invalid ID",
			"api/queue/exec-invalid-id",
			&m3uetcpb.ExecuteQueueActionRequest{
				Action: m3uetcpb.QueueAction_Q_APPEND,
				Ids:    []int64{2},
			},
			&m3uetcpb.Empty{},
			fmt.Errorf("error"),
		},
		{
			"Execute valid",
			"api/queue/exec-valid",
			&m3uetcpb.ExecuteQueueActionRequest{
				Action: m3uetcpb.QueueAction_Q_CLEAR,
			},
			&m3uetcpb.Empty{},
			fmt.Errorf("error"),
		},
	}

	svc := QueueSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			_, err := svc.ExecuteQueueAction(context.Background(), tc.req.(*m3uetcpb.ExecuteQueueActionRequest))

			assert.Equal(t, err != nil, tc.err)
		})
	}

	return
}
