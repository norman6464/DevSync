package dto

import "github.com/norman6464/devsync/backend/internal/model"

// YouTubeSearchResponse はYouTube動画検索のレスポンス。
type YouTubeSearchResponse struct {
	Videos []model.YouTubeVideo `json:"videos"`
	Query  string               `json:"query"`
	Cached bool                 `json:"cached"`
	Total  int                  `json:"total"`
}

// YouTubeRecommendResponse はおすすめ動画のレスポンス。
type YouTubeRecommendResponse struct {
	Videos    []model.YouTubeVideo `json:"videos"`
	Skills    []string             `json:"skills"`
	Available bool                 `json:"available"`
}
