package models

import (
	"log/slog"

	"github.com/jwmwalrus/m3u-etcetera/api/m3uetcpb"
	"github.com/jwmwalrus/m3u-etcetera/internal/subscription"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// PerspectiveIndex defines the perpective index.
type PerspectiveIndex int

const (
	// MusicPerspective -.
	MusicPerspective PerspectiveIndex = iota

	// RadioPerspective -.
	RadioPerspective

	// PodcastsPerspective -.
	PodcastsPerspective

	// AudiobooksPerspective -.
	AudiobooksPerspective
)

// DefaultPerspective declares the default perspective.
const DefaultPerspective = MusicPerspective

// PerspectiveIndexList returns the list of perspectives.
func PerspectiveIndexList() []PerspectiveIndex {
	return []PerspectiveIndex{
		MusicPerspective,
		RadioPerspective,
		PodcastsPerspective,
		AudiobooksPerspective,
	}
}

// PerspectiveIndexStrings Returns the string list of perspectives.
func PerspectiveIndexStrings() []string {
	return []string{"Music", "Radio", "Podcasts", "Audiobooks"}
}

func (idx PerspectiveIndex) String() string {
	return PerspectiveIndexStrings()[idx]
}

func (idx PerspectiveIndex) Description() string {
	return []string{"The Music Perspective",
		"The Radio Perspective",
		"The Podcasts Perspective",
		"The Audiobooks Perspective",
	}[idx]
}

// Activate sets the given perspective as active.
func (idx PerspectiveIndex) Activate() (err error) {
	slog.Info("Activating perspective", "idx", idx)

	s := []Perspective{}
	err = db.Where("active = 1 OR idx = ?", int(idx)).Find(&s).Error
	if err != nil {
		return
	}

	for i := 0; i < len(s); i++ {
		if s[i].Idx == int(idx) {
			s[i].Active = true
			continue
		}
		s[i].Active = false
	}

	err = db.Where("id > 0").Save(&s).Error
	if err == nil {
		subscription.Broadcast(subscription.ToPerspectiveEvent)
	}
	return

}

// Get returns the database row for the given index.
func (idx PerspectiveIndex) Get() (p Perspective) {
	db.Where("idx = ?", int(idx)).First(&p)
	return
}

// Perspective defines a perspective.
type Perspective struct {
	Model
	Idx         int    `json:"idx" gorm:"uniqueIndex:unique_idx_perspective_idx,not null"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

func (p *Perspective) Read(id int64) error {
	return p.ReadTx(db, id)
}

func (p *Perspective) ReadTx(tx *gorm.DB, id int64) error {
	return tx.
		First(p, id).
		Error
}

// GetActivePerspectiveIndex returns the index for the active perspective.
func GetActivePerspectiveIndex() (idx PerspectiveIndex) {
	var err error
	p := Perspective{}
	if err = db.Where("active = 1").First(&p).Error; err != nil {
		_ = DefaultPerspective.Activate()
		idx = DefaultPerspective
		return
	}
	idx = PerspectiveIndex(p.Idx)
	return
}

// GetActivePerspectiveName -.
func GetActivePerspectiveName() string {
	return GetActivePerspectiveIndex().String()
}

// PerspectiveDigest defines the perspective digest/summary.
// Implements the ProtoOut interface.
type PerspectiveDigest struct {
	Idx      PerspectiveIndex
	Duration int64
}

func (pd *PerspectiveDigest) ToProtobuf() proto.Message {
	return &m3uetcpb.PerspectiveDigest{
		Perspective: m3uetcpb.Perspective(pd.Idx),
		Duration:    pd.Duration,
	}
}
