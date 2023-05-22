package impexp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseM3U(t *testing.T) {
	pl := &playlist{}
	dep := &M3U{pl}

	f, err := os.Open("../../data/testing/impexp/m3u/pl1.m3u")
	assert.NoError(t, err)
	defer f.Close()

	err = dep.Parse(f)
	assert.NoError(t, err)

	assert.Len(t, dep.tracks, 2)
	assert.Equal(t, "PL1", dep.Name())

	assert.Equal(t, "Artist 1 - Title 1", dep.tracks[0].ArtistTitle)
	assert.Equal(t, "Artist 2 - Title 1", dep.tracks[1].ArtistTitle)

	assert.Equal(t, 300, int(dep.tracks[0].Duration/1e9))
	assert.Equal(t, 300, int(dep.tracks[1].Duration/1e9))
}
