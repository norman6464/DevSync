package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// learningResourceRepository は [repository.LearningResourceRepository] の GORM 実装。
//
// 旧 repository パッケージにも同じテーブルを扱う実装が残っている。ai_advice がまだそちらを
// 使っているため、移行が一巡するまで新旧のアダプタが並存する。
type learningResourceRepository struct {
	db *gorm.DB
}

// NewLearningResourceRepository は LearningResourceRepository の GORM 実装を返す。
func NewLearningResourceRepository(db *gorm.DB) repository.LearningResourceRepository {
	return &learningResourceRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningResourceRepository = (*learningResourceRepository)(nil)

// Create は新しい学習リソースを作成する。
func (r *learningResourceRepository) Create(ctx context.Context, resource *model.LearningResource) error {
	return r.db.WithContext(ctx).Create(resource).Error
}

// Update は既存の学習リソースを更新する。
func (r *learningResourceRepository) Update(ctx context.Context, resource *model.LearningResource) error {
	return r.db.WithContext(ctx).Save(resource).Error
}

// Delete は学習リソースを削除する。
func (r *learningResourceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.LearningResource{}, id).Error
}

// FindByID は指定 ID の学習リソースをユーザー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *learningResourceRepository) FindByID(ctx context.Context, id uint) (*model.LearningResource, error) {
	var resource model.LearningResource
	err := r.db.WithContext(ctx).Preload("User").First(&resource, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// FindByUserID は指定ユーザーのリソースを取得する（新しい順）。
// 一覧系の中でこれだけユーザー情報をプリロードしない（移行前からの挙動）。
func (r *learningResourceRepository) FindByUserID(ctx context.Context, userID uint, includePrivate bool, limit, offset int) ([]model.LearningResource, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.LearningResource{}).Where("user_id = ?", userID)
	if !includePrivate {
		scope = scope.Where("is_public = ?", true)
	}

	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var resources []model.LearningResource
	err := scope.Session(&gorm.Session{}).
		Order("created_at DESC").Limit(limit).Offset(offset).
		Find(&resources).Error
	return resources, total, err
}

// FindPublic は公開リソースをカテゴリ・難易度で絞り込んで取得する（いいね数降順 → 新しい順）。
func (r *learningResourceRepository) FindPublic(ctx context.Context, limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.LearningResource{}).Where("is_public = ?", true)
	if category != "" {
		scope = scope.Where("category = ?", category)
	}
	if difficulty != "" {
		scope = scope.Where("difficulty = ?", difficulty)
	}
	return r.paginatedResources(scope, "like_count DESC, created_at DESC", limit, offset)
}

// FindByDifficulty は公開リソースを難易度で絞り込んで取得する（新しい順）。
func (r *learningResourceRepository) FindByDifficulty(ctx context.Context, difficulty string, limit, offset int) ([]model.LearningResource, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.LearningResource{}).
		Where("is_public = ? AND difficulty = ?", true, difficulty)
	return r.paginatedResources(scope, "created_at DESC", limit, offset)
}

// Search は公開リソースをタイトル・説明・タグで部分一致検索する（いいね数降順 → 新しい順）。
func (r *learningResourceRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.LearningResource, int64, error) {
	pattern := escapeLikePattern(query)
	scope := r.db.WithContext(ctx).Model(&model.LearningResource{}).
		Where("is_public = ?", true).
		Where("title ILIKE ? OR description ILIKE ? OR tags ILIKE ?", pattern, pattern, pattern)
	return r.paginatedResources(scope, "like_count DESC, created_at DESC", limit, offset)
}

// FindSavedByUserID は指定ユーザーが保存したリソースを取得する（新しい順）。
func (r *learningResourceRepository) FindSavedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	saved := r.db.WithContext(ctx).Model(&model.ResourceSave{}).
		Select("resource_id").Where("user_id = ?", userID)
	scope := r.db.WithContext(ctx).Model(&model.LearningResource{}).Where("id IN (?)", saved)
	return r.paginatedResources(scope, "created_at DESC", limit, offset)
}

// paginatedResources は絞り込み済みクエリに対して総件数とページを取得する共通処理。
func (r *learningResourceRepository) paginatedResources(scope *gorm.DB, orderClause string, limit, offset int) ([]model.LearningResource, int64, error) {
	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var resources []model.LearningResource
	err := scope.Session(&gorm.Session{}).
		Preload("User").
		Order(orderClause).Limit(limit).Offset(offset).
		Find(&resources).Error
	return resources, total, err
}

// CountByUserID は指定ユーザーのリソース総数を返す。
func (r *learningResourceRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.LearningResource{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// Like はいいねを追加し、リソースのいいね数を 1 増やす。
func (r *learningResourceRepository) Like(ctx context.Context, userID, resourceID uint) error {
	db := r.db.WithContext(ctx)
	if err := db.Create(&model.ResourceLike{UserID: userID, ResourceID: resourceID}).Error; err != nil {
		return err
	}
	return db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// Unlike はいいねを取り消し、リソースのいいね数を 1 減らす（0 未満にはしない）。
func (r *learningResourceRepository) Unlike(ctx context.Context, userID, resourceID uint) error {
	db := r.db.WithContext(ctx)
	if err := db.Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Delete(&model.ResourceLike{}).Error; err != nil {
		return err
	}
	return db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

// HasLiked は指定ユーザーがいいね済みかを返す。
func (r *learningResourceRepository) HasLiked(ctx context.Context, userID, resourceID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ResourceLike{}).
		Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Count(&count).Error
	return count > 0, err
}

// Save は保存を追加し、リソースの保存数を 1 増やす。
func (r *learningResourceRepository) Save(ctx context.Context, userID, resourceID uint) error {
	db := r.db.WithContext(ctx)
	if err := db.Create(&model.ResourceSave{UserID: userID, ResourceID: resourceID}).Error; err != nil {
		return err
	}
	return db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("save_count", gorm.Expr("save_count + 1")).Error
}

// Unsave は保存を取り消し、リソースの保存数を 1 減らす（0 未満にはしない）。
func (r *learningResourceRepository) Unsave(ctx context.Context, userID, resourceID uint) error {
	db := r.db.WithContext(ctx)
	if err := db.Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Delete(&model.ResourceSave{}).Error; err != nil {
		return err
	}
	return db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("save_count", gorm.Expr("GREATEST(save_count - 1, 0)")).Error
}

// HasSaved は指定ユーザーが保存済みかを返す。
func (r *learningResourceRepository) HasSaved(ctx context.Context, userID, resourceID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ResourceSave{}).
		Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Count(&count).Error
	return count > 0, err
}
