package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSearchQuery_EmptyQuery(t *testing.T) {
	_, err := validateSearchQuery("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "検索キーワードは必須です")
}

func TestValidateSearchQuery_WhitespaceOnly(t *testing.T) {
	_, err := validateSearchQuery("   \t  ")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "検索キーワードは必須です")
}

func TestValidateSearchQuery_ValidQuery(t *testing.T) {
	q, err := validateSearchQuery("  Go言語  ")
	assert.NoError(t, err)
	assert.Equal(t, "Go言語", q)
}
