package pane

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
)

type paneData struct {
	id    string
	path  string
	setup func(signals *map[string]interface{}) error
}

type paneMap map[models.PerspectiveIndex]paneData

var (
	paneList paneMap
)

func Add(idx models.PerspectiveIndex, nb *gtk.Notebook, signals *map[string]interface{}) (err error) {
	data, ok := paneList[idx]
	if !ok {
		err = fmt.Errorf("Unsupported perspective")
		return
	}

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

	label, err := gtk.LabelNew(idx.String())
	if err != nil {
		err = fmt.Errorf("Unable to create %v label: %v", idx, err)
		return
	}

	if err = data.setup(signals); err != nil {
		return
	}

	nb.AppendPage(newPane, label)
	return
}

func init() {
	paneList = paneMap{
		models.MusicPerspective: paneData{
			id:    "music_perspective_pane",
			path:  "data/ui/pane/music.ui",
			setup: setupMusic,
		},
		models.RadioPerspective: paneData{
			id:    "radio_perspective_pane",
			path:  "data/ui/pane/radio.ui",
			setup: setupRadio,
		},
		models.PodcastsPerspective: paneData{
			id:    "podcasts_perspective_pane",
			path:  "data/ui/pane/podcasts.ui",
			setup: setupPodcasts,
		},
		models.AudiobooksPerspective: paneData{
			id:    "audiobooks_perspective_pane",
			path:  "data/ui/pane/audiobooks.ui",
			setup: setupPodcasts,
		},
	}
}
