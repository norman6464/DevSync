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
	ID uint `json:"id" gorm:"primaryKey"`
	// インデックス名は PostgreSQL のスキーマ内で一意でなければならない。
	// 旧名 idx_collection_post は post_collection_items 側に同名の索引が残っており、
	// GORM が「既に存在する」と判定して作成をスキップするため、テーブル固有の名前にしている。
	CollectionID uint      `json:"collection_id" gorm:"not null;uniqueIndex:idx_bookmark_collection_post"`
	PostID       uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_bookmark_collection_post;index"`
	Post         Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt    time.Time `json:"created_at"`
}
