package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

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

// GetByUserID は指定ユーザーの全投稿を取得する。
func (s *PostService) GetByUserID(userID uint) ([]model.Post, error) {
	return s.repo.FindByUserID(userID)
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
func (s *PostService) Update(id, userID uint, title, content string) (*model.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if post.UserID != userID {
		return nil, ErrForbidden
	}

	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
	}

	if err := s.repo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

// Delete は所有権を検証した後、投稿を削除する。
func (s *PostService) Delete(id, userID uint) error {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// Like は投稿にいいねを追加する。
func (s *PostService) Like(userID, postID uint) error {
	return s.repo.Like(userID, postID)
}

// Unlike は投稿のいいねを取り消す。
func (s *PostService) Unlike(userID, postID uint) error {
	return s.repo.Unlike(userID, postID)
}

// HasLiked は指定ユーザーが投稿にいいね済みかを判定する。
func (s *PostService) HasLiked(userID, postID uint) bool {
	return s.repo.HasLiked(userID, postID)
}

// CreateComment は投稿にコメントを作成する。
func (s *PostService) CreateComment(comment *model.Comment) error {
	return s.repo.CreateComment(comment)
}

// GetComments は指定投稿の全コメントを取得する。
func (s *PostService) GetComments(postID uint) ([]model.Comment, error) {
	return s.repo.GetComments(postID)
}

// DeleteComment はコメントを削除する。
func (s *PostService) DeleteComment(id, userID uint) error {
	return s.repo.DeleteComment(id, userID)
}

// Publish は下書き投稿を公開し、フォロワーに通知する。
func (s *PostService) Publish(id, userID uint) (*model.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if post.UserID != userID {
		return nil, ErrForbidden
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
