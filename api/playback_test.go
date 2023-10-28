package api

import (
	"context"
	"testing"
	"time"

	"github.com/go-gst/go-gst/gst"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestGetPlayback(t *testing.T) {
	table := []testCase{
		{
			"Get with no playback",
			"api/playback/get-nopb",
			&m3uetcpb.Empty{},
			&m3uetcpb.GetPlaybackResponse{},
			false,
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
			false,
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
			false,
		},
	}

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

	svc := PlaybackSvc{PbEvents: &events}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetPlaybackResponse)

			res, err := svc.GetPlayback(context.Background(), tc.req.(*m3uetcpb.Empty))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if exp.IsPlaying {
				assert.Equal(t, true, res.Playback != nil)
				assert.Equal(t, exp.Playback.Id, res.Playback.Id)
				assert.Equal(t, exp.Playback.Location, res.Playback.Location)
				assert.Equal(t, true, res.Track != nil)
				assert.Equal(t, exp.Track.Id, res.Track.Id)
				assert.Equal(t, exp.Track.Location, res.Track.Location)
				assert.Equal(t, exp.Track.Title, res.Track.Title)
				assert.Equal(t, exp.Track.Album, res.Track.Album)
				assert.Equal(t, exp.Track.Artist, res.Track.Artist)
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
			true,
		},
		{
			"Execute play invalid ID",
			"api/playback/exec-invalid-id",
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
			&m3uetcpb.ExecutePlaybackActionRequest{
				Action: m3uetcpb.PlaybackAction_PB_NEXT,
			},
			&m3uetcpb.Empty{},
			false,
		},
	}

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
		isPlaying:     false,
		hasNextStream: true,
	}

	svc := PlaybackSvc{PbEvents: &events}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			_, err := svc.ExecutePlaybackAction(context.Background(), tc.req.(*m3uetcpb.ExecutePlaybackActionRequest))

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}

	return
}

func TestPlaybackToProtobuf(t *testing.T) {
	pb := models.Playback{
		ID:        1,
		Location:  "file://somewhere",
		Played:    true,
		Skip:      500,
		CreatedAt: time.Now().UnixNano(),
		UpdatedAt: time.Now().UnixNano(),
		TrackID:   1,
	}

	out := pb.ToProtobuf()

	pbpb, ok := out.(*m3uetcpb.Playback)
	assert.True(t, ok)

	assert.Equal(t, pb.ID, pbpb.Id)
	assert.Equal(t, pb.Location, pbpb.Location)
	assert.Equal(t, pb.Played, pbpb.Played)
	assert.Equal(t, pb.Skip, pbpb.Skip)
	assert.Equal(t, pb.CreatedAt, pbpb.CreatedAt)
	assert.Equal(t, pb.UpdatedAt, pbpb.UpdatedAt)
	assert.Equal(t, pb.TrackID, pbpb.TrackId)
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

func (e *pbEventsMock) GetPlayback() (pb *models.Playback, t *models.Track) { return e.pb, e.t }

func (e *pbEventsMock) GetState() gst.State { return e.state }

func (e *pbEventsMock) HasNextStream() bool { return e.hasNextStream }

func (e *pbEventsMock) IsPaused() bool { return e.isPaused }

func (e *pbEventsMock) IsPlaying() bool { return e.isPlaying }

func (e *pbEventsMock) IsReady() bool { return e.isReady }

func (e *pbEventsMock) IsStreaming() bool { return e.isStreaming }

func (e *pbEventsMock) IsStopped() bool { return e.isStopped }

func (e *pbEventsMock) NextStream() (err error) { return e.nextStreamErr }

func (e *pbEventsMock) PauseStream(off bool) (err error) { return e.pauseStreamErr }

func (p *pbEventsMock) PlayStreams(force bool, locations []string, ids []int64) {}

func (p *pbEventsMock) PreviousStream() {}

func (p *pbEventsMock) QuitPlayingFromBar(pl *models.Playlist) {}

func (p *pbEventsMock) SeekInStream(pos int64) {}

func (p *pbEventsMock) StopAll() {}

func (p *pbEventsMock) StopStream() {}

func (p *pbEventsMock) TryPlayingFromBar(pl *models.Playlist, position int) {}
