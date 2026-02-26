package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
)

// ArticlePlatformServiceInterface は記事プラットフォーム（Qiita/Zenn等）サービスの共通インターフェース。
type ArticlePlatformServiceInterface[A any, S any] interface {
	Connect(userID uint, username string) (int, error)
	Disconnect(userID uint) error
	Sync(userID uint) (int, error)
	GetArticles(userID uint) ([]A, error)
	GetStats(userID uint) (*S, error)
}

// ArticlePlatformHandler は記事プラットフォーム連携の共通HTTPハンドラ。
// プラットフォームアカウントの接続・切断・記事同期・統計情報の取得を処理する。
type ArticlePlatformHandler[A any, S any] struct {
	service ArticlePlatformServiceInterface[A, S]
	name    string
}

// NewArticlePlatformHandler は新しいArticlePlatformHandlerインスタンスを生成する。
func NewArticlePlatformHandler[A any, S any](s ArticlePlatformServiceInterface[A, S], name string) *ArticlePlatformHandler[A, S] {
	return &ArticlePlatformHandler[A, S]{service: s, name: name}
}

// Connect はプラットフォームのユーザー名を設定し、記事を同期する。
func (h *ArticlePlatformHandler[A, S]) Connect(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.ConnectUsernameRequest](c)
	if input == nil {
		return
	}

	count, err := h.service.Connect(userID, input.Username)
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

	if err := h.service.Disconnect(userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse(h.name+" disconnected successfully"))
}

// Sync は現在のユーザーの記事を再同期する。
func (h *ArticlePlatformHandler[A, S]) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.Sync(userID)
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

	articles, err := h.service.GetArticles(userID)
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

	stats, err := h.service.GetStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}
