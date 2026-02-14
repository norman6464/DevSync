package validator

import (
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestResourceValidator_ValidateCreateResource(t *testing.T) {
	v := NewResourceValidator()

	tests := []struct {
		name        string
		title       string
		description string
		url         string
		category    string
		difficulty  string
		wantErr     bool
	}{
		{"有効なリソース", "リソースタイトル", "説明", "https://example.com", "book", "beginner", false},
		{"有効（説明なし）", "リソースタイトル", "", "https://example.com", "book", "beginner", false},
		{"有効（難易度なし）", "リソースタイトル", "説明", "https://example.com", "book", "", false},
		{"無効（タイトルが空）", "", "説明", "https://example.com", "book", "beginner", true},
		{"無効（URLが空）", "リソースタイトル", "説明", "", "book", "beginner", true},
		{"無効（カテゴリーが空）", "リソースタイトル", "説明", "https://example.com", "", "beginner", true},
		{"無効（カテゴリーが不正）", "リソースタイトル", "説明", "https://example.com", "invalid", "beginner", true},
		{"無効（難易度が不正）", "リソースタイトル", "説明", "https://example.com", "book", "expert", true},
		{"無効（説明が長すぎる）", "リソースタイトル", strings.Repeat("a", 1001), "https://example.com", "book", "beginner", true},
		{"無効（URLが不正）", "リソースタイトル", "説明", "not-a-url", "book", "beginner", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateCreateResource(tt.title, tt.description, tt.url, tt.category, tt.difficulty)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, domain.IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResourceValidator_ValidateUpdateResource(t *testing.T) {
	v := NewResourceValidator()

	// ValidateUpdateResourceはValidateCreateResourceと同じロジック
	err := v.ValidateUpdateResource("タイトル", "説明", "https://example.com", "book", "beginner")
	assert.NoError(t, err)

	err = v.ValidateUpdateResource("", "説明", "https://example.com", "book", "beginner")
	assert.Error(t, err)
}
