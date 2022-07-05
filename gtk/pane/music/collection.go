package musicpane

import (
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/gtk/util"
	log "github.com/sirupsen/logrus"
)

type onMusicCollections struct {
	*onContext

	hierarchy        *gtk.ComboBoxText
	hierarchyGrouped *gtk.ToggleButton
}

func createMusicCollections() (omc *onMusicCollections, err error) {
	log.Info("Creating music collections view and model")

	omc = &onMusicCollections{
		onContext: &onContext{ct: collectionContext},
	}

	omc.hierarchy, err = builder.GetComboBoxText("collections_hierarchy")
	if err != nil {
		return
	}

	omc.hierarchyGrouped, err = builder.GetToggleButton("collections_hierarchy_grouped")
	if err != nil {
		return
	}

	omc.view, err = builder.GetTreeView("collections_view")
	if err != nil {
		return
	}

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	qcols := []int{
		int(store.CColTree),
	}

	for _, i := range qcols {
		var col *gtk.TreeViewColumn
		col, err = gtk.TreeViewColumnNewWithAttribute(
			store.CTreeColumn[i].Name,
			renderer,
			"text",
			i,
		)
		if err != nil {
			return
		}
		omc.view.InsertColumn(col, -1)
	}

	model, err := store.CreateCollectionTreeModel(store.ArtistYearAlbumTree)
	if err != nil {
		return
	}

	omc.view.SetModel(model)
	return
}

func (omc *onMusicCollections) context(tv *gtk.TreeView, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() == gdk.BUTTON_SECONDARY {
		menu, err := builder.GetMenu("collections_view_context")
		if err != nil {
			log.Error(err)
			return
		}
		menu.PopupAtPointer(event)
	}
}

func (omc *onMusicCollections) contextAppend(mi *gtk.MenuItem) {
	ids := omc.getSelection()
	if len(ids) == 0 {
		return
	}

	plID := playlists.GetFocused(m3uetcpb.Perspective_MUSIC)

	if plID > 0 {
		action := m3uetcpb.PlaylistTrackAction_PT_APPEND
		if strings.Contains(mi.GetLabel(), "Prepend") {
			action = m3uetcpb.PlaylistTrackAction_PT_PREPEND
		}
		req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
			PlaylistId: plID,
			Action:     action,
			TrackIds:   ids,
		}

		if err := dialer.ExecutePlaylistTrackAction(req); err != nil {
			log.Error(err)
			return
		}
	} else {
		action := m3uetcpb.QueueAction_Q_APPEND
		if strings.Contains(mi.GetLabel(), "Prepend") {
			action = m3uetcpb.QueueAction_Q_PREPEND
		}
		req := &m3uetcpb.ExecuteQueueActionRequest{
			Action: action,
			Ids:    ids,
		}

		if err := dialer.ExecuteQueueAction(req); err != nil {
			log.Error(err)
			return
		}
	}
}

func (omc *onMusicCollections) contextPlayNow(mi *gtk.MenuItem) {
	ids := omc.getSelection()
	if len(ids) == 0 {
		return
	}

	req := &m3uetcpb.ExecutePlaybackActionRequest{
		Action: m3uetcpb.PlaybackAction_PB_PLAY,
		Force:  true,
		Ids:    ids,
	}

	if err := dialer.ExecutePlaybackAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (omc *onMusicCollections) dblClicked(tv *gtk.TreeView,
	path *gtk.TreePath, col *gtk.TreeViewColumn) {

	values, err := store.GetTreeStoreValues(
		tv,
		path,
		[]store.ModelColumn{store.CColTree, store.CColTreeIDList},
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Doouble-clicked column value: %v", values[store.CColTree])

	ids, err := util.StringToIDList(values[store.CColTreeIDList].(string))
	if err != nil {
		log.Error(err)
		return
	}

	plID := playlists.GetFocused(m3uetcpb.Perspective_MUSIC)

	if plID > 0 {
		req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
			PlaylistId: plID,
			Action:     m3uetcpb.PlaylistTrackAction_PT_APPEND,
			TrackIds:   ids,
		}

		if err := dialer.ExecutePlaylistTrackAction(req); err != nil {
			log.Error(err)
			return
		}
	} else {
		req := &m3uetcpb.ExecuteQueueActionRequest{
			Action: m3uetcpb.QueueAction_Q_APPEND,
			Ids:    ids,
		}

		if err := dialer.ExecuteQueueAction(req); err != nil {
			log.Error(err)
			return
		}
	}
}

func (omc *onMusicCollections) filtered(se *gtk.SearchEntry) {
	text, err := se.GetText()
	if err != nil {
		log.Error(err)
		return
	}
	store.FilterCollectionTreeBy(text)
}

func (omc *onMusicCollections) hierarchyChanged(cbt *gtk.ComboBoxText) {
	id := cbt.GetActiveID()
	log.WithField("activeText", id).
		Info("Collection hierarchy changed")

	grouped := omc.hierarchyGrouped.GetActive()

	go func() {
		store.CData.SwitchHierarchyTo(id, grouped)
	}()
}

func (omc *onMusicCollections) hierarchyGroupToggled(cb *gtk.ToggleButton) {
	grouped := cb.GetActive()
	log.WithField("grouped", grouped).
		Info("Collection hierarchy group toggled")

	id := omc.hierarchy.GetActiveID()

	go func() {
		store.CData.SwitchHierarchyTo(id, grouped)
	}()
}
