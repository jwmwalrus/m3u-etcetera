package store

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
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

// ModelColumn defines a tree model column ID
type ModelColumn int

// CCol* column definition
const (
	CColCollectionID ModelColumn = iota
	CColName
	CColDescription
	CColLocation
	CColRemoteLocation
	CColPerspective
	CColDisabled
	CColRemote
	CColScanned
	CColTracks
	CColTracksView
	CColRescan

	CColsN
)

// CColTree* column definition
const (
	CColTree ModelColumn = iota
	CColTreeIDList
	CColTreeKeywords
)

// TCol* column  definition
const (
	TColTrackID ModelColumn = iota
	TColCollectionID

	TColLocation
	TColFormat
	TColType
	TColTitle
	TColAlbum
	TColArtist
	TColAlbumartist
	TColComposer
	TColGenre

	TColYear
	TColTracknumber
	TColTracktotal
	TColDiscnumber
	TColDisctotal
	TColLyrics
	TColComment
	TColPlaycount

	TColRating
	TColDuration
	TColRemote
	TColLastplayed
	TColNumber
	TColToggleSelect
	TColPosition
	TColLastPosition
	TColDynamic
	TColFontWeight

	TColsN
)

// QCol*: queue-track/track column
const (
	QColQueueTrackID ModelColumn = iota
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

// QYCol*: query column
const (
	QYColQueryID ModelColumn = iota
	QYColName
	QYColDescription
	QYColRandom
	QYColRating
	QYColLimit
	QYColParams
	QYColFrom
	QYColTo
	QYColCollectionIDs

	QYColsN
)

// QYColTree*: query tree column
const (
	QYColTree ModelColumn = iota
	QYColTreeIDList
	QYColTreeKeywords
)

// PLColTree*: query tree column
const (
	PLColTree ModelColumn = iota
	PLColTreeIDList
	PLColTreeKeywords
	PLColTreeIsGroup
)

// PGCol* column definition
const (
	PGColPlaylistGroupID ModelColumn = iota
	PGColName
	PGColDescription
	PGColPerspective

	PGColsN
)

var (
	wg                    sync.WaitGroup
	wgplayback            sync.WaitGroup
	wgqueue               sync.WaitGroup
	wgcollection          sync.WaitGroup
	wgplaybar             sync.WaitGroup
	wgquery               sync.WaitGroup
	forceExit             bool
	perspectivesList      []m3uetcpb.Perspective
	perspectiveQueuesList []m3uetcpb.Perspective

	// CColumns collection columns
	CColumns storeColumns

	// TColumns tracks columns
	TColumns storeColumns

	// QColumns queue columns
	QColumns storeColumns

	// QYColumns query columns
	QYColumns storeColumns

	// PGColumns query columns
	PGColumns storeColumns

	// CTreeColumn collection tree column
	CTreeColumn storeColumns

	// QYTreeColumn query tree column
	QYTreeColumn storeColumns

	// PLTreeColumn playlist tree column
	PLTreeColumn storeColumns
)

// IDListToString converts an ID list to a comma-separaterd string
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

// SetForceExit sets forceExit to true
func SetForceExit() {
	forceExit = true
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
	onerror.Panic(pbdata.setPlaybackUI())

	onerror.Panic(sanityCheck())

	log.Info("Subscribing")

	wg.Add(5)

	wgplayback.Add(1)
	go subscribeToPlayback()

	wgqueue.Add(1)
	go subscribeToQueueStore()

	wgcollection.Add(1)
	go subscribeToCollectionStore()

	wgplaybar.Add(1)
	go subscribeToPlaybarStore()

	wgquery.Add(1)
	go subscribeToQueryStore()

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
	unsubscribeFromPlaybarStore()
	unsubscribeFromQueryStore()
	wgplayback.Wait()
	wgqueue.Wait()
	wgcollection.Wait()
	wgplaybar.Wait()
	wgquery.Wait()

	alive.Serve(
		alive.ServeOptions{
			TurnOff: true,
			NoWait:  !forceExit,
			Force:   forceExit,
		},
	)

	log.Info("Done unsubscribing")
}

func getClientConn() (*grpc.ClientConn, error) {
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	return grpc.Dial(auth, opts...)
}

func getClientConn1() (*grpc.ClientConn, error) {
	if err := sanityCheck(); err != nil {
		return nil, err
	}
	return getClientConn()
}

func sanityCheck() (err error) {
	log.Info("Checking server status")
	err = alive.CheckServerStatus()
	fmt.Println("sanityCheck:", err)
	switch err.(type) {
	case *alive.ServerAlreadyRunning,
		*alive.ServerStarted:
		log.Info(err)
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

	perspectiveQueuesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}

	barTree.pplt = map[m3uetcpb.Perspective]playlistTree{
		m3uetcpb.Perspective_MUSIC:      {},
		m3uetcpb.Perspective_RADIO:      {},
		m3uetcpb.Perspective_PODCASTS:   {},
		m3uetcpb.Perspective_AUDIOBOOKS: {},
	}

	CTreeColumn = storeColumns{
		columnDef{"Tree", glib.TYPE_STRING},
		columnDef{"ID List", glib.TYPE_STRING},
		columnDef{"Keywords", glib.TYPE_STRING},
	}

	CColumns = make(storeColumns, CColsN)
	CColumns[CColCollectionID] = columnDef{"ID", glib.TYPE_INT64}
	CColumns[CColName] = columnDef{"Name", glib.TYPE_STRING}
	CColumns[CColDescription] = columnDef{"Description", glib.TYPE_STRING}
	CColumns[CColLocation] = columnDef{"Location", glib.TYPE_STRING}
	CColumns[CColRemoteLocation] = columnDef{"Remote Location", glib.TYPE_STRING}
	CColumns[CColPerspective] = columnDef{"Perspective", glib.TYPE_STRING}
	CColumns[CColDisabled] = columnDef{"Disabled", glib.TYPE_BOOLEAN}
	CColumns[CColRemote] = columnDef{"Remote", glib.TYPE_BOOLEAN}
	CColumns[CColScanned] = columnDef{"Scanned", glib.TYPE_INT}
	CColumns[CColTracks] = columnDef{"# Tracks", glib.TYPE_INT64}

	CColumns[CColTracksView] = columnDef{"# Tracks", glib.TYPE_STRING}
	CColumns[CColRescan] = columnDef{"Re-scan", glib.TYPE_BOOLEAN}

	TColumns = make(storeColumns, TColsN)
	TColumns[TColTrackID] = columnDef{"ID", glib.TYPE_INT64}
	TColumns[TColCollectionID] = columnDef{"Collection ID", glib.TYPE_INT64}
	TColumns[TColLocation] = columnDef{"Location", glib.TYPE_STRING}
	TColumns[TColFormat] = columnDef{"Format", glib.TYPE_STRING}
	TColumns[TColType] = columnDef{"Type", glib.TYPE_STRING}
	TColumns[TColTitle] = columnDef{"Title", glib.TYPE_STRING}
	TColumns[TColAlbum] = columnDef{"Album", glib.TYPE_STRING}
	TColumns[TColArtist] = columnDef{"Artist", glib.TYPE_STRING}
	TColumns[TColAlbumartist] = columnDef{"Album Artist", glib.TYPE_STRING}
	TColumns[TColComposer] = columnDef{"Composer", glib.TYPE_STRING}
	TColumns[TColGenre] = columnDef{"Genre", glib.TYPE_STRING}

	TColumns[TColYear] = columnDef{"Year", glib.TYPE_INT}
	TColumns[TColTracknumber] = columnDef{"Track Number", glib.TYPE_INT}
	TColumns[TColTracktotal] = columnDef{"Track Total", glib.TYPE_INT}
	TColumns[TColDiscnumber] = columnDef{"Disc Number", glib.TYPE_INT}
	TColumns[TColDisctotal] = columnDef{"Disc Total", glib.TYPE_INT}
	TColumns[TColLyrics] = columnDef{"Lyrics", glib.TYPE_STRING}
	TColumns[TColComment] = columnDef{"Comment", glib.TYPE_STRING}
	TColumns[TColPlaycount] = columnDef{"Play Count", glib.TYPE_INT}

	TColumns[TColRating] = columnDef{"Rating", glib.TYPE_INT}
	TColumns[TColDuration] = columnDef{"Duration", glib.TYPE_STRING}
	TColumns[TColRemote] = columnDef{"Remote (T)", glib.TYPE_BOOLEAN}
	TColumns[TColLastplayed] = columnDef{"Last Played", glib.TYPE_INT64}
	TColumns[TColNumber] = columnDef{"#", glib.TYPE_INT}
	TColumns[TColToggleSelect] = columnDef{"Select", glib.TYPE_BOOLEAN}
	TColumns[TColPosition] = columnDef{"#", glib.TYPE_INT}
	TColumns[TColLastPosition] = columnDef{"#", glib.TYPE_INT}
	TColumns[TColDynamic] = columnDef{"Dynamic", glib.TYPE_BOOLEAN}
	TColumns[TColFontWeight] = columnDef{"Font weight", glib.TYPE_INT}

	QColumns = make(storeColumns, QColsN)
	QColumns[QColQueueTrackID] = columnDef{"ID", glib.TYPE_INT64}
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

	QYTreeColumn = storeColumns{
		columnDef{"Tree", glib.TYPE_STRING},
		columnDef{"ID List", glib.TYPE_STRING},
		columnDef{"Keywords", glib.TYPE_STRING},
	}

	// NOTE: Will I ever use this?
	QYColumns = make(storeColumns, QYColsN)
	QYColumns[QYColQueryID] = columnDef{"ID", glib.TYPE_INT64}
	QYColumns[QYColName] = columnDef{"Name", glib.TYPE_STRING}
	QYColumns[QYColDescription] = columnDef{"Description", glib.TYPE_STRING}
	QYColumns[QYColRandom] = columnDef{"Random", glib.TYPE_BOOLEAN}
	QYColumns[QYColRating] = columnDef{"Rating", glib.TYPE_INT}
	QYColumns[QYColLimit] = columnDef{"Limit", glib.TYPE_INT64}
	QYColumns[QYColParams] = columnDef{"Params", glib.TYPE_STRING}
	QYColumns[QYColFrom] = columnDef{"From", glib.TYPE_INT64}
	QYColumns[QYColTo] = columnDef{"To", glib.TYPE_INT64}
	QYColumns[QYColCollectionIDs] = columnDef{"Collection IDs", glib.TYPE_INT64}

	PLTreeColumn = storeColumns{
		columnDef{"Tree", glib.TYPE_STRING},
		columnDef{"ID List", glib.TYPE_STRING},
		columnDef{"Keywords", glib.TYPE_STRING},
		columnDef{"Is Group", glib.TYPE_BOOLEAN},
	}

	PGColumns = make(storeColumns, PGColsN)
	PGColumns[PGColPlaylistGroupID] = columnDef{"ID", glib.TYPE_INT64}
	PGColumns[PGColName] = columnDef{"Name", glib.TYPE_STRING}
	PGColumns[PGColDescription] = columnDef{"Description", glib.TYPE_STRING}
	PGColumns[PGColPerspective] = columnDef{"Perspective", glib.TYPE_STRING}

}
