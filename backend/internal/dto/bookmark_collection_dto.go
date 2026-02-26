package dto

import "github.com/norman6464/devsync/backend/internal/model"

// BookmarkCollectionRequest はブックマークコレクション作成・更新リクエスト。
type BookmarkCollectionRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Color       string `json:"color" binding:"omitempty,max=50"`
}

// BookmarkCollectionPostsResponse はコレクション内投稿一覧レスポンス。
type BookmarkCollectionPostsResponse struct {
	Posts []model.Post `json:"posts"`
	Total int64        `json:"total"`
}
