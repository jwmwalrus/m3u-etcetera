package seeds

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

func SeedPlaybar(tx *gorm.DB) (err error) {
	create := func(idx models.PerspectiveIndex) error {
		p := models.Perspective{}
		if err := tx.Where("idx = ?", int(idx)).First(&p).Error; err != nil {
			return err
		}
		bar := models.Playbar{PerspectiveID: p.ID}
		if err := tx.Where(&bar).FirstOrCreate(&models.Playbar{}).Error; err != nil {
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
