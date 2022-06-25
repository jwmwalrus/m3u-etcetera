package mpris

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
)

// MediaPlayer2 -
type MediaPlayer2 struct {
	*Instance
}

func (*MediaPlayer2) introspectInterface() introspect.Interface {
	return introspect.Interface{
		Name: RootInterface,
		Properties: []introspect.Property{
			{Name: "CanQuit", Type: "b", Access: "read"},
			{Name: "Fullscreen", Type: "b", Access: "readwrite"},
			{Name: "CanSetFullscreen", Type: "b", Access: "read"},
			{Name: "CanRaise", Type: "b", Access: "read"},
			{Name: "HasTrackList", Type: "b", Access: "read"},
			{Name: "Identity", Type: "s", Access: "read"},
			{Name: "DesktopEntry", Type: "s", Access: "read"},
			{Name: "SupportedUriSchemes", Type: "as", Access: "read"},
			{Name: "SupportedMimeTypes", Type: "as", Access: "read"},
		},
		Methods: []introspect.Method{
			{Name: "Raise"},
			{Name: "Quit"},
		},
	}
}

func (mp2 *MediaPlayer2) properties() map[string]*prop.Prop {
	return map[string]*prop.Prop{
		"CanQuit":             {Value: mp2.CanQuit(), Emit: prop.EmitTrue},
		"Fullscreen":          {Value: false, Writable: true, Emit: prop.EmitTrue},
		"CanSetFullscreen":    {Value: mp2.CanSetFullscreen(), Emit: prop.EmitTrue},
		"CanRaise":            {Value: mp2.CanRaise(), Emit: prop.EmitTrue},
		"HasTrackList":        {Value: mp2.HasTrackList(), Emit: prop.EmitTrue},
		"Identity":            {Value: mp2.Identity(), Emit: prop.EmitTrue},
		"DesktopEntry":        {Value: mp2.DesktopEntry(), Emit: prop.EmitTrue},
		"SupportedUriSchemes": {Value: mp2.SupportedUriSchemes(), Emit: prop.EmitTrue},
		"SupportedMimeTypes":  {Value: mp2.SupportedMimeTypes(), Emit: prop.EmitTrue},
	}
}

// Quit -
func (mp2 *MediaPlayer2) Quit() *dbus.Error {
	// TODO: implement
	return nil
}

// Raise -
func (mp2 *MediaPlayer2) Raise() *dbus.Error {
	return nil
}

// CanQuit -
func (mp2 *MediaPlayer2) CanQuit() bool {
	return true
}

// Fullscreen -
func (mp2 *MediaPlayer2) Fullscreen(b bool) (bool, *dbus.Error) {
	return false, nil
}

// CanSetFullscreen -
func (mp2 *MediaPlayer2) CanSetFullscreen() bool {
	return false
}

// CanRaise -
func (mp2 *MediaPlayer2) CanRaise() bool {
	return false
}

// HasTrackList -
func (mp2 *MediaPlayer2) HasTrackList() bool {
	return false
}

// Identity -
func (mp2 *MediaPlayer2) Identity() string {
	return base.AppName
}

// DesktopEntry -
func (mp2 *MediaPlayer2) DesktopEntry() string {
	return "m3uetc-server"
}

// SupportedUriSchemes -
func (mp2 *MediaPlayer2) SupportedUriSchemes() []string {
	return base.SupportedURISchemes
}

// SupportedMimeTypes -
func (mp2 *MediaPlayer2) SupportedMimeTypes() []string {
	return base.SupportedMIMETypes
}
