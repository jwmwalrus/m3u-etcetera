package store

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/api/middleware"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// CollectionsTree defines the collections-tree type
type CollectionsTree int

// Album: Album > Title
// Artist: Artist > Album > Title
// Album Artist: Album Artist > Album > Title
// Genre - Artist: Genre > Artist > Album > Title
// Genre - Album Artist: Genre > Album Artist > Album > Title
// Genre - Album: Genre > Album  > Artist > Title
// Year - Artist: Year > Artist > Album > Title
// Year - Album Artist: Year > Album Artist > Album > Title
// Year - Album: Year > Album  > Artist > Title
//  Artist - (Year - Album): Artist > (Year - Album) > Title
//  Album Artist - (Year - Album): Album Artist > (Year - Album) > Title
const (
	ArtistYearAlbumTree CollectionsTree = iota
	CollectionArtistYearAlbumTree
)

func (tree CollectionsTree) String() string {
	return []string{
		"Artist - (Year - Album)",
		"Collection - Artist - (Year - Album)",
	}[tree]
}

var (
	collectionsModel     *gtk.TreeStore
	collectionsFilter    *gtk.TreeModelFilter
	collectionsFilterVal string

	currentTree CollectionsTree

	// CStore collection store
	CStore struct {
		subscriptionID  string
		Mu              sync.Mutex
		CollectionTrack []*m3uetcpb.CollectionTrack
		Collection      []*m3uetcpb.Collection
		Track           []*m3uetcpb.Track
	}
)

// CreateCollectionsModel creates a collection model
func CreateCollectionsModel(tree CollectionsTree) (model *gtk.TreeStore, err error) {
	log.WithField("tree", tree).
		Info("Creating collections model")

	collectionsModel, err = gtk.TreeStoreNew(TreeColumn.getTypes()...)
	if err != nil {
		return
	}

	currentTree = tree
	model = collectionsModel
	return
}

// FilterCollectionsBy filters the collections by the given value
func FilterCollectionsBy(val string) {
	collectionsFilterVal = val
	updateCollections()
}

// GetCollectionsModel returns the current collections model
func GetCollectionsModel() *gtk.TreeStore {
	return collectionsModel
}

func subscribeToCollectionStore() {
	log.Info("Subscribing to collection store")

	defer wgcollections.Done()

	var wgdone bool

	cc, err := getClientConn()
	if err != nil {
		log.Errorf("Error obtaining client connection: %v", err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewCollectionSvcClient(cc)
	stream, err := cl.SubscribeToCollectionStore(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		log.Errorf("Error subscribing to collection store: %v", err)
		return
	}

	appendItem := func(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
		switch res.Item.(type) {
		case *m3uetcpb.SubscribeToCollectionStoreResponse_CollectionTrack:
			CStore.CollectionTrack = append(
				CStore.CollectionTrack,
				res.GetCollectionTrack(),
			)
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
			CStore.Collection = append(
				CStore.Collection,
				res.GetCollection(),
			)
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
			CStore.Track = append(CStore.Track, res.GetTrack())
		default:
		}
	}

	changeItem := func(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
		switch res.Item.(type) {
		case *m3uetcpb.SubscribeToCollectionStoreResponse_CollectionTrack:
			ct := res.GetCollectionTrack()
			for i := range CStore.CollectionTrack {
				if CStore.CollectionTrack[i].Id == ct.Id {
					CStore.CollectionTrack[i] = ct
					break
				}
			}
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
			c := res.GetCollection()
			for i := range CStore.Collection {
				if CStore.Collection[i].Id == c.Id {
					CStore.Collection[i] = c
					break
				}
			}
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
			t := res.GetTrack()
			for i := range CStore.Track {
				if CStore.Track[i].Id == t.Id {
					CStore.Track[i] = t
					break
				}
			}
		}
	}

	removeItem := func(res *m3uetcpb.SubscribeToCollectionStoreResponse) {
		switch res.Item.(type) {
		case *m3uetcpb.SubscribeToCollectionStoreResponse_CollectionTrack:
			ct := res.GetCollectionTrack()
			n := len(CStore.CollectionTrack)
			for i := range CStore.CollectionTrack {
				if CStore.CollectionTrack[i].Id == ct.Id {
					CStore.CollectionTrack[i] = CStore.CollectionTrack[n-1]
					CStore.CollectionTrack = CStore.CollectionTrack[:n-1]
					break
				}
			}
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
			c := res.GetCollection()
			n := len(CStore.Collection)
			for i := range CStore.Collection {
				if CStore.Collection[i].Id == c.Id {
					CStore.Collection[i] = CStore.Collection[n-1]
					CStore.Collection = CStore.Collection[:n-1]
					break
				}
			}
		case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
			t := res.GetTrack()
			n := len(CStore.Track)
			for i := range CStore.Track {
				if CStore.Track[i].Id == t.Id {
					CStore.Track[i] = CStore.Track[n-1]
					CStore.Track = CStore.Track[:n-1]
					break
				}
			}
		}
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Infof("Subscription closed by server: %v", err)
			break
		}

		CStore.Mu.Lock()

		if CStore.subscriptionID == "" {
			CStore.subscriptionID = res.SubscriptionId
		}

		switch res.Event {
		case m3uetcpb.CollectionEvent_CE_INITIAL:
			CStore.CollectionTrack = []*m3uetcpb.CollectionTrack{}
			CStore.Collection = []*m3uetcpb.Collection{}
			CStore.Track = []*m3uetcpb.Track{}
		case m3uetcpb.CollectionEvent_CE_INITIAL_ITEM:
			appendItem(res)
		case m3uetcpb.CollectionEvent_CE_INITIAL_DONE:
			// pass
		case m3uetcpb.CollectionEvent_CE_ITEM_ADDED:
			appendItem(res)
		case m3uetcpb.CollectionEvent_CE_ITEM_CHANGED:
			changeItem(res)
		case m3uetcpb.CollectionEvent_CE_ITEM_REMOVED:
			removeItem(res)
		}

		CStore.Mu.Unlock()
		if res.Event != m3uetcpb.CollectionEvent_CE_INITIAL &&
			res.Event != m3uetcpb.CollectionEvent_CE_INITIAL_ITEM {
			glib.IdleAdd(updateCollections)
		}
		if !wgdone {
			wg.Done()
			wgdone = true
		}
	}
}

func unsubscribeFromCollectionStore() {
	log.Info("Unsubscribing from collection store")

	CStore.Mu.Lock()
	id := CStore.subscriptionID
	CStore.Mu.Unlock()

	var cc *grpc.ClientConn
	var err error
	opts := middleware.GetClientOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		log.Error(err)
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewCollectionSvcClient(cc)
	_, err = cl.UnsubscribeFromCollectionStore(
		context.Background(),
		&m3uetcpb.UnsubscribeFromCollectionStoreRequest{
			SubscriptionId: id,
		},
	)
	if err != nil {
		log.Error(err)
		return
	}
}

func updateCollections() bool {
	switch currentTree {
	case ArtistYearAlbumTree:
		updateArtistYearAlbumTree()
	default:
	}
	return false
}

func updateArtistYearAlbumTree() {
	model := collectionsModel
	if model == nil {
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

	if model.GetNColumns() == 0 {
		return
	}

	_, ok := model.GetIterFirst()
	if ok {
		model.Clear()
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
		if collectionsFilterVal != "" {
			kw := getKeywords(t)
			match := false
			for _, s := range strings.Split(collectionsFilterVal, " ") {
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
		artIter := model.Append(nil)
		model.SetValue(
			artIter,
			CColTree,
			all[i].artist+" ("+strconv.Itoa(len(all[i].ids))+")",
		)
		model.SetValue(artIter, CColTreeIDList, idListToString(all[i].ids))
		for _, ya := range all[i].yearAlbum {
			yaIter := model.Append(artIter)
			model.SetValue(
				yaIter,
				CColTree,
				ya.yearAlbum+" ("+strconv.Itoa(len(ya.ids))+")",
			)
			model.SetValue(yaIter, CColTreeIDList, idListToString(ya.ids))
			for _, t := range ya.title {
				tIter := model.Append(yaIter)
				model.SetValue(tIter, CColTree, t.title)
				model.SetValue(tIter, CColTreeIDList, strconv.FormatInt(t.id, 10))
			}
		}
	}
}
