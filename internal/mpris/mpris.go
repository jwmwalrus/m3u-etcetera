package mpris

import (
	"github.com/godbus/dbus/v5"
	"github.com/jwmwalrus/onerror"
)

const (
	rootPath        = "/org/mpris/MediaPlayer2"
	rootInterface   = "org.mpris.MediaPlayer2"
	playerInterface = rootInterface + ".Player"

	serverName = rootInterface + ".M3UEtcetera"
)

// New returns a new dbus instance implementing MPRIS
func New() *Instance {
	var err error

	ins := &Instance{name: serverName}
	ins.conn, err = dbus.ConnectSessionBus()
	onerror.Panic(err)

	return ins
}
