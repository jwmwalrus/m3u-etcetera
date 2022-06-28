package pane

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	audiobookspane "github.com/jwmwalrus/m3u-etcetera/gtk/pane/audiobooks"
	musicpane "github.com/jwmwalrus/m3u-etcetera/gtk/pane/music"
	podcastspane "github.com/jwmwalrus/m3u-etcetera/gtk/pane/podcasts"
	radiopane "github.com/jwmwalrus/m3u-etcetera/gtk/pane/radio"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type paneData struct {
	id    string
	path  string
	setup func(signals *map[string]interface{}) error
}

type paneMap map[m3uetcpb.Perspective]paneData

var (
	paneList paneMap
)

// Add adds pane to notebook
func Add(idx m3uetcpb.Perspective, nb *gtk.Notebook,
	signals *map[string]interface{}) (err error) {

	log.WithField("idx", idx.String()).
		Info("Adding perspective to notebook")

	data, ok := paneList[idx]
	if !ok {
		err = fmt.Errorf("Unsupported perspective")
		return
	}

	if err = builder.AddFromFile(data.path); err != nil {
		err = fmt.Errorf("Unable to add file %v to builder: %v", data.path, err)
		return
	}

	newPane, err := builder.GetPane(data.id)
	if err != nil {
		err = fmt.Errorf("Unable to create %v pane: %v", idx, err)
		return
	}

	label, err := gtk.LabelNew(cases.Title(language.English).String(idx.String()))
	if err != nil {
		err = fmt.Errorf("Unable to create %v label: %v", idx, err)
		return
	}

	nb.AppendPage(newPane, label)

	if err = data.setup(signals); err != nil {
		return
	}
	return
}

func init() {
	paneList = paneMap{
		m3uetcpb.Perspective_MUSIC: paneData{
			id:    "music_perspective_pane",
			path:  "ui/pane/music.ui",
			setup: musicpane.Setup,
		},
		m3uetcpb.Perspective_RADIO: paneData{
			id:    "radio_perspective_pane",
			path:  "ui/pane/radio.ui",
			setup: radiopane.Setup,
		},
		m3uetcpb.Perspective_PODCASTS: paneData{
			id:    "podcasts_perspective_pane",
			path:  "ui/pane/podcasts.ui",
			setup: podcastspane.Setup,
		},
		m3uetcpb.Perspective_AUDIOBOOKS: paneData{
			id:    "audiobooks_perspective_pane",
			path:  "ui/pane/audiobooks.ui",
			setup: audiobookspane.Setup,
		},
	}
}
