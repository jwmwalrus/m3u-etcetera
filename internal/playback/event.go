package playback

import (
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
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
)

// GetPlayback returns a copy of the current playback
func GetPlayback() (pb *models.Playback) {
	if eng.pb == nil {
		return
	}

	log.Info("Obtaining current playback")

	pb = &models.Playback{}
	if err := pb.Read(eng.pb.ID); err != nil {
		log.Error(err)
	}
	return
}

// GetState returns the current state of the playback
func GetState() gst.StateOptions {
	return eng.state
}

// IsPaused checks if a stream is paused right now
func IsPaused() bool {
	return eng.state == gst.StatePaused
}

// IsPlaying checks if a stream is playing right now
func IsPlaying() bool {
	return eng.state == gst.StatePlaying
}

// IsStopped checks if a stream is playing right now
func IsStopped() bool {
	return eng.lastEvent == stopAllEvent
}

// NextStream plays the next stream in playlist
func NextStream() (err error) {
	eng.lastEvent = nextEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	if eng.playbin == nil || !(eng.state == gst.StatePlaying || eng.state == gst.StateReady) {
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

	eng.goingBack = true
	defer func() { eng.goingBack = false }()

	prevInHistory <- struct{}{}
	StopStream()
	return
}

// SeekInStream seek a position in the current stream
func SeekInStream(pos int64) {
	eng.lastEvent = seekEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	if eng.playbin == nil || eng.state != gst.StatePlaying {
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

	base.RegisterUnloader(base.Unloader{
		Description: "StopEngine",
		Callback: func() error {
			StopEngine()
			return nil
		},
	})
	return
}

// StopAll stops all playback
func StopAll() {
	eng.lastEvent = stopAllEvent
	log.Infof("Firing the %v event", eng.lastEvent)
	StopStream()
}

// StopEngine stops the engine
func StopEngine() {
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

// StopStream stops the current stream
func StopStream() {
	log.Info("Stopping current playback")

	if eng.lastEvent != stopAllEvent {
		eng.lastEvent = stopStreamEvent
		log.Infof("Firing the %v event", eng.lastEvent)
	}

	if eng.playbin == nil || !IsPlaying() {
		return
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
