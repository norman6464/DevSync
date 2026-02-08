package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

type RankingHandler struct {
	service *service.RankingService
}

func NewRankingHandler(s *service.RankingService) *RankingHandler {
	return &RankingHandler{service: s}
}

func (h *RankingHandler) ContributionRanking(c *gin.Context) {
	period := c.DefaultQuery("period", "weekly")
	entries, err := h.service.ContributionRanking(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

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

func (h *RankingHandler) AvailableLanguages(c *gin.Context) {
	languages, err := h.service.AvailableLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, languages)
}
