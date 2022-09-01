package playback

import (
	"sync/atomic"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/onerror"

	// "github.com/notedit/gst"
	log "github.com/sirupsen/logrus"
	"github.com/tinyzimmer/go-gst/gst"
)

type engineEvent struct {
	atomic.Int32
}

// engine events
const (
	noLoopEvent = iota
	loopEvent
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
)

func (ee *engineEvent) Store(val int32) {
	ee.Int32.Store(val)
	log.Infof("Firing the %v event", ee.String())
}

// eeString returns the engine event string
func (ee *engineEvent) String() string {
	return [...]string{
		"NO-LOOP",
		"LOOP",
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
	}[ee.Load()]
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

	unloader = &base.Unloader{
		Description: "StopEngine",
		Callback: func() error {
			stopEngine()
			return nil
		},
	}
)

// GetPlayback returns a copy of the current playback
func GetPlayback() (pb *models.Playback, t *models.Track) {
	if eng.pb.Load() == nil {
		return
	}

	log.Debug("Obtaining current playback")

	pbcopy := *eng.pb.Load()
	pb = &pbcopy
	if pb.TrackID > 0 {
		t = &models.Track{}
		if eng.t.Load() != nil {
			t = eng.t.Load()
		} else {
			if err := t.Read(pb.TrackID); err != nil {
				log.Error(err)
				return
			}
			eng.t.Store(t)
		}
	} else {
		var err error
		if t, err = models.ReadTagsForLocation(pb.Location); err != nil {
			log.Error(err)
			return
		}
	}

	if t != nil && t.Duration == 0 {
		t.Duration = eng.duration.Load()
	}
	pb.Skip = eng.lastPosition.Load()
	return
}

// GetState returns the current state of the playback
func GetState() gst.State {
	return eng.state.Load()
}

// HasNextStream returns true if there is a playback/queue/playlist track after
// the current one
func HasNextStream() bool {
	pb := &models.Playback{}
	if pb.GetNextToPlay() != nil {
		return true
	}

	q, _ := models.GetActivePerspectiveIndex().GetPerspectiveQueue()
	if !q.IsEmpty() {
		return true
	}

	if eng.pt.Load() != nil {
		if _, err := eng.pt.Load().GetTrackAfter(false); err == nil {
			return true
		}
	}

	return false
}

// IsPaused checks if a stream is paused right now
func IsPaused() bool {
	return eng.playbin.Load() != nil && eng.state.Load() == gst.StatePaused
}

// IsPlaying checks if a stream is playing right now
func IsPlaying() bool {
	return eng.playbin.Load() != nil && eng.state.Load() == gst.StatePlaying
}

// IsReady checks if a stream is ready
func IsReady() bool {
	return eng.playbin.Load() != nil && eng.state.Load() == gst.StateReady
}

// IsStreaming check if a stream is playing or paused
func IsStreaming() bool {
	return IsPaused() || IsPlaying()
}

// IsStopped checks if a stream is playing right now
func IsStopped() bool {
	return eng.lastEvent.Load() == stopAllEvent
}

// NextStream plays the next stream in playlist
func NextStream() (err error) {
	eng.lastEvent.Store(nextEvent)

	if !(IsStreaming() || IsReady()) {
		return
	}

	eng.setPlaybackHint(hintNextInPlaylist)
	StopStream()
	models.TriggerPlaybackChange()
	return
}

// PauseStream pauses the current stream
func PauseStream(off bool) (err error) {
	if eng.playbin.Load() == nil {
		if !IsStopped() {
			return
		}
		eng.lastEvent.Store(resumeAllEvent)
		models.TriggerPlaybackChange()
		return
	}

	eng.lastEvent.Store(pauseEvent)

	if off {
		if !IsPaused() {
			return
		}
		eng.state.Store(gst.StatePlaying)
		err = eng.playbin.Load().SetState(eng.state.Load())
	} else {
		if !IsPlaying() {
			return
		}
		eng.state.Store(gst.StatePaused)
		err = eng.playbin.Load().SetState(eng.state.Load())
	}

	subscription.Broadcast(subscription.ToPlaybackEvent)
	return
}

// PlayStreams starts playback for the given streams
func PlayStreams(force bool, locations []string, ids []int64) {
	eng.lastEvent.Store(playEvent)

	entry := log.WithFields(log.Fields{
		"locations": locations,
		"ids":       ids,
	})
	entry.Infof("Playing streams")

	for _, v := range locations {
		models.AddPlaybackLocation(v)
	}
	for _, v := range ids {
		t := &models.Track{}
		err := t.Read(v)
		if err != nil {
			entry.Error(err)
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
	if time.Duration(eng.lastPosition.Load())*time.Nanosecond >=
		time.Duration(base.PlaybackPlayedThreshold)*time.Second {
		SeekInStream(0)
		return
	}

	eng.lastEvent.Store(previousEvent)

	eng.clearPendingPlayback()

	var hint playbackHint
	if IsStreaming() {
		if eng.pt.Load() != nil {
			hint = hintPrevInPlaylist
		} else {
			hint = hintPrevInHistory
		}
	} else {
		hint = hintPrevInHistory
	}
	eng.setPlaybackHint(hint)
	StopStream()
	models.TriggerPlaybackChange()
}

// QuitPlayingFromBar stops reproducing a playlist
func QuitPlayingFromBar(pl *models.Playlist) {
	eng.lastEvent.Store(stopPlaybarEvent)
	log.WithField("pl", *pl).
		Infof("Quit playing from bar")

	if !(pl.Open && pl.Active) {
		return
	}

	bar := pl.Playbar
	bar.DeactivateEntry(pl)
	quitPlayingFromList()
}

// SeekInStream seek a position in the current stream
func SeekInStream(pos int64) {
	eng.lastEvent.Store(seekEvent)

	if !IsPlaying() {
		return
	}

	if eng.seekable.Load() {
		seek := gst.NewSeekEvent(
			1.0,
			gst.FormatTime,
			gst.SeekFlagFlush|gst.SeekFlagKeyUnit,
			gst.SeekTypeSet,
			pos,
			gst.SeekTypeNone,
			-1,
		)

		if !eng.playbin.Load().SendEvent(seek) {
			log.Errorf("Error sending playback event: %v", gst.EventTypeSeek)
		}
	}
}

// SetMode sets engine mode
func SetMode(mode EngineMode) {
	eng.mode = mode
}

// StartEngine starts the playback engine
func StartEngine() *base.Unloader {
	log.Info("Starting playback engine")

	if eng.mode == TestMode {
		mockEngineLoop()
		return unloader
	}

	gst.Init(nil)

	eng.resumeActivePlaylist()
	go eng.engineLoop()
	models.TriggerPlaybackChange()

	return unloader
}

// StopAll stops all playback
func StopAll() {
	eng.lastEvent.Store(stopAllEvent)
	StopStream()

	eng.updateMPRIS(true)
}

// StopStream stops the current stream
func StopStream() {
	log.Info("Stopping current playback")

	if eng.lastEvent.Load() != stopAllEvent {
		eng.lastEvent.Store(stopStreamEvent)
	}

	if !IsStreaming() {
		return
	}

	// NOTE: sending the EOS event has between 1.5 and 3 seconds of
	// latency, and conflicts with the paused state, so we are
	// just ending things here
	if !IsPaused() {
		eng.state.Store(gst.StatePaused)
		onerror.Log(eng.playbin.Load().SetState(eng.state.Load()))
	}
	eng.wrapUp()
	eng.mainLoop.Quit()
}

// TryPlayingFromBar starts a playlist in the playbar
func TryPlayingFromBar(pl *models.Playlist, position int) {
	eng.lastEvent.Store(playbarEvent)

	entry := log.WithField("pl", pl)
	entry.Infof("Try playing from bar")

	bar := models.Playbar{}
	if err := bar.Read(pl.PlaybarID); err != nil {
		entry.Error(err)
		return
	}

	bar.ActivateEntry(pl)
	tryPlayingFromList(pl, position)
}

func mockEngineLoop() {
	aux := &models.Playback{}
	aux.GetNextToPlay()

	if aux.ID != 0 {
		eng.pb.Store(aux)
	}
}

func quitPlayingFromList() {
	eng.lastEvent.Store(stopPlaylistEvent)

	eng.setPlaybackHint(hintStopPlaylist)
	StopStream()
	models.TriggerPlaybackChange()

	eng.updateMPRIS(true)
}

func stopEngine() {
	if eng.mode == TestMode {
		return
	}

	if eng.lastEvent.Load() == noLoopEvent {
		return
	}
	log.Info("Stopping engine")

	eng.freezePlayback.Store(true)
	StopAll()

	for i := 0; i < base.ServerWaitTimeout; i++ {
		if eng.playbin.Load() != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	quitEngineLoop <- struct{}{}

	for i := 0; i < base.ServerWaitTimeout; i++ {
		if eng.lastEvent.Load() != noLoopEvent {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
}

func tryPlayingFromList(pl *models.Playlist, position int) {
	eng.lastEvent.Store(playlistEvent)

	pt, err := pl.GetTrackAt(position)
	if err != nil {
		log.Error(err)
		return
	}

	eng.pt.Store(pt)

	eng.setPlaybackHint(hintStartPlaylist)
	StopStream()
	models.TriggerPlaybackChange()
}
