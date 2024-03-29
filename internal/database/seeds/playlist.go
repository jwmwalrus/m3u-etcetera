package seeds

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

func SeedPlaylist(tx *gorm.DB) (err error) {
	create := func(idx models.PlaylistGroupIndex, p models.PerspectiveIndex) error {
		plg := models.PlaylistGroup{
			Idx:           int(idx),
			Name:          idx.String(),
			Hidden:        true,
			PerspectiveID: p.Get().ID,
		}
		if err := tx.Where(&plg).FirstOrCreate(&models.PlaylistGroup{}).Error; err != nil {
			return err
		}
		return nil
	}

	err = create(models.MusicPlaylistGroup, models.MusicPerspective)
	if err != nil {
		return
	}

	err = create(models.RadioPlaylistGroup, models.RadioPerspective)
	if err != nil {
		return
	}

	err = create(models.PodcastsPlaylistGroup, models.PodcastsPerspective)
	if err != nil {
		return
	}

	err = create(models.AudiobooksPlaylistGroup, models.AudiobooksPerspective)
	if err != nil {
		return
	}

	return
}
