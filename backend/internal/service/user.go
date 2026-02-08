package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// UserService はユーザー情報管理のビジネスロジックを提供する。
type UserService struct {
	repo repository.UserRepositoryInterface
}

// NewUserService は新しいUserServiceインスタンスを生成する。
func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

// GetAll は全ユーザーを取得する。検索クエリが指定された場合はフィルタリングする。
func (s *UserService) GetAll(query string) ([]model.User, error) {
	if query != "" {
		return s.repo.Search(query)
	}
	return s.repo.FindAll()
}

// GetByID は指定IDのユーザーを取得する。
func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

// FindByID は指定IDのユーザーを取得する（リポジトリ互換エイリアス）。
func (s *UserService) FindByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

// Update はユーザー情報を更新する。
func (s *UserService) Update(user *model.User) error {
	return s.repo.Update(user)
}
