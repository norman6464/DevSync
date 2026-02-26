package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestLearningGoalCreate_WhitespaceOnlyTitle(t *testing.T) {
	svc, _ := newTestLearningGoalService()

	// 空白のみのタイトル → エラーになるべき
	goal := &model.LearningGoal{UserID: 1, Title: "   "}
	err := svc.Create(goal)
	assert.Error(t, err)
}

func TestLearningGoalCreate_EmptyTitle(t *testing.T) {
	svc, _ := newTestLearningGoalService()

	// 空のタイトル → エラーになるべき
	goal := &model.LearningGoal{UserID: 1, Title: ""}
	err := svc.Create(goal)
	assert.Error(t, err)
}

// ============================================================
// 学習目標取得テスト
// ============================================================

func TestLearningGoalGetByID_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	expected := &model.LearningGoal{Title: "Go学習", UserID: 1}
	expected.ID = 1
	repo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := svc.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Go学習", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByID_Forbidden(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goal := &model.LearningGoal{Title: "Go学習", UserID: 1}
	goal.ID = 1
	repo.On("FindByID", uint(1)).Return(goal, nil)

	result, err := svc.GetByID(1, 999)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrForbidden, err)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByID_NotFound(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999, 1)
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
	repo.On("GetByUserID", uint(1), 20, 0).Return(goals, int64(2), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByUserID_Error(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("GetByUserID", uint(1), 20, 0).Return([]model.LearningGoal{}, int64(0), errors.New("db error"))

	result, _, err := svc.GetByUserID(1, 20, 0)
	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByUserID_Page2(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("GetByUserID", uint(1), 10, 10).Return([]model.LearningGoal{}, int64(15), nil)

	result, total, err := svc.GetByUserID(1, 10, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(15), total)
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

func TestLearningGoalUpdate_RepoError(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	existing := &model.LearningGoal{Title: "Goal", UserID: 1, Status: model.GoalStatusActive}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(errors.New("db error"))
	updates := &model.LearningGoal{Title: "New Title"}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLearningGoalDelete_NotFound(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.Delete(99, 1)
	assert.Error(t, err)
}

func TestLearningGoalUpdate_WithTargetDate(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	existing := &model.LearningGoal{Title: "Goal", UserID: 1, Status: model.GoalStatusActive}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)
	targetDate := time.Now().AddDate(0, 1, 0)
	updates := &model.LearningGoal{
		Title:       "Updated",
		Description: "Desc",
		Category:    model.GoalCategoryOther,
		TargetDate:  &targetDate,
	}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Title)
	assert.Equal(t, "Desc", result.Description)
	assert.NotNil(t, result.TargetDate)
}

// ============================================================
// カテゴリ別取得テスト
// ============================================================

func TestLearningGoalGetByCategory_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	expected := []model.LearningGoal{
		{ID: 1, UserID: 1, Title: "Go習得", Category: model.GoalCategoryLanguage},
		{ID: 2, UserID: 1, Title: "Rust入門", Category: model.GoalCategoryLanguage},
	}
	repo.On("GetByCategory", uint(1), "language").Return(expected, nil)

	result, err := svc.GetByCategory(1, "language")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByCategory_EmptyResult(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	repo.On("GetByCategory", uint(1), "framework").Return([]model.LearningGoal{}, nil)

	result, err := svc.GetByCategory(1, "framework")
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByCategory_InvalidCategory(t *testing.T) {
	svc, _ := newTestLearningGoalService()

	_, err := svc.GetByCategory(1, "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "無効なカテゴリ")
}

func TestLearningGoalGetByCategory_RepoError(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	repo.On("GetByCategory", uint(1), "language").Return([]model.LearningGoal{}, errors.New("db error"))

	_, err := svc.GetByCategory(1, "language")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// ステータス別目標取得テスト
// ============================================================

func TestLearningGoalGetByStatus_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	expected := []model.LearningGoal{
		{UserID: 1, Title: "Go学習", Status: model.GoalStatusCompleted},
		{UserID: 1, Title: "React学習", Status: model.GoalStatusCompleted},
	}
	repo.On("GetByStatus", uint(1), "completed").Return(expected, nil)

	result, err := svc.GetByStatus(1, "completed")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByStatus_EmptyResult(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	repo.On("GetByStatus", uint(1), "paused").Return([]model.LearningGoal{}, nil)

	result, err := svc.GetByStatus(1, "paused")
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetByStatus_InvalidStatus(t *testing.T) {
	svc, _ := newTestLearningGoalService()

	_, err := svc.GetByStatus(1, "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "無効なステータスです")
}

func TestLearningGoalGetByStatus_RepoError(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	repo.On("GetByStatus", uint(1), "active").Return([]model.LearningGoal{}, errors.New("db error"))

	_, err := svc.GetByStatus(1, "active")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// 空白バイパス脆弱性テスト
// ============================================================

func TestLearningGoalUpdate_WhitespaceTitle(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	existing := &model.LearningGoal{Title: "Original Title", Description: "Desc", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.LearningGoal{Title: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Original Title", result.Title)
}

func TestLearningGoalUpdate_WhitespaceDescription(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	existing := &model.LearningGoal{Title: "Title", Description: "Original Desc", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.LearningGoal{Description: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Original Desc", result.Description)
}

func TestLearningGoalUpdate_WhitespaceCategory(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	existing := &model.LearningGoal{Title: "Title", Category: model.GoalCategoryLanguage, UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.LearningGoal{Category: "   "})
	assert.NoError(t, err)
	assert.Equal(t, model.GoalCategoryLanguage, result.Category)
}

func TestLearningGoalUpdate_WhitespaceStatus(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	existing := &model.LearningGoal{Title: "Title", Status: model.GoalStatusActive, UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.LearningGoal{Status: "   "})
	assert.NoError(t, err)
	assert.Equal(t, model.GoalStatusActive, result.Status)
}

func TestLearningGoalCreate_TitleTooLong(t *testing.T) {
	svc, _ := newTestLearningGoalService()

	goal := &model.LearningGoal{UserID: 1, Title: strings.Repeat("あ", 201)}
	err := svc.Create(goal)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "タイトルは200文字以下")
}

func TestLearningGoalUpdate_TitleTooLong(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	existing := &model.LearningGoal{Title: "Title", Status: model.GoalStatusActive, UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Update(1, 1, &model.LearningGoal{Title: strings.Repeat("あ", 201)})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "タイトルは200文字以下")
}

func TestLearningGoalUpdate_DescriptionTooLong(t *testing.T) {
	svc, repo := newTestLearningGoalService()
	existing := &model.LearningGoal{Title: "Title", Status: model.GoalStatusActive, UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Update(1, 1, &model.LearningGoal{Description: strings.Repeat("あ", 1001)})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "説明は1000文字以下")
}

// ============================================================
// 学習目標 複製テスト
// ============================================================

func TestLearningGoalDuplicate_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{
		UserID:      1,
		Title:       "Goマスター",
		Description: "Go言語を完全習得する",
		Category:    model.GoalCategoryLanguage,
		Progress:    80,
		Status:      model.GoalStatusCompleted,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Create", mock.MatchedBy(func(g *model.LearningGoal) bool {
		return g.Title == "Goマスター (コピー)" &&
			g.Description == "Go言語を完全習得する" &&
			g.Category == model.GoalCategoryLanguage &&
			g.Progress == 0 &&
			g.Status == model.GoalStatusActive &&
			g.UserID == 1 &&
			g.CompletedAt == nil
	})).Return(nil)

	result, err := svc.Duplicate(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Goマスター (コピー)", result.Title)
	assert.Equal(t, 0, result.Progress)
	assert.Equal(t, model.GoalStatusActive, result.Status)
	assert.Nil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

func TestLearningGoalDuplicate_Forbidden(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{UserID: 99, Title: "他人の目標"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Duplicate(1, 1)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
}

func TestLearningGoalDuplicate_NotFound(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.Duplicate(99, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLearningGoalDuplicate_CreateError(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{UserID: 1, Title: "目標"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Create", mock.AnythingOfType("*model.LearningGoal")).Return(errors.New("db error"))

	result, err := svc.Duplicate(1, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 共有トグルテスト
// ============================================================

func TestLearningGoalToggleShare_MakePublic(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goal := &model.LearningGoal{UserID: 1, Title: "公開テスト", IsPublic: false}
	goal.ID = 1
	repo.On("FindByID", uint(1)).Return(goal, nil)
	repo.On("Update", mock.MatchedBy(func(g *model.LearningGoal) bool {
		return g.IsPublic == true
	})).Return(nil)

	result, err := svc.ToggleShare(1, 1)
	assert.NoError(t, err)
	assert.True(t, result.IsPublic)
	repo.AssertExpectations(t)
}

func TestLearningGoalToggleShare_MakePrivate(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goal := &model.LearningGoal{UserID: 1, Title: "非公開テスト", IsPublic: true}
	goal.ID = 1
	repo.On("FindByID", uint(1)).Return(goal, nil)
	repo.On("Update", mock.MatchedBy(func(g *model.LearningGoal) bool {
		return g.IsPublic == false
	})).Return(nil)

	result, err := svc.ToggleShare(1, 1)
	assert.NoError(t, err)
	assert.False(t, result.IsPublic)
	repo.AssertExpectations(t)
}

func TestLearningGoalToggleShare_Forbidden(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goal := &model.LearningGoal{UserID: 999, Title: "他人の目標"}
	goal.ID = 1
	repo.On("FindByID", uint(1)).Return(goal, nil)

	_, err := svc.ToggleShare(1, 1)
	assert.Error(t, err)
}

func TestLearningGoalToggleShare_NotFound(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	_, err := svc.ToggleShare(99, 1)
	assert.Error(t, err)
}

func TestLearningGoalToggleShare_UpdateError(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goal := &model.LearningGoal{UserID: 1, Title: "テスト目標", IsPublic: false}
	goal.ID = 1
	repo.On("FindByID", uint(1)).Return(goal, nil)
	repo.On("Update", mock.Anything).Return(errors.New("db error"))

	_, err := svc.ToggleShare(1, 1)
	assert.Error(t, err)
}

func TestLearningGoalGetPublicGoals_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goals := []model.LearningGoal{
		{Title: "公開目標1", IsPublic: true},
		{Title: "公開目標2", IsPublic: true},
	}
	repo.On("GetPublicGoals", 20, 0).Return(goals, int64(2), nil)

	result, total, err := svc.GetPublicGoals(20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetPublicByUserID_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	goals := []model.LearningGoal{{Title: "公開目標", IsPublic: true, UserID: 5}}
	repo.On("GetPublicByUserID", uint(5), 20, 0).Return(goals, int64(1), nil)

	result, total, err := svc.GetPublicByUserID(5, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

// ============================================================
// CalculateGoalForecast テスト
// ============================================================

func TestCalculateGoalForecast_NoTargetDate_NoDailyAvg(t *testing.T) {
	goal := &model.LearningGoal{Title: "Go学習", Progress: 30}
	goal.ID = 1

	f := CalculateGoalForecast(goal, 0, 0, time.Now())
	assert.Equal(t, uint(1), f.GoalID)
	assert.Equal(t, "Go学習", f.Title)
	assert.Equal(t, 30, f.CurrentProgress)
	assert.Equal(t, -1, f.EstimatedDaysLeft)
	assert.Equal(t, -1, f.DaysUntilDeadline)
	assert.Equal(t, "unknown", f.Difficulty)
	assert.False(t, f.OnTrack)
}

func TestCalculateGoalForecast_WithDeadline_OnTrack(t *testing.T) {
	deadline := time.Now().Add(30 * 24 * time.Hour)
	goal := &model.LearningGoal{
		Title:       "React学習",
		TargetHours: 10,
		TargetDate:  &deadline,
	}
	goal.ID = 2

	// 10時間 = 600分目標、既に300分学習済み、日平均30分 → 残り10日
	f := CalculateGoalForecast(goal, 300, 30, time.Now())
	assert.Equal(t, 10, f.EstimatedDaysLeft)
	assert.Equal(t, 30, f.DaysUntilDeadline)
	assert.True(t, f.OnTrack)
	assert.Equal(t, "easy", f.Difficulty) // 10/30 ≈ 0.33 ≤ 0.5
}

func TestCalculateGoalForecast_WithDeadline_OffTrack(t *testing.T) {
	deadline := time.Now().Add(5 * 24 * time.Hour)
	goal := &model.LearningGoal{
		Title:       "難しい目標",
		TargetHours: 100,
		TargetDate:  &deadline,
	}
	goal.ID = 3

	// 100時間 = 6000分、0分学習済み、日平均30分 → 200日必要
	f := CalculateGoalForecast(goal, 0, 30, time.Now())
	assert.Equal(t, 200, f.EstimatedDaysLeft)
	assert.Equal(t, 5, f.DaysUntilDeadline)
	assert.False(t, f.OnTrack)
	assert.Equal(t, "hard", f.Difficulty) // 200/5 = 40.0 > 1.0
}

func TestCalculateGoalForecast_AlreadyCompleted(t *testing.T) {
	deadline := time.Now().Add(10 * 24 * time.Hour)
	goal := &model.LearningGoal{
		Title:       "完了済み",
		TargetHours: 5,
		TargetDate:  &deadline,
	}
	goal.ID = 4

	// 5時間 = 300分目標、既に400分学習済み → 残り0日
	f := CalculateGoalForecast(goal, 400, 60, time.Now())
	assert.Equal(t, 0, f.EstimatedDaysLeft)
	assert.True(t, f.OnTrack)
	assert.Equal(t, "easy", f.Difficulty)
}

func TestCalculateGoalForecast_NoDailyAvg_WithTargetHours(t *testing.T) {
	goal := &model.LearningGoal{
		Title:       "学習開始前",
		TargetHours: 10,
	}
	goal.ID = 5

	f := CalculateGoalForecast(goal, 0, 0, time.Now())
	assert.Equal(t, -1, f.EstimatedDaysLeft)
	assert.Equal(t, "hard", f.Difficulty)
}

func TestCalculateGoalForecast_MediumDifficulty(t *testing.T) {
	deadline := time.Now().Add(10 * 24 * time.Hour)
	goal := &model.LearningGoal{
		Title:       "中程度",
		TargetHours: 5,
		TargetDate:  &deadline,
	}
	goal.ID = 6

	// 300分目標、0分学習済み、日平均40分 → 8日必要、10日猶予 → ratio=0.8
	f := CalculateGoalForecast(goal, 0, 40, time.Now())
	assert.Equal(t, 8, f.EstimatedDaysLeft)
	assert.True(t, f.OnTrack)
	assert.Equal(t, "medium", f.Difficulty) // 8/10 = 0.8, 0.5 < 0.8 ≤ 1.0
}

func TestCalculateGoalForecast_PastDeadline(t *testing.T) {
	deadline := time.Now().Add(-2 * 24 * time.Hour)
	goal := &model.LearningGoal{
		Title:      "期限切れ",
		TargetDate: &deadline,
	}
	goal.ID = 7

	f := CalculateGoalForecast(goal, 0, 0, time.Now())
	assert.Equal(t, 0, f.DaysUntilDeadline) // 負の値は0に正規化
}
