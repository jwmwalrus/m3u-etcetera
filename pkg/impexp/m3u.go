package impexp

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jwmwalrus/bnp/urlstr"
)

type M3U struct {
	*playlist
}

func (mi *M3U) Format(w io.StringWriter) (n int, err error) {
	out := strings.Builder{}
	_, err = out.WriteString("#EXTM3U\n")
	if err != nil {
		return
	}

	if mi.name != "" {
		_, err = out.WriteString(
			fmt.Sprintf(
				"#PLAYLIST: %v\n",
				mi.name,
			),
		)
		if err != nil {
			return
		}
	}

	for _, t := range mi.tracks {
		un := t.Location
		if strings.HasPrefix(t.Location, "file://") {
			un, err = urlstr.URLToPath(t.Location)
			if err != nil {
				return
			}
		}

		dur := time.Duration(t.Duration) * time.Nanosecond
		dur = dur.Truncate(time.Second)
		_, err = out.WriteString(
			fmt.Sprintf(
				"#EXTINF:%v,%v\n",
				dur.Seconds(),
				t.ArtistTitle,
			),
		)
		if err != nil {
			return
		}

		_, err = out.WriteString(un + "\n")
		if err != nil {
			return
		}
	}

	n, err = w.WriteString(out.String())
	return
}

func (mi *M3U) Parse(f io.Reader) (err error) {
	bv, err := io.ReadAll(f)
	if err != nil {
		return
	}

	s := strings.ReplaceAll(string(bv), "\r", "")
	lines := strings.Split(s, "\n")

	directre := regexp.MustCompile(`^(#EXT\\w{3}\\b|#PLAYLIST\\b|#EXTGENRE\\b)`)
	extinf := regexp.MustCompile(`^(\\d+),\\s*([^,]*)(,.*)?`)
	extalb := regexp.MustCompile(`^(.+)\\s*([(]\\d{4}[)])?\\s*$`)

	var encoding, name string
	ti := []TrackInfo{{}}
	i := -1
	piv := 0
	validM3U := false
	for {
		i++
		if i >= len(lines) {
			break
		}
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}

		lines[i] = strings.TrimSpace(lines[i])

		if match := directre.FindStringSubmatch(lines[i]); len(match) > 1 {
			lines[i] = strings.TrimPrefix(lines[i], match[1])
			lines[i] = strings.TrimPrefix(lines[i], ":")
			switch match[1] {
			case "#EXTM3U":
				if i != 0 {
					err = fmt.Errorf(
						"The %v directive is not at the first line",
						match[1],
					)
					return
				}
				validM3U = true
			case "#EXTENC":
				if i != 1 {
					err = fmt.Errorf(
						"The %v directive is not at the second line",
						match[1],
					)
					return
				}
				encoding = lines[i]
			case "#EXTINF":
				subm := extinf.FindStringSubmatch(lines[i])
				if len(subm) < 2 {
					err = fmt.Errorf("Invalid %s at line %d", match[1], i+1)
					return
				}

				var sec int64
				sec, err = strconv.ParseInt(subm[1], 10, 64)
				if err != nil {
					return
				}
				ti[piv].Duration = sec * 1e9

				if len(subm) > 3 {
					ti[piv].ArtistTitle = strings.TrimSpace(subm[2])
				}
			case "#PLAYLIST":
				name = strings.TrimSpace(lines[i])
			case "#EXTALB":
				subm := extalb.FindStringSubmatch(lines[i])
				if len(subm) < 2 {
					err = fmt.Errorf("Invalid %s at line %d", match[1], i+1)
				}

				ti[piv].Album = subm[1]

				if len(subm) > 2 {
					year, err2 := strconv.ParseInt(subm[2], 10, 64)
					if err2 == nil {
						ti[piv].Year = int(year)
					}
				}
			case "#EXTART":
				ti[piv].Albumartist = lines[i]
			case "#EXTGENRE":
				ti[piv].Genre = lines[i]

			case "#EXTGRP",
				"#EXTM3A",
				"#EXTBYT",
				"#EXTBIN",
				"#EXTIMG":
			}
			continue
		} else if strings.HasPrefix(lines[i], "#") {
			continue
		}

		var u string
		u, err = getURL(lines[i])
		if err != nil {
			return
		}

		ti[piv].Location = u
		ti = append(ti, TrackInfo{})
		piv = len(ti) - 1
	}

	if !validM3U {
		err = fmt.Errorf("Required M3U directive not found")
		return
	}

	mi.encoding = encoding
	mi.name = name
	mi.tracks = ti[:len(ti)-1]

	return
}

func (*M3U) Type() string {
	return M3UPlaylist.String()
}

func getURL(s string) (string, error) {
	u := s

	if strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") {
		return s, nil
	} else if strings.HasPrefix(s, "file://") {
		u = strings.TrimPrefix(u, "file://")
	}

	return urlstr.PathToURL(u)
}
