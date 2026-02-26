package model

import "time"

// Post はユーザーの投稿（学習報告など）を表す。
type Post struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Title        string    `json:"title" gorm:"not null"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	ImageURLs    string    `json:"image_urls" gorm:"type:text"` // カンマ区切りの画像URL
	IsDraft      bool      `json:"is_draft" gorm:"default:false;index"` // 下書きフラグ
	LikeCount         int            `json:"like_count" gorm:"default:0"`
	CommentCount      int            `json:"comment_count" gorm:"default:0"`
	BookmarkCount     int            `json:"bookmark_count" gorm:"default:0"`
	ViewCount         int            `json:"view_count" gorm:"default:0"`
	EstimatedReadTime int            `json:"estimated_read_time" gorm:"default:0"`
	ScheduledAt       *time.Time     `json:"scheduled_at,omitempty" gorm:"index"`
	CodeSnippets      []CodeSnippet  `json:"code_snippets,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// Like は投稿への「いいね」を記録する。
// uniqueIndex制約でユーザーごとに1投稿1いいねを保証する。
type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post_like"`
	PostID    uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post_like;index"`
	CreatedAt time.Time `json:"created_at"`
}

// Bookmark は投稿のブックマーク（保存）を記録する。
// uniqueIndex制約でユーザーごとに1投稿1ブックマークを保証する。
type Bookmark struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post_bookmark"`
	PostID    uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post_bookmark;index"`
	Post      Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"created_at"`
}

// Reaction は投稿への絵文字リアクションを記録する。
// uniqueIndex制約でユーザーごとに1投稿1絵文字1リアクションを保証する。
type Reaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post_emoji"`
	PostID    uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post_emoji;index"`
	Emoji     string    `json:"emoji" gorm:"not null;uniqueIndex:idx_user_post_emoji;size:10"`
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
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PostSeriesItem はシリーズ内の投稿を表す。
// uniqueIndex制約でシリーズ内に同じ投稿が重複しないことを保証する。
type PostSeriesItem struct {
	ID         uint `json:"id" gorm:"primaryKey"`
	SeriesID   uint `json:"series_id" gorm:"not null;uniqueIndex:idx_series_post;index"`
	PostID     uint `json:"post_id" gorm:"not null;uniqueIndex:idx_series_post;index"`
	Post       Post `json:"post,omitempty" gorm:"foreignKey:PostID"`
	OrderIndex int  `json:"order_index" gorm:"not null;default:0"`
}

// PostCollection はテーマ別の投稿コレクションを表す。
// ユーザーが自他の投稿をテーマ別にまとめて管理・公開できる。
type PostCollection struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PostCollectionItem はコレクション内の投稿を表す。
// uniqueIndex制約でコレクション内に同じ投稿が重複しないことを保証する。
type PostCollectionItem struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CollectionID uint      `json:"collection_id" gorm:"not null;uniqueIndex:idx_collection_post;index"`
	PostID       uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_collection_post;index"`
	Post         Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Note         string    `json:"note" gorm:"type:text"`
	OrderIndex   int       `json:"order_index" gorm:"not null;default:0"`
	CreatedAt    time.Time `json:"created_at"`
}

// PostTag は投稿に付与されたタグを表す。
// uniqueIndex制約で同一投稿に同じタグが重複しないことを保証する。
type PostTag struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	PostID uint   `json:"post_id" gorm:"not null;uniqueIndex:idx_post_tag;index"`
	Tag    string `json:"tag" gorm:"not null;uniqueIndex:idx_post_tag;index;size:50"`
}

// TagCount はタグとその使用回数を表す集計結果。
type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// PostPin はプロフィールにピン留めされた投稿を表す。
// uniqueIndex制約でユーザーごとに同じ投稿が重複してピン留めされないことを保証する。
type PostPin struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post_pin;index"`
	PostID    uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post_pin;index"`
	Post      Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	PinOrder  int       `json:"pin_order" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at"`
}

// PostView は投稿の閲覧記録を表す。
// uniqueIndex制約でユーザーごとに1投稿1閲覧レコードを保証する。
type PostView struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post_view;index"`
	PostID    uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post_view;index"`
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
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	User      User       `json:"user" gorm:"foreignKey:UserID"`
	PostID    uint       `json:"post_id" gorm:"not null;index"`
	ParentID  *uint      `json:"parent_id,omitempty" gorm:"index"`
	Content   string     `json:"content" gorm:"type:text;not null"`
	LikeCount int        `json:"like_count" gorm:"default:0"`
	IsHidden  bool       `json:"is_hidden" gorm:"default:false"`
	Replies   []Comment  `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CommentLike はコメントへの「いいね」を記録する。
// uniqueIndex制約でユーザーごとに1コメント1いいねを保証する。
type CommentLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_comment_like"`
	CommentID uint      `json:"comment_id" gorm:"not null;uniqueIndex:idx_user_comment_like;index"`
	CreatedAt time.Time `json:"created_at"`
}
