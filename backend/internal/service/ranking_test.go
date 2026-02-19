package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRankingService_ContributionRanking(t *testing.T) {
	t.Run("正常にコントリビューションランキングを取得", func(t *testing.T) {
		repo := new(MockRankingRepository)
		svc := NewRankingService(repo)

		expected := []model.RankingEntry{
			{UserID: 1, Name: "alice", Score: 100},
			{UserID: 2, Name: "bob", Score: 80},
		}
		repo.On("ContributionRanking", "weekly").Return(expected, nil)

		result, err := svc.ContributionRanking("weekly")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		repo := new(MockRankingRepository)
		svc := NewRankingService(repo)

		repo.On("ContributionRanking", "monthly").Return([]model.RankingEntry(nil), errors.New("db error"))

		result, err := svc.ContributionRanking("monthly")
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestRankingService_LanguageRanking(t *testing.T) {
	t.Run("正常に言語別ランキングを取得", func(t *testing.T) {
		repo := new(MockRankingRepository)
		svc := NewRankingService(repo)

		expected := []model.RankingEntry{
			{UserID: 1, Name: "alice", Score: 50},
		}
		repo.On("LanguageRanking", "Go", "weekly").Return(expected, nil)

		result, err := svc.LanguageRanking("Go", "weekly")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		repo := new(MockRankingRepository)
		svc := NewRankingService(repo)

		repo.On("LanguageRanking", "Python", "monthly").Return([]model.RankingEntry(nil), errors.New("db error"))

		result, err := svc.LanguageRanking("Python", "monthly")
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestRankingService_LevelRanking(t *testing.T) {
	t.Run("正常にレベルランキングを取得", func(t *testing.T) {
		repo := new(MockRankingRepository)
		svc := NewRankingService(repo)

		expected := []model.RankingEntry{
			{UserID: 3, Name: "charlie", Score: 200},
			{UserID: 1, Name: "alice", Score: 150},
		}
		repo.On("LevelRanking").Return(expected, nil)

		result, err := svc.LevelRanking()
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		repo := new(MockRankingRepository)
		svc := NewRankingService(repo)

		repo.On("LevelRanking").Return([]model.RankingEntry(nil), errors.New("db error"))

		result, err := svc.LevelRanking()
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestRankingService_AvailableLanguages(t *testing.T) {
	t.Run("正常に利用可能言語一覧を取得", func(t *testing.T) {
		repo := new(MockRankingRepository)
		svc := NewRankingService(repo)

		expected := []string{"Go", "Python", "TypeScript"}
		repo.On("AvailableLanguages").Return(expected, nil)

		result, err := svc.AvailableLanguages()
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		repo := new(MockRankingRepository)
		svc := NewRankingService(repo)

		repo.On("AvailableLanguages").Return([]string(nil), errors.New("db error"))

		result, err := svc.AvailableLanguages()
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestRankingService_ContributionRanking_EmptyList(t *testing.T) {
	repo := new(MockRankingRepository)
	svc := NewRankingService(repo)

	repo.On("ContributionRanking", "weekly").Return([]model.RankingEntry{}, nil)

	result, err := svc.ContributionRanking("weekly")
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestRankingService_LanguageRanking_EmptyList(t *testing.T) {
	repo := new(MockRankingRepository)
	svc := NewRankingService(repo)

	repo.On("LanguageRanking", "Rust", "weekly").Return([]model.RankingEntry{}, nil)

	result, err := svc.LanguageRanking("Rust", "weekly")
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestRankingService_LevelRanking_EmptyList(t *testing.T) {
	repo := new(MockRankingRepository)
	svc := NewRankingService(repo)

	repo.On("LevelRanking").Return([]model.RankingEntry{}, nil)

	result, err := svc.LevelRanking()
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestRankingService_AvailableLanguages_EmptyList(t *testing.T) {
	repo := new(MockRankingRepository)
	svc := NewRankingService(repo)

	repo.On("AvailableLanguages").Return([]string{}, nil)

	result, err := svc.AvailableLanguages()
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestRankingService_ContributionRanking_LargeScores(t *testing.T) {
	repo := new(MockRankingRepository)
	svc := NewRankingService(repo)

	expected := []model.RankingEntry{
		{UserID: 1, Name: "top-user", Score: 999999},
		{UserID: 2, Name: "second-user", Score: 0},
	}
	repo.On("ContributionRanking", "monthly").Return(expected, nil)

	result, err := svc.ContributionRanking("monthly")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(999999), result[0].Score)
	assert.Equal(t, int64(0), result[1].Score)
	repo.AssertExpectations(t)
}
