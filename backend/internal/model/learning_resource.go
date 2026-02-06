package model

import (
	"time"

	"gorm.io/gorm"
)

type ResourceCategory string

const (
	ResourceCategoryBook     ResourceCategory = "book"
	ResourceCategoryVideo    ResourceCategory = "video"
	ResourceCategoryArticle  ResourceCategory = "article"
	ResourceCategoryCourse   ResourceCategory = "course"
	ResourceCategoryTutorial ResourceCategory = "tutorial"
	ResourceCategoryPodcast  ResourceCategory = "podcast"
	ResourceCategoryTool     ResourceCategory = "tool"
	ResourceCategoryOther    ResourceCategory = "other"
)

type ResourceDifficulty string

const (
	ResourceDifficultyBeginner     ResourceDifficulty = "beginner"
	ResourceDifficultyIntermediate ResourceDifficulty = "intermediate"
	ResourceDifficultyAdvanced     ResourceDifficulty = "advanced"
)

type LearningResource struct {
	ID          uint               `json:"id" gorm:"primaryKey"`
	UserID      uint               `json:"user_id" gorm:"not null;index"`
	User        User               `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title       string             `json:"title" gorm:"not null;size:300"`
	Description string             `json:"description" gorm:"type:text"`
	URL         string             `json:"url" gorm:"size:500"`
	Category    ResourceCategory   `json:"category" gorm:"size:50;not null"`
	Difficulty  ResourceDifficulty `json:"difficulty" gorm:"size:50"`
	Tags        string             `json:"tags" gorm:"type:text"` // JSON array of tags
	ImageURL    string             `json:"image_url" gorm:"size:500"`
	IsPublic    bool               `json:"is_public" gorm:"default:true"`
	LikeCount   int                `json:"like_count" gorm:"default:0"`
	SaveCount   int                `json:"save_count" gorm:"default:0"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `json:"-" gorm:"index"`
}

// ResourceLike tracks likes on learning resources
type ResourceLike struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_resource_like"`
	ResourceID uint      `json:"resource_id" gorm:"not null;uniqueIndex:idx_resource_like"`
	CreatedAt  time.Time `json:"created_at"`
}

// ResourceSave tracks saved/bookmarked resources by users
type ResourceSave struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_resource_save"`
	ResourceID uint      `json:"resource_id" gorm:"not null;uniqueIndex:idx_resource_save"`
	CreatedAt  time.Time `json:"created_at"`
}
