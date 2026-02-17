package dto_test

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestConnectUsernameRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.ConnectUsernameRequest
		wantErr bool
	}{
		{
			name:    "有効なリクエスト",
			request: dto.ConnectUsernameRequest{Username: "testuser"},
			wantErr: false,
		},
		{
			name:    "ユーザー名が空",
			request: dto.ConnectUsernameRequest{Username: ""},
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
