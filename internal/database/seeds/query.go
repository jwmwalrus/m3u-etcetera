package seeds

import (
	"github.com/jwmwalrus/m3u-etcetera/internal/database/models"
	"gorm.io/gorm"
)

func SeedQuery(tx *gorm.DB) (err error) {
	create := func(idx models.QueryIndex) error {
		qy := models.Query{
			Idx:         int(idx),
			Name:        idx.String(),
			Description: idx.Description(),
		}
		if err := tx.Where(&qy).FirstOrCreate(&models.Query{}).Error; err != nil {
			return err
		}

		return nil
	}

	if err = create(models.HistoryQuery); err != nil {
		return
	}
	if err = create(models.TopTracksQuery); err != nil {
		return
	}
	if err = create(models.Gimme20RandomsQuery); err != nil {
		return
	}
	if err = create(models.Gimme50RandomsQuery); err != nil {
		return
	}
	if err = create(models.Gimme100RandomsQuery); err != nil {
		return
	}

	return
}
