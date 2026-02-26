package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// BookmarkCollectionService はブックマークコレクションのビジネスロジックを提供する。
type BookmarkCollectionService struct {
	repo repository.BookmarkCollectionRepositoryInterface
}

// NewBookmarkCollectionService は新しいBookmarkCollectionServiceインスタンスを生成する。
func NewBookmarkCollectionService(repo repository.BookmarkCollectionRepositoryInterface) *BookmarkCollectionService {
	return &BookmarkCollectionService{repo: repo}
}

// Create は新しいブックマークコレクションを作成する。
func (s *BookmarkCollectionService) Create(collection *model.BookmarkCollection) error {
	if err := domain.ValidateStringLength(collection.Name, 1, 100, "コレクション名"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(collection.Description, 0, 500, "説明"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(collection.Color, 0, 20, "カラー"); err != nil {
		return err
	}
	collection.Name = strings.TrimSpace(collection.Name)
	collection.Description = strings.TrimSpace(collection.Description)
	collection.Color = strings.TrimSpace(collection.Color)
	return s.repo.Create(collection)
}

// GetByUserID はユーザーのコレクション一覧を取得する。
func (s *BookmarkCollectionService) GetByUserID(userID uint) ([]model.BookmarkCollection, error) {
	return s.repo.FindByUserID(userID)
}

// Update はコレクションを更新する。所有権を検証する。
func (s *BookmarkCollectionService) Update(id, userID uint, updates *model.BookmarkCollection) (*model.BookmarkCollection, error) {
	collection, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(updates.Name) != "" {
		if err := domain.ValidateStringLength(updates.Name, 1, 100, "コレクション名"); err != nil {
			return nil, err
		}
		collection.Name = strings.TrimSpace(updates.Name)
	}
	if strings.TrimSpace(updates.Description) != "" {
		if err := domain.ValidateStringLength(updates.Description, 1, 500, "説明"); err != nil {
			return nil, err
		}
		collection.Description = strings.TrimSpace(updates.Description)
	}
	if strings.TrimSpace(updates.Color) != "" {
		if err := domain.ValidateStringLength(updates.Color, 1, 20, "カラー"); err != nil {
			return nil, err
		}
		collection.Color = strings.TrimSpace(updates.Color)
	}

	if err := s.repo.Update(collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// Delete はコレクションを削除する。所有権を検証する。
func (s *BookmarkCollectionService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// AddPost はコレクションにブックマークを追加する。
func (s *BookmarkCollectionService) AddPost(collectionID, postID, userID uint) error {
	if _, err := s.findAndCheckOwnership(collectionID, userID); err != nil {
		return err
	}

	exists, err := s.repo.HasPost(collectionID, postID)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewError(domain.ErrCodeConflict, "この投稿は既にコレクションに追加されています", nil)
	}

	return s.repo.AddPost(&model.BookmarkCollectionItem{
		CollectionID: collectionID,
		PostID:       postID,
	})
}

// RemovePost はコレクションからブックマークを削除する。
func (s *BookmarkCollectionService) RemovePost(collectionID, postID, userID uint) error {
	if _, err := s.findAndCheckOwnership(collectionID, userID); err != nil {
		return err
	}
	return s.repo.RemovePost(collectionID, postID)
}

// GetPosts はコレクション内の投稿一覧を取得する。
func (s *BookmarkCollectionService) GetPosts(collectionID uint, limit, offset int) ([]model.Post, int64, error) {
	return s.repo.GetPosts(collectionID, limit, offset)
}

// findAndCheckOwnership はコレクションを取得し、所有権を検証する。
func (s *BookmarkCollectionService) findAndCheckOwnership(id, userID uint) (*model.BookmarkCollection, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(c *model.BookmarkCollection) uint { return c.UserID })
}

// CountByUserID は指定ユーザーのコレクション総数を返す。
func (s *BookmarkCollectionService) CountByUserID(userID uint) (int64, error) {
	return s.repo.CountByUserID(userID)
}
