package playlists

import (
	"strconv"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
)

// EditPlaylist edits a playlist properties
func EditPlaylist(id int64) (err error) {
	pl := store.GetPlaylist(id)
	nameIn := ""
	if !pl.Transient {
		nameIn = pl.Name
	}
	descrIn := pl.Description
	pgID := pl.PlaylistGroupId

	nameent, err := builder.GetEntry("playlist_dialog_name")
	if err != nil {
		return
	}
	descrent, err := builder.GetEntry("playlist_dialog_descr")
	if err != nil {
		return
	}
	groups, err := builder.GetComboBoxText("playlist_dialog_pg")
	if err != nil {
		return
	}

	updBtn, err := builder.GetButton("playlist_dialog_btn_apply")
	if err != nil {
		return
	}

	nameent.SetText(nameIn)
	descrent.SetText(descrIn)

	groups.RemoveAll()
	groups.Append("0", "--")
	activeIdx := 0
	count := 0
	store.BData.Mu.Lock()
	for _, pg := range store.BData.PlaylistGroup {
		groups.Append(strconv.FormatInt(pg.Id, 10), pg.Name)
		count++
		if pgID == pg.Id {
			activeIdx = count
		}
	}
	store.BData.Mu.Unlock()

	groups.SetActiveID(strconv.FormatInt(pgID, 10))
	groups.SetActive(activeIdx)

	updBtn.SetSensitive(nameIn != "")
	nameent.Connect("changed", func(e *gtk.Entry) {
		name, _ := e.GetText()
		if name == "" {
			updBtn.SetSensitive(false)
			return
		}
		if name == nameIn {
			updBtn.SetSensitive(true)
			return
		}
		updBtn.SetSensitive(!store.PlaylistAlreadyExists(name))
	})

	res := playlistDlg.Run()
	defer playlistDlg.Hide()

	switch res {
	case gtk.RESPONSE_APPLY:
		var name, descr string
		name, err = nameent.GetText()
		if err != nil {
			return
		}
		descr, err = descrent.GetText()
		if err != nil {
			return
		}
		pgActive := groups.GetActiveID()
		pgID, err = strconv.ParseInt(pgActive, 10, 64)
		if err != nil {
			return
		}
		if pgID == 0 {
			pgID = -1
		}
		req := &m3uetcpb.ExecutePlaylistActionRequest{
			Action:          m3uetcpb.PlaylistAction_PL_UPDATE,
			Id:              id,
			Name:            name,
			Description:     descr,
			PlaylistGroupId: pgID,
		}
		if descr == "" {
			req.ResetDescription = true
		}
		_, err = store.ExecutePlaylistAction(req)
		onerror.Log(err)
	case gtk.RESPONSE_CANCEL:
	default:
	}
	return
}
