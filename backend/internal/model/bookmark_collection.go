package model

import "time"

// BookmarkCollection はブックマークを整理するコレクション（フォルダ）を表す。
type BookmarkCollection struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"default:'blue'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BookmarkCollectionItem はコレクション内のブックマークアイテムを表す。
type BookmarkCollectionItem struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CollectionID uint      `json:"collection_id" gorm:"not null;uniqueIndex:idx_collection_post"`
	PostID       uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_collection_post;index"`
	Post         Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt    time.Time `json:"created_at"`
}
