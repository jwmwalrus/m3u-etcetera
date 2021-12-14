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
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Infof("Listening on port %v", port)

	opts := middleware.GetServerOpts()
	s := grpc.NewServer(opts...)

	m3uetcpb.RegisterRootSvcServer(s, &api.RootSvc{})
	m3uetcpb.RegisterPlaybackSvcServer(s, &api.PlaybackSvc{})
	m3uetcpb.RegisterQueueSvcServer(s, &api.QueueSvc{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
