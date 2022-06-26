package mpris

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/jwmwalrus/onerror"

	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

const (
	RootPath            = "/org/mpris/MediaPlayer2"
	RootInterface       = "org.mpris.MediaPlayer2"
	PlayerInterface     = RootInterface + ".Player"
	PropertiesInterface = "org.freedesktop.DBus.Properties"

	name       = "M3UEtcetera"
	serverName = RootInterface + "." + name
)

// New returns a new dbus instance implementing MPRIS
func New() *Instance {
	var err error

	ins := &Instance{Name: serverName}
	ins.Conn, err = dbus.ConnectSessionBus()
	onerror.Panic(err)

	return ins
}

// Instance -
type Instance struct {
	Name string
	Conn *dbus.Conn
	// Obj   *dbus.Object
	props *prop.Properties
}

func (*Instance) introspectInterface() introspect.Interface {
	return introspect.IntrospectData
}

func (i *Instance) Delete() error {
	i.Conn.ReleaseName(serverName)
	i.Conn.Close()
	return nil
}

func (i *Instance) Setup(p Player) (err error) {
	mp2 := &MediaPlayer2{i}
	err = i.Conn.Export(mp2, RootPath, RootInterface)
	if err != nil {
		return
	}

	err = i.Conn.Export(p, RootPath, PlayerInterface)
	if err != nil {
		return
	}

	err = i.Conn.Export(
		introspect.NewIntrospectable(&introspect.Node{
			Name: serverName,
			Interfaces: []introspect.Interface{
				i.introspectInterface(),
				mp2.introspectInterface(),
				p.IntrospectInterface(),
			},
		}),
		RootPath,
		"org.freedesktop.DBus.Introspectable",
	)
	if err != nil {
		return
	}

	i.props, err = prop.Export(i.Conn, RootPath, map[string]map[string]*prop.Prop{
		RootInterface:   mp2.properties(),
		PlayerInterface: p.Properties(),
	})
	if err != nil {
		return
	}

	reply, err := i.Conn.RequestName(
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

	// i.Obj = i.Conn.Object(serverName, RootPath).(*dbus.Object)
	return
}
