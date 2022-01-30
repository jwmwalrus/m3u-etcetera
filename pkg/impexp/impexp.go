package impexp

import (
	"fmt"
	"io"
	"path/filepath"
)

type PlaylistType int

const (
	M3UPlaylist PlaylistType = iota
	PLSPlaylist
)

func (plt PlaylistType) String() string {
	return []string{
		"m3u",
		"pls",
	}[plt]
}

type PlaylistData interface {
	Add(ti []TrackInfo)
	Name() string
	Reset()
	Tracks() []TrackInfo
}

type PlaylistDef interface {
	PlaylistData
	Type() string
	Parse(io.Reader) error
	Format() (io.Writer, error)
}

var (
	extToType = map[string]PlaylistType{
		".m3u":  M3UPlaylist,
		".m3u8": M3UPlaylist,
		".pls":  PLSPlaylist,
	}
)

func New(plt PlaylistType, props ...PlaylistProp) (PlaylistDef, error) {
	pl := &playlist{}
	pl.setProps(props)

	switch plt {
	case M3UPlaylist:
		return &M3U{playlist: pl}, nil
	case PLSPlaylist:
		return &PLS{playlist: pl}, nil
	default:
		return nil, nil
	}
}

func NewFromPath(path string, props ...PlaylistProp) (PlaylistDef, error) {
	plt, ok := extToType[filepath.Ext(path)]
	if !ok {
		return nil, fmt.Errorf("Unsupported playlist type")
	}

	return New(plt, props...)
}
