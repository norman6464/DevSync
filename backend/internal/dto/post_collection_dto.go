package dto

// CreatePostCollectionRequest はコレクション作成のリクエストボディ。
type CreatePostCollectionRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

// UpdatePostCollectionRequest はコレクション更新のリクエストボディ。
type UpdatePostCollectionRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

// AddPostToCollectionRequest はコレクションへの投稿追加リクエスト。
type AddPostToCollectionRequest struct {
	PostID uint   `json:"post_id" binding:"required"`
	Note   string `json:"note"`
}
