package model

import (
	"time"

	"gorm.io/gorm"
)

type BookReview struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title     string         `json:"title" gorm:"not null;size:300"`
	Author    string         `json:"author" gorm:"size:200"`
	ISBN      string         `json:"isbn" gorm:"size:20"`
	Rating    int            `json:"rating" gorm:"not null"` // 1-5
	Review    string         `json:"review" gorm:"type:text"`
	ImageURL  string         `json:"image_url" gorm:"size:500"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
