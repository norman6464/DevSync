package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// PostServiceInterface はPostServiceが実装すべきインターフェース。
type PostServiceInterface interface {
	Create(post *model.Post) (*model.Post, error)
	GetByID(id uint) (*model.Post, error)
	GetAll(page, limit int) ([]model.Post, error)
	GetByUserID(userID uint) ([]model.Post, error)
	GetDrafts(userID uint) ([]model.Post, error)
	Timeline(userID uint, page, limit int) ([]model.Post, error)
	Update(id, userID uint, title, content, imageUrls string) (*model.Post, error)
	Delete(id, userID uint) error
	Like(userID, postID uint) error
	Unlike(userID, postID uint) error
	HasLiked(userID, postID uint) bool
	CreateComment(comment *model.Comment) error
	GetComments(postID uint) ([]model.Comment, error)
	DeleteComment(id, userID uint) error
	Publish(id, userID uint) (*model.Post, error)
	Bookmark(userID, postID uint) error
	Unbookmark(userID, postID uint) error
	HasBookmarked(userID, postID uint) bool
	GetBookmarks(userID uint, page, limit int) ([]model.Post, int64, error)
}

// CodeSnippetServiceInterface はCodeSnippetServiceが実装すべきインターフェース。
type CodeSnippetServiceInterface interface {
	Create(snippet *model.CodeSnippet) (*model.CodeSnippet, error)
}

// PostHandler は投稿関連のHTTPハンドラ。
// 投稿のCRUD・いいね・コメント・タイムラインを処理する。
type PostHandler struct {
	service        PostServiceInterface
	snippetService CodeSnippetServiceInterface
}

// NewPostHandler は新しいPostHandlerインスタンスを生成する。
func NewPostHandler(s PostServiceInterface, snippetService CodeSnippetServiceInterface) *PostHandler {
	return &PostHandler{service: s, snippetService: snippetService}
}

// Create は新しい投稿を作成する。
func (h *PostHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	input := bindJSON[dto.CreatePostRequest](c)
	if input == nil {
		return
	}

	post := &model.Post{
		UserID:    userID,
		Title:     input.Title,
		Content:   input.Content,
		ImageURLs: input.ImageURLs,
		IsDraft:   input.IsDraft,
	}
	created, err := h.service.Create(post)
	if err != nil {
		respondError(c, err)
		return
	}

	// コードスニペットを一括作成
	if len(input.CodeSnippets) > 0 && h.snippetService != nil {
		for _, s := range input.CodeSnippets {
			if s.Language == "" || s.Code == "" {
				continue
			}
			snippet := &model.CodeSnippet{
				PostID:   created.ID,
				UserID:   userID,
				Language: s.Language,
				FileName: s.FileName,
				Code:     s.Code,
			}
			h.snippetService.Create(snippet)
		}
		// スニペット付きで再取得
		if updated, err := h.service.GetByID(created.ID); err == nil {
			created = updated
		}
	}

	respondCreated(c, created)
}

// GetAll は投稿一覧をページネーション付きで返す。
func (h *PostHandler) GetAll(c *gin.Context) {
	page, limit := parsePagination(c)

	posts, err := h.service.GetAll(page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, posts)
}

// GetByID は指定IDの投稿を返す。いいね済みフラグも付与する。
func (h *PostHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	post, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	userID := c.GetUint("userID")
	respondOK(c, dto.PostDetailResponse{
		Post:       *post,
		Liked:      h.service.HasLiked(userID, post.ID),
		Bookmarked: h.service.HasBookmarked(userID, post.ID),
	})
}

// Update は投稿を更新する。所有者のみ更新可能。
func (h *PostHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdatePostRequest](c)
	if input == nil {
		return
	}

	post, err := h.service.Update(id, userID, input.Title, input.Content, input.ImageURLs)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, post)
}

// Delete は投稿を削除する。所有者のみ削除可能。
func (h *PostHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Delete(id, userID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// Timeline はフォロー中ユーザーの投稿タイムラインを返す。
func (h *PostHandler) Timeline(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	posts, err := h.service.Timeline(userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, posts)
}

// GetUserPosts は指定ユーザーの投稿一覧を返す。
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	posts, err := h.service.GetByUserID(id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, posts)
}

// Like は投稿にいいねする。
func (h *PostHandler) Like(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Like(userID, id); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("liked"))
}

// Unlike は投稿のいいねを取り消す。
func (h *PostHandler) Unlike(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Unlike(userID, id); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("unliked"))
}

// GetComments は投稿のコメント一覧を返す。
func (h *PostHandler) GetComments(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	comments, err := h.service.GetComments(id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, comments)
}

// CreateComment は投稿にコメントを作成する。
func (h *PostHandler) CreateComment(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.CreateCommentRequest](c)
	if input == nil {
		return
	}

	comment := &model.Comment{UserID: userID, PostID: id, Content: input.Content}
	if err := h.service.CreateComment(comment); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, comment)
}

// DeleteComment はコメントを削除する。所有者のみ削除可能。
func (h *PostHandler) DeleteComment(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.DeleteComment(commentID, userID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// GetDrafts は現在のユーザーの下書き投稿一覧を返す。
func (h *PostHandler) GetDrafts(c *gin.Context) {
	userID := c.GetUint("userID")

	drafts, err := h.service.GetDrafts(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, drafts)
}

// Publish は下書き投稿を公開する。所有者のみ公開可能。
func (h *PostHandler) Publish(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	post, err := h.service.Publish(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, post)
}

// Bookmark は投稿をブックマークする。
func (h *PostHandler) Bookmark(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Bookmark(userID, id); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("bookmarked"))
}

// Unbookmark は投稿のブックマークを解除する。
func (h *PostHandler) Unbookmark(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Unbookmark(userID, id); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("unbookmarked"))
}

// GetBookmarks は現在のユーザーのブックマーク済み投稿一覧を返す。
func (h *PostHandler) GetBookmarks(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	posts, total, err := h.service.GetBookmarks(userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, dto.BookmarkedPostsResponse{Posts: posts, Total: total})
}
