package store

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
)

// GetActivePerspective returns the active perspective.
func GetActivePerspective() m3uetcpb.Perspective {
	// TODO: implement
	return m3uetcpb.Perspective_MUSIC
}

type perspectiveData struct {
	res *m3uetcpb.SubscribeToPerspectiveResponse

	uiSet bool
	combo *gtk.ComboBoxText

	mu sync.RWMutex
}

var (
	// PerspData perspective data.
	PerspData = &perspectiveData{}
)

func (pd *perspectiveData) SubscriptionID() string {
	pd.mu.RLock()
	defer pd.mu.RUnlock()

	return pd.res.SubscriptionId
}

func (pd *perspectiveData) SetPerspectiveUI() (err error) {
	pd.combo, err = builder.GetComboBoxText("perspective")
	if err != nil {
		return
	}

	pd.uiSet = true
	return
}

func (pd *perspectiveData) updateActivePerspective() bool {
	slog.Debug("Updating active perspective")

	pd.mu.RLock()
	defer pd.mu.RUnlock()

	active := pd.res.ActivePerspective

	id := strings.ToLower(active.String()) + "_perspective"
	if !pd.combo.SetActiveID(id) {
		slog.Error("Error setting active perspective")
	}

	return false
}

func (pd *perspectiveData) ProcessSubscriptionResponse(res *m3uetcpb.SubscribeToPerspectiveResponse) {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	pd.res = res
	glib.IdleAdd(pd.updateActivePerspective)
}
