package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// StudyCircleRepository はスタディサークルの永続化に対する、usecase 側が要求する契約。
// サークル CRUD・メンバー管理・ステップ管理・進捗・チェックイン・ストリークランキングを提供する。
type StudyCircleRepository interface {
	// サークル CRUD
	// CreateWithOwner はサークル行の作成とオーナーのメンバー登録を不可分に行う。
	// オーナー登録に失敗した場合はサークル行も残さない。
	CreateWithOwner(ctx context.Context, circle *model.StudyCircle) error
	// FindByID は指定 ID のサークルを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.StudyCircle, error)
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.StudyCircle, int64, error)
	Update(ctx context.Context, circle *model.StudyCircle) error
	Delete(ctx context.Context, id uint) error

	// フィルタリング / 検索
	GetByStatus(ctx context.Context, userID uint, status string) ([]model.StudyCircle, error)
	Search(ctx context.Context, query string, limit, offset int) ([]model.StudyCircle, int64, error)

	// メンバー管理
	AddMember(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error
	RemoveMember(ctx context.Context, circleID, userID uint) error
	GetMembers(ctx context.Context, circleID uint) ([]model.StudyCircleMember, error)
	IsMember(ctx context.Context, circleID, userID uint) (bool, error)
	GetMemberCount(ctx context.Context, circleID uint) (int, error)
	UpdateMemberRole(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error
	CountByUserID(ctx context.Context, userID uint) (int64, error)

	// ステップ CRUD
	CreateStep(ctx context.Context, step *model.StudyCircleStep) error
	UpdateStep(ctx context.Context, step *model.StudyCircleStep) error
	DeleteStep(ctx context.Context, stepID uint) error
	// FindStepByID は指定 ID のステップを返す。不在の場合は (nil, nil) を返す。
	FindStepByID(ctx context.Context, stepID uint) (*model.StudyCircleStep, error)
	ReorderSteps(ctx context.Context, circleID uint, stepOrders []model.StepOrder) error

	// 進捗管理
	UpsertProgress(ctx context.Context, progress *model.StudyCircleMemberProgress) error
	GetProgress(ctx context.Context, circleID uint) ([]model.StudyCircleMemberProgress, error)

	// チェックイン
	CreateCheckin(ctx context.Context, checkin *model.StudyCircleCheckin) error
	GetCheckins(ctx context.Context, circleID uint) ([]model.StudyCircleCheckin, error)
	HasCheckedInToday(ctx context.Context, circleID, userID uint) (bool, error)

	// ストリークランキング
	GetStreakRanking(ctx context.Context, circleID uint) ([]model.CircleMemberStreak, error)
}
