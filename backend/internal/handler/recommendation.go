package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// RecommendationServiceInterface はRecommendationHandlerが依存するサービスのインターフェース。
type RecommendationServiceInterface interface {
	GetRecommendedUsers(userID uint) ([]model.RecommendedUser, error)
	GetTrendingPosts() ([]model.Post, error)
	GetTrendingResources() ([]model.LearningResource, error)
}

// RecommendationHandler はレコメンド関連のHTTPリクエストを処理する。
type RecommendationHandler struct {
	service RecommendationServiceInterface
}

// NewRecommendationHandler は新しいRecommendationHandlerインスタンスを生成する。
func NewRecommendationHandler(service RecommendationServiceInterface) *RecommendationHandler {
	return &RecommendationHandler{service: service}
}

// GetRecommendedUsers はスキルマッチングに基づくおすすめユーザーを返す。
func (h *RecommendationHandler) GetRecommendedUsers(c *gin.Context) {
	userID := c.GetUint("userID")
	users, err := h.service.GetRecommendedUsers(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, users)
}

// GetTrendingPosts は直近7日間の人気投稿を返す。
func (h *RecommendationHandler) GetTrendingPosts(c *gin.Context) {
	posts, err := h.service.GetTrendingPosts()
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, posts)
}

// GetTrendingResources は直近30日間の人気学習リソースを返す。
func (h *RecommendationHandler) GetTrendingResources(c *gin.Context) {
	resources, err := h.service.GetTrendingResources()
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, resources)
}
