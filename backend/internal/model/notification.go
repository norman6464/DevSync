package model

import "time"

// NotificationType は通知の種別を表すカスタム文字列型。
type NotificationType string

// 通知種別の定数群。
const (
	NotificationTypePost    NotificationType = "post"    // 新規投稿
	NotificationTypeMessage NotificationType = "message" // 新規メッセージ
	NotificationTypeLike    NotificationType = "like"    // いいね
	NotificationTypeComment NotificationType = "comment" // コメント
	NotificationTypeFollow  NotificationType = "follow"  // フォロー
	NotificationTypeAnswer  NotificationType = "answer"  // Q&A回答
	NotificationTypeBadge   NotificationType = "badge"    // バッジ獲得
	NotificationTypeLevelUp NotificationType = "level_up" // レベルアップ
	NotificationTypeMention NotificationType = "mention"  // メンション
)

// Notification はユーザーへの通知を表す。
// ActorID は通知を発生させたユーザー、UserID は通知の受信者を示す。
type Notification struct {
	ID         uint             `json:"id" gorm:"primaryKey"`
	UserID     uint             `json:"user_id" gorm:"not null;index"`             // 通知の受信者
	User       User             `json:"user" gorm:"foreignKey:UserID"`
	Type       NotificationType `json:"type" gorm:"not null"`                      // 通知の種別
	ActorID    uint             `json:"actor_id" gorm:"not null"`                  // 通知を発生させたユーザー
	Actor      User             `json:"actor" gorm:"foreignKey:ActorID"`
	PostID     *uint            `json:"post_id" gorm:"index"`                      // 関連投稿（該当時のみ）
	Post       *Post            `json:"post,omitempty" gorm:"foreignKey:PostID"`
	QuestionID *uint            `json:"question_id" gorm:"index"`                  // 関連質問（該当時のみ）
	Question   *Question        `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	BadgeID    *string          `json:"badge_id,omitempty" gorm:"size:50"`         // 関連バッジ（該当時のみ）
	Read       bool             `json:"read" gorm:"default:false"`                 // 既読フラグ
	CreatedAt  time.Time        `json:"created_at"`
}
