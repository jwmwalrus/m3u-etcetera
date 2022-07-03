package mpris

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

// Player -
type Player interface {
	IntrospectInterface() introspect.Interface
	Properties() map[string]*prop.Prop

	Next() *dbus.Error
	Previous() *dbus.Error
	Pause() *dbus.Error
	PlayPause() *dbus.Error
	Stop() *dbus.Error
	Play() *dbus.Error
	Seek(x int64) *dbus.Error
	SetPosition(o string, x int64) *dbus.Error
	OpenUri(s string) *dbus.Error

	// Seeked(x int64) *dbus.Error

	PlaybackStatus() string
	LoopStatus(s string) (string, *dbus.Error)
	Rate(in float64) (float64, *dbus.Error)
	Shuffle(b bool) (bool, *dbus.Error)
	Metadata() map[string]dbus.Variant
	Volume(in float64) (float64, *dbus.Error)
	Position() int64
	MinimumRate() float64
	MaximumRate() float64
	CanGoNext() bool
	CanGoPrevious() bool
	CanPlay() bool
	CanPause() bool
	CanSeek() bool
	CanControl() bool
}

// PlayerIntrospectInterface returns the instrospection for the player
func PlayerIntrospectInterface() introspect.Interface {
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
			{Name: "CanPause", Type: "b", Access: "read"},
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
			{
				Name: "OpenUri",
				Args: []introspect.Arg{
					{Name: "Uri", Type: "s", Direction: "in"},
				},
			},
		},
	}
}
