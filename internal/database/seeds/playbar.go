package seeds

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

func seedPlaybar(db *gorm.DB) (err error) {
	create := func(idx models.PerspectiveIndex) error {
		p := models.Perspective{}
		if err := db.Where("idx = ?", int(idx)).First(&p).Error; err != nil {
			return err
		}
		bar := models.Playbar{PerspectiveID: p.ID}
		if err := db.Create(&bar).Error; err != nil {
			return err
		}

		return nil
	}

	if err = create(models.MusicPerspective); err != nil {
		return
	}
	if err = create(models.RadioPerspective); err != nil {
		return
	}
	if err = create(models.PodcastsPerspective); err != nil {
		return
	}
	if err = create(models.AudiobooksPerspective); err != nil {
		return
	}

	return
}
