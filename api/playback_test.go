package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/stretchr/testify/assert"
	"github.com/tinyzimmer/go-gst/gst"
)

func TestGetPlayback(t *testing.T) {
	events := pbEventsMock{
		pb: &models.Playback{
			ID:       1,
			Location: "./data/testing/audio1/track01.ogg",
		},
		t: &models.Track{
			ID:       1,
			Location: "./data/testing/audio1/track01.ogg",
			Title:    "track",
			Album:    "tracks",
			Artist:   "tracker",
		},
		isPlaying: true,
	}

	table := []testCase{
		{
			"Get with no playback",
			"api/playback/get-nopb",
			&m3uetcpb.Empty{},
			&m3uetcpb.GetPlaybackResponse{},
			nil,
		},
		{
			"Get with playback, TrackID=0",
			"api/playback/get-notrackid",
			&m3uetcpb.Empty{},
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
			nil,
		},
		{
			"Get with playback, TrackID>0",
			"api/playback/get-trackid",
			&m3uetcpb.Empty{},
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
			nil,
		},
	}

	svc := PlaybackSvc{PbEvents: &events}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetPlaybackResponse)

			res, err := svc.GetPlayback(context.Background(), tc.req.(*m3uetcpb.Empty))

			if tc.err != nil {
				assert.NotNil(t, err)
			}
			if tc.err == nil {
				assert.Nil(t, err)

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
			&m3uetcpb.ExecutePlaybackActionRequest{
				Action:    m3uetcpb.PlaybackAction_PB_PLAY,
				Locations: []string{"2"},
			},
			&m3uetcpb.Empty{},
			fmt.Errorf("error"),
		},
		{
			"Execute play invalid ID",
			"api/playback/exec-invalid-id",
			&m3uetcpb.ExecutePlaybackActionRequest{
				Action: m3uetcpb.PlaybackAction_PB_PLAY,
				Ids:    []int64{2},
			},
			&m3uetcpb.Empty{},
			fmt.Errorf("error"),
		},
		{
			"Execute valid",
			"api/playback/exec-valid",
			&m3uetcpb.ExecutePlaybackActionRequest{
				Action: m3uetcpb.PlaybackAction_PB_NEXT,
			},
			&m3uetcpb.Empty{},
			nil,
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

type pbEventsMock struct {
	pb             *models.Playback
	t              *models.Track
	state          gst.State
	hasNextStream  bool
	isPaused       bool
	isPlaying      bool
	isReady        bool
	isStreaming    bool
	isStopped      bool
	nextStreamErr  error
	pauseStreamErr error
}

// GetPlayback implements the playback.IEvents interface
func (e *pbEventsMock) GetPlayback() (pb *models.Playback, t *models.Track) {
	return e.pb, e.t
}

// GetState implements the playback.IEvents interface
func (e *pbEventsMock) GetState() gst.State {
	return e.state
}

// HasNextStream implements the playback.IEvents interface
func (e *pbEventsMock) HasNextStream() bool {
	return e.hasNextStream
}

// IsPaused implements the playback.IEvents interface
func (e *pbEventsMock) IsPaused() bool {
	return e.isPaused
}

// IsPlaying implements the playback.IEvents interface
func (e *pbEventsMock) IsPlaying() bool {
	return e.isPlaying
}

// IsReady implements the playback.IEvents interface
func (e *pbEventsMock) IsReady() bool {
	return e.isReady
}

// IsStreaming implements the playback.IEvents interface
func (e *pbEventsMock) IsStreaming() bool {
	return e.isStreaming
}

// IsStopped implements the playback.IEvents interface
func (e *pbEventsMock) IsStopped() bool {
	return e.isStopped
}

// NextStream implements the playback.IEvents interface
func (e *pbEventsMock) NextStream() (err error) {
	return e.nextStreamErr
}

// PauseStream implements the playback.IEvents interface
func (e *pbEventsMock) PauseStream(off bool) (err error) {
	return e.pauseStreamErr
}

// PlayStreams implements the playback.IEvents interface
func (p *pbEventsMock) PlayStreams(force bool, locations []string, ids []int64) {}

// PreviousStream implements the playback.IEvents interface
func (p *pbEventsMock) PreviousStream() {}

// QuitPlayingFromBar implements the playback.IEvents interface
func (p *pbEventsMock) QuitPlayingFromBar(pl *models.Playlist) {}

// SeekInStream implements the playback.IEvents interface
func (p *pbEventsMock) SeekInStream(pos int64) {}

// StopAll implements the playback.IEvents interface
func (p *pbEventsMock) StopAll() {}

// StopStream implements the playback.IEvents interface
func (p *pbEventsMock) StopStream() {}

// TryPlayingFromBar implements the playback.IEvents interface
func (p *pbEventsMock) TryPlayingFromBar(pl *models.Playlist, position int) {}
