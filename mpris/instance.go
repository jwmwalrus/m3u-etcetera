package mpris

import (
	"github.com/godbus/dbus/v5"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"

	"github.com/godbus/dbus/v5/introspect"
)

// Instance -
type Instance struct {
	name     string
	conn     *dbus.Conn
	Unloader base.Unloader
}

func (*Instance) introspectInterface() introspect.Interface {
	return introspect.IntrospectData
}
