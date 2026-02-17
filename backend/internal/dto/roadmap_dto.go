package dto

import "github.com/norman6464/devsync/backend/internal/service"

// CreateRoadmapRequest はロードマップ作成リクエストのDTO。
type CreateRoadmapRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsPublic    bool   `json:"is_public"`
}

// UpdateRoadmapRequest はロードマップ更新リクエストのDTO。
type UpdateRoadmapRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	IsPublic    *bool   `json:"is_public"`
	Status      *string `json:"status"`
}

// CreateRoadmapStepRequest はロードマップステップ作成リクエストのDTO。
type CreateRoadmapStepRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ResourceURL string `json:"resource_url"`
	OrderIndex  *int   `json:"order_index"`
}

// UpdateRoadmapStepRequest はロードマップステップ更新リクエストのDTO。
type UpdateRoadmapStepRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	ResourceURL *string `json:"resource_url"`
	IsCompleted *bool   `json:"is_completed"`
}

// ReorderRoadmapStepsRequest はロードマップステップ並べ替えリクエストのDTO。
type ReorderRoadmapStepsRequest struct {
	Orders []service.StepOrder `json:"orders" binding:"required"`
}
