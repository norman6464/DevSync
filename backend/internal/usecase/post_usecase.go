package usecase

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// EstimateReadTime はコンテンツの文字数から推定読了時間（分）を計算する純粋関数。
// 日本語基準で約 500 文字 / 分として計算し、最低 1 分を返す。
func EstimateReadTime(content string) int {
	minutes := utf8.RuneCountInString(content) / 500
	if minutes < 1 {
		return 1
	}
	return minutes
}

// CreatePostUseCase は投稿を作成する。
type CreatePostUseCase struct {
	posts     repository.PostRepository
	followers *NotifyFollowersUseCase
}

// NewCreatePostUseCase は CreatePostUseCase を生成する。
func NewCreatePostUseCase(posts repository.PostRepository, followers *NotifyFollowersUseCase) *CreatePostUseCase {
	return &CreatePostUseCase{posts: posts, followers: followers}
}

// Execute は入力を検証して投稿を作成し、下書きでなければフォロワーへ非同期に通知する。
// 作成後はアソシエーション付きで取得し直して返す（取得に失敗した場合は作成した投稿をそのまま返す）。
func (uc *CreatePostUseCase) Execute(ctx context.Context, post *model.Post) (*model.Post, error) {
	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)

	if err := validator.NewPostValidator().ValidateCreatePost(post.Title, post.Content, post.ImageURLs, nil); err != nil {
		return nil, err
	}

	post.EstimatedReadTime = EstimateReadTime(post.Content)

	if err := uc.posts.Create(ctx, post); err != nil {
		return nil, err
	}

	if !post.IsDraft {
		uc.followers.Notify(ctx, post.UserID, post.ID, model.NotificationTypePost)
	}

	// 作成自体は既に成功しているため、再取得の失敗だけで全体を失敗にはしない。
	created, err := uc.posts.FindByID(ctx, post.ID)
	if err != nil || created == nil {
		return post, nil //nolint:nilerr // 再取得失敗時は手元の値へ意図的にフォールバックする
	}
	return created, nil
}

// GetPostUseCase は投稿を 1 件取得する。
type GetPostUseCase struct {
	posts repository.PostRepository
}

// NewGetPostUseCase は GetPostUseCase を生成する。
func NewGetPostUseCase(posts repository.PostRepository) *GetPostUseCase {
	return &GetPostUseCase{posts: posts}
}

// Execute は投稿を投稿者・コードスニペット付きで返す。
func (uc *GetPostUseCase) Execute(ctx context.Context, id uint) (*model.Post, error) {
	post, err := uc.posts.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, domain.ErrNotFound
	}
	return post, nil
}

// ListPostsUseCase は公開済み投稿の一覧を取得する。
type ListPostsUseCase struct {
	posts repository.PostRepository
}

// NewListPostsUseCase は ListPostsUseCase を生成する。
func NewListPostsUseCase(posts repository.PostRepository) *ListPostsUseCase {
	return &ListPostsUseCase{posts: posts}
}

// Execute は公開済み投稿を新しい順に返す。
func (uc *ListPostsUseCase) Execute(ctx context.Context, page, limit int) ([]model.Post, error) {
	return uc.posts.FindAll(ctx, page, limit)
}

// CountPostsUseCase は公開済み投稿の総数を取得する。
type CountPostsUseCase struct {
	posts repository.PostRepository
}

// NewCountPostsUseCase は CountPostsUseCase を生成する。
func NewCountPostsUseCase(posts repository.PostRepository) *CountPostsUseCase {
	return &CountPostsUseCase{posts: posts}
}

// Execute は公開済み投稿の総数を返す。
func (uc *CountPostsUseCase) Execute(ctx context.Context) (int64, error) {
	return uc.posts.CountAll(ctx)
}

// ListUserPostsUseCase は指定ユーザーの公開済み投稿を取得する。
type ListUserPostsUseCase struct {
	posts repository.PostRepository
}

// NewListUserPostsUseCase は ListUserPostsUseCase を生成する。
func NewListUserPostsUseCase(posts repository.PostRepository) *ListUserPostsUseCase {
	return &ListUserPostsUseCase{posts: posts}
}

// Execute は指定ユーザーの公開済み投稿と総件数を返す。
func (uc *ListUserPostsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.Post, int64, error) {
	return uc.posts.FindByUserID(ctx, userID, limit, offset)
}

// ListDraftPostsUseCase は下書き一覧を取得する。
type ListDraftPostsUseCase struct {
	posts repository.PostRepository
}

// NewListDraftPostsUseCase は ListDraftPostsUseCase を生成する。
func NewListDraftPostsUseCase(posts repository.PostRepository) *ListDraftPostsUseCase {
	return &ListDraftPostsUseCase{posts: posts}
}

// Execute は指定ユーザーの下書きを返す。
func (uc *ListDraftPostsUseCase) Execute(ctx context.Context, userID uint) ([]model.Post, error) {
	return uc.posts.FindDraftsByUserID(ctx, userID)
}

// ListScheduledPostsUseCase は公開予約済み投稿の一覧を取得する。
type ListScheduledPostsUseCase struct {
	posts repository.PostRepository
}

// NewListScheduledPostsUseCase は ListScheduledPostsUseCase を生成する。
func NewListScheduledPostsUseCase(posts repository.PostRepository) *ListScheduledPostsUseCase {
	return &ListScheduledPostsUseCase{posts: posts}
}

// Execute は指定ユーザーの公開予約済み投稿を返す。
func (uc *ListScheduledPostsUseCase) Execute(ctx context.Context, userID uint) ([]model.Post, error) {
	return uc.posts.FindScheduledByUserID(ctx, userID)
}

// GetTimelineUseCase はタイムラインを取得する。
type GetTimelineUseCase struct {
	posts repository.PostRepository
}

// NewGetTimelineUseCase は GetTimelineUseCase を生成する。
func NewGetTimelineUseCase(posts repository.PostRepository) *GetTimelineUseCase {
	return &GetTimelineUseCase{posts: posts}
}

// Execute はフォロー中ユーザーと自分の公開済み投稿を返す。
func (uc *GetTimelineUseCase) Execute(ctx context.Context, userID uint, page, limit int) ([]model.Post, error) {
	return uc.posts.Timeline(ctx, userID, page, limit)
}

// CountUserPostsUseCase は指定ユーザーの公開済み投稿数を取得する。
type CountUserPostsUseCase struct {
	posts repository.PostRepository
}

// NewCountUserPostsUseCase は CountUserPostsUseCase を生成する。
func NewCountUserPostsUseCase(posts repository.PostRepository) *CountUserPostsUseCase {
	return &CountUserPostsUseCase{posts: posts}
}

// Execute は公開済み投稿数を返す。
func (uc *CountUserPostsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.posts.CountByUserID(ctx, userID)
}

// CountUserDraftsUseCase は指定ユーザーの下書き数を取得する。
type CountUserDraftsUseCase struct {
	posts repository.PostRepository
}

// NewCountUserDraftsUseCase は CountUserDraftsUseCase を生成する。
func NewCountUserDraftsUseCase(posts repository.PostRepository) *CountUserDraftsUseCase {
	return &CountUserDraftsUseCase{posts: posts}
}

// Execute は下書き数を返す。
func (uc *CountUserDraftsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.posts.CountDraftsByUserID(ctx, userID)
}

// CountUserScheduledPostsUseCase は指定ユーザーの公開予約済み投稿数を取得する。
type CountUserScheduledPostsUseCase struct {
	posts repository.PostRepository
}

// NewCountUserScheduledPostsUseCase は CountUserScheduledPostsUseCase を生成する。
func NewCountUserScheduledPostsUseCase(posts repository.PostRepository) *CountUserScheduledPostsUseCase {
	return &CountUserScheduledPostsUseCase{posts: posts}
}

// Execute は公開予約済み投稿数を返す。
func (uc *CountUserScheduledPostsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.posts.CountScheduledByUserID(ctx, userID)
}

// UpdatePostUseCase は投稿を更新する。
type UpdatePostUseCase struct {
	posts repository.PostRepository
}

// NewUpdatePostUseCase は UpdatePostUseCase を生成する。
func NewUpdatePostUseCase(posts repository.PostRepository) *UpdatePostUseCase {
	return &UpdatePostUseCase{posts: posts}
}

// Execute は所有権と入力を検証したうえで投稿を更新する。
// 空文字（空白のみを含む）の項目は「変更なし」として扱う。
func (uc *UpdatePostUseCase) Execute(ctx context.Context, id, userID uint, title, content, imageURLs string) (*model.Post, error) {
	post, err := ensurePostOwner(ctx, uc.posts, id, userID)
	if err != nil {
		return nil, err
	}

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	imageURLs = strings.TrimSpace(imageURLs)

	if err := validator.NewPostValidator().ValidateUpdatePost(title, content, imageURLs); err != nil {
		return nil, err
	}

	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
		post.EstimatedReadTime = EstimateReadTime(content)
	}
	if imageURLs != "" {
		post.ImageURLs = imageURLs
	}

	if err := uc.posts.Update(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

// DeletePostUseCase は投稿を削除する。
type DeletePostUseCase struct {
	posts repository.PostRepository
}

// NewDeletePostUseCase は DeletePostUseCase を生成する。
func NewDeletePostUseCase(posts repository.PostRepository) *DeletePostUseCase {
	return &DeletePostUseCase{posts: posts}
}

// Execute は所有権を検証したうえで投稿を削除する。
func (uc *DeletePostUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensurePostOwner(ctx, uc.posts, id, userID); err != nil {
		return err
	}
	return uc.posts.Delete(ctx, id)
}

// PublishPostUseCase は下書きを公開する。
type PublishPostUseCase struct {
	posts     repository.PostRepository
	followers *NotifyFollowersUseCase
}

// NewPublishPostUseCase は PublishPostUseCase を生成する。
func NewPublishPostUseCase(posts repository.PostRepository, followers *NotifyFollowersUseCase) *PublishPostUseCase {
	return &PublishPostUseCase{posts: posts, followers: followers}
}

// Execute は所有権を検証したうえで下書きを公開し、フォロワーへ非同期に通知する。
func (uc *PublishPostUseCase) Execute(ctx context.Context, id, userID uint) (*model.Post, error) {
	post, err := ensurePostOwner(ctx, uc.posts, id, userID)
	if err != nil {
		return nil, err
	}
	if !post.IsDraft {
		return nil, domain.ErrBadRequest
	}

	post.IsDraft = false
	// 公開した時点で公開予約は役目を終える。残すと下書きへ戻したときに予約が復活し、
	// 予約一覧・件数に再び現れてしまう。
	post.ScheduledAt = nil
	if err := uc.posts.Update(ctx, post); err != nil {
		return nil, err
	}

	uc.followers.Notify(ctx, post.UserID, post.ID, model.NotificationTypePost)
	return post, nil
}

// UnpublishPostUseCase は公開済み投稿を下書きに戻す。
type UnpublishPostUseCase struct {
	posts repository.PostRepository
}

// NewUnpublishPostUseCase は UnpublishPostUseCase を生成する。
func NewUnpublishPostUseCase(posts repository.PostRepository) *UnpublishPostUseCase {
	return &UnpublishPostUseCase{posts: posts}
}

// Execute は所有権を検証したうえで投稿を下書きに戻す。
func (uc *UnpublishPostUseCase) Execute(ctx context.Context, id, userID uint) (*model.Post, error) {
	post, err := ensurePostOwner(ctx, uc.posts, id, userID)
	if err != nil {
		return nil, err
	}
	if post.IsDraft {
		return nil, domain.ErrBadRequest
	}

	post.IsDraft = true
	if err := uc.posts.Update(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

// SchedulePostPublishUseCase は下書きに公開予定日時を設定する。
type SchedulePostPublishUseCase struct {
	posts repository.PostRepository
}

// NewSchedulePostPublishUseCase は SchedulePostPublishUseCase を生成する。
func NewSchedulePostPublishUseCase(posts repository.PostRepository) *SchedulePostPublishUseCase {
	return &SchedulePostPublishUseCase{posts: posts}
}

// Execute は所有権を検証したうえで公開予定日時を設定する。
// 公開済みの投稿・過去の日時は受け付けない。
func (uc *SchedulePostPublishUseCase) Execute(ctx context.Context, id, userID uint, scheduledAt time.Time) (*model.Post, error) {
	post, err := ensurePostOwner(ctx, uc.posts, id, userID)
	if err != nil {
		return nil, err
	}
	if !post.IsDraft {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "公開済みの投稿はスケジュールできません", nil)
	}
	if scheduledAt.Before(time.Now()) {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "過去の日時は指定できません", nil)
	}

	post.ScheduledAt = &scheduledAt
	if err := uc.posts.Update(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

// CancelPostScheduleUseCase は公開予約を解除する。
type CancelPostScheduleUseCase struct {
	posts repository.PostRepository
}

// NewCancelPostScheduleUseCase は CancelPostScheduleUseCase を生成する。
func NewCancelPostScheduleUseCase(posts repository.PostRepository) *CancelPostScheduleUseCase {
	return &CancelPostScheduleUseCase{posts: posts}
}

// Execute は所有権を検証したうえで公開予約を解除し、通常の下書きに戻す。
func (uc *CancelPostScheduleUseCase) Execute(ctx context.Context, id, userID uint) (*model.Post, error) {
	post, err := ensurePostOwner(ctx, uc.posts, id, userID)
	if err != nil {
		return nil, err
	}
	if post.ScheduledAt == nil {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "この投稿はスケジュールされていません", nil)
	}

	post.ScheduledAt = nil
	if err := uc.posts.Update(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

// AutoSaveDraftUseCase は下書きをサーバー側で自動保存する。
type AutoSaveDraftUseCase struct {
	posts repository.PostRepository
}

// NewAutoSaveDraftUseCase は AutoSaveDraftUseCase を生成する。
func NewAutoSaveDraftUseCase(posts repository.PostRepository) *AutoSaveDraftUseCase {
	return &AutoSaveDraftUseCase{posts: posts}
}

// Execute は draftID が 0 なら新規下書きを作成し、0 以外なら既存下書きを更新する。
// 作成途中の投稿を許容するため、タイトル・本文は空でも許可する（長さ上限のみ検証）。
func (uc *AutoSaveDraftUseCase) Execute(ctx context.Context, userID, draftID uint, title, content, imageURLs string) (*model.Post, error) {
	if err := domain.ValidateStringLength(title, 0, 200, "タイトル"); err != nil {
		return nil, err
	}
	if err := domain.ValidateStringLength(content, 0, 50000, "本文"); err != nil {
		return nil, err
	}

	if draftID == 0 {
		post := &model.Post{
			UserID:            userID,
			Title:             title,
			Content:           content,
			ImageURLs:         imageURLs,
			IsDraft:           true,
			EstimatedReadTime: EstimateReadTime(content),
		}
		if err := uc.posts.Create(ctx, post); err != nil {
			return nil, err
		}
		return post, nil
	}

	post, err := ensurePostOwner(ctx, uc.posts, draftID, userID)
	if err != nil {
		return nil, err
	}
	if !post.IsDraft {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "公開済みの投稿は自動保存できません", nil)
	}

	post.Title = title
	post.Content = content
	if imageURLs != "" {
		post.ImageURLs = imageURLs
	}
	post.EstimatedReadTime = EstimateReadTime(content)

	if err := uc.posts.Update(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

// LikePostUseCase は投稿にいいねする。
type LikePostUseCase struct {
	likes repository.PostLikeRepository
	posts repository.PostAuthorReader
}

// NewLikePostUseCase は LikePostUseCase を生成する。
func NewLikePostUseCase(likes repository.PostLikeRepository, posts repository.PostAuthorReader) *LikePostUseCase {
	return &LikePostUseCase{likes: likes, posts: posts}
}

// Execute は投稿者を検証したうえでいいねする。自分の投稿にはいいねできない。
func (uc *LikePostUseCase) Execute(ctx context.Context, userID, postID uint) error {
	if err := ensureNotOwnPost(ctx, uc.posts, userID, postID); err != nil {
		return err
	}
	return uc.likes.Like(ctx, userID, postID)
}

// UnlikePostUseCase は投稿のいいねを取り消す。
type UnlikePostUseCase struct {
	likes repository.PostLikeRepository
	posts repository.PostAuthorReader
}

// NewUnlikePostUseCase は UnlikePostUseCase を生成する。
func NewUnlikePostUseCase(likes repository.PostLikeRepository, posts repository.PostAuthorReader) *UnlikePostUseCase {
	return &UnlikePostUseCase{likes: likes, posts: posts}
}

// Execute は投稿者を検証したうえでいいねを取り消す。
// 自分の投稿にはそもそもいいねできないため、取り消しも同じ条件で弾く。
func (uc *UnlikePostUseCase) Execute(ctx context.Context, userID, postID uint) error {
	if err := ensureNotOwnPost(ctx, uc.posts, userID, postID); err != nil {
		return err
	}
	return uc.likes.Unlike(ctx, userID, postID)
}

// HasLikedPostUseCase は投稿にいいね済みかを判定する。
type HasLikedPostUseCase struct {
	likes repository.PostLikeRepository
}

// NewHasLikedPostUseCase は HasLikedPostUseCase を生成する。
func NewHasLikedPostUseCase(likes repository.PostLikeRepository) *HasLikedPostUseCase {
	return &HasLikedPostUseCase{likes: likes}
}

// Execute は指定ユーザーが投稿にいいね済みかを返す。
func (uc *HasLikedPostUseCase) Execute(ctx context.Context, userID, postID uint) (bool, error) {
	return uc.likes.HasLiked(ctx, userID, postID)
}

// ensurePostOwner は投稿を取得し、userID が投稿者であることを検証する。
// 不在は 404、他人の投稿は 403 を返す。
func ensurePostOwner(ctx context.Context, posts repository.PostRepository, id, userID uint) (*model.Post, error) {
	return ensureOwner(ctx, posts.FindByID, id, userID, func(p *model.Post) uint { return p.UserID })
}
