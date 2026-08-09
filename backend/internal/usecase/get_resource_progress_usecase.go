package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetResourceProgressUseCase は指定ユーザー・リソースの進捗を取得する。
type GetResourceProgressUseCase struct {
	progress repository.ResourceProgressRepository
}

// NewGetResourceProgressUseCase は GetResourceProgressUseCase を生成する。
func NewGetResourceProgressUseCase(progress repository.ResourceProgressRepository) *GetResourceProgressUseCase {
	return &GetResourceProgressUseCase{progress: progress}
}

// Execute は指定ユーザー・リソースの進捗を返す。
func (uc *GetResourceProgressUseCase) Execute(ctx context.Context, userID, resourceID uint) (*model.ResourceProgress, error) {
	return uc.progress.FindByUserAndResource(ctx, userID, resourceID)
}
