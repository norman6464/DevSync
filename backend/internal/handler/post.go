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
	CountAll() (int64, error)
	GetByUserID(userID uint, limit, offset int) ([]model.Post, int64, error)
	GetDrafts(userID uint) ([]model.Post, error)
	Timeline(userID uint, page, limit int) ([]model.Post, error)
	Update(id, userID uint, title, content, imageUrls string) (*model.Post, error)
	Delete(id, userID uint) error
	Like(userID, postID uint) error
	Unlike(userID, postID uint) error
	HasLiked(userID, postID uint) bool
	CreateComment(comment *model.Comment) error
	GetComments(postID uint) ([]model.Comment, error)
	GetReplies(parentID uint) ([]model.Comment, error)
	DeleteComment(id, userID uint) error
	Publish(id, userID uint) (*model.Post, error)
	Unpublish(id, userID uint) (*model.Post, error)
	Bookmark(userID, postID uint) error
	Unbookmark(userID, postID uint) error
	HasBookmarked(userID, postID uint) bool
	GetBookmarks(userID uint, page, limit int) ([]model.Post, int64, error)
	AddReaction(userID, postID uint, emoji string) error
	RemoveReaction(userID, postID uint, emoji string) error
	GetReactionsByPostID(postID uint) ([]model.ReactionCount, error)
	GetUserReactions(userID, postID uint) ([]string, error)
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

	total, err := h.service.CountAll()
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, posts, total, page, limit)
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

// GetUserPosts は指定ユーザーの投稿一覧をページネーション付きで返す。
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	posts, total, err := h.service.GetByUserID(id, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.PostListResponse{
		Posts:  ensureSlice(posts),
		Total: total,
		Limit: limit,
		Offset: offset,
	})
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

	comment := &model.Comment{UserID: userID, PostID: id, Content: input.Content, ParentID: input.ParentID}
	if err := h.service.CreateComment(comment); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, comment)
}

// GetReplies は指定コメントへの返信一覧を返す。
func (h *PostHandler) GetReplies(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}

	replies, err := h.service.GetReplies(commentID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, replies)
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

// Unpublish は公開済みの投稿を下書きに戻す。所有者のみ操作可能。
func (h *PostHandler) Unpublish(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	post, err := h.service.Unpublish(id, userID)
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

// AddReaction は投稿にリアクション（絵文字）を追加する。
func (h *PostHandler) AddReaction(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.ReactionRequest](c)
	if input == nil {
		return
	}

	if err := h.service.AddReaction(userID, id, input.Emoji); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("reaction_added"))
}

// RemoveReaction は投稿のリアクションを削除する。
func (h *PostHandler) RemoveReaction(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.ReactionRequest](c)
	if input == nil {
		return
	}

	if err := h.service.RemoveReaction(userID, id, input.Emoji); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("reaction_removed"))
}

// GetReactions は投稿のリアクション一覧を返す。
func (h *PostHandler) GetReactions(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	reactions, err := h.service.GetReactionsByPostID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	userReactions, err := h.service.GetUserReactions(userID, id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ReactionResponse{
		Reactions:     reactions,
		UserReactions: userReactions,
	})
}
