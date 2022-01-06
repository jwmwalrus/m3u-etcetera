package store

import (
	"strconv"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type storeColumns []columnDef

func (sc storeColumns) getTypes() (s []glib.Type) {
	for _, v := range sc {
		s = append(s, v.ColType)
	}
	return
}

type columnDef struct {
	Name    string
	ColType glib.Type
}

// TreeModelColumn defines a tree model column ID
type TreeModelColumn int

// CColTree column definition
const (
	CColTree TreeModelColumn = iota
	CColTreeIDList
	CColTreeKeywords
)

// CCol* column definition
const (
	CColCollectionID TreeModelColumn = iota
	CColName
	CColDescription
	CColLocation
	CColHidden
	CColDisabled
	CColRemote
	CColScanned
	CColTracks

	CColsN
)

// CTCol* column  definition
const (
	CTColCollectionTrackID TreeModelColumn = iota
	CTColCollectionID
	CTColTrackID

	CTColLocation
	CTColFormat
	CTColType
	CTColTitle
	CTColAlbum
	CTColArtist
	CTColAlbumartist
	CTColComposer
	CTColGenre

	CTColYear
	CTColTracknumber
	CTColTracktotal
	CTColDiscnumber
	CTColDisctotal
	CTColLyrics
	CTColComment
	CTColPlaycount

	CTColRating
	CTColDuration
	CTColRemote
	CTColLastplayed

	CTColsN
)

// QCol*: queue-track/track column
const (
	QColQueueTrackID TreeModelColumn = iota
	QColPosition
	QColLastPosition
	QColPlayed
	QColLocation
	QColPerspective
	QColTrackID

	QColTrackLocation
	QColFormat
	QColType
	QColTitle
	QColAlbum
	QColArtist
	QColAlbumartist
	QColComposer
	QColGenre

	QColYear
	QColTracknumber
	QColTracktotal
	QColDiscnumber
	QColDisctotal
	QColLyrics
	QColComment
	QColPlaycount

	QColRating
	QColDuration
	QColRemote
	QColLastplayed

	QColsN
)

var (
	wg               sync.WaitGroup
	wgplayback       sync.WaitGroup
	wgqueue          sync.WaitGroup
	wgcollections    sync.WaitGroup
	perspectivesList []m3uetcpb.Perspective

	// TreeColumn collections tree column
	TreeColumn storeColumns

	// CColumns collections columns
	CColumns storeColumns

	// CTColumns collection-tracks columns
	CTColumns storeColumns

	// QColumns queue columns
	QColumns storeColumns
)

// GetClientConn returns a valid client connection to the server
func GetClientConn() (*grpc.ClientConn, error) {
	if err := sanityCheck(); err != nil {
		return nil, err
	}
	return getClientConn()
}

func IDListToString(ids []int64) (s string) {
	if len(ids) < 1 {
		return
	}
	s = strconv.FormatInt(ids[0], 10)
	for i := 1; i < len(ids); i++ {
		s += "," + strconv.FormatInt(ids[i], 10)
	}
	return
}

// StringToIDList parses the IDList column
func StringToIDList(s string) (ids []int64, err error) {
	if len(s) == 0 {
		return
	}
	list := strings.Split(s, ",")
	for _, l := range list {
		var id int64
		id, err = strconv.ParseInt(l, 10, 64)
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}

// Subscribe start subscriptions to the server
func Subscribe() {
	err := sanityCheck()
	onerror.Panic(err)

	log.Info("Subscribing")

	wg.Add(3)

	wgplayback.Add(1)
	go subscribeToPlayback()

	wgqueue.Add(1)
	go subscribeToQueueStore()

	wgcollections.Add(1)
	go subscribeToCollectionStore()

	wg.Wait()
	log.Info("Done subscribing")
}

// Unsubscribe finish all subscriptions to the server
func Unsubscribe() {
	sanityCheck()

	log.Info("Unubscribing")

	unsubscribeFromPlayback()
	unsubscribeFromQueueStore()
	unsubscribeFromCollectionStore()
	wgplayback.Wait()
	wgqueue.Wait()
	wgcollections.Wait()

	alive.Serve(alive.ServeOptions{TurnOff: true, NoWait: true})

	log.Info("Done unsubscribing")
}

func getClientConn() (*grpc.ClientConn, error) {
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	return grpc.Dial(auth, opts...)
}

func sanityCheck() (err error) {
	log.Info("Checking server status")
	err = alive.CheckServerStatus()
	switch err.(type) {
	case *alive.ServerAlreadyRunning,
		*alive.ServerStarted:
		err = nil
	default:
	}
	return
}

func init() {
	perspectivesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_RADIO,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}

	TreeColumn = storeColumns{
		columnDef{"Tree", glib.TYPE_STRING},
		columnDef{"ID List", glib.TYPE_STRING},
		columnDef{"Keywords", glib.TYPE_STRING},
	}

	CColumns = make(storeColumns, CColsN)
	CColumns[CColCollectionID] = columnDef{"Collection ID", glib.TYPE_INT64}
	CColumns[CColName] = columnDef{"Name", glib.TYPE_STRING}
	CColumns[CColDescription] = columnDef{"Description", glib.TYPE_STRING}
	CColumns[CColLocation] = columnDef{"Location", glib.TYPE_STRING}
	CColumns[CColHidden] = columnDef{"Hidden", glib.TYPE_BOOLEAN}
	CColumns[CColDisabled] = columnDef{"Disabled", glib.TYPE_BOOLEAN}
	CColumns[CColRemote] = columnDef{"Remote", glib.TYPE_BOOLEAN}
	CColumns[CColScanned] = columnDef{"Scanned", glib.TYPE_INT}
	CColumns[CColTracks] = columnDef{"# Tracks", glib.TYPE_INT}

	// FIXME: Why doesn't make work here?
	// CTColumns = make(storeColumns, CTColsN)
	for i := 0; i < int(CTColsN); i++ {
		CColumns = append(CColumns, columnDef{})
	}
	CColumns[CTColCollectionTrackID] = columnDef{"CollectionTrack ID", glib.TYPE_INT64}
	CColumns[CTColCollectionID] = columnDef{"Collection ID", glib.TYPE_INT64}
	CColumns[CTColTrackID] = columnDef{"Track ID", glib.TYPE_INT64}

	CColumns[CTColLocation] = columnDef{"Location", glib.TYPE_STRING}
	CColumns[CTColFormat] = columnDef{"Format", glib.TYPE_STRING}
	CColumns[CTColType] = columnDef{"Type", glib.TYPE_STRING}
	CColumns[CTColTitle] = columnDef{"Title", glib.TYPE_STRING}
	CColumns[CTColAlbum] = columnDef{"Album", glib.TYPE_STRING}
	CColumns[CTColArtist] = columnDef{"Artist", glib.TYPE_STRING}
	CColumns[CTColAlbumartist] = columnDef{"Album Artist", glib.TYPE_STRING}
	CColumns[CTColComposer] = columnDef{"Composer", glib.TYPE_STRING}
	CColumns[CTColGenre] = columnDef{"Genre", glib.TYPE_STRING}

	CColumns[CTColYear] = columnDef{"Year", glib.TYPE_INT}
	CColumns[CTColTracknumber] = columnDef{"Track Number", glib.TYPE_INT}
	CColumns[CTColTracktotal] = columnDef{"Track Total", glib.TYPE_INT}
	CColumns[CTColDiscnumber] = columnDef{"Disc Number", glib.TYPE_INT}
	CColumns[CTColDisctotal] = columnDef{"Disc Total", glib.TYPE_INT}
	CColumns[CTColLyrics] = columnDef{"Lyrics", glib.TYPE_STRING}
	CColumns[CTColComment] = columnDef{"Comment", glib.TYPE_STRING}
	CColumns[CTColPlaycount] = columnDef{"Play Count", glib.TYPE_INT}

	CColumns[CTColRating] = columnDef{"Rating", glib.TYPE_INT}
	CColumns[CTColDuration] = columnDef{"Duration", glib.TYPE_INT64}
	CColumns[CTColRemote] = columnDef{"Remote (T)", glib.TYPE_BOOLEAN}
	CColumns[CTColLastplayed] = columnDef{"Last Played", glib.TYPE_INT64}

	QColumns = make(storeColumns, QColsN)
	QColumns[QColQueueTrackID] = columnDef{"QueueTrack ID", glib.TYPE_INT64}
	QColumns[QColPosition] = columnDef{"Position", glib.TYPE_INT}
	QColumns[QColLastPosition] = columnDef{"Last Position", glib.TYPE_INT}
	QColumns[QColPlayed] = columnDef{"Played", glib.TYPE_BOOLEAN}
	QColumns[QColLocation] = columnDef{"Location (QT)", glib.TYPE_STRING}
	QColumns[QColPerspective] = columnDef{"Perspective", glib.TYPE_INT}
	QColumns[QColTrackID] = columnDef{"Track ID", glib.TYPE_INT64}

	QColumns[QColTrackLocation] = columnDef{"Location", glib.TYPE_STRING}
	QColumns[QColFormat] = columnDef{"Format", glib.TYPE_STRING}
	QColumns[QColType] = columnDef{"Type", glib.TYPE_STRING}
	QColumns[QColTitle] = columnDef{"Title", glib.TYPE_STRING}
	QColumns[QColAlbum] = columnDef{"Album", glib.TYPE_STRING}
	QColumns[QColArtist] = columnDef{"Artist", glib.TYPE_STRING}
	QColumns[QColAlbumartist] = columnDef{"Album Artist", glib.TYPE_STRING}
	QColumns[QColComposer] = columnDef{"Composer", glib.TYPE_STRING}
	QColumns[QColGenre] = columnDef{"Genre", glib.TYPE_STRING}

	QColumns[QColYear] = columnDef{"Year", glib.TYPE_INT}
	QColumns[QColTracknumber] = columnDef{"Track Number", glib.TYPE_INT}
	QColumns[QColTracktotal] = columnDef{"Track Total", glib.TYPE_INT}
	QColumns[QColDiscnumber] = columnDef{"Disc Number", glib.TYPE_INT}
	QColumns[QColDisctotal] = columnDef{"Disc Total", glib.TYPE_INT}
	QColumns[QColLyrics] = columnDef{"Lyrics", glib.TYPE_STRING}
	QColumns[QColComment] = columnDef{"Comment", glib.TYPE_STRING}
	QColumns[QColPlaycount] = columnDef{"Play Count", glib.TYPE_INT}

	QColumns[QColRating] = columnDef{"Rating", glib.TYPE_INT}
	QColumns[QColDuration] = columnDef{"Duration", glib.TYPE_STRING}
	QColumns[QColRemote] = columnDef{"Remote (T)", glib.TYPE_BOOLEAN}
	QColumns[QColLastplayed] = columnDef{"Last Played", glib.TYPE_INT64}
}
