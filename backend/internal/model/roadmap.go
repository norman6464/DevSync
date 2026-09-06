package model

import "time"

// StepOrder はステップの並び替え情報を表す。
// ロードマップ・スタディサークルのステップ並び替えに使用する。
type StepOrder struct {
	StepID     uint `json:"step_id"`     // ステップID
	OrderIndex int  `json:"order_index"` // 新しい表示順序
}

// RoadmapCategory は学習ロードマップのカテゴリを表す型。
type RoadmapCategory string

// 学習ロードマップのカテゴリ定数群。
const (
	RoadmapCategoryLanguage  RoadmapCategory = "language"  // プログラミング言語
	RoadmapCategoryFramework RoadmapCategory = "framework" // フレームワーク
	RoadmapCategorySkill     RoadmapCategory = "skill"     // スキル
	RoadmapCategoryProject   RoadmapCategory = "project"   // プロジェクト
	RoadmapCategoryOther     RoadmapCategory = "other"     // その他
)

// RoadmapStatus は学習ロードマップの進行状態を表す型。
type RoadmapStatus string

// 学習ロードマップのステータス定数群。
const (
	RoadmapStatusActive    RoadmapStatus = "active"    // 進行中
	RoadmapStatusCompleted RoadmapStatus = "completed" // 完了
)

// Roadmap はユーザーの学習ロードマップを表す。
// 複数の RoadmapStep を持ち、Progress は全ステップの完了率から自動計算される。
// IsPublic がtrueの場合、他のユーザーにも公開される。
type Roadmap struct {
	ID                 uint            `json:"id"`
	UserID             uint            `json:"user_id"`
	User               User            `json:"user,omitempty"`
	Title              string          `json:"title"`
	Description        string          `json:"description"`
	Category           RoadmapCategory `json:"category"`
	IsPublic           bool            `json:"is_public"`            // 公開フラグ
	IsTemplate         bool            `json:"is_template"`          // テンプレートフラグ
	StepCount          int             `json:"step_count"`           // 全ステップ数
	CompletedStepCount int             `json:"completed_step_count"` // 完了済みステップ数
	Progress           int             `json:"progress"`             // 進捗率（0〜100、自動計算）
	Status             RoadmapStatus   `json:"status"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	CompletedAt        *time.Time      `json:"completed_at"` // 完了日時（完了時に自動設定）

	// リレーション
	Steps []RoadmapStep `json:"steps,omitempty"`
}

// RoadmapStep は学習ロードマップ内の個別ステップを表す。
// OrderIndex で表示順序を管理し、IsCompleted で完了状態を追跡する。
type RoadmapStep struct {
	ID          uint       `json:"id"`
	RoadmapID   uint       `json:"roadmap_id"` // 所属するロードマップのID
	Title       string     `json:"title"`
	Description string     `json:"description"`
	OrderIndex  int        `json:"order_index"`  // 表示順序（0始まり）
	IsCompleted bool       `json:"is_completed"` // 完了フラグ
	CompletedAt *time.Time `json:"completed_at"` // 完了日時
	ResourceURL string     `json:"resource_url"` // 参考リソースのURL
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// RoadmapStats はユーザーのロードマップ統計情報を表す。
// DBテーブルには対応せず、集計結果を格納するDTO。
type RoadmapStats struct {
	TotalRoadmaps     int `json:"total_roadmaps"`     // ロードマップ総数
	ActiveRoadmaps    int `json:"active_roadmaps"`    // 進行中のロードマップ数
	CompletedRoadmaps int `json:"completed_roadmaps"` // 完了済みロードマップ数
	TotalSteps        int `json:"total_steps"`        // 全ステップ総数
	CompletedSteps    int `json:"completed_steps"`    // 完了済みステップ総数
}
