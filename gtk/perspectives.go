package gtkui

import (
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/pane"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

var (
	perspectivesList []m3uetcpb.Perspective
	notebook         *gtk.Notebook
)

func addPerspectives(signals *map[string]interface{}) (err error) {
	log.Info("Adding perspectives")

	_, err = builder.GetComboBoxText("perspective")
	if err != nil {
		return
	}
	(*signals)["on_perspective_changed"] = onPerspectiveChanged

	notebook, err = builder.GetNotebook("perspective_panes")
	if err != nil {
		return
	}

	for i := 0; i < notebook.GetNPages(); i++ {
		notebook.RemovePage(0)
	}

	for _, v := range perspectivesList {
		err = pane.Add(v, notebook, signals)
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

	go func() {
		req := &m3uetcpb.SetActivePerspectiveRequest{
			Perspective: m3uetcpb.Perspective(m3uetcpb.Perspective_value[strings.ToUpper(text)]),
		}
		onerror.Log(store.SetActivePerspective(req))
	}()
}
