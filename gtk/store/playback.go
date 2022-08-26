package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/bnp/ing2"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type playbackData struct {
	res *m3uetcpb.SubscribeToPlaybackResponse

	trackID                      int64
	uiSet                        bool
	lastDir                      string
	coverFiles                   []string
	headerbar                    *gtk.HeaderBar
	cover                        *gtk.Image
	logoPixbuf                   *gdk.Pixbuf
	playBtn                      *gtk.ToolButton
	title, artist, source, extra *gtk.Label
	prog                         *gtk.ProgressBar

	mu sync.Mutex
}

const (
	defaultSubtitle = "A playlist-centric music player"
)

var (
	// PbData playback data
	PbData = &playbackData{}
)

func (pbd *playbackData) GetCurrentPlayback() (pb *m3uetcpb.Playback,
	t *m3uetcpb.Track, duration int64, status map[string]bool) {
	pbd.mu.Lock()
	defer pbd.mu.Unlock()

	pb = pbd.res.Playback
	t = pbd.res.Track
	duration = pbd.res.Track.Duration
	status = map[string]bool{
		"is-ready":     pbd.res.IsReady,
		"is-paused":    pbd.res.IsPaused,
		"is-playing":   pbd.res.IsPlaying,
		"is-stopped":   pbd.res.IsStopped,
		"is-streaming": pbd.res.IsStreaming,
	}
	return
}

func (pbd *playbackData) GetSubscriptionID() string {
	pbd.mu.Lock()
	defer pbd.mu.Unlock()

	id := pbd.res.SubscriptionId
	return id
}

func (pbd *playbackData) ProcessSubscriptionResponse(res *m3uetcpb.SubscribeToPlaybackResponse) {
	pbd.mu.Lock()
	defer pbd.mu.Unlock()

	pbd.res = res

	glib.IdleAdd(pbd.updatePlayback)
	glib.IdleAdd(pbd.setCover)
}

func (pbd *playbackData) SetPlaybackUI() (err error) {
	pbd.headerbar, err = builder.GetHeaderBar("window_headerbar")
	if err != nil {
		return
	}

	pbd.cover, err = builder.GetImage("cover")
	if err != nil {
		return
	}

	pbd.logoPixbuf, err = builder.PixbufNewFromFile("images/m3u-etcetera.png")
	if err != nil {
		return
	}

	pbd.title, err = builder.GetLabel("playback_title")
	if err != nil {
		return
	}
	pbd.artist, err = builder.GetLabel("playback_artist")
	if err != nil {
		return
	}
	pbd.source, err = builder.GetLabel("playback_source")
	if err != nil {
		return
	}
	pbd.extra, err = builder.GetLabel("playback_extra")
	if err != nil {
		return
	}
	pbd.prog, err = builder.GetProgressBar("progress")
	if err != nil {
		return
	}
	pbd.playBtn, err = builder.GetToolButton("control_play")
	if err != nil {
		return
	}

	for _, v := range base.Conf.GTK.Playback.CoverFilenames {
		for _, ext := range []string{".jpeg", ".jpg", ".png"} {
			pbd.coverFiles = append(pbd.coverFiles, v+ext)
			pbd.coverFiles = append(pbd.coverFiles,
				cases.Title(language.English).String(v)+ext)
		}
	}

	pbd.uiSet = true
	return
}

func (pbd *playbackData) getTrackID() int64 {
	pbd.mu.Lock()
	defer pbd.mu.Unlock()

	return pbd.trackID
}

func (pbd *playbackData) setCover() bool {
	pbd.mu.Lock()
	defer pbd.mu.Unlock()

	if pbd.res.IsStreaming {
		un, err := urlstr.URLToPath(pbd.res.Playback.Location)
		if err != nil {
			return false
		}
		dir := filepath.Dir(un)
		if dir != pbd.lastDir {
			pbd.lastDir = dir
			fp := ""

			trackCover := pbd.res.Track.Cover
			coverFiles := pbd.coverFiles

			for _, v := range coverFiles {
				dirfile := filepath.Join(pbd.lastDir, v)
				if _, err := os.Stat(dirfile); !os.IsNotExist(err) {
					fp = dirfile
					break
				}
			}

			if fp == "" && trackCover != "" {
				trackCover = filepath.Join(base.CoversDir, trackCover)
				if _, err := os.Stat(trackCover); !os.IsNotExist(err) {
					fp = trackCover
				}
			}

			if fp == "" {
				pbd.cover.SetFromPixbuf(pbd.logoPixbuf)
				return false
			}

			pixbuf, err := gdk.PixbufNewFromFileAtScale(fp, 150, 150, true)
			if err != nil {
				return false
			}
			pbd.cover.SetFromPixbuf(pixbuf)
		}
		return false
	}

	pbd.lastDir = ""
	pbd.cover.SetFromPixbuf(pbd.logoPixbuf)
	return false
}

func (pbd *playbackData) updatePlayback() bool {
	log.Debug("Updating playback")

	iconName := "media-playback-pause"

	pbd.mu.Lock()
	if pbd.res.IsPaused {
		iconName = "media-playback-start"
	}
	pbd.playBtn.SetIconName(iconName)

	var location, title, artist, album string
	var duration, position int64

	oldTrackID := pbd.trackID
	if pbd.res.IsStreaming {
		pbd.trackID = pbd.res.Track.Id
		location = pbd.res.Playback.Location

		title = pbd.res.Track.Title
		artist = pbd.res.Track.Artist
		album = pbd.res.Track.Album
		duration = pbd.res.Track.Duration
		position = pbd.res.Playback.Skip
	} else {
		pbd.trackID = 0
		location = ""
		title, artist, album = "", "", ""
	}

	if duration > 0 {
		pos := time.Duration(position) * time.Nanosecond
		dur := time.Duration(duration) * time.Nanosecond
		pbd.prog.SetFraction(float64(position) / float64(duration))
		pbd.prog.SetText(
			fmt.Sprintf(
				"%v / %v",
				pos.Truncate(time.Second),
				dur.Truncate(time.Second),
			),
		)
	} else {
		pbd.prog.SetFraction(float64(0))
		pbd.prog.SetText("Not Playing")
	}

	maxLen := 45
	subtitle := ing2.TruncateText(title, maxLen)
	if title == "" {
		title = "Not Playing"
	}
	if artist != "" {
		artist = "by " + artist
		if subtitle != "" {
			subtitle += " (" + ing2.TruncateText(artist, maxLen) + ")"
		}
	}
	if album != "" {
		location = "from " + album
	} else {
		path, err := urlstr.URLToPath(location)
		if err == nil {
			location = path
		}
	}

	if subtitle != "" {
		pbd.headerbar.SetSubtitle(subtitle)
	} else {
		pbd.headerbar.SetSubtitle(defaultSubtitle)
	}

	pbd.title.SetText(ing2.TruncateText(title, maxLen))
	pbd.title.SetTooltipText(title)
	pbd.artist.SetText(ing2.TruncateText(artist, maxLen))
	pbd.artist.SetTooltipText(artist)
	pbd.source.SetText(ing2.TruncateText(location, maxLen))
	pbd.source.SetTooltipText(location)
	pbd.mu.Unlock()

	if oldTrackID != pbd.getTrackID() {
		BData.updatePlaybarModel()
	}
	return false
}
