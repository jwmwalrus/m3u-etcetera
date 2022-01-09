package builder

import (
	"github.com/gotk3/gotk3/gtk"
)

var (
	app *gtk.Builder
)

func Setup(b *gtk.Builder) {
	app = b
}

func AddFromFile(path string) error {
	return app.AddFromFile(path)
}
