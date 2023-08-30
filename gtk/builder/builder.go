package builder

import (
	"embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	rtc "github.com/jwmwalrus/rtcycler"
)

var (
	app  *gtk.Builder
	data *embed.FS
)

// AddFromFile - adds a new resource from the given embedded path.
func AddFromFile(path string) error {
	bv, err := data.ReadFile(path)
	if err != nil {
		return err
	}
	return app.AddFromString(string(bv))
}

// ConnectSignals connects the signals map.
func ConnectSignals(signals map[string]interface{}) {
	app.ConnectSignals(signals)
}

// GetApplicationWindow returns the main application window.
func GetApplicationWindow() (window *gtk.ApplicationWindow, err error) {
	obj, err := app.GetObject("window")
	if err != nil {
		rtc.Fatal("Unable to get window", "error", err)
	}
	window, ok := obj.(*gtk.ApplicationWindow)
	if !ok {
		rtc.Fatal("Unable to create window", "error", err)
	}
	return
}

// PixbufNewFromFile creates a pixbuf from the given file path.
func PixbufNewFromFile(path string) (*gdk.Pixbuf, error) {
	bv, err := data.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return gdk.PixbufNewFromDataOnly(bv)
}

// Setup -.
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
