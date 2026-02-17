package dto

import "github.com/norman6464/devsync/backend/internal/model"

// BookReviewListResponse は書籍レビュー一覧レスポンス。
type BookReviewListResponse struct {
	Reviews []model.BookReview `json:"reviews"`
	Total   int64              `json:"total"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
}
