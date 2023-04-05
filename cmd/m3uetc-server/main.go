package main

import (
	"fmt"
	"net"

	"github.com/jwmwalrus/m3u-etcetera/api"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	base.Load()
	base.StartIdler()

	base.RegisterUnloader(database.Open())

	base.RegisterUnloader(playback.StartEngine())

	log.Info("Starting server...")

	port := base.Conf.Server.Port
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Infof("Listening on port %v", port)

	opts := middleware.GetServerOpts()
	s := grpc.NewServer(opts...)

	pbEvents := playback.GetEventsInstance()

	m3uetcpb.RegisterRootSvcServer(s, &api.RootSvc{})
	m3uetcpb.RegisterPlaybackSvcServer(s, &api.PlaybackSvc{PbEvents: pbEvents})
	m3uetcpb.RegisterQueueSvcServer(s, &api.QueueSvc{})
	m3uetcpb.RegisterCollectionSvcServer(s, &api.CollectionSvc{})
	m3uetcpb.RegisterQuerySvcServer(s, &api.QuerySvc{})
	m3uetcpb.RegisterPlaybarSvcServer(s, &api.PlaybarSvc{PbEvents: pbEvents})
	m3uetcpb.RegisterPerspectiveSvcServer(s, &api.PerspectiveSvc{})

	reflection.Register(s)

	base.RegisterUnloader(&base.Unloader{
		Description: "StopServer",
		Callback: func() error {
			s.Stop()
			listener.Close()
			return nil
		},
	})

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	base.RegisterUnloader(subscription.Unloader)

	<-base.InterruptSignal

	base.Unload()

	fmt.Printf("\nBye %v from %v\n", base.OS, base.AppInstance)
}
