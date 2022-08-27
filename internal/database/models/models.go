package models

import (
	"context"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

var (
	db *gorm.DB

	playbackTrackNeeded = make(chan int64, 1)

	queueTrackNeeded = make(chan int64, 1)

	// make sure we only do heavy collection tasks one at a time
	storageGuard chan struct{} = make(chan struct{}, 1)

	// PlaybackChanged is the AfterCreate-hook channel for QueueTrack and Playback
	PlaybackChanged chan struct{} = make(chan struct{}, 1)
)

// DataCreator defines a DML interface of CRUD to create
type DataCreator interface {
	Create() error
}

// DataCreatorTx defines a DML interface of CRUD to create, with transactions
type DataCreatorTx interface {
	CreateTx(*gorm.DB) error
}

// DataDeleter defines a DML interface of CRUD to delete
type DataDeleter interface {
	Delete() error
}

// DataDeleterTx defines a DML interface of CRUD to delete, with transactions
type DataDeleterTx interface {
	DeleteTx(*gorm.DB) error
}

// DataReader defines a DML interface of CRUD for reading
type DataReader interface {
	Read(int64) error
}

// DataReaderTx defines a DML interface of CRUD for reading, with transactions
type DataReaderTx interface {
	ReadTx(*gorm.DB, int64) error
}

// DataUpdater defines a DML interface of CRUD to update
type DataUpdater interface {
	Save() error
}

// DataUpdaterTx defines a DML interface of CRUD to update, with transactions
type DataUpdaterTx interface {
	SaveTx(*gorm.DB) error
}

// ProtoIn defines an interface to convert from protocol buffers
type ProtoIn interface {
	FromProtobuf(proto.Message)
}

// ProtoOut defines an interface to convert to protocol buffers
type ProtoOut interface {
	ToProtobuf() proto.Message
}

// SetUp sets the database used by the models and starts some listeners
func SetUp(ctx context.Context, conn *gorm.DB) {
	db = conn

	doInitialCleanup()

	go findPlaybackTrack(ctx)
	go findQueueTrack(ctx)
}

// TearDown unsets the models listeners
func TearDown() {
	close(playbackTrackNeeded)
	close(queueTrackNeeded)
}

// TriggerPlaybackChange signals the PlaybackChanged channel
func TriggerPlaybackChange() {
	if len(PlaybackChanged) < 1 {
		PlaybackChanged <- struct{}{}
	}
}

func doInitialCleanup() {
	tx := db.Session(&gorm.Session{SkipHooks: true})

	// Clean playback
	tx.Where("played = 1").Delete(&Playback{})

	// Clean queue
	tx.Where("played = 1").Delete(&Queue{})

	// Clean playlists
	pls := []Playlist{}
	tx.Where("open = 0 and transient = 1").Find(&pls)
	for _, pl := range pls {
		err := tx.Where("playlist_id = ?", pl.ID).Delete(&PlaylistTrack{}).Error
		if err != nil {
			log.Error(err)
			break
		}
		tx.Where("id = ?", pl.ID).Delete(&Playlist{})
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
