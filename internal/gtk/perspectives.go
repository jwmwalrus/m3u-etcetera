package gtkui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/pane"
)

var (
	perspective     *gtk.ComboBoxText
	perspectiveList map[string]models.PerspectiveIndex
	notebook        *gtk.Notebook
)

// AddPerspectives sets the perspective panes
func AddPerspectives(b *gtk.Builder, signals *map[string]interface{}) (err error) {
	builder.SetupApp(b)

	perspective, err = builder.GetComboBoxText("perspective")
	if err != nil {
		return
	}

	(*signals)["on_perspective_changed"] = onPerspectiveChanged

	if err = setupPlayback(signals); err != nil {
		return
	}

	notebook, err = builder.GetNotebook("perspective_panes")
	if err != nil {
		return
	}

	for i := 0; i < notebook.GetNPages(); i++ {
		notebook.RemovePage(0)
	}

	for _, idx := range models.PerspectiveIndexList() {
		if err = pane.Add(idx, notebook, signals); err != nil {
			return
		}
	}

	glib.IdleAdd(updateCollection)
	go subscribeToPlayback()

	return
}

func onPerspectiveChanged(cbt *gtk.ComboBoxText) {
	text := cbt.GetActiveText()
	if idx, ok := perspectiveList[text]; ok {
		notebook.SetCurrentPage(int(idx))
	}
}

func init() {
	perspectiveList = map[string]models.PerspectiveIndex{
		models.MusicPerspective.String():      models.MusicPerspective,
		models.RadioPerspective.String():      models.RadioPerspective,
		models.PodcastsPerspective.String():   models.PodcastsPerspective,
		models.AudiobooksPerspective.String(): models.AudiobooksPerspective,
	}
}
