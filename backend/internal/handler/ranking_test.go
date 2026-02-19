package handler

import (
	"fmt"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ---------- ContributionRanking ----------

func TestRankingContribution_Success(t *testing.T) {
	h, svc := setupRankingHandler()
	entries := []model.RankingEntry{{UserID: 1, Score: 100}}
	svc.On("ContributionRanking", "weekly").Return(entries, nil)

	r := newRouter(1)
	r.GET("/rankings/contributions", h.ContributionRanking)
	w := doRequest(r, "GET", "/rankings/contributions", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestRankingContribution_WithPeriod(t *testing.T) {
	h, svc := setupRankingHandler()
	entries := []model.RankingEntry{{UserID: 1, Score: 50}}
	svc.On("ContributionRanking", "monthly").Return(entries, nil)

	r := newRouter(1)
	r.GET("/rankings/contributions", h.ContributionRanking)
	w := doRequest(r, "GET", "/rankings/contributions?period=monthly", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestRankingContribution_ServiceError(t *testing.T) {
	h, svc := setupRankingHandler()
	svc.On("ContributionRanking", "weekly").Return([]model.RankingEntry{}, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/rankings/contributions", h.ContributionRanking)
	w := doRequest(r, "GET", "/rankings/contributions", nil)

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}

// ---------- LanguageRanking ----------

func TestRankingLanguage_Success(t *testing.T) {
	h, svc := setupRankingHandler()
	entries := []model.RankingEntry{{UserID: 1, Score: 80}}
	svc.On("LanguageRanking", "Go", "weekly").Return(entries, nil)

	r := newRouter(1)
	r.GET("/rankings/languages/:lang", h.LanguageRanking)
	w := doRequest(r, "GET", "/rankings/languages/Go", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

// ---------- LevelRanking ----------

func TestRankingLevel_Success(t *testing.T) {
	h, svc := setupRankingHandler()
	entries := []model.RankingEntry{{UserID: 1, Score: 500}}
	svc.On("LevelRanking").Return(entries, nil)

	r := newRouter(1)
	r.GET("/rankings/levels", h.LevelRanking)
	w := doRequest(r, "GET", "/rankings/levels", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

// ---------- AvailableLanguages ----------

func TestRankingAvailableLanguages_Success(t *testing.T) {
	h, svc := setupRankingHandler()
	languages := []string{"Go", "TypeScript", "Python"}
	svc.On("AvailableLanguages").Return(languages, nil)

	r := newRouter(1)
	r.GET("/rankings/languages", h.AvailableLanguages)
	w := doRequest(r, "GET", "/rankings/languages", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestRankingAvailableLanguages_ServiceError(t *testing.T) {
	h, svc := setupRankingHandler()
	svc.On("AvailableLanguages").Return([]string{}, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/rankings/languages", h.AvailableLanguages)
	w := doRequest(r, "GET", "/rankings/languages", nil)

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}
