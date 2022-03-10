package builder

import (
	"embed"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var (
	app  *gtk.Builder
	data *embed.FS
)

func AddFromFile(path string) error {
	bv, err := data.ReadFile(path)
	if err != nil {
		return err
	}
	return app.AddFromString(string(bv))
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

func PixbufNewFromFile(path string) (*gdk.Pixbuf, error) {
	bv, err := data.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return gdk.PixbufNewFromDataOnly(bv)
}

func Setup(fs *embed.FS) (b *gtk.Builder, err error) {
	data = fs
	var bv []byte
	if bv, err = data.ReadFile("ui/appwindow.ui"); err != nil {
		return
	}
	app, err = gtk.BuilderNewFromString(string(bv))
	if err == nil {
		b = app
	}
	return
}
