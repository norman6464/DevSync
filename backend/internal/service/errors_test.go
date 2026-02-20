package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// validateRequiredID テスト
// ============================================================

func TestValidateRequiredID_ValidID(t *testing.T) {
	err := validateRequiredID(1, "userID")
	assert.NoError(t, err)
}

func TestValidateRequiredID_LargeID(t *testing.T) {
	err := validateRequiredID(999999, "postID")
	assert.NoError(t, err)
}

func TestValidateRequiredID_ZeroID(t *testing.T) {
	err := validateRequiredID(0, "userID")
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
	assert.Contains(t, domainErr.Message, "userIDは必須です")
}

func TestValidateRequiredID_ZeroID_PostID(t *testing.T) {
	err := validateRequiredID(0, "postID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "postIDは必須です")
}
