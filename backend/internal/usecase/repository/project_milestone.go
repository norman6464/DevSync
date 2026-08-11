package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ProjectMilestoneRepository はプロジェクトマイルストーンの永続化に対する、usecase 側が要求する契約。
type ProjectMilestoneRepository interface {
	Create(ctx context.Context, milestone *model.ProjectMilestone) error
	// FindByID は指定 ID のマイルストーンを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.ProjectMilestone, error)
	FindByProjectID(ctx context.Context, projectID uint) ([]model.ProjectMilestone, error)
	Update(ctx context.Context, milestone *model.ProjectMilestone) error
	Delete(ctx context.Context, id uint) error
}

// ProjectReader は所有権チェックに必要なプロジェクト読み取りだけを切り出した最小 port（-er）。
type ProjectReader interface {
	// FindByID は指定 ID のプロジェクトを返す。不在の場合は (nil, nil) を返す。
	FindByID(ctx context.Context, id uint) (*model.Project, error)
}
