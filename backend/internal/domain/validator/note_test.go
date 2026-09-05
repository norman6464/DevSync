package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoteValidator_ValidateTitle(t *testing.T) {
	v := NewNoteValidator()

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"有効（通常のタイトル）", "学習ノート", false},
		{"有効（長いタイトル）", "TypeScriptの型システムについての学習ノート", false},
		{"無効（空文字）", "", true},
		{"無効（スペースのみ）", "   ", true},
		{"無効（長すぎる）", string(make([]byte, 201)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateTitle(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteValidator_ValidateContent(t *testing.T) {
	v := NewNoteValidator()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"有効（通常のマークダウン）", "# 見出し\n\n本文です", false},
		{"有効（空文字）", "", false},                           // 本文は空でもOK
		{"無効（長すぎる）", string(make([]byte, 100001)), true}, // 100KB超
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateContent(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteValidator_ValidateTags(t *testing.T) {
	v := NewNoteValidator()

	tests := []struct {
		name    string
		tags    string
		wantErr bool
	}{
		{"有効（通常のタグ）", "Go,TypeScript,React", false},
		{"有効（空文字）", "", false}, // タグは空でもOK
		{"有効（スペース含む）", "Go, TypeScript, React", false},
		{"無効（長すぎる）", string(make([]byte, 501)), true}, // 500文字超
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateTags(tt.tags)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteValidator_ValidateCreateNote(t *testing.T) {
	v := NewNoteValidator()

	tests := []struct {
		name    string
		title   string
		content string
		tags    string
		wantErr bool
	}{
		{"有効（完全なノート）", "学習ノート", "# 見出し\n本文", "Go,TDD", false},
		{"有効（タグなし）", "学習ノート", "本文", "", false},
		{"無効（タイトル空）", "", "本文", "", true},
		{"無効（タイトル長すぎ）", string(make([]byte, 201)), "本文", "", true},
		{"無効（本文長すぎ）", "タイトル", string(make([]byte, 100001)), "", true},
		{"無効（タグ長すぎ）", "タイトル", "本文", string(make([]byte, 501)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateCreateNote(tt.title, tt.content, tt.tags)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteValidator_ValidateUpdateNote(t *testing.T) {
	v := NewNoteValidator()

	tests := []struct {
		name    string
		title   string
		content string
		tags    string
		wantErr bool
	}{
		{"有効（完全な更新）", "更新タイトル", "更新内容", "Go", false},
		{"有効（タイトルのみ）", "更新タイトル", "", "", false}, // 部分更新OK
		{"有効（本文のみ）", "", "更新内容", "", false},     // 部分更新OK
		{"有効（タグのみ）", "", "", "React", false},    // 部分更新OK
		{"無効（タイトル長すぎ）", string(make([]byte, 201)), "", "", true},
		{"無効（本文長すぎ）", "", string(make([]byte, 100001)), "", true},
		{"無効（タグ長すぎ）", "", "", string(make([]byte, 501)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateUpdateNote(tt.title, tt.content, tt.tags)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
