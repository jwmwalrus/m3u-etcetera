package api

import (
	"testing"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestPlaylistToProtobuf(t *testing.T) {
	tests.SetupTest(t, fixturesDir("api/playbar/playlist-to-protobuf"))
	t.Cleanup(func() { tests.TeardownTest(t) })

	pl := models.Playlist{}
	pl.Read(1)

	out := pl.ToProtobuf()

	plpb, ok := out.(*m3uetcpb.Playlist)
	assert.True(t, ok)

	assert.Equal(t, pl.ID, plpb.Id)
	assert.Equal(t, pl.Name, plpb.Name)
	assert.Equal(t, pl.Description, plpb.Description)
	assert.Equal(t, pl.Open, plpb.Open)
	assert.Equal(t, pl.Active, plpb.Active)
	assert.Equal(t, pl.Transient, plpb.Transient)
	assert.Equal(t, pl.CreatedAt, plpb.CreatedAt)
	assert.Equal(t, pl.UpdatedAt, plpb.UpdatedAt)
	assert.Equal(t, pl.PlaylistGroupID, plpb.PlaylistGroupId)
}

func TestPlaylistGroupToProtobuf(t *testing.T) {
	pg := models.PlaylistGroup{
		ID:            1,
		Idx:           1,
		Name:          "playlist group 1",
		Description:   "playlist group 1 description",
		CreatedAt:     time.Now().UnixNano(),
		UpdatedAt:     time.Now().UnixNano(),
		PerspectiveID: 1,
	}

	out := pg.ToProtobuf()

	pgpb, ok := out.(*m3uetcpb.PlaylistGroup)
	assert.True(t, ok)

	assert.Equal(t, pg.ID, pgpb.Id)
	assert.Equal(t, pg.Name, pgpb.Name)
	assert.Equal(t, pg.Description, pgpb.Description)
	assert.Equal(t, pg.CreatedAt, pgpb.CreatedAt)
	assert.Equal(t, pg.UpdatedAt, pgpb.UpdatedAt)
}

func TestPlaylistTrackToProtobuf(t *testing.T) {
	pt := models.PlaylistTrack{
		ID:            1,
		Position:      1,
		Dynamic:       true,
		CreatedAt:     time.Now().UnixNano(),
		UpdatedAt:     time.Now().UnixNano(),
		PlaylistID:    1,
		TrackID:       1,
		Lastplayedfor: 500000000000,
	}

	out := pt.ToProtobuf()

	ptpb, ok := out.(*m3uetcpb.PlaylistTrack)
	assert.True(t, ok)

	assert.Equal(t, pt.ID, ptpb.Id)
	assert.Equal(t, int32(pt.Position), ptpb.Position)
	assert.Equal(t, pt.Dynamic, ptpb.Dynamic)
	assert.Equal(t, pt.Lastplayedfor, ptpb.Lastplayedfor)
	assert.Equal(t, pt.CreatedAt, ptpb.CreatedAt)
	assert.Equal(t, pt.UpdatedAt, ptpb.UpdatedAt)
	assert.Equal(t, pt.PlaylistID, ptpb.PlaylistId)
	assert.Equal(t, pt.TrackID, ptpb.TrackId)
}
