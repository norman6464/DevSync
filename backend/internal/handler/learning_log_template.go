package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningLogTemplateServiceInterface はLearningLogTemplateServiceのインターフェース。
type LearningLogTemplateServiceInterface interface {
	Create(template *model.LearningLogTemplate) error
	GetByID(id, userID uint) (*model.LearningLogTemplate, error)
	GetByUserID(userID uint) ([]model.LearningLogTemplate, error)
	GetDefaultByUserID(userID uint) (*model.LearningLogTemplate, error)
	Update(id, userID uint, name, defaultTitle, defaultContent string, defaultCategory model.LogCategory, defaultDuration *int, isDefault *bool) (*model.LearningLogTemplate, error)
	Delete(id, userID uint) error
	UseTemplate(id, userID uint) (*model.LearningLog, error)
	CountByUserID(userID uint) (int64, error)
}

// LearningLogTemplateHandler は学習ログテンプレート関連のHTTPハンドラ。
type LearningLogTemplateHandler struct {
	service LearningLogTemplateServiceInterface
}

// NewLearningLogTemplateHandler は新しいLearningLogTemplateHandlerインスタンスを生成する。
func NewLearningLogTemplateHandler(s LearningLogTemplateServiceInterface) *LearningLogTemplateHandler {
	return &LearningLogTemplateHandler{service: s}
}

// Create は新しいテンプレートを作成する。
func (h *LearningLogTemplateHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	input := bindJSON[dto.CreateLearningLogTemplateRequest](c)
	if input == nil {
		return
	}

	template := &model.LearningLogTemplate{
		UserID:          userID,
		Name:            input.Name,
		DefaultTitle:    input.DefaultTitle,
		DefaultContent:  input.DefaultContent,
		DefaultCategory: model.LogCategory(input.DefaultCategory),
		DefaultDuration: input.DefaultDuration,
		IsDefault:       input.IsDefault,
	}

	if err := h.service.Create(template); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, template)
}

// GetByID は指定IDのテンプレートを取得する。
func (h *LearningLogTemplateHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	template, err := h.service.GetByID(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, template)
}

// GetByUserID は現在のユーザーのテンプレート一覧を取得する。
func (h *LearningLogTemplateHandler) GetByUserID(c *gin.Context) {
	userID := c.GetUint("userID")

	templates, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(templates))
}

// GetDefault は現在のユーザーのデフォルトテンプレートを取得する。
func (h *LearningLogTemplateHandler) GetDefault(c *gin.Context) {
	userID := c.GetUint("userID")

	template, err := h.service.GetDefaultByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, template)
}

// Update はテンプレートを更新する。
func (h *LearningLogTemplateHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdateLearningLogTemplateRequest](c)
	if input == nil {
		return
	}

	template, err := h.service.Update(id, userID, input.Name, input.DefaultTitle, input.DefaultContent, model.LogCategory(input.DefaultCategory), input.DefaultDuration, input.IsDefault)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, template)
}

// GetMyCount は認証ユーザー自身のテンプレート総数を返す。
func (h *LearningLogTemplateHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.service.CountByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}

// Delete はテンプレートを削除する。
func (h *LearningLogTemplateHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Delete(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// UseTemplate はテンプレートから学習ログを作成する。
func (h *LearningLogTemplateHandler) UseTemplate(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	log, err := h.service.UseTemplate(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, log)
}
