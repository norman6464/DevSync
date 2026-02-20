package dto

// ErrorResponse は標準エラーレスポンス
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// MessageResponse は汎用メッセージレスポンス
type MessageResponse struct {
	Message string `json:"message"`
}

// URLResponse はURLを返すレスポンス（OAuth等で使用）
type URLResponse struct {
	URL string `json:"url"`
}

// StatusResponse はステータスレスポンス
type StatusResponse struct {
	Status string `json:"status"`
}

// CountResponse はカウントレスポンス
type CountResponse struct {
	Count int64 `json:"count"`
}

// LikeStatusResponse はいいね状態レスポンス
type LikeStatusResponse struct {
	Liked bool  `json:"liked"`
	Count int64 `json:"count"`
}

// ViewCountResponse は閲覧数レスポンス
type ViewCountResponse struct {
	PostID    uint  `json:"post_id"`
	ViewCount int64 `json:"view_count"`
}

