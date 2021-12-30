package pane

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	log "github.com/sirupsen/logrus"
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
func Add(idx m3uetcpb.Perspective, nb *gtk.Notebook, signals *map[string]interface{}) (err error) {
	log.WithField("idx", idx.String()).
		Info("Adding perspective to notebook")

	data, ok := paneList[idx]
	if !ok {
		err = fmt.Errorf("Unsupported perspective")
		return
	}

	if err = builder.AddFromFile(data.path); err != nil {
		err = fmt.Errorf("Unable to add file %v to builder: %w", data.path, err)
		return
	}

	newPane, err := builder.GetPane(data.id)
	if err != nil {
		err = fmt.Errorf("Unable to create %v pane: %v", idx, err)
		return
	}

	/*
		b, err := gtk.BuilderNewFromFile(data.path)
		if err != nil {
			err = fmt.Errorf("Unable to create builder for %v: %w", idx, err)
			return
		}

		obj, err := b.GetObject(data.id)
		if err != nil {
			err = fmt.Errorf("Unable to get %v: %v", data.id, err)
			return
		}

		newPane, ok := obj.(*gtk.Paned)
		if !ok {
			err = fmt.Errorf("Unable to create pane: %v", err)
			return
		}
	*/

	label, err := gtk.LabelNew(strings.Title(idx.String()))
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
			path:  "data/ui/pane/music.ui",
			setup: setupMusic,
		},
		m3uetcpb.Perspective_RADIO: paneData{
			id:    "radio_perspective_pane",
			path:  "data/ui/pane/radio.ui",
			setup: setupRadio,
		},
		m3uetcpb.Perspective_PODCASTS: paneData{
			id:    "podcasts_perspective_pane",
			path:  "data/ui/pane/podcasts.ui",
			setup: setupPodcasts,
		},
		m3uetcpb.Perspective_AUDIOBOOKS: paneData{
			id:    "audiobooks_perspective_pane",
			path:  "data/ui/pane/audiobooks.ui",
			setup: setupPodcasts,
		},
	}
}
