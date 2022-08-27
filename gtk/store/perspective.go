package store

import (
	"strings"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/gtk/builder"
	log "github.com/sirupsen/logrus"
)

// GetActivePerspective returns the active perspective
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
	// PerspData perspective data
	PerspData = &perspectiveData{}
)

func (pd *perspectiveData) GetSubscriptionID() string {
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
	log.Debug("Updating active perspective")

	pd.mu.RLock()
	defer pd.mu.RUnlock()

	active := pd.res.ActivePerspective

	id := strings.ToLower(active.String()) + "_perspective"
	if !pd.combo.SetActiveID(id) {
		log.Error("Error setting active perspectie")
	}

	return false
}

func (pd *perspectiveData) ProcessSubscriptionResponse(res *m3uetcpb.SubscribeToPerspectiveResponse) {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	pd.res = res
	glib.IdleAdd(pd.updateActivePerspective)
}
