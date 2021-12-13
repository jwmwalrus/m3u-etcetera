package models

type PerspectiveIndex int

const (
	MusicPerspective PerspectiveIndex = iota
	RadioPerspective
	PodcastsPerspective
	AudiobooksPerspective
)

const DefaultPerspective = MusicPerspective

func PerspectiveIndexStrings() []string {
	return []string{"Music", "Radio", "Podcasts", "Audiobooks"}
}

func (idx PerspectiveIndex) String() string {
	return PerspectiveIndexStrings()[idx]
}

func (idx PerspectiveIndex) Activate() (err error) {
	s := []Perspective{}
	if err = db.Where("active = 1 OR idx = ?", int(idx)).Find(&s).Error; err != nil {
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
	return

}
func (idx PerspectiveIndex) Get() (p Perspective) {
	db.Where("idx = ?", int(idx)).First(&p)
	return
}

type Perspective struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Idx         int    `json:"idx" gorm:"uniqueIndex:unique_idx_perspective_index,not null"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	CreatedAt   int64  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   int64  `json:"updatedAt" gorm:"autoUpdateTime"`
}

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

func GetActivePerspectiveName() string {
	return GetActivePerspectiveIndex().String()
}
