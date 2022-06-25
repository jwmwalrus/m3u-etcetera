package playback

import (
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	"github.com/jwmwalrus/m3u-etcetera/internal/mpris"
	"github.com/jwmwalrus/onerror"
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
		return PlaybackStatusPlaying
	}
	if IsPaused() {
		return PlaybackStatusPaused
	}
	return PlaybackStatusStopped
}

// LoopStatus -
func (p *Player) LoopStatus(s string) (string, *dbus.Error) {
	// TODO: implement
	return "None", nil
}

// Rate -
func (p *Player) Rate(in float64) (float64, *dbus.Error) {
	// TODO: implement
	return float64(1.), nil
}

// LoopStatus -
func (p *Player) Shuffle(b bool) (bool, *dbus.Error) {
	// TODO: implement
	return false, nil
}

// Metadata -
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

// Volume -
func (p *Player) Volume(in float64) (float64, *dbus.Error) {
	// TODO: implement
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

func (p *Player) PropertiesChanged() {
	err := p.Conn.Emit(mpris.RootPath,
		mpris.PropertiesInterface+".PropertiesChanged",
		map[string]dbus.Variant{
			"PlaybackStatus": dbus.MakeVariant(p.PlaybackStatus()),
			"Metadata":       dbus.MakeVariant(p.Metadata()),
			"CanGoNext":      dbus.MakeVariant(p.CanGoNext()),
		},
		[]string{},
	)
	onerror.Log(err)
}
