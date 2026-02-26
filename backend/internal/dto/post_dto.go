package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CodeSnippetInput はコードスニペットの入力データ。
type CodeSnippetInput struct {
	Language string `json:"language" binding:"omitempty,max=100"`
	FileName string `json:"file_name" binding:"omitempty,max=255"`
	Code     string `json:"code" binding:"omitempty,max=50000"`
}

// CreatePostRequest は投稿作成リクエスト。
type CreatePostRequest struct {
	Title        string             `json:"title" binding:"required,min=1,max=200"`
	Content      string             `json:"content" binding:"required,min=1,max=50000"`
	ImageURLs    string             `json:"image_urls" binding:"omitempty,max=2000"`
	IsDraft      bool               `json:"is_draft"`
	CodeSnippets []CodeSnippetInput `json:"code_snippets" binding:"omitempty,max=20"`
}

// UpdatePostRequest は投稿更新リクエスト。
type UpdatePostRequest struct {
	Title     string `json:"title" binding:"omitempty,min=1,max=200"`
	Content   string `json:"content" binding:"omitempty,min=1,max=50000"`
	ImageURLs string `json:"image_urls" binding:"omitempty,max=2000"`
}

// CreateCommentRequest はコメント作成リクエスト。
// ParentIDが指定された場合は返信コメントとして作成する。
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=5000"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

// UpdateCommentRequest はコメント編集リクエスト。
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"`
}

// PostDetailResponse は投稿詳細レスポンス（いいね済み・ブックマーク済みフラグ付き）。
type PostDetailResponse struct {
	model.Post
	Liked      bool `json:"liked"`
	Bookmarked bool `json:"bookmarked"`
}

// PostListResponse は投稿一覧レスポンス（ページネーション付き）。
type PostListResponse struct {
	Posts  []model.Post `json:"posts"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

// AutoSaveDraftRequest は下書き自動保存リクエスト。
// IDが0の場合は新規作成、0以外の場合は既存下書きの更新。
type AutoSaveDraftRequest struct {
	ID        uint   `json:"id"`
	Title     string `json:"title" binding:"omitempty,max=200"`
	Content   string `json:"content" binding:"omitempty,max=50000"`
	ImageURLs string `json:"image_urls" binding:"omitempty,max=2000"`
}

// AutoSaveDraftResponse は下書き自動保存レスポンス。
type AutoSaveDraftResponse struct {
	ID        uint   `json:"id"`
	UpdatedAt string `json:"updated_at"`
}

// ReactionRequest はリアクション追加/削除リクエスト。
type ReactionRequest struct {
	Emoji string `json:"emoji" binding:"required,max=10"`
}

// ReactionResponse はリアクション一覧レスポンス。
type ReactionResponse struct {
	Reactions     []model.ReactionCount `json:"reactions"`
	UserReactions []string              `json:"user_reactions"`
}

// BatchReactionRequest はリアクション一括取得リクエスト。
type BatchReactionRequest struct {
	PostIDs []uint `json:"post_ids" binding:"required"`
}

// BatchReactionResponse はリアクション一括取得レスポンス。
type BatchReactionResponse struct {
	Reactions map[uint]ReactionResponse `json:"reactions"`
}
