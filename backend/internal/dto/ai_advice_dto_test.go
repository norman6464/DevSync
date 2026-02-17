package dto_test

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestAIChatRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.AIChatRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト（新規会話）",
			request: dto.AIChatRequest{
				Message: "Go言語について教えてください",
			},
			wantErr: false,
		},
		{
			name: "有効なリクエスト（既存会話）",
			request: dto.AIChatRequest{
				Message:        "続きを教えてください",
				ConversationID: 42,
			},
			wantErr: false,
		},
		{
			name: "メッセージが空",
			request: dto.AIChatRequest{
				Message: "",
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
