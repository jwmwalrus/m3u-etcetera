package playback

import (
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/notedit/gst"
	log "github.com/sirupsen/logrus"
)

type engineEvent int

const (
	noLoopEvent engineEvent = iota
	playEvent
	pauseEvent
	previousEvent
	nextEvent
	seekEvent
	stopStreamEvent
	playbarEvent
	stopPlaybarEvent
	playlistEvent
	stopPlaylistEvent
	stopAllEvent
	resumeAllEvent
	watchQueue
)

func (ee engineEvent) String() string {
	return [...]string{
		"NO-LOOP",
		"PLAY",
		"PAUSE",
		"PREVIOUS",
		"NEXT",
		"SEEK",
		"STOP-STREAM",
		"PLAYBAR",
		"STOP-PLAYBAR",
		"PLAYLIST",
		"STOP-PLAYLIST",
		"STOP-ALL",
		"RESUME-ALL",
		"WATCH-QUEUE",
	}[ee]
}

// EngineMode defines the engine mode type
type EngineMode int

const (
	// NormalMode regular operation mode
	NormalMode EngineMode = iota

	// TestMode testing operation mode
	TestMode
)

func (em EngineMode) String() string {
	return [...]string{"normal", "testing"}[em]
}

var (
	eng *engine

	// Unloader -
	Unloader = base.Unloader{
		Description: "StopEngine",
		Callback: func() error {
			stopEngine()
			return nil
		},
	}
)

// GetPlayback returns a copy of the current playback
func GetPlayback() (pb *models.Playback, t *models.Track) {
	if eng.pb == nil {
		return
	}

	log.Info("Obtaining current playback")

	pbcopy := *eng.pb
	pb = &pbcopy
	if pb.TrackID > 0 {
		t = &models.Track{}
		if eng.t != nil {
			t = eng.t
		} else {
			if err := t.Read(pb.TrackID); err != nil {
				log.Error(err)
				return
			}
			eng.t = t
		}
	} else {
		var err error
		if t, err = models.ReadTagsForLocation(pb.Location); err != nil {
			log.Error(err)
			return
		}
	}

	if t != nil && t.Duration == 0 {
		t.Duration = eng.duration
	}
	pb.Skip = eng.lastPosition
	return
}

// GetState returns the current state of the playback
func GetState() gst.StateOptions {
	return eng.state
}

// IsPaused checks if a stream is paused right now
func IsPaused() bool {
	return eng.playbin != nil && eng.state == gst.StatePaused
}

// IsPlaying checks if a stream is playing right now
func IsPlaying() bool {
	return eng.playbin != nil && eng.state == gst.StatePlaying
}

// IsReady checks if a stream is ready
func IsReady() bool {
	return eng.playbin != nil && eng.state == gst.StateReady
}

// IsStreaming check if a stream is playing or paused
func IsStreaming() bool {
	return IsPaused() || IsPlaying()
}

// IsStopped checks if a stream is playing right now
func IsStopped() bool {
	return eng.lastEvent == stopAllEvent
}

// NextStream plays the next stream in playlist
func NextStream() (err error) {
	eng.lastEvent = nextEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	if !(IsStreaming() || IsReady()) {
		return
	}

	StopStream()
	return
}

// PauseStream pauses the current stream
func PauseStream(off bool) (err error) {
	if eng.playbin == nil {
		if !IsStopped() {
			return
		}
		eng.lastEvent = resumeAllEvent
		log.Infof("Firing the %v event", eng.lastEvent)
		models.PlaybackChanged <- struct{}{}
		return
	}

	eng.lastEvent = pauseEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	if off {
		if !IsPaused() {
			return
		}
		eng.state = gst.StatePlaying
		eng.playbin.SetState(eng.state)
	} else {
		if !IsPlaying() {
			return
		}
		eng.state = gst.StatePaused
		eng.playbin.SetState(eng.state)
	}

	subscription.Broadcast(subscription.ToPlaybackEvent)
	return
}

// PlayStreams starts playback for the given streams
func PlayStreams(force bool, locations []string, ids []int64) {
	eng.lastEvent = playEvent
	log.WithFields(log.Fields{
		"locations": locations,
		"ids":       ids,
	}).
		Infof("Firing the %v event", eng.lastEvent)

	for _, v := range locations {
		models.AddPlaybackLocation(v)
	}
	for _, v := range ids {
		t := &models.Track{}
		err := t.Read(v)
		if err != nil {
			log.Error(err)
			continue
		}
		models.AddPlaybackTrack(t)
	}

	if force {
		StopStream()
	}
}

// PreviousStream plays the previous stream in history
func PreviousStream() {
	eng.lastEvent = previousEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	eng.clearPendingPlayback()

	if IsStreaming() {
		/*
			if eng.playingFromList {
				hint = hintPrevInPlaylist
			}
		*/
		eng.setPlaybackHint(hintPrevInHistory)
	} else {
		eng.setPlaybackHint(hintPrevInHistory)
	}
	StopStream()
	models.PlaybackChanged <- struct{}{}
}

// SeekInStream seek a position in the current stream
func SeekInStream(pos int64) {
	eng.lastEvent = seekEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	if !IsPlaying() {
		return
	}

	if eng.seekEnabled && time.Duration(eng.lastPosition) > 10*time.Second {
		eng.playbin.SeekSimple(gst.FormatTime,
			gst.SeekFlagFlush|gst.SeekFlagKeyUnit, time.Duration(pos))
	}

	return
}

// SetMode sets engine mode
func SetMode(mode EngineMode) {
	eng.mode = mode
}

// StartEngine starts the playback engine
func StartEngine() {
	log.Info("Starting playback engine")

	if eng.mode == TestMode {
		mockEngineLoop()
		return
	}

	go eng.engineLoop()
	models.PlaybackChanged <- struct{}{}
	return
}

// StopAll stops all playback
func StopAll() {
	eng.lastEvent = stopAllEvent
	log.Infof("Firing the %v event", eng.lastEvent)
	StopStream()
}

// StopStream stops the current stream
func StopStream() {
	log.Info("Stopping current playback")

	if eng.lastEvent != stopAllEvent {
		eng.lastEvent = stopStreamEvent
		log.Infof("Firing the %v event", eng.lastEvent)
	}

	if !IsStreaming() {
		return
	}

	// FIXME: this is a hack to avoid hanging/freezing
	if IsPaused() {
		eng.state = gst.StatePlaying
		eng.playbin.SetState(eng.state)
	}

	eos := gst.NewEosEvent()
	eng.playbin.SendEvent(eos)

	return
}

func mockEngineLoop() {
	aux := &models.Playback{}
	aux.GetNextToPlay()
	if aux != nil {
		eng.pb = aux
	}
}

// stopEngine stops the engine
func stopEngine() {
	if eng.mode == TestMode {
		return
	}

	if eng.lastEvent == noLoopEvent {
		return
	}
	log.Info("Stopping engine")

	StopAll()

	for i := 0; i < base.ServerWaitTimeout; i++ {
		if eng.playbin != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	quitEngineLoop <- struct{}{}

	for i := 0; i < base.ServerWaitTimeout; i++ {
		if eng.lastEvent != noLoopEvent {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return
}
