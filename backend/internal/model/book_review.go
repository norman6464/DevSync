package model

import (
	"time"

	"gorm.io/gorm"
)

// BookReview はユーザーが投稿した書籍レビューを表す。
// Rating は1〜5の整数で評価を表し、Review にレビュー本文を格納する。
type BookReview struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title     string         `json:"title" gorm:"not null;size:300"`          // 書籍タイトル
	Author    string         `json:"author" gorm:"size:200"`                  // 書籍の著者名
	ISBN      string         `json:"isbn" gorm:"size:20"`                     // ISBNコード
	Rating    int            `json:"rating" gorm:"not null"`                  // 評価（1〜5の整数）
	Review    string         `json:"review" gorm:"type:text"`                 // レビュー本文
	ImageURL   string         `json:"image_url" gorm:"size:500"`               // 書籍カバー画像URL
	IsArchived bool           `json:"is_archived" gorm:"default:false"`        // アーカイブ済みフラグ
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 論理削除用
}
