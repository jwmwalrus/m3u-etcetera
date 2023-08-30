package subscription

import (
	"log/slog"
	"sync"

	"github.com/jwmwalrus/bnp/ing2"
	rtc "github.com/jwmwalrus/rtcycler"
)

// Type defines the subscription type.
type Type int

// Subscription type events.
const (
	ToNone = iota
	ToPlaybackEvent
	ToQueueStoreEvent
	ToCollectionStoreEvent
	ToQueryStoreEvent
	ToPlaybarStoreEvent
	ToPerspectiveEvent
)

func (st Type) String() string {
	return []string{
		"subscribed-to-none",
		"subscribed-to-playback-event",
		"subscribed-to-queue-store-event",
		"subscribed-to-collection-store-event",
		"subscribed-to-query-store-event",
		"subscribed-to-playbar-store-event",
		"subscribed-to-perspective-event",
	}[st]
}

// Subscriber defines the subscriber type.
type Subscriber struct {
	st          Type
	id          string
	Event       chan Event
	Unsubscribe func()
}

// Event defines the event type.
type Event struct {
	// Idx custom even index/type, shall be a non-negative enum.
	Idx int

	// Data event data.
	Data interface{}
}

type unsubscribeData struct {
	st Type
	id string
}

var (
	subscriptors struct {
		s  []Subscriber
		mu sync.Mutex
	}
	unloading = false

	// Unloader declares the subscription unloader.
	Unloader = &rtc.Unloader{
		Description: "UnsubscribeAll",
		Callback:    unloadSubscriptions,
	}
)

// Subscribe subscribes a client to the given event|store|model.
func Subscribe(st Type) (*Subscriber, string) {
	slog.Info("Subscribing", "st", st)

	rl, _ := ing2.GetRandomLetters(16)
	s := Subscriber{
		st:    st,
		id:    rl,
		Event: make(chan Event),
	}
	s.Unsubscribe = func() {
		removeSubscription(s.id)
	}

	subscriptors.mu.Lock()
	defer subscriptors.mu.Unlock()

	subscriptors.s = append(subscriptors.s, s)
	return &s, s.id
}

// Broadcast sends data to all subscribers.
func Broadcast(st Type, es ...Event) {
	rtc.Trace("Broadcasting", "st", st, "es", es)

	if st == ToNone && len(es) > 0 {
		id, ok := es[0].Data.(string)
		if !ok {
			return
		}

		subscriptors.mu.Lock()
		defer subscriptors.mu.Unlock()

		for i := range subscriptors.s {
			if subscriptors.s[i].id != id {
				continue
			}
			evt := Event{
				Data: unsubscribeData{st: st, id: id},
			}
			subscriptors.s[i].Event <- evt
			break
		}
		return
	}

	list := es
	if len(list) == 0 {
		list = []Event{{}}
	}

	subscriptors.mu.Lock()
	defer subscriptors.mu.Unlock()
	for i := range subscriptors.s {
		if subscriptors.s[i].st == st {
			go func(i int) {
				for _, x := range list {
					subscriptors.s[i].Event <- x
				}
			}(i)
		}
	}
}

// MustUnsubscribe checks if the event means unsubscribing.
func (s *Subscriber) MustUnsubscribe(e Event) bool {
	rtc.Trace("Checking if must unsubscribe", "e", e)

	ud, ok := e.Data.(unsubscribeData)
	if !ok {
		return false
	}
	return ud.st == ToNone && ud.id == s.id
}

func removeSubscription(id string) {
	slog.Info("Removing subscription", "id", id)

	if unloading {
		return
	}

	subscriptors.mu.Lock()
	defer subscriptors.mu.Unlock()

	for k := range subscriptors.s {
		if subscriptors.s[k].id == id {
			subscriptors.s[k] = subscriptors.s[len(subscriptors.s)-1]
			subscriptors.s = subscriptors.s[:len(subscriptors.s)-1]
			break
		}
	}
}

func unloadSubscriptions() error {
	unloading = true

	subscriptors.mu.Lock()
	defer subscriptors.mu.Unlock()

	for i := range subscriptors.s {
		evt := Event{
			Data: unsubscribeData{st: ToNone, id: subscriptors.s[i].id},
		}
		subscriptors.s[i].Event <- evt
	}
	return nil
}
