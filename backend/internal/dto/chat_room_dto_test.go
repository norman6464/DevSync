package dto_test

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreateChatRoomRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.CreateChatRoomRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト",
			request: dto.CreateChatRoomRequest{
				Name:        "テストルーム",
				Description: "説明文",
				MemberIDs:   []uint{1, 2},
			},
			wantErr: false,
		},
		{
			name: "名前のみ（最小限）",
			request: dto.CreateChatRoomRequest{
				Name: "ルーム",
			},
			wantErr: false,
		},
		{
			name: "名前が空",
			request: dto.CreateChatRoomRequest{
				Name: "",
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

func TestAddChatRoomMemberRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.AddChatRoomMemberRequest
		wantErr bool
	}{
		{
			name:    "有効なリクエスト",
			request: dto.AddChatRoomMemberRequest{UserID: 1},
			wantErr: false,
		},
		{
			name:    "UserIDがゼロ",
			request: dto.AddChatRoomMemberRequest{UserID: 0},
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

func TestSendChatRoomMessageRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.SendChatRoomMessageRequest
		wantErr bool
	}{
		{
			name:    "有効なリクエスト",
			request: dto.SendChatRoomMessageRequest{Content: "こんにちは"},
			wantErr: false,
		},
		{
			name:    "コンテンツが空",
			request: dto.SendChatRoomMessageRequest{Content: ""},
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
