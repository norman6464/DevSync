// Package external は外部サービス（HTTP API）へのアクセスを提供する adapter。
// usecase 側で宣言された port を実装し、HTTP の詳細をこの層に閉じ込める。
package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// atcoderRequestTimeout は AtCoder への 1 リクエストのタイムアウト。
const atcoderRequestTimeout = 10 * time.Second

// atcoderRatingHistoryURL はレーティング履歴を返す AtCoder のエンドポイント。
const atcoderRatingHistoryURL = "https://atcoder.jp/users/%s/history/json"

// atcoderClient は [repository.AtCoderRatingFetcher] の HTTP 実装。
type atcoderClient struct {
	client *http.Client
}

// NewAtCoderClient は AtCoderRatingFetcher の HTTP 実装を返す。
func NewAtCoderClient() repository.AtCoderRatingFetcher {
	return &atcoderClient{client: &http.Client{Timeout: atcoderRequestTimeout}}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.AtCoderRatingFetcher = (*atcoderClient)(nil)

// FetchRatingHistory は指定ユーザーのレーティング履歴を取得する。
func (c *atcoderClient) FetchRatingHistory(ctx context.Context, username string) ([]model.AtCoderRatingEntry, error) {
	resp, err := c.get(ctx, username)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "AtCoder APIリクエストに失敗", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.NewError(domain.ErrCodeNotFound, fmt.Sprintf("AtCoderユーザーが見つかりません: %s", username), nil)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, fmt.Sprintf("AtCoder APIエラー: ステータスコード %d", resp.StatusCode), nil)
	}

	var history []model.AtCoderRatingEntry
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, domain.NewError(domain.ErrCodeServiceUnavailable, "AtCoder APIレスポンスのパースに失敗", err)
	}
	return history, nil
}

// UserExists は指定ユーザーのページを取得できるかどうかを返す。
// レスポンスの中身は見ないため、本文が壊れていても取得できれば存在するものとして扱う。
func (c *atcoderClient) UserExists(ctx context.Context, username string) bool {
	resp, err := c.get(ctx, username)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// get はレーティング履歴のエンドポイントへ GET する。
// ユーザー名は usecase 側でも検証しているが、パスに埋め込む前にエスケープしておく。
func (c *atcoderClient) get(ctx context.Context, username string) (*http.Response, error) {
	endpoint := fmt.Sprintf(atcoderRatingHistoryURL, url.PathEscape(username))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}
