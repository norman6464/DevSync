package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// PostTemplateHandler は投稿テンプレート関連の HTTP ハンドラ。
type PostTemplateHandler struct {
	create *usecase.CreatePostTemplateUseCase
	get    *usecase.GetPostTemplateUseCase
	list   *usecase.ListPostTemplatesUseCase
	update *usecase.UpdatePostTemplateUseCase
	delete *usecase.DeletePostTemplateUseCase
}

// NewPostTemplateHandler は PostTemplateHandler を生成する。
func NewPostTemplateHandler(
	create *usecase.CreatePostTemplateUseCase,
	get *usecase.GetPostTemplateUseCase,
	list *usecase.ListPostTemplatesUseCase,
	update *usecase.UpdatePostTemplateUseCase,
	deleteUC *usecase.DeletePostTemplateUseCase,
) *PostTemplateHandler {
	return &PostTemplateHandler{
		create: create,
		get:    get,
		list:   list,
		update: update,
		delete: deleteUC,
	}
}

// createPostTemplateRequest は投稿テンプレート作成リクエスト。
type createPostTemplateRequest struct {
	Name            string `json:"name" binding:"required,min=1,max=100"`
	TitleTemplate   string `json:"title_template" binding:"omitempty,max=200"`
	ContentTemplate string `json:"content_template" binding:"required,min=1,max=50000"`
}

// Create は新しい投稿テンプレートを作成する。
func (h *PostTemplateHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[createPostTemplateRequest](c)
	if input == nil {
		return
	}

	tmpl := &model.PostTemplate{
		UserID:          userID,
		Name:            input.Name,
		TitleTemplate:   input.TitleTemplate,
		ContentTemplate: input.ContentTemplate,
	}

	if err := h.create.Execute(c.Request.Context(), tmpl); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, tmpl)
}

// postTemplateListResponse は投稿テンプレート一覧レスポンス（ページネーション付き）。
type postTemplateListResponse struct {
	Templates []model.PostTemplate `json:"templates"`
	Total     int64                `json:"total"`
	Limit     int                  `json:"limit"`
	Offset    int                  `json:"offset"`
}

// GetMyTemplates は認証ユーザーの投稿テンプレート一覧を取得する。
func (h *PostTemplateHandler) GetMyTemplates(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	templates, total, err := h.list.Execute(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, postTemplateListResponse{
		Templates: templates,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetByID は指定IDの投稿テンプレートを取得する。
func (h *PostTemplateHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	handleGetByID(c, func(id, userID uint) (*model.PostTemplate, error) {
		return h.get.Execute(ctx, id, userID)
	})
}

// updatePostTemplateRequest は投稿テンプレート更新リクエスト。
type updatePostTemplateRequest struct {
	Name            *string `json:"name" binding:"omitempty,max=100"`
	TitleTemplate   *string `json:"title_template" binding:"omitempty,max=200"`
	ContentTemplate *string `json:"content_template" binding:"omitempty,max=50000"`
}

// Update は指定された投稿テンプレートを更新する。
func (h *PostTemplateHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[updatePostTemplateRequest](c)
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

	tmpl, err := h.update.Execute(c.Request.Context(), id, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, tmpl)
}

// Delete は指定された投稿テンプレートを削除する。
func (h *PostTemplateHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	handleDelete(c, func(id, userID uint) error {
		return h.delete.Execute(ctx, id, userID)
	})
}
