package dto

import "github.com/norman6464/devsync/backend/internal/repository"

// CreateStudyCircleRequest はスタディサークル作成リクエストのDTO。
type CreateStudyCircleRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Topic       string `json:"topic" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	MaxMembers  int    `json:"max_members" binding:"omitempty,min=0,max=1000"`
	MemberIDs   []uint `json:"member_ids" binding:"omitempty,max=100"`
}

// UpdateStudyCircleRequest はスタディサークル更新リクエストのDTO。
type UpdateStudyCircleRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=200"`
	Topic       *string `json:"topic" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
}

// AddStudyCircleMemberRequest はスタディサークルメンバー追加リクエストのDTO。
type AddStudyCircleMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// CreateStudyCircleStepRequest はスタディサークルステップ作成リクエストのDTO。
type CreateStudyCircleStepRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	ResourceURL string `json:"resource_url" binding:"omitempty,max=2000"`
	OrderIndex  int    `json:"order_index" binding:"omitempty,min=0"`
}

// UpdateStudyCircleStepRequest はスタディサークルステップ更新リクエストのDTO。
type UpdateStudyCircleStepRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
}

// ReorderStudyCircleStepsRequest はスタディサークルステップ並べ替えリクエストのDTO。
type ReorderStudyCircleStepsRequest struct {
	Orders []repository.StepOrder `json:"orders" binding:"required"`
}

// UpdateStudyCircleProgressRequest はスタディサークル進捗更新リクエストのDTO。
type UpdateStudyCircleProgressRequest struct {
	IsCompleted bool `json:"is_completed"`
}

// CreateStudyCircleCheckinRequest はスタディサークルチェックイン作成リクエストのDTO。
type CreateStudyCircleCheckinRequest struct {
	Content string `json:"content" binding:"required,max=5000"`
}
