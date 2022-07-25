package api

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/stretchr/testify/assert"
)

func TestGetPlayback(t *testing.T) {
	table := []testCase{
		{
			"Get with no playback",
			"api/playback/get-nopb",
			false,
			0,
			&empty.Empty{},
			&m3uetcpb.GetPlaybackResponse{},
			false,
		},
		{
			"Get with playback, TrackID=0",
			"api/playback/get-notrackid",
			true,
			0,
			&empty.Empty{},
			&m3uetcpb.GetPlaybackResponse{
				IsPlaying: true,
				Playback: &m3uetcpb.Playback{
					Id:       1,
					Location: "./data/testing/audio1/track01.ogg",
				},
				Track: &m3uetcpb.Track{
					Location: "./data/testing/audio1/track01.ogg",
				},
			},
			false,
		},
		{
			"Get with playback, TrackID>0",
			"api/playback/get-trackid",
			true,
			0,
			&empty.Empty{},
			&m3uetcpb.GetPlaybackResponse{
				IsPlaying: true,
				Playback: &m3uetcpb.Playback{
					Id:       1,
					Location: "./data/testing/audio1/track01.ogg",
				},
				Track: &m3uetcpb.Track{
					Id:       1,
					Location: "./data/testing/audio1/track01.ogg",
					Title:    "track",
					Album:    "tracks",
					Artist:   "tracker",
				},
			},
			false,
		},
	}

	svc := PlaybackSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetPlaybackResponse)

			res, err := svc.GetPlayback(context.Background(), tc.req.(*empty.Empty))

			assert.Equal(t, err != nil, tc.err)
			if exp.IsPlaying {
				assert.Equal(t, res.Playback != nil, true)
				assert.Equal(t, res.Playback.Id, exp.Playback.Id)
				assert.Equal(t, res.Playback.Location, exp.Playback.Location)
				assert.Equal(t, res.Track != nil, true)
				assert.Equal(t, res.Track.Id, exp.Track.Id)
				assert.Equal(t, res.Track.Location, exp.Track.Location)
				assert.Equal(t, res.Track.Title, exp.Track.Title)
				assert.Equal(t, res.Track.Album, exp.Track.Album)
				assert.Equal(t, res.Track.Artist, exp.Track.Artist)
			}
		})
	}
	return
}

func TestExecutePlaybackAction(t *testing.T) {
	table := []testCase{
		{
			"Execute play invalid location",
			"api/playback/exec-invalid-loc",
			false,
			0,
			&m3uetcpb.ExecutePlaybackActionRequest{
				Action:    m3uetcpb.PlaybackAction_PB_PLAY,
				Locations: []string{"2"},
			},
			&empty.Empty{},
			true,
		},
		{
			"Execute play invalid ID",
			"api/playback/exec-invalid-id",
			false,
			0,
			&m3uetcpb.ExecutePlaybackActionRequest{
				Action: m3uetcpb.PlaybackAction_PB_PLAY,
				Ids:    []int64{2},
			},
			&empty.Empty{},
			true,
		},
		{
			"Execute valid",
			"api/playback/exec-valid",
			false,
			0,
			&m3uetcpb.ExecutePlaybackActionRequest{
				Action: m3uetcpb.PlaybackAction_PB_NEXT,
			},
			&empty.Empty{},
			false,
		},
	}

	svc := PlaybackSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			_, err := svc.ExecutePlaybackAction(context.Background(), tc.req.(*m3uetcpb.ExecutePlaybackActionRequest))

			assert.Equal(t, err != nil, tc.err)
		})
	}

	return
}
