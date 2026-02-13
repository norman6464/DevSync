package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"有効なメール", "user@example.com", false},
		{"有効なメール（数字付き）", "user123@example.com", false},
		{"有効なメール（ドット付き）", "user.name@example.co.jp", false},
		{"無効なメール（@なし）", "userexample.com", true},
		{"無効なメール（ドメインなし）", "user@", true},
		{"無効なメール（空）", "", true},
		{"無効なメール（スペースのみ）", "   ", true},
		{"無効なメール（長すぎる）", strings.Repeat("a", 250) + "@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"有効なパスワード", "password123", false},
		{"有効なパスワード（記号付き）", "pass@word!", false},
		{"有効なパスワード（最小長）", "pass1!", false},
		{"無効なパスワード（短すぎる）", "pass", true},
		{"無効なパスワード（数字記号なし）", "password", true},
		{"無効なパスワード（空）", "", true},
		{"無効なパスワード（長すぎる）", strings.Repeat("a", 130), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"有効なユーザー名", "user_name", false},
		{"有効なユーザー名（ハイフン）", "user-name", false},
		{"有効なユーザー名（数字）", "user123", false},
		{"無効なユーザー名（短すぎる）", "u", true},
		{"無効なユーザー名（記号）", "user@name", true},
		{"無効なユーザー名（スペース）", "user name", true},
		{"無効なユーザー名（空）", "", true},
		{"無効なユーザー名（長すぎる）", strings.Repeat("a", 35), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"有効なタイトル", "My Title", false},
		{"有効なタイトル（長い）", strings.Repeat("a", 200), false},
		{"無効なタイトル（空）", "", true},
		{"無効なタイトル（スペースのみ）", "   ", true},
		{"無効なタイトル（長すぎる）", strings.Repeat("a", 201), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"有効なコンテンツ", "Some content", false},
		{"有効なコンテンツ（長い）", strings.Repeat("a", 10000), false},
		{"無効なコンテンツ（空）", "", true},
		{"無効なコンテンツ（スペースのみ）", "   ", true},
		{"無効なコンテンツ（長すぎる）", strings.Repeat("a", 10001), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContent(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"有効なURL（https）", "https://example.com", false},
		{"有効なURL（http）", "http://example.com/path", false},
		{"有効なURL（空・オプショナル）", "", false},
		{"無効なURL（httpなし）", "example.com", true},
		{"無効なURL（長すぎる）", "https://" + strings.Repeat("a", 2050), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTags(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		wantErr bool
	}{
		{"有効なタグ", []string{"tag1", "tag2"}, false},
		{"有効なタグ（空配列）", []string{}, false},
		{"無効なタグ（多すぎる）", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}, true},
		{"無効なタグ（空文字）", []string{"tag1", ""}, true},
		{"無効なタグ（長すぎる）", []string{strings.Repeat("a", 31)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTags(tt.tags)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
