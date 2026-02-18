package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestLearningGoalService はLearningGoalServiceのテスト用インスタンスを生成するヘルパー。
func newTestLearningGoalService() (*LearningGoalService, *MockLearningGoalRepository) {
	repo := new(MockLearningGoalRepository)
	svc := NewLearningGoalService(repo)
	return svc, repo
}

// ============================================================
// 学習目標作成テスト
// ============================================================

func TestLearningGoalCreate_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goal := &model.LearningGoal{
		UserID:   1,
		Title:    "Goを学ぶ",
		Category: model.GoalCategoryOther,
	}

	repo.On("Create", goal).Return(nil)

	err := svc.Create(goal)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningGoalCreate_Error(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goal := &model.LearningGoal{UserID: 1, Title: "Test"}
	repo.On("Create", goal).Return(errors.New("db error"))

	err := svc.Create(goal)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// 学習目標取得テスト
// ============================================================

func TestLearningGoalGetByID_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	expected := &model.LearningGoal{Title: "Go学習", UserID: 1}
	expected.ID = 1
	repo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Go学習", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByID_NotFound(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByUserID_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goals := []model.LearningGoal{
		{Title: "Goal 1", UserID: 1},
		{Title: "Goal 2", UserID: 1},
	}
	repo.On("GetByUserID", uint(1)).Return(goals, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByUserID_Error(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, errors.New("db error"))

	result, err := svc.GetByUserID(1)
	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetActiveByUserID_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goals := []model.LearningGoal{
		{Title: "Active Goal", UserID: 1, Status: model.GoalStatusActive},
	}
	repo.On("GetActiveByUserID", uint(1)).Return(goals, nil)

	result, err := svc.GetActiveByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, model.GoalStatusActive, result[0].Status)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetActiveByUserID_Error(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("GetActiveByUserID", uint(1)).Return([]model.LearningGoal{}, errors.New("db error"))

	result, err := svc.GetActiveByUserID(1)
	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetStats_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	stats := &model.LearningGoalStats{
		TotalGoals:     5,
		ActiveGoals:    3,
		CompletedGoals: 2,
	}
	repo.On("GetStats", uint(1)).Return(stats, nil)

	result, err := svc.GetStats(1)
	assert.NoError(t, err)
	assert.Equal(t, 5, result.TotalGoals)
	assert.Equal(t, 3, result.ActiveGoals)
	assert.Equal(t, 2, result.CompletedGoals)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetStats_Error(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("GetStats", uint(1)).Return(nil, errors.New("db error"))

	result, err := svc.GetStats(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 学習目標更新テスト
// ============================================================

func TestLearningGoalUpdate_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{
		Title:    "Old Title",
		UserID:   1,
		Status:   model.GoalStatusActive,
		Progress: 50,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.LearningGoal{Title: "New Title"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.LearningGoal{Title: "New Title"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_AutoCompleteAt100(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{
		Title:    "Goal",
		UserID:   1,
		Status:   model.GoalStatusActive,
		Progress: 80,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 進捗100%で自動完了
	updates := &model.LearningGoal{Progress: 100}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, 100, result.Progress)
	assert.Equal(t, model.GoalStatusCompleted, result.Status)
	assert.NotNil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_ProgressClampTo100(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{
		Title:    "Goal",
		UserID:   1,
		Status:   model.GoalStatusActive,
		Progress: 90,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 進捗101はクランプされて100になる
	updates := &model.LearningGoal{Progress: 150}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, 100, result.Progress)
	assert.Equal(t, model.GoalStatusCompleted, result.Status)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_SetCompletedStatus(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{
		Title:  "Goal",
		UserID: 1,
		Status: model.GoalStatusActive,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// ステータスを直接completedに変更
	updates := &model.LearningGoal{Status: model.GoalStatusCompleted}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, model.GoalStatusCompleted, result.Status)
	assert.NotNil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_NotFound(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	updates := &model.LearningGoal{Title: "New Title"}
	result, err := svc.Update(999, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 学習目標削除テスト
// ============================================================

func TestLearningGoalDelete_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningGoalDelete_Forbidden(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// デッドラインステータス判定テスト（純粋関数）
// ============================================================

func TestDeadlineStatus_NoTargetDate(t *testing.T) {
	now := time.Now()
	goal := model.LearningGoal{Status: model.GoalStatusActive}
	assert.Equal(t, "", DeadlineStatus(&goal, now))
}

func TestDeadlineStatus_CompletedGoal(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	goal := model.LearningGoal{
		Status:     model.GoalStatusCompleted,
		TargetDate: &yesterday,
	}
	assert.Equal(t, "", DeadlineStatus(&goal, now))
}

func TestDeadlineStatus_PausedGoal(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	goal := model.LearningGoal{
		Status:     model.GoalStatusPaused,
		TargetDate: &yesterday,
	}
	assert.Equal(t, "", DeadlineStatus(&goal, now))
}

func TestDeadlineStatus_Overdue(t *testing.T) {
	now := time.Now()
	pastDate := now.AddDate(0, 0, -2)
	goal := model.LearningGoal{
		Status:     model.GoalStatusActive,
		TargetDate: &pastDate,
	}
	assert.Equal(t, "overdue", DeadlineStatus(&goal, now))
}

func TestDeadlineStatus_Approaching_Today(t *testing.T) {
	now := time.Now()
	today := now
	goal := model.LearningGoal{
		Status:     model.GoalStatusActive,
		TargetDate: &today,
	}
	assert.Equal(t, "approaching", DeadlineStatus(&goal, now))
}

func TestDeadlineStatus_Approaching_3Days(t *testing.T) {
	now := time.Now()
	threeDays := now.AddDate(0, 0, 3)
	goal := model.LearningGoal{
		Status:     model.GoalStatusActive,
		TargetDate: &threeDays,
	}
	assert.Equal(t, "approaching", DeadlineStatus(&goal, now))
}

func TestDeadlineStatus_Safe_4Days(t *testing.T) {
	now := time.Now()
	fourDays := now.AddDate(0, 0, 4)
	goal := model.LearningGoal{
		Status:     model.GoalStatusActive,
		TargetDate: &fourDays,
	}
	assert.Equal(t, "", DeadlineStatus(&goal, now))
}

// ============================================================
// DaysUntilDeadline テスト
// ============================================================

func TestDaysUntilDeadline_NoTargetDate(t *testing.T) {
	now := time.Now()
	goal := model.LearningGoal{}
	assert.Equal(t, -1, DaysUntilDeadline(&goal, now))
}

func TestDaysUntilDeadline_Tomorrow(t *testing.T) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	goal := model.LearningGoal{TargetDate: &tomorrow}
	assert.Equal(t, 1, DaysUntilDeadline(&goal, now))
}

func TestDaysUntilDeadline_Yesterday(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	goal := model.LearningGoal{TargetDate: &yesterday}
	assert.Equal(t, -1, DaysUntilDeadline(&goal, now))
}

// ============================================================
// GetDeadlineAlerts サービステスト
// ============================================================

func TestGetDeadlineAlerts_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	now := time.Now()
	overdue := now.AddDate(0, 0, -2)
	approaching := now.AddDate(0, 0, 2)
	safe := now.AddDate(0, 0, 10)

	goals := []model.LearningGoal{
		{Title: "Overdue Goal", Status: model.GoalStatusActive, TargetDate: &overdue},
		{Title: "Approaching Goal", Status: model.GoalStatusActive, TargetDate: &approaching},
		{Title: "Safe Goal", Status: model.GoalStatusActive, TargetDate: &safe},
		{Title: "No Deadline", Status: model.GoalStatusActive},
	}
	repo.On("GetActiveByUserID", uint(1)).Return(goals, nil)

	alerts, err := svc.GetDeadlineAlerts(1)
	assert.NoError(t, err)
	assert.Len(t, alerts, 2)
	assert.Equal(t, "overdue", alerts[0].Status)
	assert.Equal(t, "Overdue Goal", alerts[0].Goal.Title)
	assert.Equal(t, "approaching", alerts[1].Status)
	assert.Equal(t, "Approaching Goal", alerts[1].Goal.Title)
	repo.AssertExpectations(t)
}

func TestGetDeadlineAlerts_Empty(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("GetActiveByUserID", uint(1)).Return([]model.LearningGoal{}, nil)

	alerts, err := svc.GetDeadlineAlerts(1)
	assert.NoError(t, err)
	assert.Empty(t, alerts)
	repo.AssertExpectations(t)
}

func TestGetDeadlineAlerts_RepoError(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("GetActiveByUserID", uint(1)).Return([]model.LearningGoal{}, assert.AnError)

	alerts, err := svc.GetDeadlineAlerts(1)
	assert.Error(t, err)
	assert.Nil(t, alerts)
	repo.AssertExpectations(t)
}
