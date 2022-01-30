package mpris

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
)

// MediaPlayer2 -
type MediaPlayer2 struct {
	*Instance
}

func (*MediaPlayer2) introspectInterface() introspect.Interface {
	return introspect.Interface{
		Name: "org.mpris.MediaPlayer2",
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
func (p *Player) Fullscreen(b bool) (bool, *dbus.Error) {
	// TODO: implement
	return false, nil
}

// CanSetFullscreen -
func (p *Player) CanSetFullscreen(b bool) bool {
	// TODO: implement
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
