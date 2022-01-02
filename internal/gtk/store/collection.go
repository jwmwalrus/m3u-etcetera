package store

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
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
	collectionsModel *gtk.TreeStore
	currentTree      CollectionsTree

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

// GetCollectionsModel returns the current collections model
func GetCollectionsModel() *gtk.TreeStore {
	return collectionsModel
}

func subscribeToCollectionStore() {
	log.Info("Subscribing to collection store")

	defer wgcollections.Done()

	var wgdone bool
	var cc *grpc.ClientConn
	var err error
	opts := alive.GetGrpcDialOpts()
	auth := base.Conf.Server.GetAuthority()
	if cc, err = grpc.Dial(auth, opts...); err != nil {
		return
	}
	defer cc.Close()

	cl := m3uetcpb.NewCollectionSvcClient(cc)
	stream, err := cl.SubscribeToCollectionStore(context.Background(), &m3uetcpb.Empty{})
	if err != nil {
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			break
		}

		CStore.Mu.Lock()

		if CStore.subscriptionID == "" {
			CStore.subscriptionID = res.SubscriptionId
		}

		if res.Event == m3uetcpb.CollectionEvent_CE_INITIAL {
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
		} else if res.Event == m3uetcpb.CollectionEvent_CE_INITIAL_DONE {
		} else {
			switch res.Item.(type) {
			case *m3uetcpb.SubscribeToCollectionStoreResponse_CollectionTrack:
				ct := res.GetCollectionTrack()
				switch res.Event {
				case m3uetcpb.CollectionEvent_CE_ITEM_ADDED:
					CStore.CollectionTrack = append(CStore.CollectionTrack, ct)
				case m3uetcpb.CollectionEvent_CE_ITEM_CHANGED:
					for i := range CStore.CollectionTrack {
						if CStore.CollectionTrack[i].Id == ct.Id {
							CStore.CollectionTrack[i] = ct
							break
						}
					}
				case m3uetcpb.CollectionEvent_CE_ITEM_REMOVED:
					n := len(CStore.CollectionTrack)
					for i := range CStore.CollectionTrack {
						if CStore.CollectionTrack[i].Id == ct.Id {
							CStore.CollectionTrack[i] = CStore.CollectionTrack[n-1]
							CStore.CollectionTrack = CStore.CollectionTrack[:n-1]
							break
						}
					}
				}
			case *m3uetcpb.SubscribeToCollectionStoreResponse_Collection:
				c := res.GetCollection()
				switch res.Event {
				case m3uetcpb.CollectionEvent_CE_ITEM_ADDED:
					CStore.Collection = append(CStore.Collection, c)
				case m3uetcpb.CollectionEvent_CE_ITEM_CHANGED:
					for i := range CStore.Collection {
						if CStore.Collection[i].Id == c.Id {
							CStore.Collection[i] = c
							break
						}
					}
				case m3uetcpb.CollectionEvent_CE_ITEM_REMOVED:
					n := len(CStore.Collection)
					for i := range CStore.Collection {
						if CStore.Collection[i].Id == c.Id {
							CStore.Collection[i] = CStore.Collection[n-1]
							CStore.Collection = CStore.Collection[:n-1]
							break
						}
					}
				}
			case *m3uetcpb.SubscribeToCollectionStoreResponse_Track:
				t := res.GetTrack()
				switch res.Event {
				case m3uetcpb.CollectionEvent_CE_ITEM_ADDED:
					CStore.Track = append(CStore.Track, t)
				case m3uetcpb.CollectionEvent_CE_ITEM_CHANGED:
					for i := range CStore.Track {
						if CStore.Track[i].Id == t.Id {
							CStore.Track[i] = t
							break
						}
					}
				case m3uetcpb.CollectionEvent_CE_ITEM_REMOVED:
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
		}
		CStore.Mu.Unlock()
		if res.Event != m3uetcpb.CollectionEvent_CE_INITIAL {
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
	opts := alive.GetGrpcDialOpts()
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
		id     int64
		title  string
		number int
	}

	type yearAlbumType struct {
		yearAlbum string
		ids       []int64
		title     []titleType
	}

	type artistType struct {
		artist    string
		ids       []int64
		yearAlbum []yearAlbumType
	}

	if model.GetNColumns() == 0 {
		return
	}

	iter, ok := model.GetIterFirst()
	for ok {
		model.Remove(iter)
		ok = model.IterNext(iter)
	}

	artistM := map[string]int{}
	all := []artistType{}

	CStore.Mu.Lock()
	for _, t := range CStore.Track {
		artist := t.Albumartist
		if artist == "" {
			artist = t.Artist
		}
		aidx, ok := artistM[artist]
		if !ok {
			all = append(all, artistType{})
			aidx = len(all) - 1
			all[aidx].artist = artist
			artistM[artist] = aidx
		}

		yearAlbum := strconv.Itoa(int(t.Year)) + " - " + t.Album
		yaidx := -1
		for i, ya := range all[aidx].yearAlbum {
			if yearAlbum == ya.yearAlbum {
				yaidx = i
				break
			}
		}
		if yaidx < 0 {
			all[aidx].yearAlbum = append(
				all[aidx].yearAlbum,
				yearAlbumType{},
			)
			yaidx = len(all[aidx].yearAlbum) - 1
			all[aidx].yearAlbum[yaidx].yearAlbum = yearAlbum
		}

		all[aidx].yearAlbum[yaidx].title = append(
			all[aidx].yearAlbum[yaidx].title,
			titleType{t.Id, t.Title, int(t.Tracknumber)},
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

	idListToString := func(ids []int64) (s string) {
		if len(ids) < 1 {
			return
		}
		s = strconv.FormatInt(ids[0], 10)
		for i := 1; i < len(ids); i++ {
			s += "," + strconv.FormatInt(ids[i], 10)
		}
		return
	}

	for i := range all {
		artIter := model.Append(nil)
		model.SetValue(
			artIter,
			CColTree,
			all[i].artist+" ("+strconv.Itoa(len(all[i].ids))+")",
		)
		model.SetValue(
			artIter,
			CColIDList,
			idListToString(all[i].ids),
		)
		for _, ya := range all[i].yearAlbum {
			yaIter := model.Append(artIter)
			model.SetValue(
				yaIter,
				CColTree,
				ya.yearAlbum+" ("+strconv.Itoa(len(ya.ids))+")",
			)
			model.SetValue(
				yaIter,
				CColIDList,
				idListToString(ya.ids),
			)
			for _, t := range ya.title {
				tIter := model.Append(yaIter)
				model.SetValue(tIter, CColTree, t.title)
				model.SetValue(tIter, CColIDList, strconv.FormatInt(t.id, 10))
			}
		}
	}
}
