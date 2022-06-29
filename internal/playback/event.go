package playback

import (
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"github.com/jwmwalrus/onerror"

	// "github.com/notedit/gst"
	log "github.com/sirupsen/logrus"
	"github.com/tinyzimmer/go-gst/gst"
)

type engineEvent int

const (
	noLoopEvent engineEvent = iota
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

func (ee engineEvent) String() string {
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

	log.Debug("Obtaining current playback")

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
func GetState() gst.State {
	return eng.state
}

func HasNextStream() bool {
	pb := &models.Playback{}
	if pb.GetNextToPlay() != nil {
		return true
	}

	q, _ := models.GetActivePerspectiveIndex().GetPerspectiveQueue()
	if !q.IsEmpty() {
		return true
	}

	if eng.pt != nil {
		if _, err := eng.pt.GetTrackAfter(false); err == nil {
			return true
		}
	}

	return false
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

	eng.setPlaybackHint(hintNextInPlaylist)
	StopStream()
	models.PlaybackChanged <- struct{}{}
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
		err = eng.playbin.SetState(eng.state)
	} else {
		if !IsPlaying() {
			return
		}
		eng.state = gst.StatePaused
		err = eng.playbin.SetState(eng.state)
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
	if time.Duration(eng.lastPosition)*time.Nanosecond >=
		time.Duration(base.PlaybackPlayedThreshold)*time.Second {
		SeekInStream(0)
		return
	}

	eng.lastEvent = previousEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	eng.clearPendingPlayback()

	var hint playbackHint
	if IsStreaming() {
		if eng.pt != nil {
			hint = hintPrevInPlaylist
		} else {
			hint = hintPrevInHistory
		}
	} else {
		hint = hintPrevInHistory
	}
	eng.setPlaybackHint(hint)
	StopStream()
	models.PlaybackChanged <- struct{}{}
}

// QuitPlayingFromBar stops reproducing a playlist
func QuitPlayingFromBar(pl *models.Playlist) {
	eng.lastEvent = stopPlaybarEvent
	log.WithField("pl", *pl).
		Infof("Firing the %v event", eng.lastEvent)

	if !(pl.Open && pl.Active) {
		return
	}

	bar := pl.Playbar
	bar.DeactivateEntry(pl)
	quitPlayingFromList()
}

// SeekInStream seek a position in the current stream
func SeekInStream(pos int64) {
	eng.lastEvent = seekEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	if !IsPlaying() {
		return
	}

	if eng.seekable {
		seek := gst.NewSeekEvent(
			1.0,
			gst.FormatTime,
			gst.SeekFlagFlush|gst.SeekFlagKeyUnit,
			gst.SeekTypeSet,
			pos,
			gst.SeekTypeNone,
			-1,
		)
		if !eng.playbin.SendEvent(seek) {
			log.Errorf("Error sending playback event: %v", gst.EventTypeSeek)
		}
	}
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

	gst.Init(nil)

	eng.resumeActivePlaylist()
	go eng.engineLoop()
	models.PlaybackChanged <- struct{}{}
}

// StopAll stops all playback
func StopAll() {
	eng.lastEvent = stopAllEvent
	log.Infof("Firing the %v event", eng.lastEvent)
	StopStream()

	eng.updateMPRIS(true)
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

	// NOTE: sending the EOS event has between 1.5 and 3 seconds of
	// latency, and conflicts with the paused state, so we are
	// just ending things here
	if !IsPaused() {
		eng.state = gst.StatePaused
		onerror.Log(eng.playbin.SetState(eng.state))
	}
	eng.wrapUp()
	eng.mainLoop.Quit()
}

// TryPlayingFromBar starts a playlist in the playbar
func TryPlayingFromBar(pl *models.Playlist, position int) {
	eng.lastEvent = playbarEvent
	log.WithField("pl", pl).
		Infof("Firing the %v event", eng.lastEvent)

	bar := models.Playbar{}
	if err := bar.Read(pl.PlaybarID); err != nil {
		log.Error(err)
		return
	}

	bar.ActivateEntry(pl)
	tryPlayingFromList(pl, position)
}

func mockEngineLoop() {
	aux := &models.Playback{}
	aux.GetNextToPlay()
	if aux.ID != 0 {
		eng.pb = aux
	}
}

func quitPlayingFromList() {
	eng.lastEvent = stopPlaylistEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	eng.setPlaybackHint(hintStopPlaylist)
	StopStream()
	models.PlaybackChanged <- struct{}{}

	eng.updateMPRIS(true)
}

func stopEngine() {
	if eng.mode == TestMode {
		return
	}

	if eng.lastEvent == noLoopEvent {
		return
	}
	log.Info("Stopping engine")

	eng.freezePlayback = true
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
}

func tryPlayingFromList(pl *models.Playlist, position int) {
	eng.lastEvent = playlistEvent
	log.Infof("Firing the %v event", eng.lastEvent)

	pt, err := pl.GetTrackAt(position)
	if err != nil {
		log.Error(err)
		return
	}
	eng.pt = pt
	eng.setPlaybackHint(hintStartPlaylist)
	StopStream()
	models.PlaybackChanged <- struct{}{}
}
