package builder

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

func AddFromFile(path string) error {
	return app.AddFromFile(path)
}

func GetPane(id string) (p *gtk.Paned, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get pane: %v", err)
		return
	}

	p, ok := obj.(*gtk.Paned)
	if !ok {
		err = fmt.Errorf("Unable to create pane: %v", err)
		return
	}
	return
}
