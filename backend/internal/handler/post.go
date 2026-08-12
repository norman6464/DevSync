package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
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
	Publish(id, userID uint) (*model.Post, error)
	Unpublish(id, userID uint) (*model.Post, error)
	SchedulePublish(id, userID uint, scheduledAt time.Time) (*model.Post, error)
	CancelSchedule(id, userID uint) (*model.Post, error)
	GetScheduled(userID uint) ([]model.Post, error)
	AutoSaveDraft(userID, draftID uint, title, content, imageURLs string) (*model.Post, error)
	CountByUserID(userID uint) (int64, error)
	CountDraftsByUserID(userID uint) (int64, error)
	CountScheduledByUserID(userID uint) (int64, error)
}

// PostHandler は投稿関連のHTTPハンドラ。
// 投稿のCRUD・いいね・コメント・タイムラインを処理する。
type PostHandler struct {
	service       PostServiceInterface
	createSnippet *usecase.CreateCodeSnippetUseCase
	autoTags      *usecase.SetAutoPostTagsUseCase

	// リアクション・コメント・ブックマークは DIP へ移行済み。他の操作は順次 usecase へ移していく。
	addReaction    *usecase.AddPostReactionUseCase
	removeReaction *usecase.RemovePostReactionUseCase
	getReactions   *usecase.GetPostReactionsUseCase
	reactionsBatch *usecase.GetPostReactionsBatchUseCase

	bookmark       *usecase.BookmarkPostUseCase
	unbookmark     *usecase.UnbookmarkPostUseCase
	hasBookmarked  *usecase.HasBookmarkedPostUseCase
	listBookmarks  *usecase.ListBookmarkedPostsUseCase
	countBookmarks *usecase.CountBookmarkedPostsUseCase

	createComment *usecase.CreatePostCommentUseCase
	listComments  *usecase.ListPostCommentsUseCase
	listReplies   *usecase.ListCommentRepliesUseCase
	editComment   *usecase.EditPostCommentUseCase
	deleteComment *usecase.DeletePostCommentUseCase
	hideComment   *usecase.HidePostCommentUseCase
	unhideComment *usecase.UnhidePostCommentUseCase
}

// NewPostHandler は新しいPostHandlerインスタンスを生成する。
func NewPostHandler(
	s PostServiceInterface,
	createSnippet *usecase.CreateCodeSnippetUseCase,
	addReaction *usecase.AddPostReactionUseCase,
	removeReaction *usecase.RemovePostReactionUseCase,
	getReactions *usecase.GetPostReactionsUseCase,
	reactionsBatch *usecase.GetPostReactionsBatchUseCase,
	createComment *usecase.CreatePostCommentUseCase,
	listComments *usecase.ListPostCommentsUseCase,
	listReplies *usecase.ListCommentRepliesUseCase,
	editComment *usecase.EditPostCommentUseCase,
	deleteComment *usecase.DeletePostCommentUseCase,
	hideComment *usecase.HidePostCommentUseCase,
	unhideComment *usecase.UnhidePostCommentUseCase,
	bookmark *usecase.BookmarkPostUseCase,
	unbookmark *usecase.UnbookmarkPostUseCase,
	hasBookmarked *usecase.HasBookmarkedPostUseCase,
	listBookmarks *usecase.ListBookmarkedPostsUseCase,
	countBookmarks *usecase.CountBookmarkedPostsUseCase,
) *PostHandler {
	return &PostHandler{
		service:        s,
		createSnippet:  createSnippet,
		addReaction:    addReaction,
		removeReaction: removeReaction,
		getReactions:   getReactions,
		reactionsBatch: reactionsBatch,
		createComment:  createComment,
		listComments:   listComments,
		listReplies:    listReplies,
		editComment:    editComment,
		deleteComment:  deleteComment,
		hideComment:    hideComment,
		unhideComment:  unhideComment,
		bookmark:       bookmark,
		unbookmark:     unbookmark,
		hasBookmarked:  hasBookmarked,
		listBookmarks:  listBookmarks,
		countBookmarks: countBookmarks,
	}
}

// SetAutoTagsUseCase はオプショナルな自動タグ設定を注入する。
// 設定すると、投稿の作成・更新時にコンテンツからハッシュタグを自動抽出してタグを設定する。
func (h *PostHandler) SetAutoTagsUseCase(autoTags *usecase.SetAutoPostTagsUseCase) {
	h.autoTags = autoTags
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
	if len(input.CodeSnippets) > 0 && h.createSnippet != nil {
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
			h.createSnippet.Execute(c.Request.Context(), snippet)
		}
		// スニペット付きで再取得
		if updated, err := h.service.GetByID(created.ID); err == nil {
			created = updated
		}
	}

	// コンテンツからハッシュタグを自動抽出してタグを設定
	if h.autoTags != nil {
		_ = h.autoTags.Execute(c.Request.Context(), created.ID, userID, input.Content)
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
		Bookmarked: h.isBookmarked(c, userID, post.ID),
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

	// コンテンツ更新時にハッシュタグを自動再抽出
	if h.autoTags != nil && input.Content != "" {
		_ = h.autoTags.Execute(c.Request.Context(), id, userID, input.Content)
	}

	respondOK(c, post)
}

// Delete は投稿を削除する。所有者のみ削除可能。
func (h *PostHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
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
	respondOK(c, ensureSlice(posts))
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
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetMyPosts は認証ユーザー自身の投稿一覧を取得する。
func (h *PostHandler) GetMyPosts(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	posts, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.PostListResponse{
		Posts:  ensureSlice(posts),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetMyCount は認証ユーザーの投稿数を返す。
func (h *PostHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.CountByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// GetDraftsCount は認証ユーザーの下書き投稿数を返す。
func (h *PostHandler) GetDraftsCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.CountDraftsByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// GetScheduledCount は認証ユーザーのスケジュール済み投稿数を返す。
func (h *PostHandler) GetScheduledCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.CountScheduledByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// Like は投稿にいいねする。
func (h *PostHandler) Like(c *gin.Context) {
	handleToggleAction(c, h.service.Like, "liked")
}

// Unlike は投稿のいいねを取り消す。
func (h *PostHandler) Unlike(c *gin.Context) {
	handleToggleAction(c, h.service.Unlike, "unliked")
}

// GetComments は投稿のコメント一覧を返す。
func (h *PostHandler) GetComments(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	comments, err := h.listComments.Execute(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(comments))
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

	comment, err := h.createComment.Execute(c.Request.Context(), userID, id, input.Content, input.ParentID)
	if err != nil {
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

	replies, err := h.listReplies.Execute(c.Request.Context(), commentID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(replies))
}

// EditComment はコメントを編集する。所有者のみ編集可能。
func (h *PostHandler) EditComment(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdateCommentRequest](c)
	if input == nil {
		return
	}

	comment, err := h.editComment.Execute(c.Request.Context(), commentID, userID, input.Content)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, comment)
}

// DeleteComment はコメントを削除する。所有者のみ削除可能。
func (h *PostHandler) DeleteComment(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.deleteComment.Execute(c.Request.Context(), commentID, userID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// HideComment はコメントを非表示にする。
func (h *PostHandler) HideComment(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.hideComment.Execute(c.Request.Context(), commentID, userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("コメントを非表示にしました"))
}

// UnhideComment はコメントの非表示を解除する。
func (h *PostHandler) UnhideComment(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.unhideComment.Execute(c.Request.Context(), commentID, userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("コメントの非表示を解除しました"))
}

// GetDrafts は現在のユーザーの下書き投稿一覧を返す。
func (h *PostHandler) GetDrafts(c *gin.Context) {
	userID := c.GetUint("userID")

	drafts, err := h.service.GetDrafts(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(drafts))
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

// isBookmarked は投稿詳細に載せるブックマーク済みフラグを返す。
// 取得に失敗しても投稿詳細自体は返したいため、移行前と同じく false にフォールバックする。
func (h *PostHandler) isBookmarked(c *gin.Context, userID, postID uint) bool {
	bookmarked, err := h.hasBookmarked.Execute(c.Request.Context(), userID, postID)
	if err != nil {
		return false
	}
	return bookmarked
}

// Bookmark は投稿をブックマークする。
func (h *PostHandler) Bookmark(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.bookmark.Execute(c.Request.Context(), userID, id)
	}, "bookmarked")
}

// Unbookmark は投稿のブックマークを解除する。
func (h *PostHandler) Unbookmark(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.unbookmark.Execute(c.Request.Context(), userID, id)
	}, "unbookmarked")
}

// GetBookmarks は現在のユーザーのブックマーク済み投稿一覧を返す。
func (h *PostHandler) GetBookmarks(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	posts, total, err := h.listBookmarks.Execute(c.Request.Context(), userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, dto.PostListResponse{
		Posts:  ensureSlice(posts),
		Total:  total,
		Limit:  limit,
		Offset: (page - 1) * limit,
	})
}

// GetBookmarksCount は認証ユーザーのブックマーク済み投稿数を返す。
func (h *PostHandler) GetBookmarksCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.countBookmarks.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
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

	if err := h.addReaction.Execute(c.Request.Context(), userID, id, input.Emoji); err != nil {
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

	if err := h.removeReaction.Execute(c.Request.Context(), userID, id, input.Emoji); err != nil {
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

	reactions, userReactions, err := h.getReactions.Execute(c.Request.Context(), userID, id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ReactionResponse{
		Reactions:     reactions,
		UserReactions: userReactions,
	})
}

// SchedulePublish は投稿のスケジュール公開日時を設定する。
func (h *PostHandler) SchedulePublish(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.SchedulePublishRequest](c)
	if input == nil {
		return
	}

	scheduledAt, ok := parseDateTimeRFC3339(input.ScheduledAt)
	if !ok {
		respondBadRequest(c, "日時のフォーマットが不正です（RFC3339形式で指定してください）")
		return
	}

	post, err := h.service.SchedulePublish(id, userID, scheduledAt)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, post)
}

// CancelSchedule はスケジュールを解除して下書きに戻す。
func (h *PostHandler) CancelSchedule(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	post, err := h.service.CancelSchedule(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, post)
}

// GetScheduled はユーザーのスケジュール済み投稿一覧を返す。
func (h *PostHandler) GetScheduled(c *gin.Context) {
	userID := c.GetUint("userID")

	posts, err := h.service.GetScheduled(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(posts))
}

// AutoSaveDraft は下書きをサーバーサイドで自動保存する。
func (h *PostHandler) AutoSaveDraft(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.AutoSaveDraftRequest](c)
	if req == nil {
		return
	}

	result, err := h.service.AutoSaveDraft(userID, req.ID, req.Title, req.Content, req.ImageURLs)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.AutoSaveDraftResponse{
		ID:        result.ID,
		UpdatedAt: result.UpdatedAt.Format(time.RFC3339),
	})
}

// GetReactionsBatch は複数投稿のリアクション情報を一括取得する。
func (h *PostHandler) GetReactionsBatch(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.BatchReactionRequest](c)
	if input == nil {
		return
	}

	reactions, userReactions, err := h.reactionsBatch.Execute(c.Request.Context(), userID, input.PostIDs)
	if err != nil {
		respondError(c, err)
		return
	}

	result := make(map[uint]dto.ReactionResponse, len(input.PostIDs))
	for _, id := range input.PostIDs {
		result[id] = dto.ReactionResponse{
			Reactions:     reactions[id],
			UserReactions: userReactions[id],
		}
	}

	respondOK(c, dto.BatchReactionResponse{Reactions: result})
}
