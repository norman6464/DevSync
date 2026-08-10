package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockRankingRepo は usecase/repository.RankingRepository のモック（ctx 付き）。
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

// setupRankingHandler は本物の usecase と port モックで RankingHandler を組む。
func setupRankingHandler() (*RankingHandler, *mockRankingRepo) {
	repo := new(mockRankingRepo)
	h := NewRankingHandler(
		usecase.NewGetContributionRankingUseCase(repo),
		usecase.NewGetLanguageRankingUseCase(repo),
		usecase.NewGetLevelRankingUseCase(repo),
		usecase.NewListRankingLanguagesUseCase(repo),
	)
	return h, repo
}

// ---------- ContributionRanking ----------

func TestRankingContribution_Success(t *testing.T) {
	h, repo := setupRankingHandler()
	entries := []model.RankingEntry{{UserID: 1, Score: 100}}
	repo.On("ContributionRanking", mock.Anything, "weekly").Return(entries, nil)

	r := newRouter(1)
	r.GET("/rankings/contributions", h.ContributionRanking)
	w := doRequest(r, "GET", "/rankings/contributions", nil)

	assertStatus(t, w, 200)
	repo.AssertExpectations(t)
}

func TestRankingContribution_WithPeriod(t *testing.T) {
	h, repo := setupRankingHandler()
	entries := []model.RankingEntry{{UserID: 1, Score: 50}}
	repo.On("ContributionRanking", mock.Anything, "monthly").Return(entries, nil)

	r := newRouter(1)
	r.GET("/rankings/contributions", h.ContributionRanking)
	w := doRequest(r, "GET", "/rankings/contributions?period=monthly", nil)

	assertStatus(t, w, 200)
	repo.AssertExpectations(t)
}

// period が weekly / monthly 以外なら repo を呼ばず 400 を返す。
func TestRankingContribution_InvalidPeriod(t *testing.T) {
	h, repo := setupRankingHandler()

	r := newRouter(1)
	r.GET("/rankings/contributions", h.ContributionRanking)
	w := doRequest(r, "GET", "/rankings/contributions?period=yearly", nil)

	assertStatus(t, w, 400)
	repo.AssertNotCalled(t, "ContributionRanking")
}

func TestRankingContribution_ServiceError(t *testing.T) {
	h, repo := setupRankingHandler()
	repo.On("ContributionRanking", mock.Anything, "weekly").Return([]model.RankingEntry{}, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/rankings/contributions", h.ContributionRanking)
	w := doRequest(r, "GET", "/rankings/contributions", nil)

	assertStatus(t, w, 500)
	repo.AssertExpectations(t)
}

// ---------- LanguageRanking ----------

func TestRankingLanguage_Success(t *testing.T) {
	h, repo := setupRankingHandler()
	entries := []model.RankingEntry{{UserID: 1, Score: 80}}
	repo.On("LanguageRanking", mock.Anything, "Go", "weekly").Return(entries, nil)

	r := newRouter(1)
	r.GET("/rankings/languages/:lang", h.LanguageRanking)
	w := doRequest(r, "GET", "/rankings/languages/Go", nil)

	assertStatus(t, w, 200)
	repo.AssertExpectations(t)
}

// ---------- LevelRanking ----------

func TestRankingLevel_Success(t *testing.T) {
	h, repo := setupRankingHandler()
	entries := []model.RankingEntry{{UserID: 1, Score: 500}}
	repo.On("LevelRanking", mock.Anything).Return(entries, nil)

	r := newRouter(1)
	r.GET("/rankings/levels", h.LevelRanking)
	w := doRequest(r, "GET", "/rankings/levels", nil)

	assertStatus(t, w, 200)
	repo.AssertExpectations(t)
}

func TestRankingLevel_ServiceError(t *testing.T) {
	h, repo := setupRankingHandler()
	repo.On("LevelRanking", mock.Anything).Return([]model.RankingEntry{}, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/rankings/levels", h.LevelRanking)
	w := doRequest(r, "GET", "/rankings/levels", nil)

	assertStatus(t, w, 500)
	repo.AssertExpectations(t)
}

func TestRankingLanguage_ServiceError(t *testing.T) {
	h, repo := setupRankingHandler()
	repo.On("LanguageRanking", mock.Anything, "Go", "weekly").Return([]model.RankingEntry{}, fmt.Errorf("db error"))

	r := newRouter(1)
	r.GET("/rankings/languages/:lang", h.LanguageRanking)
	w := doRequest(r, "GET", "/rankings/languages/Go", nil)

	assertStatus(t, w, 500)
	repo.AssertExpectations(t)
}

// ---------- AvailableLanguages ----------

func TestRankingAvailableLanguages_Success(t *testing.T) {
	h, repo := setupRankingHandler()
	languages := []string{"Go", "TypeScript", "Python"}
	repo.On("AvailableLanguages", mock.Anything).Return(languages, nil)

	r := newRouter(1)
	r.GET("/rankings/languages", h.AvailableLanguages)
	w := doRequest(r, "GET", "/rankings/languages", nil)

	assertStatus(t, w, 200)
	repo.AssertExpectations(t)
}

func TestRankingAvailableLanguages_ServiceError(t *testing.T) {
	h, repo := setupRankingHandler()
	repo.On("AvailableLanguages", mock.Anything).Return([]string{}, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/rankings/languages", h.AvailableLanguages)
	w := doRequest(r, "GET", "/rankings/languages", nil)

	assertStatus(t, w, 500)
	repo.AssertExpectations(t)
}
