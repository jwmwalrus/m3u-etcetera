package models

import (
	"github.com/jwmwalrus/bnp/pointers"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CollectionQuery Defines a collection query
type CollectionQuery struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	CreatedAt    int64      `json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt    int64      `json:"updatedAt" gorm:"autoUpdateTime:nano"`
	CollectionID int64      `json:"collectionId" gorm:"index:idx_collection_query_collection_id,not null"`
	QueryID      int64      `json:"queryId" gorm:"index:idx_collection_query_query_id,not null"`
	Collection   Collection `json:"collection" gorm:"foreignKey:CollectionID"`
	Query        Query      `json:"query" gorm:"foreignKey:QueryID"`
}

// Save implements the DataUpdater interface
func (cq *CollectionQuery) Save() error {
	return db.Save(cq).Error
}

// DeleteTx implements QueryBoundaryTx interface
func (cq *CollectionQuery) DeleteTx(tx *gorm.DB) error {
	return tx.Delete(cq).Error
}

// FindTracksTx implements QueryBoundaryTx interface
func (cq *CollectionQuery) FindTracksTx(tx *gorm.DB) []*Track {
	ts := []Track{}

	err := tx.
		Joins(
			"JOIN collection ON track.collection_id = collection.id AND track.collection_id = ?",
			cq.CollectionID,
		).
		Find(&ts).
		Error
	if err != nil {
		log.Error(err)
		return []*Track{}
	}
	return pointers.FromSlice(ts)
}

// GetQueryID implements QueryBoundaryID interface
func (cq *CollectionQuery) GetQueryID() int64 {
	return cq.QueryID
}

// SaveTx implements QueryBoundaryTx interface
func (cq *CollectionQuery) SaveTx(tx *gorm.DB) error {
	return tx.Save(cq).Error
}

// CollectionsToBoundaries adds forward support for CollectionQuery
func CollectionsToBoundaries(cts []*CollectionQuery) (qbs []QueryBoundaryTx) {
	for i := range cts {
		var x interface{} = cts[i]
		qbs = append(qbs, x.(QueryBoundaryTx))
	}
	return
}

// CreateCollectionQueries -
func CreateCollectionQueries(ids []int64) (cqs []*CollectionQuery) {
	cqs = []*CollectionQuery{}
	for _, id := range ids {
		c := CollectionQuery{CollectionID: id}
		cqs = append(cqs, &c)
	}
	return
}

// DeleteCollectionQueries deletes all the collection queries associated
// to the given query
func DeleteCollectionQueries(queryID int64) (err error) {
	cqs := []CollectionQuery{}
	if err = db.Where("query_id = ?", queryID).Find(&cqs).Error; err != nil {
		return
	}
	if len(cqs) < 1 {
		return
	}
	err = db.Where("id > 0").Delete(&cqs).Error
	return
}

// FilterCollectionQueryBoundaries -
func FilterCollectionQueryBoundaries(ids []int64) (qbs []QueryBoundaryID) {
	cqs := []CollectionQuery{}
	if err := db.Where("collection_id in ?", ids).Find(&cqs).Error; err != nil {
		return
	}

	for _, x := range cqs {
		var i interface{} = &x
		qbs = append(qbs, i.(QueryBoundaryID))
	}
	return
}

// GetApplicableCollectionQueries returns all the collections that can be
// applied to the given query
func GetApplicableCollectionQueries(qy *Query, ids ...int64) []*CollectionQuery {
	cqs := []CollectionQuery{}
	var err error

	if qy.ID > 0 {
		s := qy.GetCollections()
		if len(s) > 0 {
			err = db.
				Joins("JOIN collection on collection_query.collection_id = "+
					"collection.id and collection.hidden = 0 and collection.disabled = 0").
				Where("query_id = ?", qy.ID).
				Find(&cqs).
				Error
		} else {
			cs := []Collection{}
			err = db.Where("hidden = 0 and disabled = 0").Find(&cs).Error
			if err != nil {
				return []*CollectionQuery{}
			}
			for _, x := range cs {
				c := CollectionQuery{CollectionID: x.ID, QueryID: qy.ID}
				cqs = append(cqs, c)
			}
		}
	} else {
		cs := []Collection{}
		if len(ids) > 0 {
			err = db.Where("hidden = 0 and disabled = 0 and id in ?", ids).
				Find(&cs).
				Error
			if err != nil {
				return []*CollectionQuery{}
			}
			for _, x := range cs {
				c := CollectionQuery{CollectionID: x.ID}
				cqs = append(cqs, c)
			}
		} else {
			err = db.Where("hidden = 0 and disabled = 0").
				Find(&cs).
				Error
			if err != nil {
				return []*CollectionQuery{}
			}
			for _, x := range cs {
				c := CollectionQuery{CollectionID: x.ID}
				cqs = append(cqs, c)
			}
		}
	}
	if err != nil {
		log.Error(err)
		return []*CollectionQuery{}
	}

	return pointers.FromSlice(cqs)
}
