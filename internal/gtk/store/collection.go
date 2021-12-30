package store

import (
	"context"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/alive"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	titleCol int = iota
	albumCol
	artistCol
	albumArtistCol
	genreCol
	yearCol
)

type collectionData struct {
	subscriptionID  string
	Mu              sync.Mutex
	CollectionTrack []*m3uetcpb.CollectionTrack
	Collection      []*m3uetcpb.Collection
	Track           []*m3uetcpb.Track
}

var (
	// CStore collection store
	CStore collectionData
)

// Album -> Album > Title
// Artist -> Artist > Album > Title
// Album Artist -> Album Artist > Album > Title

// Genre - Artist -> Genre > Artist > Album > Title
// Genre - Album Artist -> Genre > Album Artist > Album > Title
// Genre - Album -> Genre > Album  > Artist > Title

// Year - Artist -> Year > Artist > Album > Title
// Year - Album Artist -> Year > Album Artist > Album > Title
// Year - Album -> Year > Album  > Artist > Title

//  Artist - (Year - Album) -> Artist > (Year - Album) > Title
//  Album Artist - (Year - Album) -> Album Artist > (Year - Album) > Title

// GetArtistYearAlbumModel -
func GetArtistYearAlbumModel(useAlbumArtist bool) (s *gtk.TreeStore, err error) {
	s, err = gtk.TreeStoreNew(
		glib.TYPE_STRING,
	)
	if err != nil {
		return
	}

	/*
		var artist, yearAlbum, title *gtk.TreeIter

		for _, v := range cts {
			artist, err = store.GetIterFromString(v.Track.Artist)
			if err != nil {
				continue
			} else if artist != nil {
				var parent *gtk.TreeIter
				store.IterParent(parent, artist)
				if parent != nil {
					continue
				}
			} else {
				artist = store.Append(nil)
				err = store.SetValue(artist, 0, v.Track.Artist)
			}
			if err != nil {
				return
			}
		}
	*/

	return
}

func subscribeToCollectionStore() {
	log.Info("Subscribing to collection store")

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
		glib.IdleAdd(updateCollections)
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
		return
	}
}

func getCollectionTree() {
	const path = "/collections/tree"

	/*
		uri := base.Conf.Server.GetURL() + path

		res, err := http.Get(uri)
		if err != nil || !httpstatus.IsSuccess(res) {
			err = onerror.LogHTTP(err, res, false)
			return
		}
		defer res.Body.Close()

		r, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error(err)
		}

		cts = []models.CollectionTrack{}
		err = json.Unmarshal(r, cts)

		view, _ := createViewAndModel()
	*/

	return
}

func updateCollections() bool {
	log.Info("Updating collections")

	return false
}

func createViewAndModel() (view *gtk.TreeView, err error) {
	/*
		view, err = gtk.TreeViewNew()
		renderer, err := gtk.CellRendererTextNew()

		cols := []struct {
			name string
			col  int
		}{
			{"Title", titleCol},
			// {"Album", albumCol},
			// {"Artist", artistCol},
			// {"Album Artist", albumArtistCol},
			// {"Genre", genreCol},
			// {"Year", yearCol},
		}

		for _, c := range cols {
			var col *gtk.TreeViewColumn
			col, err = gtk.TreeViewColumnNewWithAttribute(
				c.name,
				renderer,
				"text",
				c.col,
			)
			if err != nil {
				return
			}
			view.InsertColumn(col, -1)
		}

		model, err := store.GetQueueStore(m3uetcpb.Perspective_MUSIC)
		view.SetModel(model)
	*/
	return
}

// int
// main (int argc, char **argv)
// {
//   gtk_init (&argc, &argv);

//   GtkWidget *window = gtk_window_new (GTK_WINDOW_TOPLEVEL);
//   g_signal_connect (window, "destroy", gtk_main_quit, NULL);

//   GtkWidget *view = create_view_and_model ();

//   gtk_container_add (GTK_CONTAINER (window), view);

//   gtk_widget_show_all (window);

//   gtk_main ();

//   return 0;
// }
// Copy
