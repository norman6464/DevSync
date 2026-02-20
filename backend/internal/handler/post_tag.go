package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// PostTagServiceInterface はPostTagHandlerが依存するサービスインターフェース。
type PostTagServiceInterface interface {
	SetTags(postID, userID uint, tags []string) error
	GetByPostID(postID uint) ([]string, error)
	FindPostsByTag(tag string, limit, offset int) ([]model.Post, int64, error)
	GetPopularTags(limit int) ([]model.TagCount, error)
}

// PostTagHandler は投稿タグのHTTPハンドラー。
type PostTagHandler struct {
	service PostTagServiceInterface
}

// NewPostTagHandler は新しいPostTagHandlerを生成する。
func NewPostTagHandler(service PostTagServiceInterface) *PostTagHandler {
	return &PostTagHandler{service: service}
}

// SetTags は投稿のタグを設定する。
func (h *PostTagHandler) SetTags(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	req := bindJSON[struct {
		Tags []string `json:"tags"`
	}](c)
	if req == nil {
		return
	}

	if err := h.service.SetTags(postID, userID, req.Tags); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("タグを更新しました"))
}

// GetByPostID は投稿のタグ一覧を取得する。
func (h *PostTagHandler) GetByPostID(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	tags, err := h.service.GetByPostID(postID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, dto.TagsResponse{Tags: tags})
}

// FindPostsByTag はタグで投稿を検索する。
func (h *PostTagHandler) FindPostsByTag(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		respondBadRequest(c, "tagパラメータが必要です")
		return
	}

	page, limit := parsePagination(c)
	offset := (page - 1) * limit

	posts, total, err := h.service.FindPostsByTag(tag, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	respondPaginated(c, posts, total, page, limit)
}

// GetPopularTags は人気タグ一覧を取得する。
func (h *PostTagHandler) GetPopularTags(c *gin.Context) {
	tags, err := h.service.GetPopularTags(20)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, dto.TagsResponse{Tags: tags})
}
