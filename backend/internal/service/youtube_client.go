package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
)

// YouTubeClientInterface はYouTube Data API v3クライアントの契約を定義する。
// テスト時にモックに差し替え可能。
type YouTubeClientInterface interface {
	SearchVideos(query string, maxResults int, language string) ([]model.YouTubeVideo, error)
}

// YouTubeClient はYouTube Data API v3のクライアント実装。
type YouTubeClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewYouTubeClient は新しいYouTubeClientインスタンスを生成する。
func NewYouTubeClient(apiKey string) *YouTubeClient {
	return &YouTubeClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL: "https://www.googleapis.com/youtube/v3",
	}
}

// youtubeSearchResponse はYouTube Search APIのレスポンス構造体。
type youtubeSearchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelID    string `json:"channelId"`
			ChannelTitle string `json:"channelTitle"`
			PublishedAt  string `json:"publishedAt"`
			Thumbnails   struct {
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// SearchVideos はYouTube Search APIで動画を検索する。
func (c *YouTubeClient) SearchVideos(query string, maxResults int, language string) ([]model.YouTubeVideo, error) {
	if maxResults <= 0 || maxResults > 50 {
		maxResults = 10
	}
	if language == "" {
		language = "ja"
	}

	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("q", query)
	params.Set("type", "video")
	params.Set("maxResults", fmt.Sprintf("%d", maxResults))
	params.Set("relevanceLanguage", language)
	params.Set("order", "relevance")
	params.Set("key", c.apiKey)

	reqURL := fmt.Sprintf("%s/search?%s", c.baseURL, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "HTTPリクエストの作成に失敗", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "YouTube APIの呼び出しに失敗", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "レスポンスの読み取りに失敗", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable,
			fmt.Sprintf("YouTube APIエラー (ステータス %d): %s", resp.StatusCode, string(body)), nil)
	}

	var apiResp youtubeSearchResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "レスポンスのパースに失敗", err)
	}

	if apiResp.Error != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable,
			fmt.Sprintf("YouTube APIエラー: %s", apiResp.Error.Message), nil)
	}

	var videos []model.YouTubeVideo
	for _, item := range apiResp.Items {
		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		thumbnailURL := item.Snippet.Thumbnails.High.URL
		if thumbnailURL == "" {
			thumbnailURL = item.Snippet.Thumbnails.Medium.URL
		}
		videos = append(videos, model.YouTubeVideo{
			VideoID:      item.ID.VideoID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			ChannelID:    item.Snippet.ChannelID,
			ChannelTitle: item.Snippet.ChannelTitle,
			ThumbnailURL: thumbnailURL,
			PublishedAt:  publishedAt,
		})
	}

	return videos, nil
}
