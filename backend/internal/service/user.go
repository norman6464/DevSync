package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
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

// GetByUsername は指定ユーザー名のユーザーを取得する。
func (s *UserService) GetByUsername(username string) (*model.User, error) {
	return s.repo.FindByUsername(username)
}

// FindByID は指定IDのユーザーを取得する（リポジトリ互換エイリアス）。
func (s *UserService) FindByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

// UpdateByOwner は所有権を検証した後、ユーザー情報を更新する。
func (s *UserService) UpdateByOwner(id, userID uint, user *model.User) error {
	if id != userID {
		return domain.ErrForbidden
	}
	return s.Update(user)
}

// Update はユーザー情報を更新する。
func (s *UserService) Update(user *model.User) error {
	if err := domain.ValidateStringLength(user.Name, 1, 100, "名前"); err != nil {
		return err
	}
	if len(user.Bio) > 500 {
		return domain.NewError(domain.ErrCodeValidation, "自己紹介は500文字以下である必要があります", nil)
	}
	if len(user.SkillsLanguages) > 500 {
		return domain.NewError(domain.ErrCodeValidation, "プログラミング言語スキルは500文字以下である必要があります", nil)
	}
	if len(user.SkillsFrameworks) > 500 {
		return domain.NewError(domain.ErrCodeValidation, "フレームワークスキルは500文字以下である必要があります", nil)
	}
	if len(user.AvatarURL) > 2000 {
		return domain.NewError(domain.ErrCodeValidation, "アバターURLは2000文字以下である必要があります", nil)
	}
	return s.repo.Update(user)
}

// ProfileCompleteness はプロフィール完成度の計算結果を表す。
type ProfileCompleteness struct {
	Percentage    int      `json:"percentage"`
	MissingFields []string `json:"missing_fields"`
}

// GetProfileCompleteness は指定ユーザーのプロフィール完成度を計算する。
func (s *UserService) GetProfileCompleteness(userID uint) (*ProfileCompleteness, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	totalFields := 4
	completed := 0
	var missing []string

	if user.AvatarURL != "" {
		completed++
	} else {
		missing = append(missing, "avatar")
	}

	if user.Bio != "" {
		completed++
	} else {
		missing = append(missing, "bio")
	}

	if user.GitHubConnected {
		completed++
	} else {
		missing = append(missing, "github")
	}

	if user.SkillsLanguages != "" || user.SkillsFrameworks != "" {
		completed++
	} else {
		missing = append(missing, "skills")
	}

	percentage := (completed * 100) / totalFields
	return &ProfileCompleteness{
		Percentage:    percentage,
		MissingFields: missing,
	}, nil
}
