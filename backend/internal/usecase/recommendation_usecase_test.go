package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockRecommendationRepo は usecase/repository.RecommendationRepository のモック。
type mockRecommendationRepo struct{ mock.Mock }

func (m *mockRecommendationRepo) GetRecommendedUsers(ctx context.Context, userID uint, skills []string, limit int) ([]model.RecommendedUser, error) {
	args := m.Called(ctx, userID, skills, limit)
	u, _ := args.Get(0).([]model.RecommendedUser)
	return u, args.Error(1)
}

func (m *mockRecommendationRepo) GetTrendingPosts(ctx context.Context, limit, days int) ([]model.Post, error) {
	args := m.Called(ctx, limit, days)
	p, _ := args.Get(0).([]model.Post)
	return p, args.Error(1)
}

func (m *mockRecommendationRepo) GetTrendingResources(ctx context.Context, limit, days int) ([]model.LearningResource, error) {
	args := m.Called(ctx, limit, days)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Error(1)
}

// mockUserSkillsReader は usecase/repository.UserSkillsReader のモック。
type mockUserSkillsReader struct{ mock.Mock }

func (m *mockUserSkillsReader) FindByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

// ============================================================
// おすすめユーザー
// ============================================================

func TestGetRecommendedUsersUseCase_Execute(t *testing.T) {
	t.Run("プロフィールのスキルで検索する", func(t *testing.T) {
		recs := new(mockRecommendationRepo)
		users := new(mockUserSkillsReader)
		users.On("FindByID", mock.Anything, uint(1)).
			Return(&model.User{ID: 1, SkillsLanguages: "Go,TypeScript", SkillsFrameworks: "Gin"}, nil)
		recs.On("GetRecommendedUsers", mock.Anything, uint(1), []string{"Go", "TypeScript", "Gin"}, 10).
			Return([]model.RecommendedUser{{MatchScore: 2}}, nil)
		uc := usecase.NewGetRecommendedUsersUseCase(recs, users)

		result, err := uc.Execute(context.Background(), 1)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		recs.AssertExpectations(t)
	})

	t.Run("スキルが無ければ空配列を返し、検索もしない", func(t *testing.T) {
		for _, u := range []*model.User{
			{ID: 1},
			{ID: 1, SkillsLanguages: "  ", SkillsFrameworks: ",,"},
		} {
			recs := new(mockRecommendationRepo)
			users := new(mockUserSkillsReader)
			users.On("FindByID", mock.Anything, uint(1)).Return(u, nil)
			uc := usecase.NewGetRecommendedUsersUseCase(recs, users)

			result, err := uc.Execute(context.Background(), 1)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Empty(t, result)
			recs.AssertNotCalled(t, "GetRecommendedUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		}
	})

	t.Run("不在のユーザーは DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		recs := new(mockRecommendationRepo)
		users := new(mockUserSkillsReader)
		users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetRecommendedUsersUseCase(recs, users)

		_, err := uc.Execute(context.Background(), 1)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
		recs.AssertNotCalled(t, "GetRecommendedUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("ユーザー取得の DB 障害はそのまま伝播する", func(t *testing.T) {
		recs := new(mockRecommendationRepo)
		users := new(mockUserSkillsReader)
		users.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetRecommendedUsersUseCase(recs, users)

		_, err := uc.Execute(context.Background(), 1)

		assert.EqualError(t, err, "db error")
	})

	t.Run("検索の DB 障害はそのまま伝播する", func(t *testing.T) {
		recs := new(mockRecommendationRepo)
		users := new(mockUserSkillsReader)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, SkillsLanguages: "Go"}, nil)
		recs.On("GetRecommendedUsers", mock.Anything, uint(1), []string{"Go"}, 10).
			Return([]model.RecommendedUser(nil), errors.New("db error"))
		uc := usecase.NewGetRecommendedUsersUseCase(recs, users)

		_, err := uc.Execute(context.Background(), 1)

		assert.EqualError(t, err, "db error")
	})
}

func TestParseSkills(t *testing.T) {
	cases := []struct {
		name       string
		languages  string
		frameworks string
		want       []string
	}{
		{"言語とフレームワークを連結する", "Go,Python", "Gin,React", []string{"Go", "Python", "Gin", "React"}},
		{"前後の空白とタブを落とす", " Go , Python\t", "\tGin ", []string{"Go", "Python", "Gin"}},
		{"空の要素は捨てる", "Go,,Python,", ",", []string{"Go", "Python"}},
		{"両方空なら nil", "", "", nil},
		{"空白のみなら nil", "  ", " \t ", nil},
		{"区切りが無ければ 1 件", "Go", "", []string{"Go"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, usecase.ParseSkills(c.languages, c.frameworks))
		})
	}
}

// ============================================================
// 人気投稿 / 人気リソース
// ============================================================

func TestGetTrendingPostsUseCase_Execute(t *testing.T) {
	t.Run("直近 7 日を 10 件で取得する", func(t *testing.T) {
		recs := new(mockRecommendationRepo)
		recs.On("GetTrendingPosts", mock.Anything, 10, 7).
			Return([]model.Post{{ID: 1}}, nil)
		uc := usecase.NewGetTrendingPostsUseCase(recs)

		posts, err := uc.Execute(context.Background())

		require.NoError(t, err)
		assert.Len(t, posts, 1)
		recs.AssertExpectations(t)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		recs := new(mockRecommendationRepo)
		recs.On("GetTrendingPosts", mock.Anything, 10, 7).
			Return([]model.Post(nil), errors.New("db error"))
		uc := usecase.NewGetTrendingPostsUseCase(recs)

		_, err := uc.Execute(context.Background())

		assert.EqualError(t, err, "db error")
	})
}

func TestGetTrendingResourcesUseCase_Execute(t *testing.T) {
	t.Run("直近 30 日を 10 件で取得する", func(t *testing.T) {
		recs := new(mockRecommendationRepo)
		recs.On("GetTrendingResources", mock.Anything, 10, 30).
			Return([]model.LearningResource{{ID: 1}}, nil)
		uc := usecase.NewGetTrendingResourcesUseCase(recs)

		resources, err := uc.Execute(context.Background())

		require.NoError(t, err)
		assert.Len(t, resources, 1)
		recs.AssertExpectations(t)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		recs := new(mockRecommendationRepo)
		recs.On("GetTrendingResources", mock.Anything, 10, 30).
			Return([]model.LearningResource(nil), errors.New("db error"))
		uc := usecase.NewGetTrendingResourcesUseCase(recs)

		_, err := uc.Execute(context.Background())

		assert.EqualError(t, err, "db error")
	})
}
