package model

import "time"

// BookmarkCollection はブックマークを整理するコレクション（フォルダ）を表す。
type BookmarkCollection struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BookmarkCollectionItem はコレクション内のブックマークアイテムを表す。
type BookmarkCollectionItem struct {
	ID uint `json:"id"`
	// インデックス名は PostgreSQL のスキーマ内で一意でなければならない。
	// 旧名 idx_collection_post は post_collection_items 側に同名の索引が残っており、
	// GORM が「既に存在する」と判定して作成をスキップするため、テーブル固有の名前にしている。
	CollectionID uint      `json:"collection_id"`
	PostID       uint      `json:"post_id"`
	Post         Post      `json:"post,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
