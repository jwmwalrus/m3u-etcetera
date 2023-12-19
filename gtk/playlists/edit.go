package playlists

import (
	"strconv"

	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/dialer"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
)

// EditPlaylist edits a playlist properties.
func EditPlaylist(id int64) (err error) {
	pl := store.BData.GetPlaylist(id)
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
	pgnames := store.BData.PlaylistGroupNames()

	for k, v := range pgnames {
		groups.Append(strconv.FormatInt(k, 10), v)
		count++
		if pgID == k {
			activeIdx = count
		}
	}

	groups.SetActiveID(strconv.FormatInt(pgID, 10))
	groups.SetActive(activeIdx)

	updBtn.SetSensitive(nameIn != "")
	nameent.Connect("changed", func(e *gtk.Entry) {
		name := e.Text()
		if name == "" {
			updBtn.SetSensitive(false)
			return
		}
		if name == nameIn {
			updBtn.SetSensitive(true)
			return
		}
		updBtn.SetSensitive(!store.BData.PlaylistAlreadyExists(name))
	})

	res := playlistDlg.Run()
	defer playlistDlg.Hide()

	switch gtk.ResponseType(res) {
	case gtk.ResponseApply:
		var name, descr string
		name = nameent.Text()
		descr = descrent.Text()
		pgActive := groups.ActiveID()
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
		_, err = dialer.ExecutePlaylistAction(req)
		onerror.Log(err)
	case gtk.ResponseCancel:
	default:
	}
	return
}
