package playlists

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

// EditPlaylist edits a playlist properties
func EditPlaylist(id int64) (err error) {
	pl := store.GetOpenPlaylist(id)
	name := ""
	if !pl.Transient {
		name = pl.Name
	}
	descr := pl.Description
	pgID := int64(1) // FIXME: implement

	nameent, err := builder.GetEntry("playlist_dialog_name")
	if err != nil {
		return
	}
	descrent, err := builder.GetEntry("playlist_dialog_descr")
	if err != nil {
		return
	}
	nameent.SetText(name)
	descrent.SetText(descr)

	res := playlistDlg.Run()
	switch res {
	case gtk.RESPONSE_APPLY:
		name, err = nameent.GetText()
		if err != nil {
			log.Error(err)
			return
		}
		descr, err = descrent.GetText()
		if err != nil {
			log.Error(err)
			return
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
	playlistDlg.Hide()
	return
}
