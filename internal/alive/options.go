package alive

type Option func(*aliveSrv)

// WithTurnOff sets the `turnOff` option to true or to the passed argument value.
// If the resulting option is true, it requests the running server to turn off.
func WithTurnOff(turnOff ...bool) Option {
	return func(a *aliveSrv) {
		a.turnOff = true
		if len(turnOff) > 0 {
			a.turnOff = turnOff[0]
		}
	}
}

// WithForceOff sets the `forceOff` option to true or to the passed argument value.
// This option implies WithTurnOff(). If the resulting option is true, it sends
// the `force` flag when requesting the running server to turn ogg.
func WithForceOff(forceOff ...bool) Option {
	return func(a *aliveSrv) {
		a.forceOff = true
		if len(forceOff) > 0 {
			a.forceOff = forceOff[0]
		}
	}
}

// WithNoWait sets the `noWait` option to true or to the passed argument value.
// If the resulting value is true, the serve function will process the request
// without waiting to confirm if the server is running/stopped.
func WithNoWait(noWait ...bool) Option {
	return func(a *aliveSrv) {
		a.noWait = true
		if len(noWait) > 0 {
			a.noWait = noWait[0]
		}
	}
}
