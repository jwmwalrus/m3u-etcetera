package impexp

import (
	"os"
	"testing"
)

func TestParseM3U(t *testing.T) {
	pl := &playlist{}
	dep := &M3U{pl}
	f, err := os.Open("../../data/testing/impexp/m3u/pl1.m3u")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	err = dep.Parse(f)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("Playlist 1")
	for k, v := range dep.tracks {
		t.Logf("\t%d -> %+v", k, v)
	}

}
