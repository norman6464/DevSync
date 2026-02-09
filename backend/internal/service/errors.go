// Package service はDevSyncアプリケーションのビジネスロジック層を提供する。
// 各Serviceはリポジトリインターフェースを通じてデータアクセスを行い、
// 権限チェック・ビジネスルール・通知連携などのロジックを実装する。
package service

import "errors"

// サービス層の共通エラー定義。
// ハンドラー層でこれらのエラーを判定し、適切なHTTPステータスコードに変換する。
var (
	ErrNotFound   = errors.New("not found")   // リソースが見つからない（404）
	ErrForbidden  = errors.New("forbidden")   // アクセス権限がない（403）
	ErrBadRequest = errors.New("bad request") // リクエストが不正（400）
	ErrConflict          = errors.New("conflict")            // リソースの競合（409）
	ErrRateLimitExceeded = errors.New("rate limit exceeded") // レート制限超過（429）
	ErrLLMNotConfigured  = errors.New("LLM not configured") // LLM未設定（503）
)
