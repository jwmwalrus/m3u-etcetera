package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	gtkui "github.com/jwmwalrus/m3u-etcetera/internal/gtk"
	"github.com/jwmwalrus/m3u-etcetera/internal/gtk/store"
	log "github.com/sirupsen/logrus"
)

var (
	window  *gtk.ApplicationWindow
	builder *gtk.Builder
)

func main() {
	var (
		err error
		obj glib.IObject
		ok  bool
	)

	gtk.Init(nil)

	base.Load()

	if builder, err = gtk.BuilderNewFromFile("data/ui/appwindow.ui"); err != nil {
		log.Fatalf("Unable to create builder: %v", err)
	}

	if obj, err = builder.GetObject("window"); err != nil {
		log.Fatalf("Unable to get window: %v", err)
	}
	if window, ok = obj.(*gtk.ApplicationWindow); !ok {
		log.Fatalf("Unable to create window: %v", err)
	}

	window.Connect("destroy", func() {
		// store.Unsubscribe()
		gtk.MainQuit()
	})

	signals := make(map[string]interface{})

	err = gtkui.AddPerspectives(builder, &signals)
	onerror.Panic(err)

	builder.ConnectSignals(signals)

	window.ShowAll()

	store.Subscribe()

	gtk.Main()
}
