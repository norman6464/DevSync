package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoteFolderValidator_ValidateName(t *testing.T) {
	v := NewNoteFolderValidator()

	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"有効（通常のフォルダ名）", "マイフォルダ", false},
		{"有効（英数字）", "My Folder 2024", false},
		{"有効（最大長）", strings.Repeat("a", 100), false},
		{"無効（空文字）", "", true},
		{"無効（スペースのみ）", "   ", true},
		{"無効（長すぎる）", strings.Repeat("a", 101), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateName(tt.input)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteFolderValidator_ValidateCreate(t *testing.T) {
	v := NewNoteFolderValidator()

	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"有効なフォルダ作成", "新規フォルダ", false},
		{"無効（空文字）", "", true},
		{"無効（長すぎる）", strings.Repeat("a", 101), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateCreate(tt.input)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteFolderValidator_ValidateUpdate(t *testing.T) {
	v := NewNoteFolderValidator()

	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"有効な更新", "更新後の名前", false},
		{"有効（空文字 - 部分更新）", "", false}, // 更新時は空を許容
		{"無効（長すぎる）", strings.Repeat("a", 101), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateUpdate(tt.input)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
