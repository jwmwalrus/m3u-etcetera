package seeds

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

func SeedPerspective(tx *gorm.DB) (err error) {
	create := func(idx models.PerspectiveIndex) error {
		persp := models.Perspective{
			Idx:         int(idx),
			Description: idx.Description(),
		}
		err := tx.Where(&persp).FirstOrCreate(&models.Perspective{}).Error
		if err != nil {
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
