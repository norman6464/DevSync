package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// LearningLogTemplateHandler は学習ログテンプレート関連のHTTPハンドラ。
type LearningLogTemplateHandler struct {
	create     *usecase.CreateLearningLogTemplateUseCase
	get        *usecase.GetLearningLogTemplateUseCase
	list       *usecase.ListLearningLogTemplatesUseCase
	getDefault *usecase.GetDefaultLearningLogTemplateUseCase
	update     *usecase.UpdateLearningLogTemplateUseCase
	remove     *usecase.DeleteLearningLogTemplateUseCase
	createLog  *usecase.CreateLearningLogFromTemplateUseCase
	count      *usecase.CountLearningLogTemplatesUseCase
}

// NewLearningLogTemplateHandler は新しいLearningLogTemplateHandlerインスタンスを生成する。
func NewLearningLogTemplateHandler(
	create *usecase.CreateLearningLogTemplateUseCase,
	get *usecase.GetLearningLogTemplateUseCase,
	list *usecase.ListLearningLogTemplatesUseCase,
	getDefault *usecase.GetDefaultLearningLogTemplateUseCase,
	update *usecase.UpdateLearningLogTemplateUseCase,
	remove *usecase.DeleteLearningLogTemplateUseCase,
	createLog *usecase.CreateLearningLogFromTemplateUseCase,
	count *usecase.CountLearningLogTemplatesUseCase,
) *LearningLogTemplateHandler {
	return &LearningLogTemplateHandler{
		create: create, get: get, list: list, getDefault: getDefault,
		update: update, remove: remove, createLog: createLog, count: count,
	}
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

	if err := h.create.Execute(c.Request.Context(), template); err != nil {
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

	template, err := h.get.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, template)
}

// GetByUserID は現在のユーザーのテンプレート一覧を取得する。
func (h *LearningLogTemplateHandler) GetByUserID(c *gin.Context) {
	userID := c.GetUint("userID")

	templates, err := h.list.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(templates))
}

// GetDefault は現在のユーザーのデフォルトテンプレートを取得する。
func (h *LearningLogTemplateHandler) GetDefault(c *gin.Context) {
	userID := c.GetUint("userID")

	template, err := h.getDefault.Execute(c.Request.Context(), userID)
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

	template, err := h.update.Execute(c.Request.Context(), usecase.UpdateLearningLogTemplateInput{
		ID:              id,
		UserID:          userID,
		Name:            input.Name,
		DefaultTitle:    input.DefaultTitle,
		DefaultContent:  input.DefaultContent,
		DefaultCategory: model.LogCategory(input.DefaultCategory),
		DefaultDuration: input.DefaultDuration,
		IsDefault:       input.IsDefault,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, template)
}

// GetMyCount は認証ユーザー自身のテンプレート総数を返す。
func (h *LearningLogTemplateHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.count.Execute(c.Request.Context(), userID)
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

	if err := h.remove.Execute(c.Request.Context(), id, userID); err != nil {
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

	log, err := h.createLog.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, log)
}
