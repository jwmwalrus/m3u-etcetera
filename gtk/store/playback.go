package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/bnp/stringing"
	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type playbackData struct {
	mu  sync.Mutex
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
}

const (
	defaultSubtitle = "A playlist-centric music player"
)

var (
	pbdata = &playbackData{}
)

// ExecutePlaybackAction -
func ExecutePlaybackAction(req *m3uetcpb.ExecutePlaybackActionRequest) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	return
}

func subscribeToPlayback() {
	log.Info("Subscribing to playback")

	defer wgplayback.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	stream, err := cl.SubscribeToPlayback(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Errorf("Error subscribing to playback: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		pbdata.mu.Lock()
		pbdata.res = res
		pbdata.mu.Unlock()
		glib.IdleAdd(pbdata.updatePlayback)
		glib.IdleAdd(pbdata.setCover)
		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromPlayback() {
	log.Info("Unsubscribing from playback")

	pbdata.mu.Lock()
	id := pbdata.res.SubscriptionId
	pbdata.mu.Unlock()

	var cc *grpc.ClientConn
	var err error
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.UnsubscribeFromPlayback(
		context.Background(),
		&m3uetcpb.UnsubscribeFromPlaybackRequest{
			SubscriptionId: id,
		},
	)
	onerror.Log(err)
}

func (pbd *playbackData) setCover() bool {
	pbd.mu.Lock()
	defer func() { pbd.mu.Unlock() }()

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

func (pbd *playbackData) setPlaybackUI() (err error) {
	pbd.headerbar, err = builder.GetHeaderBar("window_headerbar")
	if err != nil {
		return
	}

	pbd.cover, err = builder.GetImage("cover")
	if err != nil {
		return
	}

	pbd.logoPixbuf, err = gdk.PixbufNewFromFile("data/images/m3u-etcetera.png")
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
			pbd.coverFiles = append(pbd.coverFiles, strings.Title(v)+ext)
		}
	}

	pbd.uiSet = true
	return
}

func (pbd *playbackData) updatePlayback() bool {
	log.Debug("Updating playback")

	pbd.mu.Lock()
	iconName := "media-playback-pause"
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
	pbd.mu.Unlock()

	if duration > 0 {
		pbd.prog.SetFraction(float64(position) / float64(duration))
		pbd.prog.SetText(
			fmt.Sprintf(
				"%v / %v",
				time.Duration(position)*time.Nanosecond,
				time.Duration(duration)*time.Nanosecond,
			),
		)
	} else {
		pbd.prog.SetFraction(float64(0))
		pbd.prog.SetText("Not Playing")
	}

	maxLen := 45
	subtitle := stringing.TruncateText(title, maxLen)
	if title == "" {
		title = "Not Playing"
	}
	if artist != "" {
		artist = "by " + artist
		if subtitle != "" {
			subtitle += " (" + stringing.TruncateText(artist, maxLen) + ")"
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

	pbd.title.SetText(stringing.TruncateText(title, maxLen))
	pbd.title.SetTooltipText(title)
	pbd.artist.SetText(stringing.TruncateText(artist, maxLen))
	pbd.artist.SetTooltipText(artist)
	pbd.source.SetText(stringing.TruncateText(location, maxLen))
	pbd.source.SetTooltipText(location)

	if oldTrackID != pbd.trackID {
		updatePlaybarModel()
	}
	return false
}
