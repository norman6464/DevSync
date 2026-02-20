package dto

import "github.com/norman6464/devsync/backend/internal/model"

// FollowListResponse はフォロワー/フォロー中一覧レスポンス（ページネーション付き）。
type FollowListResponse struct {
	Users  []model.User `json:"users"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}
