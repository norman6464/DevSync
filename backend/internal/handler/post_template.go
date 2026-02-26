package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// PostTemplateServiceInterface はPostTemplateHandlerが依存するサービスのインターフェース。
type PostTemplateServiceInterface interface {
	Create(tmpl *model.PostTemplate) error
	GetByID(id, userID uint) (*model.PostTemplate, error)
	GetByUserID(userID uint, limit, offset int) ([]model.PostTemplate, int64, error)
	Update(id, userID uint, updates *model.PostTemplate) (*model.PostTemplate, error)
	Delete(id, userID uint) error
}

// PostTemplateHandler は投稿テンプレート関連のHTTPハンドラ。
type PostTemplateHandler struct {
	service PostTemplateServiceInterface
}

// NewPostTemplateHandler は新しいPostTemplateHandlerインスタンスを生成する。
func NewPostTemplateHandler(s PostTemplateServiceInterface) *PostTemplateHandler {
	return &PostTemplateHandler{service: s}
}

// Create は新しい投稿テンプレートを作成する。
func (h *PostTemplateHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.CreatePostTemplateRequest](c)
	if input == nil {
		return
	}

	tmpl := &model.PostTemplate{
		UserID:          userID,
		Name:            input.Name,
		TitleTemplate:   input.TitleTemplate,
		ContentTemplate: input.ContentTemplate,
	}

	if err := h.service.Create(tmpl); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, tmpl)
}

// GetMyTemplates は認証ユーザーの投稿テンプレート一覧を取得する。
func (h *PostTemplateHandler) GetMyTemplates(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	templates, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.PostTemplateListResponse{
		Templates: templates,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetByID は指定IDの投稿テンプレートを取得する。
func (h *PostTemplateHandler) GetByID(c *gin.Context) {
	handleGetByID(c, h.service.GetByID)
}

// Update は指定された投稿テンプレートを更新する。
func (h *PostTemplateHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.UpdatePostTemplateRequest](c)
	if input == nil {
		return
	}

	updates := &model.PostTemplate{}
	if input.Name != nil {
		updates.Name = *input.Name
	}
	if input.TitleTemplate != nil {
		updates.TitleTemplate = *input.TitleTemplate
	}
	if input.ContentTemplate != nil {
		updates.ContentTemplate = *input.ContentTemplate
	}

	tmpl, err := h.service.Update(id, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, tmpl)
}

// Delete は指定された投稿テンプレートを削除する。
func (h *PostTemplateHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
}
