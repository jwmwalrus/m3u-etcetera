package musicpane

import (
	"log/slog"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gdk/v3"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/gtk/util"
)

type onMusicCollections struct {
	*onContext

	hierarchy        *gtk.ComboBoxText
	hierarchyGrouped *gtk.ToggleButton
}

func createMusicCollections() (omc *onMusicCollections, err error) {
	slog.Info("Creating music collections view and model")

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

	renderer := gtk.NewCellRendererText()

	qcols := []int{
		int(store.CColTree),
	}

	for _, i := range qcols {
		col := gtk.NewTreeViewColumn()
		col.SetTitle(store.CTreeColumn[i].Name)
		col.PackStart(renderer, true)
		col.AddAttribute(
			renderer,
			"text",
			i,
		)
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
	btn := event.AsButton()
	if btn.Button() == gdk.BUTTON_SECONDARY {
		menu, err := builder.GetMenu("collections_view_context")
		if err != nil {
			slog.With(
				"menu", "collections_view_context",
				"error", err,
			).Error("Failed to get menu from builder")
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
		if strings.Contains(mi.Label(), "Prepend") {
			action = m3uetcpb.PlaylistTrackAction_PT_PREPEND
		}
		req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
			PlaylistId: plID,
			Action:     action,
			TrackIds:   ids,
		}

		if err := dialer.ExecutePlaylistTrackAction(req); err != nil {
			slog.Error("Failed to execute playlist track action", "error", err)
			return
		}
	} else {
		action := m3uetcpb.QueueAction_Q_APPEND
		if strings.Contains(mi.Label(), "Prepend") {
			action = m3uetcpb.QueueAction_Q_PREPEND
		}
		req := &m3uetcpb.ExecuteQueueActionRequest{
			Action: action,
			Ids:    ids,
		}

		if err := dialer.ExecuteQueueAction(req); err != nil {
			slog.Error("Failed to execute queue action", "error", err)
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
		slog.Error("Failed to execute playback action", "error", err)
		return
	}
}

func (omc *onMusicCollections) dblClicked(tv *gtk.TreeView,
	path *gtk.TreePath, col *gtk.TreeViewColumn) {

	values, err := store.GetTreeViewTreePathValues(
		tv,
		path,
		[]store.ModelColumn{store.CColTree, store.CColTreeIDList},
	)
	if err != nil {
		slog.Error("Failed to get tree-view's tree-path values", "error", err)
		return
	}
	slog.Debug("Doouble-clicked column value", "value", values[store.CColTree])

	ids, err := util.StringToIDList(values[store.CColTreeIDList].(string))
	if err != nil {
		slog.Error("Failed to convert string to ID list", "error", err)
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
			slog.Error("Failed to execute playlist track action", "error", err)
			return
		}
	} else {
		req := &m3uetcpb.ExecuteQueueActionRequest{
			Action: m3uetcpb.QueueAction_Q_APPEND,
			Ids:    ids,
		}

		if err := dialer.ExecuteQueueAction(req); err != nil {
			slog.Error("Failed to execute queue action", "error", err)
			return
		}
	}
}

func (omc *onMusicCollections) filtered(se *gtk.SearchEntry) {
	text := se.Text()
	store.FilterCollectionTreeBy(text)
}

func (omc *onMusicCollections) hierarchyChanged(cbt *gtk.ComboBoxText) {
	id := cbt.ActiveID()
	slog.With("activeText", id).
		Info("Collection hierarchy changed")

	grouped := omc.hierarchyGrouped.Active()

	go func() {
		store.CData.SwitchHierarchyTo(id, grouped)
	}()
}

func (omc *onMusicCollections) hierarchyGroupToggled(cb *gtk.ToggleButton) {
	grouped := cb.Active()
	slog.With("grouped", grouped).
		Info("Collection hierarchy group toggled")

	id := omc.hierarchy.ActiveID()

	go func() {
		store.CData.SwitchHierarchyTo(id, grouped)
	}()
}
