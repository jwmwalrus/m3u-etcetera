package api

import (
	"testing"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestQueryToProtobuf(t *testing.T) {
	tests.SetupTest(t, fixturesDir("api/query/query-to-protobuf"))
	t.Cleanup(func() { tests.TeardownTest(t) })

	qy := models.Query{}
	qy.Read(1)

	out := qy.ToProtobuf()

	qypb, ok := out.(*m3uetcpb.Query)
	assert.True(t, ok)

	assert.Equal(t, qy.ID, qypb.Id)
	assert.Equal(t, qy.Name, qypb.Name)
	assert.Equal(t, qy.Description, qypb.Description)
	assert.Equal(t, qy.Random, qypb.Random)
	assert.Equal(t, int32(qy.Rating), qypb.Rating)
	assert.Equal(t, int32(qy.Limit), qypb.Limit)
	assert.Equal(t, qy.From, qypb.From)
	assert.Equal(t, qy.To, qypb.To)
	assert.Equal(t, qy.CreatedAt, qypb.CreatedAt)
	assert.Equal(t, qy.UpdatedAt, qypb.UpdatedAt)
	assert.True(t, qypb.ReadOnly)
}
