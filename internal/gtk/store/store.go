package store

import (
	"sync"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	log "github.com/sirupsen/logrus"
)

var (
	wg               sync.WaitGroup
	perspectivesList []m3uetcpb.Perspective
)

// Subscribe start subscriptions to the server
func Subscribe() {
	log.Info("Subscribing")
	wg.Add(3)
	go subscribeToPlayback()
	go subscribeToQueueStore()
	go subscribeToCollectionStore()
	wg.Wait()
	log.Info("Done subscribing")
}

// Unsubscribe finish all subscriptions to the server
func Unsubscribe() {
	log.Info("Unubscribing")
	unsubscribeFromPlayback()
	unsubscribeFromQueueStore()
	unsubscribeFromCollectionStore()
	log.Info("Done unsubscribing")
}

func init() {
	perspectivesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_RADIO,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}

}
