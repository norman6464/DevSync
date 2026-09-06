package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// ProjectMilestoneHandler はマイルストーン関連のHTTPハンドラ。
type ProjectMilestoneHandler struct {
	create *usecase.CreateProjectMilestoneUseCase
	list   *usecase.ListProjectMilestonesUseCase
	update *usecase.UpdateProjectMilestoneUseCase
	remove *usecase.DeleteProjectMilestoneUseCase
}

// NewProjectMilestoneHandler は新しいProjectMilestoneHandlerインスタンスを生成する。
func NewProjectMilestoneHandler(
	create *usecase.CreateProjectMilestoneUseCase,
	list *usecase.ListProjectMilestonesUseCase,
	update *usecase.UpdateProjectMilestoneUseCase,
	remove *usecase.DeleteProjectMilestoneUseCase,
) *ProjectMilestoneHandler {
	return &ProjectMilestoneHandler{create: create, list: list, update: update, remove: remove}
}

// createMilestoneRequest はマイルストーン作成のリクエストボディ。
type createMilestoneRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=5000"`
	DueDate     string `json:"due_date" binding:"omitempty,max=20"`
}

// Create はマイルストーンを作成する。
func (h *ProjectMilestoneHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	projectID, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[createMilestoneRequest](c)
	if req == nil {
		return
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		d, ok := parseDateParam(req.DueDate)
		if !ok {
			respondBadRequest(c, "日付の形式が不正です（YYYY-MM-DD）")
			return
		}
		dueDate = &d
	}

	if err := h.create.Execute(c.Request.Context(), usecase.CreateProjectMilestoneInput{
		UserID:      userID,
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
	}); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, domain.NewMessageResponse("マイルストーンを作成しました"))
}

// milestoneListResponse はマイルストーン一覧レスポンス。
type milestoneListResponse struct {
	Milestones []model.ProjectMilestone `json:"milestones"`
}

// GetByProjectID はプロジェクトのマイルストーン一覧を取得する。
func (h *ProjectMilestoneHandler) GetByProjectID(c *gin.Context) {
	projectID, ok := parseID(c, "id")
	if !ok {
		return
	}

	milestones, err := h.list.Execute(c.Request.Context(), projectID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, milestoneListResponse{
		Milestones: ensureSlice(milestones),
	})
}

// updateMilestoneRequest はマイルストーン更新のリクエストボディ。
type updateMilestoneRequest struct {
	Title       string `json:"title" binding:"omitempty,max=200"`
	Description string `json:"description" binding:"omitempty,max=5000"`
	DueDate     string `json:"due_date" binding:"omitempty,max=20"`
	Status      string `json:"status" binding:"omitempty"`
}

// Update はマイルストーンを更新する。
func (h *ProjectMilestoneHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	milestoneID, ok := parseID(c, "milestoneId")
	if !ok {
		return
	}

	req := bindJSON[updateMilestoneRequest](c)
	if req == nil {
		return
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		d, ok := parseDateParam(req.DueDate)
		if !ok {
			respondBadRequest(c, "日付の形式が不正です（YYYY-MM-DD）")
			return
		}
		dueDate = &d
	}

	result, err := h.update.Execute(c.Request.Context(), usecase.UpdateProjectMilestoneInput{
		UserID:      userID,
		MilestoneID: milestoneID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
		Status:      req.Status,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, result)
}

// Delete はマイルストーンを削除する。
func (h *ProjectMilestoneHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	milestoneID, ok := parseID(c, "milestoneId")
	if !ok {
		return
	}

	if err := h.remove.Execute(c.Request.Context(), userID, milestoneID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("マイルストーンを削除しました"))
}
