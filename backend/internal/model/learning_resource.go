package model

import (
	"time"

	"gorm.io/gorm"
)

// ResourceCategory は学習リソースの分類を表す型。
type ResourceCategory string

// 学習リソースのカテゴリ定数群。
const (
	ResourceCategoryBook     ResourceCategory = "book"     // 書籍
	ResourceCategoryVideo    ResourceCategory = "video"    // 動画
	ResourceCategoryArticle  ResourceCategory = "article"  // 記事
	ResourceCategoryCourse   ResourceCategory = "course"   // 講座・コース
	ResourceCategoryTutorial ResourceCategory = "tutorial" // チュートリアル
	ResourceCategoryPodcast  ResourceCategory = "podcast"  // ポッドキャスト
	ResourceCategoryTool     ResourceCategory = "tool"     // ツール
	ResourceCategoryOther    ResourceCategory = "other"    // その他
)

// ResourceDifficulty は学習リソースの難易度を表す型。
type ResourceDifficulty string

// 学習リソースの難易度定数群。
const (
	ResourceDifficultyBeginner     ResourceDifficulty = "beginner"     // 初級
	ResourceDifficultyIntermediate ResourceDifficulty = "intermediate" // 中級
	ResourceDifficultyAdvanced     ResourceDifficulty = "advanced"     // 上級
)

// LearningResource はユーザーが登録した学習リソースを表す。
// IsPublic フラグで公開/非公開を制御する。
type LearningResource struct {
	ID          uint               `json:"id" gorm:"primaryKey"`
	UserID      uint               `json:"user_id" gorm:"not null;index"`
	User        User               `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title       string             `json:"title" gorm:"not null;size:300"`
	Description string             `json:"description" gorm:"type:text"`
	URL         string             `json:"url" gorm:"size:500"`
	Category    ResourceCategory   `json:"category" gorm:"size:50;not null"`
	Difficulty  ResourceDifficulty `json:"difficulty" gorm:"size:50"`
	Tags        string             `json:"tags" gorm:"type:text"`          // JSON配列形式のタグ
	ImageURL    string             `json:"image_url" gorm:"size:500"`
	IsPublic    bool               `json:"is_public" gorm:"default:true"` // 他ユーザーに公開するか
	LikeCount   int                `json:"like_count" gorm:"default:0"`
	SaveCount   int                `json:"save_count" gorm:"default:0"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `json:"-" gorm:"index"` // 論理削除用
}

// ResourceLike は学習リソースへの「いいね」を記録する。
type ResourceLike struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_resource_like"`
	ResourceID uint      `json:"resource_id" gorm:"not null;uniqueIndex:idx_resource_like"`
	CreatedAt  time.Time `json:"created_at"`
}

// ResourceSave は学習リソースのブックマーク（保存）を記録する。
type ResourceSave struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_resource_save"`
	ResourceID uint      `json:"resource_id" gorm:"not null;uniqueIndex:idx_resource_save"`
	CreatedAt  time.Time `json:"created_at"`
}
