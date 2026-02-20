package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// ---------- GetRecommendedUsers ----------

func TestRecommendationGetUsers_Success(t *testing.T) {
	h, recRepo, userRepo := setupRecommendationHandlerRepo()
	r := newRouter(1)
	r.GET("/recommendations/users", h.GetRecommendedUsers)

	userRepo.On("FindByID", uint(1)).Return(&model.User{
		SkillsLanguages:  "Go,TypeScript",
		SkillsFrameworks: "React,Gin",
	}, nil)

	recRepo.On("GetRecommendedUsers", uint(1), []string{"Go", "TypeScript", "React", "Gin"}, 10).Return([]model.RecommendedUser{
		{
			User:         model.User{Name: "Alice"},
			CommonSkills: []string{"Go", "React"},
			MatchScore:   2,
		},
		{
			User:         model.User{Name: "Bob"},
			CommonSkills: []string{"TypeScript"},
			MatchScore:   1,
		},
	}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/users", nil)
	assertStatus(t, w, http.StatusOK)

	var users []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &users)
	assert.Len(t, users, 2)
}

func TestRecommendationGetUsers_EmptySkills(t *testing.T) {
	h, _, userRepo := setupRecommendationHandlerRepo()
	r := newRouter(1)
	r.GET("/recommendations/users", h.GetRecommendedUsers)

	// スキル未設定のユーザー
	userRepo.On("FindByID", uint(1)).Return(&model.User{
		SkillsLanguages:  "",
		SkillsFrameworks: "",
	}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/users", nil)
	assertStatus(t, w, http.StatusOK)

	var users []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &users)
	assert.Len(t, users, 0)
}

func TestRecommendationGetUsers_UserNotFound(t *testing.T) {
	h, _, userRepo := setupRecommendationHandlerRepo()
	r := newRouter(1)
	r.GET("/recommendations/users", h.GetRecommendedUsers)

	userRepo.On("FindByID", uint(1)).Return(nil, assert.AnError)

	w := doRequest(r, http.MethodGet, "/recommendations/users", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- GetTrendingPosts ----------

func TestRecommendationGetTrendingPosts_Success(t *testing.T) {
	h, recRepo, _ := setupRecommendationHandlerRepo()
	r := newRouter(1)
	r.GET("/recommendations/posts", h.GetTrendingPosts)

	recRepo.On("GetTrendingPosts", 10, 7).Return([]model.Post{
		{Title: "Popular Post 1", LikeCount: 50, CommentCount: 10},
		{Title: "Popular Post 2", LikeCount: 30, CommentCount: 20},
	}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/posts", nil)
	assertStatus(t, w, http.StatusOK)

	var posts []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &posts)
	assert.Len(t, posts, 2)
}

func TestRecommendationGetTrendingPosts_Empty(t *testing.T) {
	h, recRepo, _ := setupRecommendationHandlerRepo()
	r := newRouter(1)
	r.GET("/recommendations/posts", h.GetTrendingPosts)

	recRepo.On("GetTrendingPosts", 10, 7).Return([]model.Post{}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/posts", nil)
	assertStatus(t, w, http.StatusOK)

	var posts []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &posts)
	assert.Len(t, posts, 0)
}

// ---------- GetTrendingResources ----------

func TestRecommendationGetTrendingResources_Success(t *testing.T) {
	h, recRepo, _ := setupRecommendationHandlerRepo()
	r := newRouter(1)
	r.GET("/recommendations/resources", h.GetTrendingResources)

	recRepo.On("GetTrendingResources", 10, 30).Return([]model.LearningResource{
		{Title: "Popular Resource 1", LikeCount: 40, SaveCount: 15},
		{Title: "Popular Resource 2", LikeCount: 25, SaveCount: 10},
		{Title: "Popular Resource 3", LikeCount: 20, SaveCount: 5},
	}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/resources", nil)
	assertStatus(t, w, http.StatusOK)

	var resources []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resources)
	assert.Len(t, resources, 3)
}

func TestRecommendationGetTrendingPosts_ServiceError(t *testing.T) {
	h, recRepo, _ := setupRecommendationHandlerRepo()
	r := newRouter(1)
	r.GET("/recommendations/posts", h.GetTrendingPosts)

	recRepo.On("GetTrendingPosts", 10, 7).Return([]model.Post{}, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/recommendations/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestRecommendationGetTrendingResources_Empty(t *testing.T) {
	h, recRepo, _ := setupRecommendationHandlerRepo()
	r := newRouter(1)
	r.GET("/recommendations/resources", h.GetTrendingResources)

	recRepo.On("GetTrendingResources", 10, 30).Return([]model.LearningResource{}, nil)

	w := doRequest(r, http.MethodGet, "/recommendations/resources", nil)
	assertStatus(t, w, http.StatusOK)

	var resources []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resources)
	assert.Len(t, resources, 0)
}

func TestRecommendationGetTrendingResources_ServiceError(t *testing.T) {
	h, recRepo, _ := setupRecommendationHandlerRepo()
	r := newRouter(1)
	r.GET("/recommendations/resources", h.GetTrendingResources)

	recRepo.On("GetTrendingResources", 10, 30).Return([]model.LearningResource{}, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/recommendations/resources", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}
