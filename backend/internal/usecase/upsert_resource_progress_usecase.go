package usecase

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// validProgressStatuses は有効なリソース進捗ステータスの一覧。
var validProgressStatuses = map[string]bool{
	string(model.ResourceProgressNotStarted): true,
	string(model.ResourceProgressInProgress): true,
	string(model.ResourceProgressCompleted):  true,
}

// UpsertResourceProgressUseCase は学習リソース進捗を作成/更新する。
type UpsertResourceProgressUseCase struct {
	progress  repository.ResourceProgressRepository
	resources repository.LearningResourceReader
}

// NewUpsertResourceProgressUseCase は UpsertResourceProgressUseCase を生成する。
func NewUpsertResourceProgressUseCase(
	progress repository.ResourceProgressRepository,
	resources repository.LearningResourceReader,
) *UpsertResourceProgressUseCase {
	return &UpsertResourceProgressUseCase{progress: progress, resources: resources}
}

// Execute はリソース存在確認・入力検証を行い進捗を upsert し、最新状態を返す。
func (uc *UpsertResourceProgressUseCase) Execute(ctx context.Context, userID, resourceID uint, status string, completionPercent int, note string) (*model.ResourceProgress, error) {
	if _, err := uc.resources.FindByID(ctx, resourceID); err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "リソースが見つかりません", err)
	}

	if !validProgressStatuses[status] {
		return nil, domain.NewError(domain.ErrCodeValidation, "無効なステータスです。not_started, in_progress, completed のいずれかを指定してください", nil)
	}

	if completionPercent < 0 || completionPercent > 100 {
		return nil, domain.NewError(domain.ErrCodeValidation, "進捗率は0〜100の範囲で指定してください", nil)
	}

	if err := domain.ValidateStringLength(note, 0, 1000, "メモ"); err != nil {
		return nil, err
	}

	progress := &model.ResourceProgress{
		UserID:            userID,
		ResourceID:        resourceID,
		Status:            model.ResourceProgressStatus(status),
		CompletionPercent: completionPercent,
		Note:              note,
	}

	// ステータスに応じたタイムスタンプ設定
	now := time.Now()
	if status == string(model.ResourceProgressInProgress) {
		progress.StartedAt = &now
	}
	if status == string(model.ResourceProgressCompleted) {
		progress.StartedAt = &now
		progress.CompletedAt = &now
	}

	if err := uc.progress.Upsert(ctx, progress); err != nil {
		return nil, err
	}

	return uc.progress.FindByUserAndResource(ctx, userID, resourceID)
}
