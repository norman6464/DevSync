package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
)

// LLMClientInterface はLLM APIクライアントの契約を定義する。
// テスト時にモックに差し替え可能。
type LLMClientInterface interface {
	Complete(messages []ChatMessage) (*ChatResponse, error)
}

// ChatMessage はLLM APIに送信するメッセージを表す。
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse はLLM APIからのレスポンスを表す。
type ChatResponse struct {
	Content    string `json:"content"`
	TokensUsed int    `json:"tokens_used"`
}

// OpenAIClient はOpenAI Chat Completions APIのクライアント実装。
// net/httpのみを使用し、外部ライブラリに依存しない。
type OpenAIClient struct {
	apiKey     string
	httpClient *http.Client
	model      string
	baseURL    string
}

// NewOpenAIClient は新しいOpenAIClientインスタンスを生成する。
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		model:   "gpt-4o-mini",
		baseURL: "https://api.openai.com/v1/chat/completions",
	}
}

// openAIRequest はOpenAI APIリクエストボディを表す。
type openAIRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

// openAIResponse はOpenAI APIレスポンスを表す。
type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Complete はOpenAI Chat Completions APIを呼び出し、LLMの応答を返す。
func (c *OpenAIClient) Complete(messages []ChatMessage) (*ChatResponse, error) {
	reqBody := openAIRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "リクエストのシリアライズに失敗", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "HTTPリクエストの作成に失敗", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "OpenAI APIの呼び出しに失敗", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "レスポンスの読み取りに失敗", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[WARN] OpenAI APIエラー (ステータス %d): %s", resp.StatusCode, string(body))
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "AI サービスが一時的に利用できません", nil)
	}

	var apiResp openAIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "レスポンスのパースに失敗", err)
	}

	if apiResp.Error != nil {
		log.Printf("[WARN] OpenAI APIエラー: %s", apiResp.Error.Message)
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "AI サービスが一時的に利用できません", nil)
	}

	if len(apiResp.Choices) == 0 {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "OpenAI APIからの応答がありません", nil)
	}

	return &ChatResponse{
		Content:    apiResp.Choices[0].Message.Content,
		TokensUsed: apiResp.Usage.TotalTokens,
	}, nil
}
