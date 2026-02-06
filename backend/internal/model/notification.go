package model

import "time"

type NotificationType string

const (
	NotificationTypePost    NotificationType = "post"
	NotificationTypeMessage NotificationType = "message"
)

type Notification struct {
	ID        uint             `json:"id" gorm:"primaryKey"`
	UserID    uint             `json:"user_id" gorm:"not null;index"`
	User      User             `json:"user" gorm:"foreignKey:UserID"`
	Type      NotificationType `json:"type" gorm:"not null"`
	ActorID   uint             `json:"actor_id" gorm:"not null"`
	Actor     User             `json:"actor" gorm:"foreignKey:ActorID"`
	PostID    *uint            `json:"post_id" gorm:"index"`
	Post      *Post            `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Read      bool             `json:"read" gorm:"default:false"`
	CreatedAt time.Time        `json:"created_at"`
}
