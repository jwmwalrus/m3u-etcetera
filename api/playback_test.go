package api

import (
	"context"
	"testing"

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
			&m3uetcpb.Empty{},
			&m3uetcpb.GetPlaybackResponse{
				Playing: &m3uetcpb.GetPlaybackResponse_Empty{
					Empty: &m3uetcpb.Empty{},
				},
			},
			false,
		},
		{
			"Get with playback, TrackID=0",
			"api/playback/get-notrackid",
			true,
			0,
			&m3uetcpb.Empty{},
			&m3uetcpb.GetPlaybackResponse{
				Playing: &m3uetcpb.GetPlaybackResponse_Playback{
					Playback: &m3uetcpb.Playback{
						Id:       1,
						Location: "./data/testing/audio1/track01.ogg",
					},
				},
			},
			false,
		},
		{
			"Get with playback, TrackID>0",
			"api/playback/get-trackid",
			true,
			0,
			&m3uetcpb.Empty{},
			&m3uetcpb.GetPlaybackResponse{
				Playing: &m3uetcpb.GetPlaybackResponse_Track{
					Track: &m3uetcpb.Track{
						Id:       1,
						Location: "./data/testing/audio1/track01.ogg",
						Title:    "track",
						Album:    "tracks",
						Artist:   "tracker",
					},
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

			res, err := svc.GetPlayback(context.Background(), tc.req.(*m3uetcpb.Empty))

			assert.Equal(t, err != nil, tc.err)
			if exp.GetEmpty() != nil {
				assert.Equal(t, res.GetEmpty() != nil, true)
			}
			if exp.GetPlayback() != nil {
				assert.Equal(t, res.GetPlayback() != nil, true)
				assert.Equal(t, res.GetPlayback().Id, exp.GetPlayback().Id)
				assert.Equal(t, res.GetPlayback().Location, exp.GetPlayback().Location)
			}
			if exp.GetTrack() != nil {
				assert.Equal(t, res.GetTrack() != nil, true)
				assert.Equal(t, res.GetTrack().Id, exp.GetTrack().Id)
				assert.Equal(t, res.GetTrack().Location, exp.GetTrack().Location)
				assert.Equal(t, res.GetTrack().Title, exp.GetTrack().Title)
				assert.Equal(t, res.GetTrack().Album, exp.GetTrack().Album)
				assert.Equal(t, res.GetTrack().Artist, exp.GetTrack().Artist)
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
			&m3uetcpb.Empty{},
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
			&m3uetcpb.Empty{},
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
			&m3uetcpb.Empty{},
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
