package model

// BookReviewStats はユーザーの書籍レビュー集計統計を表す。
type BookReviewStats struct {
	TotalReviews   int64   `json:"total_reviews"`
	AverageRating  float64 `json:"average_rating"`
	MaxRating      int     `json:"max_rating"`
	MinRating      int     `json:"min_rating"`
	FiveStarCount  int64   `json:"five_star_count"`
}
