package dto

import "github.com/norman6464/devsync/backend/internal/service"

// BadgesResponse はバッジ一覧レスポンス
type BadgesResponse struct {
	Badges []service.BadgeResult `json:"badges"`
}

// NotifyBadgeEarnedRequest はバッジ獲得通知リクエスト。
type NotifyBadgeEarnedRequest struct {
	BadgeID string `json:"badge_id" binding:"required" validate:"required"`
}
