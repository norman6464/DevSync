package dto

// CreatePostSeriesRequest はシリーズ作成のリクエストボディ。
type CreatePostSeriesRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
}

// UpdatePostSeriesRequest はシリーズ更新のリクエストボディ。
type UpdatePostSeriesRequest struct {
	Title       string `json:"title" binding:"max=200"`
	Description string `json:"description"`
}

// AddPostToSeriesRequest はシリーズへの投稿追加リクエスト。
type AddPostToSeriesRequest struct {
	PostID     uint `json:"post_id" binding:"required"`
	OrderIndex int  `json:"order_index"`
}
