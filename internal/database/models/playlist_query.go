package models

// PlaylistQuery Defines a playlist query
type PlaylistQuery struct {
	ID         int64    `json:"id" gorm:"primaryKey"`
	CreatedAt  int64    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  int64    `json:"updatedAt" gorm:"autoUpdateTime"`
	PlaylistID int64    `json:"playlistId" gorm:"index:idx_playlist_query_playlist_id,not null"`
	QueryID    int64    `json:"queryId" gorm:"index:idx_playlist_query_query_id,not null"`
	Playlist   Playlist `json:"playlist" gorm:"foreignKey:PlaylistID"`
	Query      Query    `json:"query" gorm:"foreignKey:QueryID"`
}

// Delete implements the DataDeleter interface
func (pqy *PlaylistQuery) Delete() error {
	return db.Delete(pqy).Error
}
