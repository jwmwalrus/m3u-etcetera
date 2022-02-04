package impexp

import (
	"os"
	"testing"
)

func TestParsePLS(t *testing.T) {
	pl := &playlist{}
	dep := &PLS{pl}
	f, err := os.Open("../../data/testing/impexp/m3u/pl1.pls")
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
