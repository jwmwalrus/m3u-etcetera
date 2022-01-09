package gtkui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/store"
	log "github.com/sirupsen/logrus"
)

type onSettingsMenu struct {
	dlg           *gtk.Dialog
	discoverBtn   *gtk.ToggleToolButton
	updateTagsBtn *gtk.ToggleToolButton
	pm            *gtk.PopoverMenu
}

func (osm *onSettingsMenu) editCollections(btn *gtk.Button) {
	osm.hide()

	osm.resetCollectionsToggles()

	res := osm.dlg.Run()
	switch res {
	case gtk.RESPONSE_APPLY:
		store.ApplyCollectionChanges(osm.getCollectionsToggles())
	case gtk.RESPONSE_CANCEL:
	default:
	}
	osm.dlg.Hide()
}

func (osm *onSettingsMenu) getCollectionsToggles() (opts store.CollectionsOptions) {
	opts.Discover = osm.discoverBtn.GetActive()
	opts.UpdateTags = osm.updateTagsBtn.GetActive()
	return
}

func (osm *onSettingsMenu) hide() {
	if osm.pm == nil {
		var err error
		osm.pm, err = builder.GetPopoverMenu("settings_menu")
		if err != nil {
			log.Error(err)
			return
		}
	}
	osm.pm.Popdown()
}

func (osm *onSettingsMenu) quitAll(btn *gtk.Button) {
	osm.hide()
	window, err := builder.GetWindow()
	if err != nil {
		log.Error(err)
		return
	}
	// window.Emit("delete-event", gdk.EVENT_DELETE)
	store.SetForceExit()
	window.Destroy()
}

func (osm *onSettingsMenu) resetCollectionsToggles() {
	opts := store.CollectionsOptions{}
	opts.SetDefaults()

	osm.discoverBtn.SetActive(opts.Discover)
	osm.updateTagsBtn.SetActive(opts.UpdateTags)
}

func init() {
	perspectivesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_RADIO,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}

}
