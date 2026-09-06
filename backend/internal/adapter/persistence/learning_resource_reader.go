package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// learningResourceReader は [repository.LearningResourceReader] の sqlc(pgx) 実装。
// レビューの存在確認に必要な読み取りだけを担う最小 adapter。
type learningResourceReader struct {
	q *sqlcgen.Queries
}

// NewLearningResourceReader は LearningResourceReader の sqlc(pgx) 実装を返す。
func NewLearningResourceReader(q *sqlcgen.Queries) repository.LearningResourceReader {
	return &learningResourceReader{q: q}
}

// コンパイル時に port を満たすことを保証する。
var _ repository.LearningResourceReader = (*learningResourceReader)(nil)

// FindByID は指定 ID の学習リソースを取得する。不在の場合は (nil, nil) を返す。
// 存在確認のみに使うため関連の Preload は行わない。
func (r *learningResourceReader) FindByID(ctx context.Context, id uint) (*model.LearningResource, error) {
	row, err := r.q.GetLearningResourceByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &model.LearningResource{
		ID:          uint(row.ID),
		UserID:      uint(row.UserID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		URL:         fromStringPtr(row.Url),
		Category:    model.ResourceCategory(row.Category),
		Difficulty:  model.ResourceDifficulty(fromStringPtr(row.Difficulty)),
		Tags:        fromStringPtr(row.Tags),
		ImageURL:    fromStringPtr(row.ImageUrl),
		IsPublic:    row.IsPublic,
		LikeCount:   int(fromInt64PtrValue(row.LikeCount)),
		SaveCount:   int(fromInt64PtrValue(row.SaveCount)),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}, nil
}
