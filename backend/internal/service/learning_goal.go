package service

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningGoalService は学習目標のビジネスロジックを提供する。
// 目標のCRUD操作と、進捗100%時の自動完了ロジックを担当する。
type LearningGoalService struct {
	repo repository.LearningGoalRepositoryInterface
}

// NewLearningGoalService は新しいLearningGoalServiceインスタンスを生成する。
func NewLearningGoalService(repo repository.LearningGoalRepositoryInterface) *LearningGoalService {
	return &LearningGoalService{repo: repo}
}

// Create は新しい学習目標を作成する。
func (s *LearningGoalService) Create(goal *model.LearningGoal) error {
	return s.repo.Create(goal)
}

// GetByID は指定IDの学習目標を取得する。
func (s *LearningGoalService) GetByID(id uint) (*model.LearningGoal, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの全学習目標を取得する。
func (s *LearningGoalService) GetByUserID(userID uint) ([]model.LearningGoal, error) {
	return s.repo.GetByUserID(userID)
}

// GetActiveByUserID は指定ユーザーのアクティブな学習目標のみを取得する。
func (s *LearningGoalService) GetActiveByUserID(userID uint) ([]model.LearningGoal, error) {
	return s.repo.GetActiveByUserID(userID)
}

// GetStats は指定ユーザーの学習目標統計情報を取得する。
func (s *LearningGoalService) GetStats(userID uint) (*model.LearningGoalStats, error) {
	return s.repo.GetStats(userID)
}

// Update は所有権を検証した後、学習目標を更新する。
// 進捗が100%に達した場合、ステータスを自動的に「完了」に変更する。
func (s *LearningGoalService) Update(id, userID uint, updates *model.LearningGoal) (*model.LearningGoal, error) {
	goal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if goal.UserID != userID {
		return nil, ErrForbidden
	}

	if updates.Title != "" {
		goal.Title = updates.Title
	}
	if updates.Description != "" {
		goal.Description = updates.Description
	}
	if updates.Category != "" {
		goal.Category = updates.Category
	}
	if updates.TargetDate != nil {
		goal.TargetDate = updates.TargetDate
	}
	if updates.Progress >= 0 {
		progress := updates.Progress
		if progress > 100 {
			progress = 100
		}
		goal.Progress = progress

		// 進捗100%達成時に自動完了
		if progress == 100 && goal.Status == model.GoalStatusActive {
			goal.Status = model.GoalStatusCompleted
			now := time.Now()
			goal.CompletedAt = &now
		}
	}
	if updates.Status != "" {
		goal.Status = updates.Status
		if goal.Status == model.GoalStatusCompleted && goal.CompletedAt == nil {
			now := time.Now()
			goal.CompletedAt = &now
		}
	}

	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

// Delete は所有権を検証した後、学習目標を削除する。
func (s *LearningGoalService) Delete(id, userID uint) error {
	goal, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if goal.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}
