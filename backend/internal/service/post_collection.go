package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// PostCollectionService は投稿コレクションに関するビジネスロジックを提供する。
// テーマ別の投稿コレクションのCRUD操作と投稿管理を担当する。
type PostCollectionService struct {
	repo repository.PostCollectionRepositoryInterface
}

// NewPostCollectionService は新しいPostCollectionServiceインスタンスを生成する。
func NewPostCollectionService(repo repository.PostCollectionRepositoryInterface) *PostCollectionService {
	return &PostCollectionService{repo: repo}
}

// Create は新しい投稿コレクションを作成する。
func (s *PostCollectionService) Create(collection *model.PostCollection) (*model.PostCollection, error) {
	if collection.Title == "" {
		return nil, domain.NewError(domain.ErrCodeValidation, "タイトルは必須です", nil)
	}
	if err := s.repo.Create(collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// GetByID は指定IDのコレクションを取得する。
func (s *PostCollectionService) GetByID(id uint) (*model.PostCollection, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの全コレクションを取得する。
func (s *PostCollectionService) GetByUserID(userID uint) ([]model.PostCollection, error) {
	return s.repo.FindByUserID(userID)
}

// GetPublicByUserID は指定ユーザーの公開コレクションを取得する。
func (s *PostCollectionService) GetPublicByUserID(userID uint) ([]model.PostCollection, error) {
	return s.repo.FindPublicByUserID(userID)
}

// Update は所有権を検証した後、コレクションを更新する。
func (s *PostCollectionService) Update(id, userID uint, title, description string, isPublic bool) (*model.PostCollection, error) {
	collection, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if collection.UserID != userID {
		return nil, ErrForbidden
	}

	collection.Title = title
	collection.Description = description
	collection.IsPublic = isPublic

	if err := s.repo.Update(collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// Delete は所有権を検証した後、コレクションを削除する。
func (s *PostCollectionService) Delete(id, userID uint) error {
	collection, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if collection.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// AddPost は所有権を検証した後、コレクションに投稿を追加する。
func (s *PostCollectionService) AddPost(collectionID, userID, postID uint, note string) error {
	collection, err := s.repo.FindByID(collectionID)
	if err != nil {
		return err
	}
	if collection.UserID != userID {
		return ErrForbidden
	}

	item := &model.PostCollectionItem{
		CollectionID: collectionID,
		PostID:       postID,
		Note:         note,
	}
	return s.repo.AddPost(item)
}

// RemovePost は所有権を検証した後、コレクションから投稿を削除する。
func (s *PostCollectionService) RemovePost(collectionID, userID, postID uint) error {
	collection, err := s.repo.FindByID(collectionID)
	if err != nil {
		return err
	}
	if collection.UserID != userID {
		return ErrForbidden
	}
	return s.repo.RemovePost(collectionID, postID)
}

// GetPosts は指定コレクションの投稿一覧を取得する。
func (s *PostCollectionService) GetPosts(collectionID uint) ([]model.PostCollectionItem, error) {
	return s.repo.GetPostsByCollectionID(collectionID)
}
