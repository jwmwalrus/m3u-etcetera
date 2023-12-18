package store

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
)

type queueData struct {
	res *m3uetcpb.SubscribeToQueueStoreResponse

	mu sync.RWMutex
}

var (
	// QData queue store.
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
	qd.mu.RLock()
	defer qd.mu.RUnlock()

	for _, dig := range qd.res.Digest {
		if dig.Perspective == idx {
			return dig
		}
	}
	return nil
}

func (qd *queueData) GetQueueTracksCount(idx m3uetcpb.Perspective) int64 {
	qd.mu.RLock()
	defer qd.mu.RUnlock()

	var count int64
	for _, qt := range qd.res.QueueTracks {
		if qt.Perspective == idx {
			count++
		}
	}
	return count
}

func (qd *queueData) GetSubscriptionID() string {
	qd.mu.RLock()
	defer qd.mu.RUnlock()

	return qd.res.SubscriptionId
}

func (qd *queueData) ProcessSubscriptionResponse(res *m3uetcpb.SubscribeToQueueStoreResponse) {
	qd.mu.Lock()
	defer qd.mu.Unlock()

	qd.res = res

	glib.IdleAdd(qd.updateQueueModels)
}

func (qd *queueData) updateQueueModels() bool {
	slog.Info("Updating all queue models")

	for _, idx := range perspectiveQueuesList {
		model := GetQueueModel(idx)
		if model == nil {
			continue
		}

		_, ok := model.IterFirst()
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

			if model.NColumns() == 0 {
				continue
			}

			count := 0
			for _, qt := range qd.res.QueueTracks {
				if qt.Perspective != idx {
					continue
				}
				count++
			}

			var iter *gtk.TreeIter
			for _, qt := range qd.res.QueueTracks {
				if qt.Perspective != idx {
					continue
				}
				iter = model.Append()
				model.Set(
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
					[]glib.Value{
						*glib.NewValue(qt.Id),
						*glib.NewValue(int(qt.Position)),
						*glib.NewValue(count),
						*glib.NewValue(qt.Played),
						*glib.NewValue(qt.Location),
						*glib.NewValue(int(qt.Perspective)),
						*glib.NewValue(qt.TrackId),
					},
				)
				for _, t := range qd.res.Tracks {
					if qt.Location == t.Location {
						dur := time.Duration(t.Duration) * time.Nanosecond
						model.Set(
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
							[]glib.Value{
								*glib.NewValue(t.Format),
								*glib.NewValue(t.Type),
								*glib.NewValue(t.Title),
								*glib.NewValue(t.Album),
								*glib.NewValue(t.Artist),
								*glib.NewValue(t.Albumartist),
								*glib.NewValue(t.Composer),
								*glib.NewValue(t.Genre),
								*glib.NewValue(int(t.Year)),
								*glib.NewValue(int(t.Tracknumber)),

								*glib.NewValue(int(t.Tracktotal)),
								*glib.NewValue(int(t.Discnumber)),
								*glib.NewValue(int(t.Disctotal)),
								*glib.NewValue(t.Lyrics),
								*glib.NewValue(t.Comment),
								*glib.NewValue(int(t.Playcount)),
								*glib.NewValue(int(t.Rating)),
								*glib.NewValue(fmt.Sprint(dur.Truncate(time.Second))),
								*glib.NewValue(t.Remote),
								*glib.NewValue(time.Unix(0, t.Lastplayed).Format(lastPlayedLayout)),
							},
						)
						break
					}
				}
			}
		}
	}
	return false
}

// CreateQueueModel -.
func CreateQueueModel(idx m3uetcpb.Perspective) (
	model *gtk.ListStore, err error) {

	slog.Info("Creating queue model", "idx", idx)

	model = gtk.NewListStore(QColumns.getTypes())
	if model == nil {
		err = fmt.Errorf("failed to create list-store")
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

// GetQueueModel returns the queue model for the given perspective.
func GetQueueModel(idx m3uetcpb.Perspective) *gtk.ListStore {
	slog.Debug("Returning queue model", "idx", idx)

	switch idx {
	case m3uetcpb.Perspective_MUSIC:
		return musicQueueModel
	case m3uetcpb.Perspective_PODCASTS:
		// TODO: return podcastsQueueModel
		return musicQueueModel
	case m3uetcpb.Perspective_AUDIOBOOKS:
		// TODO: return audiobooksQueueModel
		return musicQueueModel
	default:
		return nil
	}
}
