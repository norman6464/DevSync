package model

// BookmarkStats はユーザーのブックマーク集計統計を表す。
type BookmarkStats struct {
	TotalBookmarksMade     int64 `json:"total_bookmarks_made"`
	TotalBookmarksReceived int64 `json:"total_bookmarks_received"`
	BookmarksThisMonth     int64 `json:"bookmarks_this_month"`
}
