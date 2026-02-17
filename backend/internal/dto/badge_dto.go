package dto

// NotifyBadgeEarnedRequest はバッジ獲得通知リクエスト。
type NotifyBadgeEarnedRequest struct {
	BadgeID string `json:"badge_id" binding:"required" validate:"required"`
}
