package api

import (
	"context"
	"testing"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestGetCollection(t *testing.T) {
	table := []testCase{
		{
			"Get with ID, invalid",
			"api/collection/get-id-invalid",
			&m3uetcpb.GetCollectionRequest{},
			&m3uetcpb.GetCollectionResponse{},
			true,
		},
		{
			"Get with ID, not found",
			"api/collection/get-id-not-found",
			&m3uetcpb.GetCollectionRequest{Id: 2},
			&m3uetcpb.GetCollectionResponse{},
			true,
		},
		{
			"Get with ID, success",
			"api/collection/get-id-success",
			&m3uetcpb.GetCollectionRequest{Id: 1},
			&m3uetcpb.GetCollectionResponse{
				Collection: &m3uetcpb.Collection{
					Id:       1,
					Name:     "local:audio",
					Location: "./data/testing/audio1/",
				},
			},
			false,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetCollectionResponse)

			res, err := svc.GetCollection(context.Background(), tc.req.(*m3uetcpb.GetCollectionRequest))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, exp.Collection.Id, res.Collection.Id)
			assert.Equal(t, exp.Collection.Name, res.Collection.Name)
			assert.Equal(t, exp.Collection.Location, res.Collection.Location)
		})
	}

	return
}

func TestGetAllCollections(t *testing.T) {
	table := []testCase{
		{
			"Get all",
			"api/collection/get-all",
			&m3uetcpb.Empty{},
			&m3uetcpb.GetAllCollectionsResponse{
				Collections: []*m3uetcpb.Collection{{
					Id:       1,
					Name:     "local:audio",
					Location: "./data/testing/audio1/",
				}},
			},
			false,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetAllCollectionsResponse)

			res, err := svc.GetAllCollections(context.Background(), tc.req.(*m3uetcpb.Empty))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(exp.Collections), len(res.Collections))
			assert.Equal(t, exp.Collections[0].Id, res.Collections[0].Id)
			assert.Equal(t, exp.Collections[0].Name, res.Collections[0].Name)
			assert.Equal(t, exp.Collections[0].Location, res.Collections[0].Location)
		})
	}

	return
}

func TestAddCollection(t *testing.T) {
	table := []testCase{
		{
			"Add with empty location",
			fixturesDir("api/collection/add-no-loc"),
			&m3uetcpb.AddCollectionRequest{},
			&m3uetcpb.AddCollectionResponse{},
			true,
		},
		{
			"Add with existing location",
			fixturesDir("api/collection/add-existing-loc"),
			&m3uetcpb.AddCollectionRequest{
				Name:     "new collection",
				Location: "./data/testing/audio1/",
			},
			&m3uetcpb.AddCollectionResponse{},
			true,
		},
		{
			"Add collection, success",
			fixturesDir("api/collection/add-success"),
			&m3uetcpb.AddCollectionRequest{
				Name:     "new collection",
				Location: "./data/testing/audio2/",
			},
			&m3uetcpb.AddCollectionResponse{
				Id: 2,
			},
			false,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			exp := tc.res.(*m3uetcpb.AddCollectionResponse)

			res, err := svc.AddCollection(context.Background(), tc.req.(*m3uetcpb.AddCollectionRequest))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, exp.Id, res.Id)
		})
	}

	return
}

func TestRemoveCollection(t *testing.T) {
	table := []testCase{
		{
			"Remove with ID, invalid",
			"api/collection/rem-id-invalid",
			&m3uetcpb.RemoveCollectionRequest{},
			&m3uetcpb.Empty{},
			true,
		},
		{
			"Remove with ID, not found",
			"api/collection/rem-id-not-found",
			&m3uetcpb.RemoveCollectionRequest{Id: 2},
			&m3uetcpb.Empty{},
			true,
		},
		{
			"Remove with ID, success",
			"api/collection/rem-id-success",
			&m3uetcpb.RemoveCollectionRequest{Id: 1},
			&m3uetcpb.Empty{},
			false,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			_, err := svc.RemoveCollection(context.Background(), tc.req.(*m3uetcpb.RemoveCollectionRequest))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}

	return
}

func TestScanCollection(t *testing.T) {
	table := []testCase{
		{
			"Scan with ID, invalid",
			"api/collection/scan-id-invalid",
			&m3uetcpb.ScanCollectionRequest{},
			&m3uetcpb.Empty{},
			true,
		},
		{
			"Scan with ID, not found",
			"api/collection/scan-id-not-found",
			&m3uetcpb.ScanCollectionRequest{Id: 2},
			&m3uetcpb.Empty{},
			true,
		},
		{
			"Scan with ID, success",
			"api/collection/scan-id-success",
			&m3uetcpb.ScanCollectionRequest{Id: 1},
			&m3uetcpb.Empty{},
			false,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			_, err := svc.ScanCollection(context.Background(), tc.req.(*m3uetcpb.ScanCollectionRequest))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}

	return
}

func TestDiscoverCollection(t *testing.T) {
	table := []testCase{
		{
			"Discover, success",
			"api/collection/discover",
			&m3uetcpb.Empty{},
			&m3uetcpb.Empty{},
			false,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			tests.SetupTest(t, tc.fixturesDir)
			t.Cleanup(func() { tests.TeardownTest(t) })

			_, err := svc.DiscoverCollections(context.Background(), tc.req.(*m3uetcpb.Empty))

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}

	return
}

func TestCollectionToProtobuf(t *testing.T) {
	c := models.Collection{
		ID:             1,
		Name:           "test",
		Location:       "file://somewhere",
		Remotelocation: "file://somewhere-else",
		Remote:         true,
		Scanned:        100,
		Perspective:    models.Perspective{Idx: 1},
		CreatedAt:      time.Now().UnixNano(),
		UpdatedAt:      time.Now().UnixNano(),
	}

	out := c.ToProtobuf()

	cpb, ok := out.(*m3uetcpb.Collection)
	assert.True(t, ok)

	assert.Equal(t, c.ID, cpb.Id)
	assert.Equal(t, c.Name, cpb.Name)
	assert.Equal(t, c.Location, cpb.Location)
	assert.Equal(t, c.Remotelocation, cpb.RemoteLocation)
	assert.Equal(t, c.Remote, cpb.Remote)
	assert.Equal(t, int32(c.Scanned), cpb.Scanned)
	assert.Equal(t, m3uetcpb.Perspective(c.Perspective.Idx), cpb.Perspective)
	assert.Equal(t, c.CreatedAt, cpb.CreatedAt)
	assert.Equal(t, c.UpdatedAt, cpb.UpdatedAt)
}

func TestTrackToProtobuf(t *testing.T) {
	track := models.Track{
		ID:           1,
		Location:     "file://some-location",
		Title:        "song title",
		Album:        "tiniest hits",
		Artist:       "somebody, someone",
		Albumartist:  "Various",
		Genre:        "A New One",
		Year:         time.Now().Year(),
		Tracknumber:  1,
		Tracktotal:   10,
		Discnumber:   1,
		Disctotal:    1,
		Duration:     270000000000,
		Playcount:    13,
		Lastplayed:   time.Now().UnixNano(),
		CollectionID: 1,
		CreatedAt:    time.Now().UnixNano(),
		UpdatedAt:    time.Now().UnixNano(),
	}

	out := track.ToProtobuf()
	tpb, ok := out.(*m3uetcpb.Track)
	assert.True(t, ok)

	assert.Equal(t, track.ID, tpb.Id)
	assert.Equal(t, track.Location, tpb.Location)
	assert.Equal(t, track.Title, tpb.Title)
	assert.Equal(t, track.Album, tpb.Album)
	assert.Equal(t, track.Artist, tpb.Artist)
	assert.Equal(t, track.Albumartist, tpb.Albumartist)
	assert.Equal(t, track.Genre, tpb.Genre)
	assert.Equal(t, int32(track.Year), tpb.Year)
	assert.Equal(t, int32(track.Tracknumber), tpb.Tracknumber)
	assert.Equal(t, int32(track.Tracktotal), tpb.Tracktotal)
	assert.Equal(t, int32(track.Discnumber), tpb.Discnumber)
	assert.Equal(t, int32(track.Disctotal), tpb.Disctotal)
	assert.Equal(t, track.Duration, tpb.Duration)
	assert.Equal(t, int32(track.Playcount), tpb.Playcount)
	assert.Equal(t, track.Lastplayed, tpb.Lastplayed)
	assert.Equal(t, track.CollectionID, tpb.CollectionId)
	assert.Equal(t, track.CreatedAt, tpb.CreatedAt)
	assert.Equal(t, track.UpdatedAt, tpb.UpdatedAt)

	assert.True(t, tpb.Dangling)
}
