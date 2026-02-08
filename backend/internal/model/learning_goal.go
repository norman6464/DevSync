package model

import "time"

// GoalStatus は学習目標の進行状態を表す型。
type GoalStatus string

// 学習目標のステータス定数群。
const (
	GoalStatusActive    GoalStatus = "active"    // 進行中
	GoalStatusCompleted GoalStatus = "completed" // 完了
	GoalStatusPaused    GoalStatus = "paused"    // 一時停止
)

// GoalCategory は学習目標のカテゴリを表す型。
type GoalCategory string

// 学習目標のカテゴリ定数群。
const (
	GoalCategoryLanguage  GoalCategory = "language"  // プログラミング言語
	GoalCategoryFramework GoalCategory = "framework" // フレームワーク
	GoalCategorySkill     GoalCategory = "skill"     // スキル
	GoalCategoryProject   GoalCategory = "project"   // プロジェクト
	GoalCategoryOther     GoalCategory = "other"     // その他
)

// LearningGoal はユーザーの学習目標を表す。
// Progress が100%に達すると、Serviceが自動的にStatusをcompletedに変更する。
type LearningGoal struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	UserID      uint         `json:"user_id" gorm:"not null;index"`
	Title       string       `json:"title" gorm:"not null"`
	Description string       `json:"description"`
	Category    GoalCategory `json:"category" gorm:"default:'other'"`
	TargetDate  *time.Time   `json:"target_date"`                       // 目標達成予定日
	Progress    int          `json:"progress" gorm:"default:0"`         // 0〜100の達成率
	Status      GoalStatus   `json:"status" gorm:"default:'active'"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt *time.Time   `json:"completed_at"`                      // 完了日時（完了時に自動設定）
}

// LearningGoalStats はユーザーの学習目標統計情報を表す。
type LearningGoalStats struct {
	TotalGoals      int `json:"total_goals"`
	ActiveGoals     int `json:"active_goals"`
	CompletedGoals  int `json:"completed_goals"`
	AverageProgress int `json:"average_progress"` // アクティブ目標の平均進捗率
}
