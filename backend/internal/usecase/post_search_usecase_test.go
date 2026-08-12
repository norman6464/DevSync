package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockPostSearchRepo は usecase/repository.PostSearchRepository のモック。
type mockPostSearchRepo struct{ mock.Mock }

func (m *mockPostSearchRepo) SearchWithFilter(ctx context.Context, params model.PostSearchParams) ([]model.Post, int64, error) {
	args := m.Called(ctx, params)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Get(1).(int64), args.Error(2)
}

func TestSearchPostsUseCase(t *testing.T) {
	t.Run("検索結果と件数を返す", func(t *testing.T) {
		posts := new(mockPostSearchRepo)
		posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
			return p.Query == "go" && p.Limit == 10 && p.Offset == 5
		})).Return([]model.Post{{ID: 1}}, int64(3), nil)

		got, err := usecase.NewSearchPostsUseCase(posts).Execute(context.Background(), model.PostSearchParams{
			Query: "go", Limit: 10, Offset: 5, SortBy: model.SearchSortByLatest,
		})
		require.NoError(t, err)
		assert.Len(t, got.Posts, 1)
		assert.Equal(t, int64(3), got.Total)
		assert.Equal(t, 10, got.Limit)
		assert.Equal(t, 5, got.Offset)
		posts.AssertExpectations(t)
	})

	t.Run("検索クエリが空なら 400", func(t *testing.T) {
		posts := new(mockPostSearchRepo)

		_, err := usecase.NewSearchPostsUseCase(posts).Execute(context.Background(), model.PostSearchParams{})
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
		posts.AssertNotCalled(t, "SearchWithFilter", mock.Anything, mock.Anything)
	})

	t.Run("limit 未指定なら 20 件にする", func(t *testing.T) {
		posts := new(mockPostSearchRepo)
		posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
			return p.Limit == 20
		})).Return([]model.Post{}, int64(0), nil)

		got, err := usecase.NewSearchPostsUseCase(posts).Execute(context.Background(), model.PostSearchParams{Query: "go"})
		require.NoError(t, err)
		assert.Equal(t, 20, got.Limit)
	})

	t.Run("limit は 100 で頭打ちにする", func(t *testing.T) {
		posts := new(mockPostSearchRepo)
		posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
			return p.Limit == 100
		})).Return([]model.Post{}, int64(0), nil)

		got, err := usecase.NewSearchPostsUseCase(posts).Execute(context.Background(), model.PostSearchParams{Query: "go", Limit: 500})
		require.NoError(t, err)
		assert.Equal(t, 100, got.Limit)
	})

	t.Run("ソート順が空なら最新順にする", func(t *testing.T) {
		posts := new(mockPostSearchRepo)
		posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
			return p.SortBy == model.SearchSortByLatest
		})).Return([]model.Post{}, int64(0), nil)

		_, err := usecase.NewSearchPostsUseCase(posts).Execute(context.Background(), model.PostSearchParams{Query: "go"})
		require.NoError(t, err)
		posts.AssertExpectations(t)
	})

	t.Run("無効なソート順は 400", func(t *testing.T) {
		posts := new(mockPostSearchRepo)

		_, err := usecase.NewSearchPostsUseCase(posts).Execute(context.Background(), model.PostSearchParams{
			Query: "go", SortBy: model.SearchSortBy("unknown"),
		})
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
		posts.AssertNotCalled(t, "SearchWithFilter", mock.Anything, mock.Anything)
	})

	t.Run("タグと日付範囲はそのまま渡す", func(t *testing.T) {
		posts := new(mockPostSearchRepo)
		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
		posts.On("SearchWithFilter", mock.Anything, mock.MatchedBy(func(p model.PostSearchParams) bool {
			return len(p.Tags) == 2 && p.DateFrom.Equal(from) && p.DateTo.Equal(to)
		})).Return([]model.Post{}, int64(0), nil)

		_, err := usecase.NewSearchPostsUseCase(posts).Execute(context.Background(), model.PostSearchParams{
			Query: "go", Tags: []string{"golang", "beginner"}, DateFrom: &from, DateTo: &to,
		})
		require.NoError(t, err)
		posts.AssertExpectations(t)
	})

	t.Run("検索エラーはそのまま返す", func(t *testing.T) {
		posts := new(mockPostSearchRepo)
		dbErr := errors.New("db error")
		posts.On("SearchWithFilter", mock.Anything, mock.Anything).Return(nil, int64(0), dbErr)

		_, err := usecase.NewSearchPostsUseCase(posts).Execute(context.Background(), model.PostSearchParams{Query: "go"})
		assert.ErrorIs(t, err, dbErr)
	})
}
