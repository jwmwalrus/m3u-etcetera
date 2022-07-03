package dialer

import (
	"sync"

	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	wg            sync.WaitGroup
	wgplayback    sync.WaitGroup
	wgqueue       sync.WaitGroup
	wgcollection  sync.WaitGroup
	wgplaybar     sync.WaitGroup
	wgquery       sync.WaitGroup
	wgperspective sync.WaitGroup
	forceExit     bool
)

// SetForceExit sets forceExit to true
func SetForceExit() {
	forceExit = true
}

// Subscribe start subscriptions to the server
func Subscribe() {
	onerror.Panic(store.PbData.SetPlaybackUI())
	onerror.Panic(store.PerspData.SetPerspectiveUI())

	onerror.Panic(sanityCheck())

	log.Info("Subscribing")

	wg.Add(6)

	wgplayback.Add(1)
	go subscribeToPlayback()

	wgqueue.Add(1)
	go subscribeToQueueStore()

	wgcollection.Add(1)
	go subscribeToCollectionStore()

	wgplaybar.Add(1)
	go subscribeToPlaybarStore()

	wgquery.Add(1)
	go subscribeToQueryStore()

	wgperspective.Add(1)
	go subscribeToPerspective()

	wg.Wait()
	log.Info("Done subscribing")
}

// Unsubscribe finish all subscriptions to the server
func Unsubscribe() {
	sanityCheck()

	log.Info("Unubscribing")

	unsubscribeFromPlayback()
	unsubscribeFromQueueStore()
	unsubscribeFromCollectionStore()
	unsubscribeFromPlaybarStore()
	unsubscribeFromQueryStore()
	unsubscribeFromPerspective()
	wgplayback.Wait()
	wgqueue.Wait()
	wgcollection.Wait()
	wgplaybar.Wait()
	wgquery.Wait()
	wgperspective.Wait()

	alive.Serve(
		alive.ServeOptions{
			TurnOff: true,
			NoWait:  !forceExit,
			Force:   forceExit,
		},
	)

	log.Info("Done unsubscribing")
}

func getClientConn() (*grpc.ClientConn, error) {
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	return grpc.Dial(auth, opts...)
}

func getClientConn1() (*grpc.ClientConn, error) {
	if err := sanityCheck(); err != nil {
		return nil, err
	}
	return getClientConn()
}

func sanityCheck() (err error) {
	log.Info("Checking server status")
	err = alive.CheckServerStatus()
	switch err.(type) {
	case *alive.ServerAlreadyRunning,
		*alive.ServerStarted:
		log.Info(err)
		err = nil
	default:
	}
	return
}
