package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/jwmwalrus/m3u-etcetera/api"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	base.Load()
	database.Open()
	playback.StartEngine()

	log.Info("Starting server...")

	port := base.Conf.Server.Port
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Infof("Listening on port %v", port)

	opts := middleware.GetServerOpts()
	s := grpc.NewServer(opts...)

	m3uetcpb.RegisterRootSvcServer(s, &api.RootSvc{})
	m3uetcpb.RegisterPlaybackSvcServer(s, &api.PlaybackSvc{})
	m3uetcpb.RegisterQueueSvcServer(s, &api.QueueSvc{})
	m3uetcpb.RegisterCollectionSvcServer(s, &api.CollectionSvc{})
	m3uetcpb.RegisterQuerySvcServer(s, &api.QuerySvc{})

	reflection.Register(s)

	base.RegisterUnloader(base.Unloader{
		Description: "StopServer",
		Callback: func() error {
			s.Stop()
			lsnr.Close()
			return nil
		},
	})

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-base.InterruptSignal

	base.Unload()

	fmt.Printf("\nBye %v from %v\n", base.OS, filepath.Base(os.Args[0]))
}
