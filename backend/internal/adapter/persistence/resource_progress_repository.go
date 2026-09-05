package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// resourceProgressRepository は [repository.ResourceProgressRepository] の sqlc(pgx) 実装。
type resourceProgressRepository struct {
	q *sqlcgen.Queries
}

// NewResourceProgressRepository は ResourceProgressRepository の sqlc(pgx) 実装を返す。
func NewResourceProgressRepository(q *sqlcgen.Queries) repository.ResourceProgressRepository {
	return &resourceProgressRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ResourceProgressRepository = (*resourceProgressRepository)(nil)

// toModelResourceProgress は sqlc の生成行を model.ResourceProgress へ変換する（Resource は含まない）。
func toModelResourceProgress(row sqlcgen.ResourceProgress) model.ResourceProgress {
	return model.ResourceProgress{
		ID:                uint(row.ID),
		UserID:            uint(row.UserID),
		ResourceID:        uint(row.ResourceID),
		Status:            model.ResourceProgressStatus(fromStringPtr(row.Status)),
		CompletionPercent: int(fromInt64PtrValue(row.CompletionPercent)),
		Note:              fromStringPtr(row.Note),
		StartedAt:         fromTimestamptz(row.StartedAt),
		CompletedAt:       fromTimestamptz(row.CompletedAt),
		CreatedAt:         timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:         timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Upsert は (user_id, resource_id) をキーに進捗を作成または更新する。
func (r *resourceProgressRepository) Upsert(ctx context.Context, progress *model.ResourceProgress) error {
	status := string(progress.Status)
	row, err := r.q.UpsertResourceProgress(ctx, sqlcgen.UpsertResourceProgressParams{
		UserID:            int64(progress.UserID),
		ResourceID:        int64(progress.ResourceID),
		Status:            &status,
		CompletionPercent: toInt64Ptr(progress.CompletionPercent),
		Note:              &progress.Note,
		StartedAt:         toTimestamptz(progress.StartedAt),
		CompletedAt:       toTimestamptz(progress.CompletedAt),
	})
	if err != nil {
		return err
	}
	*progress = toModelResourceProgress(row)
	return nil
}

// FindByUserAndResource は指定ユーザー・リソースの進捗を取得する。
func (r *resourceProgressRepository) FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceProgress, error) {
	row, err := r.q.GetResourceProgressByUserAndResource(ctx, sqlcgen.GetResourceProgressByUserAndResourceParams{
		UserID:     int64(userID),
		ResourceID: int64(resourceID),
	})
	if err != nil {
		return nil, err
	}
	progress := toModelResourceProgress(row)
	return &progress, nil
}

// FindByUserID は指定ユーザーの進捗一覧を取得する。status が空でなければ絞り込む。
func (r *resourceProgressRepository) FindByUserID(ctx context.Context, userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error) {
	var statusFilter *string
	if status != "" {
		statusFilter = &status
	}

	total, err := r.q.CountResourceProgressByUser(ctx, sqlcgen.CountResourceProgressByUserParams{
		UserID: int64(userID),
		Status: statusFilter,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListResourceProgressByUser(ctx, sqlcgen.ListResourceProgressByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
		Status: statusFilter,
	})
	if err != nil {
		return nil, 0, err
	}

	progresses := make([]model.ResourceProgress, len(rows))
	for i, row := range rows {
		progresses[i] = toModelResourceProgress(row.ResourceProgress)
		if row.ResourceID2 != nil {
			progresses[i].Resource = &model.LearningResource{
				ID:          uint(*row.ResourceID2),
				UserID:      uint(fromInt64PtrValue(row.ResourceUserID)),
				Title:       fromStringPtr(row.ResourceTitle),
				Description: fromStringPtr(row.ResourceDescription),
				URL:         fromStringPtr(row.ResourceUrl),
				Category:    model.ResourceCategory(fromStringPtr(row.ResourceCategory)),
				Difficulty:  model.ResourceDifficulty(fromStringPtr(row.ResourceDifficulty)),
				Tags:        fromStringPtr(row.ResourceTags),
				ImageURL:    fromStringPtr(row.ResourceImageUrl),
				IsPublic:    fromBoolPtr(row.ResourceIsPublic),
				LikeCount:   int(fromInt64PtrValue(row.ResourceLikeCount)),
				SaveCount:   int(fromInt64PtrValue(row.ResourceSaveCount)),
				CreatedAt:   timeValue(fromTimestamptz(row.ResourceCreatedAt)),
				UpdatedAt:   timeValue(fromTimestamptz(row.ResourceUpdatedAt)),
			}
		}
	}
	return progresses, total, nil
}
