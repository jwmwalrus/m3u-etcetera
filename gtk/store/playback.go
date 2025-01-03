package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/bnp/chars"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	rtc "github.com/jwmwalrus/rtcycler"
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
	logoPixbuf                   *gdkpixbuf.Pixbuf
	playBtn                      *gtk.ToolButton
	title, artist, source, extra *gtk.Label
	prog                         *gtk.ProgressBar

	mu sync.RWMutex
}

const (
	defaultSubtitle = "A playlist-centric music player"
)

var (
	// PbData playback data.
	PbData = &playbackData{}
)

func (pbd *playbackData) CurrentPlayback() (pb *m3uetcpb.Playback,
	t *m3uetcpb.Track, duration int64, status map[string]bool) {
	pbd.mu.RLock()
	defer pbd.mu.RUnlock()

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

func (pbd *playbackData) SubscriptionID() string {
	pbd.mu.RLock()
	defer pbd.mu.RUnlock()

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

	pbd.logoPixbuf, err = builder.PixbufNewFromFile("images/m3u-etcetera-logo.png")
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
	pbd.mu.RLock()
	defer pbd.mu.RUnlock()

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
				if _, err := os.Stat(dirfile); err == nil {
					fp = dirfile
					break
				}
			}

			if fp == "" && trackCover != "" {
				trackCover = filepath.Join(base.CoversDir(), trackCover)
				if _, err := os.Stat(trackCover); err == nil {
					fp = trackCover
				}
			}

			if fp == "" {
				pbd.cover.SetFromPixbuf(pbd.logoPixbuf)
				return false
			}

			pixbuf, err := gdkpixbuf.NewPixbufFromFileAtScale(fp, 150, 150, true)
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
	rtc.Trace("Updating playback")

	iconName := "media-playback-pause"

	var oldTrackID int64
	pbd.mu.Lock()
	{
		if pbd.res.IsPaused {
			iconName = "media-playback-start"
		}
		pbd.playBtn.SetIconName(iconName)

		var location, title, artist, album string
		var duration, position int64

		oldTrackID = pbd.trackID
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
		subtitle := chars.Truncate(title, maxLen)
		if title == "" {
			title = "Not Playing"
		}
		if artist != "" {
			artist = "by " + artist
			if subtitle != "" {
				subtitle += " (" + chars.Truncate(artist, maxLen) + ")"
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

		pbd.title.SetText(chars.Truncate(title, maxLen))
		pbd.title.SetTooltipText(title)
		pbd.artist.SetText(chars.Truncate(artist, maxLen))
		pbd.artist.SetTooltipText(artist)
		pbd.source.SetText(chars.Truncate(location, maxLen))
		pbd.source.SetTooltipText(location)
	}
	pbd.mu.Unlock()

	if oldTrackID != pbd.getTrackID() {
		BData.updatePlaybarModel()
	}
	return false
}
