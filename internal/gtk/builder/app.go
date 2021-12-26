package builder

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

var (
	app *gtk.Builder
)

func SetupApp(b *gtk.Builder) {
	app = b
}

func GetComboBoxText(id string) (cbt *gtk.ComboBoxText, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get combo-box-text: %v", err)
		return
	}
	cbt, ok := obj.(*gtk.ComboBoxText)
	if !ok {
		err = fmt.Errorf("Unable to create combo-box-text: %v", err)
		return
	}
	return
}

func GetNotebook(id string) (nb *gtk.Notebook, err error) {
	obj, err := app.GetObject("perspective_panes")
	if err != nil {
		err = fmt.Errorf("Unable to get notebook: %v", err)
		return
	}
	nb, ok := obj.(*gtk.Notebook)
	if !ok {
		err = fmt.Errorf("Unable to create notebook: %v", err)
		return
	}
	return
}

func GetTextView(id string) (tv *gtk.TextView, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get text-view object: %w", err)
		return
	}

	tv, ok := obj.(*gtk.TextView)
	if !ok {
		err = fmt.Errorf("Unable to create text view: %w", err)
		return
	}
	return
}

func GetToolButton(id string) (btn *gtk.ToolButton, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get tool-button object: %w", err)
		return
	}

	btn, ok := obj.(*gtk.ToolButton)
	if !ok {
		err = fmt.Errorf("Unable to create tool-button: %w", err)
		return
	}
	return
}

func SetTextView(id, val string) (err error) {
	tv, err := GetTextView(id)
	if err != nil {
		return
	}

	buf, err := tv.GetBuffer()
	if err != nil {
		err = fmt.Errorf("Unable to get buffer from text view: %w", err)
		return
	}

	buf.SetText(val)
	return
}
