package model

import (
	"time"

	"gorm.io/gorm"
)

type BookReview struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title     string         `json:"title" gorm:"not null;size:255"`
	Author    string         `json:"author" gorm:"not null;size:255"`
	CoverURL  string         `json:"cover_url" gorm:"size:500"`
	Rating    int            `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Content   string         `json:"content" gorm:"type:text"`
	LikeCount int            `json:"like_count" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type BookReviewLike struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_review"`
	BookReviewID uint      `json:"book_review_id" gorm:"not null;uniqueIndex:idx_user_review"`
	CreatedAt    time.Time `json:"created_at"`
}
