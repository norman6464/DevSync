// Package domain provides domain-level types and errors for the DevSync application.
package domain

import (
	"errors"
	"fmt"
	// DomainError.HTTPStatus がステータス定数を使うための例外。
	// この対応表は本来 handler 側に置くべきで、移設は別途対応する。
	"net/http" //archlint:allow
)

// ErrorCode represents application-specific error codes.
type ErrorCode string

const (
	// Authentication errors
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"

	// Resource errors
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeConflict      ErrorCode = "CONFLICT"

	// Validation errors
	ErrCodeValidation ErrorCode = "VALIDATION_ERROR"
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"

	// System errors
	ErrCodeInternal           ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase           ErrorCode = "DATABASE_ERROR"
	ErrCodeRateLimitExceeded  ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// DomainError represents a domain-level error with additional context.
type DomainError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface.
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error.
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is は errors.Is で同じエラーコードの DomainError を一致させるために使用する。
func (e *DomainError) Is(target error) bool {
	t, ok := target.(*DomainError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// HTTPStatus returns the appropriate HTTP status code for the error.
func (e *DomainError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeAlreadyExists, ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeValidation, ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeRateLimitExceeded:
		return http.StatusTooManyRequests
	case ErrCodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrCodeDatabase, ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// NewError creates a new DomainError.
func NewError(code ErrorCode, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined domain errors
var (
	ErrUnauthorized       = NewError(ErrCodeUnauthorized, "認証が必要です", nil)
	ErrForbidden          = NewError(ErrCodeForbidden, "この操作を実行する権限がありません", nil)
	ErrNotFound           = NewError(ErrCodeNotFound, "リソースが見つかりません", nil)
	ErrBadRequest         = NewError(ErrCodeBadRequest, "不正なリクエストです", nil)
	ErrInternal           = NewError(ErrCodeInternal, "内部エラーが発生しました", nil)
	ErrAlreadyExists      = NewError(ErrCodeAlreadyExists, "リソースが既に存在します", nil)
	ErrConflict           = NewError(ErrCodeConflict, "競合が発生しました", nil)
	ErrValidation         = NewError(ErrCodeValidation, "バリデーションエラー", nil)
	ErrDatabase           = NewError(ErrCodeDatabase, "データベースエラーが発生しました", nil)
	ErrRateLimitExceeded  = NewError(ErrCodeRateLimitExceeded, "レート制限を超過しました", nil)
	ErrServiceUnavailable = NewError(ErrCodeServiceUnavailable, "サービスが利用できません", nil)
)

// IsDomainError checks if an error is a DomainError.
func IsDomainError(err error) bool {
	var domainErr *DomainError
	return errors.As(err, &domainErr)
}

// GetDomainError extracts the DomainError from an error.
func GetDomainError(err error) *DomainError {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr
	}
	return nil
}
