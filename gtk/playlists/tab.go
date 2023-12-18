package playlists

import (
	"fmt"
	"log/slog"

	"github.com/diamondburned/gotk4/pkg/gdk/v3"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
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

		for ipage := 0; ipage < nb.NPages(); ipage++ {
			page := nb.NthPage(ipage)
			if page == nil {
				slog.Warn("Failed to get page from notebook", "page", ipage)
				continue
			}
			header := nb.TabLabel(page)
			pageName := gtk.BaseWidget(header).Name()
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
	btn := event.AsButton()
	if btn.Button() == gdk.BUTTON_SECONDARY {
		ot.pageMenu.PopupAtPointer(event)
	}
}

func (ot *onTab) contextUpdate(mi *gtk.MenuItem) {
	onerror.Log(EditPlaylist(ot.id))
}

func (ot *onTab) createContextMenus() (err error) {
	ctxMenu := gtk.NewMenu()
	ctxMenu.SetVisible(true)

	miSuffix, _ := chars.GetRandomLetters(6)

	miPlayNow := gtk.NewMenuItemWithLabel("Play now")
	if miPlayNow == nil {
		err = fmt.Errorf("failed to create menu item: play-now")
		return
	}
	miPlayNow.SetVisible(true)
	miPlayNow.Connect("activate", ot.ContextPlayNow)
	miPlayNow.SetName(fmt.Sprintf("menuitem-%s-%s", "playnow", miSuffix))
	ctxMenu.Add(miPlayNow)

	miEnqueue := gtk.NewMenuItemWithLabel("Enqueue")
	if miEnqueue == nil {
		err = fmt.Errorf("failed to create menu item: enqueue")
		return
	}
	miEnqueue.SetVisible(true)
	miEnqueue.Connect("activate", ot.ContextEnqueue)
	miEnqueue.SetName(fmt.Sprintf("menuitem-%s-%s", "enqueue", miSuffix))
	ctxMenu.Add(miEnqueue)

	sep1 := gtk.NewSeparatorMenuItem()
	sep1.SetVisible(true)
	ctxMenu.Add(sep1)

	miTop := gtk.NewMenuItemWithLabel("Move to top")
	if miTop == nil {
		err = fmt.Errorf("failed to create menu item: move-to-top")
		return
	}
	miTop.SetVisible(true)
	miTop.Connect("activate", ot.ContextMove)
	miTop.SetName(fmt.Sprintf("menuitem-%s-%s", "top", miSuffix))
	ctxMenu.Add(miTop)

	miUp := gtk.NewMenuItemWithLabel("Move up")
	if miUp == nil {
		err = fmt.Errorf("failed to create menu item: move-up")
		return
	}
	miUp.SetVisible(true)
	miUp.Connect("activate", ot.ContextMove)
	miUp.SetName(fmt.Sprintf("menuitem-%s-%s", "up", miSuffix))
	ctxMenu.Add(miUp)

	miDown := gtk.NewMenuItemWithLabel("Move down")
	if miDown == nil {
		err = fmt.Errorf("failed to create menu item: move-down")
		return
	}
	miDown.SetVisible(true)
	miDown.Connect("activate", ot.ContextMove)
	miDown.SetName(fmt.Sprintf("menuitem-%s-%s", "down", miSuffix))
	ctxMenu.Add(miDown)

	miBottom := gtk.NewMenuItemWithLabel("Move to bottom")
	if miBottom == nil {
		err = fmt.Errorf("failed to create menu item: move-to-bottom")
		return
	}
	miBottom.SetVisible(true)
	miBottom.Connect("activate", ot.ContextMove)
	miBottom.SetName(fmt.Sprintf("menuitem-%s-%s", "bottom", miSuffix))
	ctxMenu.Add(miBottom)

	sep2 := gtk.NewSeparatorMenuItem()
	sep2.SetVisible(true)
	ctxMenu.Add(sep2)

	miDelete := gtk.NewMenuItemWithLabel("Remove from playlist")
	if miDelete == nil {
		err = fmt.Errorf("failed to create menu item: remove-from-playlist")
		return
	}
	miDelete.SetVisible(true)
	miDelete.Connect("activate", ot.ContextDelete)
	miDelete.SetName(fmt.Sprintf("menuitem-%s-%s", "delete", miSuffix))
	ctxMenu.Add(miDelete)

	ctxMenu.Connect("hide", ot.ContextHide)
	ctxMenu.Connect("popped-up", ot.ContextPoppedUp)
	ot.ctxMenu = ctxMenu

	pageMenu := gtk.NewMenu()
	pageMenu.SetVisible(true)

	pl := store.BData.GetOpenPlaylist(ot.id)
	if pl == nil {
		err = fmt.Errorf("playlist no longer available")
		return
	}
	label := "Edit properties"
	if pl.Transient {
		label = "Save as..."
	}
	miUpdate := gtk.NewMenuItemWithLabel(label)
	if miUpdate == nil {
		err = fmt.Errorf("failed to create menu item: %s", label)
		return
	}
	miUpdate.SetVisible(true)
	miUpdate.Connect("activate", ot.contextUpdate)
	pageMenu.Add(miUpdate)

	sep3 := gtk.NewSeparatorMenuItem()
	sep3.SetVisible(true)
	pageMenu.Add(sep3)

	miClear := gtk.NewMenuItemWithLabel("Clear playlist")
	if miClear == nil {
		err = fmt.Errorf("failed to create menu item: clear-playlist")
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
	ot.img = gtk.NewImageFromIconName(
		"media-playback-start",
		int(gtk.IconSizeMenu),
	)
	if ot.img == nil {
		return
	}
	ot.img.SetVisible(false)

	ot.label = gtk.NewLabel("")
	ot.label.SetVisible(true)
	ot.updateLabel()
	return
}

func (ot *onTab) setTabContent() (vbox *gtk.Box, err error) {
	err = ot.setTreeView()
	if err != nil {
		return
	}

	scrolled := gtk.NewScrolledWindow(nil, nil)
	scrolled.SetVisible(true)

	vbox = gtk.NewBox(gtk.OrientationVertical, 0)
	vbox.SetVisible(true)

	scrolled.Add(ot.view)
	vbox.PackStart(scrolled, true, true, 0)
	return
}

func (ot *onTab) setTabHeader() (hbox *gtk.Box, err error) {
	if err = ot.setLabel(); err != nil {
		return
	}

	ebox := gtk.NewEventBox()
	ebox.SetVisible(true)
	ebox.Add(ot.label)
	ebox.Connect("button-press-event", ot.contextEvent)

	closeBtn := gtk.NewButton()
	closeBtn.SetVisible(true)
	closeBtn.SetCanDefault(false)
	closeBtn.SetCanFocus(false)
	closeBtn.SetRelief(gtk.ReliefNone)

	img := gtk.NewImageFromIconName("window-close", int(gtk.IconSizeMenu))
	if err != nil {
		return
	}
	img.SetVisible(true)
	closeBtn.SetImage(img)
	closeBtn.Connect("clicked", ot.doClose)

	hbox = gtk.NewBox(gtk.OrientationHorizontal, 0)
	hbox.SetVisible(true)
	hbox.PackStart(ot.img, false, false, 0)
	hbox.PackStart(ebox, true, true, 0)
	hbox.PackEnd(closeBtn, false, false, 0)
	ot.headerName = fmt.Sprintf("playlist-%d", ot.id)
	hbox.SetName(ot.headerName)
	return
}

func (ot *onTab) setTreeView() (err error) {
	ot.view = gtk.NewTreeView()
	ot.view.SetVisible(true)

	renderer := gtk.NewCellRendererText()

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
		col := gtk.NewTreeViewColumn()
		col.SetTitle(store.TColumns[c.colID].Name)
		col.PackStart(renderer, true)
		col.AddAttribute(
			renderer,
			"text",
			int(c.colID),
		)
		col.AddAttribute(renderer, "weight", int(store.TColFontWeight))
		col.SetSizing(gtk.TreeViewColumnAutosize)
		col.SetResizable(true)
		col.SetVisible(c.visible)
		ot.view.InsertColumn(col, -1)
	}

	model, err := store.CreatePlaylistModel(ot.id)
	if err != nil {
		return
	}
	ot.view.SetModel(model)

	sel := ot.view.Selection()
	sel.SetMode(gtk.SelectionMultiple)
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
