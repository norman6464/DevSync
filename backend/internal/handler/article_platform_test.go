package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectUsernameRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request connectUsernameRequest
		wantErr bool
	}{
		{
			name:    "有効なリクエスト",
			request: connectUsernameRequest{Username: "testuser"},
			wantErr: false,
		},
		{
			name:    "ユーザー名が空",
			request: connectUsernameRequest{Username: ""},
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
