package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	musicQueueModel      *gtk.ListStore
	podcastsQueueModel   *gtk.ListStore
	audiobooksQueueModel *gtk.ListStore

	// QData queue store
	QData struct {
		Mu  sync.Mutex
		res *m3uetcpb.SubscribeToQueueStoreResponse
	}
)

// CreateQueueModel -
func CreateQueueModel(idx m3uetcpb.Perspective) (model *gtk.ListStore, err error) {
	log.WithField("idx", idx).
		Info("Creating queue model")

	model, err = gtk.ListStoreNew(QColumns.getTypes()...)
	if err != nil {
		return
	}

	switch idx {
	case m3uetcpb.Perspective_MUSIC:
		musicQueueModel = model
	case m3uetcpb.Perspective_PODCASTS:
		podcastsQueueModel = model
	case m3uetcpb.Perspective_AUDIOBOOKS:
		audiobooksQueueModel = model
	}
	return
}

// ExecuteQueueAction sends an ExecuteQueueAction request
func ExecuteQueueAction(req *m3uetcpb.ExecuteQueueActionRequest) (err error) {
	cc, err := GetClientConn()
	if err != nil {
		return
	}
	defer cc.Close()
	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.ExecuteQueueAction(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		log.Error(s.Message())
		return
	}
	return
}

// GetQueueModel returns the queue model for the given perspective
func GetQueueModel(idx m3uetcpb.Perspective) *gtk.ListStore {
	log.WithField("idx", idx).
		Debug("Returning queue model")

	switch idx {
	case m3uetcpb.Perspective_MUSIC:
		return musicQueueModel
	case m3uetcpb.Perspective_PODCASTS:
		return musicQueueModel
	case m3uetcpb.Perspective_AUDIOBOOKS:
		return musicQueueModel
	default:
		return nil
	}
}

func subscribeToQueueStore() {
	log.Info("Subscribing to queue store")

	defer wgqueue.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	stream, err := cl.SubscribeToQueueStore(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Errorf("Error subscribing to queue store: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		QData.Mu.Lock()
		QData.res = res
		QData.Mu.Unlock()

		glib.IdleAdd(updateQueueModels)

		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromQueueStore() {
	log.Info("Unsuubscribing from queue store")

	QData.Mu.Lock()
	id := QData.res.SubscriptionId
	QData.Mu.Unlock()

	var cc *grpc.ClientConn
	var err error
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	_, err = cl.UnsubscribeFromQueueStore(
		context.Background(),
		&m3uetcpb.UnsubscribeFromQueueStoreRequest{
			SubscriptionId: id,
		},
	)
	onerror.Log(err)
}

func updateQueueModels() bool {
	log.Info("Updating all queue models")

	for _, idx := range perspectiveQueuesList {
		model := GetQueueModel(idx)
		if model == nil {
			continue
		}

		_, ok := model.GetIterFirst()
		if ok {
			model.Clear()
		}
	}

	QData.Mu.Lock()
	if QData.res != nil {
		for _, idx := range perspectiveQueuesList {
			model := GetQueueModel(idx)
			if model == nil {
				continue
			}

			if model.GetNColumns() == 0 {
				continue
			}

			count := 0
			for _, qt := range QData.res.QueueTracks {
				if qt.Perspective != idx {
					continue
				}
				count++
			}

			// model.Clear()
			var iter *gtk.TreeIter
			for _, qt := range QData.res.QueueTracks {
				if qt.Perspective != idx {
					continue
				}
				iter = model.Append()
				err := model.Set(
					iter,
					[]int{
						int(QColQueueTrackID),
						int(QColPosition),
						int(QColLastPosition),
						int(QColPlayed),
						int(QColLocation),
						int(QColPerspective),
						int(QColTrackID),
					},
					[]interface{}{
						qt.Id,
						int(qt.Position),
						int(count),
						qt.Played,
						qt.Location,
						int(qt.Perspective),
						qt.TrackId,
					},
				)
				if err != nil {
					log.Error(err)
					continue
				}
				for _, t := range QData.res.Tracks {
					if qt.TrackId == t.Id {
						err = model.Set(
							iter,
							[]int{
								int(QColFormat),
								int(QColType),
								int(QColTitle),
								int(QColAlbum),
								int(QColArtist),
								int(QColAlbumartist),
								int(QColComposer),
								int(QColGenre),
								int(QColYear),
								int(QColTracknumber),

								int(QColTracktotal),
								int(QColDiscnumber),
								int(QColDisctotal),
								int(QColLyrics),
								int(QColComment),
								int(QColPlaycount),
								int(QColRating),
								int(QColDuration),
								int(QColRemote),
								int(QColLastplayed),
							},
							[]interface{}{
								t.Format,
								t.Type,
								t.Title,
								t.Album,
								t.Artist,
								t.Albumartist,
								t.Composer,
								t.Genre,

								int(t.Year),
								int(t.Tracknumber),
								int(t.Tracktotal),
								int(t.Discnumber),
								int(t.Disctotal),
								t.Lyrics,
								t.Comment,
								int(t.Playcount),

								int(t.Rating),
								fmt.Sprint(time.Duration(t.Duration) * time.Nanosecond),
								t.Remote,
								t.Lastplayed,
							},
						)
					}
				}
			}
		}
	}
	QData.Mu.Unlock()
	return false
}
