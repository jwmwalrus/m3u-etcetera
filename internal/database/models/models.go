package models

import (
	"context"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

var (
	db *gorm.DB

	playbackTrackNeeded chan int64
	queueTrackNeeded    chan int64

	jsonUnmarshaler = protojson.UnmarshalOptions{DiscardUnknown: true}

	// storageGuard make sure we only do heavy collection tasks one at a time.
	storageGuard chan struct{} = make(chan struct{}, 1)

	// PlaybackChanged is the AfterCreate-hook channel for QueueTrack and Playback.
	PlaybackChanged chan struct{} = make(chan struct{}, 1)
)

// DataCreator defines a DML interface of CRUD to create.
type DataCreator interface {
	Create() error
}

// DataCreatorTx defines a DML interface of CRUD to create, with transactions.
type DataCreatorTx interface {
	CreateTx(*gorm.DB) error
}

// DataDeleter defines a DML interface of CRUD to delete.
type DataDeleter interface {
	Delete() error
}

// DataDeleterTx defines a DML interface of CRUD to delete, with transactions.
type DataDeleterTx interface {
	DeleteTx(*gorm.DB) error
}

// DataReader defines a DML interface of CRUD for reading.
type DataReader interface {
	Read(int64) error
}

// DataReaderTx defines a DML interface of CRUD for reading, with transactions.
type DataReaderTx interface {
	ReadTx(*gorm.DB, int64) error
}

// DataUpdater defines a DML interface of CRUD to update.
type DataUpdater interface {
	Save() error
}

// DataUpdaterTx defines a DML interface of CRUD to update, with transactions.
type DataUpdaterTx interface {
	SaveTx(*gorm.DB) error
}

// ProtoIn defines an interface to convert from protocol buffers.
type ProtoIn interface {
	FromProtobuf(proto.Message)
}

// ProtoOut defines an interface to convert to protocol buffers.
type ProtoOut interface {
	ToProtobuf() proto.Message
}

// DoInitialCleanup perform initial cleanup of models tables.
func DoInitialCleanup() {
	tx := db.Session(&gorm.Session{SkipHooks: true})

	// Clean playback
	tx.Where("played = 1").Delete(&Playback{})

	// Clean queue
	tx.Where("played = 1").Delete(&QueueTrack{})

	// Clean deleted playlists
	pls := []Playlist{}
	tx.Where("open = 0 and transient = 1").Find(&pls)
	for _, pl := range pls {
		err := tx.Where("playlist_id = ?", pl.ID).Delete(&PlaylistTrack{}).Error
		if err != nil {
			log.Error(err)
			continue
		}

		log.Debugf("Removing delete playlist, ID=%d", pl.ID)
		tx.Where("id = ?", pl.ID).Delete(&Playlist{})
	}

	// Activate default perspective
	tx.Model(&Perspective{}).Where("idx = ?", int(DefaultPerspective)).Update("active", true)

	// Set collections' scanned to 100
	tx.Model(&Collection{}).Where("id > 0").Update("scanned", 100)

	// Remove unused transient tracks
	trc, _ := TransientCollection.Get()
	if trc != nil {
		func() {
			var ts []struct{ ID int64 }
			tx.Model(&Track{}).Where("collection_id = ?", trc.ID).Find(&ts)
			if len(ts) == 0 {
				return
			}

			var tset []int64
			for _, t := range ts {
				tset = append(tset, t.ID)
			}

			var pts []struct{ TrackID int64 }
			tx.Model(&PlaylistTrack{}).Where("track_id IN (?)", tset).Find(&pts)

			var ptset []int64
			for _, pt := range pts {
				ptset = append(ptset, pt.TrackID)
			}
			ptset = slices.Compact(ptset)

			var diff []int64
			for _, id := range tset {
				if slices.Contains(ptset, id) {
					continue
				}
				diff = append(diff, id)
			}

			log.Debugf("Number of transient tracks: %d", len(tset))
			log.Debugf("Number of transient tracks in use: %d", len(ptset))
			log.Debugf("Number of unused transient tracks: %d", len(diff))
			if len(diff) == 0 {
				return
			}

			tx.Where("id IN (?)", diff).Delete(&Track{})
		}()
	}
}

// SetUp sets the database used by the models and starts some listeners.
func SetUp(ctx context.Context, conn *gorm.DB) {
	db = conn

	playbackTrackNeeded = make(chan int64, 1)
	queueTrackNeeded = make(chan int64, 1)

	go findPlaybackTrack(ctx)
	go findQueueTrack(ctx)
}

// TearDown unsets the models listeners.
func TearDown() {
	storageGuard <- struct{}{}
	defer func() { <-storageGuard }()

	close(playbackTrackNeeded)
	close(queueTrackNeeded)
}

// TriggerPlaybackChange signals the PlaybackChanged channel.
func TriggerPlaybackChange() {
	if len(PlaybackChanged) < 1 {
		PlaybackChanged <- struct{}{}
	}
}

func getSuffler(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed.Shuffle(n, func(i, j int) { s[i], s[j] = s[j], s[i] })

	return s
}
