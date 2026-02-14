package validator

import (
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestQuestionValidator_ValidateCreateQuestion(t *testing.T) {
	v := NewQuestionValidator()

	tests := []struct {
		name    string
		title   string
		body    string
		tags    string
		wantErr bool
	}{
		{"有効な質問", "質問タイトル", "質問本文", "tag1,tag2", false},
		{"有効な質問（タグなし）", "質問タイトル", "質問本文", "", false},
		{"無効（タイトルが空）", "", "質問本文", "", true},
		{"無効（本文が空）", "質問タイトル", "", "", true},
		{"無効（タイトルが長すぎる）", strings.Repeat("a", 501), "質問本文", "", true},
		{"無効（本文が長すぎる）", "質問タイトル", strings.Repeat("a", 10001), "", true},
		{"無効（タグが長すぎる）", "質問タイトル", "質問本文", strings.Repeat("a", 301), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateCreateQuestion(tt.title, tt.body, tt.tags)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, domain.IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQuestionValidator_ValidateUpdateQuestion(t *testing.T) {
	v := NewQuestionValidator()

	// ValidateUpdateQuestionはValidateCreateQuestionと同じロジック
	err := v.ValidateUpdateQuestion("タイトル", "本文", "")
	assert.NoError(t, err)

	err = v.ValidateUpdateQuestion("", "本文", "")
	assert.Error(t, err)
}

func TestQuestionValidator_ValidateVote(t *testing.T) {
	v := NewQuestionValidator()

	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"有効（1）", 1, false},
		{"有効（-1）", -1, false},
		{"無効（0）", 0, true},
		{"無効（2）", 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateVote(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, domain.IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
