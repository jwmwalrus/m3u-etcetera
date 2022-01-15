package seeds

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

func seedPlaylist(db *gorm.DB) (err error) {
	p := models.Perspective{}
	if err = db.Where("idx = ?", int(models.DefaultPerspective)).First(&p).Error; err != nil {
		return
	}

	create := func(idx models.PlaylistGroupIndex) error {
		plg := models.PlaylistGroup{
			Idx:           int(idx),
			Name:          idx.String(),
			Hidden:        true,
			PerspectiveID: p.ID,
		}
		if err := db.Create(&plg).Error; err != nil {
			return err
		}
		return nil
	}

	if err = create(models.DefaultPlaylistGroup); err != nil {
		return
	}
	if err = create(models.TransientPlaylistGroup); err != nil {
		return
	}

	return
}
