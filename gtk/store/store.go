package store

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
)

type storeColumns []columnDef

func (sc storeColumns) getTypes() (s []glib.Type) {
	for _, v := range sc {
		s = append(s, v.colType)
	}
	return
}

func (sc storeColumns) GetActivatableColumns() (s []ModelColumn) {
	for k, v := range sc {
		if !v.activatable {
			continue
		}
		s = append(s, ModelColumn(k))
	}
	return
}

func (sc storeColumns) GetEditableColumns() (s []ModelColumn) {
	for k, v := range sc {
		if !v.editable {
			continue
		}
		s = append(s, ModelColumn(k))
	}
	return
}

type columnDef struct {
	Name        string
	colType     glib.Type
	editable    bool
	activatable bool
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

	CColsN
)

// CActionCol* column definition
const (
	CActionColCollectionID ModelColumn = iota
	CActionColName
	CActionColRescan
	CActionColRemove

	CActionColsN
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
	perspectivesList      []m3uetcpb.Perspective
	perspectiveQueuesList []m3uetcpb.Perspective

	// CColumns collection columns
	CColumns storeColumns

	// CActionColumns collection actions columns
	CActionColumns storeColumns

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
		columnDef{Name: "Tree", colType: glib.TYPE_STRING},
		columnDef{Name: "ID List", colType: glib.TYPE_STRING},
		columnDef{Name: "Keywords", colType: glib.TYPE_STRING},
	}

	CColumns = make(storeColumns, CColsN)
	CColumns[CColCollectionID] = columnDef{Name: "ID", colType: glib.TYPE_INT64}
	CColumns[CColName] = columnDef{Name: "Name", colType: glib.TYPE_STRING, editable: true}
	CColumns[CColDescription] = columnDef{Name: "Description", colType: glib.TYPE_STRING, editable: true}
	CColumns[CColLocation] = columnDef{Name: "Location", colType: glib.TYPE_STRING}
	CColumns[CColRemoteLocation] = columnDef{Name: "Remote Location", colType: glib.TYPE_STRING, editable: true}
	CColumns[CColPerspective] = columnDef{Name: "Perspective", colType: glib.TYPE_STRING}
	CColumns[CColDisabled] = columnDef{Name: "Disabled", colType: glib.TYPE_BOOLEAN, activatable: true}
	CColumns[CColRemote] = columnDef{Name: "Remote", colType: glib.TYPE_BOOLEAN, activatable: true}
	CColumns[CColScanned] = columnDef{Name: "Scanned", colType: glib.TYPE_INT, activatable: true}
	CColumns[CColTracks] = columnDef{Name: "# Tracks", colType: glib.TYPE_INT64}

	CColumns[CColTracksView] = columnDef{Name: "# Tracks", colType: glib.TYPE_STRING}

	CActionColumns = make(storeColumns, CActionColsN)
	CActionColumns[CActionColCollectionID] = columnDef{Name: "ID", colType: glib.TYPE_INT64}
	CActionColumns[CActionColName] = columnDef{Name: "Name", colType: glib.TYPE_STRING}
	CActionColumns[CActionColRescan] = columnDef{Name: "Re-scan", colType: glib.TYPE_BOOLEAN, activatable: true}
	CActionColumns[CActionColRemove] = columnDef{Name: "Remove", colType: glib.TYPE_BOOLEAN, activatable: true}

	TColumns = make(storeColumns, TColsN)
	TColumns[TColTrackID] = columnDef{Name: "ID", colType: glib.TYPE_INT64}
	TColumns[TColCollectionID] = columnDef{Name: "Collection ID", colType: glib.TYPE_INT64}
	TColumns[TColLocation] = columnDef{Name: "Location", colType: glib.TYPE_STRING}
	TColumns[TColFormat] = columnDef{Name: "Format", colType: glib.TYPE_STRING}
	TColumns[TColType] = columnDef{Name: "Type", colType: glib.TYPE_STRING}
	TColumns[TColTitle] = columnDef{Name: "Title", colType: glib.TYPE_STRING}
	TColumns[TColAlbum] = columnDef{Name: "Album", colType: glib.TYPE_STRING}
	TColumns[TColArtist] = columnDef{Name: "Artist", colType: glib.TYPE_STRING}
	TColumns[TColAlbumartist] = columnDef{Name: "Album Artist", colType: glib.TYPE_STRING}
	TColumns[TColComposer] = columnDef{Name: "Composer", colType: glib.TYPE_STRING}
	TColumns[TColGenre] = columnDef{Name: "Genre", colType: glib.TYPE_STRING}

	TColumns[TColYear] = columnDef{Name: "Year", colType: glib.TYPE_INT}
	TColumns[TColTracknumber] = columnDef{Name: "Track Number", colType: glib.TYPE_INT}
	TColumns[TColTracktotal] = columnDef{Name: "Track Total", colType: glib.TYPE_INT}
	TColumns[TColDiscnumber] = columnDef{Name: "Disc Number", colType: glib.TYPE_INT}
	TColumns[TColDisctotal] = columnDef{Name: "Disc Total", colType: glib.TYPE_INT}
	TColumns[TColLyrics] = columnDef{Name: "Lyrics", colType: glib.TYPE_STRING}
	TColumns[TColComment] = columnDef{Name: "Comment", colType: glib.TYPE_STRING}
	TColumns[TColPlaycount] = columnDef{Name: "Play Count", colType: glib.TYPE_INT}

	TColumns[TColRating] = columnDef{Name: "Rating", colType: glib.TYPE_INT}
	TColumns[TColDuration] = columnDef{Name: "Duration", colType: glib.TYPE_STRING}
	TColumns[TColRemote] = columnDef{Name: "Remote (T)", colType: glib.TYPE_BOOLEAN}
	TColumns[TColLastplayed] = columnDef{Name: "Last Played", colType: glib.TYPE_INT64}
	TColumns[TColNumber] = columnDef{Name: "#", colType: glib.TYPE_INT}
	TColumns[TColToggleSelect] = columnDef{Name: "Select", colType: glib.TYPE_BOOLEAN, activatable: true}
	TColumns[TColPosition] = columnDef{Name: "#", colType: glib.TYPE_INT}
	TColumns[TColLastPosition] = columnDef{Name: "#", colType: glib.TYPE_INT}
	TColumns[TColDynamic] = columnDef{Name: "Dynamic", colType: glib.TYPE_BOOLEAN}
	TColumns[TColFontWeight] = columnDef{Name: "Font weight", colType: glib.TYPE_INT}

	QColumns = make(storeColumns, QColsN)
	QColumns[QColQueueTrackID] = columnDef{Name: "ID", colType: glib.TYPE_INT64}
	QColumns[QColPosition] = columnDef{Name: "Position", colType: glib.TYPE_INT}
	QColumns[QColLastPosition] = columnDef{Name: "Last Position", colType: glib.TYPE_INT}
	QColumns[QColPlayed] = columnDef{Name: "Played", colType: glib.TYPE_BOOLEAN}
	QColumns[QColLocation] = columnDef{Name: "Location (QT)", colType: glib.TYPE_STRING}
	QColumns[QColPerspective] = columnDef{Name: "Perspective", colType: glib.TYPE_INT}
	QColumns[QColTrackID] = columnDef{Name: "Track ID", colType: glib.TYPE_INT64}

	QColumns[QColTrackLocation] = columnDef{Name: "Location", colType: glib.TYPE_STRING}
	QColumns[QColFormat] = columnDef{Name: "Format", colType: glib.TYPE_STRING}
	QColumns[QColType] = columnDef{Name: "Type", colType: glib.TYPE_STRING}
	QColumns[QColTitle] = columnDef{Name: "Title", colType: glib.TYPE_STRING}
	QColumns[QColAlbum] = columnDef{Name: "Album", colType: glib.TYPE_STRING}
	QColumns[QColArtist] = columnDef{Name: "Artist", colType: glib.TYPE_STRING}
	QColumns[QColAlbumartist] = columnDef{Name: "Album Artist", colType: glib.TYPE_STRING}
	QColumns[QColComposer] = columnDef{Name: "Composer", colType: glib.TYPE_STRING}
	QColumns[QColGenre] = columnDef{Name: "Genre", colType: glib.TYPE_STRING}

	QColumns[QColYear] = columnDef{Name: "Year", colType: glib.TYPE_INT}
	QColumns[QColTracknumber] = columnDef{Name: "Track Number", colType: glib.TYPE_INT}
	QColumns[QColTracktotal] = columnDef{Name: "Track Total", colType: glib.TYPE_INT}
	QColumns[QColDiscnumber] = columnDef{Name: "Disc Number", colType: glib.TYPE_INT}
	QColumns[QColDisctotal] = columnDef{Name: "Disc Total", colType: glib.TYPE_INT}
	QColumns[QColLyrics] = columnDef{Name: "Lyrics", colType: glib.TYPE_STRING}
	QColumns[QColComment] = columnDef{Name: "Comment", colType: glib.TYPE_STRING}
	QColumns[QColPlaycount] = columnDef{Name: "Play Count", colType: glib.TYPE_INT}

	QColumns[QColRating] = columnDef{Name: "Rating", colType: glib.TYPE_INT}
	QColumns[QColDuration] = columnDef{Name: "Duration", colType: glib.TYPE_STRING}
	QColumns[QColRemote] = columnDef{Name: "Remote (T)", colType: glib.TYPE_BOOLEAN}
	QColumns[QColLastplayed] = columnDef{Name: "Last Played", colType: glib.TYPE_INT64}

	QYTreeColumn = storeColumns{
		columnDef{Name: "Tree", colType: glib.TYPE_STRING},
		columnDef{Name: "ID List", colType: glib.TYPE_STRING},
		columnDef{Name: "Keywords", colType: glib.TYPE_STRING},
	}

	// NOTE: Will I ever use this?
	QYColumns = make(storeColumns, QYColsN)
	QYColumns[QYColQueryID] = columnDef{Name: "ID", colType: glib.TYPE_INT64}
	QYColumns[QYColName] = columnDef{Name: "Name", colType: glib.TYPE_STRING}
	QYColumns[QYColDescription] = columnDef{Name: "Description", colType: glib.TYPE_STRING}
	QYColumns[QYColRandom] = columnDef{Name: "Random", colType: glib.TYPE_BOOLEAN}
	QYColumns[QYColRating] = columnDef{Name: "Rating", colType: glib.TYPE_INT}
	QYColumns[QYColLimit] = columnDef{Name: "Limit", colType: glib.TYPE_INT64}
	QYColumns[QYColParams] = columnDef{Name: "Params", colType: glib.TYPE_STRING}
	QYColumns[QYColFrom] = columnDef{Name: "From", colType: glib.TYPE_INT64}
	QYColumns[QYColTo] = columnDef{Name: "To", colType: glib.TYPE_INT64}
	QYColumns[QYColCollectionIDs] = columnDef{Name: "Collection IDs", colType: glib.TYPE_INT64}

	PLTreeColumn = storeColumns{
		columnDef{Name: "Tree", colType: glib.TYPE_STRING},
		columnDef{Name: "ID List", colType: glib.TYPE_STRING},
		columnDef{Name: "Keywords", colType: glib.TYPE_STRING},
		columnDef{Name: "Is Group", colType: glib.TYPE_BOOLEAN},
	}

	PGColumns = make(storeColumns, PGColsN)
	PGColumns[PGColPlaylistGroupID] = columnDef{Name: "ID", colType: glib.TYPE_INT64}
	PGColumns[PGColName] = columnDef{Name: "Name", colType: glib.TYPE_STRING, editable: true}
	PGColumns[PGColDescription] = columnDef{Name: "Description", colType: glib.TYPE_STRING, editable: true}
	PGColumns[PGColPerspective] = columnDef{Name: "Perspective", colType: glib.TYPE_STRING}

}
