package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// projectRepository は [repository.ProjectRepository] の GORM 実装。
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository は ProjectRepository の GORM 実装を返す。
func NewProjectRepository(db *gorm.DB) repository.ProjectRepository {
	return &projectRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ProjectRepository = (*projectRepository)(nil)

// Create は新しいプロジェクトを作成する。
func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

// Update は既存のプロジェクトを更新する。
func (r *projectRepository) Update(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

// Delete はプロジェクトを削除する（モデルが論理削除を持つため soft delete になる）。
func (r *projectRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Project{}, id).Error
}

// FindByID は指定 ID のプロジェクトをユーザー・GitHub リポジトリ付きで取得する。
// 不在の場合は (nil, nil) を返す。
func (r *projectRepository) FindByID(ctx context.Context, id uint) (*model.Project, error) {
	var project model.Project
	err := r.db.WithContext(ctx).Preload("User").Preload("GithubRepo").First(&project, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// FindByUserID はユーザーのプロジェクトを注目優先・作成日の新しい順で取得し、総数も返す。
func (r *projectRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	return r.pagedProjects(ctx,
		r.db.WithContext(ctx).Model(&model.Project{}).Where("user_id = ?", userID),
		"featured DESC, created_at DESC", false, limit, offset)
}

// FindArchivedByUserID はアーカイブ済みのプロジェクトを更新日の新しい順で取得し、総数も返す。
func (r *projectRepository) FindArchivedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	return r.pagedProjects(ctx,
		r.db.WithContext(ctx).Model(&model.Project{}).Where("user_id = ? AND is_archived = ?", userID, true),
		"updated_at DESC", false, limit, offset)
}

// FindAll は全プロジェクトを作成日の新しい順で取得し、総数も返す。
func (r *projectRepository) FindAll(ctx context.Context, limit, offset int) ([]model.Project, int64, error) {
	return r.pagedProjects(ctx,
		r.db.WithContext(ctx).Model(&model.Project{}),
		"created_at DESC", true, limit, offset)
}

// Search はタイトル・説明・技術スタックへの部分一致で検索する（大文字小文字を区別しない）。
func (r *projectRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Project, int64, error) {
	like := escapeLikePattern(query)
	scope := r.db.WithContext(ctx).Model(&model.Project{}).
		Where("title ILIKE ? OR description ILIKE ? OR tech_stack ILIKE ?", like, like, like)
	return r.pagedProjects(ctx, scope, "created_at DESC", true, limit, offset)
}

// pagedProjects は絞り込み済みのクエリから総数とページを取得する共通処理。
// withUser が true のときだけユーザーを Preload する（移行前のエンドポイントごとの差を維持している）。
func (r *projectRepository) pagedProjects(
	ctx context.Context, scope *gorm.DB, order string, withUser bool, limit, offset int,
) ([]model.Project, int64, error) {
	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := scope.Session(&gorm.Session{}).Preload("GithubRepo")
	if withUser {
		query = query.Preload("User")
	}

	var projects []model.Project
	err := query.Order(order).Limit(limit).Offset(offset).Find(&projects).Error
	return projects, total, err
}

// FindFeaturedByUserID は注目のプロジェクトを作成日の新しい順で取得する。
func (r *projectRepository) FindFeaturedByUserID(ctx context.Context, userID uint) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.WithContext(ctx).Preload("GithubRepo").
		Where("user_id = ? AND featured = ?", userID, true).
		Order("created_at DESC").
		Find(&projects).Error
	return projects, err
}

// CountByUserID はユーザーのプロジェクト総数を返す。
func (r *projectRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Project{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// Archive はプロジェクトをアーカイブする。
func (r *projectRepository) Archive(ctx context.Context, id uint) error {
	return r.setArchived(ctx, id, true)
}

// Unarchive はプロジェクトのアーカイブを解除する。
func (r *projectRepository) Unarchive(ctx context.Context, id uint) error {
	return r.setArchived(ctx, id, false)
}

// setArchived はアーカイブ状態を更新する共通処理。
func (r *projectRepository) setArchived(ctx context.Context, id uint, archived bool) error {
	return r.db.WithContext(ctx).Model(&model.Project{}).
		Where("id = ?", id).
		Update("is_archived", archived).Error
}
