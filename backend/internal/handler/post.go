package handler

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// PostUseCases は PostHandler が依存する投稿系 usecase をまとめる。
// 投稿スライスは操作数が多いため、コンストラクタの引数を 1 つにまとめて対応関係を明示する。
type PostUseCases struct {
	Create         *usecase.CreatePostUseCase
	Get            *usecase.GetPostUseCase
	List           *usecase.ListPostsUseCase
	Count          *usecase.CountPostsUseCase
	ListByUser     *usecase.ListUserPostsUseCase
	ListDrafts     *usecase.ListDraftPostsUseCase
	ListScheduled  *usecase.ListScheduledPostsUseCase
	Timeline       *usecase.GetTimelineUseCase
	Update         *usecase.UpdatePostUseCase
	Delete         *usecase.DeletePostUseCase
	Publish        *usecase.PublishPostUseCase
	Unpublish      *usecase.UnpublishPostUseCase
	Schedule       *usecase.SchedulePostPublishUseCase
	CancelSchedule *usecase.CancelPostScheduleUseCase
	AutoSaveDraft  *usecase.AutoSaveDraftUseCase
	CountByUser    *usecase.CountUserPostsUseCase
	CountDrafts    *usecase.CountUserDraftsUseCase
	CountScheduled *usecase.CountUserScheduledPostsUseCase

	Like     *usecase.LikePostUseCase
	Unlike   *usecase.UnlikePostUseCase
	HasLiked *usecase.HasLikedPostUseCase

	CreateSnippet *usecase.CreateCodeSnippetUseCase

	AddReaction    *usecase.AddPostReactionUseCase
	RemoveReaction *usecase.RemovePostReactionUseCase
	GetReactions   *usecase.GetPostReactionsUseCase
	ReactionsBatch *usecase.GetPostReactionsBatchUseCase

	ProcessMentions       *usecase.ProcessMentionsUseCase
	DeleteCommentMentions *usecase.DeleteCommentMentionsUseCase

	CreateComment *usecase.CreatePostCommentUseCase
	ListComments  *usecase.ListPostCommentsUseCase
	ListReplies   *usecase.ListCommentRepliesUseCase
	EditComment   *usecase.EditPostCommentUseCase
	DeleteComment *usecase.DeletePostCommentUseCase
	HideComment   *usecase.HidePostCommentUseCase
	UnhideComment *usecase.UnhidePostCommentUseCase

	Bookmark       *usecase.BookmarkPostUseCase
	Unbookmark     *usecase.UnbookmarkPostUseCase
	HasBookmarked  *usecase.HasBookmarkedPostUseCase
	ListBookmarks  *usecase.ListBookmarkedPostsUseCase
	CountBookmarks *usecase.CountBookmarkedPostsUseCase
}

// PostHandler は投稿関連のHTTPハンドラ。
// 投稿のCRUD・いいね・コメント・リアクション・ブックマーク・タイムラインを処理する。
type PostHandler struct {
	uc       PostUseCases
	autoTags *usecase.SetAutoPostTagsUseCase
}

// NewPostHandler は新しいPostHandlerインスタンスを生成する。
func NewPostHandler(uc PostUseCases) *PostHandler {
	return &PostHandler{uc: uc}
}

// SetAutoTagsUseCase はオプショナルな自動タグ設定を注入する。
// 設定すると、投稿の作成・更新時にコンテンツからハッシュタグを自動抽出してタグを設定する。
func (h *PostHandler) SetAutoTagsUseCase(autoTags *usecase.SetAutoPostTagsUseCase) {
	h.autoTags = autoTags
}

// createPostRequest は投稿作成リクエスト。
type createPostRequest struct {
	Title        string             `json:"title" binding:"required,min=1,max=200"`
	Content      string             `json:"content" binding:"required,min=1,max=50000"`
	ImageURLs    string             `json:"image_urls" binding:"omitempty,max=2000"`
	IsDraft      bool               `json:"is_draft"`
	CodeSnippets []codeSnippetInput `json:"code_snippets" binding:"omitempty,max=20"`
}

// Create は新しい投稿を作成する。
func (h *PostHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	input := bindJSON[createPostRequest](c)
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
	created, err := h.uc.Create.Execute(c.Request.Context(), post)
	if err != nil {
		respondError(c, err)
		return
	}

	// コードスニペットを一括作成
	if len(input.CodeSnippets) > 0 && h.uc.CreateSnippet != nil {
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
			if _, err := h.uc.CreateSnippet.Execute(c.Request.Context(), snippet); err != nil {
				log.Printf("[ERROR] コードスニペット作成に失敗しました (post_id=%d): %v", created.ID, err)
			}
		}
		// スニペット付きで再取得
		if updated, err := h.uc.Get.Execute(c.Request.Context(), created.ID); err == nil {
			created = updated
		}
	}

	// コンテンツからハッシュタグを自動抽出してタグを設定
	if h.autoTags != nil {
		_ = h.autoTags.Execute(c.Request.Context(), created.ID, userID, input.Content)
	}

	h.processPostMentions(c, userID, created.ID, input.Content)

	respondCreated(c, created)
}

// processPostMentions は投稿本文のメンションを記録して相手に通知する。
// メンションの成否は投稿の成否と切り離す（本文は既に保存済みのため）。
func (h *PostHandler) processPostMentions(c *gin.Context, actorID, postID uint, content string) {
	if h.uc.ProcessMentions == nil || content == "" {
		return
	}
	_ = h.uc.ProcessMentions.Execute(c.Request.Context(), usecase.ProcessMentionsInput{
		ActorID:      actorID,
		Text:         content,
		PostID:       &postID,
		NotifyPostID: &postID,
	})
}

// processCommentMentions はコメント本文のメンションを記録して相手に通知する。
// 記録はコメントに紐づけ、通知からは元の投稿へ辿れるようにする。
func (h *PostHandler) processCommentMentions(c *gin.Context, actorID, commentID, postID uint, content string) {
	if h.uc.ProcessMentions == nil || content == "" {
		return
	}
	_ = h.uc.ProcessMentions.Execute(c.Request.Context(), usecase.ProcessMentionsInput{
		ActorID:      actorID,
		Text:         content,
		CommentID:    &commentID,
		NotifyPostID: &postID,
	})
}

// GetAll は投稿一覧をページネーション付きで返す。
func (h *PostHandler) GetAll(c *gin.Context) {
	page, limit := parsePagination(c)

	posts, err := h.uc.List.Execute(c.Request.Context(), page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.uc.Count.Execute(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, posts, total, page, limit)
}

// postDetailResponse は投稿詳細レスポンス（いいね済み・ブックマーク済みフラグ付き）。
type postDetailResponse struct {
	model.Post
	Liked      bool `json:"liked"`
	Bookmarked bool `json:"bookmarked"`
}

// GetByID は指定IDの投稿を返す。いいね済みフラグも付与する。
func (h *PostHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	post, err := h.uc.Get.Execute(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	userID := c.GetUint("userID")
	respondOK(c, postDetailResponse{
		Post:       *post,
		Liked:      h.hasLiked(c, userID, post.ID),
		Bookmarked: h.isBookmarked(c, userID, post.ID),
	})
}

// updatePostRequest は投稿更新リクエスト。
type updatePostRequest struct {
	Title     string `json:"title" binding:"omitempty,min=1,max=200"`
	Content   string `json:"content" binding:"omitempty,min=1,max=50000"`
	ImageURLs string `json:"image_urls" binding:"omitempty,max=2000"`
}

// Update は投稿を更新する。所有者のみ更新可能。
func (h *PostHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[updatePostRequest](c)
	if input == nil {
		return
	}

	post, err := h.uc.Update.Execute(c.Request.Context(), id, userID, input.Title, input.Content, input.ImageURLs)
	if err != nil {
		respondError(c, err)
		return
	}

	// コンテンツ更新時にハッシュタグを自動再抽出
	if h.autoTags != nil && input.Content != "" {
		_ = h.autoTags.Execute(c.Request.Context(), id, userID, input.Content)
	}

	if input.Content != "" {
		h.processPostMentions(c, userID, id, input.Content)
	}

	respondOK(c, post)
}

// Delete は投稿を削除する。所有者のみ削除可能。
func (h *PostHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.uc.Delete.Execute(c.Request.Context(), id, userID)
	})
}

// Timeline はフォロー中ユーザーの投稿タイムラインを返す。
func (h *PostHandler) Timeline(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	posts, err := h.uc.Timeline.Execute(c.Request.Context(), userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(posts))
}

// postListResponse は投稿一覧レスポンス（ページネーション付き）。
type postListResponse struct {
	Posts  []model.Post `json:"posts"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

// GetUserPosts は指定ユーザーの投稿一覧をページネーション付きで返す。
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	posts, total, err := h.uc.ListByUser.Execute(c.Request.Context(), id, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, postListResponse{
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

	posts, total, err := h.uc.ListByUser.Execute(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, postListResponse{
		Posts:  ensureSlice(posts),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetMyCount は認証ユーザーの投稿数を返す。
func (h *PostHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.uc.CountByUser.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// GetDraftsCount は認証ユーザーの下書き投稿数を返す。
func (h *PostHandler) GetDraftsCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.uc.CountDrafts.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// GetScheduledCount は認証ユーザーのスケジュール済み投稿数を返す。
func (h *PostHandler) GetScheduledCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.uc.CountScheduled.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// Like は投稿にいいねする。
func (h *PostHandler) Like(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.uc.Like.Execute(c.Request.Context(), userID, id)
	}, "liked")
}

// Unlike は投稿のいいねを取り消す。
func (h *PostHandler) Unlike(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.uc.Unlike.Execute(c.Request.Context(), userID, id)
	}, "unliked")
}

// GetComments は投稿のコメント一覧を返す。
func (h *PostHandler) GetComments(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	comments, err := h.uc.ListComments.Execute(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(comments))
}

// createCommentRequest はコメント作成リクエスト。
// ParentIDが指定された場合は返信コメントとして作成する。
type createCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=5000"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

// CreateComment は投稿にコメントを作成する。
func (h *PostHandler) CreateComment(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[createCommentRequest](c)
	if input == nil {
		return
	}

	comment, err := h.uc.CreateComment.Execute(c.Request.Context(), userID, id, input.Content, input.ParentID)
	if err != nil {
		respondError(c, err)
		return
	}

	h.processCommentMentions(c, userID, comment.ID, id, input.Content)

	respondCreated(c, comment)
}

// GetReplies は指定コメントへの返信一覧を返す。
func (h *PostHandler) GetReplies(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}

	replies, err := h.uc.ListReplies.Execute(c.Request.Context(), commentID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(replies))
}

// updateCommentRequest はコメント編集リクエスト。
type updateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"`
}

// EditComment はコメントを編集する。所有者のみ編集可能。
func (h *PostHandler) EditComment(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[updateCommentRequest](c)
	if input == nil {
		return
	}

	comment, err := h.uc.EditComment.Execute(c.Request.Context(), commentID, userID, input.Content)
	if err != nil {
		respondError(c, err)
		return
	}

	h.processCommentMentions(c, userID, commentID, comment.PostID, input.Content)

	respondOK(c, comment)
}

// DeleteComment はコメントを削除する。所有者のみ削除可能。
func (h *PostHandler) DeleteComment(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.uc.DeleteComment.Execute(c.Request.Context(), commentID, userID); err != nil {
		respondError(c, err)
		return
	}

	// コメントが消えたら、そのコメントを指すメンションも残さない
	if h.uc.DeleteCommentMentions != nil {
		_ = h.uc.DeleteCommentMentions.Execute(c.Request.Context(), commentID)
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

	if err := h.uc.HideComment.Execute(c.Request.Context(), commentID, userID); err != nil {
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

	if err := h.uc.UnhideComment.Execute(c.Request.Context(), commentID, userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("コメントの非表示を解除しました"))
}

// GetDrafts は現在のユーザーの下書き投稿一覧を返す。
func (h *PostHandler) GetDrafts(c *gin.Context) {
	userID := c.GetUint("userID")

	drafts, err := h.uc.ListDrafts.Execute(c.Request.Context(), userID)
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

	post, err := h.uc.Publish.Execute(c.Request.Context(), id, userID)
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

	post, err := h.uc.Unpublish.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, post)
}

// hasLiked は投稿詳細に載せるいいね済みフラグを返す。
// 取得に失敗しても投稿詳細自体は返したいため、移行前と同じく false にフォールバックする。
func (h *PostHandler) hasLiked(c *gin.Context, userID, postID uint) bool {
	liked, err := h.uc.HasLiked.Execute(c.Request.Context(), userID, postID)
	if err != nil {
		return false
	}
	return liked
}

// isBookmarked は投稿詳細に載せるブックマーク済みフラグを返す。
// 取得に失敗しても投稿詳細自体は返したいため、移行前と同じく false にフォールバックする。
func (h *PostHandler) isBookmarked(c *gin.Context, userID, postID uint) bool {
	bookmarked, err := h.uc.HasBookmarked.Execute(c.Request.Context(), userID, postID)
	if err != nil {
		return false
	}
	return bookmarked
}

// Bookmark は投稿をブックマークする。
func (h *PostHandler) Bookmark(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.uc.Bookmark.Execute(c.Request.Context(), userID, id)
	}, "bookmarked")
}

// Unbookmark は投稿のブックマークを解除する。
func (h *PostHandler) Unbookmark(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.uc.Unbookmark.Execute(c.Request.Context(), userID, id)
	}, "unbookmarked")
}

// GetBookmarks は現在のユーザーのブックマーク済み投稿一覧を返す。
func (h *PostHandler) GetBookmarks(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	posts, total, err := h.uc.ListBookmarks.Execute(c.Request.Context(), userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, postListResponse{
		Posts:  ensureSlice(posts),
		Total:  total,
		Limit:  limit,
		Offset: (page - 1) * limit,
	})
}

// GetBookmarksCount は認証ユーザーのブックマーク済み投稿数を返す。
func (h *PostHandler) GetBookmarksCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.uc.CountBookmarks.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// reactionRequest はリアクション追加/削除リクエスト。
type reactionRequest struct {
	Emoji string `json:"emoji" binding:"required,max=10"`
}

// AddReaction は投稿にリアクション（絵文字）を追加する。
func (h *PostHandler) AddReaction(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[reactionRequest](c)
	if input == nil {
		return
	}

	if err := h.uc.AddReaction.Execute(c.Request.Context(), userID, id, input.Emoji); err != nil {
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

	input := bindJSON[reactionRequest](c)
	if input == nil {
		return
	}

	if err := h.uc.RemoveReaction.Execute(c.Request.Context(), userID, id, input.Emoji); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("reaction_removed"))
}

// reactionResponse はリアクション一覧レスポンス。
type reactionResponse struct {
	Reactions     []model.ReactionCount `json:"reactions"`
	UserReactions []string              `json:"user_reactions"`
}

// GetReactions は投稿のリアクション一覧を返す。
func (h *PostHandler) GetReactions(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	reactions, userReactions, err := h.uc.GetReactions.Execute(c.Request.Context(), userID, id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, reactionResponse{
		Reactions:     reactions,
		UserReactions: userReactions,
	})
}

// schedulePublishRequest はスケジュール公開リクエスト。
type schedulePublishRequest struct {
	ScheduledAt string `json:"scheduled_at" binding:"required"`
}

// SchedulePublish は投稿のスケジュール公開日時を設定する。
func (h *PostHandler) SchedulePublish(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[schedulePublishRequest](c)
	if input == nil {
		return
	}

	scheduledAt, ok := parseDateTimeRFC3339(input.ScheduledAt)
	if !ok {
		respondBadRequest(c, "日時のフォーマットが不正です（RFC3339形式で指定してください）")
		return
	}

	post, err := h.uc.Schedule.Execute(c.Request.Context(), id, userID, scheduledAt)
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

	post, err := h.uc.CancelSchedule.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, post)
}

// GetScheduled はユーザーのスケジュール済み投稿一覧を返す。
func (h *PostHandler) GetScheduled(c *gin.Context) {
	userID := c.GetUint("userID")

	posts, err := h.uc.ListScheduled.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(posts))
}

// autoSaveDraftRequest は下書き自動保存リクエスト。
// IDが0の場合は新規作成、0以外の場合は既存下書きの更新。
type autoSaveDraftRequest struct {
	ID        uint   `json:"id"`
	Title     string `json:"title" binding:"omitempty,max=200"`
	Content   string `json:"content" binding:"omitempty,max=50000"`
	ImageURLs string `json:"image_urls" binding:"omitempty,max=2000"`
}

// autoSaveDraftResponse は下書き自動保存レスポンス。
type autoSaveDraftResponse struct {
	ID        uint   `json:"id"`
	UpdatedAt string `json:"updated_at"`
}

// AutoSaveDraft は下書きをサーバーサイドで自動保存する。
func (h *PostHandler) AutoSaveDraft(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[autoSaveDraftRequest](c)
	if req == nil {
		return
	}

	result, err := h.uc.AutoSaveDraft.Execute(c.Request.Context(), userID, req.ID, req.Title, req.Content, req.ImageURLs)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, autoSaveDraftResponse{
		ID:        result.ID,
		UpdatedAt: result.UpdatedAt.Format(time.RFC3339),
	})
}

// batchReactionRequest はリアクション一括取得リクエスト。
type batchReactionRequest struct {
	PostIDs []uint `json:"post_ids" binding:"required"`
}

// batchReactionResponse はリアクション一括取得レスポンス。
type batchReactionResponse struct {
	Reactions map[uint]reactionResponse `json:"reactions"`
}

// GetReactionsBatch は複数投稿のリアクション情報を一括取得する。
func (h *PostHandler) GetReactionsBatch(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[batchReactionRequest](c)
	if input == nil {
		return
	}

	reactions, userReactions, err := h.uc.ReactionsBatch.Execute(c.Request.Context(), userID, input.PostIDs)
	if err != nil {
		respondError(c, err)
		return
	}

	result := make(map[uint]reactionResponse, len(input.PostIDs))
	for _, id := range input.PostIDs {
		result[id] = reactionResponse{
			Reactions:     reactions[id],
			UserReactions: userReactions[id],
		}
	}

	respondOK(c, batchReactionResponse{Reactions: result})
}

// codeSnippetInput はコードスニペットの入力データ。
type codeSnippetInput struct {
	Language string `json:"language" binding:"omitempty,max=100"`
	FileName string `json:"file_name" binding:"omitempty,max=255"`
	Code     string `json:"code" binding:"omitempty,max=50000"`
}
