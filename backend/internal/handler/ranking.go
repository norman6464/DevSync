package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type RankingHandler struct {
	repo *repository.RankingRepository
}

func NewRankingHandler(repo *repository.RankingRepository) *RankingHandler {
	return &RankingHandler{repo: repo}
}

func (h *RankingHandler) ContributionRanking(c *gin.Context) {
	period := c.DefaultQuery("period", "weekly")
	entries, err := h.repo.ContributionRanking(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *RankingHandler) LanguageRanking(c *gin.Context) {
	lang := c.Param("lang")
	period := c.DefaultQuery("period", "weekly")
	entries, err := h.repo.LanguageRanking(lang, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *RankingHandler) AvailableLanguages(c *gin.Context) {
	languages, err := h.repo.AvailableLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, languages)
}
