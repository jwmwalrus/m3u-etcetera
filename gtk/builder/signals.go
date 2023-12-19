package builder

import (
	"log/slog"

	rtc "github.com/jwmwalrus/rtcycler"
)

type signalDetail struct {
	Signal  string
	Handler interface{}
}
type Signals map[string][]signalDetail

// AddDetail adds the given signal to the map.
func (s Signals) AddDetail(id, signal string, handler interface{}) {
	list, ok := s[id]
	if ok {
		list = append(list, signalDetail{signal, handler})
		s[id] = list
		return
	}

	s[id] = []signalDetail{{signal, handler}}
}

// ConnectSignals connects the signals map.
func ConnectSignals(signals *Signals) {
	for k, list := range *signals {
		for _, v := range list {
			slog.With(
				"id", k,
				"signal-detail", v.Signal,
			).Debug("Connecting signal")

			obj := app.GetObject(k)
			if obj == nil {
				rtc.With(
					"id", k,
					"signal-detail", v.Signal,
				).Fatal("Unable to connect signal")
			}

			obj.Connect(v.Signal, v.Handler)
		}
	}
}
