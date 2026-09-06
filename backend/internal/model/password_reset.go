package model

import "time"

// PasswordResetToken はパスワードリセット用のトークンを表す。
// トークンは1時間有効で、使用済みフラグで再利用を防止する。
type PasswordResetToken struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user"`
	Token     string    `json:"token"`      // ランダム生成された64文字のHex文字列
	ExpiresAt time.Time `json:"expires_at"` // トークンの有効期限
	Used      bool      `json:"used"`       // 使用済みフラグ
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
