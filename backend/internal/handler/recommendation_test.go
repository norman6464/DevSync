package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockRecommendationRepo は usecase/repository.RecommendationRepository のモック（ctx 付き）。
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

// mockUserSkillsReader は usecase/repository.UserSkillsReader のモック（ctx 付き）。
type mockUserSkillsReader struct{ mock.Mock }

func (m *mockUserSkillsReader) FindByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

// recommendationPorts は RecommendationHandler が使う port モックの束。
type recommendationPorts struct {
	Recommendations *mockRecommendationRepo
	Users           *mockUserSkillsReader
}

// newTestRecommendationHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestRecommendationHandler() (*RecommendationHandler, recommendationPorts) {
	p := recommendationPorts{
		Recommendations: new(mockRecommendationRepo),
		Users:           new(mockUserSkillsReader),
	}
	h := NewRecommendationHandler(
		usecase.NewGetRecommendedUsersUseCase(p.Recommendations, p.Users),
		usecase.NewGetTrendingPostsUseCase(p.Recommendations),
		usecase.NewGetTrendingResourcesUseCase(p.Recommendations),
	)
	return h, p
}

// ============================================================
// おすすめユーザー
// ============================================================

func TestRecommendationHandler_GetRecommendedUsers(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/users", h.GetRecommendedUsers)

	p.Users.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, SkillsLanguages: "Go,TypeScript", SkillsFrameworks: "Gin"}, nil)
	p.Recommendations.On("GetRecommendedUsers", mock.Anything, uint(1), []string{"Go", "TypeScript", "Gin"}, 10).
		Return([]model.RecommendedUser{
			{User: model.User{ID: 2, Username: "gopher"}, CommonSkills: []string{"Go"}, MatchScore: 1},
		}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/users", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "gopher")
	p.Recommendations.AssertExpectations(t)
}

// スキルが登録されていなければ空配列を返し、検索も走らない。
func TestRecommendationHandler_GetRecommendedUsers_EmptySkills(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/users", h.GetRecommendedUsers)

	p.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/users", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
	p.Recommendations.AssertNotCalled(t, "GetRecommendedUsers",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// 0 件でも null ではなく空配列を返す。
func TestRecommendationHandler_GetRecommendedUsers_Empty(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/users", h.GetRecommendedUsers)

	p.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, SkillsLanguages: "Go"}, nil)
	p.Recommendations.On("GetRecommendedUsers", mock.Anything, uint(1), []string{"Go"}, 10).
		Return([]model.RecommendedUser(nil), nil)

	w := doRequest(r, http.MethodGet, "/recommendations/users", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

// 不在のユーザーは 500 になる（移行前から変わらない挙動）。
func TestRecommendationHandler_GetRecommendedUsers_UserNotFoundIs500(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/users", h.GetRecommendedUsers)

	p.Users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/users", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	p.Recommendations.AssertNotCalled(t, "GetRecommendedUsers",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestRecommendationHandler_GetRecommendedUsers_RepositoryError(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/users", h.GetRecommendedUsers)

	p.Users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, SkillsLanguages: "Go"}, nil)
	p.Recommendations.On("GetRecommendedUsers", mock.Anything, uint(1), []string{"Go"}, 10).
		Return([]model.RecommendedUser(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/recommendations/users", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// 人気投稿
// ============================================================

func TestRecommendationHandler_GetTrendingPosts(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/posts", h.GetTrendingPosts)

	p.Recommendations.On("GetTrendingPosts", mock.Anything, 10, 7).
		Return([]model.Post{{ID: 1, Title: "人気投稿"}}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/posts", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "人気投稿")
	p.Recommendations.AssertExpectations(t)
}

func TestRecommendationHandler_GetTrendingPosts_Empty(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/posts", h.GetTrendingPosts)

	p.Recommendations.On("GetTrendingPosts", mock.Anything, 10, 7).Return([]model.Post(nil), nil)

	w := doRequest(r, http.MethodGet, "/recommendations/posts", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestRecommendationHandler_GetTrendingPosts_RepositoryError(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/posts", h.GetTrendingPosts)

	p.Recommendations.On("GetTrendingPosts", mock.Anything, 10, 7).
		Return([]model.Post(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/recommendations/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// 人気学習リソース
// ============================================================

func TestRecommendationHandler_GetTrendingResources(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/resources", h.GetTrendingResources)

	p.Recommendations.On("GetTrendingResources", mock.Anything, 10, 30).
		Return([]model.LearningResource{{ID: 1, Title: "人気リソース"}}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/resources", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "人気リソース")
	p.Recommendations.AssertExpectations(t)
}

func TestRecommendationHandler_GetTrendingResources_Empty(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/resources", h.GetTrendingResources)

	p.Recommendations.On("GetTrendingResources", mock.Anything, 10, 30).
		Return([]model.LearningResource(nil), nil)

	w := doRequest(r, http.MethodGet, "/recommendations/resources", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestRecommendationHandler_GetTrendingResources_RepositoryError(t *testing.T) {
	h, p := newTestRecommendationHandler()
	r := newRouter(1)
	r.GET("/recommendations/resources", h.GetTrendingResources)

	p.Recommendations.On("GetTrendingResources", mock.Anything, 10, 30).
		Return([]model.LearningResource(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/recommendations/resources", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}
