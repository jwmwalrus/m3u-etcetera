package playback

import (
	"context"
	"fmt"
	"time"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/onerror"

	// "github.com/notedit/gst"
	log "github.com/sirupsen/logrus"
	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

const (
	positionThreshold = 1e9
)

var (
	quitEngineLoop chan struct{}
)

type playbackHint int

const (
	hintNone playbackHint = iota
	hintPrevInHistory
	hintStopPlaylist
	hintStartPlaylist
	hintNextInPlaylist
	hintPrevInPlaylist
)

type engine struct {
	freezePlayback bool
	terminate      bool
	seekEnabled    bool
	seekDone       bool
	lastPosition   int64
	duration       int64
	buffering      int
	mode           EngineMode
	lastEvent      engineEvent
	hint           playbackHint
	pt             *models.PlaylistTrack
	pb             *models.Playback
	t              *models.Track

	mainLoop  *glib.MainLoop
	prevState gst.State
	state     gst.State
	playbin   *gst.Element
}

func (e *engine) addPlaybackFromQueue(qt *models.QueueTrack) (pb *models.Playback) {
	log.WithField("qt", qt).
		Info("Adding playback from queue")

	if qt.TrackID > 0 {
		t := &models.Track{}
		if err := t.Read(qt.TrackID); err == nil {
			pb = models.AddPlaybackTrack(t)
			return
		}
	}

	pb = models.AddPlaybackLocation(qt.Location)
	return
}

func (e *engine) clearPendingPlayback() {
	if e.pb == nil {
		return
	}

	e.pb.ClearPending()
}

func (e *engine) engineLoop() {
	log.Info("Starting engine loop")

loop:
	for {
		pb := &models.Playback{}

	signals:
		select {
		case <-quitEngineLoop:
			break loop
		case <-models.PlaybackChanged:
			switch e.lastEvent {
			case stopAllEvent:
				pb = &models.Playback{}
				break signals
			default:
			}

			var err error
			switch e.getPlaybackHint() {
			case hintPrevInHistory:
				pb, err = e.getPrevInHistory()
			case hintStartPlaylist:
				e.clearPendingPlayback()
				pb = e.getFirstInPlaylist()
				break signals
			case hintStopPlaylist:
				e.pt = nil
				continue loop
			case hintPrevInPlaylist:
				pb = e.getNextInPlaylist(true)
				break signals
			case hintNextInPlaylist:
				pb = e.getNextInPlaylist(false)
				break signals
			default:
				err = pb.GetNextToPlay()
			}

			if err != nil {
				q, _ := models.GetActivePerspectiveIndex().GetPerspectiveQueue()
				if qt := q.Pop(); qt != nil {
					log.Debug("Popped successfully")
					pb = e.addPlaybackFromQueue(qt)
				} else if e.pt != nil {
					pb = e.getNextInPlaylist(false)
				}
			}
			break signals
		}

		if pb != nil && pb.ID > 0 && pb.Location != "" {
			log.Debug("There is a playback")
			e.playStream(pb)
			go func() {
				models.PlaybackChanged <- struct{}{}
			}()
			continue loop
		}

		if !e.freezePlayback {
			models.DeactivatePlaybars()
		}
		subscription.Broadcast(subscription.ToPlaybackEvent)
	}

	e.lastEvent = noLoopEvent
	log.Infof("Firing the %v event", e.lastEvent)
}

func (e *engine) getFirstInPlaylist() (pb *models.Playback) {
	if e.pt == nil {
		log.Error("There is no playlist-track available")
		return
	}
	pl := e.pt.Playlist
	t := e.pt.Track
	if pl.ID == 0 || t.ID == 0 {
		log.Error("There is no list or track to play from")
		return
	}

	log.WithField("pt", *e.pt).
		Info("Obtaining first track in playlist")

	pb = models.AddPlaybackTrack(&t)
	return
}

func (e *engine) getNextInPlaylist(goingBack bool) (pb *models.Playback) {
	if e.pt == nil {
		log.Error("There is no playlist-track available")
		return
	}
	pt := *e.pt
	pl := e.pt.Playlist
	if pl.ID == 0 {
		log.Error("There is no list to play from")
		return
	}

	log.WithField("pt", *e.pt).
		Info("Obtaining next track in playlist")

	newpt, err := pl.GetTrackAfter(pt, goingBack)
	if err != nil {
		bar := models.Playbar{}
		if errb := bar.Read(pl.PlaybarID); errb == nil {
			bar.DeactivateEntry(&pl)
		}
		e.pt = nil
		log.Info(err)
		return
	}

	e.pt = newpt
	newt := newpt.Track
	pb = models.AddPlaybackTrack(&newt)

	return
}

func (e *engine) getPlaybackHint(keep ...bool) (h playbackHint) {
	h = e.hint

	reset := true
	if len(keep) > 0 {
		reset = !keep[0]
	}
	if reset {
		e.hint = hintNone
	}
	return
}

func (e *engine) getPrevInHistory() (pb *models.Playback, err error) {
	log.Info("Obtaining previous track in history")

	h := models.PlaybackHistory{}

	if err = h.ReadLast(); err != nil {
		log.Error(err)
		return
	}

	if h.TrackID > 0 {
		t := &models.Track{}
		if err = t.Read(h.TrackID); err != nil {
			log.Error(err)
			return
		}
		pb = models.AddPlaybackTrack(t)
	} else if h.Location != "" {
		pb = models.AddPlaybackLocation(h.Location)
	} else {
		err = fmt.Errorf("History entry lacks both track and location")
	}
	return
}

func (e *engine) handleBusMessage(msg *gst.Message, pb *models.Playback) bool {
	log.Debug(msg.String())

	switch msg.Type() {
	case gst.MessageEOS:
		log.Debugf("End of stream: %v", pb.Location)
		if e.getPlaybackHint(true) != hintPrevInHistory {
			go pb.AddToHistory(e.lastPosition, e.duration, e.freezePlayback)
		}
		e.terminate = true
		e.mainLoop.Quit()
		break

	case gst.MessageError:
		log.Error(msg.String())
		e.terminate = true
	case gst.MessageWarning:
		log.Warning(msg.String())
	case gst.MessageInfo:
		log.Debug(msg.String())
	case gst.MessageStateChanged:
		e.prevState, e.state = msg.ParseStateChanged()
		log.WithFields(log.Fields{
			"previousState": e.prevState,
			"newState":      e.state,
		}).
			Debug("Pipeline state changed")

		if e.state == gst.StatePlaying {
			// We just moved to PLAYING. Check if seeking is possible
			var start, end int64
			q := gst.NewSeekingQuery(gst.FormatTime)
			if e.playbin.Query(q) {
				var format gst.Format
				format, e.seekEnabled, start, end = q.ParseSeeking()
				if e.seekEnabled {
					log.WithFields(log.Fields{
						"format": format,
						"start":  start,
						"end":    end,
					}).
						Debug("Seeking is ENABLED")
				} else {
					log.Debug("Seeking is DISABLED for this stream")
				}
			} else {
				log.Debug("Seeking query failed")
			}
		}

	case gst.MessageDurationChanged:
		e.duration = 0

	case gst.MessageBuffering:
		e.buffering = msg.ParseBuffering()
		if e.buffering < 100 {
			e.state = gst.StatePaused
		} else {
			e.state = gst.StatePlaying
		}
		onerror.Log(e.playbin.SetState(e.state))

	default:
	}

	return true
}

func (e *engine) playStream(pb *models.Playback) {
	log.WithField("pb", *pb).
		Info("Starting playStream")

	e.pb = pb
	e.terminate = false
	defer func() { e.pb = nil; e.t = nil }()

	base.GetBusy(base.IdleStatusEngineLoop)
	defer func() { base.GetFree(base.IdleStatusEngineLoop) }()

	// check if playback is valid
	if e.pb == nil || pb.Location == "" {
		return
	}
	log.Debug("Playback is valid")

	var err error

	e.mainLoop = glib.NewMainLoop(glib.MainContextDefault(), false)

	e.playbin, err = gst.NewElementWithName("playbin", "m3uetc-playbin")
	if err != nil {
		log.Error(err)
		return
	}
	defer func() { e.playbin = nil }()
	log.Debug("Playbin created")

	if e.playbin == nil {
		log.Error("Not all elements could be created")
		return
	}

	subscription.Broadcast(subscription.ToPlaybackEvent)

	if e.mode == TestMode {
		// Location is relative
		loc, _ := urlstr.PathToURL(pb.Location)
		e.playbin.Set("uri", loc)
	} else {
		e.playbin.Set("uri", pb.Location)
	}

	bus := e.playbin.GetBus()

	e.state = gst.StatePlaying
	if err := e.playbin.SetState(e.state); err != nil {
		log.Errorf("Unable to start playback: %v", err)
		return
	}
	log.Debugf("State changed: %d\n", gst.StatePlaying)

	qposctx, cancelqpos := context.WithCancel(context.Background())
	go e.queryPosition(qposctx)

	bus.AddWatch(func(msg *gst.Message) bool {
		return e.handleBusMessage(msg, pb)
	})

	e.mainLoop.Run()

	bus.RemoveWatch()

	cancelqpos()

	log.Debug("End of playback")
	e.state = gst.StateNull
	e.playbin.SetState(e.state)
	return
}

func (e *engine) queryPosition(ctx context.Context) {
	tick := time.NewTicker(time.Duration(positionThreshold) * time.Nanosecond).C
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			if e.terminate {
				break
			}

			if e.state != gst.StatePlaying {
				continue
			}

			// Query the current position of the stream
			var position int64
			var ok bool
			if ok, position = e.playbin.QueryPosition(gst.FormatTime); !ok {
				log.Warn("Could not query current position")
			}

			if position > 0 {
				subscription.Broadcast(subscription.ToPlaybackEvent)
				e.lastPosition = position
			}

			// If we didn't know it yet, query the stream duration
			if e.duration == 0 {
				var duration int64
				if ok, duration = e.playbin.QueryDuration(gst.FormatTime); !ok {
					log.Warn("Could not query current duration")
				}
				e.duration = duration
				subscription.Broadcast(subscription.ToPlaybackEvent)
			}

			// NOTE: If seeking is enabled, we have not done it yet, and the time is right, seek
			//
			// if e.seekEnabled && !e.seekDone && e.lastPosition > 10 * time.Second {
			//   g_print ("\nReached 10s, performing seek...\n");
			//   gst_element_seek_simple (data.playbin, GST_FORMAT_TIME,
			//       GST_SEEK_FLAG_FLUSH | GST_SEEK_FLAG_KEY_UNIT, 30 * GST_SECOND);
			//   data.seek_done = TRUE;
			// }
			//
		}
	}
}

func (e *engine) resumeActivePlaylist() {
	e.pt = models.GetActivePlaylistTrack()
	if e.pt != nil {
		log.Info("Resuming playback for active playlist")
	}
}

func (e *engine) setPlaybackHint(h playbackHint) {
	e.hint = h
}

func init() {
	eng = &engine{
		state:     gst.StateNull,
		lastEvent: noLoopEvent,
		mode:      NormalMode,
		hint:      hintNone,
	}
	quitEngineLoop = make(chan struct{})
}
