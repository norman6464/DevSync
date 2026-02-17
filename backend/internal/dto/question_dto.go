package dto

import "github.com/norman6464/devsync/backend/internal/model"

// QuestionListResponse は質問一覧レスポンス。
type QuestionListResponse struct {
	Questions []model.Question `json:"questions"`
	Total     int64            `json:"total"`
	Limit     int              `json:"limit"`
	Offset    int              `json:"offset"`
}

// QuestionDetailResponse は質問詳細レスポンス（ユーザー投票情報付き）。
type QuestionDetailResponse struct {
	Question model.Question `json:"question"`
	UserVote int            `json:"user_vote"`
}
