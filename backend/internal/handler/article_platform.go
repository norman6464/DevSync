package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
)

// ArticlePlatformOps は記事プラットフォーム（Qiita/Zenn等）連携ハンドラーが呼び出す操作の集合。
// 移行済みのスライスは usecase のメソッド値をそのまま渡せる。
type ArticlePlatformOps[A any, S any] struct {
	Connect     func(ctx context.Context, userID uint, username string) (int, error)
	Disconnect  func(ctx context.Context, userID uint) error
	Sync        func(ctx context.Context, userID uint) (int, error)
	GetArticles func(ctx context.Context, userID uint) ([]A, error)
	GetStats    func(ctx context.Context, userID uint) (*S, error)
}

// ArticlePlatformHandler は記事プラットフォーム連携の共通HTTPハンドラ。
// プラットフォームアカウントの接続・切断・記事同期・統計情報の取得を処理する。
type ArticlePlatformHandler[A any, S any] struct {
	ops  ArticlePlatformOps[A, S]
	name string
}

// NewArticlePlatformHandler は新しいArticlePlatformHandlerインスタンスを生成する。
func NewArticlePlatformHandler[A any, S any](name string, ops ArticlePlatformOps[A, S]) *ArticlePlatformHandler[A, S] {
	return &ArticlePlatformHandler[A, S]{ops: ops, name: name}
}

// Connect はプラットフォームのユーザー名を設定し、記事を同期する。
func (h *ArticlePlatformHandler[A, S]) Connect(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.ConnectUsernameRequest](c)
	if input == nil {
		return
	}

	count, err := h.ops.Connect(c.Request.Context(), userID, input.Username)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ConnectSyncResponse{
		Message:       h.name + " connected successfully",
		ArticlesCount: count,
	})
}

// Disconnect はプラットフォームのユーザー名を削除し、キャッシュされた記事を削除する。
func (h *ArticlePlatformHandler[A, S]) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.ops.Disconnect(c.Request.Context(), userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse(h.name+" disconnected successfully"))
}

// Sync は現在のユーザーの記事を再同期する。
func (h *ArticlePlatformHandler[A, S]) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.ops.Sync(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ConnectSyncResponse{
		Message:       h.name + " synced successfully",
		ArticlesCount: count,
	})
}

// GetArticles は指定ユーザーの全記事を返す。
func (h *ArticlePlatformHandler[A, S]) GetArticles(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	articles, err := h.ops.GetArticles(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(articles))
}

// GetStats は指定ユーザーの統計情報を返す。
func (h *ArticlePlatformHandler[A, S]) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	stats, err := h.ops.GetStats(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
