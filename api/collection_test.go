package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/stretchr/testify/assert"
)

func TestGetCollection(t *testing.T) {
	table := []testCase{
		{
			"Get with ID, invalid",
			"api/collection/get-id-invalid",
			&m3uetcpb.GetCollectionRequest{},
			&m3uetcpb.GetCollectionResponse{},
			fmt.Errorf("error"),
		},
		{
			"Get with ID, not found",
			"api/collection/get-id-not-found",
			&m3uetcpb.GetCollectionRequest{Id: 2},
			&m3uetcpb.GetCollectionResponse{},
			fmt.Errorf("error"),
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
			nil,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetCollectionResponse)

			res, err := svc.GetCollection(context.Background(), tc.req.(*m3uetcpb.GetCollectionRequest))

			if tc.err != nil {
				assert.NotNil(t, err)
			}
			if tc.err == nil {
				assert.Nil(t, err)
				assert.Equal(t, res.Collection.Id, exp.Collection.Id)
				assert.Equal(t, res.Collection.Name, exp.Collection.Name)
				assert.Equal(t, res.Collection.Location, exp.Collection.Location)
			}
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
			nil,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			exp := tc.res.(*m3uetcpb.GetAllCollectionsResponse)

			res, err := svc.GetAllCollections(context.Background(), tc.req.(*m3uetcpb.Empty))

			if tc.err != nil {
				assert.NotNil(t, err)
			}
			if tc.err == nil {
				assert.Nil(t, err)
				assert.Equal(t, len(res.Collections), len(exp.Collections))
				assert.Equal(t, res.Collections[0].Id, exp.Collections[0].Id)
				assert.Equal(t, res.Collections[0].Name, exp.Collections[0].Name)
				assert.Equal(t, res.Collections[0].Location, exp.Collections[0].Location)
			}
		})
	}

	return
}

func TestAddCollection(t *testing.T) {
	table := []testCase{
		{
			"Add with empty location",
			"api/collection/add-no-loc",
			&m3uetcpb.AddCollectionRequest{},
			&m3uetcpb.AddCollectionResponse{},
			fmt.Errorf("error"),
		},
		{
			"Add with existing location",
			"api/collection/add-existing-loc",
			&m3uetcpb.AddCollectionRequest{
				Name:     "new collection",
				Location: "./data/testing/audio1/",
			},
			&m3uetcpb.AddCollectionResponse{},
			fmt.Errorf("error"),
		},
		{
			"Add collection, success",
			"api/collection/add-success",
			&m3uetcpb.AddCollectionRequest{
				Name:     "new collection",
				Location: "./data/testing/audio2/",
			},
			&m3uetcpb.AddCollectionResponse{
				Id: 2,
			},
			nil,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			exp := tc.res.(*m3uetcpb.AddCollectionResponse)

			res, err := svc.AddCollection(context.Background(), tc.req.(*m3uetcpb.AddCollectionRequest))

			if tc.err != nil {
				assert.NotNil(t, err)
			}
			if tc.err == nil {
				assert.Nil(t, err)
				assert.Equal(t, res.Id, exp.Id)
			}
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
			fmt.Errorf("error"),
		},
		{
			"Remove with ID, not found",
			"api/collection/rem-id-not-found",
			&m3uetcpb.RemoveCollectionRequest{Id: 2},
			&m3uetcpb.Empty{},
			fmt.Errorf("error"),
		},
		{
			"Remove with ID, success",
			"api/collection/rem-id-success",
			&m3uetcpb.RemoveCollectionRequest{Id: 1},
			&m3uetcpb.Empty{},
			nil,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			_, err := svc.RemoveCollection(context.Background(), tc.req.(*m3uetcpb.RemoveCollectionRequest))

			if tc.err != nil {
				assert.NotNil(t, err)
			}
			if tc.err == nil {
				assert.Nil(t, err)
			}
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
			fmt.Errorf("error"),
		},
		{
			"Scan with ID, not found",
			"api/collection/scan-id-not-found",
			&m3uetcpb.ScanCollectionRequest{Id: 2},
			&m3uetcpb.Empty{},
			fmt.Errorf("error"),
		},
		{
			"Scan with ID, success",
			"api/collection/scan-id-success",
			&m3uetcpb.ScanCollectionRequest{Id: 1},
			&m3uetcpb.Empty{},
			nil,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			_, err := svc.ScanCollection(context.Background(), tc.req.(*m3uetcpb.ScanCollectionRequest))

			if tc.err != nil {
				assert.NotNil(t, err)
			}
			if tc.err == nil {
				assert.Nil(t, err)
			}
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
			nil,
		},
	}

	svc := CollectionSvc{}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			setupTest(t, tc)
			t.Cleanup(func() { teardownTest(t) })

			_, err := svc.DiscoverCollections(context.Background(), tc.req.(*m3uetcpb.Empty))

			if tc.err != nil {
				assert.NotNil(t, err)
			}
			if tc.err == nil {
				assert.Nil(t, err)
			}
		})
	}

	return
}
