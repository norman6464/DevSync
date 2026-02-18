package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CreateProjectRequest はプロジェクト作成のリクエストボディ。
type CreateProjectRequest struct {
	Title        string `json:"title" binding:"required,max=200" validate:"required,max=200"`
	Description  string `json:"description" binding:"omitempty,max=5000"`
	TechStack    string `json:"tech_stack" binding:"omitempty,max=500"`
	DemoURL      string `json:"demo_url" binding:"omitempty,max=2000"`
	GithubURL    string `json:"github_url" binding:"omitempty,max=2000"`
	ImageURL     string `json:"image_url" binding:"omitempty,max=2000"`
	Role         string `json:"role" binding:"omitempty,max=200"`
	StartDate    string `json:"start_date" binding:"omitempty,max=20"`
	EndDate      string `json:"end_date" binding:"omitempty,max=20"`
	Featured     bool   `json:"featured"`
	GithubRepoID *uint  `json:"github_repo_id"`
}

// UpdateProjectRequest はプロジェクト更新のリクエストボディ。
type UpdateProjectRequest struct {
	Title        string `json:"title" binding:"omitempty,max=200" validate:"omitempty,max=200"`
	Description  string `json:"description" binding:"omitempty,max=5000"`
	TechStack    string `json:"tech_stack" binding:"omitempty,max=500"`
	DemoURL      string `json:"demo_url" binding:"omitempty,max=2000"`
	GithubURL    string `json:"github_url" binding:"omitempty,max=2000"`
	ImageURL     string `json:"image_url" binding:"omitempty,max=2000"`
	Role         string `json:"role" binding:"omitempty,max=200"`
	StartDate    string `json:"start_date" binding:"omitempty,max=20"`
	EndDate      string `json:"end_date" binding:"omitempty,max=20"`
	Featured     *bool  `json:"featured"`
	GithubRepoID *uint  `json:"github_repo_id"`
}

// ProjectListResponse はプロジェクト一覧レスポンス。
type ProjectListResponse struct {
	Projects []model.Project `json:"projects"`
	Total    int64           `json:"total"`
	Limit    int             `json:"limit"`
	Offset   int             `json:"offset"`
}
