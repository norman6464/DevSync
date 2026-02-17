package dto_test

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreateCodeSnippetRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.CreateCodeSnippetRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト",
			request: dto.CreateCodeSnippetRequest{
				Language: "go",
				FileName: "main.go",
				Code:     "package main",
			},
			wantErr: false,
		},
		{
			name: "ファイル名なし（任意フィールド）",
			request: dto.CreateCodeSnippetRequest{
				Language: "go",
				Code:     "package main",
			},
			wantErr: false,
		},
		{
			name: "言語が空",
			request: dto.CreateCodeSnippetRequest{
				Language: "",
				Code:     "package main",
			},
			wantErr: true,
		},
		{
			name: "コードが空",
			request: dto.CreateCodeSnippetRequest{
				Language: "go",
				Code:     "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateCodeSnippetRequest_Validation(t *testing.T) {
	// UpdateCodeSnippetRequestは全フィールドが任意のためバリデーションエラーにならない
	tests := []struct {
		name    string
		request dto.UpdateCodeSnippetRequest
		wantErr bool
	}{
		{
			name: "全フィールド指定",
			request: dto.UpdateCodeSnippetRequest{
				Language: "python",
				FileName: "app.py",
				Code:     "print('hello')",
			},
			wantErr: false,
		},
		{
			name:    "全フィールド空（任意）",
			request: dto.UpdateCodeSnippetRequest{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateSnippetCommentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.CreateSnippetCommentRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト",
			request: dto.CreateSnippetCommentRequest{
				LineNumber: 10,
				Content:    "ここにバグがあります",
			},
			wantErr: false,
		},
		{
			name: "コンテンツが空",
			request: dto.CreateSnippetCommentRequest{
				LineNumber: 10,
				Content:    "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
