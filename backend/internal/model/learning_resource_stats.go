package model

// LearningResourceStats はユーザーの学習リソース活動集計統計を表す。
type LearningResourceStats struct {
	TotalResources int64 `json:"total_resources"`
	TotalLikes     int64 `json:"total_likes"`
	TotalSaves     int64 `json:"total_saves"`
	CategoryCount  int64 `json:"category_count"`
}
