package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ListResourceProgressUseCase は指定ユーザーの進捗一覧を取得する。
type ListResourceProgressUseCase struct {
	progress repository.ResourceProgressRepository
}

// NewListResourceProgressUseCase は ListResourceProgressUseCase を生成する。
func NewListResourceProgressUseCase(progress repository.ResourceProgressRepository) *ListResourceProgressUseCase {
	return &ListResourceProgressUseCase{progress: progress}
}

// Execute はユーザーの進捗一覧と総件数を返す。status が空でなければ絞り込む。
func (uc *ListResourceProgressUseCase) Execute(ctx context.Context, userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error) {
	return uc.progress.FindByUserID(ctx, userID, status, limit, offset)
}
