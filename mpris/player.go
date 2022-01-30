package mpris

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
)

// Player -
type Player struct {
	*Instance
}

func (*Player) introspectInterface() introspect.Interface {
	return introspect.Interface{
		Name: "org.mpris.MediaPlayer2.Player",
		Properties: []introspect.Property{
			{Name: "PlaybackStatus", Type: "s", Access: "read"},
			{Name: "LoopStatus", Type: "s", Access: "readwrite"},
			{Name: "Rate", Type: "d", Access: "readwrite"},
			{Name: "Shuffle", Type: "b", Access: "readwrite"},
			{Name: "Metadata", Type: "a{sv}", Access: "read"},
			{Name: "Volume", Type: "d", Access: "readwrite"},
			{Name: "Position", Type: "x", Access: "read"},
			{Name: "MinimumRate", Type: "d", Access: "read"},
			{Name: "MaximumRate", Type: "d", Access: "read"},
			{Name: "CanGoNext", Type: "b", Access: "read"},
			{Name: "CanGoPrevious", Type: "b", Access: "read"},
			{Name: "CanPlay", Type: "b", Access: "read"},
			{Name: "CanSeek", Type: "b", Access: "read"},
			{Name: "CanControl", Type: "b", Access: "read"},
		},
		Signals: []introspect.Signal{
			{
				Name: "Seeked",
				Args: []introspect.Arg{
					{Name: "Position", Type: "x"},
				},
			},
		},
		Methods: []introspect.Method{
			{Name: "Next"},
			{Name: "Previous"},
			{Name: "Pause"},
			{Name: "PlayPause"},
			{Name: "Stop"},
			{Name: "Play"},
			{
				Name: "Seek",
				Args: []introspect.Arg{
					{Name: "Offset", Type: "x", Direction: "in"},
				},
			},
			{
				Name: "SetPosition",
				Args: []introspect.Arg{
					{Name: "TrackId", Type: "o", Direction: "in"},
					{Name: "Position", Type: "x", Direction: "in"},
				},
			},
		},
	}
}

// Next -
func (p *Player) Next() *dbus.Error {
	playback.NextStream()
	return nil
}

// Previous -
func (p *Player) Previous() *dbus.Error {
	playback.PreviousStream()
	return nil
}

// Pause -
func (p *Player) Pause() *dbus.Error {
	playback.PauseStream(false)
	return nil
}

// PlayPause -
func (p *Player) PlayPause() *dbus.Error {
	if playback.IsPlaying() {
		return p.Pause()
	} else if playback.IsPaused() {
		return p.Play()
	}
	return nil
}

// Stop -
func (p *Player) Stop() *dbus.Error {
	playback.StopAll()
	return nil
}

// Play -
func (p *Player) Play() *dbus.Error {
	playback.PauseStream(true)
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
	if playback.IsPlaying() {
		return "Playing"
	}
	if playback.IsPaused() {
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
	pb, _ := playback.GetPlayback()
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
