package builder

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func GetButton(id string) (btn *gtk.Button, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get button: %v", err)
		return
	}
	btn, ok := obj.(*gtk.Button)
	if !ok {
		err = fmt.Errorf("Unable to create button: %v", err)
		return
	}
	return
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

func GetDialog(id string) (dlg *gtk.Dialog, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get dialog: %v", err)
		return
	}
	dlg, ok := obj.(*gtk.Dialog)
	if !ok {
		err = fmt.Errorf("Unable to create dialog: %v", err)
		return
	}
	return
}

func GetLabel(id string) (l *gtk.Label, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get label object: %w", err)
		return
	}

	l, ok := obj.(*gtk.Label)
	if !ok {
		err = fmt.Errorf("Unable to create label: %w", err)
		return
	}
	return
}

func GetListStore(id string) (s *gtk.ListStore, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get list-store object: %w", err)
		return
	}

	s, ok := obj.(*gtk.ListStore)
	if !ok {
		err = fmt.Errorf("Unable to create list-store: %w", err)
		return
	}
	return
}

func GetMenu(id string) (m *gtk.Menu, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get menu: %v", err)
		return
	}
	m, ok := obj.(*gtk.Menu)
	if !ok {
		err = fmt.Errorf("Unable to create menu: %v", err)
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

func GetPopoverMenu(id string) (pm *gtk.PopoverMenu, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get popover-menu object: %w", err)
		return
	}

	pm, ok := obj.(*gtk.PopoverMenu)
	if !ok {
		err = fmt.Errorf("Unable to create popover-menu: %w", err)
		return
	}
	return
}

func GetProgressBar(id string) (p *gtk.ProgressBar, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get progress-bar object: %w", err)
		return
	}

	p, ok := obj.(*gtk.ProgressBar)
	if !ok {
		err = fmt.Errorf("Unable to create progress bar: %w", err)
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

func GetToggleToolButton(id string) (btn *gtk.ToggleToolButton, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get button: %v", err)
		return
	}
	btn, ok := obj.(*gtk.ToggleToolButton)
	if !ok {
		err = fmt.Errorf("Unable to create button: %v", err)
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

func GetTreeStore(id string) (s *gtk.TreeStore, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get tree-store object: %w", err)
		return
	}

	s, ok := obj.(*gtk.TreeStore)
	if !ok {
		err = fmt.Errorf("Unable to create tree-store: %w", err)
		return
	}
	return
}

func GetTreeView(id string) (s *gtk.TreeView, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get tree-view object: %w", err)
		return
	}

	s, ok := obj.(*gtk.TreeView)
	if !ok {
		err = fmt.Errorf("Unable to create tree-view: %w", err)
		return
	}
	return
}

func GetWindow() (window *gtk.ApplicationWindow, err error) {
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
