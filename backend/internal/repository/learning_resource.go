package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// LearningResourceRepository は学習リソースデータへのアクセスを提供するリポジトリ実装。
type LearningResourceRepository struct {
	db *gorm.DB
}

// NewLearningResourceRepository は新しいLearningResourceRepositoryインスタンスを生成する。
func NewLearningResourceRepository(db *gorm.DB) *LearningResourceRepository {
	return &LearningResourceRepository{db: db}
}

// Create は新しい学習リソースをデータベースに作成する。
func (r *LearningResourceRepository) Create(resource *model.LearningResource) error {
	return r.db.Create(resource).Error
}

// FindByID は指定IDの学習リソースをユーザー情報付きで取得する。
func (r *LearningResourceRepository) FindByID(id uint) (*model.LearningResource, error) {
	var resource model.LearningResource
	err := r.db.Preload("User").First(&resource, id).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// FindByUserID は指定ユーザーの学習リソースを取得する。
// includePrivateがfalseの場合、公開リソースのみを返す。
func (r *LearningResourceRepository) FindByUserID(userID uint, includePrivate bool, limit, offset int) ([]model.LearningResource, int64, error) {
	var resources []model.LearningResource
	var total int64
	query := r.db.Where("user_id = ?", userID)
	if !includePrivate {
		query = query.Where("is_public = ?", true)
	}
	query.Model(&model.LearningResource{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&resources).Error
	return resources, total, err
}

// FindPublic は公開学習リソースをフィルタ・ページネーション付きで取得する。
// categoryやdifficultyが指定された場合、それぞれでフィルタする。
func (r *LearningResourceRepository) FindPublic(limit, offset int, category string, difficulty string) ([]model.LearningResource, int64, error) {
	var resources []model.LearningResource
	var total int64

	query := r.db.Model(&model.LearningResource{}).Where("is_public = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	query.Count(&total)

	err := query.Preload("User").
		Order("like_count DESC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&resources).Error

	return resources, total, err
}

// Update は既存の学習リソースを更新する。
func (r *LearningResourceRepository) Update(resource *model.LearningResource) error {
	return r.db.Save(resource).Error
}

// Delete は指定IDの学習リソースを削除する。
func (r *LearningResourceRepository) Delete(id uint) error {
	return r.db.Delete(&model.LearningResource{}, id).Error
}

// Search は公開学習リソースをタイトル・説明・タグで全文検索する。
func (r *LearningResourceRepository) Search(query string, limit, offset int) ([]model.LearningResource, int64, error) {
	var resources []model.LearningResource
	var total int64

	searchQuery := EscapeLikePattern(query)
	dbQuery := r.db.Model(&model.LearningResource{}).
		Where("is_public = ?", true).
		Where("title ILIKE ? OR description ILIKE ? OR tags ILIKE ?", searchQuery, searchQuery, searchQuery)

	dbQuery.Count(&total)

	err := dbQuery.Preload("User").
		Order("like_count DESC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&resources).Error

	return resources, total, err
}

// Like は学習リソースにいいねを追加し、like_countをインクリメントする。
func (r *LearningResourceRepository) Like(userID, resourceID uint) error {
	like := &model.ResourceLike{
		UserID:     userID,
		ResourceID: resourceID,
	}
	err := r.db.Create(like).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// Unlike は学習リソースのいいねを取り消し、like_countをデクリメントする。
func (r *LearningResourceRepository) Unlike(userID, resourceID uint) error {
	err := r.db.Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Delete(&model.ResourceLike{}).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

// HasLiked は指定ユーザーが学習リソースにいいね済みかどうかを判定する。
func (r *LearningResourceRepository) HasLiked(userID, resourceID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.ResourceLike{}).
		Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Count(&count).Error
	return count > 0, err
}

// Save は学習リソースをブックマーク（保存）し、save_countをインクリメントする。
func (r *LearningResourceRepository) Save(userID, resourceID uint) error {
	save := &model.ResourceSave{
		UserID:     userID,
		ResourceID: resourceID,
	}
	err := r.db.Create(save).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("save_count", gorm.Expr("save_count + 1")).Error
}

// Unsave は学習リソースのブックマークを解除し、save_countをデクリメントする。
func (r *LearningResourceRepository) Unsave(userID, resourceID uint) error {
	err := r.db.Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Delete(&model.ResourceSave{}).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("save_count", gorm.Expr("GREATEST(save_count - 1, 0)")).Error
}

// HasSaved は指定ユーザーが学習リソースをブックマーク済みかどうかを判定する。
func (r *LearningResourceRepository) HasSaved(userID, resourceID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.ResourceSave{}).
		Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Count(&count).Error
	return count > 0, err
}

// FindSavedByUserID は指定ユーザーがブックマークした学習リソースをページネーション付きで取得する。
func (r *LearningResourceRepository) FindSavedByUserID(userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	var resources []model.LearningResource
	var total int64

	subQuery := r.db.Model(&model.ResourceSave{}).Select("resource_id").Where("user_id = ?", userID)

	r.db.Model(&model.LearningResource{}).Where("id IN (?)", subQuery).Count(&total)

	err := r.db.Preload("User").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&resources).Error

	return resources, total, err
}

// FindByDifficulty は公開リソースを難易度でフィルタリングして取得する。
func (r *LearningResourceRepository) FindByDifficulty(difficulty string, limit, offset int) ([]model.LearningResource, int64, error) {
	var resources []model.LearningResource
	var total int64

	query := r.db.Where("is_public = ? AND difficulty = ?", true, difficulty)
	query.Model(&model.LearningResource{}).Count(&total)

	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&resources).Error

	return resources, total, err
}
