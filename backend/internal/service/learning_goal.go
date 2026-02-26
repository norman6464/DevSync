package service

import (
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
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
	if err := domain.ValidateStringLength(goal.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	goal.Title = strings.TrimSpace(goal.Title)
	return s.repo.Create(goal)
}

// GetByID は指定IDの学習目標を取得する。所有権を検証する。
func (s *LearningGoalService) GetByID(id, userID uint) (*model.LearningGoal, error) {
	return s.findAndCheckOwnership(id, userID)
}

// GetByUserID は指定ユーザーの学習目標をページネーション付きで取得する。
func (s *LearningGoalService) GetByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

// GetActiveByUserID は指定ユーザーのアクティブな学習目標のみを取得する。
func (s *LearningGoalService) GetActiveByUserID(userID uint) ([]model.LearningGoal, error) {
	return s.repo.GetActiveByUserID(userID)
}

// validGoalCategories は有効なGoalCategoryの集合。
var validGoalCategories = map[string]bool{
	string(model.GoalCategoryLanguage):  true,
	string(model.GoalCategoryFramework): true,
	string(model.GoalCategorySkill):     true,
	string(model.GoalCategoryProject):   true,
	string(model.GoalCategoryOther):     true,
}

// GetByCategory は指定ユーザーの学習目標をカテゴリでフィルタリングして取得する。
func (s *LearningGoalService) GetByCategory(userID uint, category string) ([]model.LearningGoal, error) {
	if !validGoalCategories[category] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なカテゴリです", nil)
	}
	return s.repo.GetByCategory(userID, category)
}

// validGoalStatuses は有効なGoalStatusの集合。
var validGoalStatuses = map[string]bool{
	string(model.GoalStatusActive):    true,
	string(model.GoalStatusCompleted): true,
	string(model.GoalStatusPaused):    true,
}

// GetByStatus は指定ユーザーの学習目標をステータスでフィルタリングして取得する。
func (s *LearningGoalService) GetByStatus(userID uint, status string) ([]model.LearningGoal, error) {
	if !validGoalStatuses[status] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なステータスです", nil)
	}
	return s.repo.GetByStatus(userID, status)
}

// GetStats は指定ユーザーの学習目標統計情報を取得する。
func (s *LearningGoalService) GetStats(userID uint) (*model.LearningGoalStats, error) {
	return s.repo.GetStats(userID)
}

// findAndCheckOwnership は学習目標を取得し、指定ユーザーが所有者かを検証する。
func (s *LearningGoalService) findAndCheckOwnership(id, userID uint) (*model.LearningGoal, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(g *model.LearningGoal) uint { return g.UserID })
}

// Update は所有権を検証した後、学習目標を更新する。
// 進捗が100%に達した場合、ステータスを自動的に「完了」に変更する。
func (s *LearningGoalService) Update(id, userID uint, updates *model.LearningGoal) (*model.LearningGoal, error) {
	goal, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(updates.Title) != "" {
		if len(strings.TrimSpace(updates.Title)) > 200 {
			return nil, domain.NewError(domain.ErrCodeValidation, "タイトルは200文字以下である必要があります", nil)
		}
		goal.Title = strings.TrimSpace(updates.Title)
	}
	if strings.TrimSpace(updates.Description) != "" {
		if len(strings.TrimSpace(updates.Description)) > 1000 {
			return nil, domain.NewError(domain.ErrCodeValidation, "説明は1000文字以下である必要があります", nil)
		}
		goal.Description = strings.TrimSpace(updates.Description)
	}
	if strings.TrimSpace(string(updates.Category)) != "" {
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
	if strings.TrimSpace(string(updates.Status)) != "" {
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

// DeadlineStatus はアクティブな目標のデッドライン状態を判定する純粋関数。
// "overdue"（期限超過）、"approaching"（3日以内）、""（安全/対象外）を返す。
func DeadlineStatus(goal *model.LearningGoal, now time.Time) string {
	if goal.TargetDate == nil || goal.Status != model.GoalStatusActive {
		return ""
	}
	days := int(goal.TargetDate.Sub(now).Hours() / 24)
	if goal.TargetDate.Before(now) && days < 0 {
		return "overdue"
	}
	if days <= 3 {
		return "approaching"
	}
	return ""
}

// DaysUntilDeadline は目標の期限までの残り日数を返す。
// TargetDateがnilまたは期限超過の場合は-1を返す。
func DaysUntilDeadline(goal *model.LearningGoal, now time.Time) int {
	if goal.TargetDate == nil {
		return -1
	}
	days := int(goal.TargetDate.Sub(now).Hours() / 24)
	if days < 0 {
		return -1
	}
	return days
}

// GetDeadlineAlerts はユーザーのアクティブ目標からデッドラインアラートを取得する。
func (s *LearningGoalService) GetDeadlineAlerts(userID uint) ([]model.GoalDeadlineAlert, error) {
	goals, err := s.repo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var alerts []model.GoalDeadlineAlert
	for _, g := range goals {
		status := DeadlineStatus(&g, now)
		if status != "" {
			alerts = append(alerts, model.GoalDeadlineAlert{
				Goal:     g,
				Status:   status,
				DaysLeft: DaysUntilDeadline(&g, now),
			})
		}
	}
	return alerts, nil
}

// Duplicate は所有権を検証した後、学習目標を複製する。
// 進捗は0%にリセットし、ステータスはactiveに戻す。タイトルに「(コピー)」を付与する。
func (s *LearningGoalService) Duplicate(id, userID uint) (*model.LearningGoal, error) {
	goal, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	newGoal := &model.LearningGoal{
		UserID:      goal.UserID,
		Title:       goal.Title + " (コピー)",
		Description: goal.Description,
		Category:    goal.Category,
		TargetDate:  goal.TargetDate,
		Progress:    0,
		Status:      model.GoalStatusActive,
	}

	if err := s.repo.Create(newGoal); err != nil {
		return nil, err
	}
	return newGoal, nil
}

// ToggleShare は学習目標の公開/非公開を切り替える。所有者のみ操作可能。
func (s *LearningGoalService) ToggleShare(id, userID uint) (*model.LearningGoal, error) {
	goal, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	goal.IsPublic = !goal.IsPublic
	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

// GetPublicGoals は全ユーザーの公開済み学習目標をページネーション付きで取得する。
func (s *LearningGoalService) GetPublicGoals(limit, offset int) ([]model.LearningGoal, int64, error) {
	return s.repo.GetPublicGoals(limit, offset)
}

// GetPublicByUserID は指定ユーザーの公開済み学習目標をページネーション付きで取得する。
func (s *LearningGoalService) GetPublicByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	return s.repo.GetPublicByUserID(userID, limit, offset)
}

// Delete は所有権を検証した後、学習目標を削除する。
func (s *LearningGoalService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
