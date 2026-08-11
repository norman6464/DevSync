package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// msgUserNotFound はユーザーを取得できなかったときのメッセージ。
const msgUserNotFound = "ユーザーが見つかりません"

// errProfileUserNotFound はプロフィール完成度の算出でユーザーが不在だったときに返すエラー。
// DomainError ではないため handler では 500 になり、取得エラーをそのまま返していた
// 移行前の挙動と一致する。
var errProfileUserNotFound = errors.New(msgUserNotFound)

// profileCompletenessFields はプロフィール完成度の評価項目数。
const profileCompletenessFields = 4

// validEmailLanguages はメール配信で受け付ける言語コード。
var validEmailLanguages = map[string]bool{
	"ja": true, "en": true, "ko": true, "zh-CN": true, "zh-TW": true,
	"es": true, "fr": true, "de": true, "pt": true, "ru": true,
}

// findUserOr404 はユーザーを取得する。取得できない場合は 404 を返す
// （不在も DB 障害も 404 に潰す移行前の挙動を維持している）。
func findUserOr404(ctx context.Context, users repository.UserSkillsReader, id uint) (*model.User, error) {
	user, err := users.FindByID(ctx, id)
	if err != nil || user == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, msgUserNotFound, err)
	}
	return user, nil
}

// validateUserProfile はプロフィール各項目の文字数を検証する。
func validateUserProfile(user *model.User) error {
	checks := []struct {
		value    string
		min, max int
		label    string
	}{
		{user.Name, 1, 100, "名前"},
		{user.Bio, 0, 500, "自己紹介"},
		{user.SkillsLanguages, 0, 500, "プログラミング言語スキル"},
		{user.SkillsFrameworks, 0, 500, "フレームワークスキル"},
		{user.AvatarURL, 0, 2000, "アバターURL"},
	}
	for _, c := range checks {
		if err := domain.ValidateStringLength(c.value, c.min, c.max, c.label); err != nil {
			return err
		}
	}
	return nil
}

// ListUsersUseCase はユーザー一覧を取得する。
type ListUsersUseCase struct {
	users repository.UserRepository
}

// NewListUsersUseCase は ListUsersUseCase を生成する。
func NewListUsersUseCase(users repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{users: users}
}

// Execute は検索キーワードがあれば検索し、無ければ全件を返す。
func (uc *ListUsersUseCase) Execute(ctx context.Context, query string) ([]model.User, error) {
	if query != "" {
		return uc.users.Search(ctx, query)
	}
	return uc.users.FindAll(ctx)
}

// GetUserUseCase は指定 ID のユーザーを取得する。
type GetUserUseCase struct {
	users repository.UserRepository
}

// NewGetUserUseCase は GetUserUseCase を生成する。
func NewGetUserUseCase(users repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{users: users}
}

// Execute はユーザーを返す。取得できない場合は 404。
func (uc *GetUserUseCase) Execute(ctx context.Context, id uint) (*model.User, error) {
	return findUserOr404(ctx, uc.users, id)
}

// GetUserByUsernameUseCase は指定ユーザー名のユーザーを取得する。
type GetUserByUsernameUseCase struct {
	users repository.UserRepository
}

// NewGetUserByUsernameUseCase は GetUserByUsernameUseCase を生成する。
func NewGetUserByUsernameUseCase(users repository.UserRepository) *GetUserByUsernameUseCase {
	return &GetUserByUsernameUseCase{users: users}
}

// Execute はユーザーを返す。取得できない場合は 404。
func (uc *GetUserByUsernameUseCase) Execute(ctx context.Context, username string) (*model.User, error) {
	user, err := uc.users.FindByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, msgUserNotFound, err)
	}
	return user, nil
}

// UpdateProfileInput はプロフィール更新の入力。
// ポインタのフィールドは nil なら「変更なし」を表す。
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

// UpdateUserProfileUseCase は本人のプロフィールを更新する。
type UpdateUserProfileUseCase struct {
	users repository.UserRepository
}

// NewUpdateUserProfileUseCase は UpdateUserProfileUseCase を生成する。
func NewUpdateUserProfileUseCase(users repository.UserRepository) *UpdateUserProfileUseCase {
	return &UpdateUserProfileUseCase{users: users}
}

// Execute は本人であることを確認し、渡された項目でプロフィールを更新する。
// 名前は空でなければ更新し、自己紹介とアバター URL はトリム後の値で必ず上書きする。
func (uc *UpdateUserProfileUseCase) Execute(ctx context.Context, id, userID uint, in *UpdateProfileInput) (*model.User, error) {
	if id != userID {
		return nil, domain.ErrForbidden
	}

	user, err := findUserOr404(ctx, uc.users, id)
	if err != nil {
		return nil, err
	}

	if name := strings.TrimSpace(in.Name); name != "" {
		user.Name = name
	}
	user.Bio = strings.TrimSpace(in.Bio)
	user.AvatarURL = strings.TrimSpace(in.AvatarURL)

	optionals := []struct {
		value  *string
		target *string
	}{
		{in.SkillsLanguages, &user.SkillsLanguages},
		{in.SkillsFrameworks, &user.SkillsFrameworks},
		{in.AtCoderUsername, &user.AtCoderUsername},
		{in.PaizaRank, &user.PaizaRank},
	}
	for _, o := range optionals {
		if o.value != nil {
			*o.target = strings.TrimSpace(*o.value)
		}
	}
	if in.OnboardingCompleted != nil {
		user.OnboardingCompleted = *in.OnboardingCompleted
	}

	if err := validateUserProfile(user); err != nil {
		return nil, err
	}
	if err := uc.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// ProfileCompleteness はプロフィール完成度の計算結果。
type ProfileCompleteness struct {
	Percentage    int      `json:"percentage"`
	MissingFields []string `json:"missing_fields"`
}

// GetProfileCompletenessUseCase はプロフィール完成度を取得する。
type GetProfileCompletenessUseCase struct {
	users repository.UserRepository
}

// NewGetProfileCompletenessUseCase は GetProfileCompletenessUseCase を生成する。
func NewGetProfileCompletenessUseCase(users repository.UserRepository) *GetProfileCompletenessUseCase {
	return &GetProfileCompletenessUseCase{users: users}
}

// Execute は完成度を返す。ユーザーを取得できない場合は DomainError ではないエラーを返す
// （他の取得と違い 404 に変換せず 500 になる移行前の挙動を維持している）。
func (uc *GetProfileCompletenessUseCase) Execute(ctx context.Context, userID uint) (*ProfileCompleteness, error) {
	user, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errProfileUserNotFound
	}
	return CalculateProfileCompleteness(user), nil
}

// CalculateProfileCompleteness はアバター・自己紹介・GitHub 連携・スキルの 4 項目から
// 完成度を算出する純粋関数。未設定の項目名も返す。
func CalculateProfileCompleteness(user *model.User) *ProfileCompleteness {
	items := []struct {
		filled bool
		name   string
	}{
		{user.AvatarURL != "", "avatar"},
		{user.Bio != "", "bio"},
		{user.GitHubConnected, "github"},
		{user.SkillsLanguages != "" || user.SkillsFrameworks != "", "skills"},
	}

	completed := 0
	var missing []string
	for _, item := range items {
		if item.filled {
			completed++
			continue
		}
		missing = append(missing, item.name)
	}

	return &ProfileCompleteness{
		Percentage:    (completed * 100) / profileCompletenessFields,
		MissingFields: missing,
	}
}

// GetEmailPreferencesUseCase はメール配信設定を取得する。
type GetEmailPreferencesUseCase struct {
	users repository.UserRepository
}

// NewGetEmailPreferencesUseCase は GetEmailPreferencesUseCase を生成する。
func NewGetEmailPreferencesUseCase(users repository.UserRepository) *GetEmailPreferencesUseCase {
	return &GetEmailPreferencesUseCase{users: users}
}

// Execute は設定を保持するユーザーを返す。取得できない場合は 404。
func (uc *GetEmailPreferencesUseCase) Execute(ctx context.Context, userID uint) (*model.User, error) {
	return findUserOr404(ctx, uc.users, userID)
}

// UpdateEmailPreferencesUseCase はメール配信設定を更新する。
type UpdateEmailPreferencesUseCase struct {
	users repository.UserRepository
}

// NewUpdateEmailPreferencesUseCase は UpdateEmailPreferencesUseCase を生成する。
func NewUpdateEmailPreferencesUseCase(users repository.UserRepository) *UpdateEmailPreferencesUseCase {
	return &UpdateEmailPreferencesUseCase{users: users}
}

// Execute は週次レポートの受信可否と言語を更新する。言語は許可リスト以外なら 400。
func (uc *UpdateEmailPreferencesUseCase) Execute(ctx context.Context, userID uint, weeklyReport *bool, language *string) (*model.User, error) {
	user, err := findUserOr404(ctx, uc.users, userID)
	if err != nil {
		return nil, err
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

	if err := uc.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
