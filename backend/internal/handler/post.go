package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
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
			h.uc.CreateSnippet.Execute(c.Request.Context(), snippet)
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
	respondOK(c, dto.PostDetailResponse{
		Post:       *post,
		Liked:      h.hasLiked(c, userID, post.ID),
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

	posts, total, err := h.uc.ListByUser.Execute(c.Request.Context(), userID, limit, offset)
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

	count, err := h.uc.CountBookmarks.Execute(c.Request.Context(), userID)
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

	input := bindJSON[dto.ReactionRequest](c)
	if input == nil {
		return
	}

	if err := h.uc.RemoveReaction.Execute(c.Request.Context(), userID, id, input.Emoji); err != nil {
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

	reactions, userReactions, err := h.uc.GetReactions.Execute(c.Request.Context(), userID, id)
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

// AutoSaveDraft は下書きをサーバーサイドで自動保存する。
func (h *PostHandler) AutoSaveDraft(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.AutoSaveDraftRequest](c)
	if req == nil {
		return
	}

	result, err := h.uc.AutoSaveDraft.Execute(c.Request.Context(), userID, req.ID, req.Title, req.Content, req.ImageURLs)
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

	reactions, userReactions, err := h.uc.ReactionsBatch.Execute(c.Request.Context(), userID, input.PostIDs)
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
