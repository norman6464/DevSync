package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockStreakFreezeRepo は usecase/repository.StreakFreezeRepository のモック。
type mockStreakFreezeRepo struct{ mock.Mock }

func (m *mockStreakFreezeRepo) CreateWithinLimits(ctx context.Context, freeze *model.StreakFreeze, maxPerMonth int) (repository.FreezeUseOutcome, error) {
	args := m.Called(ctx, freeze, maxPerMonth)
	return args.Get(0).(repository.FreezeUseOutcome), args.Error(1)
}

func (m *mockStreakFreezeRepo) GetByUserIDAndMonth(ctx context.Context, userID uint, year, month int) ([]model.StreakFreeze, error) {
	args := m.Called(ctx, userID, year, month)
	f, _ := args.Get(0).([]model.StreakFreeze)
	return f, args.Error(1)
}

func (m *mockStreakFreezeRepo) HasFreezeOnDate(ctx context.Context, userID uint, date string) (bool, error) {
	args := m.Called(ctx, userID, date)
	return args.Bool(0), args.Error(1)
}

// freezeNowParts は usecase が使う「今日」「今年」「今月」を返す。
func freezeNowParts() (string, int, int) {
	now := time.Now()
	return now.Format("2006-01-02"), now.Year(), int(now.Month())
}

// assertDomainCode はエラーが指定の DomainError コードであることを確認する。
func assertDomainCode(t *testing.T, err error, want domain.ErrorCode) {
	t.Helper()
	if !assert.Error(t, err) {
		return
	}
	var de *domain.DomainError
	if assert.ErrorAs(t, err, &de) {
		assert.Equal(t, want, de.Code)
	}
}

func TestUseStreakFreezeUseCase_Execute(t *testing.T) {
	today, year, month := freezeNowParts()

	t.Run("未使用かつ上限未満なら記録する", func(t *testing.T) {
		repo := new(mockStreakFreezeRepo)
		repo.On("CreateWithinLimits", mock.Anything, mock.MatchedBy(func(f *model.StreakFreeze) bool {
			return f.UserID == 1 && f.UsedDate == today && f.Year == year && f.Month == month
		}), model.MaxFreezesPerMonth).Return(repository.FreezeUseCreated, nil)
		uc := usecase.NewUseStreakFreezeUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1))
		repo.AssertExpectations(t)
	})

	t.Run("当日使用済みは Conflict", func(t *testing.T) {
		repo := new(mockStreakFreezeRepo)
		repo.On("CreateWithinLimits", mock.Anything, mock.Anything, model.MaxFreezesPerMonth).
			Return(repository.FreezeUseDuplicateDay, nil)
		uc := usecase.NewUseStreakFreezeUseCase(repo)

		err := uc.Execute(context.Background(), 1)

		assertDomainCode(t, err, domain.ErrCodeConflict)
		var de *domain.DomainError
		if assert.ErrorAs(t, err, &de) {
			assert.Equal(t, "今日は既にフリーズを使用済みです", de.Message)
		}
	})

	t.Run("月次上限に達していれば BadRequest", func(t *testing.T) {
		repo := new(mockStreakFreezeRepo)
		repo.On("CreateWithinLimits", mock.Anything, mock.Anything, model.MaxFreezesPerMonth).
			Return(repository.FreezeUseMonthlyLimitReached, nil)
		uc := usecase.NewUseStreakFreezeUseCase(repo)

		assertDomainCode(t, uc.Execute(context.Background(), 1), domain.ErrCodeBadRequest)
	})

	t.Run("DB 障害を伝播する", func(t *testing.T) {
		repo := new(mockStreakFreezeRepo)
		repo.On("CreateWithinLimits", mock.Anything, mock.Anything, model.MaxFreezesPerMonth).
			Return(repository.FreezeUseCreated, errors.New("db error"))
		uc := usecase.NewUseStreakFreezeUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 1))
	})
}

func TestGetStreakFreezeStatusUseCase_Execute(t *testing.T) {
	today, year, month := freezeNowParts()

	t.Run("使用日一覧と残り回数を組み立てる", func(t *testing.T) {
		repo := new(mockStreakFreezeRepo)
		repo.On("GetByUserIDAndMonth", mock.Anything, uint(1), year, month).
			Return([]model.StreakFreeze{{UsedDate: "2026-02-01"}}, nil)
		repo.On("HasFreezeOnDate", mock.Anything, uint(1), today).Return(false, nil)
		uc := usecase.NewGetStreakFreezeStatusUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, model.MaxFreezesPerMonth, got.MaxFreezes)
		assert.Equal(t, 1, got.UsedFreezes)
		assert.Equal(t, model.MaxFreezesPerMonth-1, got.Remaining)
		assert.Equal(t, []string{"2026-02-01"}, got.UsedDates)
		assert.False(t, got.TodayUsed)
		assert.True(t, got.CanUseToday)
		repo.AssertExpectations(t)
	})

	t.Run("当日使用済みなら CanUseToday は false", func(t *testing.T) {
		repo := new(mockStreakFreezeRepo)
		repo.On("GetByUserIDAndMonth", mock.Anything, uint(1), year, month).
			Return([]model.StreakFreeze{{UsedDate: today}}, nil)
		repo.On("HasFreezeOnDate", mock.Anything, uint(1), today).Return(true, nil)
		uc := usecase.NewGetStreakFreezeStatusUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.True(t, got.TodayUsed)
		assert.False(t, got.CanUseToday)
		repo.AssertExpectations(t)
	})

	// 上限を超えて記録されていても残り回数は負にならない
	t.Run("上限超過でも残り回数は 0 で下限を切る", func(t *testing.T) {
		repo := new(mockStreakFreezeRepo)
		repo.On("GetByUserIDAndMonth", mock.Anything, uint(1), year, month).
			Return(make([]model.StreakFreeze, model.MaxFreezesPerMonth+3), nil)
		repo.On("HasFreezeOnDate", mock.Anything, uint(1), today).Return(false, nil)
		uc := usecase.NewGetStreakFreezeStatusUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 0, got.Remaining)
		assert.False(t, got.CanUseToday)
		repo.AssertExpectations(t)
	})

	t.Run("使用が無ければ空の一覧を返す（nil ではない）", func(t *testing.T) {
		repo := new(mockStreakFreezeRepo)
		repo.On("GetByUserIDAndMonth", mock.Anything, uint(1), year, month).Return([]model.StreakFreeze{}, nil)
		repo.On("HasFreezeOnDate", mock.Anything, uint(1), today).Return(false, nil)
		uc := usecase.NewGetStreakFreezeStatusUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, got.UsedDates)
		assert.Empty(t, got.UsedDates)
		repo.AssertExpectations(t)
	})

	t.Run("当日確認の DB 障害を伝播する", func(t *testing.T) {
		repo := new(mockStreakFreezeRepo)
		repo.On("GetByUserIDAndMonth", mock.Anything, uint(1), year, month).Return([]model.StreakFreeze{}, nil)
		repo.On("HasFreezeOnDate", mock.Anything, uint(1), today).Return(false, errors.New("db error"))
		uc := usecase.NewGetStreakFreezeStatusUseCase(repo)

		_, err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
