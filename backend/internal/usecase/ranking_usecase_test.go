package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockRankingRepo は usecase/repository.RankingRepository のモック。
type mockRankingRepo struct{ mock.Mock }

func (m *mockRankingRepo) ContributionRanking(ctx context.Context, period string) ([]model.RankingEntry, error) {
	args := m.Called(ctx, period)
	e, _ := args.Get(0).([]model.RankingEntry)
	return e, args.Error(1)
}

func (m *mockRankingRepo) LanguageRanking(ctx context.Context, language, period string) ([]model.RankingEntry, error) {
	args := m.Called(ctx, language, period)
	e, _ := args.Get(0).([]model.RankingEntry)
	return e, args.Error(1)
}

func (m *mockRankingRepo) LevelRanking(ctx context.Context) ([]model.RankingEntry, error) {
	args := m.Called(ctx)
	e, _ := args.Get(0).([]model.RankingEntry)
	return e, args.Error(1)
}

func (m *mockRankingRepo) AvailableLanguages(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	l, _ := args.Get(0).([]string)
	return l, args.Error(1)
}

func TestGetContributionRankingUseCase_Execute(t *testing.T) {
	entries := []model.RankingEntry{{UserID: 1, Score: 100}}

	for _, period := range []string{"weekly", "monthly"} {
		t.Run("period="+period+" は repo へ渡す", func(t *testing.T) {
			repo := new(mockRankingRepo)
			repo.On("ContributionRanking", mock.Anything, period).Return(entries, nil)
			uc := usecase.NewGetContributionRankingUseCase(repo)

			got, err := uc.Execute(context.Background(), period)

			assert.NoError(t, err)
			assert.Equal(t, entries, got)
			repo.AssertExpectations(t)
		})
	}

	t.Run("不正な period は 400（repo を呼ばない）", func(t *testing.T) {
		repo := new(mockRankingRepo)
		uc := usecase.NewGetContributionRankingUseCase(repo)

		_, err := uc.Execute(context.Background(), "yearly")

		assert.Error(t, err)
		repo.AssertNotCalled(t, "ContributionRanking")
	})

	t.Run("repo のエラーをそのまま返す", func(t *testing.T) {
		repo := new(mockRankingRepo)
		repo.On("ContributionRanking", mock.Anything, "weekly").
			Return([]model.RankingEntry(nil), errors.New("db error"))
		uc := usecase.NewGetContributionRankingUseCase(repo)

		_, err := uc.Execute(context.Background(), "weekly")

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetLanguageRankingUseCase_Execute(t *testing.T) {
	entries := []model.RankingEntry{{UserID: 1, Score: 80}}

	t.Run("言語と period を repo へ渡す", func(t *testing.T) {
		repo := new(mockRankingRepo)
		repo.On("LanguageRanking", mock.Anything, "Go", "monthly").Return(entries, nil)
		uc := usecase.NewGetLanguageRankingUseCase(repo)

		got, err := uc.Execute(context.Background(), "Go", "monthly")

		assert.NoError(t, err)
		assert.Equal(t, entries, got)
		repo.AssertExpectations(t)
	})

	t.Run("不正な period は 400（repo を呼ばない）", func(t *testing.T) {
		repo := new(mockRankingRepo)
		uc := usecase.NewGetLanguageRankingUseCase(repo)

		_, err := uc.Execute(context.Background(), "Go", "daily")

		assert.Error(t, err)
		repo.AssertNotCalled(t, "LanguageRanking")
	})

	t.Run("言語名が 50 文字超なら 400（repo を呼ばない）", func(t *testing.T) {
		repo := new(mockRankingRepo)
		uc := usecase.NewGetLanguageRankingUseCase(repo)

		_, err := uc.Execute(context.Background(), strings.Repeat("a", 51), "weekly")

		assert.Error(t, err)
		repo.AssertNotCalled(t, "LanguageRanking")
	})

	t.Run("言語名が 50 文字ちょうどなら通す", func(t *testing.T) {
		repo := new(mockRankingRepo)
		lang := strings.Repeat("a", 50)
		repo.On("LanguageRanking", mock.Anything, lang, "weekly").Return(entries, nil)
		uc := usecase.NewGetLanguageRankingUseCase(repo)

		_, err := uc.Execute(context.Background(), lang, "weekly")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetLevelRankingUseCase_Execute(t *testing.T) {
	t.Run("レベルランキングを返す", func(t *testing.T) {
		repo := new(mockRankingRepo)
		entries := []model.RankingEntry{{UserID: 1, Score: 500}}
		repo.On("LevelRanking", mock.Anything).Return(entries, nil)
		uc := usecase.NewGetLevelRankingUseCase(repo)

		got, err := uc.Execute(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, entries, got)
		repo.AssertExpectations(t)
	})

	t.Run("repo のエラーをそのまま返す", func(t *testing.T) {
		repo := new(mockRankingRepo)
		repo.On("LevelRanking", mock.Anything).Return([]model.RankingEntry(nil), errors.New("db error"))
		uc := usecase.NewGetLevelRankingUseCase(repo)

		_, err := uc.Execute(context.Background())

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestListRankingLanguagesUseCase_Execute(t *testing.T) {
	t.Run("言語一覧を返す", func(t *testing.T) {
		repo := new(mockRankingRepo)
		langs := []string{"Go", "Python"}
		repo.On("AvailableLanguages", mock.Anything).Return(langs, nil)
		uc := usecase.NewListRankingLanguagesUseCase(repo)

		got, err := uc.Execute(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, langs, got)
		repo.AssertExpectations(t)
	})

	t.Run("repo のエラーをそのまま返す", func(t *testing.T) {
		repo := new(mockRankingRepo)
		repo.On("AvailableLanguages", mock.Anything).Return([]string(nil), errors.New("db error"))
		uc := usecase.NewListRankingLanguagesUseCase(repo)

		_, err := uc.Execute(context.Background())

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
