package dto

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// CreateRoadmapRequest はロードマップ作成リクエストのDTO。
type CreateRoadmapRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Category    string `json:"category" binding:"omitempty,max=100"`
	IsPublic    bool   `json:"is_public"`
}

// UpdateRoadmapRequest はロードマップ更新リクエストのDTO。
type UpdateRoadmapRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	Category    *string `json:"category" binding:"omitempty,max=100"`
	IsPublic    *bool   `json:"is_public"`
	Status      *string `json:"status" binding:"omitempty,max=50"`
}

// CreateRoadmapStepRequest はロードマップステップ作成リクエストのDTO。
type CreateRoadmapStepRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	ResourceURL string `json:"resource_url" binding:"omitempty,max=2000"`
	OrderIndex  *int   `json:"order_index"`
}

// UpdateRoadmapStepRequest はロードマップステップ更新リクエストのDTO。
type UpdateRoadmapStepRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	ResourceURL *string `json:"resource_url" binding:"omitempty,max=2000"`
	IsCompleted *bool   `json:"is_completed"`
}

// ReorderRoadmapStepsRequest はロードマップステップ並べ替えリクエストのDTO。
type ReorderRoadmapStepsRequest struct {
	Orders []service.StepOrder `json:"orders" binding:"required"`
}

// RoadmapListResponse はロードマップ一覧レスポンス。
type RoadmapListResponse struct {
	Roadmaps []model.Roadmap `json:"roadmaps"`
	Total    int64           `json:"total"`
}
