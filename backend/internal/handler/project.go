package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// parseDate は日付文字列を "2006-01-02" 形式でパースする。
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// ProjectServiceInterface はProjectHandlerが依存するサービスメソッドを定義する。
type ProjectServiceInterface interface {
	Create(project *model.Project) error
	GetByID(id, userID uint) (*model.Project, error)
	GetByUserID(userID uint, limit, offset int) ([]model.Project, int64, error)
	GetFeaturedByUserID(userID uint) ([]model.Project, error)
	GetAll(limit, offset int) ([]model.Project, int64, error)
	Update(id, userID uint, updates *model.Project) (*model.Project, error)
	UpdateFeatured(id, userID uint, featured bool) (*model.Project, error)
	Delete(id, userID uint) error
	Archive(id, userID uint) error
	Unarchive(id, userID uint) error
	GetArchivedByUserID(userID uint, limit, offset int) ([]model.Project, int64, error)
}

// ProjectHandler はプロジェクト関連のHTTPハンドラ。
// プロジェクトのCRUD・注目プロジェクト取得・一覧取得を処理する。
type ProjectHandler struct {
	service ProjectServiceInterface
}

// NewProjectHandler は新しいProjectHandlerインスタンスを生成する。
func NewProjectHandler(s ProjectServiceInterface) *ProjectHandler {
	return &ProjectHandler{service: s}
}

// Create は新しいプロジェクトを作成する。
func (h *ProjectHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.CreateProjectRequest](c)
	if req == nil {
		return
	}

	project := &model.Project{
		UserID:       userID,
		Title:        req.Title,
		Description:  req.Description,
		TechStack:    req.TechStack,
		DemoURL:      req.DemoURL,
		GithubURL:    req.GithubURL,
		ImageURL:     req.ImageURL,
		Role:         req.Role,
		Featured:     req.Featured,
		GithubRepoID: req.GithubRepoID,
	}

	if req.StartDate != "" {
		startDate, err := parseDate(req.StartDate)
		if err == nil {
			project.StartDate = &startDate
		}
	}
	if req.EndDate != "" {
		endDate, err := parseDate(req.EndDate)
		if err == nil {
			project.EndDate = &endDate
		}
	}

	if err := h.service.Create(project); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, project)
}

// GetByID は指定IDのプロジェクトを取得する。
func (h *ProjectHandler) GetByID(c *gin.Context) {
	handleGetByID(c, h.service.GetByID)
}

// GetByUserID は指定ユーザーのプロジェクト一覧をページネーション付きで取得する。
func (h *ProjectHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	projects, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ProjectListResponse{
		Projects: projects,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	})
}

// GetFeatured は指定ユーザーの注目プロジェクト一覧を取得する。
func (h *ProjectHandler) GetFeatured(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	projects, err := h.service.GetFeaturedByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, projects)
}

// Update は指定IDのプロジェクトを更新する。
func (h *ProjectHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdateProjectRequest](c)
	if req == nil {
		return
	}

	updates := &model.Project{}
	if req.Title != "" {
		updates.Title = req.Title
	}
	if req.Description != "" {
		updates.Description = req.Description
	}
	if req.TechStack != "" {
		updates.TechStack = req.TechStack
	}
	if req.DemoURL != "" {
		updates.DemoURL = req.DemoURL
	}
	if req.GithubURL != "" {
		updates.GithubURL = req.GithubURL
	}
	if req.ImageURL != "" {
		updates.ImageURL = req.ImageURL
	}
	if req.Role != "" {
		updates.Role = req.Role
	}
	if req.GithubRepoID != nil {
		updates.GithubRepoID = req.GithubRepoID
	}
	if req.StartDate != "" {
		startDate, err := parseDate(req.StartDate)
		if err == nil {
			updates.StartDate = &startDate
		}
	}
	if req.EndDate != "" {
		endDate, err := parseDate(req.EndDate)
		if err == nil {
			updates.EndDate = &endDate
		}
	}

	project, err := h.service.Update(id, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	// featuredはboolポインタのため別途処理する
	if req.Featured != nil {
		project, err = h.service.UpdateFeatured(id, userID, *req.Featured)
		if err != nil {
			respondError(c, err)
			return
		}
	}

	respondOK(c, project)
}

// Delete は指定IDのプロジェクトを削除する。
func (h *ProjectHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
}

// Archive はプロジェクトをアーカイブする。
func (h *ProjectHandler) Archive(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Archive(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "archived"})
}

// Unarchive はプロジェクトのアーカイブを解除する。
func (h *ProjectHandler) Unarchive(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Unarchive(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "unarchived"})
}

// GetArchived はアーカイブ済みプロジェクト一覧を取得する。
func (h *ProjectHandler) GetArchived(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	projects, total, err := h.service.GetArchivedByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ProjectListResponse{
		Projects: projects,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	})
}

// GetAll はプロジェクトの一覧をページネーション付きで取得する。
func (h *ProjectHandler) GetAll(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	projects, total, err := h.service.GetAll(limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ProjectListResponse{
		Projects: projects,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	})
}
