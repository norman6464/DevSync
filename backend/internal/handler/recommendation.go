package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// RecommendationHandler はレコメンド関連のHTTPリクエストを処理する。
type RecommendationHandler struct {
	users     *usecase.GetRecommendedUsersUseCase
	posts     *usecase.GetTrendingPostsUseCase
	resources *usecase.GetTrendingResourcesUseCase
}

// NewRecommendationHandler は新しいRecommendationHandlerインスタンスを生成する。
func NewRecommendationHandler(
	users *usecase.GetRecommendedUsersUseCase,
	posts *usecase.GetTrendingPostsUseCase,
	resources *usecase.GetTrendingResourcesUseCase,
) *RecommendationHandler {
	return &RecommendationHandler{users: users, posts: posts, resources: resources}
}

// GetRecommendedUsers はスキルマッチングに基づくおすすめユーザーを返す。
func (h *RecommendationHandler) GetRecommendedUsers(c *gin.Context) {
	userID := c.GetUint("userID")
	users, err := h.users.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(users))
}

// GetTrendingPosts は直近7日間の人気投稿を返す。
func (h *RecommendationHandler) GetTrendingPosts(c *gin.Context) {
	posts, err := h.posts.Execute(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(posts))
}

// GetTrendingResources は直近30日間の人気学習リソースを返す。
func (h *RecommendationHandler) GetTrendingResources(c *gin.Context) {
	resources, err := h.resources.Execute(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(resources))
}
