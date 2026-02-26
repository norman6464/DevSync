package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// RoadmapServiceInterface はRoadmapHandlerが依存するサービスメソッドを定義する。
type RoadmapServiceInterface interface {
	Create(roadmap *model.Roadmap) error
	GetByID(id, userID uint) (*model.Roadmap, error)
	GetByUserID(userID uint, limit, offset int) ([]model.Roadmap, int64, error)
	GetByStatus(userID uint, status string) ([]model.Roadmap, error)
	GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error)
	Update(id, userID uint, updates *model.Roadmap) (*model.Roadmap, error)
	UpdateVisibility(id, userID uint, isPublic bool) (*model.Roadmap, error)
	Delete(id, userID uint) error
	CopyRoadmap(roadmapID, userID uint) (*model.Roadmap, error)
	GetTemplates() ([]model.Roadmap, error)
	CreateFromTemplate(templateID, userID uint) (*model.Roadmap, error)
	CreateStep(roadmapID, userID uint, step *model.RoadmapStep) error
	UpdateStep(roadmapID, stepID, userID uint, updates *model.RoadmapStep) (*model.RoadmapStep, error)
	UpdateStepCompletion(roadmapID, stepID, userID uint, isCompleted bool) (*model.RoadmapStep, error)
	BatchCompleteSteps(roadmapID, userID uint, stepIDs []uint) (*model.Roadmap, error)
	DeleteStep(roadmapID, stepID, userID uint) error
	ReorderSteps(roadmapID, userID uint, orders []model.StepOrder) error
	GetStats(userID uint) (*model.RoadmapStats, error)
	CountByUserID(userID uint) (int64, error)
}

// RoadmapHandler はロードマップ関連のHTTPハンドラ。
// ロードマップとステップのCRUD・公開一覧・コピー・並べ替えを処理する。
type RoadmapHandler struct {
	service RoadmapServiceInterface
}

// NewRoadmapHandler は新しいRoadmapHandlerインスタンスを生成する。
func NewRoadmapHandler(s RoadmapServiceInterface) *RoadmapHandler {
	return &RoadmapHandler{service: s}
}

// === ロードマップエンドポイント ===

// Create は新しいロードマップを作成する。
func (h *RoadmapHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.CreateRoadmapRequest](c)
	if input == nil {
		return
	}

	roadmap := &model.Roadmap{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		Category:    model.RoadmapCategory(input.Category),
		IsPublic:    input.IsPublic,
		Status:      model.RoadmapStatusActive,
	}

	if input.Category == "" {
		roadmap.Category = model.RoadmapCategoryOther
	}

	if err := h.service.Create(roadmap); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, roadmap)
}

// GetMyRoadmaps は現在のユーザーのロードマップ一覧を取得する。
func (h *RoadmapHandler) GetMyRoadmaps(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	roadmaps, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.RoadmapListResponse{
		Roadmaps: roadmaps,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	})
}

// GetByStatus はステータス別にロードマップを取得する。
func (h *RoadmapHandler) GetByStatus(c *gin.Context) {
	userID := c.GetUint("userID")
	status := c.Param("status")

	roadmaps, err := h.service.GetByStatus(userID, status)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(roadmaps))
}

// GetPublicRoadmaps は公開ロードマップの一覧をページネーション付きで取得する。
func (h *RoadmapHandler) GetPublicRoadmaps(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	roadmaps, total, err := h.service.GetPublicRoadmaps(limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.RoadmapListResponse{
		Roadmaps: roadmaps,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	})
}

// GetByID は指定IDのロードマップをステップ付きで取得する。
func (h *RoadmapHandler) GetByID(c *gin.Context) {
	handleGetByID(c, h.service.GetByID)
}

// Update は指定IDのロードマップを更新する。
func (h *RoadmapHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.UpdateRoadmapRequest](c)
	if input == nil {
		return
	}

	updates := &model.Roadmap{}
	if input.Title != nil {
		updates.Title = *input.Title
	}
	if input.Description != nil {
		updates.Description = *input.Description
	}
	if input.Category != nil {
		updates.Category = model.RoadmapCategory(*input.Category)
	}
	if input.Status != nil {
		updates.Status = model.RoadmapStatus(*input.Status)
	}

	roadmap, err := h.service.Update(roadmapID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	// IsPublicが指定されている場合は別途処理する
	if input.IsPublic != nil {
		roadmap, err = h.service.UpdateVisibility(roadmapID, userID, *input.IsPublic)
		if err != nil {
			respondError(c, err)
			return
		}
	}

	respondOK(c, roadmap)
}

// Delete は指定IDのロードマップを削除する。
func (h *RoadmapHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
}

// CopyRoadmap は公開ロードマップをテンプレートとしてコピーする。
func (h *RoadmapHandler) CopyRoadmap(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	copied, err := h.service.CopyRoadmap(roadmapID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, copied)
}

// GetTemplates はテンプレートロードマップの一覧を取得する。
func (h *RoadmapHandler) GetTemplates(c *gin.Context) {
	templates, err := h.service.GetTemplates()
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, templates)
}

// CreateFromTemplate はテンプレートからユーザー用ロードマップを作成する。
func (h *RoadmapHandler) CreateFromTemplate(c *gin.Context) {
	userID := c.GetUint("userID")
	templateID, ok := parseID(c, "id")
	if !ok {
		return
	}

	roadmap, err := h.service.CreateFromTemplate(templateID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, roadmap)
}

// === ロードマップステップエンドポイント ===

// CreateStep はロードマップに新しいステップを作成する。
func (h *RoadmapHandler) CreateStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.CreateRoadmapStepRequest](c)
	if input == nil {
		return
	}

	step := &model.RoadmapStep{
		Title:       input.Title,
		Description: input.Description,
		ResourceURL: input.ResourceURL,
	}
	if input.OrderIndex != nil {
		step.OrderIndex = *input.OrderIndex
	}

	if err := h.service.CreateStep(roadmapID, userID, step); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, step)
}

// UpdateStep はロードマップのステップを更新する。
func (h *RoadmapHandler) UpdateStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}
	stepID, ok := parseID(c, "stepId")
	if !ok {
		return
	}

	input := bindJSON[dto.UpdateRoadmapStepRequest](c)
	if input == nil {
		return
	}

	// 完了ステータスの変更を別途処理する
	if input.IsCompleted != nil {
		step, err := h.service.UpdateStepCompletion(roadmapID, stepID, userID, *input.IsCompleted)
		if err != nil {
			respondError(c, err)
			return
		}
		// 完了ステータスのみの更新の場合は早期リターンする
		if input.Title == nil && input.Description == nil && input.ResourceURL == nil {
			respondOK(c, step)
			return
		}
	}

	updates := &model.RoadmapStep{}
	if input.Title != nil {
		updates.Title = *input.Title
	}
	if input.Description != nil {
		updates.Description = *input.Description
	}
	if input.ResourceURL != nil {
		updates.ResourceURL = *input.ResourceURL
	}

	step, err := h.service.UpdateStep(roadmapID, stepID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, step)
}

// BatchCompleteSteps はロードマップの複数ステップを一括で完了にする。
func (h *RoadmapHandler) BatchCompleteSteps(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.BatchCompleteStepsRequest](c)
	if input == nil {
		return
	}

	roadmap, err := h.service.BatchCompleteSteps(roadmapID, userID, input.StepIDs)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, roadmap)
}

// DeleteStep はロードマップのステップを削除する。
func (h *RoadmapHandler) DeleteStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}
	stepID, ok := parseID(c, "stepId")
	if !ok {
		return
	}

	if err := h.service.DeleteStep(roadmapID, stepID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// ReorderSteps はロードマップ内のステップの並び順を変更する。
func (h *RoadmapHandler) ReorderSteps(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.ReorderRoadmapStepsRequest](c)
	if input == nil {
		return
	}

	if err := h.service.ReorderSteps(roadmapID, userID, input.Orders); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("steps reordered"))
}

// GetMyCount は認証ユーザーのロードマップ総数を返す。
func (h *RoadmapHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.service.CountByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}

// GetMyStats は認証ユーザーのロードマップ統計情報を取得する。
func (h *RoadmapHandler) GetMyStats(c *gin.Context) {
	userID := c.GetUint("userID")

	stats, err := h.service.GetStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
