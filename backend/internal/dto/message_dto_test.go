package dto_test

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestSendMessageRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.SendMessageRequest
		wantErr bool
	}{
		{
			name:    "有効なリクエスト",
			request: dto.SendMessageRequest{Content: "こんにちは"},
			wantErr: false,
		},
		{
			name:    "コンテンツが空",
			request: dto.SendMessageRequest{Content: ""},
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
