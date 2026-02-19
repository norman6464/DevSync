package model

// ProjectStats はユーザーのプロジェクト活動集計統計を表す。
type ProjectStats struct {
	TotalProjects     int64 `json:"total_projects"`
	FeaturedProjects  int64 `json:"featured_projects"`
	OngoingProjects   int64 `json:"ongoing_projects"`
	CompletedProjects int64 `json:"completed_projects"`
}
