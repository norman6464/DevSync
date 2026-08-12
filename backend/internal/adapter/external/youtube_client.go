package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// youtubeRequestTimeout は YouTube への 1 リクエストのタイムアウト。
const youtubeRequestTimeout = 15 * time.Second

// youtubeAPIBaseURL は YouTube Data API v3 のベース URL。
const youtubeAPIBaseURL = "https://www.googleapis.com/youtube/v3"

// 検索件数の既定値と上限。
const (
	youtubeDefaultMaxResults = 10
	youtubeMaxAllowedResults = 50
)

// youtubeSearchResponse は YouTube Search API のレスポンス構造。
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

// youTubeClient は [repository.YouTubeVideoSearcher] の HTTP 実装。
type youTubeClient struct {
	apiKey  string
	client  *http.Client
	baseURL string
}

// NewYouTubeClient は YouTubeVideoSearcher の HTTP 実装を返す。
func NewYouTubeClient(apiKey string) repository.YouTubeVideoSearcher {
	return &youTubeClient{
		apiKey:  apiKey,
		client:  &http.Client{Timeout: youtubeRequestTimeout},
		baseURL: youtubeAPIBaseURL,
	}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.YouTubeVideoSearcher = (*youTubeClient)(nil)

// SearchVideos は YouTube Search API で動画を検索する。
func (c *youTubeClient) SearchVideos(ctx context.Context, query string, maxResults int, language string) ([]model.YouTubeVideo, error) {
	if maxResults <= 0 || maxResults > youtubeMaxAllowedResults {
		maxResults = youtubeDefaultMaxResults
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/search?"+params.Encode(), nil)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "HTTPリクエストの作成に失敗", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "YouTube APIの呼び出しに失敗", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "レスポンスの読み取りに失敗", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[WARN] YouTube APIエラー (ステータス %d): %s", resp.StatusCode, string(body))
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "YouTube APIが一時的に利用できません", nil)
	}

	var apiResp youtubeSearchResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "レスポンスのパースに失敗", err)
	}

	if apiResp.Error != nil {
		log.Printf("[WARN] YouTube APIエラー: %s", apiResp.Error.Message)
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "YouTube APIが一時的に利用できません", nil)
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
