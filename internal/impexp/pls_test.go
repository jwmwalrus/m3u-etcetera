package impexp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePLS(t *testing.T) {
	pl := &playlist{}
	dep := &PLS{pl}

	f, err := os.Open("../../data/testing/impexp/pls/pl1.pls")
	assert.NoError(t, err)
	defer f.Close()

	err = dep.Parse(f)
	assert.NoError(t, err)

	assert.Len(t, dep.tracks, 2)

	assert.Equal(t, "Title 1", dep.tracks[0].Title)
	assert.Equal(t, "Title 2", dep.tracks[1].Title)

	assert.Equal(t, 300, int(dep.tracks[0].Duration/1e9))
	assert.Equal(t, 300, int(dep.tracks[1].Duration/1e9))
}
