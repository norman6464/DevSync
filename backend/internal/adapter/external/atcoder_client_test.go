package external

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// roundTripFunc はテスト用のHTTPラウンドトリッパー。
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// newTestAtCoderClient は任意のレスポンスを返すクライアントを生成する。
func newTestAtCoderClient(fn roundTripFunc) *atcoderClient {
	return &atcoderClient{client: &http.Client{Transport: fn}}
}

// jsonResponse はステータスと本文だけを持つレスポンスを返す。
func jsonResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body))}
}

func TestAtCoderClient_FetchRatingHistory_Success(t *testing.T) {
	var requestedURL string
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		requestedURL = req.URL.String()
		return jsonResponse(http.StatusOK, `[{"IsRated":true,"NewRating":1500,"ContestName":"ABC300"}]`), nil
	})

	history, err := c.FetchRatingHistory(context.Background(), "testuser")
	require.NoError(t, err)
	require.Len(t, history, 1)
	assert.Equal(t, 1500, history[0].NewRating)
	assert.Equal(t, "ABC300", history[0].ContestName)
	assert.True(t, history[0].IsRated)
	assert.Equal(t, "https://atcoder.jp/users/testuser/history/json", requestedURL)
}

func TestAtCoderClient_FetchRatingHistory_EmptyHistory(t *testing.T) {
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, "[]"), nil
	})

	history, err := c.FetchRatingHistory(context.Background(), "newuser")
	require.NoError(t, err)
	assert.Empty(t, history)
}

func TestAtCoderClient_FetchRatingHistory_NotFound(t *testing.T) {
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusNotFound, ""), nil
	})

	history, err := c.FetchRatingHistory(context.Background(), "unknown")
	assert.Nil(t, history)
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	assert.Contains(t, domainErr.Message, "unknown")
}

func TestAtCoderClient_FetchRatingHistory_ServerError(t *testing.T) {
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusInternalServerError, ""), nil
	})

	_, err := c.FetchRatingHistory(context.Background(), "testuser")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
	assert.Contains(t, domainErr.Message, "500")
}

func TestAtCoderClient_FetchRatingHistory_NetworkError(t *testing.T) {
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	_, err := c.FetchRatingHistory(context.Background(), "testuser")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestAtCoderClient_FetchRatingHistory_InvalidJSON(t *testing.T) {
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, "invalid json"), nil
	})

	_, err := c.FetchRatingHistory(context.Background(), "testuser")
	require.Error(t, err)

	var domainErr *domain.DomainError
	require.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

// ユーザー名はパスに埋め込む前にエスケープする。
func TestAtCoderClient_FetchRatingHistory_EscapesUsername(t *testing.T) {
	var requestedPath string
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		requestedPath = req.URL.EscapedPath()
		return jsonResponse(http.StatusOK, "[]"), nil
	})

	_, err := c.FetchRatingHistory(context.Background(), "a/../b")
	require.NoError(t, err)
	assert.Equal(t, "/users/a%2F..%2Fb/history/json", requestedPath)
}

// ctxKey はテストで ctx の伝播を確認するためのキー。
type ctxKey struct{}

// ctx は HTTP リクエストへ伝播する。
func TestAtCoderClient_FetchRatingHistory_PropagatesContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey{}, "propagated")

	var got interface{}
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		got = req.Context().Value(ctxKey{})
		return jsonResponse(http.StatusOK, "[]"), nil
	})

	_, err := c.FetchRatingHistory(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, "propagated", got)
}

func TestAtCoderClient_UserExists(t *testing.T) {
	tests := []struct {
		name   string
		status int
		err    error
		want   bool
	}{
		{"200 なら存在する", http.StatusOK, nil, true},
		{"404 なら存在しない", http.StatusNotFound, nil, false},
		{"500 なら存在しない扱い", http.StatusInternalServerError, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
				return jsonResponse(tt.status, ""), tt.err
			})
			assert.Equal(t, tt.want, c.UserExists(context.Background(), "testuser"))
		})
	}
}

func TestAtCoderClient_UserExists_NetworkError(t *testing.T) {
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	assert.False(t, c.UserExists(context.Background(), "testuser"))
}

// 本文が壊れていても取得できれば存在するものとして扱う（移行前と同じ）。
func TestAtCoderClient_UserExists_IgnoresBody(t *testing.T) {
	c := newTestAtCoderClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, "invalid json"), nil
	})

	assert.True(t, c.UserExists(context.Background(), "testuser"))
}
