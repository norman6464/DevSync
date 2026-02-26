package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// ResourceProgressServiceInterface はリソース進捗サービスの抽象インターフェース。
type ResourceProgressServiceInterface interface {
	UpsertProgress(userID, resourceID uint, status string, completionPercent int, note string) (*model.ResourceProgress, error)
	GetProgress(userID, resourceID uint) (*model.ResourceProgress, error)
	GetProgressList(userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error)
}

// ResourceProgressHandler はリソース進捗関連のHTTPハンドラ。
type ResourceProgressHandler struct {
	service ResourceProgressServiceInterface
}

// NewResourceProgressHandler は新しいResourceProgressHandlerインスタンスを生成する。
func NewResourceProgressHandler(s ResourceProgressServiceInterface) *ResourceProgressHandler {
	return &ResourceProgressHandler{service: s}
}

// Upsert はリソース進捗をUPSERT（作成/更新）する。
func (h *ResourceProgressHandler) Upsert(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.UpsertResourceProgressRequest](c)
	if req == nil {
		return
	}

	result, err := h.service.UpsertProgress(userID, req.ResourceID, req.Status, req.CompletionPercent, req.Note)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ResourceProgressResponse{Progress: *result})
}

// GetByResource は指定リソースの進捗を取得する。
func (h *ResourceProgressHandler) GetByResource(c *gin.Context) {
	userID := c.GetUint("userID")
	resourceID, ok := parseID(c, "resourceId")
	if !ok {
		return
	}

	progress, err := h.service.GetProgress(userID, resourceID)
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

	progresses, total, err := h.service.GetProgressList(userID, status, limit, offset)
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
