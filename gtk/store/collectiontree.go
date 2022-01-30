package store

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	log "github.com/sirupsen/logrus"
)

// collectionTree defines the collections-tree type
type collectionTreeHierarchy int

const (
	// ArtistYearAlbumTree - Artist > Year - Album > Title
	ArtistYearAlbumTree collectionTreeHierarchy = iota

	// ArtistAlbumTree - Artist > Album > Title
	ArtistAlbumTree

	// AlbumTree - Album > Title
	AlbumTree

	// GenreArtistAlbumTree - Genre > Artist > Album > Title
	GenreArtistAlbumTree

	// YearArtistAlbumTree - Year > Artist > Album > Title
	YearArtistAlbumTree
)

func (h collectionTreeHierarchy) getGuide(groupByColl bool) map[int]collectionEntryType {
	switch h {
	case ArtistYearAlbumTree:
		return map[int]collectionEntryType{1: artistEntry, 2: yearAlbumEntry, 3: titleEntry}
	case ArtistAlbumTree:
		return map[int]collectionEntryType{1: artistEntry, 2: albumEntry, 3: titleEntry}
	case AlbumTree:
		return map[int]collectionEntryType{1: albumEntry, 2: titleEntry}
	case GenreArtistAlbumTree:
		return map[int]collectionEntryType{1: genreEntry, 2: artistEntry, 3: albumEntry, 4: titleEntry}
	case YearArtistAlbumTree:
		return map[int]collectionEntryType{1: yearEntry, 2: artistEntry, 3: albumEntry, 4: titleEntry}
	}
	return nil
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

func (te *collectionTreeEntry) appendNode(model *gtk.TreeStore, iter *gtk.TreeIter) {
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
	model.SetValue(citer, int(CColTreeIDList), IDListToString(te.ids))

	for i := range te.child {
		te.child[i].appendNode(model, citer)
	}
}

func (te *collectionTreeEntry) completeTree(level int, guide map[int]collectionEntryType, t *m3uetcpb.Track) {
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

func (te *collectionTreeEntry) fillValues(et collectionEntryType, label, kw string, t *m3uetcpb.Track) {
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
	// rootIndex         map[string]int
	// root              []treeEntry
}

func (tree *collectionTree) rebuild() {
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

	CData.Mu.Lock()
	for _, t := range CData.Track {
		if tree.filterVal != "" {
			kw := getKeywords(t)
			match := true
			for _, s := range strings.Split(strings.ToLower(tree.filterVal), " ") {
				match = match && strings.Contains(kw, s)
			}
			if !match {
				continue
			}
		}

		kw := getKeywords(t)

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
	CData.Mu.Unlock()

	sort.Slice(root, func(i, j int) bool {
		return root[i].label < root[j].label
	})

	for i := range root {
		root[i].sort()
	}

	for i := range root {
		root[i].appendNode(tree.model, nil)
	}
	log.Info("Tree built in ", time.Now().Sub(start))
}

func (tree *collectionTree) update() bool {
	if tree.lastEvent == m3uetcpb.CollectionEvent_CE_ITEM_CHANGED {
		return false
	}
	tree.rebuild()
	return false
}
