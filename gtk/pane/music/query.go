package musicpane

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/playlists"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/gtk/util"
)

type onMusicQuery struct {
	*onContext

	dlg                               *gtk.Dialog
	name, id, descr, params, from, to *gtk.Entry
	rating, limit                     *gtk.SpinButton
	random                            *gtk.CheckButton
	resultsLabel                      *gtk.Label
}

func createMusicQueries() (omqy *onMusicQuery, err error) {
	slog.Info("Creating music queries view and model")

	omqy = &onMusicQuery{
		onContext: &onContext{ct: queryContext},
	}

	omqy.view, err = builder.GetTreeView("queries_view")
	if err != nil {
		return
	}

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	qcols := []int{
		int(store.QYColTree),
	}

	for _, i := range qcols {
		var col *gtk.TreeViewColumn
		col, err = gtk.TreeViewColumnNewWithAttribute(
			store.QYTreeColumn[i].Name,
			renderer,
			"text",
			i,
		)
		if err != nil {
			return
		}
		omqy.view.InsertColumn(col, -1)
	}

	model, err := store.CreateQueryTreeModel()
	if err != nil {
		return
	}

	omqy.view.SetModel(model)
	return
}

func (omqy *onMusicQuery) context(tv *gtk.TreeView, event *gdk.Event) {
	btn := gdk.EventButtonNewFromEvent(event)
	if btn.Button() != gdk.BUTTON_SECONDARY {
		return
	}

	ids := omqy.getSelection(true)
	if len(ids) != 1 {
		return
	}

	menu, err := builder.GetMenu("queries_view_context")
	if err != nil {
		slog.With(
			"menu", "queries_view_context",
			"error", err,
		).Error("Failed to get menu from builder")
		return
	}

	qy := store.QYData.GetQuery(ids[0])

	miEdit, err := builder.GetMenuItem("queries_view_context_edit")
	if err != nil {
		slog.With(
			"menu-item", "queries_view_context_edit",
			"error", err,
		).Error("Failed to get menu item from builder")
		return
	}
	miEdit.SetSensitive(!qy.ReadOnly)

	miDelete, err := builder.GetMenuItem("queries_view_context_delete")
	if err != nil {
		slog.With(
			"menu-item", "queries_view_context_delete",
			"error", err,
		).Error("Failed to get menu item from builder")
		return
	}
	miDelete.SetSensitive(!qy.ReadOnly)

	miToPl, err := builder.GetMenuItem("queries_view_context_to_playlist")
	if err != nil {
		slog.With(
			"menu-item", "queries_view_context_to_playlist",
			"error", err,
		).Error("Failed to get menu item from builder")
		return
	}
	id := playlists.GetFocused(m3uetcpb.Perspective_MUSIC)
	miToPl.SetSensitive(id > 0)

	menu.PopupAtPointer(event)
}

func (omqy *onMusicQuery) contextAppend(mi *gtk.MenuItem) {
	ids := omqy.getSelection()
	if len(ids) != 1 {
		slog.Error("Query selection vanished?")
		return
	}

	label := mi.GetLabel()
	if strings.Contains(label, "playlist") {
		req := &m3uetcpb.QueryInPlaylistRequest{
			Id:         ids[0],
			PlaylistId: playlists.GetFocused(m3uetcpb.Perspective_MUSIC),
		}

		if _, err := dialer.QueryInPlaylist(req); err != nil {
			slog.Error("Failed to get query in playlist", "error", err)
		}
		return
	}

	req := &m3uetcpb.QueryInQueueRequest{
		Id: ids[0],
	}

	if err := dialer.QueryInQueue(req); err != nil {
		slog.Error("Failed to get query in queue", "error", err)
	}
}

func (omqy *onMusicQuery) contextDelete(mi *gtk.MenuItem) {
	ids := omqy.getSelection()
	if len(ids) != 1 {
		slog.Error("Query selection vanished?")
		return
	}

	req := &m3uetcpb.RemoveQueryRequest{
		Id: ids[0],
	}

	if err := dialer.RemoveQuery(req); err != nil {
		slog.Error("Failed to remove query", "error", err)
	}
}

func (omqy *onMusicQuery) contextEdit(mi *gtk.MenuItem) {
	ids := omqy.getSelection()
	if len(ids) != 1 {
		slog.Error("Query selection vanished?")
		return
	}

	if err := omqy.edit(ids[0]); err != nil {
		slog.Error("Failed to edit query", "error", err)
		return
	}
}

func (omqy *onMusicQuery) contextNewPlaylist(mi *gtk.MenuItem) {
	ids := omqy.getSelection()
	if len(ids) != 1 {
		slog.Error("Query selection vanished?")
		return
	}

	if err := omqy.newPlaylist(ids[0]); err != nil {
		slog.Error("Failed to create playlist from query", "error", err)
	}
}

func (omqy *onMusicQuery) createDialog() (err error) {
	slog.Info("Creating query dialog")

	err = builder.AddFromFile("ui/pane/query-dialog.ui")
	if err != nil {
		err = fmt.Errorf("Unable to add query-dialog file to builder: %v", err)
		return
	}

	omqy.dlg, err = builder.GetDialog("query_dialog")
	if err != nil {
		return
	}

	view, err := builder.GetTreeView("query_dialog_results_view")
	if err != nil {
		return
	}

	model, err := store.CreateQueryResultsModel()
	if err != nil {
		return
	}

	qyr := store.Renderer{Model: model, Columns: store.TColumns}

	textro, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	togglerw, err := qyr.GetActivatable(store.TColToggleSelect)
	if err != nil {
		return
	}

	cols := []struct {
		idx  store.ModelColumn
		rend gtk.ICellRenderer
		rsz  bool
	}{
		{store.TColNumber, textro, false},
		{store.TColTrackID, textro, false},
		{store.TColTitle, textro, true},
		{store.TColArtist, textro, true},
		{store.TColAlbum, textro, true},
		{store.TColDuration, textro, false},
		{store.TColToggleSelect, togglerw, false},
	}

	for _, v := range cols {
		var col *gtk.TreeViewColumn
		if renderer, ok := v.rend.(*gtk.CellRendererToggle); ok {
			col, err = gtk.TreeViewColumnNewWithAttribute(
				store.TColumns[v.idx].Name,
				renderer,
				"active",
				int(v.idx),
			)
		} else if renderer, ok := v.rend.(*gtk.CellRendererText); ok {
			col, err = gtk.TreeViewColumnNewWithAttribute(
				store.TColumns[v.idx].Name,
				renderer,
				"text",
				int(v.idx),
			)
		} else {
			slog.Error("¿Cómo sabré si es pez o iguana?")
			continue
		}
		if err != nil {
			return
		}
		col.SetResizable(v.rsz)
		view.InsertColumn(col, -1)
	}

	view.SetModel(model)

	omqy.name, err = builder.GetEntry("query_dialog_name")
	if err != nil {
		slog.With(
			"entry", "query_dialog_name",
			"error", err,
		).Error("Failed to get entry from builder")
		return
	}
	omqy.id, err = builder.GetEntry("query_dialog_id")
	if err != nil {
		slog.With(
			"entry", "query_dialog_id",
			"error", err,
		).Error("Failed to get entry from builder")
		return
	}
	omqy.descr, err = builder.GetEntry("query_dialog_description")
	if err != nil {
		slog.With(
			"entry", "query_dialog_description",
			"error", err,
		).Error("Failed to get entry from builder")
		return
	}
	omqy.params, err = builder.GetEntry("query_dialog_params")
	if err != nil {
		slog.With(
			"entry", "query_dialog_params",
			"error", err,
		).Error("Failed to get entry from builder")
		return
	}
	omqy.from, err = builder.GetEntry("query_dialog_from")
	if err != nil {
		slog.With(
			"entry", "query_dialog_from",
			"error", err,
		).Error("Failed to get entry from builder")
		return
	}
	omqy.to, err = builder.GetEntry("query_dialog_to")
	if err != nil {
		slog.With(
			"entry", "query_dialog_to",
			"error", err,
		).Error("Failed to get entry from builder")
		return
	}

	omqy.rating, err = builder.GetSpinButton("query_dialog_rating")
	if err != nil {
		slog.With(
			"entry", "query_dialog_rating",
			"error", err,
		).Error("Failed to get spin button from builder")
		return
	}
	omqy.limit, err = builder.GetSpinButton("query_dialog_limit")
	if err != nil {
		slog.With(
			"entry", "query_dialog_limit",
			"error", err,
		).Error("Failed to get spin button from builder")
		return
	}

	omqy.random, err = builder.GetCheckButton("query_dialog_random")
	if err != nil {
		slog.With(
			"entry", "query_dialog_random",
			"error", err,
		).Error("Failed to get check button from builder")
		return
	}

	omqy.resultsLabel, err = builder.GetLabel("query_dialog_results_count")
	if err != nil {
		slog.With(
			"entry", "query_dialog_results_count",
			"error", err,
		).Error("Failed to get check button from builder")
		return
	}
	omqy.resultsLabel.SetVisible(false)

	return
}

func (omqy *onMusicQuery) dblClicked(tv *gtk.TreeView,
	path *gtk.TreePath, col *gtk.TreeViewColumn) {

	values, err := store.GetTreeViewTreePathValues(
		tv,
		path,
		[]store.ModelColumn{store.QYColTree, store.QYColTreeIDList},
	)
	if err != nil {
		slog.Error("Failed to get tree-view's tree-path values", "error", err)
		return
	}
	slog.Debug("Doouble-clicked column value", "value", values[store.CColTree])

	ids, err := util.StringToIDList(values[store.QYColTreeIDList].(string))
	if err != nil {
		slog.Error("Failed to convert string to ID list", "error", err)
		return
	}

	if len(ids) != 1 {
		slog.Error("Length of ids is different from 1", "IDs", ids)
		return
	}

	if err := omqy.newPlaylist(ids[0]); err != nil {
		slog.Error("Failed to edit query", "error", err)
		return
	}
}

func (omqy *onMusicQuery) defineQuery(btn *gtk.ToolButton) {
	omqy.resetDialog()

	omqy.id.SetText("0")

	res := omqy.dlg.Run()
	switch res {
	case gtk.RESPONSE_APPLY:
		qy, err := omqy.getQuery()
		if err != nil {
			slog.Error("Failed to get query", "error", err)
		} else {
			req := &m3uetcpb.AddQueryRequest{Query: qy}
			dialer.AddQuery(req)
		}
	case gtk.RESPONSE_CANCEL:
	default:
	}
	omqy.dlg.Hide()
}

func (omqy *onMusicQuery) doSearch(btn *gtk.Button) {
	qy, err := omqy.getQuery()
	if err != nil {
		slog.Error("Failed to get query", "error", err)
		return
	}
	qy.Name = ""
	req := &m3uetcpb.QueryByRequest{Query: qy}
	count, err := dialer.QueryBy(req)
	if err != nil {
		slog.Error("Failed to query-by", "error", err)
		return
	}

	omqy.resultsLabel.SetText(fmt.Sprintf("Results: %v", count))
	omqy.resultsLabel.SetVisible(true)
}

func (omqy *onMusicQuery) edit(id int64) (err error) {
	qy := store.QYData.GetQuery(id)
	if qy == nil {
		slog.Error("Query returned from store is nil")
		return
	}

	if err = omqy.setQuery(qy); err != nil {
		return
	}

	res := omqy.dlg.Run()
	defer omqy.dlg.Hide()
	switch res {
	case gtk.RESPONSE_APPLY:
		qy, err = omqy.getQuery()
		if err != nil {
			return
		}
		req := &m3uetcpb.UpdateQueryRequest{Query: qy}
		dialer.UpdateQuery(req)
	case gtk.RESPONSE_CANCEL:
	default:
	}
	return
}

func (omqy *onMusicQuery) filtered(se *gtk.SearchEntry) {
	text, err := se.GetText()
	if err != nil {
		slog.Error("Failed to get search entry text", "error", err)
		return
	}
	store.FilterQueryTreeBy(text)

}

func (omqy *onMusicQuery) getQuery() (qy *m3uetcpb.Query, err error) {
	name, err := omqy.name.GetText()
	if err != nil {
		slog.Error("Failed to get query name's text", "error", err)
		return
	}
	descr, err := omqy.descr.GetText()
	if err != nil {
		slog.Error("Failed to get query description's text", "error", err)
		return
	}
	params, err := omqy.params.GetText()
	if err != nil {
		slog.Error("Failed to get query param's text", "error", err)
		return
	}
	ids, err := store.GetQueryResultsSelections()
	onerror.Log(err)
	if len(ids) > 0 {
		if params != "" {
			params += " and "
		}
		params += "id=" + util.IDListToString(ids)
	}

	idTxt, err := omqy.id.GetText()
	if err != nil {
		slog.Error("Failed to get query ID's text", "error", err)
		return
	}
	id, err := strconv.ParseInt(idTxt, 10, 64)
	if err != nil {
		slog.With(
			"ID", idTxt,
			"error", err,
		).Error("Failed to parse query ID")
		return
	}

	var from, to int64

	fromTxt, err := omqy.from.GetText()
	if err != nil {
		slog.Error("Failed to get query `from`'s text", "error", err)
		return
	}
	if fromTxt != "" && fromTxt != "0" {
		var ft time.Time
		ft, err = time.Parse("2006/01/02", fromTxt+"/01/01")
		if err != nil {
			slog.With(
				"from", fromTxt,
				"error", err,
			).Error("Failed to parse query `from`")
		} else {
			from = ft.UnixNano()
		}
	}

	toTxt, err := omqy.to.GetText()
	if err != nil {
		slog.Error("Failed to get query `to`'s text", "error", err)
		return
	}
	if toTxt != "" && toTxt != "0" {
		var tt time.Time
		tt, err = time.Parse("2006/01/02", toTxt+"/01/01")
		if err != nil {
			slog.With(
				"to", toTxt,
				"error", err,
			).Error("Failed to parse query `to`")
		} else {
			to = tt.UnixNano()
		}
	}

	qy = &m3uetcpb.Query{
		Name:        name,
		Id:          id,
		Description: descr,
		Params:      params,
		From:        from,
		To:          to,
		Rating:      int32(omqy.rating.GetValue()),
		Limit:       int32(omqy.limit.GetValue()),
		Random:      omqy.random.GetActive(),
	}
	return
}

func (omqy *onMusicQuery) newPlaylist(id int64) error {
	req := &m3uetcpb.QueryInPlaylistRequest{
		Id: id,
	}

	var playlistID int64
	var err error
	if playlistID, err = dialer.QueryInPlaylist(req); err != nil {
		slog.Error("failed to get query in playlist", "error", err)

		reqbar := &m3uetcpb.ExecutePlaybarActionRequest{
			Action: m3uetcpb.PlaybarAction_BAR_CLOSE,
			Ids:    []int64{playlistID},
		}

		if err = dialer.ExecutePlaybarAction(reqbar); err != nil {
			return err
		}

		return nil
	}

	playlists.RequestFocus(m3uetcpb.Perspective_MUSIC, playlistID)
	return nil
}

func (omqy *onMusicQuery) resetDialog() error {
	omqy.resultsLabel.SetVisible(false)

	qy := &m3uetcpb.Query{}
	return omqy.setQuery(qy)
}

func (omqy *onMusicQuery) setQuery(qy *m3uetcpb.Query) (err error) {
	store.ClearQueryResults()

	if qy == nil {
		err = fmt.Errorf("Received nil query to set")
		return
	}

	omqy.name.SetText(qy.Name)
	omqy.id.SetText(strconv.FormatInt(qy.Id, 10))
	omqy.descr.SetText(qy.Description)
	omqy.params.SetText(qy.Params)
	omqy.rating.SetValue(float64(qy.Rating))
	omqy.limit.SetValue(float64(qy.Limit))
	omqy.random.SetActive(qy.Random)

	var from, to string
	if qy.From > 0 {
		from = strconv.Itoa(time.UnixMicro(qy.From / 1e3).Year())
	} else {
		from = "0"
	}
	omqy.from.SetText(from)

	if qy.To > 0 {
		to = strconv.Itoa(time.UnixMicro(qy.To / 1e3).Year())
	} else {
		to = "0"
	}
	omqy.to.SetText(to)

	return
}

func (omqy *onMusicQuery) toggleSelection(btn *gtk.Button) {
	store.ToggleQueryResultsSelection()
}
