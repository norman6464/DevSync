package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// YouTubeHandler はYouTube関連のHTTPハンドラ。
type YouTubeHandler struct {
	search       *usecase.SearchYouTubeVideosUseCase
	recommend    *usecase.RecommendYouTubeVideosUseCase
	availability *usecase.CheckYouTubeAvailabilityUseCase
}

// NewYouTubeHandler は新しいYouTubeHandlerインスタンスを生成する。
func NewYouTubeHandler(
	search *usecase.SearchYouTubeVideosUseCase,
	recommend *usecase.RecommendYouTubeVideosUseCase,
	availability *usecase.CheckYouTubeAvailabilityUseCase,
) *YouTubeHandler {
	return &YouTubeHandler{search: search, recommend: recommend, availability: availability}
}

// Search はキーワードでYouTube動画を検索する。
func (h *YouTubeHandler) Search(c *gin.Context) {
	query := c.Query("q")
	language := c.DefaultQuery("lang", "ja")

	if query == "" {
		respondBadRequest(c, "検索キーワードが必要です")
		return
	}

	videos, cached, err := h.search.Execute(c.Request.Context(), query, language)
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

	videos, skills, err := h.recommend.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.YouTubeRecommendResponse{
		Videos:    ensureSlice(videos),
		Skills:    ensureSlice(skills),
		Available: h.availability.Execute(),
	})
}

// Status はYouTube API機能の利用可能状態を返す。
func (h *YouTubeHandler) Status(c *gin.Context) {
	respondOK(c, dto.AvailabilityResponse{Available: h.availability.Execute()})
}
