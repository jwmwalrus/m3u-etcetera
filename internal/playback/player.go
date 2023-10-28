package playback

import (
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	"github.com/jwmwalrus/m3u-etcetera/internal/mpris"
)

// Defined PlaybackStatuses.
const (
	PlaybackStatusPlaying = "Playing"
	PlaybackStatusPaused  = "Paused"
	PlaybackStatusStopped = "Stopped"
)

// Player implements the mpris.Player interface.
type Player struct {
	*mpris.Instance
	lastPlaybackStatus string
}

func (*Player) IntrospectInterface() introspect.Interface {
	return mpris.PlayerIntrospectInterface()
}

func (p *Player) Properties() map[string]*prop.Prop {
	return map[string]*prop.Prop{
		"PlaybackStatus": {Value: p.PlaybackStatus(), Emit: prop.EmitTrue},
		"LoopStatus":     {Value: "None", Writable: true, Emit: prop.EmitTrue},
		"Rate":           {Value: float64(1.), Writable: true, Emit: prop.EmitTrue},
		"Shuffle":        {Value: false, Writable: true, Emit: prop.EmitTrue},
		"Metadata":       {Value: p.Metadata(), Emit: prop.EmitTrue},
		"Volume":         {Value: float64(1.), Writable: true, Emit: prop.EmitTrue},
		"Position":       {Value: p.Position(), Emit: prop.EmitTrue},
		"MinimumRate":    {Value: p.MinimumRate(), Emit: prop.EmitTrue},
		"MaximumRate":    {Value: p.MaximumRate(), Emit: prop.EmitTrue},
		"CanGoNext":      {Value: p.CanGoNext(), Emit: prop.EmitTrue},
		"CanGoPrevious":  {Value: p.CanGoPrevious(), Emit: prop.EmitTrue},
		"CanPlay":        {Value: p.CanPlay(), Emit: prop.EmitTrue},
		"CanPause":       {Value: p.CanPause(), Emit: prop.EmitTrue},
		"CanSeek":        {Value: p.CanSeek(), Emit: prop.EmitTrue},
		"CanControl":     {Value: p.CanControl(), Emit: prop.EmitTrue},
	}
}

func (*Player) Next() *dbus.Error {
	err := GetEventsInstance().NextStream()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (*Player) Previous() *dbus.Error {
	GetEventsInstance().PreviousStream()
	return nil
}

func (*Player) Pause() *dbus.Error {
	err := GetEventsInstance().PauseStream(false)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (p *Player) PlayPause() *dbus.Error {
	if GetEventsInstance().IsPlaying() {
		return p.Pause()
	} else if GetEventsInstance().IsPaused() {
		return p.Play()
	}
	return nil
}

func (*Player) Stop() *dbus.Error {
	GetEventsInstance().StopAll()
	return nil
}

func (*Player) Play() *dbus.Error {
	err := GetEventsInstance().PauseStream(true)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (*Player) Seek(x int64) *dbus.Error {
	return nil
}

func (*Player) SetPosition(o string, x int64) *dbus.Error {
	return nil
}

func (*Player) OpenUri(s string) *dbus.Error {
	return nil
}

func (*Player) PlaybackStatus() string {
	if GetEventsInstance().IsPlaying() {
		return PlaybackStatusPlaying
	}
	if GetEventsInstance().IsPaused() {
		return PlaybackStatusPaused
	}
	return PlaybackStatusStopped
}

func (*Player) LoopStatus(s string) (string, *dbus.Error) {
	// TODO: implement
	return "None", nil
}

func (*Player) Rate(in float64) (float64, *dbus.Error) {
	// TODO: implement
	return float64(1.), nil
}

func (*Player) Shuffle(b bool) (bool, *dbus.Error) {
	// TODO: implement
	return false, nil
}

func (*Player) Metadata() map[string]dbus.Variant {
	pb, t := GetEventsInstance().GetPlayback()
	if t != nil {
		return map[string]dbus.Variant{
			"xesam:album":          dbus.MakeVariant(t.Album),
			"xesam:title":          dbus.MakeVariant(t.Title),
			"xesam:url":            dbus.MakeVariant(t.Location),
			"xesam:contentCreated": dbus.MakeVariant(t.Year),
			"xesam:albumArtist":    dbus.MakeVariant([]string{t.Albumartist}),
			"xesam:artist":         dbus.MakeVariant([]string{t.Artist}),
			"xesam:genre":          dbus.MakeVariant([]string{t.Genre}),
			"xesam:composer":       dbus.MakeVariant([]string{t.Composer}),
			"xesam:trackNumber":    dbus.MakeVariant(t.Tracknumber),
			"xesam:discNumber":     dbus.MakeVariant(t.Discnumber),
			"mpris:artUrl":         dbus.MakeVariant(t.Cover),
			"mpris:length":         dbus.MakeVariant(time.Duration(t.Duration) / time.Microsecond),
			"mpris:trackid":        dbus.MakeVariant(t.ID),
		}
	}
	return map[string]dbus.Variant{
		"mpris:trackid": dbus.MakeVariant(pb.ID),
		"xesam:url":     dbus.MakeVariant(pb.Location),
	}
}

func (*Player) Volume(in float64) (float64, *dbus.Error) {
	// TODO: implement
	return float64(1.), nil
}

func (*Player) Position() int64 {
	pb, _ := GetEventsInstance().GetPlayback()
	return pb.Skip
}

func (*Player) MinimumRate() float64 {
	return float64(1.)
}

func (*Player) MaximumRate() float64 {
	return float64(1.)
}

func (*Player) CanGoNext() bool {
	return GetEventsInstance().HasNextStream()
}

func (*Player) CanGoPrevious() bool {
	return true
}

func (*Player) CanPlay() bool {
	return true
}

func (*Player) CanPause() bool {
	return true
}

func (*Player) CanSeek() bool {
	return false
}

func (*Player) CanControl() bool {
	return true
}
