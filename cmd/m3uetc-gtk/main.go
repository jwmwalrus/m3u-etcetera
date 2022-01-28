package main

import (
	"github.com/gotk3/gotk3/gtk"
	gtkui "github.com/jwmwalrus/m3u-etcetera/gtk"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

var (
	window *gtk.ApplicationWindow
	b      *gtk.Builder
)

func main() {
	var err error

	gtk.Init(nil)

	base.Load()

	if b, err = gtk.BuilderNewFromFile("data/ui/appwindow.ui"); err != nil {
		log.Fatalf("Unable to create builder: %v", err)
	}

	builder.Setup(b)

	if window, err = builder.GetApplicationWindow(); err != nil {
		log.Fatalf("Unable to obtaain the application window: %v", err)
	}

	window.Connect("destroy", func() {
		store.Unsubscribe()
		gtk.MainQuit()
	})

	signals := make(map[string]interface{})

	err = gtkui.Setup(window, &signals)
	onerror.Panic(err)

	builder.ConnectSignals(signals)

	store.Subscribe()

	window.ShowAll()

	gtk.Main()
}
