package builder

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v3"
)

// GetButton -.
func GetButton(id string) (btn *gtk.Button, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get button")
		return
	}
	btn, ok := obj.Cast().(*gtk.Button)
	if !ok {
		err = fmt.Errorf("Unable to create button")
		return
	}
	return
}

// GetCheckButton -.
func GetCheckButton(id string) (c *gtk.CheckButton, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get check-button")
		return
	}
	c, ok := obj.Cast().(*gtk.CheckButton)
	if !ok {
		err = fmt.Errorf("Unable to create check-button")
		return
	}
	return
}

// GetComboBoxText -.
func GetComboBoxText(id string) (cbt *gtk.ComboBoxText, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get combo-box-text")
		return
	}
	cbt, ok := obj.Cast().(*gtk.ComboBoxText)
	if !ok {
		err = fmt.Errorf("Unable to create combo-box-text")
		return
	}
	return
}

// GetDialog -.
func GetDialog(id string) (dlg *gtk.Dialog, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get dialog")
		return
	}
	dlg, ok := obj.Cast().(*gtk.Dialog)
	if !ok {
		err = fmt.Errorf("Unable to create dialog")
		return
	}
	return
}

// GetEntry -.
func GetEntry(id string) (e *gtk.Entry, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get entry object")
		return
	}

	e, ok := obj.Cast().(*gtk.Entry)
	if !ok {
		err = fmt.Errorf("Unable to create entry")
		return
	}
	return
}

// GetFileChooserButton -.
func GetFileChooserButton(id string) (fcb *gtk.FileChooserButton, err error) {
	obj := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get file-chooser-button object")
		return
	}

	fcb, ok := obj.Cast().(*gtk.FileChooserButton)
	if !ok {
		err = fmt.Errorf("Unable to create file-chooser-button")
		return
	}
	return
}

// GetHeaderBar -.
func GetHeaderBar(id string) (hb *gtk.HeaderBar, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get header-bar object")
		return
	}

	hb, ok := obj.Cast().(*gtk.HeaderBar)
	if !ok {
		err = fmt.Errorf("Unable to create header-bar")
		return
	}

	return
}

// GetImage -.
func GetImage(id string) (im *gtk.Image, err error) {
	obj := app.GetObject(id)
	if err != nil {
		err = fmt.Errorf("Unable to get image object: %v", err)
		return
	}

	im, ok := obj.Cast().(*gtk.Image)
	if !ok {
		err = fmt.Errorf("Unable to create image")
		return
	}
	return
}

// GetLabel -.
func GetLabel(id string) (l *gtk.Label, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get label object")
		return
	}

	l, ok := obj.Cast().(*gtk.Label)
	if !ok {
		err = fmt.Errorf("Unable to create label")
		return
	}
	return
}

// GetListStore -.
func GetListStore(id string) (s *gtk.ListStore, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get list-store object")
		return
	}

	s, ok := obj.Cast().(*gtk.ListStore)
	if !ok {
		err = fmt.Errorf("Unable to create list-store")
		return
	}
	return
}

// GetMenu -.
func GetMenu(id string) (m *gtk.Menu, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get menu")
		return
	}
	m, ok := obj.Cast().(*gtk.Menu)
	if !ok {
		err = fmt.Errorf("Unable to create menu")
		return
	}

	return
}

// GetMenuItem -.
func GetMenuItem(id string) (mi *gtk.MenuItem, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get menu item")
		return
	}
	mi, ok := obj.Cast().(*gtk.MenuItem)
	if !ok {
		err = fmt.Errorf("Unable to create menu item")
		return
	}

	return
}

// GetNotebook -.
func GetNotebook(id string) (nb *gtk.Notebook, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get notebook")
		return
	}
	nb, ok := obj.Cast().(*gtk.Notebook)
	if !ok {
		err = fmt.Errorf("Unable to create notebook")
		return
	}
	return
}

// GetPane -.
func GetPane(id string) (p *gtk.Paned, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get pane")
		return
	}

	p, ok := obj.Cast().(*gtk.Paned)
	if !ok {
		err = fmt.Errorf("Unable to create pane")
		return
	}
	return
}

// GetPopoverMenu -.
func GetPopoverMenu(id string) (pm *gtk.PopoverMenu, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get popover-menu object")
		return
	}

	pm, ok := obj.Cast().(*gtk.PopoverMenu)
	if !ok {
		err = fmt.Errorf("Unable to create popover-menu")
		return
	}
	return
}

// GetProgressBar -.
func GetProgressBar(id string) (p *gtk.ProgressBar, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get progress-bar object")
		return
	}

	p, ok := obj.Cast().(*gtk.ProgressBar)
	if !ok {
		err = fmt.Errorf("Unable to create progress bar")
		return
	}
	return
}

// GetSpinButton -.
func GetSpinButton(id string) (sb *gtk.SpinButton, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get spin-button object")
		return
	}

	sb, ok := obj.Cast().(*gtk.SpinButton)
	if !ok {
		err = fmt.Errorf("Unable to create spin-button")
		return
	}
	return
}

// GetStatusBar -.
func GetStatusBar(id string) (sb *gtk.Statusbar, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get status-bar object")
		return
	}

	sb, ok := obj.Cast().(*gtk.Statusbar)
	if !ok {
		err = fmt.Errorf("Unable to create status-bar")
		return
	}
	return
}

// GetTextView -.
func GetTextView(id string) (tv *gtk.TextView, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get text-view object")
		return
	}

	tv, ok := obj.Cast().(*gtk.TextView)
	if !ok {
		err = fmt.Errorf("Unable to create text view")
		return
	}
	return
}

// GetToggleButton -.
func GetToggleButton(id string) (btn *gtk.ToggleButton, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get button")
		return
	}
	btn, ok := obj.Cast().(*gtk.ToggleButton)
	if !ok {
		err = fmt.Errorf("Unable to create button")
		return
	}
	return
}

// GetToggleToolButton -.
func GetToggleToolButton(id string) (btn *gtk.ToggleToolButton, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get button")
		return
	}
	btn, ok := obj.Cast().(*gtk.ToggleToolButton)
	if !ok {
		err = fmt.Errorf("Unable to create button")
		return
	}
	return
}

// GetToolButton -.
func GetToolButton(id string) (btn *gtk.ToolButton, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get tool-button object")
		return
	}

	btn, ok := obj.Cast().(*gtk.ToolButton)
	if !ok {
		err = fmt.Errorf("Unable to create tool-button")
		return
	}
	return
}

// GetTreeStore -.
func GetTreeStore(id string) (s *gtk.TreeStore, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get tree-store object")
		return
	}

	s, ok := obj.Cast().(*gtk.TreeStore)
	if !ok {
		err = fmt.Errorf("Unable to create tree-store")
		return
	}
	return
}

// GetTreeView -.
func GetTreeView(id string) (s *gtk.TreeView, err error) {
	obj := app.GetObject(id)
	if obj == nil {
		err = fmt.Errorf("Unable to get tree-view object")
		return
	}

	s, ok := obj.Cast().(*gtk.TreeView)
	if !ok {
		err = fmt.Errorf("Unable to create tree-view")
		return
	}
	return
}
