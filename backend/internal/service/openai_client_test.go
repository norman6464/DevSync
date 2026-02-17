package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

// newTestOpenAIClient はテスト用のOpenAIClientを生成する。
func newTestOpenAIClient(fn roundTripFunc) *OpenAIClient {
	return &OpenAIClient{
		apiKey:     "test-key",
		httpClient: &http.Client{Transport: fn},
		model:      "gpt-4o-mini",
		baseURL:    "https://api.openai.com/v1/chat/completions",
	}
}

func TestComplete_Success(t *testing.T) {
	client := newTestOpenAIClient(func(req *http.Request) (*http.Response, error) {
		// リクエストヘッダーを検証
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-key", req.Header.Get("Authorization"))

		body := `{"choices":[{"message":{"content":"こんにちは！"}}],"usage":{"total_tokens":25}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	messages := []ChatMessage{{Role: "user", Content: "挨拶して"}}
	resp, err := client.Complete(messages)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "こんにちは！", resp.Content)
	assert.Equal(t, 25, resp.TokensUsed)
}

func TestComplete_NetworkError(t *testing.T) {
	client := newTestOpenAIClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	messages := []ChatMessage{{Role: "user", Content: "test"}}
	resp, err := client.Complete(messages)
	assert.Nil(t, resp)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestComplete_ServerError(t *testing.T) {
	client := newTestOpenAIClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("Internal Server Error")),
		}, nil
	})

	messages := []ChatMessage{{Role: "user", Content: "test"}}
	resp, err := client.Complete(messages)
	assert.Nil(t, resp)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
	assert.Contains(t, domainErr.Message, "500")
}

func TestComplete_RateLimited(t *testing.T) {
	client := newTestOpenAIClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"rate limit exceeded"}}`)),
		}, nil
	})

	messages := []ChatMessage{{Role: "user", Content: "test"}}
	resp, err := client.Complete(messages)
	assert.Nil(t, resp)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
	assert.Contains(t, domainErr.Message, "429")
}

func TestComplete_InvalidJSON(t *testing.T) {
	client := newTestOpenAIClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("not json")),
		}, nil
	})

	messages := []ChatMessage{{Role: "user", Content: "test"}}
	resp, err := client.Complete(messages)
	assert.Nil(t, resp)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestComplete_APIError(t *testing.T) {
	client := newTestOpenAIClient(func(req *http.Request) (*http.Response, error) {
		body := `{"error":{"message":"invalid api key"},"choices":[]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	messages := []ChatMessage{{Role: "user", Content: "test"}}
	resp, err := client.Complete(messages)
	assert.Nil(t, resp)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
	assert.Contains(t, domainErr.Message, "invalid api key")
}

func TestComplete_EmptyChoices(t *testing.T) {
	client := newTestOpenAIClient(func(req *http.Request) (*http.Response, error) {
		body := `{"choices":[],"usage":{"total_tokens":0}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	messages := []ChatMessage{{Role: "user", Content: "test"}}
	resp, err := client.Complete(messages)
	assert.Nil(t, resp)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestNewOpenAIClient(t *testing.T) {
	client := NewOpenAIClient("sk-test-key")
	assert.NotNil(t, client)
	assert.Equal(t, "sk-test-key", client.apiKey)
	assert.Equal(t, "gpt-4o-mini", client.model)
	assert.Equal(t, "https://api.openai.com/v1/chat/completions", client.baseURL)
	assert.NotNil(t, client.httpClient)
}
