// Package service はDevSyncアプリケーションのビジネスロジック層を提供する。
// 各Serviceはリポジトリインターフェースを通じてデータアクセスを行い、
// 権限チェック・ビジネスルール・通知連携などのロジックを実装する。
package service

import "github.com/norman6464/devsync/backend/internal/domain"

// サービス層の共通エラー定義。
// domain.DomainErrorを使用して統一的なエラーハンドリングを実現する。
var (
	ErrNotFound          = domain.ErrNotFound          // リソースが見つからない（404）
	ErrForbidden         = domain.ErrForbidden         // アクセス権限がない（403）
	ErrBadRequest        = domain.ErrBadRequest        // リクエストが不正（400）
	ErrUnauthorized      = domain.ErrUnauthorized      // 認証が必要（401）
	ErrConflict          = domain.ErrConflict          // リソースの競合（409）
	ErrRateLimitExceeded = domain.ErrRateLimitExceeded // レート制限超過（429）
	ErrLLMNotConfigured  = domain.ErrServiceUnavailable // LLM未設定（503）
)

// validateRequiredID はIDが0でないことを検証する。
// 0の場合はBadRequestエラーを返す。
func validateRequiredID(id uint, name string) error {
	if id == 0 {
		return domain.NewError(domain.ErrCodeBadRequest, name+"は必須です", nil)
	}
	return nil
}
