package model

import (
	"time"
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
	ID          uint               `json:"id"`
	UserID      uint               `json:"user_id"`
	User        User               `json:"user,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	URL         string             `json:"url"`
	Category    ResourceCategory   `json:"category"`
	Difficulty  ResourceDifficulty `json:"difficulty"`
	Tags        string             `json:"tags"` // JSON配列形式のタグ
	ImageURL    string             `json:"image_url"`
	IsPublic    bool               `json:"is_public"` // 他ユーザーに公開するか
	LikeCount   int                `json:"like_count"`
	SaveCount   int                `json:"save_count"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ResourceLike は学習リソースへの「いいね」を記録する。
type ResourceLike struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	ResourceID uint      `json:"resource_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// ResourceSave は学習リソースのブックマーク（保存）を記録する。
type ResourceSave struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	ResourceID uint      `json:"resource_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// ResourceReview は学習リソースへのレビュー（評価＋コメント）を記録する。
// ユーザーとリソースの組み合わせでユニーク（1リソースにつき1レビュー）。
type ResourceReview struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user,omitempty"`
	ResourceID uint      `json:"resource_id"`
	Rating     int       `json:"rating"`  // 1-5の評価
	Comment    string    `json:"comment"` // レビューコメント
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
