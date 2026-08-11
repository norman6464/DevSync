package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// AtCoderRatingFetcher は AtCoder のレーティング情報を外部から取得するための最小の契約。
// HTTP など取得手段の詳細は adapter 側に閉じ込める。
type AtCoderRatingFetcher interface {
	// FetchRatingHistory は指定ユーザーのレーティング履歴を古い順に返す。
	// ユーザーが存在しない場合や取得に失敗した場合は DomainError を返す。
	FetchRatingHistory(ctx context.Context, username string) ([]model.AtCoderRatingEntry, error)
	// UserExists は指定ユーザーのページを取得できるかどうかだけを返す。
	// 取得に失敗した場合は「存在しない」として false を返す。
	UserExists(ctx context.Context, username string) bool
}
