package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// MentionHandler はメンションのHTTPハンドラー。
type MentionHandler struct {
	listMine *usecase.ListUserMentionsUseCase
	listPost *usecase.ListPostMentionsUseCase
}

// NewMentionHandler は新しいMentionHandlerを生成する。
func NewMentionHandler(
	listMine *usecase.ListUserMentionsUseCase,
	listPost *usecase.ListPostMentionsUseCase,
) *MentionHandler {
	return &MentionHandler{listMine: listMine, listPost: listPost}
}

// GetMyMentions は認証ユーザーへのメンション一覧を取得する。
func (h *MentionHandler) GetMyMentions(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	mentions, err := h.listMine.Execute(c.Request.Context(), userID, page, limit)
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

	mentions, err := h.listPost.Execute(c.Request.Context(), postID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(mentions))
}
