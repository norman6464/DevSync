package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CreateProjectRequest はプロジェクト作成のリクエストボディ。
type CreateProjectRequest struct {
	Title        string `json:"title" binding:"required,max=200" validate:"required,max=200"`
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
	Title        string `json:"title" binding:"max=200" validate:"max=200"`
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

// ProjectListResponse はプロジェクト一覧レスポンス。
type ProjectListResponse struct {
	Projects []model.Project `json:"projects"`
	Total    int64           `json:"total"`
	Limit    int             `json:"limit"`
	Offset   int             `json:"offset"`
}
