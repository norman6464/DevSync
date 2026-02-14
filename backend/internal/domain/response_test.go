package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSuccessResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	response := NewSuccessResponse(data)

	assert.Equal(t, data, response.Data)
}

func TestNewErrorResponse(t *testing.T) {
	message := "エラーメッセージ"
	code := "VALIDATION_ERROR"
	details := map[string]string{"field": "title"}

	response := NewErrorResponse(message, code, details)

	assert.Equal(t, message, response.Error)
	assert.Equal(t, code, response.Code)
	assert.Equal(t, details, response.Details)
}

func TestNewPaginatedResponse(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		total      int64
		page       int
		limit      int
		wantPages  int
	}{
		{
			name:      "10件中10件を1ページで取得",
			data:      []string{"a", "b"},
			total:     10,
			page:      1,
			limit:     10,
			wantPages: 1,
		},
		{
			name:      "100件を20件ずつ",
			data:      []string{},
			total:     100,
			page:      1,
			limit:     20,
			wantPages: 5,
		},
		{
			name:      "95件を20件ずつ（端数あり）",
			data:      []string{},
			total:     95,
			page:      1,
			limit:     20,
			wantPages: 5,
		},
		{
			name:      "0件",
			data:      []string{},
			total:     0,
			page:      1,
			limit:     10,
			wantPages: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := NewPaginatedResponse(tt.data, tt.total, tt.page, tt.limit)

			assert.Equal(t, tt.data, response.Data)
			assert.Equal(t, tt.total, response.Total)
			assert.Equal(t, tt.page, response.Page)
			assert.Equal(t, tt.limit, response.Limit)
			assert.Equal(t, tt.wantPages, response.TotalPages)
		})
	}
}

func TestNewMessageResponse(t *testing.T) {
	message := "成功しました"
	response := NewMessageResponse(message)

	assert.Equal(t, message, response.Message)
}
