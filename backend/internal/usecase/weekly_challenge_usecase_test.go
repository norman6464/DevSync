package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockWeeklyChallengeRepo は usecase/repository.WeeklyChallengeRepository のモック。
// FindByUserAndWeek は「不在」を (nil, nil) で表す契約。
type mockWeeklyChallengeRepo struct{ mock.Mock }

func (m *mockWeeklyChallengeRepo) Create(ctx context.Context, challenge *model.WeeklyChallenge) error {
	return m.Called(ctx, challenge).Error(0)
}

func (m *mockWeeklyChallengeRepo) FindByUserAndWeek(ctx context.Context, userID uint, year, week int) (*model.WeeklyChallenge, error) {
	args := m.Called(ctx, userID, year, week)
	c, _ := args.Get(0).(*model.WeeklyChallenge)
	return c, args.Error(1)
}

func (m *mockWeeklyChallengeRepo) Update(ctx context.Context, challenge *model.WeeklyChallenge) error {
	return m.Called(ctx, challenge).Error(0)
}

func TestGetCurrentWeeklyChallengeUseCase_Execute(t *testing.T) {
	year, week := time.Now().ISOWeek()

	t.Run("既存のチャレンジがあればそのまま返す（生成しない）", func(t *testing.T) {
		repo := new(mockWeeklyChallengeRepo)
		existing := &model.WeeklyChallenge{ID: 7, UserID: 1, Year: year, Week: week, TargetValue: 300}
		repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).Return(existing, nil)
		uc := usecase.NewGetCurrentWeeklyChallengeUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, existing, got)
		repo.AssertNotCalled(t, "Create")
		repo.AssertExpectations(t)
	})

	t.Run("未登録なら今週分をテンプレートから自動生成する", func(t *testing.T) {
		repo := new(mockWeeklyChallengeRepo)
		repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).
			Return((*model.WeeklyChallenge)(nil), nil)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.WeeklyChallenge")).Return(nil)
		uc := usecase.NewGetCurrentWeeklyChallengeUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), got.UserID)
		assert.Equal(t, year, got.Year)
		assert.Equal(t, week, got.Week)
		assert.NotEmpty(t, got.Description)
		assert.Positive(t, got.TargetValue)
		assert.False(t, got.IsCompleted)
		repo.AssertExpectations(t)
	})

	t.Run("DB 障害は伝播し、生成しない", func(t *testing.T) {
		repo := new(mockWeeklyChallengeRepo)
		repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).
			Return((*model.WeeklyChallenge)(nil), errors.New("db error"))
		uc := usecase.NewGetCurrentWeeklyChallengeUseCase(repo)

		_, err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Create")
		repo.AssertExpectations(t)
	})

	t.Run("生成の保存に失敗したらエラーを返す", func(t *testing.T) {
		repo := new(mockWeeklyChallengeRepo)
		repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).
			Return((*model.WeeklyChallenge)(nil), nil)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.WeeklyChallenge")).
			Return(errors.New("db error"))
		uc := usecase.NewGetCurrentWeeklyChallengeUseCase(repo)

		_, err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestUpdateWeeklyChallengeProgressUseCase_Execute(t *testing.T) {
	year, week := time.Now().ISOWeek()

	t.Run("目標未達なら完了フラグを立てない", func(t *testing.T) {
		repo := new(mockWeeklyChallengeRepo)
		challenge := &model.WeeklyChallenge{UserID: 1, Year: year, Week: week, TargetValue: 300}
		repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).Return(challenge, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(c *model.WeeklyChallenge) bool {
			return c.CurrentValue == 200 && !c.IsCompleted && c.CompletedAt == nil
		})).Return(nil)
		uc := usecase.NewUpdateWeeklyChallengeProgressUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 200)

		assert.NoError(t, err)
		assert.False(t, got.IsCompleted)
		repo.AssertExpectations(t)
	})

	t.Run("目標到達で完了フラグと完了日時を立てる", func(t *testing.T) {
		repo := new(mockWeeklyChallengeRepo)
		challenge := &model.WeeklyChallenge{UserID: 1, Year: year, Week: week, TargetValue: 300}
		repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).Return(challenge, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(c *model.WeeklyChallenge) bool {
			return c.CurrentValue == 300 && c.IsCompleted && c.CompletedAt != nil
		})).Return(nil)
		uc := usecase.NewUpdateWeeklyChallengeProgressUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 300)

		assert.NoError(t, err)
		assert.True(t, got.IsCompleted)
		assert.NotNil(t, got.CompletedAt)
		repo.AssertExpectations(t)
	})

	t.Run("完了済みなら完了日時を上書きしない", func(t *testing.T) {
		repo := new(mockWeeklyChallengeRepo)
		completedAt := time.Now().Add(-24 * time.Hour)
		challenge := &model.WeeklyChallenge{
			UserID: 1, Year: year, Week: week, TargetValue: 300,
			IsCompleted: true, CompletedAt: &completedAt,
		}
		repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).Return(challenge, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.WeeklyChallenge")).Return(nil)
		uc := usecase.NewUpdateWeeklyChallengeProgressUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 400)

		assert.NoError(t, err)
		assert.True(t, got.IsCompleted)
		assert.Equal(t, completedAt, *got.CompletedAt)
		repo.AssertExpectations(t)
	})

	t.Run("今週分が未登録ならエラー（更新しない）", func(t *testing.T) {
		repo := new(mockWeeklyChallengeRepo)
		repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).
			Return((*model.WeeklyChallenge)(nil), nil)
		uc := usecase.NewUpdateWeeklyChallengeProgressUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 100)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
		repo.AssertExpectations(t)
	})
}
