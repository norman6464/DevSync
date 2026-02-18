package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CreateBookReviewRequest は書籍レビュー作成のリクエストボディ。
type CreateBookReviewRequest struct {
	Title    string `json:"title" binding:"required,max=300" validate:"required,max=300"`
	Author   string `json:"author" binding:"max=200" validate:"max=200"`
	ISBN     string `json:"isbn" binding:"max=20" validate:"max=20"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5" validate:"required,min=1,max=5"`
	Review   string `json:"review"`
	ImageURL string `json:"image_url"`
}

// UpdateBookReviewRequest は書籍レビュー更新のリクエストボディ。
type UpdateBookReviewRequest struct {
	Title    string `json:"title" binding:"max=300" validate:"max=300"`
	Author   string `json:"author" binding:"max=200" validate:"max=200"`
	ISBN     string `json:"isbn" binding:"max=20" validate:"max=20"`
	Rating   *int   `json:"rating" binding:"omitempty,min=1,max=5" validate:"omitempty,min=1,max=5"`
	Review   string `json:"review"`
	ImageURL string `json:"image_url"`
}

// BookReviewListResponse は書籍レビュー一覧レスポンス。
type BookReviewListResponse struct {
	Reviews []model.BookReview `json:"reviews"`
	Total   int64              `json:"total"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
}
