package service

import (
	"strings"

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
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "ユーザーが見つかりません", err)
	}
	return user, nil
}

// GetByUsername は指定ユーザー名のユーザーを取得する。
func (s *UserService) GetByUsername(username string) (*model.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "ユーザーが見つかりません", err)
	}
	return user, nil
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
	if err := domain.ValidateStringLength(user.Bio, 0, 500, "自己紹介"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(user.SkillsLanguages, 0, 500, "プログラミング言語スキル"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(user.SkillsFrameworks, 0, 500, "フレームワークスキル"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(user.AvatarURL, 0, 2000, "アバターURL"); err != nil {
		return err
	}
	return s.repo.Update(user)
}

// UpdateProfileInput はプロフィール更新の入力パラメータ。
// Handler層のDTOからService層への橋渡しとして使用する。
type UpdateProfileInput struct {
	Name                string
	Bio                 string
	AvatarURL           string
	SkillsLanguages     *string
	SkillsFrameworks    *string
	AtCoderUsername     *string
	PaizaRank           *string
	OnboardingCompleted *bool
}

// UpdateProfile は所有権チェック・フィールドマッピング・バリデーション・保存を一括処理する。
func (s *UserService) UpdateProfile(id, userID uint, input *UpdateProfileInput) (*model.User, error) {
	if id != userID {
		return nil, domain.ErrForbidden
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "ユーザーが見つかりません", err)
	}

	if name := strings.TrimSpace(input.Name); name != "" {
		existing.Name = name
	}
	existing.Bio = strings.TrimSpace(input.Bio)
	existing.AvatarURL = strings.TrimSpace(input.AvatarURL)
	if input.SkillsLanguages != nil {
		existing.SkillsLanguages = strings.TrimSpace(*input.SkillsLanguages)
	}
	if input.SkillsFrameworks != nil {
		existing.SkillsFrameworks = strings.TrimSpace(*input.SkillsFrameworks)
	}
	if input.AtCoderUsername != nil {
		existing.AtCoderUsername = strings.TrimSpace(*input.AtCoderUsername)
	}
	if input.PaizaRank != nil {
		existing.PaizaRank = strings.TrimSpace(*input.PaizaRank)
	}
	if input.OnboardingCompleted != nil {
		existing.OnboardingCompleted = *input.OnboardingCompleted
	}

	if err := s.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// validEmailLanguages はメール配信で有効な言語コードの集合。
var validEmailLanguages = map[string]bool{
	"ja": true, "en": true, "ko": true, "zh-CN": true, "zh-TW": true,
	"es": true, "fr": true, "de": true, "pt": true, "ru": true,
}

// UpdateEmailPreferences はメール配信設定を更新する。
// 言語コードのバリデーションを含む。
func (s *UserService) UpdateEmailPreferences(userID uint, weeklyReport *bool, language *string) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "ユーザーが見つかりません", err)
	}

	if weeklyReport != nil {
		user.EmailWeeklyReport = *weeklyReport
	}
	if language != nil {
		if !validEmailLanguages[*language] {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なメール言語設定です", nil)
		}
		user.EmailLanguage = *language
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
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
