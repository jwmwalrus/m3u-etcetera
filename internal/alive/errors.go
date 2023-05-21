package alive

import "strconv"

// ServerAlreadyRunning returned when requested to start a server that is already running.
type ServerAlreadyRunning struct{}

func (*ServerAlreadyRunning) Error() string {
	return "Server already running!"
}

// ServerNotRunning returned when requested to stop a server that is not running.
type ServerNotRunning struct{}

func (*ServerNotRunning) Error() string {
	return "Server not running!"
}

// ServerStarted returned when successfully attempted to start a server.
// It holds an optional description, `Desc`, and the `PID` associated
// with the server.
type ServerStarted struct {
	Desc string
	PID  int
}

func (e *ServerStarted) Error() (out string) {
	out = "Server started!"
	if e.Desc != "" {
		out += " (" + e.Desc + ")"
	}
	if e.PID != 0 {
		out += " ( PID: " + strconv.Itoa(e.PID) + ")"
	}
	return
}

// ServerStopped returned when successfully attempted to stop a server.
// It holds an optional description, `Desc`.
type ServerStopped struct {
	Desc string
}

func (e *ServerStopped) Error() (out string) {
	out = "Server stopped!"
	if e.Desc != "" {
		out += " (" + e.Desc + ")"
	}
	return
}
