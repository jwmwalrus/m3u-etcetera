package store

import (
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/util"
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
	path := model.Path(pliter)
	te.path = path.String()

	model.SetValue(
		pliter,
		int(PLColTree),
		glib.NewValue(te.label+suffix),
	)
	model.SetValue(pliter, int(PLColTreeIDList), util.IDListToGValue(te.ids))
	model.SetValue(pliter, int(PLColTreeIsGroup), glib.NewValue(te.et == playlistGroupEntry))

	for i := range te.child {
		te.child[i].appendNode(model, pliter)
	}
}

func (te *playlistTreeEntry) completeTree(level int,
	guide map[int]playlistEntryType, pl *m3uetcpb.Playlist) {

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

func (te *playlistTreeEntry) fillValues(et playlistEntryType,
	label, kw string, pl *m3uetcpb.Playlist) {
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

	mu sync.RWMutex
}

func (bt *playbarTree) canBeUpdated() bool {
	bt.mu.RLock()
	defer bt.mu.RUnlock()

	return !(bt.initialMode || bt.receivingOpenItems)
}

func (bt *playbarTree) getPlaylistTree(p m3uetcpb.Perspective) playlistTree {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	pt, ok := bt.pplt[p]
	if !ok {
		slog.Warn("There is no playlist tree for perspective", "perspective", p)
	}
	return pt
}

func (bt *playbarTree) setInitialMode(initialMode bool) *playbarTree {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	bt.initialMode = initialMode
	return bt
}

func (bt *playbarTree) setPlaylistTree(p m3uetcpb.Perspective, pt playlistTree) *playbarTree {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	bt.pplt[p] = pt
	return bt
}

func (bt *playbarTree) setReceivingOpenItems(receivingOpenItems bool) *playbarTree {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	bt.receivingOpenItems = receivingOpenItems
	return bt
}

func (bt *playbarTree) update() bool {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	if bt.initialMode || bt.receivingOpenItems {
		return false
	}

	slog.Info("Updating playlist models")

	for _, p := range perspectivesList {
		tree := bt.pplt[p]
		if tree.model == nil {
			continue
		}

		if tree.model.NColumns() == 0 {
			continue
		}
		start := time.Now()

		guide := map[int]playlistEntryType{1: playlistGroupEntry, 2: playlistEntry}

		_, ok := tree.model.IterFirst()
		if ok {
			tree.model.Clear()
		}

		rootIndex := map[string]int{}
		root := []playlistTreeEntry{}

		getKeywords := func(pl *m3uetcpb.Playlist) string {
			list := strings.Split(strings.ToLower(pl.Name), " ")
			list = append(
				list,
				strings.Split(strings.ToLower(pl.Description), " ")...,
			)
			return strings.Join(list, ",")
		}

		BData.mu.RLock()
		for _, pl := range BData.playlist {
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
		BData.mu.RUnlock()

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
		slog.Info("Tree built", "took", time.Since(start))
	}

	return false
}
