package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningResourceRepository は学習リソースの永続化に対する、usecase 側が要求する契約。
type LearningResourceRepository interface {
	Create(ctx context.Context, resource *model.LearningResource) error
	Update(ctx context.Context, resource *model.LearningResource) error
	Delete(ctx context.Context, id uint) error
	// FindByID は指定 ID のリソースを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.LearningResource, error)

	// FindByUserID は指定ユーザーのリソースを返す。includePrivate が false なら公開分のみ。
	FindByUserID(ctx context.Context, userID uint, includePrivate bool, limit, offset int) ([]model.LearningResource, int64, error)
	FindPublic(ctx context.Context, limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error)
	FindByDifficulty(ctx context.Context, difficulty string, limit, offset int) ([]model.LearningResource, int64, error)
	// Search は公開リソースをタイトル・説明・タグで部分一致検索する。
	Search(ctx context.Context, query string, limit, offset int) ([]model.LearningResource, int64, error)
	FindSavedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningResource, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)

	// Like はいいねを追加し、リソースのいいね数を 1 増やす。
	Like(ctx context.Context, userID, resourceID uint) error
	// Unlike はいいねを取り消し、リソースのいいね数を 1 減らす（0 未満にはしない）。
	Unlike(ctx context.Context, userID, resourceID uint) error
	HasLiked(ctx context.Context, userID, resourceID uint) (bool, error)

	// Save は保存を追加し、リソースの保存数を 1 増やす。
	Save(ctx context.Context, userID, resourceID uint) error
	// Unsave は保存を取り消し、リソースの保存数を 1 減らす（0 未満にはしない）。
	Unsave(ctx context.Context, userID, resourceID uint) error
	HasSaved(ctx context.Context, userID, resourceID uint) (bool, error)
}
