package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// PostService handles post business logic.
type PostService struct {
	repo                repository.PostRepositoryInterface
	notificationService *NotificationService
}

// NewPostService creates a new PostService.
func NewPostService(repo repository.PostRepositoryInterface, notificationService *NotificationService) *PostService {
	return &PostService{repo: repo, notificationService: notificationService}
}

// Create creates a new post and notifies followers.
func (s *PostService) Create(post *model.Post) (*model.Post, error) {
	if err := s.repo.Create(post); err != nil {
		return nil, err
	}

	// Notify followers asynchronously
	go s.notificationService.NotifyFollowers(post.UserID, post.ID, model.NotificationTypePost)

	// Reload with associations
	created, err := s.repo.FindByID(post.ID)
	if err != nil {
		return post, nil
	}
	return created, nil
}

// GetByID returns a post by ID.
func (s *PostService) GetByID(id uint) (*model.Post, error) {
	return s.repo.FindByID(id)
}

// GetAll returns paginated posts.
func (s *PostService) GetAll(page, limit int) ([]model.Post, error) {
	return s.repo.FindAll(page, limit)
}

// GetByUserID returns all posts for a user.
func (s *PostService) GetByUserID(userID uint) ([]model.Post, error) {
	return s.repo.FindByUserID(userID)
}

// Timeline returns the timeline for a user.
func (s *PostService) Timeline(userID uint, page, limit int) ([]model.Post, error) {
	return s.repo.Timeline(userID, page, limit)
}

// Update updates a post after verifying ownership.
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

// Delete deletes a post after verifying ownership.
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

// Like likes a post.
func (s *PostService) Like(userID, postID uint) error {
	return s.repo.Like(userID, postID)
}

// Unlike unlikes a post.
func (s *PostService) Unlike(userID, postID uint) error {
	return s.repo.Unlike(userID, postID)
}

// HasLiked checks if a user has liked a post.
func (s *PostService) HasLiked(userID, postID uint) bool {
	return s.repo.HasLiked(userID, postID)
}

// CreateComment creates a comment on a post.
func (s *PostService) CreateComment(comment *model.Comment) error {
	return s.repo.CreateComment(comment)
}

// GetComments returns all comments for a post.
func (s *PostService) GetComments(postID uint) ([]model.Comment, error) {
	return s.repo.GetComments(postID)
}

// DeleteComment deletes a comment.
func (s *PostService) DeleteComment(id, userID uint) error {
	return s.repo.DeleteComment(id, userID)
}
