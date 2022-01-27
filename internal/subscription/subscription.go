package subscription

import (
	"sync"

	"github.com/jwmwalrus/bnp/stringing"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
)

// Type defines the subscription type
type Type int

// Subscription type events
const (
	ToNone = iota
	ToPlaybackEvent
	ToQueueStoreEvent
	ToCollectionStoreEvent
	ToQueryStoreEvent
	ToPlaybarStoreEvent
)

func (st Type) String() string {
	return []string{
		"subscribed-to-none",
		"subscribed-to-playback-event",
		"subscribed-to-queue-store-event",
		"subscribed-to-query-store-event",
		"subscribed-to-playbar-store-event",
	}[st]
}

// Subscriber defines the subscriber type
type Subscriber struct {
	st          Type
	id          string
	Event       chan Event
	Unsubscribe func()
}

// Event defines the event type
type Event struct {
	// Type custom even index/type, shall be a non-negative enum
	Idx int
	// Data event data
	Data interface{}
}

type unsubscribeData struct {
	st Type
	id string
}

var (
	subscriptors []Subscriber
	smu          sync.Mutex
	unloading    = false

	// Unloader declares the subscription unloader
	Unloader = base.Unloader{
		Description: "UnsubscribeAll",
		Callback:    unloadSubscriptions,
	}
)

// Subscribe subscribes a client to the given event|store|model
func Subscribe(st Type) (*Subscriber, string) {
	log.WithField("st", st).
		Info("Subscribing")

	s := Subscriber{
		st:    st,
		id:    stringing.GetRandomString(16),
		Event: make(chan Event),
	}
	s.Unsubscribe = func() {
		removeSubscription(s.id)
	}

	smu.Lock()
	subscriptors = append(subscriptors, s)
	smu.Unlock()
	return &s, s.id
}

// Broadcast sends data to all subscribers
func Broadcast(st Type, es ...Event) {
	log.WithFields(log.Fields{
		"st": st,
		"es": es,
	}).
		Trace("Broadcasting")

	if st == ToNone && len(es) > 0 {
		id, ok := es[0].Data.(string)
		if !ok {
			return
		}
		smu.Lock()
		for i := range subscriptors {
			if subscriptors[i].id != id {
				continue
			}
			evt := Event{
				Data: unsubscribeData{st: st, id: id},
			}
			subscriptors[i].Event <- evt
			break
		}
		smu.Unlock()
		return
	}

	list := es
	if len(list) == 0 {
		list = []Event{{}}
	}

	smu.Lock()
	for i := range subscriptors {
		if subscriptors[i].st == st {
			go func(i int) {
				for _, x := range list {
					subscriptors[i].Event <- x
				}
			}(i)
		}
	}
	smu.Unlock()
}

// MustUnsubscribe checks if the event means unsubscribing
func (s *Subscriber) MustUnsubscribe(e Event) bool {
	log.WithField("e", e).
		Trace("Checking if must unsubscribe")

	ud, ok := e.Data.(unsubscribeData)
	if !ok {
		return false
	}
	return ud.st == ToNone && ud.id == s.id
}

func removeSubscription(id string) {
	log.WithField("id", id).
		Info("Removing subscription")

	if unloading {
		return
	}

	smu.Lock()
	for k := range subscriptors {
		if subscriptors[k].id == id {
			subscriptors[k] = subscriptors[len(subscriptors)-1]
			subscriptors = subscriptors[:len(subscriptors)-1]
			break
		}
	}
	smu.Unlock()
}

func unloadSubscriptions() error {
	unloading = true

	smu.Lock()
	for i := range subscriptors {
		evt := Event{
			Data: unsubscribeData{st: ToNone, id: subscriptors[i].id},
		}
		subscriptors[i].Event <- evt
	}
	smu.Unlock()
	return nil
}
