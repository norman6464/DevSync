package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type LearningResourceRepository struct {
	db *gorm.DB
}

func NewLearningResourceRepository(db *gorm.DB) *LearningResourceRepository {
	return &LearningResourceRepository{db: db}
}

func (r *LearningResourceRepository) Create(resource *model.LearningResource) error {
	return r.db.Create(resource).Error
}

func (r *LearningResourceRepository) FindByID(id uint) (*model.LearningResource, error) {
	var resource model.LearningResource
	err := r.db.Preload("User").First(&resource, id).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *LearningResourceRepository) FindByUserID(userID uint, includePrivate bool) ([]model.LearningResource, error) {
	var resources []model.LearningResource
	query := r.db.Where("user_id = ?", userID)
	if !includePrivate {
		query = query.Where("is_public = ?", true)
	}
	err := query.Order("created_at DESC").Find(&resources).Error
	return resources, err
}

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

func (r *LearningResourceRepository) Update(resource *model.LearningResource) error {
	return r.db.Save(resource).Error
}

func (r *LearningResourceRepository) Delete(id uint) error {
	return r.db.Delete(&model.LearningResource{}, id).Error
}

func (r *LearningResourceRepository) Search(query string, limit, offset int) ([]model.LearningResource, int64, error) {
	var resources []model.LearningResource
	var total int64

	searchQuery := "%" + query + "%"
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

// Like operations
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

func (r *LearningResourceRepository) Unlike(userID, resourceID uint) error {
	err := r.db.Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Delete(&model.ResourceLike{}).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

func (r *LearningResourceRepository) HasLiked(userID, resourceID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.ResourceLike{}).
		Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Count(&count).Error
	return count > 0, err
}

// Save/Bookmark operations
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

func (r *LearningResourceRepository) Unsave(userID, resourceID uint) error {
	err := r.db.Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Delete(&model.ResourceSave{}).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.LearningResource{}).Where("id = ?", resourceID).
		UpdateColumn("save_count", gorm.Expr("GREATEST(save_count - 1, 0)")).Error
}

func (r *LearningResourceRepository) HasSaved(userID, resourceID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.ResourceSave{}).
		Where("user_id = ? AND resource_id = ?", userID, resourceID).
		Count(&count).Error
	return count > 0, err
}

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
