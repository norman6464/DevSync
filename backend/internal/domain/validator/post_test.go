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
		errMsg  string
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
			errMsg:  "タイトルは必須です",
		},
		{
			name:    "タイトルが長すぎる（200文字超）",
			title:   string(make([]byte, 201)),
			wantErr: true,
			errMsg:  "タイトルは200文字以内で入力してください",
		},
		{
			name:    "空白のみのタイトル",
			title:   "   ",
			wantErr: true,
			errMsg:  "タイトルは必須です",
		},
	}

	validator := NewPostValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTitle(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
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
		errMsg  string
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
			errMsg:  "本文は必須です",
		},
		{
			name:    "本文が長すぎる（10000文字超）",
			content: string(make([]byte, 10001)),
			wantErr: true,
			errMsg:  "本文は10000文字以内で入力してください",
		},
		{
			name:    "空白のみの本文",
			content: "   ",
			wantErr: true,
			errMsg:  "本文は必須です",
		},
	}

	validator := NewPostValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateContent(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
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
