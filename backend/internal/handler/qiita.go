package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// QiitaServiceInterface はQiitaサービスの抽象インターフェース。
type QiitaServiceInterface interface {
	Connect(userID uint, username string) (int, error)
	Disconnect(userID uint) error
	Sync(userID uint) (int, error)
	GetArticles(userID uint) ([]model.QiitaArticle, error)
	GetStats(userID uint) (*model.QiitaStats, error)
}

// QiitaHandler はQiita連携関連のHTTPハンドラ。
// Qiitaアカウントの接続・切断・記事同期・統計情報の取得を処理する。
type QiitaHandler struct {
	service QiitaServiceInterface
}

// NewQiitaHandler は新しいQiitaHandlerインスタンスを生成する。
func NewQiitaHandler(s QiitaServiceInterface) *QiitaHandler {
	return &QiitaHandler{service: s}
}

// Connect はQiitaユーザー名を設定し、記事を同期する。
func (h *QiitaHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")

	var req dto.ConnectUsernameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "username is required")
		return
	}

	count, err := h.service.Connect(userID, req.Username)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ConnectSyncResponse{
		Message:       "Qiita connected successfully",
		ArticlesCount: count,
	})
}

// Disconnect はQiitaユーザー名を削除し、キャッシュされた記事を削除する。
func (h *QiitaHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.service.Disconnect(userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("Qiita disconnected successfully"))
}

// Sync は現在のユーザーのQiita記事を再同期する。
func (h *QiitaHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.Sync(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ConnectSyncResponse{
		Message:       "Qiita synced successfully",
		ArticlesCount: count,
	})
}

// GetArticles は指定ユーザーの全Qiita記事を返す。
func (h *QiitaHandler) GetArticles(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	articles, err := h.service.GetArticles(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, articles)
}

// GetStats は指定ユーザーのQiita統計情報を返す。
func (h *QiitaHandler) GetStats(c *gin.Context) {
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
