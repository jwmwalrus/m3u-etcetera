package builder

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

var (
	app *gtk.Builder
)

func AddFromFile(path string) error {
	return app.AddFromFile(path)
}

func ConnectSignals(signals map[string]interface{}) {
	app.ConnectSignals(signals)
}

func GetApplicationWindow() (window *gtk.ApplicationWindow, err error) {
	obj, err := app.GetObject("window")
	if err != nil {
		log.Fatalf("Unable to get window: %v", err)
	}
	window, ok := obj.(*gtk.ApplicationWindow)
	if !ok {
		log.Fatalf("Unable to create window: %v", err)
	}
	return
}

func Setup(b *gtk.Builder) {
	app = b
}
