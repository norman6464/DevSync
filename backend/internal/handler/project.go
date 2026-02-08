package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// parseDate は日付文字列を "2006-01-02" 形式でパースする。
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// ProjectHandler はプロジェクト関連のHTTPハンドラ。
// プロジェクトのCRUD・注目プロジェクト取得・一覧取得を処理する。
type ProjectHandler struct {
	service *service.ProjectService
}

// NewProjectHandler は新しいProjectHandlerインスタンスを生成する。
func NewProjectHandler(s *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

// CreateProjectRequest はプロジェクト作成のリクエストボディ。
type CreateProjectRequest struct {
	Title        string `json:"title" binding:"required,max=200"`
	Description  string `json:"description"`
	TechStack    string `json:"tech_stack"`
	DemoURL      string `json:"demo_url"`
	GithubURL    string `json:"github_url"`
	ImageURL     string `json:"image_url"`
	Role         string `json:"role"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Featured     bool   `json:"featured"`
	GithubRepoID *uint  `json:"github_repo_id"`
}

// UpdateProjectRequest はプロジェクト更新のリクエストボディ。
type UpdateProjectRequest struct {
	Title        string `json:"title" binding:"max=200"`
	Description  string `json:"description"`
	TechStack    string `json:"tech_stack"`
	DemoURL      string `json:"demo_url"`
	GithubURL    string `json:"github_url"`
	ImageURL     string `json:"image_url"`
	Role         string `json:"role"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Featured     *bool  `json:"featured"`
	GithubRepoID *uint  `json:"github_repo_id"`
}

// Create は新しいプロジェクトを作成する。
func (h *ProjectHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetByID は指定IDのプロジェクトを取得する。
func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// GetByUserID は指定ユーザーのプロジェクト一覧を取得する。
func (h *ProjectHandler) GetByUserID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	projects, err := h.service.GetByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetFeatured は指定ユーザーの注目プロジェクト一覧を取得する。
func (h *ProjectHandler) GetFeatured(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	projects, err := h.service.GetFeaturedByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch featured projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// Update は指定IDのプロジェクトを更新する。
func (h *ProjectHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	project, err := h.service.Update(uint(id), userID, updates)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this project"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// featuredはboolポインタのため別途処理する
	if req.Featured != nil {
		project, err = h.service.UpdateFeatured(uint(id), userID, *req.Featured)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
			return
		}
	}

	c.JSON(http.StatusOK, project)
}

// Delete は指定IDのプロジェクトを削除する。
func (h *ProjectHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this project"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// GetAll はプロジェクトの一覧をページネーション付きで取得する。
func (h *ProjectHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	projects, total, err := h.service.GetAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}
