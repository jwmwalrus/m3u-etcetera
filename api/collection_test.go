package api

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/stretchr/testify/assert"
)

func TestGetCollection(t *testing.T) {
	table := []testCase{
		{
			"Get with ID, invalid",
			"api/collection/get-id-invalid",
			false,
			0,
			&m3uetcpb.GetCollectionRequest{},
			&m3uetcpb.GetCollectionResponse{},
			true,
		},
		{
			"Get with ID, not found",
			"api/collection/get-id-not-found",
			false,
			0,
			&m3uetcpb.GetCollectionRequest{Id: 2},
			&m3uetcpb.GetCollectionResponse{},
			true,
		},
		{
			"Get with ID, success",
			"api/collection/get-id-success",
			false,
			0,
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
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetCollectionResponse)

			res, err := svc.GetCollection(context.Background(), tc.req.(*m3uetcpb.GetCollectionRequest))

			assert.Equal(t, err != nil, tc.err)
			if tc.err {
				return
			}

			assert.Equal(t, res.Collection.Id, exp.Collection.Id)
			assert.Equal(t, res.Collection.Name, exp.Collection.Name)
			assert.Equal(t, res.Collection.Location, exp.Collection.Location)
		})
	}

	return
}

func TestGetAllCollections(t *testing.T) {
	table := []testCase{
		{
			"Get all",
			"api/collection/get-all",
			false,
			0,
			&empty.Empty{},
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
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetAllCollectionsResponse)

			res, err := svc.GetAllCollections(context.Background(), tc.req.(*empty.Empty))

			assert.Equal(t, err != nil, tc.err)
			if tc.err {
				return
			}

			assert.Equal(t, len(res.Collections), len(exp.Collections))
			assert.Equal(t, res.Collections[0].Id, exp.Collections[0].Id)
			assert.Equal(t, res.Collections[0].Name, exp.Collections[0].Name)
			assert.Equal(t, res.Collections[0].Location, exp.Collections[0].Location)
		})
	}

	return
}

func TestAddCollection(t *testing.T) {
	table := []testCase{
		{
			"Add with empty location",
			"api/collection/add-no-loc",
			false,
			0,
			&m3uetcpb.AddCollectionRequest{},
			&m3uetcpb.AddCollectionResponse{},
			true,
		},
		{
			"Add with existing location",
			"api/collection/add-existing-loc",
			false,
			0,
			&m3uetcpb.AddCollectionRequest{
				Name:     "new collection",
				Location: "./data/testing/audio1/",
			},
			&m3uetcpb.AddCollectionResponse{},
			true,
		},
		{
			"Add collection, success",
			"api/collection/add-success",
			false,
			0,
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
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			exp := tc.res.(*m3uetcpb.AddCollectionResponse)

			res, err := svc.AddCollection(context.Background(), tc.req.(*m3uetcpb.AddCollectionRequest))

			assert.Equal(t, err != nil, tc.err)
			if tc.err {
				return
			}

			assert.Equal(t, res.Id, exp.Id)
		})
	}

	return
}

func TestRemoveCollection(t *testing.T) {
	table := []testCase{
		{
			"Remove with ID, invalid",
			"api/collection/rem-id-invalid",
			false,
			0,
			&m3uetcpb.RemoveCollectionRequest{},
			&empty.Empty{},
			true,
		},
		{
			"Remove with ID, not found",
			"api/collection/rem-id-not-found",
			false,
			0,
			&m3uetcpb.RemoveCollectionRequest{Id: 2},
			&empty.Empty{},
			true,
		},
		{
			"Remove with ID, success",
			"api/collection/rem-id-success",
			false,
			0,
			&m3uetcpb.RemoveCollectionRequest{Id: 1},
			&empty.Empty{},
			false,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			_, err := svc.RemoveCollection(context.Background(), tc.req.(*m3uetcpb.RemoveCollectionRequest))

			assert.Equal(t, err != nil, tc.err)
		})
	}

	return
}

func TestScanCollection(t *testing.T) {
	table := []testCase{
		{
			"Scan with ID, invalid",
			"api/collection/scan-id-invalid",
			false,
			0,
			&m3uetcpb.ScanCollectionRequest{},
			&empty.Empty{},
			true,
		},
		{
			"Scan with ID, not found",
			"api/collection/scan-id-not-found",
			false,
			0,
			&m3uetcpb.ScanCollectionRequest{Id: 2},
			&empty.Empty{},
			true,
		},
		{
			"Scan with ID, success",
			"api/collection/scan-id-success",
			false,
			0,
			&m3uetcpb.ScanCollectionRequest{Id: 1},
			&empty.Empty{},
			false,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			_, err := svc.ScanCollection(context.Background(), tc.req.(*m3uetcpb.ScanCollectionRequest))

			assert.Equal(t, err != nil, tc.err)
		})
	}

	return
}

func TestDiscoverCollection(t *testing.T) {
	table := []testCase{
		{
			"Discover, success",
			"api/collection/discover",
			false,
			0,
			&empty.Empty{},
			&empty.Empty{},
			false,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			_, err := svc.DiscoverCollections(context.Background(), tc.req.(*empty.Empty))

			assert.Equal(t, err != nil, tc.err)
		})
	}

	return
}
