package gtkui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
)

var cts []models.CollectionTrack

func updateCollection() bool {
	return false
}

func GetCollectionTree() {
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

const (
	titleCol int = iota
	albumCol
	artistCol
	albumArtistCol
	genreCol
	yearCol
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

func ArtistYearAlbumModel(useAlbumArtist bool) (store *gtk.TreeStore, err error) {
	store, err = gtk.TreeStoreNew(
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

func createViewAndModel() (view *gtk.TreeView, err error) {
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

	model, err := ArtistYearAlbumModel(false)
	view.SetModel(model)
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
