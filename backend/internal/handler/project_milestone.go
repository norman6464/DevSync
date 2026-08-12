package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
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

// Create はマイルストーンを作成する。
func (h *ProjectMilestoneHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	projectID, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.CreateMilestoneRequest](c)
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

	respondOK(c, dto.MilestoneListResponse{
		Milestones: ensureSlice(milestones),
	})
}

// Update はマイルストーンを更新する。
func (h *ProjectMilestoneHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	milestoneID, ok := parseID(c, "milestoneId")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdateMilestoneRequest](c)
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
