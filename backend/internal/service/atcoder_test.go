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

// roundTripFunc はテスト用のHTTPラウンドトリッパー。
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// newTestAtCoderService はテスト用のAtCoderServiceを生成する。
func newTestAtCoderService(fn roundTripFunc) *AtCoderService {
	return &AtCoderService{
		client: &http.Client{Transport: fn},
	}
}

func TestRatingToColor(t *testing.T) {
	tests := []struct {
		name   string
		rating int
		want   string
	}{
		{"赤（2800以上）", 2800, "red"},
		{"橙（2400-2799）", 2400, "orange"},
		{"黄（2000-2399）", 2000, "yellow"},
		{"青（1600-1999）", 1600, "blue"},
		{"水色（1200-1599）", 1200, "cyan"},
		{"緑（800-1199）", 800, "green"},
		{"茶（400-799）", 400, "brown"},
		{"灰（0-399）", 399, "gray"},
		{"灰（0）", 0, "gray"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ratingToColor(tt.rating))
		})
	}
}

func TestRatingToRank(t *testing.T) {
	tests := []struct {
		name   string
		rating int
		want   string
	}{
		{"赤（2800以上）", 3000, "赤"},
		{"橙（2400-2799）", 2500, "橙"},
		{"黄（2000-2399）", 2100, "黄"},
		{"青（1600-1999）", 1700, "青"},
		{"水色（1200-1599）", 1300, "水色"},
		{"緑（800-1199）", 900, "緑"},
		{"茶（400-799）", 500, "茶"},
		{"灰（0-399）", 100, "灰"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ratingToRank(tt.rating))
		})
	}
}

func TestGetRating_Success(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		body := `[{"IsRated":true,"NewRating":1500,"ContestName":"ABC300"}]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	info, err := svc.GetRating("testuser")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", info.Username)
	assert.Equal(t, 1500, info.Rating)
	assert.Equal(t, "cyan", info.Color)
	assert.Equal(t, "水色", info.Rank)
}

func TestGetRating_EmptyHistory(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("[]")),
		}, nil
	})

	info, err := svc.GetRating("newuser")
	assert.NoError(t, err)
	assert.Equal(t, 0, info.Rating)
	assert.Equal(t, "gray", info.Color)
	assert.Equal(t, "灰", info.Rank)
}

func TestGetRating_NotFound(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	info, err := svc.GetRating("unknown")
	assert.Nil(t, info)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
}

func TestGetRating_ServerError(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	info, err := svc.GetRating("testuser")
	assert.Nil(t, info)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestGetRating_NetworkError(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	})

	info, err := svc.GetRating("testuser")
	assert.Nil(t, info)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestGetRating_InvalidJSON(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("invalid json")),
		}, nil
	})

	info, err := svc.GetRating("testuser")
	assert.Nil(t, info)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestValidateUsername_Valid(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("[]")),
		}, nil
	})

	assert.True(t, svc.ValidateUsername("validuser"))
}

func TestValidateUsername_Invalid(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	assert.False(t, svc.ValidateUsername("invaliduser"))
}

func TestValidateUsername_NetworkError(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("timeout")
	})

	assert.False(t, svc.ValidateUsername("testuser"))
}

func TestGetRating_RequestURL(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "atcoder.jp", req.URL.Host)
		assert.Equal(t, "/users/myuser/history/json", req.URL.Path)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("[]")),
		}, nil
	})

	_, err := svc.GetRating("myuser")
	assert.NoError(t, err)
}

func TestGetRating_MultipleEntries(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		body := `[
			{"IsRated":true,"NewRating":400,"ContestName":"ABC100"},
			{"IsRated":true,"NewRating":800,"ContestName":"ABC200"},
			{"IsRated":true,"NewRating":1200,"ContestName":"ABC300"}
		]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	info, err := svc.GetRating("testuser")
	assert.NoError(t, err)
	assert.Equal(t, 1200, info.Rating)
	assert.Equal(t, "cyan", info.Color)
	assert.Equal(t, "水色", info.Rank)
}

func TestGetRating_HighRating(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		body := `[{"IsRated":true,"NewRating":3500,"ContestName":"ABC999"}]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	info, err := svc.GetRating("tourist")
	assert.NoError(t, err)
	assert.Equal(t, 3500, info.Rating)
	assert.Equal(t, "red", info.Color)
	assert.Equal(t, "赤", info.Rank)
}

func TestGetRating_RateLimited(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	info, err := svc.GetRating("testuser")
	assert.Nil(t, info)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeServiceUnavailable, domainErr.Code)
}

func TestValidateUsername_ServerError(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	assert.False(t, svc.ValidateUsername("testuser"))
}

func TestValidateUsername_RequestURL(t *testing.T) {
	svc := newTestAtCoderService(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "atcoder.jp", req.URL.Host)
		assert.Equal(t, "/users/atcoderuser/history/json", req.URL.Path)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("[]")),
		}, nil
	})

	assert.True(t, svc.ValidateUsername("atcoderuser"))
}
