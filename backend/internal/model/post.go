package model

import "time"

type Post struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Title        string    `json:"title" gorm:"not null"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	LikeCount    int       `json:"like_count" gorm:"default:0"`
	CommentCount int       `json:"comment_count" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post_like"`
	PostID    uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post_like;index"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	PostID    uint      `json:"post_id" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
