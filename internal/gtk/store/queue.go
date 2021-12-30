package store

import (
	"context"
	"os"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	// QColID -
	QColID int = iota

	// QColPosition -
	QColPosition

	// QColLocation -
	QColLocation

	// QColPerspective -
	QColPerspective

	// QColTrackID -
	QColTrackID

	// QColTitle -
	QColTitle

	// QColAlbum -
	QColAlbum

	// QColArtist -
	QColArtist

	// QColAlbumartist -
	QColAlbumartist

	// QColComposer -
	QColComposer

	// QColGenre -
	QColGenre

	// QColYear -
	QColYear

	// QColTracknumber -
	QColTracknumber

	// QColTracktotal -
	QColTracktotal

	// QColDiscnumber -
	QColDiscnumber

	// QColDisctotal -
	QColDisctotal

	// QColDuration -
	QColDuration
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

	// QColumns queue columns
	QColumns = []struct {
		Name    string
		ColType glib.Type
	}{
		{"ID", glib.TYPE_INT64},
		{"Position", glib.TYPE_INT},
		{"Location", glib.TYPE_STRING},
		{"Perspective", glib.TYPE_INT},
		{"Track ID", glib.TYPE_INT},
		{"Title", glib.TYPE_STRING},
		{"Album", glib.TYPE_STRING},
		{"Artist", glib.TYPE_STRING},
		{"Album Artist", glib.TYPE_STRING},
		{"Composer", glib.TYPE_STRING},
		{"Genre", glib.TYPE_STRING},
		{"Year", glib.TYPE_INT},
		{"Tracknumber", glib.TYPE_INT},
		{"Tracktotal", glib.TYPE_INT},
		{"Discnumber", glib.TYPE_INT},
		{"Disctotal", glib.TYPE_INT},
		{"Duration", glib.TYPE_INT64},
	}
)

// CreateQueueModel -
func CreateQueueModel(idx m3uetcpb.Perspective) (model *gtk.ListStore, err error) {
	log.WithField("idx", idx).
		Info("Creating queue model")

	gtypes := []glib.Type{}
	for _, v := range QColumns {
		gtypes = append(gtypes, v.ColType)
	}

	model, err = gtk.ListStoreNew(gtypes...)
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

	var wgdone bool
	var cc *grpc.ClientConn
	var err error
	opts := alive.GetGrpcDialOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewQueueSvcClient(cc)
	stream, err := cl.SubscribeToQueueStore(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
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
	opts := alive.GetGrpcDialOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
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
						QColID,
						QColPosition,
						QColLocation,
						QColPerspective,
						QColTrackID,
					},
					[]interface{}{
						qt.Id,
						int(qt.Position),
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
								QColDuration,
							},
							[]interface{}{
								t.Title,
								t.Album,
								t.Artist,
								t.Albumartist,
								t.Composer,
								t.Genre,
								t.Year,
								t.Tracknumber,
								t.Tracktotal,
								t.Discnumber,
								t.Disctotal,
								t.Duration,
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
