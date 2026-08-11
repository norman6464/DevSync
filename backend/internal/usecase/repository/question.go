package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// QuestionRepository は Q&A 質問の永続化に対する、usecase 側が要求する契約。
// 質問の CRUD・検索・投票・ブックマークを提供する。
type QuestionRepository interface {
	Create(ctx context.Context, question *model.Question) error
	// FindByID は指定 ID の質問を返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.Question, error)
	Update(ctx context.Context, question *model.Question) error
	Delete(ctx context.Context, id uint) error

	// 一覧・検索
	FindAll(ctx context.Context, limit, offset int, tag, sort string) ([]model.Question, int64, error)
	Search(ctx context.Context, query string, limit, offset int) ([]model.Question, int64, error)
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error)
	FindSolved(ctx context.Context, limit, offset int) ([]model.Question, int64, error)
	FindUnanswered(ctx context.Context, limit, offset int) ([]model.Question, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)

	// 投票
	Vote(ctx context.Context, userID, questionID uint, value int) error
	RemoveVote(ctx context.Context, userID, questionID uint) error
	// GetUserVote は投票値を返す。未投票の場合は 0 を返す。
	GetUserVote(ctx context.Context, userID, questionID uint) (int, error)

	// ブックマーク
	Bookmark(ctx context.Context, userID, questionID uint) error
	Unbookmark(ctx context.Context, userID, questionID uint) error
	HasBookmarked(ctx context.Context, userID, questionID uint) (bool, error)
	FindBookmarkedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error)
}
