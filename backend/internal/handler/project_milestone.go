package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// ProjectMilestoneServiceInterface はマイルストーンサービスの抽象インターフェース。
type ProjectMilestoneServiceInterface interface {
	Create(userID, projectID uint, title, description string, dueDate *time.Time) error
	GetByProjectID(projectID uint) ([]model.ProjectMilestone, error)
	Update(userID, milestoneID uint, title, description string, dueDate *time.Time, status string) (*model.ProjectMilestone, error)
	Delete(userID, milestoneID uint) error
}

// ProjectMilestoneHandler はマイルストーン関連のHTTPハンドラ。
type ProjectMilestoneHandler struct {
	service ProjectMilestoneServiceInterface
}

// NewProjectMilestoneHandler は新しいProjectMilestoneHandlerインスタンスを生成する。
func NewProjectMilestoneHandler(s ProjectMilestoneServiceInterface) *ProjectMilestoneHandler {
	return &ProjectMilestoneHandler{service: s}
}

// Create はマイルストーンを作成する。
func (h *ProjectMilestoneHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	projectID, ok := parseID(c, "projectId")
	if !ok {
		return
	}

	req := bindJSON[dto.CreateMilestoneRequest](c)
	if req == nil {
		return
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		d, err := parseDate(req.DueDate)
		if err != nil {
			respondBadRequest(c, "日付の形式が不正です（YYYY-MM-DD）")
			return
		}
		dueDate = &d
	}

	if err := h.service.Create(userID, projectID, req.Title, req.Description, dueDate); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, domain.NewMessageResponse("マイルストーンを作成しました"))
}

// GetByProjectID はプロジェクトのマイルストーン一覧を取得する。
func (h *ProjectMilestoneHandler) GetByProjectID(c *gin.Context) {
	projectID, ok := parseID(c, "projectId")
	if !ok {
		return
	}

	milestones, err := h.service.GetByProjectID(projectID)
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
		d, err := parseDate(req.DueDate)
		if err != nil {
			respondBadRequest(c, "日付の形式が不正です（YYYY-MM-DD）")
			return
		}
		dueDate = &d
	}

	result, err := h.service.Update(userID, milestoneID, req.Title, req.Description, dueDate, req.Status)
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

	if err := h.service.Delete(userID, milestoneID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("マイルストーンを削除しました"))
}
