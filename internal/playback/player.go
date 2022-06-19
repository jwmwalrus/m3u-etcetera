package playback

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/jwmwalrus/m3u-etcetera/internal/mpris"
)

// Player -
type Player struct {
	*mpris.Instance
}

func (*Player) IntrospectInterface() introspect.Interface {
	return mpris.PlayerIntrospectInterface()
}

// Next -
func (p *Player) Next() *dbus.Error {
	NextStream()
	return nil
}

// Previous -
func (p *Player) Previous() *dbus.Error {
	PreviousStream()
	return nil
}

// Pause -
func (p *Player) Pause() *dbus.Error {
	PauseStream(false)
	return nil
}

// PlayPause -
func (p *Player) PlayPause() *dbus.Error {
	if IsPlaying() {
		return p.Pause()
	} else if IsPaused() {
		return p.Play()
	}
	return nil
}

// Stop -
func (p *Player) Stop() *dbus.Error {
	StopAll()
	return nil
}

// Play -
func (p *Player) Play() *dbus.Error {
	PauseStream(true)
	return nil
}

// Seek -
func (p *Player) Seek(x int64) *dbus.Error {
	return nil
}

// SetPosition -
func (p *Player) SetPosition(o string, x int64) *dbus.Error {
	return nil
}

// OpenUri -
func (p *Player) OpenUri(s string) *dbus.Error {
	return nil
}

// PlaybackStatus -
func (p *Player) PlaybackStatus() string {
	if IsPlaying() {
		return "Playing"
	}
	if IsPaused() {
		return "Paused"
	}
	return "Stopped"
}

// LoopStatus -
func (p *Player) LoopStatus(s string) (string, *dbus.Error) {
	// TODO: implement
	return "None", nil
}

// Rate -
func (p *Player) Rate(in float64) (float64, *dbus.Error) {
	return float64(1.), nil
}

// LoopStatus -
func (p *Player) Shuffle(b bool) (bool, *dbus.Error) {
	// TODO: implement
	return false, nil
}

// Metadata -
func (p *Player) Metadata() map[string]dbus.Variant {
	return map[string]dbus.Variant{}
}

// Volume -
func (p *Player) Volume(in float64) (float64, *dbus.Error) {
	return float64(1.), nil
}

// Position -
func (p *Player) Position() int64 {
	pb, _ := GetPlayback()
	return pb.Skip
}

// MinimumRate -
func (p *Player) MinimumRate() float64 {
	return float64(1.)
}

// MaximumRate -
func (p *Player) MaximumRate() float64 {
	return float64(1.)
}

// CanGoNext -
func (p *Player) CanGoNext() bool {
	return true
}

// CanGoPrevious -
func (p *Player) CanGoPrevious() bool {
	return true
}

// CanPlay -
func (p *Player) CanPlay() bool {
	return true
}

// CanPause -
func (p *Player) CanPause() bool {
	return true
}

// CanSeek -
func (p *Player) CanSeek() bool {
	return false
}

// CanControl -
func (p *Player) CanControl() bool {
	return true
}
