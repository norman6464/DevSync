package model

import "time"

// SearchSortBy は検索結果のソート順を表す。
type SearchSortBy string

const (
	SearchSortByLatest  SearchSortBy = "latest"  // 最新順
	SearchSortByPopular SearchSortBy = "popular" // 人気順（いいね数）
	SearchSortByViews   SearchSortBy = "views"   // 閲覧数順
)

// ValidSearchSortBy は有効なソート順の集合。
var ValidSearchSortBy = map[SearchSortBy]bool{
	SearchSortByLatest:  true,
	SearchSortByPopular: true,
	SearchSortByViews:   true,
}

// PostSearchParams は投稿検索のパラメータ。
type PostSearchParams struct {
	Query    string
	Tags     []string
	SortBy   SearchSortBy
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}

// PostSearchResult は投稿検索結果のレスポンス。
type PostSearchResult struct {
	Posts  []Post `json:"posts"`
	Total  int64  `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
