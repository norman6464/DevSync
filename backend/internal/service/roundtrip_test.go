package service

import "net/http"

// roundTripFunc はテスト用のHTTPラウンドトリッパー。
// 外部APIを叩くサービスのテストで、任意のレスポンスを返すために共用する。
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
