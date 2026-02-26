package service

import (
	"strings"

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
	if err := domain.ValidateStringLength(collection.Title, 1, 200, "タイトル"); err != nil {
		return nil, err
	}
	if len([]rune(strings.TrimSpace(collection.Description))) > 1000 {
		return nil, domain.NewError(domain.ErrCodeValidation, "説明は1000文字以下である必要があります", nil)
	}
	collection.Title = strings.TrimSpace(collection.Title)
	if err := s.repo.Create(collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// GetByID は指定IDのコレクションを取得する。
func (s *PostCollectionService) GetByID(id uint) (*model.PostCollection, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの全コレクションをページネーション付きで取得する。
func (s *PostCollectionService) GetByUserID(userID uint, limit, offset int) ([]model.PostCollection, int64, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

// GetPublicByUserID は指定ユーザーの公開コレクションを取得する。
func (s *PostCollectionService) GetPublicByUserID(userID uint) ([]model.PostCollection, error) {
	return s.repo.FindPublicByUserID(userID)
}

// GetCollectionsForViewer は閲覧者の立場に応じたコレクション一覧を返す。
// 自分のコレクションは全件（ページネーション付き）、他人のコレクションは公開のみ返す。
func (s *PostCollectionService) GetCollectionsForViewer(viewerID, targetUserID uint, limit, offset int) ([]model.PostCollection, int64, error) {
	if viewerID == targetUserID {
		return s.repo.FindByUserID(targetUserID, limit, offset)
	}
	collections, err := s.repo.FindPublicByUserID(targetUserID)
	if err != nil {
		return nil, 0, err
	}
	return collections, int64(len(collections)), nil
}

// findAndCheckOwnership はコレクションを取得し、指定ユーザーが所有者かを検証する。
func (s *PostCollectionService) findAndCheckOwnership(id, userID uint) (*model.PostCollection, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(c *model.PostCollection) uint { return c.UserID })
}

// Update は所有権を検証した後、コレクションを更新する。
func (s *PostCollectionService) Update(id, userID uint, title, description string, isPublic bool) (*model.PostCollection, error) {
	if err := domain.ValidateStringLength(title, 1, 200, "タイトル"); err != nil {
		return nil, err
	}
	if len([]rune(strings.TrimSpace(description))) > 1000 {
		return nil, domain.NewError(domain.ErrCodeValidation, "説明は1000文字以下である必要があります", nil)
	}
	collection, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
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
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// AddPost は所有権を検証した後、コレクションに投稿を追加する。
// 同じ投稿がすでに追加されている場合はエラーを返す。
func (s *PostCollectionService) AddPost(collectionID, userID, postID uint, note string) error {
	if _, err := s.findAndCheckOwnership(collectionID, userID); err != nil {
		return err
	}

	exists, err := s.repo.HasPost(collectionID, postID)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewError(domain.ErrCodeBadRequest, "すでに追加済みの投稿です", nil)
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
	if _, err := s.findAndCheckOwnership(collectionID, userID); err != nil {
		return err
	}
	return s.repo.RemovePost(collectionID, postID)
}

// GetPosts は指定コレクションの投稿一覧を取得する。
func (s *PostCollectionService) GetPosts(collectionID uint) ([]model.PostCollectionItem, error) {
	return s.repo.GetPostsByCollectionID(collectionID)
}
