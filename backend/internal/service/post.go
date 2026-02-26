package service

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/norman6464/devsync/backend/internal/constants"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// EstimateReadTime はコンテンツの文字数から推定読了時間（分）を計算する純粋関数。
// 日本語基準で約500文字/分として計算し、最低1分を返す。
func EstimateReadTime(content string) int {
	charCount := utf8.RuneCountInString(content)
	minutes := charCount / 500
	if minutes < 1 {
		return 1
	}
	return minutes
}

// PostService は投稿に関するビジネスロジックを提供する。
// 投稿のCRUD操作に加え、フォロワーへの通知連携を行う。
type PostService struct {
	repo                repository.PostRepositoryInterface
	notificationService NotificationServiceInterface
}

// NewPostService は新しいPostServiceインスタンスを生成する。
func NewPostService(repo repository.PostRepositoryInterface, notificationService NotificationServiceInterface) *PostService {
	return &PostService{repo: repo, notificationService: notificationService}
}

// Create は新しい投稿を作成し、下書きでない場合はフォロワーに非同期で通知する。
func (s *PostService) Create(post *model.Post) (*model.Post, error) {
	// バリデーション
	v := validator.NewPostValidator()
	// タグは空配列として渡す（post.Tagsフィールドは未確認のため）
	if err := v.ValidateCreatePost(post.Title, post.Content, post.ImageURLs, nil); err != nil {
		return nil, err
	}

	// 読了時間を推定
	post.EstimatedReadTime = EstimateReadTime(post.Content)

	if err := s.repo.Create(post); err != nil {
		return nil, err
	}

	// 下書きでない場合のみフォロワーへ非同期で通知
	if !post.IsDraft {
		go s.notificationService.NotifyFollowers(post.UserID, post.ID, model.NotificationTypePost)
	}

	// アソシエーション付きで再取得
	created, err := s.repo.FindByID(post.ID)
	if err != nil {
		return post, nil
	}
	return created, nil
}

// GetByID は指定IDの投稿を取得する。
func (s *PostService) GetByID(id uint) (*model.Post, error) {
	return s.repo.FindByID(id)
}

// GetAll は投稿一覧をページネーション付きで取得する。
func (s *PostService) GetAll(page, limit int) ([]model.Post, error) {
	return s.repo.FindAll(page, limit)
}

// CountAll は公開済み投稿の総数を取得する。
func (s *PostService) CountAll() (int64, error) {
	return s.repo.CountAll()
}

// GetByUserID は指定ユーザーの投稿をページネーション付きで取得する。
func (s *PostService) GetByUserID(userID uint, limit, offset int) ([]model.Post, int64, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

// GetDrafts は指定ユーザーの下書き投稿を取得する。
func (s *PostService) GetDrafts(userID uint) ([]model.Post, error) {
	return s.repo.FindDraftsByUserID(userID)
}

// Timeline は指定ユーザーのタイムライン（フォロー中ユーザーの投稿）を取得する。
func (s *PostService) Timeline(userID uint, page, limit int) ([]model.Post, error) {
	return s.repo.Timeline(userID, page, limit)
}

// Update は所有権を検証した後、投稿を更新する。
func (s *PostService) Update(id, userID uint, title, content, imageUrls string) (*model.Post, error) {
	post, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	// 空白のみの入力を正規化（空白のみ→変更なしとして扱う）
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	imageUrls = strings.TrimSpace(imageUrls)

	// バリデーション
	v := validator.NewPostValidator()
	if err := v.ValidateUpdatePost(title, content, imageUrls); err != nil {
		return nil, err
	}

	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
		post.EstimatedReadTime = EstimateReadTime(content)
	}
	if imageUrls != "" {
		post.ImageURLs = imageUrls
	}

	if err := s.repo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

// Delete は所有権を検証した後、投稿を削除する。
func (s *PostService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// findAndCheckOwnership は投稿を取得し、指定ユーザーが所有者かを検証する。
func (s *PostService) findAndCheckOwnership(id, userID uint) (*model.Post, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(p *model.Post) uint { return p.UserID })
}

// findAndPreventSelfAction は投稿を取得し、自己操作でないことを検証する。
func (s *PostService) findAndPreventSelfAction(userID, postID uint) error {
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return ErrNotFound
	}
	if post.UserID == userID {
		return ErrForbidden
	}
	return nil
}

// Like は投稿にいいねを追加する。
// 自分の投稿にはいいねできない。
func (s *PostService) Like(userID, postID uint) error {
	if err := s.findAndPreventSelfAction(userID, postID); err != nil {
		return err
	}
	return s.repo.Like(userID, postID)
}

// Unlike は投稿のいいねを取り消す。
// 自分の投稿のいいねは取り消せない（そもそもいいねできないため）。
func (s *PostService) Unlike(userID, postID uint) error {
	if err := s.findAndPreventSelfAction(userID, postID); err != nil {
		return err
	}
	return s.repo.Unlike(userID, postID)
}

// HasLiked は指定ユーザーが投稿にいいね済みかを判定する。
func (s *PostService) HasLiked(userID, postID uint) bool {
	return s.repo.HasLiked(userID, postID)
}

// CreateComment は投稿にコメントを作成する。
// ParentIDが指定された場合は親コメントの存在確認と深さ制限（1階層）を行う。
func (s *PostService) CreateComment(comment *model.Comment) error {
	// バリデーション
	v := validator.NewPostValidator()
	if err := v.ValidateComment(comment.Content); err != nil {
		return err
	}

	// 返信の場合: 親コメントの存在確認と深さ制限
	if comment.ParentID != nil {
		parent, err := s.repo.FindCommentByID(*comment.ParentID)
		if err != nil {
			return domain.NewError(domain.ErrCodeNotFound, "親コメントが見つかりません", err)
		}
		if parent.PostID != comment.PostID {
			return domain.NewError(domain.ErrCodeBadRequest, "親コメントが別の投稿に属しています", nil)
		}
		// 深さ制限: 返信への返信は不可（1階層のみ）
		if parent.ParentID != nil {
			return domain.NewError(domain.ErrCodeBadRequest, "返信への返信はできません", nil)
		}
	}

	return s.repo.CreateComment(comment)
}

// GetComments は指定投稿の全コメントを取得する。
func (s *PostService) GetComments(postID uint) ([]model.Comment, error) {
	return s.repo.GetComments(postID)
}

// GetReplies は指定コメントへの返信一覧を取得する。
func (s *PostService) GetReplies(parentID uint) ([]model.Comment, error) {
	return s.repo.GetReplies(parentID)
}

// findAndCheckCommentOwnership はコメントの所有権を検証する共通ヘルパー。
func (s *PostService) findAndCheckCommentOwnership(id, userID uint) (*model.Comment, error) {
	comment, err := s.repo.FindCommentByID(id)
	if err != nil {
		return nil, err
	}

	if comment.UserID != userID {
		return nil, ErrForbidden
	}

	return comment, nil
}

// EditComment は所有権を検証した後、コメント内容を更新する。
func (s *PostService) EditComment(id, userID uint, content string) (*model.Comment, error) {
	v := validator.NewPostValidator()
	if err := v.ValidateComment(content); err != nil {
		return nil, err
	}

	comment, err := s.findAndCheckCommentOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	comment.Content = content
	if err := s.repo.UpdateComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// DeleteComment は所有権を検証した後、コメントを削除する。
func (s *PostService) DeleteComment(id, userID uint) error {
	if _, err := s.findAndCheckCommentOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.DeleteComment(id)
}

// HideComment は所有権を検証した後、コメントを非表示にする。
func (s *PostService) HideComment(id, userID uint) error {
	comment, err := s.findAndCheckCommentOwnership(id, userID)
	if err != nil {
		return err
	}
	comment.IsHidden = true
	return s.repo.UpdateComment(comment)
}

// UnhideComment は所有権を検証した後、コメントの非表示を解除する。
func (s *PostService) UnhideComment(id, userID uint) error {
	comment, err := s.findAndCheckCommentOwnership(id, userID)
	if err != nil {
		return err
	}
	comment.IsHidden = false
	return s.repo.UpdateComment(comment)
}

// Bookmark は投稿をブックマークする。
// 自分の投稿へのブックマークは禁止する。
func (s *PostService) Bookmark(userID, postID uint) error {
	if err := s.findAndPreventSelfAction(userID, postID); err != nil {
		return err
	}
	return s.repo.Bookmark(userID, postID)
}

// Unbookmark は投稿のブックマークを解除する。
// 自分の投稿のブックマークは解除できない（そもそもブックマークできないため）。
func (s *PostService) Unbookmark(userID, postID uint) error {
	if err := s.findAndPreventSelfAction(userID, postID); err != nil {
		return err
	}
	return s.repo.Unbookmark(userID, postID)
}

// HasBookmarked は指定ユーザーが投稿をブックマーク済みかを判定する。
func (s *PostService) HasBookmarked(userID, postID uint) bool {
	return s.repo.HasBookmarked(userID, postID)
}

// GetBookmarks は指定ユーザーのブックマーク済み投稿一覧を取得する。
func (s *PostService) GetBookmarks(userID uint, page, limit int) ([]model.Post, int64, error) {
	return s.repo.FindBookmarkedByUserID(userID, page, limit)
}

// AddReaction は投稿にリアクション（絵文字）を追加する。
// 許可された絵文字のみ使用可能。自分の投稿への自己リアクションは禁止する。
func (s *PostService) AddReaction(userID, postID uint, emoji string) error {
	if !constants.IsAllowedReactionEmoji(emoji) {
		return domain.NewError(domain.ErrCodeBadRequest, "許可されていない絵文字です: "+emoji, nil)
	}
	if err := s.findAndPreventSelfAction(userID, postID); err != nil {
		return err
	}
	return s.repo.AddReaction(userID, postID, emoji)
}

// RemoveReaction は投稿のリアクションを削除する。
// 自分の投稿へのリアクション削除は禁止する（そもそもリアクションできないため）。
func (s *PostService) RemoveReaction(userID, postID uint, emoji string) error {
	if !constants.IsAllowedReactionEmoji(emoji) {
		return domain.NewError(domain.ErrCodeBadRequest, "許可されていない絵文字です: "+emoji, nil)
	}
	if err := s.findAndPreventSelfAction(userID, postID); err != nil {
		return err
	}
	return s.repo.RemoveReaction(userID, postID, emoji)
}

// GetReactionsByPostID は指定投稿のリアクション集計を取得する。
func (s *PostService) GetReactionsByPostID(postID uint) ([]model.ReactionCount, error) {
	return s.repo.GetReactionsByPostID(postID)
}

// GetUserReactions は指定ユーザーが投稿に付けたリアクション絵文字一覧を取得する。
func (s *PostService) GetUserReactions(userID, postID uint) ([]string, error) {
	return s.repo.GetUserReactions(userID, postID)
}

// GetReactionsBatch は複数投稿のリアクション情報を一括取得する。
// 最大50件まで。
func (s *PostService) GetReactionsBatch(userID uint, postIDs []uint) (map[uint][]model.ReactionCount, map[uint][]string, error) {
	if len(postIDs) == 0 {
		return map[uint][]model.ReactionCount{}, map[uint][]string{}, nil
	}
	if len(postIDs) > 50 {
		return nil, nil, domain.NewError(domain.ErrCodeBadRequest, "一度に取得できる投稿は50件までです", nil)
	}

	reactions, err := s.repo.GetReactionsBatch(postIDs)
	if err != nil {
		return nil, nil, err
	}

	userReactions, err := s.repo.GetUserReactionsBatch(userID, postIDs)
	if err != nil {
		return nil, nil, err
	}

	return reactions, userReactions, nil
}

// SchedulePublish は下書き投稿にスケジュール公開日時を設定する。
func (s *PostService) SchedulePublish(id, userID uint, scheduledAt time.Time) (*model.Post, error) {
	post, err := s.findAndCheckOwnership(id, userID)
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
	if err := s.repo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

// CancelSchedule はスケジュール設定を解除し、投稿を通常の下書きに戻す。
func (s *PostService) CancelSchedule(id, userID uint) (*model.Post, error) {
	post, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}
	if post.ScheduledAt == nil {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "この投稿はスケジュールされていません", nil)
	}

	post.ScheduledAt = nil
	if err := s.repo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

// GetScheduled はユーザーのスケジュール済み投稿一覧を取得する。
func (s *PostService) GetScheduled(userID uint) ([]model.Post, error) {
	return s.repo.FindScheduledByUserID(userID)
}

// AutoSaveDraft は下書きをサーバーサイドで自動保存する。
// draftIDが0の場合は新規下書きを作成し、0以外の場合は既存下書きを更新する。
// 作成途中の投稿を許容するため、タイトル・本文は空でも許可する（長さ上限のみ検証）。
func (s *PostService) AutoSaveDraft(userID, draftID uint, title, content, imageURLs string) (*model.Post, error) {
	// 長さ上限のみバリデーション
	if len([]rune(title)) > 200 {
		return nil, domain.NewError(domain.ErrCodeValidation, "タイトルは200文字以下である必要があります", nil)
	}
	if len([]rune(content)) > 50000 {
		return nil, domain.NewError(domain.ErrCodeValidation, "本文は50000文字以下である必要があります", nil)
	}

	// 新規作成
	if draftID == 0 {
		post := &model.Post{
			UserID:            userID,
			Title:             title,
			Content:           content,
			ImageURLs:         imageURLs,
			IsDraft:           true,
			EstimatedReadTime: EstimateReadTime(content),
		}
		if err := s.repo.Create(post); err != nil {
			return nil, err
		}
		return post, nil
	}

	// 既存下書きの更新
	post, err := s.findAndCheckOwnership(draftID, userID)
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

	if err := s.repo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

// Publish は下書き投稿を公開し、フォロワーに通知する。
func (s *PostService) Publish(id, userID uint) (*model.Post, error) {
	post, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}
	if !post.IsDraft {
		return nil, ErrBadRequest
	}

	post.IsDraft = false
	if err := s.repo.Update(post); err != nil {
		return nil, err
	}

	// 公開時にフォロワーへ非同期で通知
	go s.notificationService.NotifyFollowers(post.UserID, post.ID, model.NotificationTypePost)

	return post, nil
}

// Unpublish は公開済みの投稿を下書きに戻す。
// 既に下書きの場合はBadRequestエラーを返す。
func (s *PostService) Unpublish(id, userID uint) (*model.Post, error) {
	post, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}
	if post.IsDraft {
		return nil, ErrBadRequest
	}

	post.IsDraft = true
	if err := s.repo.Update(post); err != nil {
		return nil, err
	}

	return post, nil
}
