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

// Player -
type Player struct {
	*mpris.Instance
	lastPlaybackStatus string
}

// IntrospectInterface implements the mpris.Player interface
func (*Player) IntrospectInterface() introspect.Interface {
	return mpris.PlayerIntrospectInterface()
}

// Properties implements the mpris.Player interface
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

// Next implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Next() *dbus.Error {
	err := NextStream()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

// Previous implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Previous() *dbus.Error {
	PreviousStream()
	return nil
}

// Pause implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Pause() *dbus.Error {
	err := PauseStream(false)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

// PlayPause implements org.mpris.MediaPlayer2.Player interface
func (p *Player) PlayPause() *dbus.Error {
	if IsPlaying() {
		return p.Pause()
	} else if IsPaused() {
		return p.Play()
	}
	return nil
}

// Stop implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Stop() *dbus.Error {
	StopAll()
	return nil
}

// Play implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Play() *dbus.Error {
	err := PauseStream(true)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

// Seek implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Seek(x int64) *dbus.Error {
	return nil
}

// SetPosition implements org.mpris.MediaPlayer2.Player interface
func (p *Player) SetPosition(o string, x int64) *dbus.Error {
	return nil
}

// OpenUri implements org.mpris.MediaPlayer2.Player interface
//nolint: revive // Implements interface
func (p *Player) OpenUri(s string) *dbus.Error {
	return nil
}

// PlaybackStatus implements org.mpris.MediaPlayer2.Player interface
func (p *Player) PlaybackStatus() string {
	if IsPlaying() {
		return PlaybackStatusPlaying
	}
	if IsPaused() {
		return PlaybackStatusPaused
	}
	return PlaybackStatusStopped
}

// LoopStatus implements org.mpris.MediaPlayer2.Player interface
func (p *Player) LoopStatus(s string) (string, *dbus.Error) {
	// TODO: implement
	return "None", nil
}

// Rate implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Rate(in float64) (float64, *dbus.Error) {
	// TODO: implement
	return float64(1.), nil
}

// Shuffle implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Shuffle(b bool) (bool, *dbus.Error) {
	// TODO: implement
	return false, nil
}

// Metadata implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Metadata() map[string]dbus.Variant {
	pb, t := GetPlayback()
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

// Volume implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Volume(in float64) (float64, *dbus.Error) {
	// TODO: implement
	return float64(1.), nil
}

// Position implements org.mpris.MediaPlayer2.Player interface
func (p *Player) Position() int64 {
	pb, _ := GetPlayback()
	return pb.Skip
}

// MinimumRate implements org.mpris.MediaPlayer2.Player interface
func (p *Player) MinimumRate() float64 {
	return float64(1.)
}

// MaximumRate implements org.mpris.MediaPlayer2.Player interface
func (p *Player) MaximumRate() float64 {
	return float64(1.)
}

// CanGoNext implements org.mpris.MediaPlayer2.Player interface
func (p *Player) CanGoNext() bool {
	return HasNextStream()
}

// CanGoPrevious implements org.mpris.MediaPlayer2.Player interface
func (p *Player) CanGoPrevious() bool {
	return true
}

// CanPlay implements org.mpris.MediaPlayer2.Player interface
func (p *Player) CanPlay() bool {
	return true
}

// CanPause implements org.mpris.MediaPlayer2.Player interface
func (p *Player) CanPause() bool {
	return true
}

// CanSeek implements org.mpris.MediaPlayer2.Player interface
func (p *Player) CanSeek() bool {
	return false
}

// CanControl implements org.mpris.MediaPlayer2.Player interface
func (p *Player) CanControl() bool {
	return true
}
