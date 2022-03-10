package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	gtkui "github.com/jwmwalrus/m3u-etcetera/gtk"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	"github.com/jwmwalrus/m3u-etcetera/gtk/store"
	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
)

var (
	appID     = "com.github.jwmwalrus." + base.AppInstance
	app       *gtk.Application
	window    *gtk.ApplicationWindow
	b         *gtk.Builder
	activated bool

	//go:embed ui images
	data embed.FS
)

func main() {
	var err error

	base.Load()
	app, err = gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatalf("Unable to create application: %v", err)
	}

	app.Connect("activate", func() {
		if activated {
			fmt.Printf("Primary instance already active\n")
			return
		}

		log.Infof("Activating primary instance: %v", appID)

		activated = true
		if b, err = builder.Setup(&data); err != nil {
			log.Fatalf("Unable to create builder: %v", err)
		}

		if window, err = builder.GetApplicationWindow(); err != nil {
			log.Fatalf("Unable to obtaain the application window: %v", err)
		}

		icon, err := builder.PixbufNewFromFile("images/m3u-etcetera.svg")
		if err != nil {
			log.Error(err)
		} else {
			window.SetIcon(icon)
		}

		window.SetApplication(app)

		window.Connect("destroy", func() {
			store.Unsubscribe()
			fmt.Printf("\nBye %v from %v\n", base.OS, base.AppInstance)
			app.Quit()
		})

		signals := make(map[string]interface{})

		err = gtkui.Setup(window, &signals)
		onerror.Panic(err)

		builder.ConnectSignals(signals)

		store.Subscribe()

		window.ShowAll()
	})

	os.Exit(app.Run([]string{}))
}
