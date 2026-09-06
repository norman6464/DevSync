package model

import "time"

// NotificationSettings はユーザーの通知設定を表すモデル。
type NotificationSettings struct {
	ID     uint  `json:"id"`
	UserID uint  `json:"user_id"`
	User   *User `json:"user,omitempty"`

	// 通知タイプ別の設定
	EnableLikes    bool `json:"enable_likes"`    // いいね通知
	EnableComments bool `json:"enable_comments"` // コメント通知
	EnableFollows  bool `json:"enable_follows"`  // フォロー通知
	EnableMessages bool `json:"enable_messages"` // メッセージ通知
	EnableMentions bool `json:"enable_mentions"` // メンション通知

	// 通知配信方法
	EnableWebPush bool `json:"enable_web_push"` // Webプッシュ通知
	EnableEmail   bool `json:"enable_email"`    // メール通知
	EnableSound   bool `json:"enable_sound"`    // 通知音

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
