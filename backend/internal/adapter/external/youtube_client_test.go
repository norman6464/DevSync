package external

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestYouTubeClient は任意のレスポンスを返すクライアントを生成する。
func newTestYouTubeClient(fn roundTripFunc) *youTubeClient {
	return &youTubeClient{
		apiKey:  "test-key",
		client:  &http.Client{Transport: fn},
		baseURL: youtubeAPIBaseURL,
	}
}

func TestYouTubeClient_SearchVideos_Success(t *testing.T) {
	var query url.Values
	c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
		query = req.URL.Query()
		body := `{"items":[{"id":{"videoId":"abc123"},"snippet":{"title":"Go入門","description":"説明","channelId":"ch1","channelTitle":"チャンネル","publishedAt":"2026-01-01T00:00:00Z","thumbnails":{"medium":{"url":"https://img/medium.jpg"},"high":{"url":"https://img/high.jpg"}}}}]}`
		return jsonResponse(http.StatusOK, body), nil
	})

	videos, err := c.SearchVideos(context.Background(), "golang", 10, "ja")
	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "abc123", videos[0].VideoID)
	assert.Equal(t, "Go入門", videos[0].Title)
	assert.Equal(t, "ch1", videos[0].ChannelID)
	// 高解像度のサムネイルを優先する。
	assert.Equal(t, "https://img/high.jpg", videos[0].ThumbnailURL)
	assert.Equal(t, 2026, videos[0].PublishedAt.Year())

	assert.Equal(t, "snippet", query.Get("part"))
	assert.Equal(t, "golang", query.Get("q"))
	assert.Equal(t, "video", query.Get("type"))
	assert.Equal(t, "10", query.Get("maxResults"))
	assert.Equal(t, "ja", query.Get("relevanceLanguage"))
	assert.Equal(t, "relevance", query.Get("order"))
	assert.Equal(t, "test-key", query.Get("key"))
}

// 高解像度が無ければ中解像度のサムネイルを使う。
func TestYouTubeClient_SearchVideos_FallbackThumbnail(t *testing.T) {
	c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
		body := `{"items":[{"id":{"videoId":"abc"},"snippet":{"thumbnails":{"medium":{"url":"https://img/medium.jpg"},"high":{"url":""}}}}]}`
		return jsonResponse(http.StatusOK, body), nil
	})

	videos, err := c.SearchVideos(context.Background(), "golang", 10, "ja")
	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "https://img/medium.jpg", videos[0].ThumbnailURL)
}

// 件数と言語の既定値を補う。
func TestYouTubeClient_SearchVideos_Defaults(t *testing.T) {
	tests := []struct {
		name           string
		maxResults     int
		language       string
		wantMaxResults string
		wantLanguage   string
	}{
		{"0 件指定は 10 件に丸める", 0, "ja", "10", "ja"},
		{"上限超えは 10 件に丸める", 51, "ja", "10", "ja"},
		{"言語未指定は ja", 5, "", "5", "ja"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var query url.Values
			c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
				query = req.URL.Query()
				return jsonResponse(http.StatusOK, `{"items":[]}`), nil
			})

			_, err := c.SearchVideos(context.Background(), "golang", tt.maxResults, tt.language)
			require.NoError(t, err)
			assert.Equal(t, tt.wantMaxResults, query.Get("maxResults"))
			assert.Equal(t, tt.wantLanguage, query.Get("relevanceLanguage"))
		})
	}
}

func TestYouTubeClient_SearchVideos_Empty(t *testing.T) {
	c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"items":[]}`), nil
	})

	videos, err := c.SearchVideos(context.Background(), "golang", 10, "ja")
	require.NoError(t, err)
	assert.Empty(t, videos)
}

func TestYouTubeClient_SearchVideos_HTTPError(t *testing.T) {
	c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusForbidden, `{"error":{"code":403,"message":"quotaExceeded"}}`), nil
	})

	_, err := c.SearchVideos(context.Background(), "golang", 10, "ja")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
	// API キーやクォータの詳細はクライアントへ返さない。
	assert.Equal(t, "YouTube APIが一時的に利用できません", domainErr.Message)
}

// 200 でも本文にエラーが入っている場合は利用不可として扱う。
func TestYouTubeClient_SearchVideos_ErrorInBody(t *testing.T) {
	c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"items":[],"error":{"code":400,"message":"badRequest"}}`), nil
	})

	_, err := c.SearchVideos(context.Background(), "golang", 10, "ja")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestYouTubeClient_SearchVideos_NetworkError(t *testing.T) {
	c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	_, err := c.SearchVideos(context.Background(), "golang", 10, "ja")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestYouTubeClient_SearchVideos_InvalidJSON(t *testing.T) {
	c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, "invalid json"), nil
	})

	_, err := c.SearchVideos(context.Background(), "golang", 10, "ja")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

// 検索キーワードはクエリとしてエスケープされる。
func TestYouTubeClient_SearchVideos_EscapesQuery(t *testing.T) {
	var rawQuery, q string
	c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
		rawQuery = req.URL.RawQuery
		q = req.URL.Query().Get("q")
		return jsonResponse(http.StatusOK, `{"items":[]}`), nil
	})

	_, err := c.SearchVideos(context.Background(), "go&key=leaked", 10, "ja")
	require.NoError(t, err)
	assert.Contains(t, rawQuery, "q=go%26key%3Dleaked")
	assert.Equal(t, "go&key=leaked", q)
}

func TestYouTubeClient_SearchVideos_PropagatesContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey{}, "propagated")

	var got interface{}
	c := newTestYouTubeClient(func(req *http.Request) (*http.Response, error) {
		got = req.Context().Value(ctxKey{})
		return jsonResponse(http.StatusOK, `{"items":[]}`), nil
	})

	_, err := c.SearchVideos(ctx, "golang", 10, "ja")
	require.NoError(t, err)
	assert.Equal(t, "propagated", got)
}
