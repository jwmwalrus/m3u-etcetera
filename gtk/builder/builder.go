package builder

import (
	"embed"
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
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
	_, err = app.AddFromString(string(bv))
	return err
}

// GetApplicationWindow returns the main application window.
func GetApplicationWindow() (window *gtk.ApplicationWindow, err error) {
	obj := app.GetObject("window")
	if obj == nil {
		rtc.Fatal("Unable to get window")
	}
	window, ok := obj.Cast().(*gtk.ApplicationWindow)
	if !ok {
		rtc.Fatal("Unable to create window", "error", err)
	}
	return
}

// PixbufNewFromFile creates a pixbuf from the given file path.
func PixbufNewFromFile(path string) (*gdkpixbuf.Pixbuf, error) {
	bv, err := data.ReadFile(path)
	if err != nil {
		return nil, err
	}

	loader := gdkpixbuf.NewPixbufLoader()
	err = loader.Write(bv)
	if err != nil {
		return nil, err
	}
	defer loader.Close()

	return loader.Pixbuf(), nil
}

// Setup -.
func Setup(fs *embed.FS) (b *gtk.Builder, err error) {
	const file = "ui/appwindow.ui"

	data = fs
	var bv []byte
	if bv, err = data.ReadFile(file); err != nil {
		return
	}
	app = gtk.NewBuilderFromString(string(bv))
	if app == nil {
		err = fmt.Errorf("failed to create builder from file %s", file)
		return
	}
	b = app
	return
}
