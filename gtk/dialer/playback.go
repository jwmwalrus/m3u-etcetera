package dialer

import (
	"context"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

// ExecutePlaybackAction -
func ExecutePlaybackAction(req *m3uetcpb.ExecutePlaybackActionRequest) (err error) {
	cc, err := getClientConn1()
	if err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.ExecutePlaybackAction(context.Background(), req)
	return
}

// OnProgressBarClicked is the signal handler for the button-press-event on
// the event-box that wraps the progress bar
func OnProgressBarClicked(eb *gtk.EventBox, event *gdk.Event) {
	_, _, duration, status := store.PbData.GetCurrentPlayback()

	if !status["is-streaming"] {
		return
	}

	btn := gdk.EventButtonNewFromEvent(event)
	x, _ := btn.MotionVal()
	width := eb.Widget.GetAllocatedWidth()
	seek := int64(x * float64(duration) / float64(width))

	go func() {
		req := &m3uetcpb.ExecutePlaybackActionRequest{
			Action: m3uetcpb.PlaybackAction_PB_SEEK,
			Seek:   seek,
		}
		onerror.Log(ExecutePlaybackAction(req))
	}()
}

func subscribeToPlayback() {
	log.Info("Subscribing to playback")

	defer wgplayback.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	stream, err := cl.SubscribeToPlayback(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Errorf("Error subscribing to playback: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		store.PbData.ProcessSubscriptionResponse(res)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromPlayback() {
	log.Info("Unsubscribing from playback")

	id := store.PbData.GetSubscriptionID()

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewPlaybackSvcClient(cc)
	_, err = cl.UnsubscribeFromPlayback(
		context.Background(),
		&m3uetcpb.UnsubscribeFromPlaybackRequest{
			SubscriptionId: id,
		},
	)
	onerror.Log(err)
}
