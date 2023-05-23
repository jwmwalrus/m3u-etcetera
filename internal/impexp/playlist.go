package impexp

import "strconv"

type playlist struct {
	name     string
	encoding string
	tracks   []TrackInfo
}

func (pl *playlist) Add(ti []TrackInfo) {
	pl.tracks = append(pl.tracks, ti...)
}

func (pl *playlist) Name() string {
	return pl.name
}

func (pl *playlist) Reset() {
	pl.name = ""
	pl.encoding = ""
	pl.tracks = []TrackInfo{}
}

func (pl *playlist) Tracks() []TrackInfo {
	return pl.tracks
}

func (pl *playlist) setProps(props PlaylistProps) {
	for k, v := range props {
		switch k {
		case NamePropKey:
			pl.name = v
		case EncodingPropKey:
			pl.name = v
		}
	}
}

// PlaylistPropKey playlist property key type.
type PlaylistPropKey int

// PlaylistPropKey enum.
const (
	NamePropKey PlaylistPropKey = iota
	EncodingPropKey
)

func (ppk PlaylistPropKey) String() string {
	return []string{"name", "encoding"}[ppk]
}

// PlaylistProp playlist property.
type PlaylistProps map[PlaylistPropKey]string

// TrackInfo track information.
type TrackInfo struct {
	Location    string
	Title       string
	ArtistTitle string
	Album       string
	Artist      string
	Albumartist string
	Genre       string
	Duration    int64
	Year        int
}

// ToRaw converts returns the raw (i.e., frame) information about the track.
func (ti *TrackInfo) ToRaw() (raw map[string]interface{}) {
	raw = map[string]interface{}{}

	if ti.Title != "" {
		raw["TIT1"] = ti.Title
	}

	if ti.ArtistTitle != "" {
		raw["TIT3"] = ti.ArtistTitle
	}

	if ti.Album != "" {
		raw["TALB"] = ti.Album
	}

	if ti.Artist != "" {
		raw["TPE1"] = ti.Artist
	}

	if ti.Albumartist != "" {
		raw["TPE2"] = ti.Albumartist
	}

	if ti.Genre != "" {
		raw["TCON"] = ti.Genre
	}

	if ti.Duration > 0 {
		raw["TLEN"] = strconv.FormatInt(ti.Duration/1e6, 10)
	}

	if ti.Year > 0 {
		raw["TYER"] = strconv.FormatInt(int64(ti.Year), 10)
		raw["TDRC"] = raw["TYER"]
	}
	return
}
