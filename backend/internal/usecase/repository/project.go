package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ProjectRepository はプロジェクトショーケースの永続化に対する、usecase 側が要求する契約。
//
// project_milestone 用の [ProjectReader] は別スライスが所有する最小 port なので、
// こちらへは統合していない。
type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id uint) error
	// FindByID は指定 ID のプロジェクトを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.Project, error)

	// FindByUserID はユーザーのプロジェクトを注目優先・作成日の新しい順で返し、総数も返す。
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error)
	FindFeaturedByUserID(ctx context.Context, userID uint) ([]model.Project, error)
	FindArchivedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error)
	FindAll(ctx context.Context, limit, offset int) ([]model.Project, int64, error)
	// Search はタイトル・説明・技術スタックへの部分一致で検索する（大文字小文字を区別しない）。
	Search(ctx context.Context, query string, limit, offset int) ([]model.Project, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)

	Archive(ctx context.Context, id uint) error
	Unarchive(ctx context.Context, id uint) error
}
