package dto

import "github.com/norman6464/devsync/backend/internal/repository"

// CreateStudyCircleRequest はスタディサークル作成リクエストのDTO。
type CreateStudyCircleRequest struct {
	Name        string `json:"name" binding:"required"`
	Topic       string `json:"topic" binding:"required"`
	Description string `json:"description"`
	MaxMembers  int    `json:"max_members"`
	MemberIDs   []uint `json:"member_ids"`
}

// UpdateStudyCircleRequest はスタディサークル更新リクエストのDTO。
type UpdateStudyCircleRequest struct {
	Name        *string `json:"name"`
	Topic       *string `json:"topic"`
	Description *string `json:"description"`
}

// AddStudyCircleMemberRequest はスタディサークルメンバー追加リクエストのDTO。
type AddStudyCircleMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// CreateStudyCircleStepRequest はスタディサークルステップ作成リクエストのDTO。
type CreateStudyCircleStepRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ResourceURL string `json:"resource_url"`
	OrderIndex  int    `json:"order_index"`
}

// UpdateStudyCircleStepRequest はスタディサークルステップ更新リクエストのDTO。
type UpdateStudyCircleStepRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
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
	Content string `json:"content" binding:"required"`
}
