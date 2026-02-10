package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// RankingHandler はランキング関連のHTTPハンドラ。
// コントリビューションランキング・言語別ランキングの取得を処理する。
type RankingHandler struct {
	service *service.RankingService
}

// NewRankingHandler は新しいRankingHandlerインスタンスを生成する。
func NewRankingHandler(s *service.RankingService) *RankingHandler {
	return &RankingHandler{service: s}
}

// ContributionRanking はコントリビューション数によるランキングを返す。
func (h *RankingHandler) ContributionRanking(c *gin.Context) {
	period := c.DefaultQuery("period", "weekly")
	entries, err := h.service.ContributionRanking(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

// LanguageRanking は指定言語のランキングを返す。
func (h *RankingHandler) LanguageRanking(c *gin.Context) {
	lang := c.Param("lang")
	period := c.DefaultQuery("period", "weekly")
	entries, err := h.service.LanguageRanking(lang, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

// LevelRanking はXP合計に基づくレベルランキングを返す。
func (h *RankingHandler) LevelRanking(c *gin.Context) {
	entries, err := h.service.LevelRanking()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

// AvailableLanguages はランキング対象の利用可能な言語一覧を返す。
func (h *RankingHandler) AvailableLanguages(c *gin.Context) {
	languages, err := h.service.AvailableLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, languages)
}
