package seeds

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

func seedPerspective(db *gorm.DB) (err error) {
	music := models.Perspective{
		Idx:         int(models.MusicPerspective),
		Description: "The Music Perspective",
		Active:      true,
	}
	if err = db.Create(&music).Error; err != nil {
		return
	}

	radio := models.Perspective{
		Idx:         int(models.RadioPerspective),
		Description: "The Radio Perspective",
	}
	if err = db.Create(&radio).Error; err != nil {
		return
	}

	pod := models.Perspective{
		Idx:         int(models.PodcastsPerspective),
		Description: "The Podcasts Perspective",
	}
	if err = db.Create(&pod).Error; err != nil {
		return
	}

	ab := models.Perspective{
		Idx:         int(models.AudiobooksPerspective),
		Description: "The Audiobooks Perspective",
	}
	if err = db.Create(&ab).Error; err != nil {
		return
	}

	return
}
