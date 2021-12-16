package playback

import (
	"fmt"
	"time"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/notedit/gst"
	log "github.com/sirupsen/logrus"
)

var (
	quitEngineLoop chan struct{}
	prevInHistory chan struct{}
)

type engine struct {
	goingBack bool
	terminate    bool
	seekEnabled  bool
	seekDone     bool
	lastPosition int64
	duration     int64
	mode         EngineMode
	lastEvent    engineEvent
	prevState    gst.StateOptions
	state        gst.StateOptions
	playbin      *gst.Element
	pb *models.Playback
}

func (e *engine) addPlaybackFromQueue(qt *models.QueueTrack) (pb *models.Playback) {
	log.WithField("qt", qt).
		Info("Adding playback from queue")

	if qt.TrackID > 0 {
		t := &models.Track{}
		if err := t.Read(qt.TrackID); err != nil {
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

func (e *engine) debugChannel(msg string) {
	if e.mode == TestMode {
		models.DbgChan <- msg
		return
	}
	log.Debug(msg)
}

func (e *engine) engineLoop() {
	log.Info("Starting engine loop")

	idleCount := 0
	shouldBreak := func() bool {
		log.WithField("idleCount", idleCount).
			Info("Checking if I should break from engineLoop")

		if idleCount >= base.ServerIdleTimeout {
			if !base.IsItBusy() {
				e.debugChannel("It's been about 5 minutes already, so I'm off for now")
				return true
			}
		}
		time.Sleep(1 * time.Second)
		idleCount++
		return false
	}

loop:
	for {
		pb := &models.Playback{}

		select {
		case <-quitEngineLoop:
			break loop
		case <-prevInHistory:
			pb = e.getPrevInHistory()
		default:
			switch e.lastEvent {
			case stopAllEvent:
				if shouldBreak() {
					idleCount = 0
					go base.Idle(false)
				}
				continue
			default:
			}

			err := pb.GetNextToPlay()
			if err != nil {
				q, _ := models.GetActivePerspectiveIndex().GetPerspectiveQueue()
				if qt, err := q.Pop(); err == nil {
					e.debugChannel("Popped successfully")
					pb = e.addPlaybackFromQueue(qt)
				}
			}
		}

		if pb.ID > 0 && pb.Location != "" {
			e.debugChannel("There is a playback")
			e.playStream(pb)
			idleCount = 0
			continue
		}

		if shouldBreak() {
			idleCount = 0
			go base.Idle(false)
		}
	}

	e.lastEvent = noLoopEvent
	log.Infof("Firing the %v event", e.lastEvent)
}

func (e *engine) getPrevInHistory() (pb *models.Playback) {
	log.Info("Obtaining previous track in history")

	h := models.PlaybackHistory{}

	if e.playbin != nil {
		if err := h.FindLastBy(fmt.Sprintf("location <> %v", e.pb.Location)); err != nil {
			return
		}
	} else {
		if err := h.ReadLast(); err != nil {
			return
		}
	}

	if h.TrackID > 0 {
		t := &models.Track{}
		if err := t.Read(h.TrackID); err != nil {
			pb = models.AddPlaybackTrack(t)
			return
		}
	} else if h.Location != "" {
		pb = models.AddPlaybackLocation(h.Location)
	}
	return
}

func (e *engine) playStream(pb *models.Playback) {
	log.WithField("pb", *pb).
		Info("Starting playStream")

	e.pb = pb
	e.terminate = false

	base.GetBusy()
	defer func() { base.GetFree() }()

	// check if playback is valid
	if e.pb == nil || pb.Location == "" {
		return
	}
	e.debugChannel("Playback is valid")

	var err error

	e.playbin, err = gst.ElementFactoryMake("playbin", "playbin")
	if err != nil {
		log.Error(err)
		return
	}
	e.debugChannel("Playbin created")

	if e.playbin == nil {
		log.Error("Not all elements could be created")
		return
	}

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
		e.playbin = nil
		e.pb = nil
		return
	}
	log.Infof("StateChangeReturn: %d\n", gst.StatePlaying)

	bus := e.playbin.GetBus()

	for !e.terminate {
		msg := bus.TimedPopFiltered(100*time.Millisecond, gst.MessageStateChanged|gst.MessageError|gst.MessageWarning|gst.MessageInfo|gst.MessageEos|gst.MessageDurationChanged)
		// msg := bus.Pull(gst.MessageStateChanged | gst.MessageError | gst.MessageWarning | gst.MessageInfo | gst.MessageEos | gst.MessageDurationChanged)

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
		e.lastPosition = int64(position)

		/* If we didn't know it yet, query the stream duration */
		if e.duration == 0 {
			var duration time.Duration
			if duration, err = e.playbin.QueryDuration(); err != nil {
				log.Warn("Could not query current duration")
			}
			e.duration = int64(duration)
		}
	}

	e.debugChannel("End of playback")
	e.state = gst.StateNull
	e.playbin.SetState(e.state)
	e.playbin = nil
	e.pb = nil
	return
}

func (e *engine) handleMessage(msg *gst.Message, pb *models.Playback) {
	e.debugChannel(msg.GetName())

	switch msg.GetType() {
	case gst.MessageEos:
		log.Info("End of stream: " + pb.Location)
		if !e.goingBack {
			go pb.AddToHistory(e.lastPosition)
		}
		e.terminate = true
		break

	case gst.MessageDurationChanged:
		e.duration = 0

	case gst.MessageInfo:
		log.Info(msg.GetName())
	case gst.MessageError:
		log.Error(msg.GetName())
		e.terminate = true
	case gst.MessageWarning:
		log.Warning(msg.GetName())
	case gst.MessageStateChanged:
		e.prevState, e.state, _ = msg.ParseStateChanged()
		// if (GST_MESSAGE_SRC (msg) == GST_OBJECT (data->playbin)) {
		log.WithFields(log.Fields{
			"previousState": e.prevState,
			"newState":      e.state,
		}).Info("Pipeline state changed")

		if e.state == gst.StatePlaying {
			/* We just moved to PLAYING. Check if seeking is possible */
			var start, end time.Duration
			q, _ := gst.QueryNewSeeking(gst.FormatTime)
			if e.playbin.Query(q) {
				e.seekEnabled, start, end = q.ParseSeeking(nil)
				if e.seekEnabled {
					log.WithFields(log.Fields{
						"start": start,
						"end":   end,
					}).Info("Seeking is ENABLED")
				} else {
					log.Info("Seeking is DISABLED for this stream")
				}
			} else {
				log.Info("Seeking query failed")
			}
		}

	default:
	}

	msg = nil
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
	}
	quitEngineLoop = make(chan struct{})
	prevInHistory = make(chan struct{})
}
