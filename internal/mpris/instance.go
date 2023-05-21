package mpris

import (
	"fmt"
	"sync/atomic"

	"github.com/godbus/dbus/v5"
	"github.com/jwmwalrus/onerror"

	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

// MPRIs interfaces.
const (
	RootPath            = "/org/mpris/MediaPlayer2"
	RootInterface       = "org.mpris.MediaPlayer2"
	PlayerInterface     = RootInterface + ".Player"
	PropertiesInterface = "org.freedesktop.DBus.Properties"

	name       = "M3UEtcetera"
	serverName = RootInterface + "." + name
)

// New returns a new dbus instance implementing MPRIS.
func New() *Instance {
	ins := &Instance{Name: serverName}
	conn, err := dbus.ConnectSessionBus()
	onerror.Panic(err)
	ins.Conn.Store(conn)

	return ins
}

// Instance -.
type Instance struct {
	Name string
	Conn atomic.Pointer[dbus.Conn]
	// Obj   *dbus.Object
	props atomic.Pointer[prop.Properties]
}

func (*Instance) introspectInterface() introspect.Interface {
	return introspect.IntrospectData
}

// Delete closes the instance connection.
func (i *Instance) Delete() error {

	i.Conn.Load().ReleaseName(serverName)
	i.Conn.Load().Close()
	return nil
}

// Setup sets the player.
func (i *Instance) Setup(p Player) (err error) {
	mp2 := &MediaPlayer2{i}
	err = i.Conn.Load().Export(mp2, RootPath, RootInterface)
	if err != nil {
		return
	}

	err = i.Conn.Load().Export(p, RootPath, PlayerInterface)
	if err != nil {
		return
	}

	err = i.Conn.Load().Export(
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

	props, err := prop.Export(i.Conn.Load(), RootPath, map[string]map[string]*prop.Prop{
		RootInterface:   mp2.properties(),
		PlayerInterface: p.Properties(),
	})
	i.props.Store(props)
	if err != nil {
		return
	}

	reply, err := i.Conn.Load().RequestName(
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
