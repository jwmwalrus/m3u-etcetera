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

// PLS implementation.
type PLS struct {
	*playlist
}

// Format implements the PlaylistDef interface.
func (pi *PLS) Format(w io.StringWriter) (n int, err error) {
	out := strings.Builder{}
	_, err = out.WriteString("[playlist]\n\n")
	if err != nil {
		return
	}

	for k, v := range pi.tracks {
		un := v.Location
		if strings.HasPrefix(v.Location, "file://") {
			un, err = urlstr.URLToPath(v.Location)
			if err != nil {
				return
			}
		}
		_, err = out.WriteString(fmt.Sprintf("File%d=%s\n", k+1, un))
		if err != nil {
			return
		}

		_, err = out.WriteString(fmt.Sprintf("Title%d=%s\n", k+1, v.Title))
		if err != nil {
			return
		}

		dur := time.Duration(v.Duration) * time.Nanosecond
		dur = dur.Truncate(time.Second)

		_, err = out.WriteString(
			fmt.Sprintf("Length%d=%v\n\n", k+1, dur.Seconds()),
		)
		if err != nil {
			return
		}
	}

	_, err = out.WriteString(
		fmt.Sprintf("NumberOfEntries=%v\n", len(pi.tracks)),
	)
	if err != nil {
		return
	}

	n, err = w.WriteString(out.String())
	return
}

// Parse implements the PlaylistDef interface.
func (pi *PLS) Parse(f io.Reader) (err error) {
	bv, err := io.ReadAll(f)
	if err != nil {
		return
	}

	s := strings.ReplaceAll(string(bv), "\r", "")
	lines := strings.Split(s, "\n")

	headerre := regexp.MustCompile(`^(\[playlist\])$`)
	trackre := regexp.MustCompile(`^(File|Title|Length)(\d+)=\s*(.+)$`)
	footerre := regexp.MustCompile(`^(NumberOfEntries|Version)=\s*(.+)$`)

	ti := []TrackInfo{}
	i := -1
	hasHeader := false
	hasTrack := false
	hasFooter := false
	for {
		i++
		if i >= len(lines) {
			break
		}
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}

		lines[i] = strings.TrimSpace(lines[i])
		if match := headerre.FindStringSubmatch(lines[i]); len(match) > 1 {
			if i != 0 {
				err = fmt.Errorf(
					"the header directive is not the first line of the file",
				)
				return
			}
			hasHeader = true
		} else if match := trackre.FindStringSubmatch(lines[i]); len(match) > 1 {
			var idx int64
			idx, err = strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				return
			}

			for len(ti) < int(idx) {
				ti = append(ti, TrackInfo{})
			}

			switch match[1] {
			case "File":

				var u string
				u, err = getURL(match[3])
				if err != nil {
					return
				}
				ti[idx-1].Location = u
			case "Title":
				ti[idx-1].Title = match[3]
			case "Length":
				var dur int64
				dur, err = strconv.ParseInt(match[3], 10, 64)
				ti[idx-1].Duration = dur * 1e9
			}
			hasTrack = true
		} else if match := footerre.FindStringSubmatch(lines[i]); len(match) > 1 {
			switch match[1] {
			case "NumberOfEntries":
				var n int64
				n, err = strconv.ParseInt(match[2], 10, 64)
				if err != nil {
					return
				}
				if int(n) != len(ti) {
					err = fmt.Errorf("Number of entries does not match")
				}
			case "Version":
			}
			hasFooter = true
		} else {
			err = fmt.Errorf("Line %v does not correspond to a PLS file", i+1)
			return
		}
	}

	if !hasHeader {
		err = fmt.Errorf("Required PLS header not found")
		return
	}

	if !hasTrack {
		err = fmt.Errorf("Required PLS track entries not found")
		return
	}

	if !hasFooter {
		err = fmt.Errorf("Required PLS footer not found")
		return
	}

	for k, v := range ti {
		if v.Location == "" {
			err = fmt.Errorf("Mandatory location is missing for track #%d", k+1)
			return
		}
	}

	pi.tracks = ti

	return
}

// Type implements the PlaylistDef interface.
func (*PLS) Type() string {
	return PLSPlaylist.String()
}
