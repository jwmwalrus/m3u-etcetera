package playback

import (
	"fmt"
	"time"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/notedit/gst"
	log "github.com/sirupsen/logrus"
)

const (
	positionThreshold = 5e8
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
	mode           EngineMode
	lastEvent      engineEvent
	hint           playbackHint
	prevState      gst.StateOptions
	state          gst.StateOptions
	playbin        *gst.Element
	pt             *models.PlaylistTrack
	pb             *models.Playback
	t              *models.Track
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
	if e.pt == nil || e.pt.Playlist.ID == 0 || e.pt.Track.ID == 0 {
		log.Error("There is no list to play from")
		return
	}

	log.WithField("pt", *e.pt).
		Info("Obtaining first track in playlist")

	pb = models.AddPlaybackTrack(&e.pt.Track)
	return
}

func (e *engine) getNextInPlaylist(goingBack bool) (pb *models.Playback) {
	if e.pt == nil || e.pt.Playlist.ID == 0 {

		log.Error("There is no list to play from")
		return
	}

	log.WithField("pt", *eng.pt).
		Info("Obtaining next track in playlist")

	pt, err := e.pt.Playlist.GetTrackAfter(*e.pt, goingBack)
	if err != nil {
		log.Error(err)
		return
	}

	e.pt = pt
	pb = models.AddPlaybackTrack(&pt.Track)

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

	e.playbin, err = gst.ElementFactoryMake("playbin", "m3uetc-playbin")
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
		e.playbin.SetObject("uri", loc)
	} else {
		e.playbin.SetObject("uri", pb.Location)
	}

	e.state = gst.StatePlaying
	if state := e.playbin.SetState(e.state); state == gst.StateChangeFailure {
		log.Warn("Unable to start playback")
		return
	}
	log.Debugf("StateChangeReturn: %d\n", gst.StatePlaying)

	bus := e.playbin.GetBus()

	for !e.terminate {
		msg := bus.TimedPopFiltered(100*time.Millisecond, gst.MessageEos|
			gst.MessageError|gst.MessageWarning|gst.MessageInfo|
			gst.MessageStateChanged|gst.MessageDurationChanged)

		if msg != nil {
			e.handleMessage(msg, pb)
			continue
		}

		if e.state != gst.StatePlaying {
			continue
		}

		/* Query the current position of the stream */
		var position time.Duration
		if position, err = e.playbin.QueryPosition(); err != nil {
			log.Warn("Could not query current position")
		}

		if int64(position)-e.lastPosition > positionThreshold {
			go func() {
				time.Sleep(1 * time.Second)
				subscription.Broadcast(subscription.ToPlaybackEvent)
			}()
		}
		e.lastPosition = int64(position)

		/* If we didn't know it yet, query the stream duration */
		if e.duration == 0 {
			var duration time.Duration
			if duration, err = e.playbin.QueryDuration(); err != nil {
				log.Warn("Could not query current duration")
			}
			e.duration = int64(duration)
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

	log.Debug("End of playback")
	e.state = gst.StateNull
	e.playbin.SetState(e.state)
	return
}

func (e *engine) handleMessage(msg *gst.Message, pb *models.Playback) {
	log.Debug(msg.GetName())

	switch msg.GetType() {
	case gst.MessageEos:
		log.Debugf("End of stream: %v", pb.Location)
		if e.getPlaybackHint(true) != hintPrevInHistory {
			go pb.AddToHistory(e.lastPosition, e.duration, e.freezePlayback)
		}
		e.terminate = true
		break

	case gst.MessageError:
		log.Error(msg.GetName())
		e.terminate = true
	case gst.MessageWarning:
		log.Warning(msg.GetName())
	case gst.MessageInfo:
		log.Debug(msg.GetName())
	case gst.MessageStateChanged:
		e.prevState, e.state, _ = msg.ParseStateChanged()
		// if (GST_MESSAGE_SRC (msg) == GST_OBJECT (data->playbin)) {
		log.WithFields(log.Fields{
			"previousState": e.prevState,
			"newState":      e.state,
		}).
			Debug("Pipeline state changed")

		if e.state == gst.StatePlaying {
			// We just moved to PLAYING. Check if seeking is possible
			var start, end time.Duration
			q, _ := gst.QueryNewSeeking(gst.FormatTime)
			if e.playbin.Query(q) {
				e.seekEnabled, start, end = q.ParseSeeking(nil)
				if e.seekEnabled {
					log.WithFields(log.Fields{
						"start": start,
						"end":   end,
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

	default:
	}

	msg = nil
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

func (e *engine) softReset() {
	log.Info("Doing a soft reset on playback data")
	e.state = gst.StateNull
	e.playbin = nil
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
