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

// DataResponse は汎用データレスポンス
type DataResponse[T any] struct {
	Data T `json:"data"`
}

// ListResponse はリストデータレスポンス（ページネーション情報付き）
type ListResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total,omitempty"`
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}
