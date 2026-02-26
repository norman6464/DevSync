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
	TargetHours int          `json:"target_hours" gorm:"default:0"`     // 目標学習時間（時間単位、0=未設定）
	Status      GoalStatus   `json:"status" gorm:"default:'active'"`
	IsPublic    bool         `json:"is_public" gorm:"default:false"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt *time.Time   `json:"completed_at"`                      // 完了日時（完了時に自動設定）
}

// GoalDeadlineAlert はデッドラインが近い・超過した目標のアラート情報を表す。
type GoalDeadlineAlert struct {
	Goal     LearningGoal `json:"goal"`
	Status   string       `json:"status"`    // "overdue" or "approaching"
	DaysLeft int          `json:"days_left"` // 残り日数（超過時は負数）
}

// GoalProgress は学習ゴールの実績時間 vs 目標時間の進捗情報を表す。
type GoalProgress struct {
	GoalID       uint `json:"goal_id"`
	TargetHours  int  `json:"target_hours"`  // 目標時間（時間単位）
	ActualMinutes int `json:"actual_minutes"` // 実績時間（分単位）
	Percentage   int  `json:"percentage"`     // 進捗率（0-100）
}

// LearningGoalStats はユーザーの学習目標統計情報を表す。
type LearningGoalStats struct {
	TotalGoals      int `json:"total_goals"`
	ActiveGoals     int `json:"active_goals"`
	CompletedGoals  int `json:"completed_goals"`
	AverageProgress int `json:"average_progress"` // アクティブ目標の平均進捗率
}
