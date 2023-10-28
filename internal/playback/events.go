package playback

import (
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/go-gst/go-gst/gst"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	rtc "github.com/jwmwalrus/rtcycler"
)

type engineEvent struct {
	atomic.Int32
}

// engine events.
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
	slog.Info("Firing event", "event", ee)
}

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

var (
	unloader = &rtc.Unloader{
		Description: "StopEngine",
		Callback: func() error {
			stopEngine()
			return nil
		},
	}
)

// IEvents defines the events interface.
type IEvents interface {
	// getters:

	// GetPlayback returns a copy of the current playback.
	GetPlayback() (pb *models.Playback, t *models.Track)

	// GetState returns the current state of the playback.
	GetState() gst.State

	// status:

	// HasNextStream returns true if there is a playback/queue/playlist track after
	// the current one.
	HasNextStream() bool

	// IsPaused checks if a stream is paused right now.
	IsPaused() bool

	// IsPlaying checks if a stream is playing right now.
	IsPlaying() bool

	// IsReady checks if a stream is ready.
	IsReady() bool

	// IsStreaming check if a stream is playing or paused.
	IsStreaming() bool

	// IsStopped checks if a stream is playing right now.
	IsStopped() bool

	// actions:

	// NextStream plays the next stream in playlist.
	NextStream() (err error)

	// PauseStream pauses the current stream.
	PauseStream(off bool) (err error)

	// PlayStreams starts playback for the given streams.
	PlayStreams(force bool, locations []string, ids []int64)

	// PreviousStream plays the previous stream in history.
	PreviousStream()

	// QuitPlayingFromBar stops reproducing a playlist.
	QuitPlayingFromBar(pl *models.Playlist)

	// SeekInStream seek a position in the current stream.
	SeekInStream(pos int64)

	// StopAll stops all playback.
	StopAll()

	// StopStream stops the current stream.
	StopStream()

	// TryPlayingFromBar starts a playlist in the playbar.
	TryPlayingFromBar(pl *models.Playlist, position int)
}

// events implements the IEvents interface.
type events struct {
	eng *engine
}

var (
	instance *events
)

func (et *events) GetPlayback() (pb *models.Playback, t *models.Track) {
	if et.eng.pb.Load() == nil {
		return
	}

	slog.Debug("Obtaining current playback")

	pbcopy := *et.eng.pb.Load()
	pb = &pbcopy
	if pb.TrackID > 0 {
		t = &models.Track{}
		if et.eng.t.Load() != nil {
			t = et.eng.t.Load()
		} else {
			if err := t.Read(pb.TrackID); err != nil {
				slog.Error("Failed to read track", "error", err)
				return
			}
			et.eng.t.Store(t)
		}
	} else if et.eng.t.Load() != nil {
		t = et.eng.t.Load()
	} else {
		var err error
		if t, err = models.ReadTagsForLocation(pb.Location); err != nil {
			slog.With(
				"location", pb.Location,
				"error", err,
			).Error("Failed to read tags for location")
			return
		}
		et.eng.t.Store(t)
	}

	if t != nil && t.Duration == 0 {
		slog.Info("Assigning duration from playback")
		t.Duration = et.eng.duration.Load()
	}
	pb.Skip = et.eng.lastPosition.Load()
	return
}

func (et *events) GetState() gst.State {
	return et.eng.state.Load()
}

func (et *events) HasNextStream() bool {
	pb := &models.Playback{}
	if pb.GetNextToPlay() != nil {
		return true
	}

	q, _ := models.GetActivePerspectiveIndex().GetPerspectiveQueue()
	if !q.IsEmpty() {
		return true
	}

	if et.eng.pt.Load() != nil {
		if _, err := et.eng.pt.Load().GetTrackAfter(false); err == nil {
			return true
		}
	}

	return false
}

func (et *events) IsPaused() bool {
	return et.eng.playbin.Load() != nil && et.eng.state.Load() == gst.StatePaused
}

func (et *events) IsPlaying() bool {
	return et.eng.playbin.Load() != nil && et.eng.state.Load() == gst.StatePlaying
}

func (et *events) IsReady() bool {
	return et.eng.playbin.Load() != nil && et.eng.state.Load() == gst.StateReady
}

func (et *events) IsStreaming() bool {
	return et.IsPaused() || et.IsPlaying()
}

func (et *events) IsStopped() bool {
	return et.eng.lastEvent.Load() == stopAllEvent
}

func (et *events) NextStream() (err error) {
	et.eng.lastEvent.Store(nextEvent)

	if !(et.IsStreaming() || et.IsReady()) {
		return
	}

	et.eng.setPlaybackHint(hintNextInPlaylist)
	et.StopStream()
	models.TriggerPlaybackChange()
	return
}

func (et *events) PauseStream(off bool) (err error) {
	if et.eng.playbin.Load() == nil {
		if !et.IsStopped() {
			return
		}
		et.eng.lastEvent.Store(resumeAllEvent)
		models.TriggerPlaybackChange()
		return
	}

	et.eng.lastEvent.Store(pauseEvent)

	if off {
		if !et.IsPaused() {
			return
		}
		et.eng.state.Store(gst.StatePlaying)
		err = et.eng.playbin.Load().SetState(et.eng.state.Load())
	} else {
		if !et.IsPlaying() {
			return
		}
		et.eng.state.Store(gst.StatePaused)
		err = et.eng.playbin.Load().SetState(et.eng.state.Load())
	}

	broadcastToSubscribers(subscription.ToPlaybackEvent)
	return
}

func (et *events) PlayStreams(force bool, locations []string, ids []int64) {
	et.eng.lastEvent.Store(playEvent)

	logw := slog.With(
		"locations", locations,
		"ids", ids,
	)
	logw.Info("Playing streams")

	for _, v := range locations {
		models.AddPlaybackLocation(v)
	}
	for _, v := range ids {
		t := &models.Track{}
		err := t.Read(v)
		if err != nil {
			logw.With(
				"track_id", v,
				"error", err,
			).Error("Failed to read track")
			continue
		}
		models.AddPlaybackTrack(t)
	}

	if force {
		et.StopStream()
	}
}

func (et *events) PreviousStream() {
	if time.Duration(et.eng.lastPosition.Load())*time.Nanosecond >=
		time.Duration(base.PlaybackPlayedThreshold)*time.Second {
		et.SeekInStream(0)
		return
	}

	et.eng.lastEvent.Store(previousEvent)

	et.eng.clearPendingPlayback()

	var hint playbackHint
	if et.IsStreaming() {
		if et.eng.pt.Load() != nil {
			hint = hintPrevInPlaylist
		} else {
			hint = hintPrevInHistory
		}
	} else {
		hint = hintPrevInHistory
	}
	et.eng.setPlaybackHint(hint)
	et.StopStream()
	models.TriggerPlaybackChange()
}

func (et *events) QuitPlayingFromBar(pl *models.Playlist) {
	et.eng.lastEvent.Store(stopPlaybarEvent)
	slog.Info("Quit playing from bar", "pl", *pl)

	if !(pl.Open && pl.Active) {
		return
	}

	bar := pl.Playbar
	bar.DeactivateEntry(pl)
	et.quitPlayingFromList()
}

func (et *events) SeekInStream(pos int64) {
	et.eng.lastEvent.Store(seekEvent)

	if !et.IsPlaying() {
		return
	}

	if et.eng.seekable.Load() {
		seek := gst.NewSeekEvent(
			1.0,
			gst.FormatTime,
			gst.SeekFlagFlush|gst.SeekFlagKeyUnit,
			gst.SeekTypeSet,
			pos,
			gst.SeekTypeNone,
			-1,
		)

		if !et.eng.playbin.Load().SendEvent(seek) {
			slog.Error("Failed to send playback event", "event", gst.EventTypeSeek)
		}
	}
}

func (et *events) StopAll() {
	et.eng.lastEvent.Store(stopAllEvent)

	et.eng.updateMPRIS(true)
	et.StopStream()
}

func (et *events) StopStream() {
	slog.Info("Stopping current playback")

	if et.eng.lastEvent.Load() != stopAllEvent {
		et.eng.lastEvent.Store(stopStreamEvent)
	}

	if !et.IsStreaming() {
		return
	}

	// NOTE: sending the EOS event has between 1.5 and 3 seconds of
	// latency, and conflicts with the paused state, so we are
	// just ending things here
	if !et.IsPaused() {
		et.eng.state.Store(gst.StatePaused)
		onerror.Log(et.eng.playbin.Load().SetState(et.eng.state.Load()))
	}
	et.eng.wrapUp()
	et.eng.mainLoop.Quit()
}

func (et *events) TryPlayingFromBar(pl *models.Playlist, position int) {
	et.eng.lastEvent.Store(playbarEvent)

	logw := slog.With("pl", pl)
	logw.Info("Try playing from bar")

	bar := models.Playbar{}
	if err := bar.Read(pl.PlaybarID); err != nil {
		logw.With(
			"playbar_id", pl.PlaybarID,
			"error", err,
		).Error("Failed to read playbar")
		return
	}

	bar.ActivateEntry(pl)
	et.tryPlayingFromList(pl, position)
}

func (et *events) quitPlayingFromList() {
	et.eng.lastEvent.Store(stopPlaylistEvent)

	et.eng.setPlaybackHint(hintStopPlaylist)
	et.StopStream()
	models.TriggerPlaybackChange()

	et.eng.updateMPRIS(true)
}

func (et *events) tryPlayingFromList(pl *models.Playlist, position int) {
	et.eng.lastEvent.Store(playlistEvent)

	pt, err := pl.GetTrackAt(position)
	if err != nil {
		slog.With(
			"position", position,
			"error", err,
		).Error("Failed to get track at position")
		return
	}

	et.eng.pt.Store(pt)

	et.eng.setPlaybackHint(hintStartPlaylist)
	et.StopStream()
	models.TriggerPlaybackChange()
}

// GetEventsInstance returns the events instance/object.
func GetEventsInstance() IEvents {
	if instance == nil {
		instance = &events{
			&engine{
				hint: hintNone,
			},
		}
		instance.eng.state.Store(gst.StateNull)
		instance.eng.lastEvent.Store(noLoopEvent)
	}
	return instance
}

// StartEngine starts the playback engine.
func StartEngine() *rtc.Unloader {
	slog.Info("Starting playback engine")

	gst.Init(nil)

	GetEventsInstance()

	instance.eng.resumeActivePlaylist()
	go instance.eng.engineLoop()
	models.TriggerPlaybackChange()

	return unloader
}

func stopEngine() {
	if instance.eng.lastEvent.Load() == noLoopEvent {
		return
	}
	slog.Info("Stopping engine")

	instance.eng.freezePlayback.Store(true)
	instance.StopAll()

	for i := 0; i < base.ServerWaitTimeout; i++ {
		if instance.eng.playbin.Load() != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	quitEngineLoop <- struct{}{}

	for i := 0; i < base.ServerWaitTimeout; i++ {
		if instance.eng.lastEvent.Load() != noLoopEvent {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
}
