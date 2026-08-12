package external

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// openAIRequestTimeout は OpenAI への 1 リクエストのタイムアウト。
const openAIRequestTimeout = 30 * time.Second

const (
	openAIModel              = "gpt-4o-mini"
	openAIChatCompletionsURL = "https://api.openai.com/v1/chat/completions"
	// openAIMaxTokens は 1 応答あたりの最大トークン数。
	openAIMaxTokens = 1024
	// openAITemperature は応答のばらつき。
	openAITemperature = 0.7
)

// openAIClient は [repository.LLMClient] の OpenAI Chat Completions API 実装。
// net/http のみを使用し、外部ライブラリに依存しない。
type openAIClient struct {
	apiKey     string
	httpClient *http.Client
	model      string
	baseURL    string
}

// NewOpenAIClient は LLMClient の OpenAI 実装を返す。
func NewOpenAIClient(apiKey string) repository.LLMClient {
	return &openAIClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: openAIRequestTimeout},
		model:      openAIModel,
		baseURL:    openAIChatCompletionsURL,
	}
}

var _ repository.LLMClient = (*openAIClient)(nil)

// openAIRequest はOpenAI APIリクエストボディを表す。
type openAIRequest struct {
	Model       string              `json:"model"`
	Messages    []model.ChatMessage `json:"messages"`
	MaxTokens   int                 `json:"max_tokens"`
	Temperature float64             `json:"temperature"`
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

// Complete は OpenAI Chat Completions API を呼び出し、LLM の応答を返す。
func (c *openAIClient) Complete(ctx context.Context, messages []model.ChatMessage) (*model.ChatResponse, error) {
	reqBody := openAIRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   openAIMaxTokens,
		Temperature: openAITemperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "リクエストのシリアライズに失敗", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(jsonData))
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

	return &model.ChatResponse{
		Content:    apiResp.Choices[0].Message.Content,
		TokensUsed: apiResp.Usage.TotalTokens,
	}, nil
}
