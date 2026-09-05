package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostValidator_ValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "正常なタイトル",
			title:   "テスト投稿",
			wantErr: false,
		},
		{
			name:    "空のタイトル",
			title:   "",
			wantErr: true,
		},
		{
			name:    "タイトルが長すぎる（200文字超）",
			title:   string(make([]byte, 201)),
			wantErr: true,
		},
		{
			name:    "空白のみのタイトル",
			title:   "   ",
			wantErr: true,
		},
	}

	validator := NewPostValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTitle(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostValidator_ValidateContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "正常な本文",
			content: "これはテスト投稿の本文です。",
			wantErr: false,
		},
		{
			name:    "空の本文",
			content: "",
			wantErr: true,
		},
		{
			name:    "本文が長すぎる（10000文字超）",
			content: string(make([]byte, 10001)),
			wantErr: true,
		},
		{
			name:    "空白のみの本文",
			content: "   ",
			wantErr: true,
		},
	}

	validator := NewPostValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateContent(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostValidator_ValidatePost(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
		wantErr bool
	}{
		{
			name:    "正常な投稿",
			title:   "テストタイトル",
			content: "テスト本文",
			wantErr: false,
		},
		{
			name:    "タイトルが空",
			title:   "",
			content: "テスト本文",
			wantErr: true,
		},
		{
			name:    "本文が空",
			title:   "テストタイトル",
			content: "",
			wantErr: true,
		},
		{
			name:    "両方空",
			title:   "",
			content: "",
			wantErr: true,
		},
	}

	validator := NewPostValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePost(tt.title, tt.content)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostValidator_ValidateCreatePost(t *testing.T) {
	validator := NewPostValidator()

	tests := []struct {
		name      string
		title     string
		content   string
		imageURLs string
		tags      []string
		wantErr   bool
	}{
		{"有効な投稿", "タイトル", "本文", "", nil, false},
		{"有効（画像URL付き）", "タイトル", "本文", `["https://example.com/img.jpg"]`, nil, false},
		{"有効（タグ付き）", "タイトル", "本文", "", []string{"tag1", "tag2"}, false},
		{"無効（タイトルが空）", "", "本文", "", nil, true},
		{"無効（本文が空）", "タイトル", "", "", nil, true},
		{"無効（画像URL形式が不正）", "タイトル", "本文", "not-json", nil, true},
		{"無効（タグが多すぎる）", "タイトル", "本文", "", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreatePost(tt.title, tt.content, tt.imageURLs, tt.tags)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostValidator_ValidateUpdatePost(t *testing.T) {
	validator := NewPostValidator()

	tests := []struct {
		name      string
		title     string
		content   string
		imageURLs string
		wantErr   bool
	}{
		{"有効な更新", "タイトル", "本文", "", false},
		{"有効（画像URL付き）", "タイトル", "本文", `["https://example.com/img.jpg"]`, false},
		{"有効（タイトルが空）", "", "本文", "", false}, // 部分更新では空値OK
		{"有効（本文が空）", "タイトル", "", "", false}, // 部分更新では空値OK
		{"無効（不正な画像URL）", "タイトル", "本文", `["invalid"]`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUpdatePost(tt.title, tt.content, tt.imageURLs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostValidator_ValidateComment(t *testing.T) {
	validator := NewPostValidator()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"有効なコメント", "コメント内容", false},
		{"無効（空）", "", true},
		{"無効（長すぎる）", string(make([]byte, 1001)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateComment(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
