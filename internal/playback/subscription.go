package playback

import "github.com/jwmwalrus/m3u-etcetera/internal/base"

// EventSubscription defines the event subscription type
type EventSubscription int

const (
	// SubscribedToPlayback -
	SubscribedToPlayback EventSubscription = iota
)

func (st EventSubscription) String() string {
	return []string{
		"subscribed-to-playback",
	}[st]
}

// Subscriber defines the subscriber type
type Subscriber struct {
	id          string
	es          EventSubscription
	Data        chan interface{}
	Unsubscribe func()
}

var subscriptors []Subscriber

// Subscribe subscribes a client to the given event
func Subscribe(es EventSubscription) *Subscriber {
	s := Subscriber{
		id:   base.GetRandomString(16),
		es:   es,
		Data: make(chan interface{}),
	}
	s.Unsubscribe = func() {
		removeSubscription(s.id)
	}

	subscriptors = append(subscriptors, s)
	return &s
}

func broadcastEvent(es EventSubscription, data interface{}) {
	for i := range subscriptors {
		if subscriptors[i].es == es {
			go func() { subscriptors[i].Data <- data }()
		}
	}
}

func removeSubscription(id string) {
	for k := range subscriptors {
		if subscriptors[k].id == id {
			subscriptors[k] = subscriptors[len(subscriptors)-1]
			subscriptors = subscriptors[:len(subscriptors)-1]
			return
		}
	}
}
