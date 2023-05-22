package api

import (
	"context"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestGetQueue(t *testing.T) {
	table := []testCase{
		{
			"Get with empty queue",
			"api/queue/get-empty",
			&m3uetcpb.GetQueueRequest{},
			&m3uetcpb.GetQueueResponse{},
			false,
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
			false,
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
			false,
		},
	}

	svc := QueueSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetQueueResponse)

			res, err := svc.GetQueue(context.Background(), tc.req.(*m3uetcpb.GetQueueRequest))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, len(exp.QueueTracks), len(res.QueueTracks))
			if len(exp.QueueTracks) > 0 {
				assert.Equal(t, exp.QueueTracks[0].Id, res.QueueTracks[0].Id)
				assert.Equal(t, exp.QueueTracks[0].Location, res.QueueTracks[0].Location)
				assert.Equal(t, exp.QueueTracks[0].Played, res.QueueTracks[0].Played)
				assert.Equal(t, exp.QueueTracks[0].Position, res.QueueTracks[0].Position)
			}

			assert.Equal(t, len(res.Tracks), len(exp.Tracks))
			if len(exp.Tracks) > 0 {
				assert.Equal(t, exp.Tracks[0].Id, res.Tracks[0].Id)
				assert.Equal(t, exp.Tracks[0].Location, res.Tracks[0].Location)
				assert.Equal(t, exp.Tracks[0].Title, res.Tracks[0].Title)
				assert.Equal(t, exp.Tracks[0].Album, res.Tracks[0].Album)
				assert.Equal(t, exp.Tracks[0].Artist, res.Tracks[0].Artist)
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
			true,
		},
		{
			"Execute queue invalid ID",
			"api/queue/exec-invalid-id",
			&m3uetcpb.ExecuteQueueActionRequest{
				Action: m3uetcpb.QueueAction_Q_APPEND,
				Ids:    []int64{2},
			},
			&m3uetcpb.Empty{},
			true,
		},
		{
			"Execute valid",
			"api/queue/exec-valid",
			&m3uetcpb.ExecuteQueueActionRequest{
				Action: m3uetcpb.QueueAction_Q_CLEAR,
			},
			&m3uetcpb.Empty{},
			false,
		},
	}

	svc := QueueSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			_, err := svc.ExecuteQueueAction(context.Background(), tc.req.(*m3uetcpb.ExecuteQueueActionRequest))

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}

	return
}

func TestQueueTrackToProtobuf(t *testing.T) {
	tests.SetupTest(t, fixturesDir("api/queue/queue-track-to-protobuf"))
	t.Cleanup(func() { tests.TeardownTest(t) })

	qt := models.QueueTrack{}

	qt.Read(1)

	out := qt.ToProtobuf()

	qtpb, ok := out.(*m3uetcpb.QueueTrack)
	assert.True(t, ok)

	assert.Equal(t, qt.ID, qtpb.Id)
	assert.Equal(t, int32(qt.Position), qtpb.Position)
	assert.Equal(t, qt.Played, qtpb.Played)
	assert.Equal(t, qt.Location, qtpb.Location)
	assert.Equal(t, qt.CreatedAt, qtpb.CreatedAt)
	assert.Equal(t, qt.UpdatedAt, qtpb.UpdatedAt)
	assert.Equal(t, qt.TrackID, qtpb.TrackId)
	assert.Equal(t, m3uetcpb.Perspective(qt.Queue.Perspective.Idx), qtpb.Perspective)
}
