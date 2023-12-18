package models

import "gorm.io/gorm"

// PlaylistQuery Defines a playlist query.
type PlaylistQuery struct {
	Model
	PlaylistID int64    `json:"playlistId" gorm:"index:idx_playlist_query_playlist_id,not null"`
	QueryID    int64    `json:"queryId" gorm:"index:idx_playlist_query_query_id,not null"`
	Playlist   Playlist `json:"playlist" gorm:"foreignKey:PlaylistID"`
	Query      Query    `json:"query" gorm:"foreignKey:QueryID"`
}

func (pqy *PlaylistQuery) Delete() error {
	return pqy.DeleteTx(db)
}

func (pqy *PlaylistQuery) DeleteTx(tx *gorm.DB) error {
	return tx.Delete(pqy).Error
}
