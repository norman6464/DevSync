package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// RoadmapHandler はロードマップ関連のHTTPハンドラ。
// ロードマップとステップのCRUD・公開一覧・コピー・並べ替えを処理する。
type RoadmapHandler struct {
	create        *usecase.CreateRoadmapUseCase
	get           *usecase.GetRoadmapUseCase
	listByUser    *usecase.ListRoadmapsByUserUseCase
	listByStatus  *usecase.ListRoadmapsByStatusUseCase
	listPublic    *usecase.ListPublicRoadmapsUseCase
	update        *usecase.UpdateRoadmapUseCase
	updateVisible *usecase.UpdateRoadmapVisibilityUseCase
	remove        *usecase.DeleteRoadmapUseCase
	copy          *usecase.CopyRoadmapUseCase
	listTemplates *usecase.ListRoadmapTemplatesUseCase
	fromTemplate  *usecase.CreateRoadmapFromTemplateUseCase
	createStep    *usecase.CreateRoadmapStepUseCase
	updateStep    *usecase.UpdateRoadmapStepUseCase
	completeStep  *usecase.UpdateRoadmapStepCompletionUseCase
	batchComplete *usecase.BatchCompleteRoadmapStepsUseCase
	deleteStep    *usecase.DeleteRoadmapStepUseCase
	reorderSteps  *usecase.ReorderRoadmapStepsUseCase
	stats         *usecase.GetRoadmapStatsUseCase
	count         *usecase.CountRoadmapsUseCase
}

// NewRoadmapHandler は新しいRoadmapHandlerインスタンスを生成する。
func NewRoadmapHandler(
	create *usecase.CreateRoadmapUseCase,
	get *usecase.GetRoadmapUseCase,
	listByUser *usecase.ListRoadmapsByUserUseCase,
	listByStatus *usecase.ListRoadmapsByStatusUseCase,
	listPublic *usecase.ListPublicRoadmapsUseCase,
	update *usecase.UpdateRoadmapUseCase,
	updateVisible *usecase.UpdateRoadmapVisibilityUseCase,
	remove *usecase.DeleteRoadmapUseCase,
	copyRoadmap *usecase.CopyRoadmapUseCase,
	listTemplates *usecase.ListRoadmapTemplatesUseCase,
	fromTemplate *usecase.CreateRoadmapFromTemplateUseCase,
	createStep *usecase.CreateRoadmapStepUseCase,
	updateStep *usecase.UpdateRoadmapStepUseCase,
	completeStep *usecase.UpdateRoadmapStepCompletionUseCase,
	batchComplete *usecase.BatchCompleteRoadmapStepsUseCase,
	deleteStep *usecase.DeleteRoadmapStepUseCase,
	reorderSteps *usecase.ReorderRoadmapStepsUseCase,
	stats *usecase.GetRoadmapStatsUseCase,
	count *usecase.CountRoadmapsUseCase,
) *RoadmapHandler {
	return &RoadmapHandler{
		create: create, get: get, listByUser: listByUser, listByStatus: listByStatus,
		listPublic: listPublic, update: update, updateVisible: updateVisible, remove: remove,
		copy: copyRoadmap, listTemplates: listTemplates, fromTemplate: fromTemplate,
		createStep: createStep, updateStep: updateStep, completeStep: completeStep,
		batchComplete: batchComplete, deleteStep: deleteStep, reorderSteps: reorderSteps,
		stats: stats, count: count,
	}
}

// === ロードマップエンドポイント ===

// createRoadmapRequest はロードマップ作成リクエストのDTO。
type createRoadmapRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Category    string `json:"category" binding:"omitempty,max=100"`
	IsPublic    bool   `json:"is_public"`
}

// Create は新しいロードマップを作成する。
func (h *RoadmapHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[createRoadmapRequest](c)
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

	if err := h.create.Execute(c.Request.Context(), roadmap); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, roadmap)
}

// roadmapListResponse はロードマップ一覧レスポンス。
type roadmapListResponse struct {
	Roadmaps []model.Roadmap `json:"roadmaps"`
	Total    int64           `json:"total"`
	Limit    int             `json:"limit"`
	Offset   int             `json:"offset"`
}

// GetMyRoadmaps は現在のユーザーのロードマップ一覧を取得する。
func (h *RoadmapHandler) GetMyRoadmaps(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	roadmaps, total, err := h.listByUser.Execute(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, roadmapListResponse{
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

	roadmaps, err := h.listByStatus.Execute(c.Request.Context(), userID, status)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(roadmaps))
}

// GetPublicRoadmaps は公開ロードマップの一覧をページネーション付きで取得する。
func (h *RoadmapHandler) GetPublicRoadmaps(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	roadmaps, total, err := h.listPublic.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, roadmapListResponse{
		Roadmaps: roadmaps,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	})
}

// GetByID は指定IDのロードマップをステップ付きで取得する。
func (h *RoadmapHandler) GetByID(c *gin.Context) {
	handleGetByID(c, func(id, userID uint) (*model.Roadmap, error) {
		return h.get.Execute(c.Request.Context(), id, userID)
	})
}

// updateRoadmapRequest はロードマップ更新リクエストのDTO。
type updateRoadmapRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	Category    *string `json:"category" binding:"omitempty,max=100"`
	IsPublic    *bool   `json:"is_public"`
	Status      *string `json:"status" binding:"omitempty,max=50"`
}

// Update は指定IDのロードマップを更新する。
func (h *RoadmapHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[updateRoadmapRequest](c)
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

	roadmap, err := h.update.Execute(c.Request.Context(), roadmapID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	// IsPublicが指定されている場合は別途処理する
	if input.IsPublic != nil {
		roadmap, err = h.updateVisible.Execute(c.Request.Context(), roadmapID, userID, *input.IsPublic)
		if err != nil {
			respondError(c, err)
			return
		}
	}

	respondOK(c, roadmap)
}

// Delete は指定IDのロードマップを削除する。
func (h *RoadmapHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.remove.Execute(c.Request.Context(), id, userID)
	})
}

// CopyRoadmap は公開ロードマップをテンプレートとしてコピーする。
func (h *RoadmapHandler) CopyRoadmap(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	copied, err := h.copy.Execute(c.Request.Context(), roadmapID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, copied)
}

// GetTemplates はテンプレートロードマップの一覧を取得する。
func (h *RoadmapHandler) GetTemplates(c *gin.Context) {
	templates, err := h.listTemplates.Execute(c.Request.Context())
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

	roadmap, err := h.fromTemplate.Execute(c.Request.Context(), templateID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, roadmap)
}

// === ロードマップステップエンドポイント ===

// createRoadmapStepRequest はロードマップステップ作成リクエストのDTO。
type createRoadmapStepRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	ResourceURL string `json:"resource_url" binding:"omitempty,max=2000"`
	OrderIndex  *int   `json:"order_index" binding:"omitempty,min=0"`
}

// CreateStep はロードマップに新しいステップを作成する。
func (h *RoadmapHandler) CreateStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[createRoadmapStepRequest](c)
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

	if err := h.createStep.Execute(c.Request.Context(), roadmapID, userID, step); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, step)
}

// updateRoadmapStepRequest はロードマップステップ更新リクエストのDTO。
type updateRoadmapStepRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	ResourceURL *string `json:"resource_url" binding:"omitempty,max=2000"`
	IsCompleted *bool   `json:"is_completed"`
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

	input := bindJSON[updateRoadmapStepRequest](c)
	if input == nil {
		return
	}

	// 完了ステータスの変更を別途処理する
	if input.IsCompleted != nil {
		step, err := h.completeStep.Execute(c.Request.Context(), roadmapID, stepID, userID, *input.IsCompleted)
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

	step, err := h.updateStep.Execute(c.Request.Context(), roadmapID, stepID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, step)
}

// batchCompleteStepsRequest はステップ一括完了リクエストのDTO。
type batchCompleteStepsRequest struct {
	StepIDs []uint `json:"step_ids" binding:"required,min=1"`
}

// BatchCompleteSteps はロードマップの複数ステップを一括で完了にする。
func (h *RoadmapHandler) BatchCompleteSteps(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[batchCompleteStepsRequest](c)
	if input == nil {
		return
	}

	roadmap, err := h.batchComplete.Execute(c.Request.Context(), roadmapID, userID, input.StepIDs)
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

	if err := h.deleteStep.Execute(c.Request.Context(), roadmapID, stepID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// reorderRoadmapStepsRequest はロードマップステップ並べ替えリクエストのDTO。
type reorderRoadmapStepsRequest struct {
	Orders []model.StepOrder `json:"orders" binding:"required"`
}

// ReorderSteps はロードマップ内のステップの並び順を変更する。
func (h *RoadmapHandler) ReorderSteps(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[reorderRoadmapStepsRequest](c)
	if input == nil {
		return
	}

	if err := h.reorderSteps.Execute(c.Request.Context(), roadmapID, userID, input.Orders); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("steps reordered"))
}

// GetMyCount は認証ユーザーのロードマップ総数を返す。
func (h *RoadmapHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}

// GetMyStats は認証ユーザーのロードマップ統計情報を取得する。
func (h *RoadmapHandler) GetMyStats(c *gin.Context) {
	userID := c.GetUint("userID")

	stats, err := h.stats.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
