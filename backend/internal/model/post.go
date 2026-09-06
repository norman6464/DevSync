package model

import "time"

// Post はユーザーの投稿（学習報告など）を表す。
type Post struct {
	ID                uint          `json:"id"`
	UserID            uint          `json:"user_id"`
	User              User          `json:"user"`
	Title             string        `json:"title"`
	Content           string        `json:"content"`
	ImageURLs         string        `json:"image_urls"` // カンマ区切りの画像URL
	IsDraft           bool          `json:"is_draft"`   // 下書きフラグ
	LikeCount         int           `json:"like_count"`
	CommentCount      int           `json:"comment_count"`
	BookmarkCount     int           `json:"bookmark_count"`
	ViewCount         int           `json:"view_count"`
	EstimatedReadTime int           `json:"estimated_read_time"`
	ScheduledAt       *time.Time    `json:"scheduled_at,omitempty"`
	CodeSnippets      []CodeSnippet `json:"code_snippets,omitempty"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// Like は投稿への「いいね」を記録する。
// uniqueIndex制約でユーザーごとに1投稿1いいねを保証する。
type Like struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Bookmark は投稿のブックマーク（保存）を記録する。
// uniqueIndex制約でユーザーごとに1投稿1ブックマークを保証する。
type Bookmark struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	Post      Post      `json:"post,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Reaction は投稿への絵文字リアクションを記録する。
// uniqueIndex制約でユーザーごとに1投稿1絵文字1リアクションを保証する。
type Reaction struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

// ReactionCount はリアクション種別ごとの集計。
type ReactionCount struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
}

// TopReactedPost はリアクション数が多い投稿の要約。
type TopReactedPost struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	ReactionCount int    `json:"reaction_count"`
}

// ReactionSummary はユーザーのリアクションサマリー。
type ReactionSummary struct {
	EmojiCounts    []ReactionCount  `json:"emoji_counts"`
	TopPosts       []TopReactedPost `json:"top_posts"`
	TotalReactions int              `json:"total_reactions"`
}

// PostSeries は関連する投稿をシリーズとしてグループ化する。
type PostSeries struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PostSeriesItem はシリーズ内の投稿を表す。
// uniqueIndex制約でシリーズ内に同じ投稿が重複しないことを保証する。
type PostSeriesItem struct {
	ID         uint `json:"id"`
	SeriesID   uint `json:"series_id"`
	PostID     uint `json:"post_id"`
	Post       Post `json:"post,omitempty"`
	OrderIndex int  `json:"order_index"`
}

// PostCollection はテーマ別の投稿コレクションを表す。
// ユーザーが自他の投稿をテーマ別にまとめて管理・公開できる。
type PostCollection struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PostCollectionItem はコレクション内の投稿を表す。
// uniqueIndex制約でコレクション内に同じ投稿が重複しないことを保証する。
// インデックス名は PostgreSQL のスキーマ内で一意でなければならない。
// BookmarkCollectionItem も (collection_id, post_id) の複合ユニークを持つため、
// 名前が衝突しないよう別名にしている。
type PostCollectionItem struct {
	ID           uint      `json:"id"`
	CollectionID uint      `json:"collection_id"`
	PostID       uint      `json:"post_id"`
	Post         Post      `json:"post,omitempty"`
	Note         string    `json:"note"`
	OrderIndex   int       `json:"order_index"`
	CreatedAt    time.Time `json:"created_at"`
}

// PostTag は投稿に付与されたタグを表す。
// uniqueIndex制約で同一投稿に同じタグが重複しないことを保証する。
type PostTag struct {
	ID     uint   `json:"id"`
	PostID uint   `json:"post_id"`
	Tag    string `json:"tag"`
}

// TagCount はタグとその使用回数を表す集計結果。
type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// PostPin はプロフィールにピン留めされた投稿を表す。
// uniqueIndex制約でユーザーごとに同じ投稿が重複してピン留めされないことを保証する。
type PostPin struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	Post      Post      `json:"post,omitempty"`
	PinOrder  int       `json:"pin_order"`
	CreatedAt time.Time `json:"created_at"`
}

// PostView は投稿の閲覧記録を表す。
// uniqueIndex制約でユーザーごとに1投稿1閲覧レコードを保証する。
type PostView struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ViewCount は投稿の閲覧数集計結果を表す。
type ViewCount struct {
	PostID uint `json:"post_id"`
	Count  int  `json:"count"`
}

// Comment は投稿へのコメントを表す。
// ParentID が設定されている場合は返信コメント（スレッド）。
type Comment struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user"`
	PostID    uint      `json:"post_id"`
	ParentID  *uint     `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	IsHidden  bool      `json:"is_hidden"`
	Replies   []Comment `json:"replies,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentLike はコメントへの「いいね」を記録する。
// uniqueIndex制約でユーザーごとに1コメント1いいねを保証する。
type CommentLike struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	CommentID uint      `json:"comment_id"`
	CreatedAt time.Time `json:"created_at"`
}
