package models

import (
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

var (
	db *gorm.DB

	// make sure we only do heavy collection tasks one at a time
	storageGuard chan struct{}

	// DbgChan is a debug channel
	DbgChan chan map[string]interface{}

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

// DataFinder defines an interface to query the DB
type DataFinder interface {
	FindBy(query interface{}) (err error)
}

// SetConnection sets the database connection for the whole package
func SetConnection(conn *gorm.DB) {
	db = conn
}

func init() {
	DbgChan = make(chan map[string]interface{})
	storageGuard = make(chan struct{}, 1)
	PlaybackChanged = make(chan struct{}, 1)
}
