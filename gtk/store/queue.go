package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

type queueData struct {
	res *m3uetcpb.SubscribeToQueueStoreResponse

	mu sync.Mutex
}

var (
	// QData queue store
	QData = &queueData{}

	perspectiveQueuesList []m3uetcpb.Perspective

	musicQueueModel *gtk.ListStore
	//nolint: unused //TODO
	podcastsQueueModel *gtk.ListStore
	//nolint: unused //TODO
	audiobooksQueueModel *gtk.ListStore
)

func init() {
	perspectiveQueuesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}
}

func (qd *queueData) GetQueueDigest(idx m3uetcpb.Perspective) *m3uetcpb.PerspectiveDigest {
	qd.mu.Lock()
	defer qd.mu.Unlock()

	for _, dig := range qd.res.Digest {
		if dig.Perspective == idx {
			return dig
		}
	}
	return nil
}

func (qd *queueData) GetQueueTracksCount(idx m3uetcpb.Perspective) int64 {
	qd.mu.Lock()
	defer qd.mu.Unlock()

	var count int64
	for _, qt := range qd.res.QueueTracks {
		if qt.Perspective == idx {
			count++
		}
	}
	return count
}

func (qd *queueData) GetSubscriptionID() string {
	qd.mu.Lock()
	defer qd.mu.Unlock()

	return qd.res.SubscriptionId
}

func (qd *queueData) ProcessSubscriptionResponse(res *m3uetcpb.SubscribeToQueueStoreResponse) {
	qd.mu.Lock()
	defer qd.mu.Unlock()

	qd.res = res

	glib.IdleAdd(qd.updateQueueModels)
}

func (qd *queueData) updateQueueModels() bool {
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

	qd.mu.Lock()
	defer qd.mu.Unlock()

	if qd.res != nil {
		for _, idx := range perspectiveQueuesList {
			model := GetQueueModel(idx)
			if model == nil {
				continue
			}

			if model.GetNColumns() == 0 {
				continue
			}

			count := 0
			for _, qt := range qd.res.QueueTracks {
				if qt.Perspective != idx {
					continue
				}
				count++
			}

			// model.Clear()
			var iter *gtk.TreeIter
			for _, qt := range qd.res.QueueTracks {
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
						count,
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
				for _, t := range qd.res.Tracks {
					if qt.TrackId == t.Id {
						dur := time.Duration(t.Duration) * time.Nanosecond
						err := model.Set(
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
								fmt.Sprint(dur.Truncate(time.Second)),
								t.Remote,
								t.Lastplayed,
							},
						)
						onerror.Log(err)
						break
					}
				}
			}
		}
	}
	return false
}

// CreateQueueModel -
func CreateQueueModel(idx m3uetcpb.Perspective) (
	model *gtk.ListStore, err error) {

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
