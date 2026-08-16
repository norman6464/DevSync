package external

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestSpotifyClient は fake のトークンエンドポイントを向く spotifyClient を返す。
func newTestSpotifyClient(tokenURL string) *spotifyClient {
	return &spotifyClient{
		clientID:     "id",
		clientSecret: "secret",
		redirectURL:  "http://localhost/callback",
		httpClient:   &http.Client{Timeout: time.Second},
		tokenURL:     tokenURL,
	}
}

func assertSpotifyDomainCode(t *testing.T, err error, want domain.ErrorCode) {
	t.Helper()
	var de *domain.DomainError
	require.ErrorAs(t, err, &de)
	assert.Equal(t, want, de.Code)
}

func TestSpotifyRequestToken(t *testing.T) {
	t.Run("200 でトークンが返れば成功する", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"at","refresh_token":"rt","expires_in":3600}`))
		}))
		defer srv.Close()

		token, err := newTestSpotifyClient(srv.URL).ExchangeCode(context.Background(), "code")

		require.NoError(t, err)
		assert.Equal(t, "at", token.AccessToken)
	})

	// 認可コードの期限切れなどクライアント側起因は 503 でなく 400 系にする。
	t.Run("4xx（invalid_grant）は BadRequest になる", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid_grant","error_description":"Authorization code expired"}`))
		}))
		defer srv.Close()

		_, err := newTestSpotifyClient(srv.URL).ExchangeCode(context.Background(), "expired")

		assertSpotifyDomainCode(t, err, domain.ErrCodeBadRequest)
	})

	t.Run("5xx は ServiceUnavailable になる", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		_, err := newTestSpotifyClient(srv.URL).RefreshAccessToken(context.Background(), "rt")

		assertSpotifyDomainCode(t, err, domain.ErrCodeServiceUnavailable)
	})

	t.Run("ネットワーク障害は ServiceUnavailable になる", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		srv.Close() // 接続不能にする

		_, err := newTestSpotifyClient(srv.URL).ExchangeCode(context.Background(), "code")

		assertSpotifyDomainCode(t, err, domain.ErrCodeServiceUnavailable)
	})

	t.Run("200 でもトークンが空なら ServiceUnavailable（従来どおり）", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))
		}))
		defer srv.Close()

		_, err := newTestSpotifyClient(srv.URL).ExchangeCode(context.Background(), "code")

		assertSpotifyDomainCode(t, err, domain.ErrCodeServiceUnavailable)
	})
}
