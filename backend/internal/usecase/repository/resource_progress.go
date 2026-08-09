package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ResourceProgressRepository は学習リソース進捗の永続化に対する、usecase 側が要求する契約。
type ResourceProgressRepository interface {
	Upsert(ctx context.Context, progress *model.ResourceProgress) error
	FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceProgress, error)
	FindByUserID(ctx context.Context, userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error)
}

// リソース存在確認は resource_review で定義した最小 port LearningResourceReader を再利用する。
