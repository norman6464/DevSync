package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// MentionServiceInterface はMentionHandlerが依存するサービスインターフェース。
type MentionServiceInterface interface {
	ProcessMentions(actorID uint, text string, postID *uint, commentID *uint) error
	GetMentionsByUserID(userID uint, page, limit int) ([]model.Mention, error)
	GetMentionsByPostID(postID uint) ([]model.Mention, error)
	DeleteMentionsByPostID(postID uint) error
	DeleteMentionsByCommentID(commentID uint) error
}

// MentionHandler はメンションのHTTPハンドラー。
type MentionHandler struct {
	service MentionServiceInterface
}

// NewMentionHandler は新しいMentionHandlerを生成する。
func NewMentionHandler(service MentionServiceInterface) *MentionHandler {
	return &MentionHandler{service: service}
}

// GetMyMentions は認証ユーザーへのメンション一覧を取得する。
func (h *MentionHandler) GetMyMentions(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	mentions, err := h.service.GetMentionsByUserID(userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(mentions))
}

// GetPostMentions は投稿に関連するメンション一覧を取得する。
func (h *MentionHandler) GetPostMentions(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	mentions, err := h.service.GetMentionsByPostID(postID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(mentions))
}
