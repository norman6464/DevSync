package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// ProjectHandler はプロジェクト関連のHTTPハンドラ。
// プロジェクトのCRUD・注目プロジェクト取得・一覧取得を処理する。
type ProjectHandler struct {
	create         *usecase.CreateProjectUseCase
	get            *usecase.GetProjectUseCase
	listByUser     *usecase.ListProjectsByUserUseCase
	listFeatured   *usecase.ListFeaturedProjectsUseCase
	listAll        *usecase.ListAllProjectsUseCase
	listArchived   *usecase.ListArchivedProjectsUseCase
	search         *usecase.SearchProjectsUseCase
	update         *usecase.UpdateProjectUseCase
	updateFeatured *usecase.UpdateProjectFeaturedUseCase
	archive        *usecase.ArchiveProjectUseCase
	unarchive      *usecase.UnarchiveProjectUseCase
	remove         *usecase.DeleteProjectUseCase
	count          *usecase.CountProjectsUseCase
}

// NewProjectHandler は新しいProjectHandlerインスタンスを生成する。
func NewProjectHandler(
	create *usecase.CreateProjectUseCase,
	get *usecase.GetProjectUseCase,
	listByUser *usecase.ListProjectsByUserUseCase,
	listFeatured *usecase.ListFeaturedProjectsUseCase,
	listAll *usecase.ListAllProjectsUseCase,
	listArchived *usecase.ListArchivedProjectsUseCase,
	search *usecase.SearchProjectsUseCase,
	update *usecase.UpdateProjectUseCase,
	updateFeatured *usecase.UpdateProjectFeaturedUseCase,
	archive *usecase.ArchiveProjectUseCase,
	unarchive *usecase.UnarchiveProjectUseCase,
	remove *usecase.DeleteProjectUseCase,
	count *usecase.CountProjectsUseCase,
) *ProjectHandler {
	return &ProjectHandler{
		create: create, get: get,
		listByUser: listByUser, listFeatured: listFeatured,
		listAll: listAll, listArchived: listArchived, search: search,
		update: update, updateFeatured: updateFeatured,
		archive: archive, unarchive: unarchive, remove: remove, count: count,
	}
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

	if startDate, ok := parseDateParam(req.StartDate); ok && !startDate.IsZero() {
		project.StartDate = &startDate
	}
	if endDate, ok := parseDateParam(req.EndDate); ok && !endDate.IsZero() {
		project.EndDate = &endDate
	}

	if err := h.create.Execute(c.Request.Context(), project); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, project)
}

// GetByID は指定IDのプロジェクトを取得する。
func (h *ProjectHandler) GetByID(c *gin.Context) {
	handleGetByID(c, func(id, userID uint) (*model.Project, error) {
		return h.get.Execute(c.Request.Context(), id, userID)
	})
}

// GetByUserID は指定ユーザーのプロジェクト一覧をページネーション付きで取得する。
func (h *ProjectHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	projects, total, err := h.listByUser.Execute(c.Request.Context(), userID, limit, offset)
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

// GetMyProjects は認証ユーザー自身のプロジェクト一覧を取得する。
func (h *ProjectHandler) GetMyProjects(c *gin.Context) {
	userID := c.GetUint("userID")

	limit, offset := parseLimitOffset(c)

	projects, total, err := h.listByUser.Execute(c.Request.Context(), userID, limit, offset)
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

	projects, err := h.listFeatured.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(projects))
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
	if startDate, ok := parseDateParam(req.StartDate); ok && !startDate.IsZero() {
		updates.StartDate = &startDate
	}
	if endDate, ok := parseDateParam(req.EndDate); ok && !endDate.IsZero() {
		updates.EndDate = &endDate
	}

	project, err := h.update.Execute(c.Request.Context(), id, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	// featuredはboolポインタのため別途処理する
	if req.Featured != nil {
		project, err = h.updateFeatured.Execute(c.Request.Context(), id, userID, *req.Featured)
		if err != nil {
			respondError(c, err)
			return
		}
	}

	respondOK(c, project)
}

// Delete は指定IDのプロジェクトを削除する。
func (h *ProjectHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.remove.Execute(c.Request.Context(), id, userID)
	})
}

// Archive はプロジェクトをアーカイブする。
func (h *ProjectHandler) Archive(c *gin.Context) {
	handleAction(c, func(id, userID uint) error {
		return h.archive.Execute(c.Request.Context(), id, userID)
	}, "archived")
}

// Unarchive はプロジェクトのアーカイブを解除する。
func (h *ProjectHandler) Unarchive(c *gin.Context) {
	handleAction(c, func(id, userID uint) error {
		return h.unarchive.Execute(c.Request.Context(), id, userID)
	}, "unarchived")
}

// GetArchived はアーカイブ済みプロジェクト一覧を取得する。
func (h *ProjectHandler) GetArchived(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	projects, total, err := h.listArchived.Execute(c.Request.Context(), userID, limit, offset)
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

// Search はプロジェクトをキーワード検索する。
func (h *ProjectHandler) Search(c *gin.Context) {
	q, ok := parseSearchQuery(c)
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	projects, total, err := h.search.Execute(c.Request.Context(), q, limit, offset)
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

// GetMyCount は認証ユーザーのプロジェクト総数を取得する。
func (h *ProjectHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// GetAll はプロジェクトの一覧をページネーション付きで取得する。
func (h *ProjectHandler) GetAll(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	projects, total, err := h.listAll.Execute(c.Request.Context(), limit, offset)
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
