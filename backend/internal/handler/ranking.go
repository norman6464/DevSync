package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// RankingHandler はランキング関連の HTTP ハンドラ。
// コントリビューションランキング・言語別ランキングの取得を処理する。
type RankingHandler struct {
	getContribution *usecase.GetContributionRankingUseCase
	getLanguage     *usecase.GetLanguageRankingUseCase
	getLevel        *usecase.GetLevelRankingUseCase
	listLanguages   *usecase.ListRankingLanguagesUseCase
}

// NewRankingHandler は RankingHandler を生成する。
func NewRankingHandler(
	getContribution *usecase.GetContributionRankingUseCase,
	getLanguage *usecase.GetLanguageRankingUseCase,
	getLevel *usecase.GetLevelRankingUseCase,
	listLanguages *usecase.ListRankingLanguagesUseCase,
) *RankingHandler {
	return &RankingHandler{
		getContribution: getContribution,
		getLanguage:     getLanguage,
		getLevel:        getLevel,
		listLanguages:   listLanguages,
	}
}

// ContributionRanking はコントリビューション数によるランキングを返す。
func (h *RankingHandler) ContributionRanking(c *gin.Context) {
	period := c.DefaultQuery("period", "weekly")
	entries, err := h.getContribution.Execute(c.Request.Context(), period)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(entries))
}

// LanguageRanking は指定言語のランキングを返す。
func (h *RankingHandler) LanguageRanking(c *gin.Context) {
	lang := c.Param("lang")
	period := c.DefaultQuery("period", "weekly")
	entries, err := h.getLanguage.Execute(c.Request.Context(), lang, period)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(entries))
}

// LevelRanking はXP合計に基づくレベルランキングを返す。
func (h *RankingHandler) LevelRanking(c *gin.Context) {
	entries, err := h.getLevel.Execute(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(entries))
}

// AvailableLanguages はランキング対象の利用可能な言語一覧を返す。
func (h *RankingHandler) AvailableLanguages(c *gin.Context) {
	languages, err := h.listLanguages.Execute(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(languages))
}
