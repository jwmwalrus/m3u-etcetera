package store

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	musicModel      *gtk.ListStore
	radioModel      *gtk.ListStore
	podcastsModel   *gtk.ListStore
	audiobooksModel *gtk.ListStore

	// QStore queue store
	QStore struct {
		Mu   sync.Mutex
		Data *m3uetcpb.SubscribeToQueueStoreResponse
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
		musicModel = model
	case m3uetcpb.Perspective_RADIO:
		radioModel = model
	case m3uetcpb.Perspective_PODCASTS:
		podcastsModel = model
	case m3uetcpb.Perspective_AUDIOBOOKS:
		audiobooksModel = model
	}
	return
}

// GetQueueModel returns the queue model for the given perspective
func GetQueueModel(idx m3uetcpb.Perspective) *gtk.ListStore {
	log.WithField("idx", idx).
		Info("Getting queue model")

	switch idx {
	case m3uetcpb.Perspective_MUSIC:
		return musicModel
	case m3uetcpb.Perspective_RADIO:
		return musicModel
	case m3uetcpb.Perspective_PODCASTS:
		return musicModel
	case m3uetcpb.Perspective_AUDIOBOOKS:
		return musicModel
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

		QStore.Mu.Lock()
		QStore.Data = res
		QStore.Mu.Unlock()
		glib.IdleAdd(updateQueueModels)
		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromQueueStore() {
	log.Info("Unsuubscribing from queue store")

	QStore.Mu.Lock()
	id := QStore.Data.SubscriptionId
	QStore.Mu.Unlock()

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
	if err != nil {
		log.Error(err)
		return
	}
}

func updateQueueModels() bool {
	log.Info("Updating all queue models")

	for _, idx := range perspectivesList {
		model := GetQueueModel(idx)
		if model == nil {
			continue
		}

		iter, ok := model.GetIterFirst()
		for ok {
			model.Remove(iter)
			ok = model.IterNext(iter)
		}
	}

	QStore.Mu.Lock()
	if QStore.Data != nil {
		for _, idx := range perspectivesList {
			model := GetQueueModel(idx)
			if model == nil {
				continue
			}

			if model.GetNColumns() == 0 {
				continue
			}
			// model.Clear()
			var iter *gtk.TreeIter
			for _, qt := range QStore.Data.QueueTracks {
				if qt.Perspective != idx {
					continue
				}
				iter = model.Append()
				err := model.Set(
					iter,
					[]int{
						QColQueueTrackID,
						QColPosition,
						QColPlayed,
						QColLocation,
						QColPerspective,
						QColTrackID,
					},
					[]interface{}{
						qt.Id,
						int(qt.Position),
						qt.Played,
						qt.Location,
						int(qt.Perspective),
						qt.TrackId,
					},
				)
				if err != nil {
					log.Error(err)
					os.Exit(0)
					continue
				}
				for _, t := range QStore.Data.Tracks {
					if qt.TrackId == t.Id {
						err = model.Set(
							iter,
							[]int{
								QColFormat,
								QColType,
								QColTitle,
								QColAlbum,
								QColArtist,
								QColAlbumartist,
								QColComposer,
								QColGenre,

								QColYear,
								QColTracknumber,
								QColTracktotal,
								QColDiscnumber,
								QColDisctotal,
								QColLyrics,
								QColComment,
								QColPlaycount,

								QColRating,
								QColDuration,
								QColRemote,
								QColLastplayed,
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
	QStore.Mu.Unlock()
	return false
}
