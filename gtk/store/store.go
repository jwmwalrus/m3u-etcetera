package store

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

const (
	lastPlayedLayout = "02 Jan 2006 03:04 PM"
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

// ModelColumn defines a tree model column ID.
type ModelColumn int

// CCol* column definition.
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
	CColActionRescan
	CColActionRemove

	CColsN
)

// CColTree* column definition.
const (
	CColTree ModelColumn = iota
	CColTreeIDList
	CColTreeKeywords
)

// TCol* column  definition.
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
	TColTrackNumberOverTotal
	TColDiscnumber
	TColDisctotal
	TColDiscNumberOverTotal
	TColLyrics
	TColComment
	TColPlaycount

	TColRating
	TColDuration
	TColPlayedOverDuration
	TColRemote
	TColLastplayed
	TColPosition
	TColDynamic

	// NOTE: these should never be visible
	TColLastPosition
	TColNumber
	TColToggleSelect
	TColFontWeight

	TColsN
)

// QCol*: queue-track/track column.
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

// QYCol*: query column.
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

// QYColTree*: query tree column.
const (
	QYColTree ModelColumn = iota
	QYColTreeIDList
	QYColTreeKeywords
	QYColTreeSort
)

// PLColTree*: query tree column.
const (
	PLColTree ModelColumn = iota
	PLColTreeIDList
	PLColTreeKeywords
	PLColTreeIsGroup
)

// PGCol* column definition.
const (
	PGColPlaylistGroupID ModelColumn = iota
	PGColName
	PGColDescription
	PGColPerspective
	PGColActionRemove

	PGColsN
)

var (
	// CColumns collection columns.
	CColumns storeColumns

	// TColumns tracks columns.
	TColumns storeColumns

	// QColumns queue columns.
	QColumns storeColumns

	// QYColumns query columns.
	QYColumns storeColumns

	// PGColumns query columns.
	PGColumns storeColumns

	// CTreeColumn collection tree column.
	CTreeColumn storeColumns

	// QYTreeColumn query tree column.
	QYTreeColumn storeColumns

	// PLTreeColumn playlist tree column.
	PLTreeColumn storeColumns
)

func init() {
	CTreeColumn = storeColumns{
		columnDef{Name: "Tree", colType: glib.TypeString},
		columnDef{Name: "ID List", colType: glib.TypeString},
		columnDef{Name: "Keywords", colType: glib.TypeString},
	}

	CColumns = make(storeColumns, CColsN)
	CColumns[CColCollectionID] = columnDef{Name: "ID", colType: glib.TypeInt64}
	CColumns[CColName] = columnDef{Name: "Name", colType: glib.TypeString, editable: true}
	CColumns[CColDescription] = columnDef{Name: "Description", colType: glib.TypeString, editable: true}
	CColumns[CColLocation] = columnDef{Name: "Location", colType: glib.TypeString}
	CColumns[CColRemoteLocation] = columnDef{Name: "Remote Location", colType: glib.TypeString, editable: true}
	CColumns[CColPerspective] = columnDef{Name: "Perspective", colType: glib.TypeString}
	CColumns[CColDisabled] = columnDef{Name: "Disabled", colType: glib.TypeBoolean, activatable: true}
	CColumns[CColRemote] = columnDef{Name: "Remote", colType: glib.TypeBoolean, activatable: true}
	CColumns[CColScanned] = columnDef{Name: "Scanned", colType: glib.TypeInt, activatable: true}
	CColumns[CColTracks] = columnDef{Name: "# Tracks", colType: glib.TypeInt64}

	CColumns[CColTracksView] = columnDef{Name: "# Tracks", colType: glib.TypeString}

	CColumns[CColActionRescan] = columnDef{Name: "ACTION: Re-scan", colType: glib.TypeBoolean, activatable: true}
	CColumns[CColActionRemove] = columnDef{Name: "ACTION: Remove", colType: glib.TypeBoolean, activatable: true}

	TColumns = make(storeColumns, TColsN)
	TColumns[TColTrackID] = columnDef{Name: "ID", colType: glib.TypeInt64}
	TColumns[TColCollectionID] = columnDef{Name: "Collection ID", colType: glib.TypeInt64}
	TColumns[TColLocation] = columnDef{Name: "Location", colType: glib.TypeString}
	TColumns[TColFormat] = columnDef{Name: "Format", colType: glib.TypeString}
	TColumns[TColType] = columnDef{Name: "Type", colType: glib.TypeString}
	TColumns[TColTitle] = columnDef{Name: "Title", colType: glib.TypeString}
	TColumns[TColAlbum] = columnDef{Name: "Album", colType: glib.TypeString}
	TColumns[TColArtist] = columnDef{Name: "Artist", colType: glib.TypeString}
	TColumns[TColAlbumartist] = columnDef{Name: "Album Artist", colType: glib.TypeString}
	TColumns[TColComposer] = columnDef{Name: "Composer", colType: glib.TypeString}
	TColumns[TColGenre] = columnDef{Name: "Genre", colType: glib.TypeString}

	TColumns[TColYear] = columnDef{Name: "Year", colType: glib.TypeInt}
	TColumns[TColTracknumber] = columnDef{Name: "Track Number", colType: glib.TypeInt}
	TColumns[TColTracktotal] = columnDef{Name: "Track Total", colType: glib.TypeInt}
	TColumns[TColTrackNumberOverTotal] = columnDef{Name: "Track # / Total", colType: glib.TypeString}
	TColumns[TColDiscnumber] = columnDef{Name: "Disc Number", colType: glib.TypeInt}
	TColumns[TColDisctotal] = columnDef{Name: "Disc Total", colType: glib.TypeInt}
	TColumns[TColDiscNumberOverTotal] = columnDef{Name: "Disc # / Total", colType: glib.TypeString}
	TColumns[TColLyrics] = columnDef{Name: "Lyrics", colType: glib.TypeString}
	TColumns[TColComment] = columnDef{Name: "Comment", colType: glib.TypeString}
	TColumns[TColPlaycount] = columnDef{Name: "Play Count", colType: glib.TypeInt}

	TColumns[TColRating] = columnDef{Name: "Rating", colType: glib.TypeInt}
	TColumns[TColDuration] = columnDef{Name: "Duration", colType: glib.TypeString}
	TColumns[TColPlayedOverDuration] = columnDef{Name: "(Played / ) Duration", colType: glib.TypeString}
	TColumns[TColRemote] = columnDef{Name: "Remote (T)", colType: glib.TypeBoolean}
	TColumns[TColLastplayed] = columnDef{Name: "Last Played", colType: glib.TypeString}
	TColumns[TColNumber] = columnDef{Name: "#", colType: glib.TypeInt}
	TColumns[TColToggleSelect] = columnDef{Name: "Select", colType: glib.TypeBoolean, activatable: true}
	TColumns[TColPosition] = columnDef{Name: "#", colType: glib.TypeInt}
	TColumns[TColLastPosition] = columnDef{Name: "#", colType: glib.TypeInt}
	TColumns[TColDynamic] = columnDef{Name: "Dynamic", colType: glib.TypeBoolean}
	TColumns[TColFontWeight] = columnDef{Name: "Font weight", colType: glib.TypeInt}

	QColumns = make(storeColumns, QColsN)
	QColumns[QColQueueTrackID] = columnDef{Name: "ID", colType: glib.TypeInt64}
	QColumns[QColPosition] = columnDef{Name: "Position", colType: glib.TypeInt}
	QColumns[QColLastPosition] = columnDef{Name: "Last Position", colType: glib.TypeInt}
	QColumns[QColPlayed] = columnDef{Name: "Played", colType: glib.TypeBoolean}
	QColumns[QColLocation] = columnDef{Name: "Location (QT)", colType: glib.TypeString}
	QColumns[QColPerspective] = columnDef{Name: "Perspective", colType: glib.TypeInt}
	QColumns[QColTrackID] = columnDef{Name: "Track ID", colType: glib.TypeInt64}

	QColumns[QColTrackLocation] = columnDef{Name: "Location", colType: glib.TypeString}
	QColumns[QColFormat] = columnDef{Name: "Format", colType: glib.TypeString}
	QColumns[QColType] = columnDef{Name: "Type", colType: glib.TypeString}
	QColumns[QColTitle] = columnDef{Name: "Title", colType: glib.TypeString}
	QColumns[QColAlbum] = columnDef{Name: "Album", colType: glib.TypeString}
	QColumns[QColArtist] = columnDef{Name: "Artist", colType: glib.TypeString}
	QColumns[QColAlbumartist] = columnDef{Name: "Album Artist", colType: glib.TypeString}
	QColumns[QColComposer] = columnDef{Name: "Composer", colType: glib.TypeString}
	QColumns[QColGenre] = columnDef{Name: "Genre", colType: glib.TypeString}

	QColumns[QColYear] = columnDef{Name: "Year", colType: glib.TypeInt}
	QColumns[QColTracknumber] = columnDef{Name: "Track Number", colType: glib.TypeInt}
	QColumns[QColTracktotal] = columnDef{Name: "Track Total", colType: glib.TypeInt}
	QColumns[QColDiscnumber] = columnDef{Name: "Disc Number", colType: glib.TypeInt}
	QColumns[QColDisctotal] = columnDef{Name: "Disc Total", colType: glib.TypeInt}
	QColumns[QColLyrics] = columnDef{Name: "Lyrics", colType: glib.TypeString}
	QColumns[QColComment] = columnDef{Name: "Comment", colType: glib.TypeString}
	QColumns[QColPlaycount] = columnDef{Name: "Play Count", colType: glib.TypeInt}

	QColumns[QColRating] = columnDef{Name: "Rating", colType: glib.TypeInt}
	QColumns[QColDuration] = columnDef{Name: "Duration", colType: glib.TypeString}
	QColumns[QColRemote] = columnDef{Name: "Remote (T)", colType: glib.TypeBoolean}
	QColumns[QColLastplayed] = columnDef{Name: "Last Played", colType: glib.TypeString}

	QYTreeColumn = storeColumns{
		columnDef{Name: "Tree", colType: glib.TypeString},
		columnDef{Name: "ID List", colType: glib.TypeString},
		columnDef{Name: "Keywords", colType: glib.TypeString},
		columnDef{Name: "Sort", colType: glib.TypeInt},
	}

	// NOTE: Will I ever use this?.
	QYColumns = make(storeColumns, QYColsN)
	QYColumns[QYColQueryID] = columnDef{Name: "ID", colType: glib.TypeInt64}
	QYColumns[QYColName] = columnDef{Name: "Name", colType: glib.TypeString}
	QYColumns[QYColDescription] = columnDef{Name: "Description", colType: glib.TypeString}
	QYColumns[QYColRandom] = columnDef{Name: "Random", colType: glib.TypeBoolean}
	QYColumns[QYColRating] = columnDef{Name: "Rating", colType: glib.TypeInt}
	QYColumns[QYColLimit] = columnDef{Name: "Limit", colType: glib.TypeInt64}
	QYColumns[QYColParams] = columnDef{Name: "Params", colType: glib.TypeString}
	QYColumns[QYColFrom] = columnDef{Name: "From", colType: glib.TypeInt64}
	QYColumns[QYColTo] = columnDef{Name: "To", colType: glib.TypeInt64}
	QYColumns[QYColCollectionIDs] = columnDef{Name: "Collection IDs", colType: glib.TypeInt64}

	PLTreeColumn = storeColumns{
		columnDef{Name: "Tree", colType: glib.TypeString},
		columnDef{Name: "ID List", colType: glib.TypeString},
		columnDef{Name: "Keywords", colType: glib.TypeString},
		columnDef{Name: "Is Group", colType: glib.TypeBoolean},
	}

	PGColumns = make(storeColumns, PGColsN)
	PGColumns[PGColPlaylistGroupID] = columnDef{Name: "ID", colType: glib.TypeInt64}
	PGColumns[PGColName] = columnDef{Name: "Name", colType: glib.TypeString, editable: true}
	PGColumns[PGColDescription] = columnDef{Name: "Description", colType: glib.TypeString, editable: true}
	PGColumns[PGColPerspective] = columnDef{Name: "Perspective", colType: glib.TypeString}
	PGColumns[PGColActionRemove] = columnDef{Name: "ACTION: Remove", colType: glib.TypeBoolean, activatable: true}

}
