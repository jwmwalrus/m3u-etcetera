package playlists

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

type onTab struct {
	selection        interface{}
	page             int
	id               int64
	img              *gtk.Image
	label            *gtk.Label
	view             *gtk.TreeView
	ctxMenu, tabMenu *gtk.Menu
	perspective      m3uetcpb.Perspective
}

func updatePlaybarView() {
	keep := []onTab{}
	remove := []onTab{}

	store.BData.Mu.Lock()
outer:
	for i := range tabsList {
		for j := range store.BData.OpenPlaylist {
			if store.BData.OpenPlaylist[j].Id == tabsList[i].id {
				keep = append(keep, tabsList[i])
				continue outer
			}
		}
		remove = append(remove, tabsList[i])
	}
	store.BData.Mu.Unlock()

	for i := range remove {
		nb, err := builder.GetNotebook(perspToNotebook[remove[i].perspective])
		if err != nil {
			log.Error(err)
			continue
		}
		nb.RemovePage(remove[i].page)
		store.DestroyPlaylistModel(remove[i].id)
	}

	tabsList = keep

	for persp, pls := range store.PerspectiveToPlaylists {
		nb, err := builder.GetNotebook(perspToNotebook[persp])
		if err != nil {
			log.Error(err)
			continue
		}

	mid:
		for _, pl := range pls {
			for _, k := range keep {
				if pl.Id == k.id {
					//already added
					k.updateLabel()
					continue mid
				}
			}

			tab := onTab{id: pl.Id, perspective: persp}

			vbox, err := tab.setTabContent()
			if err != nil {
				log.Error(err)
				continue mid
			}

			hbox, err := tab.setTabHeader()
			if err != nil {
				log.Error(err)
				continue mid
			}

			tab.page = nb.AppendPage(vbox, hbox)

			if err = tab.createContextMenus(); err != nil {
				log.Error(err)
				continue mid
			}

			tabsList = append(tabsList, tab)
		}

	}
}

func (ot *onTab) context(tv *gtk.TreeView, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() == gdk.BUTTON_SECONDARY {
		ot.ctxMenu.PopupAtPointer(event)
	}
}

func (ot *onTab) contextClear(mi *gtk.MenuItem) {
	req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
		PlaylistId: ot.id,
		Action:     m3uetcpb.PlaylistTrackAction_PT_CLEAR,
	}

	if err := store.ExecutePlaylistTrackAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (ot *onTab) contextDelete(mi *gtk.MenuItem) {
	values := ot.getSelection()
	if len(values) == 0 {
		return
	}

	req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
		PlaylistId: ot.id,
		Action:     m3uetcpb.PlaylistTrackAction_PT_DELETE,
		Position:   int32(values[store.TColPosition].(int)),
	}

	if err := store.ExecutePlaylistTrackAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (ot *onTab) contextEnqueue(mi *gtk.MenuItem) {
	values := ot.getSelection()
	if len(values) == 0 {
		return
	}

	id := values[store.TColTrackID].(int64)
	loc := values[store.TColLocation].(string)

	req := &m3uetcpb.ExecuteQueueActionRequest{
		Action: m3uetcpb.QueueAction_Q_APPEND,
	}
	if id > 0 {
		req.Ids = []int64{id}
	} else {
		req.Locations = []string{loc}
	}

	if err := store.ExecuteQueueAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (ot *onTab) contextEvent(_ *gtk.EventBox, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() == gdk.BUTTON_SECONDARY {
		ot.tabMenu.PopupAtPointer(event)
	}
}

func (ot *onTab) contextMove(mi *gtk.MenuItem) {
	values := ot.getSelection()
	if len(values) == 0 {
		return
	}

	l := mi.GetLabel()
	fromPos := values[store.TColPosition].(int)
	var pos int
	if strings.Contains(l, "top") {
		if fromPos == 1 {
			return
		}
		pos = 1
	} else if strings.Contains(l, "up") {
		pos = fromPos - 1
	} else if strings.Contains(l, "down") {
		pos = fromPos + 1
	} else if strings.Contains(l, "bottom") {
		pos = values[store.TColLastPosition].(int)
		if fromPos == pos {
			return
		}
	} else {
		log.Error("Invalid/unsupported queue move")
		return
	}

	req := &m3uetcpb.ExecutePlaylistTrackActionRequest{
		PlaylistId:   ot.id,
		Action:       m3uetcpb.PlaylistTrackAction_PT_MOVE,
		Position:     int32(pos),
		FromPosition: int32(fromPos),
	}

	if err := store.ExecutePlaylistTrackAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (ot *onTab) contextPlayNow(mi *gtk.MenuItem) {
	values := ot.getSelection()
	if len(values) == 0 {
		return
	}

	pos := values[store.TColPosition].(int)

	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Action:   m3uetcpb.PlaybarAction_BAR_ACTIVATE,
		Position: int32(pos),
		Ids:      []int64{ot.id},
	}

	if err := store.ExecutePlaybarAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (ot *onTab) contextUpdate(mi *gtk.MenuItem) {
	onerror.Log(EditPlaylist(ot.id))
}

func (ot *onTab) createContextMenus() (err error) {
	ctxMenu, err := gtk.MenuNew()
	if err != nil {
		return
	}
	ctxMenu.SetVisible(true)

	miPlayNow, err := gtk.MenuItemNewWithLabel("Play now")
	if err != nil {
		return
	}
	miPlayNow.SetVisible(true)
	miPlayNow.Connect("activate", ot.contextPlayNow)
	ctxMenu.Add(miPlayNow)

	miEnqueue, err := gtk.MenuItemNewWithLabel("Enqueue")
	if err != nil {
		return
	}
	miEnqueue.SetVisible(true)
	miEnqueue.Connect("activate", ot.contextEnqueue)
	ctxMenu.Add(miEnqueue)

	sep1, err := gtk.SeparatorMenuItemNew()
	if err != nil {
		return
	}
	sep1.SetVisible(true)
	ctxMenu.Add(sep1)

	miTop, err := gtk.MenuItemNewWithLabel("Move to top")
	if err != nil {
		return
	}
	miTop.SetVisible(true)
	miTop.Connect("activate", ot.contextMove)
	ctxMenu.Add(miTop)

	miUp, err := gtk.MenuItemNewWithLabel("Move up")
	if err != nil {
		return
	}
	miUp.SetVisible(true)
	miUp.Connect("activate", ot.contextMove)
	ctxMenu.Add(miUp)

	miDown, err := gtk.MenuItemNewWithLabel("Move down")
	if err != nil {
		return
	}
	miDown.SetVisible(true)
	miDown.Connect("activate", ot.contextMove)
	ctxMenu.Add(miDown)

	miBottom, err := gtk.MenuItemNewWithLabel("Move to bottom")
	if err != nil {
		return
	}
	miBottom.SetVisible(true)
	miBottom.Connect("activate", ot.contextMove)
	ctxMenu.Add(miBottom)

	sep2, err := gtk.SeparatorMenuItemNew()
	if err != nil {
		return
	}
	sep2.SetVisible(true)
	ctxMenu.Add(sep2)

	miDelete, err := gtk.MenuItemNewWithLabel("Remove from playlist")
	if err != nil {
		return
	}
	miDelete.SetVisible(true)
	miDelete.Connect("activate", ot.contextDelete)
	ctxMenu.Add(miDelete)

	ot.ctxMenu = ctxMenu

	tabMenu, err := gtk.MenuNew()
	if err != nil {
		return
	}
	tabMenu.SetVisible(true)

	pl := store.GetOpenPlaylist(ot.id)
	if pl == nil {
		err = fmt.Errorf("Playlist no longer available")
		return
	}
	label := "Update"
	if pl.Transient {
		label = "Save as..."
	}
	miUpdate, err := gtk.MenuItemNewWithLabel(label)
	if err != nil {
		return
	}
	miUpdate.SetVisible(true)
	miUpdate.Connect("activate", ot.contextUpdate)
	tabMenu.Add(miUpdate)

	sep3, err := gtk.SeparatorMenuItemNew()
	if err != nil {
		return
	}
	sep3.SetVisible(true)
	tabMenu.Add(sep3)

	miClear, err := gtk.MenuItemNewWithLabel("Clear playlist")
	if err != nil {
		return
	}
	miClear.SetVisible(true)
	miClear.Connect("activate", ot.contextClear)
	tabMenu.Add(miClear)

	ot.tabMenu = tabMenu
	return
}

func (ot *onTab) dblClicked(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	values, err := store.GetListStoreValues(
		tv,
		path,
		[]store.ModelColumn{
			store.TColTrackID,
			store.TColLocation,
			store.TColPosition,
		},
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Doouble-clicked column values: %v", values)

	pos := values[store.TColPosition].(int)

	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Action:   m3uetcpb.PlaybarAction_BAR_ACTIVATE,
		Position: int32(pos),
		Ids:      []int64{ot.id},
	}

	if err := store.ExecutePlaybarAction(req); err != nil {
		log.Error(err)
		return
	}
}

func (ot *onTab) doClose(btn *gtk.Button) {
	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Action: m3uetcpb.PlaybarAction_BAR_CLOSE,
		Ids:    []int64{ot.id},
	}

	if err := store.ExecutePlaybarAction(req); err != nil {
		log.Error(err)
	}
}

func (ot *onTab) getSelection(keep ...bool) (values map[store.ModelColumn]interface{}) {
	values, ok := ot.selection.(map[store.ModelColumn]interface{})
	if !ok {
		log.Debug("There is no selection available for queue context")
		values = map[store.ModelColumn]interface{}{}
		return
	}

	reset := true
	if len(keep) > 0 {
		reset = !keep[0]
	}

	if reset {
		ot.selection = nil
	}
	return
}

func (ot *onTab) selChanged(sel *gtk.TreeSelection) {
	var err error
	ot.selection, err = store.GetTreeSelectionValues(
		sel,
		[]store.ModelColumn{
			store.TColPosition,
			store.TColLastPosition,
			store.TColTrackID,
			store.TColLocation,
		},
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Selected collection entres: %v", ot.selection)
}

func (ot *onTab) setLabel() (err error) {
	ot.img, err = gtk.ImageNewFromIconName("media-playback-start", gtk.ICON_SIZE_MENU)
	if err != nil {
		return
	}
	ot.img.SetVisible(false)

	ot.label, err = gtk.LabelNew("")
	if err != nil {
		return
	}
	ot.label.SetVisible(true)
	ot.updateLabel()
	return
}

func (ot *onTab) setTabContent() (vbox *gtk.Box, err error) {
	err = ot.setTreeView()
	if err != nil {
		return
	}

	scrolled, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return
	}
	scrolled.SetVisible(true)

	vbox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}
	vbox.SetVisible(true)

	scrolled.Add(ot.view)
	vbox.PackStart(scrolled, true, true, 0)
	return
}

func (ot *onTab) setTabHeader() (hbox *gtk.Box, err error) {
	if err = ot.setLabel(); err != nil {
		return
	}

	ebox, err := gtk.EventBoxNew()
	if err != nil {
		return
	}
	ebox.SetVisible(true)
	ebox.Add(ot.label)
	ebox.Connect("button-press-event", ot.contextEvent)

	closeBtn, err := gtk.ButtonNew()
	if err != nil {
		return
	}
	closeBtn.SetVisible(true)
	closeBtn.SetCanDefault(false)
	closeBtn.SetCanFocus(false)
	closeBtn.SetRelief(gtk.RELIEF_NONE)

	img, err := gtk.ImageNewFromIconName("window-close", gtk.ICON_SIZE_MENU)
	if err != nil {
		return
	}
	img.SetVisible(true)
	closeBtn.SetImage(img)
	closeBtn.Connect("clicked", ot.doClose)

	hbox, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}
	hbox.SetVisible(true)
	hbox.PackStart(ot.img, false, false, 0)
	hbox.PackStart(ebox, true, true, 0)
	hbox.PackEnd(closeBtn, false, false, 0)
	return
}

func (ot *onTab) setTreeView() (err error) {
	ot.view, err = gtk.TreeViewNew()
	if err != nil {
		return
	}
	ot.view.SetVisible(true)

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	cols := []int{
		int(store.TColPosition),
		int(store.TColTitle),
		int(store.TColArtist),
		int(store.TColAlbum),
		int(store.TColDuration),
		int(store.TColTrackID),
	}

	for _, c := range cols {
		var col *gtk.TreeViewColumn
		col, err = gtk.TreeViewColumnNewWithAttribute(
			store.TColumns[c].Name,
			renderer,
			"text",
			c,
		)
		col.AddAttribute(renderer, "weight", int(store.TColFontWeight))
		if err != nil {
			return
		}
		ot.view.InsertColumn(col, -1)
	}

	model, err := store.CreatePlaylistModel(ot.id)
	if err != nil {
		return
	}
	ot.view.SetModel(model)

	sel, err := ot.view.GetSelection()
	if err != nil {
		return
	}
	sel.Connect("changed", ot.selChanged)

	ot.view.Connect("row-activated", ot.dblClicked)
	ot.view.Connect("button-press-event", ot.context)
	return
}

func (ot *onTab) updateLabel() (err error) {
	pl := store.GetOpenPlaylist(ot.id)
	if pl == nil {
		log.WithField("id", ot.id).Warn("Playlist no longer available")
		return
	}
	name := strings.Title(pl.Name)
	if pl.Transient {
		name += "*"
	}
	ot.label.SetText(name)
	ot.label.SetTooltipText(pl.Description)

	store.BData.Mu.Lock()
	ot.img.SetVisible(store.BData.ActiveID == ot.id)
	store.BData.Mu.Unlock()
	return
}