package builder

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

// GetButton -
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

// GetCheckButton -
func GetCheckButton(id string) (c *gtk.CheckButton, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get check-button: %v", err)
		return
	}
	c, ok := obj.(*gtk.CheckButton)
	if !ok {
		err = fmt.Errorf("Unable to create check-button: %v", err)
		return
	}
	return
}

// GetComboBoxText -
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

// GetDialog -
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

// GetEntry -
func GetEntry(id string) (e *gtk.Entry, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get entry object: %v", err)
		return
	}

	e, ok := obj.(*gtk.Entry)
	if !ok {
		err = fmt.Errorf("Unable to create entry: %v", err)
		return
	}
	return
}

// GetFileChooserButton -
func GetFileChooserButton(id string) (fcb *gtk.FileChooserButton, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get file-chooser-button object: %v", err)
		return
	}

	fcb, ok := obj.(*gtk.FileChooserButton)
	if !ok {
		err = fmt.Errorf("Unable to create file-chooser-button: %v", err)
		return
	}
	return
}

// GetHeaderBar -
func GetHeaderBar(id string) (hb *gtk.HeaderBar, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get header-bar object: %v", err)
		return
	}

	hb, ok := obj.(*gtk.HeaderBar)
	if !ok {
		err = fmt.Errorf("Unable to create header-bar: %v", err)
		return
	}

	return
}

// GetImage -
func GetImage(id string) (im *gtk.Image, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get image object: %v", err)
		return
	}

	im, ok := obj.(*gtk.Image)
	if !ok {
		err = fmt.Errorf("Unable to create image: %v", err)
		return
	}
	return
}

// GetLabel -
func GetLabel(id string) (l *gtk.Label, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get label object: %v", err)
		return
	}

	l, ok := obj.(*gtk.Label)
	if !ok {
		err = fmt.Errorf("Unable to create label: %v", err)
		return
	}
	return
}

// GetListStore -
func GetListStore(id string) (s *gtk.ListStore, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get list-store object: %v", err)
		return
	}

	s, ok := obj.(*gtk.ListStore)
	if !ok {
		err = fmt.Errorf("Unable to create list-store: %v", err)
		return
	}
	return
}

// GetMenu -
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

// GetMenuItem -
func GetMenuItem(id string) (mi *gtk.MenuItem, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get menu item: %v", err)
		return
	}
	mi, ok := obj.(*gtk.MenuItem)
	if !ok {
		err = fmt.Errorf("Unable to create menu item: %v", err)
		return
	}

	return
}

// GetNotebook -
func GetNotebook(id string) (nb *gtk.Notebook, err error) {
	obj, err := app.GetObject(id)
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

// GetPane -
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

// GetPopoverMenu -
func GetPopoverMenu(id string) (pm *gtk.PopoverMenu, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get popover-menu object: %v", err)
		return
	}

	pm, ok := obj.(*gtk.PopoverMenu)
	if !ok {
		err = fmt.Errorf("Unable to create popover-menu: %v", err)
		return
	}
	return
}

// GetProgressBar -
func GetProgressBar(id string) (p *gtk.ProgressBar, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get progress-bar object: %v", err)
		return
	}

	p, ok := obj.(*gtk.ProgressBar)
	if !ok {
		err = fmt.Errorf("Unable to create progress bar: %v", err)
		return
	}
	return
}

// GetSpinButton -
func GetSpinButton(id string) (sb *gtk.SpinButton, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get spin-button object: %v", err)
		return
	}

	sb, ok := obj.(*gtk.SpinButton)
	if !ok {
		err = fmt.Errorf("Unable to create spin-button: %v", err)
		return
	}
	return
}

// GetStatusBar /
func GetStatusBar(id string) (sb *gtk.Statusbar, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get status-bar object: %v", err)
		return
	}

	sb, ok := obj.(*gtk.Statusbar)
	if !ok {
		err = fmt.Errorf("Unable to create status-bar: %v", err)
		return
	}
	return
}

// GetTextView -
func GetTextView(id string) (tv *gtk.TextView, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get text-view object: %v", err)
		return
	}

	tv, ok := obj.(*gtk.TextView)
	if !ok {
		err = fmt.Errorf("Unable to create text view: %v", err)
		return
	}
	return
}

// GetToggleButton -
func GetToggleButton(id string) (btn *gtk.ToggleButton, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get button: %v", err)
		return
	}
	btn, ok := obj.(*gtk.ToggleButton)
	if !ok {
		err = fmt.Errorf("Unable to create button: %v", err)
		return
	}
	return
}

// GetToggleToolButton -
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

// GetToolButton -
func GetToolButton(id string) (btn *gtk.ToolButton, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get tool-button object: %v", err)
		return
	}

	btn, ok := obj.(*gtk.ToolButton)
	if !ok {
		err = fmt.Errorf("Unable to create tool-button: %v", err)
		return
	}
	return
}

// GetTreeStore -
func GetTreeStore(id string) (s *gtk.TreeStore, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get tree-store object: %v", err)
		return
	}

	s, ok := obj.(*gtk.TreeStore)
	if !ok {
		err = fmt.Errorf("Unable to create tree-store: %v", err)
		return
	}
	return
}

// GetTreeView -
func GetTreeView(id string) (s *gtk.TreeView, err error) {
	obj, err := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get tree-view object: %v", err)
		return
	}

	s, ok := obj.(*gtk.TreeView)
	if !ok {
		err = fmt.Errorf("Unable to create tree-view: %v", err)
		return
	}
	return
}
