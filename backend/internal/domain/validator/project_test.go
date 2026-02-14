package validator

import (
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestProjectValidator_ValidateCreateProject(t *testing.T) {
	v := NewProjectValidator()

	tests := []struct {
		name        string
		title       string
		description string
		demoURL     string
		githubURL   string
		wantErr     bool
	}{
		{"有効なプロジェクト", "プロジェクトタイトル", "プロジェクト説明", "https://demo.example.com", "https://github.com/user/repo", false},
		{"有効（デモURLなし）", "プロジェクトタイトル", "プロジェクト説明", "", "https://github.com/user/repo", false},
		{"有効（GitHubURLなし）", "プロジェクトタイトル", "プロジェクト説明", "https://demo.example.com", "", false},
		{"有効（URL両方なし）", "プロジェクトタイトル", "プロジェクト説明", "", "", false},
		{"無効（タイトルが空）", "", "プロジェクト説明", "", "", true},
		{"無効（説明が空）", "プロジェクトタイトル", "", "", "", true},
		{"無効（タイトルが長すぎる）", strings.Repeat("a", 201), "プロジェクト説明", "", "", true},
		{"無効（説明が長すぎる）", "プロジェクトタイトル", strings.Repeat("a", 10001), "", "", true},
		{"無効（デモURLが不正）", "プロジェクトタイトル", "プロジェクト説明", "not-a-url", "", true},
		{"無効（GitHubURLが不正）", "プロジェクトタイトル", "プロジェクト説明", "", "not-a-url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateCreateProject(tt.title, tt.description, tt.demoURL, tt.githubURL)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, domain.IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectValidator_ValidateUpdateProject(t *testing.T) {
	v := NewProjectValidator()

	// ValidateUpdateProjectはValidateCreateProjectと同じロジック
	err := v.ValidateUpdateProject("タイトル", "説明", "", "")
	assert.NoError(t, err)

	err = v.ValidateUpdateProject("", "説明", "", "")
	assert.Error(t, err)
}
