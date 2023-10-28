package playlists

import (
	"fmt"
	"log/slog"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/bnp/chars"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type onTab struct {
	*onContext

	headerName string
	pageMenu   *gtk.Menu
	img        *gtk.Image
	label      *gtk.Label
}

func updatePlaybarView() {
	slog.Info("Updating playbar view")

	keep := []onTab{}
	remove := []onTab{}

outer:
	for i := range tabsList {
		for _, opl := range store.BData.GetOpenPlaylists() {
			if opl.Id == tabsList[i].id {
				keep = append(keep, tabsList[i])
				continue outer
			}
		}
		remove = append(remove, tabsList[i])
	}

	for i := range remove {
		nb, err := builder.GetNotebook(perspToNotebook[remove[i].perspective])
		if err != nil {
			slog.Error("Failed to get perspective notebook", "error", err)
			continue
		}

		for ipage := 0; ipage < nb.GetNPages(); ipage++ {
			page, err := nb.GetNthPage(ipage)
			if err != nil {
				slog.With(
					"page", ipage,
					"error", err,
				).Warn("Failed to get page from notebook")
				continue
			}
			header, _ := nb.GetTabLabel(page)
			pageName, _ := header.ToWidget().GetName()
			if remove[i].headerName == pageName {
				nb.RemovePage(ipage)
			}
		}

		store.DestroyPlaylistModel(remove[i].id)
	}

	tabsList = keep

	for persp, pls := range store.PerspectiveToPlaylists {
		nb, err := builder.GetNotebook(perspToNotebook[persp])
		if err != nil {
			slog.With(
				"notebook", perspToNotebook[persp],
				"error", err,
			).Error("Failed to get notebook from builder")
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

			tab := onTab{
				onContext: &onContext{
					id:          pl.Id,
					perspective: persp,
				},
			}

			vbox, err := tab.setTabContent()
			if err != nil {
				slog.Error("Failed to set tab content", "error", err)
				continue mid
			}

			hbox, err := tab.setTabHeader()
			if err != nil {
				slog.Error("Failed to set tab header", "error", err)
				continue mid
			}

			nb.AppendPage(vbox, hbox)

			if err = tab.createContextMenus(); err != nil {
				slog.Error("Failed to create context menus", "error", err)
				continue mid
			}

			tabsList = append(tabsList, tab)
		}
	}

	setFocused()
	go UpdateStatusBar(statusBarDigest)
}

func (ot *onTab) contextEvent(_ *gtk.EventBox, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() == gdk.BUTTON_SECONDARY {
		ot.pageMenu.PopupAtPointer(event)
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

	miSuffix, _ := chars.GetRandomLetters(6)

	miPlayNow, err := gtk.MenuItemNewWithLabel("Play now")
	if err != nil {
		return
	}
	miPlayNow.SetVisible(true)
	miPlayNow.Connect("activate", ot.ContextPlayNow)
	miPlayNow.SetName(fmt.Sprintf("menuitem-%s-%s", "playnow", miSuffix))
	ctxMenu.Add(miPlayNow)

	miEnqueue, err := gtk.MenuItemNewWithLabel("Enqueue")
	if err != nil {
		return
	}
	miEnqueue.SetVisible(true)
	miEnqueue.Connect("activate", ot.ContextEnqueue)
	miEnqueue.SetName(fmt.Sprintf("menuitem-%s-%s", "enqueue", miSuffix))
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
	miTop.Connect("activate", ot.ContextMove)
	miTop.SetName(fmt.Sprintf("menuitem-%s-%s", "top", miSuffix))
	ctxMenu.Add(miTop)

	miUp, err := gtk.MenuItemNewWithLabel("Move up")
	if err != nil {
		return
	}
	miUp.SetVisible(true)
	miUp.Connect("activate", ot.ContextMove)
	miUp.SetName(fmt.Sprintf("menuitem-%s-%s", "up", miSuffix))
	ctxMenu.Add(miUp)

	miDown, err := gtk.MenuItemNewWithLabel("Move down")
	if err != nil {
		return
	}
	miDown.SetVisible(true)
	miDown.Connect("activate", ot.ContextMove)
	miDown.SetName(fmt.Sprintf("menuitem-%s-%s", "down", miSuffix))
	ctxMenu.Add(miDown)

	miBottom, err := gtk.MenuItemNewWithLabel("Move to bottom")
	if err != nil {
		return
	}
	miBottom.SetVisible(true)
	miBottom.Connect("activate", ot.ContextMove)
	miBottom.SetName(fmt.Sprintf("menuitem-%s-%s", "bottom", miSuffix))
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
	miDelete.Connect("activate", ot.ContextDelete)
	miDelete.SetName(fmt.Sprintf("menuitem-%s-%s", "delete", miSuffix))
	ctxMenu.Add(miDelete)

	ctxMenu.Connect("hide", ot.ContextHide)
	ctxMenu.Connect("popped-up", ot.ContextPoppedUp)
	ot.ctxMenu = ctxMenu

	pageMenu, err := gtk.MenuNew()
	if err != nil {
		return
	}
	pageMenu.SetVisible(true)

	pl := store.BData.GetOpenPlaylist(ot.id)
	if pl == nil {
		err = fmt.Errorf("Playlist no longer available")
		return
	}
	label := "Edit properties"
	if pl.Transient {
		label = "Save as..."
	}
	miUpdate, err := gtk.MenuItemNewWithLabel(label)
	if err != nil {
		return
	}
	miUpdate.SetVisible(true)
	miUpdate.Connect("activate", ot.contextUpdate)
	pageMenu.Add(miUpdate)

	sep3, err := gtk.SeparatorMenuItemNew()
	if err != nil {
		return
	}
	sep3.SetVisible(true)
	pageMenu.Add(sep3)

	miClear, err := gtk.MenuItemNewWithLabel("Clear playlist")
	if err != nil {
		return
	}
	miClear.SetVisible(true)
	miClear.Connect("activate", ot.ContextClear)
	pageMenu.Add(miClear)

	ot.pageMenu = pageMenu
	return
}

func (ot *onTab) dblClicked(tv *gtk.TreeView, path *gtk.TreePath,
	col *gtk.TreeViewColumn) {

	values, err := store.GetTreeViewTreePathValues(
		tv,
		path,
		[]store.ModelColumn{
			store.TColTrackID,
			store.TColLocation,
			store.TColPosition,
		},
	)
	if err != nil {
		slog.Error("Failed to get tree-view's tree-path values", "error", err)
		return
	}
	slog.Debug("Doouble-clicked column values", "values", values)

	pos := values[store.TColPosition].(int)

	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Action:   m3uetcpb.PlaybarAction_BAR_ACTIVATE,
		Position: int32(pos),
		Ids:      []int64{ot.id},
	}

	if err := dialer.ExecutePlaybarAction(req); err != nil {
		slog.Error("Failed to execute playbar action", "error", err)
		return
	}
}

func (ot *onTab) doClose(btn *gtk.Button) {
	req := &m3uetcpb.ExecutePlaybarActionRequest{
		Action: m3uetcpb.PlaybarAction_BAR_CLOSE,
		Ids:    []int64{ot.id},
	}

	if err := dialer.ExecutePlaybarAction(req); err != nil {
		slog.Error("Failed to execute playbar action", "error", err)
	}
}

func (ot *onTab) setLabel() (err error) {
	ot.img, err = gtk.ImageNewFromIconName(
		"media-playback-start",
		gtk.ICON_SIZE_MENU,
	)
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
	ot.headerName = fmt.Sprintf("playlist-%d", ot.id)
	hbox.SetName(ot.headerName)
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

	isQuery := store.BData.GetPlaylist(ot.id).QueryId > 0
	hasLastPlayedFor := store.BData.HasLastPlayedFor(ot.id)

	cols := []struct {
		colID   store.ModelColumn
		visible bool
	}{
		{store.TColPosition, true},
		{store.TColTitle, true},
		{store.TColArtist, true},
		{store.TColAlbum, true},
		{store.TColAlbumartist, false},
		{store.TColComposer, false},
		{store.TColGenre, false},
		{store.TColYear, false},
		{store.TColTrackNumberOverTotal, false},
		{store.TColDiscNumberOverTotal, false},
		{store.TColPlaycount, isQuery},
		{store.TColRating, false},
		{store.TColLastplayed, isQuery},
		{store.TColDuration, !hasLastPlayedFor},
		{store.TColPlayedOverDuration, hasLastPlayedFor},
		{store.TColTrackID, true},
		{store.TColLocation, true},
	}

	for _, c := range cols {
		var col *gtk.TreeViewColumn
		col, err = gtk.TreeViewColumnNewWithAttribute(
			store.TColumns[c.colID].Name,
			renderer,
			"text",
			int(c.colID),
		)
		if err != nil {
			return
		}
		col.AddAttribute(renderer, "weight", int(store.TColFontWeight))
		col.SetSizing(gtk.TREE_VIEW_COLUMN_AUTOSIZE)
		col.SetResizable(true)
		col.SetVisible(c.visible)
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
	sel.SetMode(gtk.SELECTION_MULTIPLE)
	sel.Connect("changed", ot.SelChanged)

	ot.view.Connect("row-activated", ot.dblClicked)
	ot.view.Connect("button-press-event", ot.Context)
	ot.view.Connect("key-press-event", ot.Key)
	return
}

func (ot *onTab) updateLabel() (err error) {
	pl := store.BData.GetOpenPlaylist(ot.id)
	if pl == nil {
		slog.Warn("Playlist no longer available", "ID", ot.id)
		return
	}
	name := cases.Title(language.English).String(pl.Name)
	if pl.Transient {
		name += "*"
	}
	ot.label.SetText(name)
	ot.label.SetTooltipText(pl.Description)

	ot.img.SetVisible(store.BData.GetActiveID() == ot.id)
	return
}
