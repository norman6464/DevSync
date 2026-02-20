package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// YouTubeServiceInterface はYouTubeサービスの抽象インターフェース。
type YouTubeServiceInterface interface {
	Search(query, language string) ([]model.YouTubeVideo, bool, error)
	GetRecommendations(userID uint) ([]model.YouTubeVideo, []string, error)
	IsAvailable() bool
}

// YouTubeHandler はYouTube関連のHTTPハンドラ。
type YouTubeHandler struct {
	service YouTubeServiceInterface
}

// NewYouTubeHandler は新しいYouTubeHandlerインスタンスを生成する。
func NewYouTubeHandler(s YouTubeServiceInterface) *YouTubeHandler {
	return &YouTubeHandler{service: s}
}

// Search はキーワードでYouTube動画を検索する。
func (h *YouTubeHandler) Search(c *gin.Context) {
	query := c.Query("q")
	language := c.DefaultQuery("lang", "ja")

	if query == "" {
		respondBadRequest(c, "検索キーワードが必要です")
		return
	}

	videos, cached, err := h.service.Search(query, language)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.YouTubeSearchResponse{
		Videos: ensureSlice(videos),
		Query:  query,
		Cached: cached,
		Total:  len(videos),
	})
}

// Recommend はユーザーのスキルに基づくおすすめ動画を返す。
func (h *YouTubeHandler) Recommend(c *gin.Context) {
	userID := c.GetUint("userID")

	videos, skills, err := h.service.GetRecommendations(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	videos = ensureSlice(videos)
	skills = ensureSlice(skills)

	respondOK(c, dto.YouTubeRecommendResponse{
		Videos:    videos,
		Skills:    skills,
		Available: h.service.IsAvailable(),
	})
}

// Status はYouTube API機能の利用可能状態を返す。
func (h *YouTubeHandler) Status(c *gin.Context) {
	respondOK(c, dto.AvailabilityResponse{Available: h.service.IsAvailable()})
}
