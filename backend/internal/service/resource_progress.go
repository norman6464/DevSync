package service

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// validProgressStatuses は有効なリソース進捗ステータスの一覧。
var validProgressStatuses = map[string]bool{
	string(model.ResourceProgressNotStarted): true,
	string(model.ResourceProgressInProgress): true,
	string(model.ResourceProgressCompleted):  true,
}

// ResourceProgressService は学習リソース進捗のビジネスロジックを提供する。
type ResourceProgressService struct {
	progressRepo repository.ResourceProgressRepositoryInterface
	resourceRepo repository.LearningResourceRepositoryInterface
}

// NewResourceProgressService は新しいResourceProgressServiceインスタンスを生成する。
func NewResourceProgressService(progressRepo repository.ResourceProgressRepositoryInterface, resourceRepo repository.LearningResourceRepositoryInterface) *ResourceProgressService {
	return &ResourceProgressService{progressRepo: progressRepo, resourceRepo: resourceRepo}
}

// UpsertProgress はリソース進捗をUPSERT（作成/更新）する。
func (s *ResourceProgressService) UpsertProgress(userID, resourceID uint, status string, completionPercent int, note string) (*model.ResourceProgress, error) {
	// リソースの存在確認
	if _, err := s.resourceRepo.FindByID(resourceID); err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "リソースが見つかりません", err)
	}

	// ステータスバリデーション
	if !validProgressStatuses[status] {
		return nil, domain.NewError(domain.ErrCodeValidation, "無効なステータスです。not_started, in_progress, completed のいずれかを指定してください", nil)
	}

	// 進捗率バリデーション
	if completionPercent < 0 || completionPercent > 100 {
		return nil, domain.NewError(domain.ErrCodeValidation, "進捗率は0〜100の範囲で指定してください", nil)
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

	if err := s.progressRepo.Upsert(progress); err != nil {
		return nil, err
	}

	// 最新の状態を返す
	return s.progressRepo.FindByUserAndResource(userID, resourceID)
}

// GetProgress は指定ユーザー・リソースの進捗を取得する。
func (s *ResourceProgressService) GetProgress(userID, resourceID uint) (*model.ResourceProgress, error) {
	return s.progressRepo.FindByUserAndResource(userID, resourceID)
}

// GetProgressList は指定ユーザーの進捗一覧を取得する。
func (s *ResourceProgressService) GetProgressList(userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error) {
	return s.progressRepo.FindByUserID(userID, status, limit, offset)
}
