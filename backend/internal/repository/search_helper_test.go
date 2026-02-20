package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeLikeChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"通常の文字列", "Go言語", "Go言語"},
		{"パーセント記号をエスケープ", "100%完了", "100\\%完了"},
		{"アンダースコアをエスケープ", "user_name", "user\\_name"},
		{"バックスラッシュをエスケープ", "path\\to", "path\\\\to"},
		{"空文字列", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeLikeChars(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEscapeLikePattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"通常の検索クエリ", "Go言語", "%Go言語%"},
		{"パーセント記号をエスケープ", "100%完了", "%100\\%完了%"},
		{"アンダースコアをエスケープ", "user_name", "%user\\_name%"},
		{"バックスラッシュをエスケープ", "path\\to", "%path\\\\to%"},
		{"複数のワイルドカード", "%test_100%", "%\\%test\\_100\\%%"},
		{"空文字列", "", "%%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeLikePattern(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
