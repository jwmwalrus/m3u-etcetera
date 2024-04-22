package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/jwmwalrus/gear-pieces/idler"
	"github.com/jwmwalrus/gear-pieces/middleware"
	"github.com/jwmwalrus/m3u-etcetera/api"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/m3u-etcetera/internal/database"
	"github.com/jwmwalrus/m3u-etcetera/internal/playback"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	rtc "github.com/jwmwalrus/rtcycler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	rtc.Load(&base.Conf,
		base.AppDirName,
		rtc.WithAppName(base.AppName),
		rtc.WithDataSubdirs([]string{base.CoversDirname}),
	)

	idler.Start()

	rtc.RegisterUnloader(database.Open())

	rtc.RegisterUnloader(playback.StartEngine())

	slog.Info("Starting server...")

	port := base.Conf.Server.Port
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		rtc.Fatal("Failed to listen", "error", err)
	}
	slog.Info("Listening on port", "port", port)

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

	rtc.RegisterUnloader(&rtc.Unloader{
		Description: "StopServer",
		Callback: func() error {
			s.Stop()
			listener.Close()
			return nil
		},
	})

	go func() {
		if err := s.Serve(listener); err != nil {
			rtc.Fatal("Failed to serve", "error", err)
		}
	}()

	rtc.RegisterUnloader(subscription.Unloader)

	<-idler.InterruptSignal

	rtc.Unload()

	fmt.Printf("\nBye %v from %v\n", rtc.OS(), rtc.AppInstance())
}
