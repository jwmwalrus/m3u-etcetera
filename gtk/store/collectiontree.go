package store

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/util"
	log "github.com/sirupsen/logrus"
)

// collectionTree defines the collection-tree hierarchy
type collectionTreeHierarchy int

const (
	// ArtistYearAlbumTree - Artist > Year - Album > Title
	ArtistYearAlbumTree collectionTreeHierarchy = iota

	// ArtistAlbumTree := Artist > Album > Title
	ArtistAlbumTree

	// AlbumTree := Album > Title
	AlbumTree

	// GenreArtistAlbumTree := Genre > Artist > Album > Title
	GenreArtistAlbumTree

	// YearArtistAlbumTree := Year > Artist > Album > Title
	YearArtistAlbumTree
)

func (h collectionTreeHierarchy) String() string {
	return []string{
		"artist-year-album",
		"artist-album",
		"album",
		"genre-artist-album",
		"year-artist-album",
	}[h]
}

func (h collectionTreeHierarchy) getGuide(groupByColl bool) map[int]collectionEntryType {
	entries := []collectionEntryType{}
	if groupByColl {
		entries = append(entries, collectionEntry)
	}
	switch h {
	case ArtistYearAlbumTree:
		entries = append(entries,
			artistEntry,
			yearAlbumEntry,
			titleEntry,
		)
	case ArtistAlbumTree:
		entries = append(entries,
			artistEntry,
			albumEntry,
			titleEntry,
		)
	case AlbumTree:
		entries = append(entries,
			albumEntry,
			titleEntry,
		)
	case GenreArtistAlbumTree:
		entries = append(entries,
			genreEntry,
			artistEntry,
			albumEntry,
			titleEntry,
		)
	case YearArtistAlbumTree:
		entries = append(entries,
			yearEntry,
			artistEntry,
			albumEntry,
			titleEntry,
		)
	}
	out := make(map[int]collectionEntryType)
	for k, v := range entries {
		out[k+1] = v
	}
	return out
}

type collectionEntryType int

const (
	titleEntry collectionEntryType = iota
	albumEntry
	yearAlbumEntry
	artistEntry
	genreEntry
	yearEntry
	collectionEntry
)

func (et collectionEntryType) getLabel(t *m3uetcpb.Track) string {
	switch et {
	case titleEntry:
		return t.Title
	case albumEntry:
		return t.Album
	case yearAlbumEntry:
		return fmt.Sprintf("%v - %v", t.Year, t.Album)
	case artistEntry:
		artist := t.Albumartist
		if artist == "" {
			artist = t.Artist
		}
		return artist
	case genreEntry:
		return t.Genre
	case yearEntry:
		return fmt.Sprintf("%v", t.Year)
	case collectionEntry:
		return collectionNameMap[t.CollectionId]
	default:
	}
	return ""
}

func (et collectionEntryType) getSorts(t *m3uetcpb.Track) (int, int) {
	if et == titleEntry {
		return int(t.Discnumber), int(t.Tracknumber)
	}
	return 0, 0
}

type collectionTreeEntry struct {
	et                    collectionEntryType
	sort1, sort2          int
	label, keywords, path string
	ids                   []int64
	index                 map[string]int
	child                 []collectionTreeEntry
}

func (te *collectionTreeEntry) appendNode(model *gtk.TreeStore,
	iter *gtk.TreeIter) {

	te.ids = te.getIDs()

	suffix := ""
	if te.et != titleEntry {
		suffix = " (" + strconv.Itoa(len(te.ids)) + ")"
	}

	citer := model.Append(iter)
	path, _ := model.GetPath(citer)
	te.path = path.String()

	model.SetValue(
		citer,
		int(CColTree),
		te.label+suffix,
	)
	model.SetValue(citer, int(CColTreeIDList), util.IDListToString(te.ids))

	for i := range te.child {
		te.child[i].appendNode(model, citer)
	}
}

func (te *collectionTreeEntry) completeTree(level int,
	guide map[int]collectionEntryType, t *m3uetcpb.Track) {

	label := guide[level].getLabel(t)
	idx, ok := te.index[label]
	if !ok {
		te.child = append(te.child, collectionTreeEntry{})
		idx = len(te.child) - 1
		te.child[idx].fillValues(guide[level], label, te.keywords, t)
		te.index[label] = idx
	}

	if level < len(guide) {
		te.child[idx].completeTree(level+1, guide, t)
	}
}

func (te *collectionTreeEntry) fillValues(et collectionEntryType,
	label, kw string, t *m3uetcpb.Track) {
	sort1, sort2 := et.getSorts(t)
	*te = collectionTreeEntry{
		et:       et,
		sort1:    sort1,
		sort2:    sort2,
		label:    label,
		keywords: kw,
		index:    map[string]int{},
	}
	if et == titleEntry {
		te.ids = []int64{t.Id}
	}
}

func (te *collectionTreeEntry) getIDs() (ids []int64) {
	if len(te.child) > 0 {
		for i := range te.child {
			ids = append(ids, te.child[i].getIDs()...)
		}
		return
	}

	ids = te.ids
	return
}

func (te *collectionTreeEntry) sort() {
	if len(te.child) == 0 {
		return
	}

	sort.Slice(te.child, func(i, j int) bool {
		if te.child[i].sort1 > 0 {
			if te.child[i].sort1 != te.child[j].sort1 {
				return te.child[i].sort1 < te.child[j].sort1
			}
			return te.child[i].sort2 < te.child[j].sort2
		}
		return te.child[i].label < te.child[j].label
	})

	for i := range te.child {
		te.child[i].sort()
	}
}

type collectionTree struct {
	model             *gtk.TreeStore
	filterVal         string
	groupByCollection bool
	initialMode       bool
	scanningMode      bool
	lastEvent         m3uetcpb.CollectionEvent
	hierarchy         collectionTreeHierarchy

	mu sync.Mutex
}

func (tree *collectionTree) canBeUpdated() bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	return !(tree.initialMode || tree.scanningMode)
}

func (tree *collectionTree) getLastEvent() m3uetcpb.CollectionEvent {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	return tree.lastEvent
}

func (tree *collectionTree) getModel() *gtk.TreeStore {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	return tree.model
}

func (tree *collectionTree) isInInitialMode() bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	return tree.initialMode
}

func (tree *collectionTree) rebuild() {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	if tree.initialMode || tree.scanningMode {
		return
	}

	if tree.model == nil {
		return
	}

	if tree.model.GetNColumns() == 0 {
		return
	}
	start := time.Now()

	guide := tree.hierarchy.getGuide(tree.groupByCollection)

	_, ok := tree.model.GetIterFirst()
	if ok {
		tree.model.Clear()
	}

	rootIndex := map[string]int{}
	root := []collectionTreeEntry{}

	getKeywords := func(t *m3uetcpb.Track) string {
		list := strings.Split(strings.ToLower(t.Title), " ")
		list = append(list, strings.Split(strings.ToLower(t.Albumartist), " ")...)
		list = append(list, strings.Split(strings.ToLower(t.Album), " ")...)
		list = append(list, strconv.Itoa(int(t.Year)))
		list = append(list, strings.Split(strings.ToLower(t.Artist), " ")...)
		return strings.Join(list, ",")
	}

	CData.updateCollectionNamesMap()

	CData.mu.Lock()
	for _, t := range CData.track {
		kw := getKeywords(t)

		if tree.filterVal != "" {
			match := true
			for _, s := range strings.Split(strings.ToLower(tree.filterVal), " ") {
				match = match && strings.Contains(kw, s)
			}
			if !match {
				continue
			}
		}

		level := 1
		label := guide[level].getLabel(t)
		idx, ok := rootIndex[label]
		if !ok {
			root = append(root, collectionTreeEntry{})
			idx = len(root) - 1
			rootIndex[label] = idx
			root[idx].fillValues(guide[level], label, kw, t)
		}
		root[idx].completeTree(level+1, guide, t)
	}
	CData.mu.Unlock()

	sort.Slice(root, func(i, j int) bool {
		return root[i].label < root[j].label
	})

	for i := range root {
		root[i].sort()
	}

	for i := range root {
		root[i].appendNode(tree.model, nil)
	}
	log.Infof("Tree built in %v", time.Since(start))
}

func (tree *collectionTree) setFilterVal(val string) *collectionTree {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	tree.filterVal = val
	return tree
}

func (tree *collectionTree) setGroupByCollection(groupByCollection bool) *collectionTree {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	tree.groupByCollection = groupByCollection
	return tree
}

func (tree *collectionTree) setHierarchy(hierarchy collectionTreeHierarchy) *collectionTree {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	tree.hierarchy = hierarchy
	return tree
}

func (tree *collectionTree) setInitialMode(initialMode bool) *collectionTree {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	tree.initialMode = initialMode
	return tree
}

func (tree *collectionTree) setLastEvent(lastEvent m3uetcpb.CollectionEvent) *collectionTree {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	tree.lastEvent = lastEvent
	return tree
}

func (tree *collectionTree) setModel(model *gtk.TreeStore) *collectionTree {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	tree.model = model
	return tree
}

func (tree *collectionTree) setScanningMode(scanningMode bool) *collectionTree {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	tree.scanningMode = scanningMode
	return tree
}

func (tree *collectionTree) update() bool {
	if tree.getLastEvent() == m3uetcpb.CollectionEvent_CE_ITEM_CHANGED {
		return false
	}

	tree.rebuild()
	return false
}
