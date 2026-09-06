package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// PostTagHandler は投稿タグのHTTPハンドラー。
type PostTagHandler struct {
	setTags     *usecase.SetPostTagsUseCase
	getTags     *usecase.GetPostTagsUseCase
	findByTag   *usecase.FindPostsByTagUseCase
	popularTags *usecase.GetPopularTagsUseCase
}

// NewPostTagHandler は新しいPostTagHandlerを生成する。
func NewPostTagHandler(
	setTags *usecase.SetPostTagsUseCase,
	getTags *usecase.GetPostTagsUseCase,
	findByTag *usecase.FindPostsByTagUseCase,
	popularTags *usecase.GetPopularTagsUseCase,
) *PostTagHandler {
	return &PostTagHandler{setTags: setTags, getTags: getTags, findByTag: findByTag, popularTags: popularTags}
}

// setTagsRequest はタグ設定のリクエストボディ。
type setTagsRequest struct {
	Tags []string `json:"tags"`
}

// SetTags は投稿のタグを設定する。
func (h *PostTagHandler) SetTags(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	req := bindJSON[setTagsRequest](c)
	if req == nil {
		return
	}

	if err := h.setTags.Execute(c.Request.Context(), postID, c.GetUint("userID"), req.Tags); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("タグを更新しました"))
}

// tagsResponse はタグ一覧レスポンス
type tagsResponse struct {
	Tags interface{} `json:"tags"`
}

// GetByPostID は投稿のタグ一覧を取得する。
func (h *PostTagHandler) GetByPostID(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	tags, err := h.getTags.Execute(c.Request.Context(), postID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, tagsResponse{Tags: tags})
}

// FindPostsByTag はタグで投稿を検索する。
func (h *PostTagHandler) FindPostsByTag(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		respondBadRequest(c, "tagパラメータが必要です")
		return
	}
	if len([]rune(tag)) > maxTagQueryLen {
		respondBadRequest(c, "タグは100文字以下である必要があります")
		return
	}

	page, limit := parsePagination(c)
	offset := (page - 1) * limit

	posts, total, err := h.findByTag.Execute(c.Request.Context(), tag, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	respondPaginated(c, posts, total, page, limit)
}

// GetPopularTags は人気タグ一覧を取得する。
func (h *PostTagHandler) GetPopularTags(c *gin.Context) {
	tags, err := h.popularTags.Execute(c.Request.Context(), 20)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, tagsResponse{Tags: tags})
}
