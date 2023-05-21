package seeds

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

func SeedCollection(tx *gorm.DB) (err error) {
	create := func(idx models.CollectionIndex) error {
		coll := models.Collection{
			Idx:      int(idx),
			Name:     idx.String(),
			Location: idx.String(),
			Hidden:   true,
			// Scanned:       100,
			PerspectiveID: 1,
		}

		err := tx.Where(&coll).FirstOrCreate(&models.Collection{}).
			Error
		if err != nil {
			return err
		}
		return nil
	}

	if err = create(models.DefaultCollection); err != nil {
		return
	}
	if err = create(models.TransientCollection); err != nil {
		return
	}

	return
}
