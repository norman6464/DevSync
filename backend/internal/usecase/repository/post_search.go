package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostSearchRepository は投稿の高度な検索に対する、usecase 側が要求する契約。
type PostSearchRepository interface {
	// SearchWithFilter はキーワード・タグ・日付範囲・ソート順で投稿を検索し、該当件数も返す。
	// 下書きは対象外。タグは AND 条件（すべて付与されている投稿のみ）で絞り込む。
	// params の Limit / SortBy は呼び出し側で正規化済みであることを前提とする。
	SearchWithFilter(ctx context.Context, params model.PostSearchParams) ([]model.Post, int64, error)
}
