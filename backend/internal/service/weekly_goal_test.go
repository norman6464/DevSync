package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestWeeklyGoalService() (*WeeklyGoalService, *MockWeeklyGoalRepository) {
	repo := new(MockWeeklyGoalRepository)
	svc := NewWeeklyGoalService(repo)
	return svc, repo
}

// ============================================================
// SetGoal テスト
// ============================================================

func TestWeeklyGoalSetGoal_Success(t *testing.T) {
	svc, repo := newTestWeeklyGoalService()

	repo.On("Upsert", mock.MatchedBy(func(g *model.WeeklyGoal) bool {
		return g.UserID == 1 && g.Category == model.LogCategoryCoding && g.TargetMinutes == 300
	})).Return(nil)

	goal, err := svc.SetGoal(1, "coding", 300)
	assert.NoError(t, err)
	assert.Equal(t, model.LogCategoryCoding, goal.Category)
	assert.Equal(t, 300, goal.TargetMinutes)
	repo.AssertExpectations(t)
}

func TestWeeklyGoalSetGoal_InvalidCategory(t *testing.T) {
	svc, _ := newTestWeeklyGoalService()

	_, err := svc.SetGoal(1, "invalid", 300)
	assert.Error(t, err)
}

func TestWeeklyGoalSetGoal_NegativeMinutes(t *testing.T) {
	svc, _ := newTestWeeklyGoalService()

	_, err := svc.SetGoal(1, "coding", -1)
	assert.Error(t, err)
}

func TestWeeklyGoalSetGoal_ExceedsMax(t *testing.T) {
	svc, _ := newTestWeeklyGoalService()

	_, err := svc.SetGoal(1, "coding", 10081)
	assert.Error(t, err)
}

func TestWeeklyGoalSetGoal_RepoError(t *testing.T) {
	svc, repo := newTestWeeklyGoalService()

	repo.On("Upsert", mock.Anything).Return(errors.New("db error"))

	_, err := svc.SetGoal(1, "coding", 300)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// GetGoals テスト
// ============================================================

func TestWeeklyGoalGetGoals_Success(t *testing.T) {
	svc, repo := newTestWeeklyGoalService()

	goals := []model.WeeklyGoal{
		{ID: 1, UserID: 1, Category: model.LogCategoryCoding, TargetMinutes: 300},
		{ID: 2, UserID: 1, Category: model.LogCategoryReading, TargetMinutes: 120},
	}
	repo.On("GetByUserID", uint(1)).Return(goals, nil)

	result, err := svc.GetGoals(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

// ============================================================
// GetProgress テスト
// ============================================================

func TestWeeklyGoalGetProgress_Success(t *testing.T) {
	svc, repo := newTestWeeklyGoalService()

	goals := []model.WeeklyGoal{
		{ID: 1, UserID: 1, Category: model.LogCategoryCoding, TargetMinutes: 300},
		{ID: 2, UserID: 1, Category: model.LogCategoryReading, TargetMinutes: 120},
	}
	repo.On("GetByUserID", uint(1)).Return(goals, nil)
	repo.On("SumDurationByUserCategoryThisWeek", uint(1), "coding").Return(150, nil)
	repo.On("SumDurationByUserCategoryThisWeek", uint(1), "reading").Return(120, nil)

	result, err := svc.GetProgress(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 50, result[0].ProgressPercent)  // 150/300 = 50%
	assert.Equal(t, 100, result[1].ProgressPercent) // 120/120 = 100%
	repo.AssertExpectations(t)
}

func TestWeeklyGoalGetProgress_Empty(t *testing.T) {
	svc, repo := newTestWeeklyGoalService()

	repo.On("GetByUserID", uint(1)).Return([]model.WeeklyGoal{}, nil)

	result, err := svc.GetProgress(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestWeeklyGoalGetProgress_RepoError(t *testing.T) {
	svc, repo := newTestWeeklyGoalService()

	repo.On("GetByUserID", uint(1)).Return([]model.WeeklyGoal(nil), errors.New("db error"))

	_, err := svc.GetProgress(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestWeeklyGoalGetProgress_SumError(t *testing.T) {
	svc, repo := newTestWeeklyGoalService()

	goals := []model.WeeklyGoal{
		{ID: 1, UserID: 1, Category: model.LogCategoryCoding, TargetMinutes: 300},
	}
	repo.On("GetByUserID", uint(1)).Return(goals, nil)
	repo.On("SumDurationByUserCategoryThisWeek", uint(1), "coding").Return(0, errors.New("db error"))

	_, err := svc.GetProgress(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}
