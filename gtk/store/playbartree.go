package store

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	log "github.com/sirupsen/logrus"
)

type playlistEntryType int

const (
	playlistEntry playlistEntryType = iota
	playlistGroupEntry
)

func (et playlistEntryType) getLabel(pl *m3uetcpb.Playlist) string {
	switch et {
	case playlistGroupEntry:
		pg := playlistToPlaylistGroup[pl.Id]
		if pg.Id > 0 {
			return pg.Name
		}
		return "."
	case playlistEntry:
		return pl.Name
	}
	return ""
}

type playlistTreeEntry struct {
	et                    playlistEntryType
	label, keywords, path string
	ids                   []int64
	index                 map[string]int
	child                 []playlistTreeEntry
}

func (te *playlistTreeEntry) appendNode(model *gtk.TreeStore, iter *gtk.TreeIter) {
	te.ids = te.getIDs()

	suffix := ""
	if te.et != playlistEntry {
		suffix = " (" + strconv.Itoa(len(te.ids)) + ")"
	}

	pliter := model.Append(iter)
	path, _ := model.GetPath(pliter)
	te.path = path.String()

	model.SetValue(
		pliter,
		int(PLColTree),
		te.label+suffix,
	)
	model.SetValue(pliter, int(PLColTreeIDList), IDListToString(te.ids))
	model.SetValue(pliter, int(PLColTreeIsGroup), te.et == playlistGroupEntry)

	for i := range te.child {
		te.child[i].appendNode(model, pliter)
	}
}

func (te *playlistTreeEntry) completeTree(level int, guide map[int]playlistEntryType, pl *m3uetcpb.Playlist) {
	label := guide[level].getLabel(pl)
	idx, ok := te.index[label]
	if !ok {
		te.child = append(te.child, playlistTreeEntry{})
		idx = len(te.child) - 1
		te.child[idx].fillValues(guide[level], label, te.keywords, pl)
		te.index[label] = idx
	}

	if level < len(guide) {
		te.child[idx].completeTree(level+1, guide, pl)
	}
}

func (te *playlistTreeEntry) fillValues(et playlistEntryType, label, kw string, pl *m3uetcpb.Playlist) {
	*te = playlistTreeEntry{
		et:       et,
		label:    label,
		keywords: kw,
		index:    map[string]int{},
	}
	if et == playlistEntry {
		te.ids = []int64{pl.Id}
	}
}

func (te *playlistTreeEntry) getIDs() (ids []int64) {
	if len(te.child) > 0 {
		for i := range te.child {
			ids = append(ids, te.child[i].getIDs()...)
		}
		return
	}

	ids = te.ids
	return
}

func (te *playlistTreeEntry) sort() {
	if len(te.child) == 0 {
		return
	}

	sort.Slice(te.child, func(i, j int) bool {
		return te.child[i].label < te.child[j].label
	})

	for i := range te.child {
		te.child[i].sort()
	}
}

type playlistTree struct {
	model     *gtk.TreeStore
	filterVal string
}

type playbarTree struct {
	pplt               map[m3uetcpb.Perspective]playlistTree
	initialMode        bool
	receivingOpenItems bool
}

func (bt *playbarTree) update() bool {
	if bt.initialMode || bt.receivingOpenItems {
		return false
	}

	log.Info("Updating playlist models")

	for _, p := range perspectivesList {
		tree := bt.pplt[p]
		if tree.model == nil {
			continue
		}

		if tree.model.GetNColumns() == 0 {
			continue
		}
		start := time.Now()

		guide := map[int]playlistEntryType{1: playlistGroupEntry, 2: playlistEntry}

		_, ok := tree.model.GetIterFirst()
		if ok {
			tree.model.Clear()
		}

		rootIndex := map[string]int{}
		root := []playlistTreeEntry{}

		getKeywords := func(pl *m3uetcpb.Playlist) string {
			list := strings.Split(strings.ToLower(pl.Name), " ")
			list = append(list, strings.Split(strings.ToLower(pl.Description), " ")...)
			return strings.Join(list, ",")
		}

		BData.Mu.Lock()
		for _, pl := range BData.Playlist {
			if pl.Transient {
				continue
			}

			if tree.filterVal != "" {
				kw := getKeywords(pl)
				match := true
				for _, s := range strings.Split(strings.ToLower(tree.filterVal), " ") {
					match = match && strings.Contains(kw, s)
				}
				if !match {
					continue
				}
			}

			kw := getKeywords(pl)

			level := 1
			label := guide[level].getLabel(pl)
			idx, ok := rootIndex[label]
			if !ok {
				root = append(root, playlistTreeEntry{})
				idx = len(root) - 1
				rootIndex[label] = idx
				root[idx].fillValues(guide[level], label, kw, pl)
			}
			root[idx].completeTree(level+1, guide, pl)
		}
		BData.Mu.Unlock()

		sort.Slice(root, func(i, j int) bool {
			return root[i].label < root[j].label
		})

		for i := range root {
			root[i].sort()
		}

		for i := range root {
			if root[i].label == "." {
				continue
			}
			root[i].appendNode(tree.model, nil)
		}

		for i := range root {
			if root[i].label != "." {
				continue
			}
			for j := range root[i].child {
				root[i].child[j].appendNode(tree.model, nil)
			}
		}
		log.Info("Tree built in ", time.Now().Sub(start))
	}
	return false
}
