package model

import "time"

// NotificationSettings はユーザーの通知設定を表すモデル。
type NotificationSettings struct {
	ID     uint  `gorm:"primaryKey" json:"id"`
	UserID uint  `gorm:"uniqueIndex;not null" json:"user_id"`
	User   *User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// 通知タイプ別の設定
	EnableLikes    bool `gorm:"default:true" json:"enable_likes"`    // いいね通知
	EnableComments bool `gorm:"default:true" json:"enable_comments"` // コメント通知
	EnableFollows  bool `gorm:"default:true" json:"enable_follows"`  // フォロー通知
	EnableMessages bool `gorm:"default:true" json:"enable_messages"` // メッセージ通知
	EnableMentions bool `gorm:"default:true" json:"enable_mentions"` // メンション通知

	// 通知配信方法
	EnableWebPush bool `gorm:"default:true" json:"enable_web_push"` // Webプッシュ通知
	EnableEmail   bool `gorm:"default:true" json:"enable_email"`    // メール通知
	EnableSound   bool `gorm:"default:true" json:"enable_sound"`    // 通知音

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
