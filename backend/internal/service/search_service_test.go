package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostAdvancedSearchRepository は PostAdvancedSearchRepo のテスト用モック。
type MockPostAdvancedSearchRepository struct {
	mock.Mock
}

func (m *MockPostAdvancedSearchRepository) SearchWithFilter(
	query string,
	tags []string,
	sortBy string,
	dateFrom, dateTo *time.Time,
	limit, offset int,
) ([]model.Post, int64, error) {
	args := m.Called(query, tags, sortBy, dateFrom, dateTo, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Post), args.Get(1).(int64), args.Error(2)
}

// ============================================================
// SearchService テスト
// ============================================================

func newTestSearchService() (*SearchService, *MockPostAdvancedSearchRepository) {
	repo := new(MockPostAdvancedSearchRepository)
	svc := NewSearchService(repo)
	return svc, repo
}

// TestSearchPosts_NoFilter_Success はフィルターなしの検索が正常に動作することを確認。
func TestSearchPosts_NoFilter_Success(t *testing.T) {
	svc, repo := newTestSearchService()

	expected := []model.Post{
		{Title: "Go言語入門", Content: "Go言語の基礎を学びます"},
		{Title: "Goルーティン", Content: "並行処理を理解する"},
	}
	repo.On("SearchWithFilter", "Go", []string(nil), "latest", (*time.Time)(nil), (*time.Time)(nil), 20, 0).
		Return(expected, int64(2), nil)

	params := PostSearchParams{Query: "Go", Limit: 20, Offset: 0}
	result, err := svc.SearchPosts(params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Posts, 2)
	repo.AssertExpectations(t)
}

// TestSearchPosts_WithTagFilter_Success はタグフィルターで絞り込めることを確認。
func TestSearchPosts_WithTagFilter_Success(t *testing.T) {
	svc, repo := newTestSearchService()

	expected := []model.Post{
		{Title: "Go入門", Content: "初心者向けGo"},
	}
	tags := []string{"golang", "beginner"}
	repo.On("SearchWithFilter", "Go", tags, "latest", (*time.Time)(nil), (*time.Time)(nil), 20, 0).
		Return(expected, int64(1), nil)

	params := PostSearchParams{
		Query:  "Go",
		Tags:   tags,
		Limit:  20,
		Offset: 0,
	}
	result, err := svc.SearchPosts(params)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
	assert.Len(t, result.Posts, 1)
	repo.AssertExpectations(t)
}

// TestSearchPosts_WithSortPopular_Success は人気順ソートが正常に動作することを確認。
func TestSearchPosts_WithSortPopular_Success(t *testing.T) {
	svc, repo := newTestSearchService()

	expected := []model.Post{
		{Title: "人気記事", LikeCount: 100},
		{Title: "普通の記事", LikeCount: 10},
	}
	repo.On("SearchWithFilter", "記事", []string(nil), "popular", (*time.Time)(nil), (*time.Time)(nil), 20, 0).
		Return(expected, int64(2), nil)

	params := PostSearchParams{
		Query:  "記事",
		SortBy: SearchSortByPopular,
		Limit:  20,
		Offset: 0,
	}
	result, err := svc.SearchPosts(params)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, 100, result.Posts[0].LikeCount)
	repo.AssertExpectations(t)
}

// TestSearchPosts_WithDateRange_Success は日付範囲フィルターが正常に動作することを確認。
func TestSearchPosts_WithDateRange_Success(t *testing.T) {
	svc, repo := newTestSearchService()

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	expected := []model.Post{
		{Title: "1月の記事"},
	}
	repo.On("SearchWithFilter", "記事", []string(nil), "latest", &from, &to, 20, 0).
		Return(expected, int64(1), nil)

	params := PostSearchParams{
		Query:    "記事",
		DateFrom: &from,
		DateTo:   &to,
		Limit:    20,
		Offset:   0,
	}
	result, err := svc.SearchPosts(params)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
	repo.AssertExpectations(t)
}

// TestSearchPosts_EmptyQuery_Error はクエリが空の場合エラーを返すことを確認。
func TestSearchPosts_EmptyQuery_Error(t *testing.T) {
	svc, _ := newTestSearchService()

	params := PostSearchParams{Query: "", Limit: 20}
	result, err := svc.SearchPosts(params)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestSearchPosts_LimitCapped_Success はlimitが100を超えた場合100に制限されることを確認。
func TestSearchPosts_LimitCapped_Success(t *testing.T) {
	svc, repo := newTestSearchService()

	repo.On("SearchWithFilter", "Go", []string(nil), "latest", (*time.Time)(nil), (*time.Time)(nil), 100, 0).
		Return([]model.Post{}, int64(0), nil)

	params := PostSearchParams{Query: "Go", Limit: 200, Offset: 0}
	_, err := svc.SearchPosts(params)

	assert.NoError(t, err)
	repo.AssertCalled(t, "SearchWithFilter", "Go", []string(nil), "latest", (*time.Time)(nil), (*time.Time)(nil), 100, 0)
}

// TestSearchPosts_DefaultLimit_Success はlimitが0の場合デフォルト(20)が使われることを確認。
func TestSearchPosts_DefaultLimit_Success(t *testing.T) {
	svc, repo := newTestSearchService()

	repo.On("SearchWithFilter", "Go", []string(nil), "latest", (*time.Time)(nil), (*time.Time)(nil), 20, 0).
		Return([]model.Post{}, int64(0), nil)

	params := PostSearchParams{Query: "Go"}
	_, err := svc.SearchPosts(params)

	assert.NoError(t, err)
	repo.AssertCalled(t, "SearchWithFilter", "Go", []string(nil), "latest", (*time.Time)(nil), (*time.Time)(nil), 20, 0)
}

// TestSearchPosts_RepositoryError はリポジトリエラーが適切に伝播することを確認。
func TestSearchPosts_RepositoryError(t *testing.T) {
	svc, repo := newTestSearchService()

	repo.On("SearchWithFilter", "Go", []string(nil), "latest", (*time.Time)(nil), (*time.Time)(nil), 20, 0).
		Return(nil, int64(0), errors.New("db error"))

	params := PostSearchParams{Query: "Go", Limit: 20}
	result, err := svc.SearchPosts(params)

	assert.Error(t, err)
	assert.Nil(t, result)
}
