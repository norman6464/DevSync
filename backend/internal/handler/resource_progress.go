package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// ResourceProgressHandler はリソース進捗関連の HTTP ハンドラ。各操作は 1 責務の usecase に委譲する。
type ResourceProgressHandler struct {
	upsert *usecase.UpsertResourceProgressUseCase
	get    *usecase.GetResourceProgressUseCase
	list   *usecase.ListResourceProgressUseCase
}

// NewResourceProgressHandler は ResourceProgressHandler を生成する。
func NewResourceProgressHandler(
	upsert *usecase.UpsertResourceProgressUseCase,
	get *usecase.GetResourceProgressUseCase,
	list *usecase.ListResourceProgressUseCase,
) *ResourceProgressHandler {
	return &ResourceProgressHandler{upsert: upsert, get: get, list: list}
}

// Upsert はリソース進捗を UPSERT（作成/更新）する。
func (h *ResourceProgressHandler) Upsert(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.UpsertResourceProgressRequest](c)
	if req == nil {
		return
	}

	result, err := h.upsert.Execute(c.Request.Context(), userID, req.ResourceID, req.Status, req.CompletionPercent, req.Note)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ResourceProgressResponse{Progress: *result})
}

// GetByResource は指定リソースの進捗を取得する。
func (h *ResourceProgressHandler) GetByResource(c *gin.Context) {
	userID := c.GetUint("userID")
	resourceID, ok := parseID(c, "id")
	if !ok {
		return
	}

	progress, err := h.get.Execute(c.Request.Context(), userID, resourceID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ResourceProgressResponse{Progress: *progress})
}

// GetMyProgress は認証ユーザーの進捗一覧を取得する。
func (h *ResourceProgressHandler) GetMyProgress(c *gin.Context) {
	userID := c.GetUint("userID")
	status := c.Query("status")
	limit, offset := parseLimitOffset(c)

	progresses, total, err := h.list.Execute(c.Request.Context(), userID, status, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ResourceProgressListResponse{
		Progresses: ensureSlice(progresses),
		Total:      total,
		Limit:      limit,
		Offset:     offset,
	})
}
