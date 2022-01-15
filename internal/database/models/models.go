package models

import (
	"math/rand"
	"time"

	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// CollectionEvent defines a collection event
type CollectionEvent int

// CollectionEvent enum
const (
	CollectionEventNone CollectionEvent = iota
	CollectionEventInitial
	_
	_
	CollectionEventItemAdded
	CollectionEventItemChanged
	CollectionEventItemRemoved
	CollectionEventScanning
	CollectionEventScanningDone
)

func (ce CollectionEvent) String() string {
	return []string{
		"none",
		"initial",
		"initial-item",
		"initial-done",
		"item-added",
		"item-changed",
		"item-removed",
		"scanning",
		"scanning-done",
	}[ce]
}

var (
	db *gorm.DB

	// make sure we only do heavy collection tasks one at a time
	storageGuard          chan struct{}
	globalCollectionEvent = CollectionEventNone

	// PlaybackChanged is the AfterCreate-hook channel for QueueTrack and Playback
	PlaybackChanged chan struct{}
)

// DataCreator defines a DML interface of CRUD to insert
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

// DataDeleterTx defines a DML interface of CRUD for writing, with transactions
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

// SetConnection sets the database connection for the whole package
func SetConnection(conn *gorm.DB) {
	db = conn
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

func init() {
	storageGuard = make(chan struct{}, 1)
	PlaybackChanged = make(chan struct{}, 1)
}
