package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
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

// upsertResourceProgressRequest はリソース進捗のUPSERTリクエスト。
type upsertResourceProgressRequest struct {
	ResourceID        uint   `json:"resource_id" binding:"required"`
	Status            string `json:"status" binding:"required"`
	CompletionPercent int    `json:"completion_percent"`
	Note              string `json:"note"`
}

// resourceProgressResponse はリソース進捗レスポンス。
type resourceProgressResponse struct {
	Progress model.ResourceProgress `json:"progress"`
}

// Upsert はリソース進捗を UPSERT（作成/更新）する。
func (h *ResourceProgressHandler) Upsert(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[upsertResourceProgressRequest](c)
	if req == nil {
		return
	}

	result, err := h.upsert.Execute(c.Request.Context(), userID, req.ResourceID, req.Status, req.CompletionPercent, req.Note)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, resourceProgressResponse{Progress: *result})
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

	respondOK(c, resourceProgressResponse{Progress: *progress})
}

// resourceProgressListResponse はリソース進捗一覧レスポンス。
type resourceProgressListResponse struct {
	Progresses []model.ResourceProgress `json:"progresses"`
	Total      int64                    `json:"total"`
	Limit      int                      `json:"limit"`
	Offset     int                      `json:"offset"`
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

	respondOK(c, resourceProgressListResponse{
		Progresses: ensureSlice(progresses),
		Total:      total,
		Limit:      limit,
		Offset:     offset,
	})
}
