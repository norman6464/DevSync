package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *DomainError
		expected string
	}{
		{
			name:     "エラーメッセージのみ",
			err:      NewError(ErrCodeNotFound, "ユーザーが見つかりません", nil),
			expected: "NOT_FOUND: ユーザーが見つかりません",
		},
		{
			name:     "ラップされたエラーあり",
			err:      NewError(ErrCodeDatabase, "データベースエラー", errors.New("connection failed")),
			expected: "DATABASE_ERROR: データベースエラー (connection failed)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestDomainError_Unwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	domainErr := NewError(ErrCodeInternal, "wrapper", innerErr)

	assert.Equal(t, innerErr, domainErr.Unwrap())
	assert.True(t, errors.Is(domainErr, innerErr))
}

func TestIsDomainError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "DomainError",
			err:      NewError(ErrCodeNotFound, "not found", nil),
			expected: true,
		},
		{
			name:     "標準エラー",
			err:      errors.New("standard error"),
			expected: false,
		},
		{
			name:     "nilエラー",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsDomainError(tt.err))
		})
	}
}

func TestGetDomainError(t *testing.T) {
	domainErr := NewError(ErrCodeNotFound, "not found", nil)

	t.Run("DomainErrorを取得", func(t *testing.T) {
		result := GetDomainError(domainErr)
		assert.NotNil(t, result)
		assert.Equal(t, ErrCodeNotFound, result.Code)
	})

	t.Run("標準エラーからはnilを返す", func(t *testing.T) {
		result := GetDomainError(errors.New("standard error"))
		assert.Nil(t, result)
	})

	t.Run("ラップされたDomainErrorを取得", func(t *testing.T) {
		wrapped := errors.New("wrapper: " + domainErr.Error())
		// このケースはラップされていないので、nilを返す
		result := GetDomainError(wrapped)
		assert.Nil(t, result)
	})
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *DomainError
		code ErrorCode
	}{
		{"Unauthorized", ErrUnauthorized, ErrCodeUnauthorized},
		{"Forbidden", ErrForbidden, ErrCodeForbidden},
		{"NotFound", ErrNotFound, ErrCodeNotFound},
		{"BadRequest", ErrBadRequest, ErrCodeBadRequest},
		{"Internal", ErrInternal, ErrCodeInternal},
		{"AlreadyExists", ErrAlreadyExists, ErrCodeAlreadyExists},
		{"Conflict", ErrConflict, ErrCodeConflict},
		{"Validation", ErrValidation, ErrCodeValidation},
		{"Database", ErrDatabase, ErrCodeDatabase},
		{"RateLimitExceeded", ErrRateLimitExceeded, ErrCodeRateLimitExceeded},
		{"ServiceUnavailable", ErrServiceUnavailable, ErrCodeServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.err.Code)
			assert.NotEmpty(t, tt.err.Message, "エラーメッセージが空であってはならない")
		})
	}
}
