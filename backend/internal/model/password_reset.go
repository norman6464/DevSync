package model

import "time"

// PasswordResetToken はパスワードリセット用のトークンを表す。
// トークンは1時間有効で、使用済みフラグで再利用を防止する。
type PasswordResetToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"` // ランダム生成された64文字のHex文字列
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`        // トークンの有効期限
	Used      bool      `json:"used" gorm:"default:false"`         // 使用済みフラグ
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired はトークンが期限切れかどうかを判定する。
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsValid はトークンが有効（未使用かつ期限内）かどうかを判定する。
func (t *PasswordResetToken) IsValid() bool {
	return !t.Used && !t.IsExpired()
}
