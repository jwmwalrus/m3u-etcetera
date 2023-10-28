package playback

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/gear-pieces/idler"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/mpris"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	rtc "github.com/jwmwalrus/rtcycler"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

const (
	positionThreshold = 1e9
)

var (
	quitEngineLoop chan struct{} = make(chan struct{})

	broadcastToSubscribers = subscription.Broadcast
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

type gstState struct{ atomic.Int32 }

func (gs *gstState) Load() gst.State   { return gst.State(gs.Int32.Load()) }
func (gs *gstState) Store(s gst.State) { gs.Int32.Store(int32(s)) }

type engine struct {
	freezePlayback atomic.Bool
	terminate      atomic.Bool
	seekable       atomic.Bool
	seekableDone   atomic.Bool
	lastPosition   atomic.Int64
	duration       atomic.Int64
	buffering      atomic.Int32
	pt             atomic.Pointer[models.PlaylistTrack]
	pb             atomic.Pointer[models.Playback]
	t              atomic.Pointer[models.Track]
	playbin        atomic.Pointer[gst.Element]
	lastEvent      engineEvent
	prevState      gstState
	state          gstState

	mpris *Player

	mainLoop *glib.MainLoop
	hint     playbackHint
}

func init() {
}

func (e *engine) addPlaybackFromQueue(qt *models.QueueTrack) (pb *models.Playback) {
	slog.Info("Adding playback from queue", "qt", qt)

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
	if e.pb.Load() == nil {
		return
	}

	e.pb.Load().ClearPending()
}

func (e *engine) engineLoop() {
	slog.Info("Starting engine loop")

	e.lastEvent.Store(loopEvent)
loop:
	for {
		pb := &models.Playback{}

	signals:
		select {
		case <-quitEngineLoop:
			break loop
		case <-models.PlaybackChanged:
			switch e.lastEvent.Load() {
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
				e.pt.Store(nil)
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
					slog.Debug("Popped successfully")
					pb = e.addPlaybackFromQueue(qt)
				} else if e.pt.Load() != nil {
					pb = e.getNextInPlaylist(false)
				}
			}
			break signals
		}

		if pb != nil && pb.ID > 0 && pb.Location != "" {
			slog.Debug("There is a playback")
			e.playStream(pb)
			go func() {
				models.TriggerPlaybackChange()
			}()
			continue loop
		}

		if !e.freezePlayback.Load() {
			models.DeactivatePlaybars()
		}
		broadcastToSubscribers(subscription.ToPlaybackEvent)
	}

	e.lastEvent.Store(noLoopEvent)
	slog.Info("Firing event", "event", e.lastEvent)
}

func (e *engine) getFirstInPlaylist() (pb *models.Playback) {
	if e.pt.Load() == nil {
		slog.Error("There is no playlist-track available")
		return
	}
	pl := e.pt.Load().Playlist
	t := e.pt.Load().Track
	if pl.ID == 0 || t.ID == 0 {
		slog.Error("There is no list or track to play from")
		return
	}

	slog.With("pt", *e.pt.Load()).
		Info("Obtaining first track in playlist")

	pb = models.AddPlaybackTrack(&t)
	return
}

func (e *engine) getNextInPlaylist(goingBack bool) (pb *models.Playback) {
	if e.pt.Load() == nil {
		slog.Info("There is no playlist-track available")
		return
	}

	logw := slog.With("pt", *e.pt.Load())
	logw.Info("Obtaining next track in playlist")

	newpt, err := e.pt.Load().GetTrackAfter(goingBack)
	if err != nil {
		pl := e.pt.Load().Playlist
		if pl.ID == 0 {
			logw.Error("There is no list to play from")
			return
		}
		bar := models.Playbar{}
		if errb := bar.Read(pl.PlaybarID); errb == nil {
			bar.DeactivateEntry(&pl)
		}
		e.pt.Store(nil)
		logw.Info("Failed to get track after", "error", err)
		return
	}

	e.pt.Store(newpt)
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
	slog.Info("Obtaining previous track in history")

	h := models.PlaybackHistory{}

	if err = h.ReadLast(); err != nil {
		slog.Error("Failed to read last history entry", "error", err)
		return
	}

	if h.TrackID > 0 {
		t := &models.Track{}
		if err = t.Read(h.TrackID); err != nil {
			slog.Error("Failed to read track", "error", err)
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

func (e *engine) handleBusMessage(msg *gst.Message) bool {
	switch msg.Type() {
	case gst.MessageEOS:
		slog.Debug("End of stream", "location", e.pb.Load().Location)
		e.wrapUp()
		e.mainLoop.Quit()
	case gst.MessageError:
		slog.Error(msg.String())
		e.wrapUp()
		e.mainLoop.Quit()
	case gst.MessageWarning:
		slog.Warn(msg.String())
	case gst.MessageInfo:
		rtc.Trace(msg.String())
	case gst.MessageStateChanged:
		prevState, state := msg.ParseStateChanged()
		e.prevState.Store(prevState)
		e.state.Store(state)
		slog.With(
			"previousState", e.prevState.Load(),
			"newState", e.state.Load(),
		).
			Debug("Pipeline state changed")

		e.updateMPRIS(false)
	case gst.MessageDurationChanged:
		e.duration.Store(0)
	case gst.MessageBuffering:
		e.buffering.Store(int32(msg.ParseBuffering()))
		if e.buffering.Load() < 100 {
			e.state.Store(gst.StatePaused)
		} else {
			e.state.Store(gst.StatePlaying)
		}
		onerror.Log(e.playbin.Load().SetState(e.state.Load()))
	default:
		rtc.Trace(msg.String())
	}

	return true
}

func (e *engine) playStream(pb *models.Playback) {
	logw := slog.With("pb", *pb)
	logw.Info("Starting playStream")

	e.pb.Store(pb)
	e.terminate.Store(false)
	defer e.reset()

	idler.GetBusy(idler.StatusEngineLoop)
	defer func() { idler.GetFree(idler.StatusEngineLoop) }()

	// check if playback is valid
	if e.pb.Load() == nil || pb.Location == "" {
		return
	}
	logw.Debug("Playback is valid")

	var err error

	e.mainLoop = glib.NewMainLoop(glib.MainContextDefault(), false)

	playbin, err := gst.NewElementWithName("playbin", "m3uetc-playbin")
	e.playbin.Store(playbin)
	if err != nil {
		logw.Error("Failed to create playbin", "error", err)
		return
	}
	logw.Debug("Playbin created")

	if e.playbin.Load() == nil {
		logw.Error("Not all elements could be created")
		return
	}

	flags, err := e.playbin.Load().GetProperty("flags")
	if err != nil {
		logw.Error("Unable to get flags", "error", err)
	} else {
		eflags := flags.(uint)
		eflags = eflags &^ (1 << 0) // no video
		eflags = eflags | (1 << 1)  // yes audio
		eflags = eflags &^ (1 << 2) // no text
		e.playbin.Load().SetArg("flags", strconv.FormatInt(int64(eflags), 10))
		fflags, _ := e.playbin.Load().GetProperty("flags")
		if fflags.(uint) != eflags {
			logw.With(
				"initialFlags", flags,
				"expectedFlags", eflags,
				"finalFlags", fflags,
			).Warn("Flags could not be set")
		}
	}

	broadcastToSubscribers(subscription.ToPlaybackEvent)
	e.updateMPRIS(false)

	e.playbin.Load().Set("uri", pb.Location)

	bus := e.playbin.Load().GetBus()

	e.state.Store(gst.StatePlaying)
	if err := e.playbin.Load().SetState(e.state.Load()); err != nil {
		logw.Error("Unable to start playback", "error", err)
		pb.Blacklist()
		return
	}
	logw.Debug("State changed", "state", gst.StatePlaying)

	pqctx, cancelpq := context.WithCancel(context.Background())
	go e.performQueries(pqctx)

	bus.AddWatch(func(msg *gst.Message) bool {
		return e.handleBusMessage(msg)
	})

	e.mainLoop.Run()

	bus.RemoveWatch()

	cancelpq()

	logw.Debug("End of playback")
	e.state.Store(gst.StateNull)
	e.playbin.Load().SetState(e.state.Load())
}

func (e *engine) performQueries(ctx context.Context) {
	tick := time.NewTicker(time.Duration(positionThreshold) * time.Nanosecond)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if e.terminate.Load() {
				break
			}

			if e.state.Load() != gst.StatePlaying {
				continue
			}

			// Query the current position of the stream
			var position int64
			var ok bool
			if ok, position = e.playbin.Load().QueryPosition(gst.FormatTime); !ok {
				slog.Warn("Could not query current position")
			}

			if position > 0 {
				broadcastToSubscribers(subscription.ToPlaybackEvent)
				e.lastPosition.Store(position)
			}

			// If we didn't know it yet, query the stream duration
			if e.duration.Load() == 0 {
				var duration int64
				if ok, duration = e.playbin.Load().QueryDuration(gst.FormatTime); !ok {
					slog.Warn("Could not query current duration")
				}
				e.duration.Store(duration)
				broadcastToSubscribers(subscription.ToPlaybackEvent)
			}

			if !e.seekableDone.Load() {
				q := gst.NewSeekingQuery(gst.FormatTime)
				if e.playbin.Load().Query(q) {
					var start, end int64
					var format gst.Format
					var seekable bool
					format, seekable, start, end = q.ParseSeeking()
					e.seekable.Store(seekable)
					if e.seekable.Load() {
						slog.With(
							"format", format,
							"start", start,
							"end", end,
						).
							Debug("Seeking is ENABLED")
						go func() {
							if e.pb.Load().Skip > 0 {
								GetEventsInstance().SeekInStream(e.pb.Load().Skip)
							}
						}()
					} else {
						slog.Debug("Seeking is DISABLED for this stream")
					}
				} else {
					slog.Debug("Seeking query failed")
				}
				e.seekableDone.Store(true)
			}
		}
	}
}

func (e *engine) resumeActivePlaylist() {
	e.pt.Store(models.GetActivePlaylistTrack())
	if e.pt.Load() != nil {
		slog.Info("Resuming playback for active playlist")
	}
}

func (e *engine) reset() {
	e.pb.Store(nil)
	e.t.Store(nil)
	e.seekable.Store(false)
	e.seekableDone.Store(false)
	e.playbin.Store(nil)

	e.lastPosition.Store(0)
	e.duration.Store(0)
	e.buffering.Store(0)
}

func (e *engine) setPlaybackHint(h playbackHint) {
	e.hint = h
}

func (e *engine) updateMPRIS(destroy bool) {
	deleteMPRIS := func() {
		if e.mpris != nil {
			e.mpris.Delete()
		}
		e.mpris = nil
	}

	if destroy {
		deleteMPRIS()
		return
	}

	if e.mpris == nil {
		mprisInstance := mpris.New()
		e.mpris = &Player{mprisInstance, PlaybackStatusStopped}
		if err := mprisInstance.Setup(e.mpris); err != nil {
			slog.Error("Failed to setup mpris instance", "error", err)
			deleteMPRIS()
		}
		return
	}

	currPbStatus := e.mpris.PlaybackStatus()
	if currPbStatus == e.mpris.lastPlaybackStatus {
		return
	}

	e.mpris.lastPlaybackStatus = currPbStatus
	err := e.mpris.Conn.Load().Emit(
		mpris.RootPath,
		mpris.PropertiesInterface+".PropertiesChanged",
		mpris.PlayerInterface,
		map[string]dbus.Variant{
			"PlaybackStatus": dbus.MakeVariant(currPbStatus),
			"Metadata":       dbus.MakeVariant(e.mpris.Metadata()),
			"CanGoNext":      dbus.MakeVariant(e.mpris.CanGoNext()),
		},
		[]string{},
	)
	onerror.Warn(err)
}

func (e *engine) wrapUp() {
	defer e.terminate.Store(true)

	if e.getPlaybackHint(true) == hintPrevInHistory {
		return
	}

	go models.AddPlaybackToHistory(
		e.pb.Load().ID,
		e.lastPosition.Load(),
		e.duration.Load(),
		e.freezePlayback.Load(),
	)
}
