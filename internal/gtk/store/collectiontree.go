package store

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
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

	// CollectionArtistYearAlbumTree - Collection > Artist > Year - Album > Title
	CollectionArtistYearAlbumTree

	// CollectionArtistAlbumTree - Collection > Artist > Album > Title
	CollectionArtistAlbumTree

	// CollectionAlbumTree - Collection > Album > Title
	CollectionAlbumTree

	// CollectionGenreArtistAlbumTree - Collection > Genre > Artist > Album > Title
	CollectionGenreArtistAlbumTree

	// CollectionYearArtistAlbumTree - Collection > Year > Artist > Album > Title
	CollectionYearArtistAlbumTree
)

type collectionTree struct {
	model             *gtk.TreeStore
	filterVal         string
	hierarchy         collectionTreeHierarchy
	groupByCollection bool
}

func (h collectionTreeHierarchy) String() string {
	return []string{
		"Artist - (Year - Album)",
		"Collection - Artist - (Year - Album)",
	}[h]
}

func (tree *collectionTree) update() bool {
	switch tree.hierarchy {
	case ArtistYearAlbumTree:
		tree.updateArtistYearAlbumTree()
	default:
		tree.updateArtistYearAlbumTree()
	}
	return false
}

func (tree *collectionTree) updateArtistYearAlbumTree() {
	if tree.model == nil {
		return
	}

	type titleType struct {
		id       int64
		title    string
		number   int
		keywords string
	}

	type yearAlbumType struct {
		yearAlbum string
		ids       []int64
		title     []titleType
	}

	type artistType struct {
		artist     string
		ids        []int64
		yearAlbumM map[string]int
		yearAlbum  []yearAlbumType
	}

	if tree.model.GetNColumns() == 0 {
		return
	}

	_, ok := tree.model.GetIterFirst()
	if ok {
		tree.model.Clear()
	}

	artistM := map[string]int{}
	all := []artistType{}

	getKeywords := func(t *m3uetcpb.Track) string {
		list := strings.Split(t.Title, " ")
		list = append(list, strings.Split(strings.ToLower(t.Albumartist), " ")...)
		list = append(list, strings.Split(strings.ToLower(t.Album), " ")...)
		list = append(list, strconv.Itoa(int(t.Year)))
		list = append(list, strings.Split(strings.ToLower(t.Artist), " ")...)
		return strings.Join(list, ",")
	}

	CStore.Mu.Lock()
	for _, t := range CStore.Track {
		if tree.filterVal != "" {
			kw := getKeywords(t)
			match := false
			for _, s := range strings.Split(tree.filterVal, " ") {
				match = match || strings.Contains(kw, s)
				if match {
					break
				}
			}
			if !match {
				continue
			}
		}
		artist := t.Albumartist
		if artist == "" {
			artist = t.Artist
		}
		aidx, ok := artistM[artist]
		if !ok {
			all = append(all, artistType{})
			aidx = len(all) - 1
			all[aidx].artist = artist
			all[aidx].yearAlbumM = map[string]int{}
			artistM[artist] = aidx
		}

		yearAlbum := strconv.Itoa(int(t.Year)) + " - " + t.Album
		yaidx, ok := all[aidx].yearAlbumM[yearAlbum]
		if !ok {
			all[aidx].yearAlbum = append(
				all[aidx].yearAlbum,
				yearAlbumType{},
			)
			yaidx = len(all[aidx].yearAlbum) - 1
			all[aidx].yearAlbum[yaidx].yearAlbum = yearAlbum
			all[aidx].yearAlbumM[yearAlbum] = yaidx
		}

		all[aidx].yearAlbum[yaidx].title = append(
			all[aidx].yearAlbum[yaidx].title,
			titleType{t.Id, t.Title, int(t.Tracknumber), getKeywords(t)},
		)
	}
	CStore.Mu.Unlock()

	sort.Slice(all, func(i, j int) bool {
		return all[i].artist < all[j].artist
	})

	for a := range all {
		sort.Slice(all[a].yearAlbum, func(i, j int) bool {
			return all[a].yearAlbum[i].yearAlbum <
				all[a].yearAlbum[j].yearAlbum
		})
		for b := range all[a].yearAlbum {
			sort.Slice(all[a].yearAlbum[b].title, func(i, j int) bool {
				return all[a].yearAlbum[b].title[i].number <
					all[a].yearAlbum[b].title[j].number
			})

			for c := range all[a].yearAlbum[b].title {
				all[a].yearAlbum[b].ids = append(
					all[a].yearAlbum[b].ids,
					all[a].yearAlbum[b].title[c].id,
				)
			}

			all[a].ids = append(all[a].ids, all[a].yearAlbum[b].ids...)
		}
	}

	for i := range all {
		artIter := tree.model.Append(nil)
		tree.model.SetValue(
			artIter,
			int(CColTree),
			all[i].artist+" ("+strconv.Itoa(len(all[i].ids))+")",
		)
		tree.model.SetValue(artIter, int(CColTreeIDList), IDListToString(all[i].ids))
		for _, ya := range all[i].yearAlbum {
			yaIter := tree.model.Append(artIter)
			tree.model.SetValue(
				yaIter,
				int(CColTree),
				ya.yearAlbum+" ("+strconv.Itoa(len(ya.ids))+")",
			)
			tree.model.SetValue(yaIter, int(CColTreeIDList), IDListToString(ya.ids))
			for _, t := range ya.title {
				tIter := tree.model.Append(yaIter)
				tree.model.SetValue(tIter, int(CColTree), t.title)
				tree.model.SetValue(tIter, int(CColTreeIDList), strconv.FormatInt(t.id, 10))
			}
		}
	}
}
