package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoteTemplateValidator_ValidateName(t *testing.T) {
	v := NewNoteTemplateValidator()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"有効（通常）", "週報テンプレート", false},
		{"有効（最大長）", strings.Repeat("a", 100), false},
		{"無効（空文字）", "", true},
		{"無効（スペースのみ）", "   ", true},
		{"無効（最大長超過）", strings.Repeat("a", 101), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteTemplateValidator_ValidateContentTemplate(t *testing.T) {
	v := NewNoteTemplateValidator()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"有効（通常）", "# タイトル\n\n本文", false},
		{"有効（長文）", strings.Repeat("a", 10000), false},
		{"無効（空文字）", "", true},
		{"無効（スペースのみ）", "   \n\n  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateContentTemplate(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteTemplateValidator_ValidateDescription(t *testing.T) {
	v := NewNoteTemplateValidator()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"有効（通常）", "週報用のテンプレート", false},
		{"有効（空文字）", "", false},
		{"有効（最大長）", strings.Repeat("a", 500), false},
		{"無効（最大長超過）", strings.Repeat("a", 501), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateDescription(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteTemplateValidator_ValidateDefaultTitle(t *testing.T) {
	v := NewNoteTemplateValidator()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"有効（通常）", "週報 - {{date}}", false},
		{"有効（空文字）", "", false},
		{"有効（最大長）", strings.Repeat("a", 200), false},
		{"無効（最大長超過）", strings.Repeat("a", 201), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateDefaultTitle(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNoteTemplateValidator_ValidateCreateTemplate(t *testing.T) {
	v := NewNoteTemplateValidator()

	tests := []struct {
		name            string
		inputName       string
		contentTemplate string
		wantErr         bool
	}{
		{"有効", "週報", "# 週報\n\n## 今週やったこと", false},
		{"無効（名前が空）", "", "本文", true},
		{"無効（本文が空）", "週報", "", true},
		{"無効（両方空）", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateCreateTemplate(tt.inputName, tt.contentTemplate)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
