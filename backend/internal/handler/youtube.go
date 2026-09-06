package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
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

// youTubeSearchResponse はYouTube動画検索のレスポンス。
type youTubeSearchResponse struct {
	Videos []model.YouTubeVideo `json:"videos"`
	Query  string               `json:"query"`
	Cached bool                 `json:"cached"`
	Total  int                  `json:"total"`
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

	respondOK(c, youTubeSearchResponse{
		Videos: ensureSlice(videos),
		Query:  query,
		Cached: cached,
		Total:  len(videos),
	})
}

// youTubeRecommendResponse はおすすめ動画のレスポンス。
type youTubeRecommendResponse struct {
	Videos    []model.YouTubeVideo `json:"videos"`
	Skills    []string             `json:"skills"`
	Available bool                 `json:"available"`
}

// Recommend はユーザーのスキルに基づくおすすめ動画を返す。
func (h *YouTubeHandler) Recommend(c *gin.Context) {
	userID := c.GetUint("userID")

	videos, skills, err := h.recommend.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, youTubeRecommendResponse{
		Videos:    ensureSlice(videos),
		Skills:    ensureSlice(skills),
		Available: h.availability.Execute(),
	})
}

// availabilityResponse は利用可能状態レスポンス
type availabilityResponse struct {
	Available bool `json:"available"`
}

// Status はYouTube API機能の利用可能状態を返す。
func (h *YouTubeHandler) Status(c *gin.Context) {
	respondOK(c, availabilityResponse{Available: h.availability.Execute()})
}
