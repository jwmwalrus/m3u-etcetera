package mpris

import (
	"fmt"

	"github.com/godbus/dbus/v5"

	"github.com/godbus/dbus/v5/introspect"
)

// Instance -
type Instance struct {
	name string
	conn *dbus.Conn
}

func (*Instance) introspectInterface() introspect.Interface {
	return introspect.IntrospectData
}

func (i *Instance) Delete() error {
	i.conn.ReleaseName(serverName)
	i.conn.Close()
	return nil
}

func (i *Instance) Setup(p Player) (err error) {
	mp2 := &MediaPlayer2{i}
	err = i.conn.Export(mp2, rootPath, rootInterface)
	if err != nil {
		return
	}

	err = i.conn.Export(p, rootPath, playerInterface)
	if err != nil {
		return
	}

	err = i.conn.Export(
		introspect.NewIntrospectable(&introspect.Node{
			Name: serverName,
			Interfaces: []introspect.Interface{
				i.introspectInterface(),
				mp2.introspectInterface(),
				p.IntrospectInterface(),
			},
		}),
		rootPath,
		"org.freedesktop.DBus.Introspectable",
	)
	if err != nil {
		return
	}

	reply, err := i.conn.RequestName(
		serverName,
		dbus.NameFlagDoNotQueue|dbus.NameFlagReplaceExisting,
	)
	if err != nil {
		return
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		err = fmt.Errorf("D-Bus name already taken")
		return
	}

	return
}
