package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
)

// validateSearchQuery は検索クエリをバリデーションする汎用ヘルパー関数。
// 空白をトリミングした上で空文字をチェックし、正規化されたクエリを返す。
func validateSearchQuery(query string) (string, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return "", domain.NewError(domain.ErrCodeBadRequest, "検索キーワードは必須です", nil)
	}
	return q, nil
}
