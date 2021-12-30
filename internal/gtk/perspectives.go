package gtkui

import (
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/pane"
	log "github.com/sirupsen/logrus"
)

var (
	perspective      *gtk.ComboBoxText
	perspectivesList []m3uetcpb.Perspective
	notebook         *gtk.Notebook
)

// AddPerspectives sets the perspective panes
func AddPerspectives(b *gtk.Builder, signals *map[string]interface{}) (err error) {
	log.Info("Adding perspectives")

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

	for _, v := range perspectivesList {
		err = pane.Add(m3uetcpb.Perspective(v), notebook, signals)
		if err != nil {
			return
		}
	}

	return
}

func onPerspectiveChanged(cbt *gtk.ComboBoxText) {
	text := cbt.GetActiveText()
	log.WithField("activeText", text).
		Info("Perspective changed")

	text = strings.ToUpper(cbt.GetActiveText())
	if idx, ok := m3uetcpb.Perspective_value[text]; ok {
		notebook.SetCurrentPage(int(idx))
	}
}

func init() {
	perspectivesList = []m3uetcpb.Perspective{
		m3uetcpb.Perspective_MUSIC,
		m3uetcpb.Perspective_RADIO,
		m3uetcpb.Perspective_PODCASTS,
		m3uetcpb.Perspective_AUDIOBOOKS,
	}

}
