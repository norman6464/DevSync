package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockWeeklyGoalRepo は usecase/repository.WeeklyGoalRepository のモック。
type mockWeeklyGoalRepo struct{ mock.Mock }

func (m *mockWeeklyGoalRepo) Upsert(ctx context.Context, goal *model.WeeklyGoal) error {
	return m.Called(ctx, goal).Error(0)
}

func (m *mockWeeklyGoalRepo) GetByUserID(ctx context.Context, userID uint) ([]model.WeeklyGoal, error) {
	args := m.Called(ctx, userID)
	goals, _ := args.Get(0).([]model.WeeklyGoal)
	return goals, args.Error(1)
}

func (m *mockWeeklyGoalRepo) SumDurationByUserCategoryThisWeek(ctx context.Context, userID uint, category string) (int, error) {
	args := m.Called(ctx, userID, category)
	return args.Int(0), args.Error(1)
}

func TestSetWeeklyGoalUseCase_Execute(t *testing.T) {
	t.Run("有効なカテゴリと時間で upsert される", func(t *testing.T) {
		goals := new(mockWeeklyGoalRepo)
		goals.On("Upsert", mock.Anything, mock.MatchedBy(func(g *model.WeeklyGoal) bool {
			return g.UserID == 1 && g.Category == model.LogCategoryCoding && g.TargetMinutes == 300
		})).Return(nil)
		uc := usecase.NewSetWeeklyGoalUseCase(goals)

		goal, err := uc.Execute(context.Background(), 1, "coding", 300)

		assert.NoError(t, err)
		assert.Equal(t, model.LogCategoryCoding, goal.Category)
		assert.Equal(t, 300, goal.TargetMinutes)
		goals.AssertExpectations(t)
	})

	t.Run("無効なカテゴリは 400", func(t *testing.T) {
		goals := new(mockWeeklyGoalRepo)
		uc := usecase.NewSetWeeklyGoalUseCase(goals)

		_, err := uc.Execute(context.Background(), 1, "invalid", 300)

		assert.Error(t, err)
		goals.AssertNotCalled(t, "Upsert")
	})

	t.Run("負の目標時間は 400", func(t *testing.T) {
		goals := new(mockWeeklyGoalRepo)
		uc := usecase.NewSetWeeklyGoalUseCase(goals)

		_, err := uc.Execute(context.Background(), 1, "coding", -1)

		assert.Error(t, err)
		goals.AssertNotCalled(t, "Upsert")
	})

	t.Run("上限超過の目標時間は 400", func(t *testing.T) {
		goals := new(mockWeeklyGoalRepo)
		uc := usecase.NewSetWeeklyGoalUseCase(goals)

		_, err := uc.Execute(context.Background(), 1, "coding", 10081)

		assert.Error(t, err)
		goals.AssertNotCalled(t, "Upsert")
	})

	t.Run("リポジトリエラーは伝播する", func(t *testing.T) {
		goals := new(mockWeeklyGoalRepo)
		goals.On("Upsert", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewSetWeeklyGoalUseCase(goals)

		_, err := uc.Execute(context.Background(), 1, "coding", 300)

		assert.Error(t, err)
	})
}

func TestListWeeklyGoalsUseCase_Execute(t *testing.T) {
	goals := new(mockWeeklyGoalRepo)
	expected := []model.WeeklyGoal{{ID: 1, UserID: 1, Category: model.LogCategoryCoding, TargetMinutes: 300}}
	goals.On("GetByUserID", mock.Anything, uint(1)).Return(expected, nil)
	uc := usecase.NewListWeeklyGoalsUseCase(goals)

	got, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	goals.AssertExpectations(t)
}

func TestGetWeeklyGoalProgressUseCase_Execute(t *testing.T) {
	t.Run("達成率を算出して返す", func(t *testing.T) {
		goals := new(mockWeeklyGoalRepo)
		goals.On("GetByUserID", mock.Anything, uint(1)).Return([]model.WeeklyGoal{
			{Category: model.LogCategoryCoding, TargetMinutes: 300},
		}, nil)
		goals.On("SumDurationByUserCategoryThisWeek", mock.Anything, uint(1), "coding").Return(150, nil)
		uc := usecase.NewGetWeeklyGoalProgressUseCase(goals)

		progress, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, progress, 1)
		assert.Equal(t, 150, progress[0].ActualMinutes)
		assert.Equal(t, 50, progress[0].ProgressPercent)
		goals.AssertExpectations(t)
	})

	t.Run("目標時間 0 のとき達成率は 0", func(t *testing.T) {
		goals := new(mockWeeklyGoalRepo)
		goals.On("GetByUserID", mock.Anything, uint(1)).Return([]model.WeeklyGoal{
			{Category: model.LogCategoryCoding, TargetMinutes: 0},
		}, nil)
		goals.On("SumDurationByUserCategoryThisWeek", mock.Anything, uint(1), "coding").Return(120, nil)
		uc := usecase.NewGetWeeklyGoalProgressUseCase(goals)

		progress, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 0, progress[0].ProgressPercent)
	})

	t.Run("集計エラーは伝播する", func(t *testing.T) {
		goals := new(mockWeeklyGoalRepo)
		goals.On("GetByUserID", mock.Anything, uint(1)).Return([]model.WeeklyGoal{
			{Category: model.LogCategoryCoding, TargetMinutes: 300},
		}, nil)
		goals.On("SumDurationByUserCategoryThisWeek", mock.Anything, uint(1), "coding").Return(0, errors.New("db error"))
		uc := usecase.NewGetWeeklyGoalProgressUseCase(goals)

		_, err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
	})
}
