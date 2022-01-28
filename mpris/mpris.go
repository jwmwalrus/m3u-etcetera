package mpris

import (
	"github.com/godbus/dbus/v5"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"

	"github.com/godbus/dbus/v5/introspect"
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
	ins.conn, err = dbus.SessionBus()
	onerror.Panic(err)

	ins.Unloader = base.Unloader{
		Description: "MPRIS-Server",
		Callback: func() error {
			ins.conn.ReleaseName(serverName)
			ins.conn.Close()
			return nil
		},
	}

	mp2 := &MediaPlayer2{ins}
	err = ins.conn.Export(mp2, rootPath, rootInterface)
	onerror.Panic(err)

	player := &Player{ins}
	err = ins.conn.Export(player, rootPath, playerInterface)
	onerror.Panic(err)

	err = ins.conn.Export(
		introspect.NewIntrospectable(&introspect.Node{
			Name: serverName,
			Interfaces: []introspect.Interface{
				ins.introspectInterface(),
				mp2.introspectInterface(),
				player.introspectInterface(),
			},
		}),
		rootPath,
		"org.freedesktop.DBus.Introspectable",
	)
	onerror.Panic(err)

	reply, err := ins.conn.RequestName(serverName, dbus.NameFlagDoNotQueue|dbus.NameFlagReplaceExisting)
	onerror.Panic(err)

	if reply != dbus.RequestNameReplyPrimaryOwner {
		panic("D-Bus name already taken")
	}

	return ins
}
