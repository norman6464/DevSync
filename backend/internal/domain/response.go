// Package domain provides core business logic and domain models.
package domain

// SuccessResponse は成功レスポンスの統一構造体。
type SuccessResponse struct {
	Data interface{} `json:"data,omitempty"`
}

// ErrorResponse はエラーレスポンスの統一構造体。
type ErrorResponse struct {
	Error string      `json:"error"`
	Code  string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// PaginatedResponse はページネーション付きレスポンスの統一構造体。
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	TotalPages int         `json:"total_pages,omitempty"`
}

// MessageResponse はメッセージのみのレスポンス構造体。
type MessageResponse struct {
	Message string `json:"message"`
}

// NewSuccessResponse は成功レスポンスを生成する。
func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{Data: data}
}

// NewErrorResponse はエラーレスポンスを生成する。
func NewErrorResponse(message string, code string, details interface{}) ErrorResponse {
	return ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	}
}

// NewPaginatedResponse はページネーション付きレスポンスを生成する。
func NewPaginatedResponse(data interface{}, total int64, page, limit int) PaginatedResponse {
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// NewMessageResponse はメッセージレスポンスを生成する。
func NewMessageResponse(message string) MessageResponse {
	return MessageResponse{Message: message}
}
