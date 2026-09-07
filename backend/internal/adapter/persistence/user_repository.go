package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// userSearchLimit はユーザー検索の最大件数。
const userSearchLimit = 50

// userRepository は [repository.UserRepository] の sqlc(pgx) 実装。
type userRepository struct {
	q *sqlcgen.Queries
}

// NewUserRepository は UserRepository の sqlc(pgx) 実装を返す。
func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepository{q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.UserRepository = (*userRepository)(nil)

// 最小 port としても使えることを保証する（おすすめユーザーの算出はこちらに依存する）。
var _ repository.UserSkillsReader = (*userRepository)(nil)

// 認証が必要とする操作（作成・削除・パスワード更新）も同じ実装で満たす。
var _ repository.AuthUserRepository = (*userRepository)(nil)

// NewAuthUserRepository は AuthUserRepository の sqlc(pgx) 実装を返す。
func NewAuthUserRepository(pool *pgxpool.Pool) repository.AuthUserRepository {
	return &userRepository{q: sqlcgen.New(pool)}
}

// FindAll は全ユーザーを取得する。
func (r *userRepository) FindAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.q.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]model.User, len(rows))
	for i, row := range rows {
		users[i] = toModelUser(row)
	}
	return users, nil
}

// FindByID は指定 ID のユーザーを取得する。不在の場合は (nil, nil) を返す。
func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	row, err := r.q.GetUserByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := toModelUser(row)
	return &user, nil
}

// FindByUsername は指定ユーザー名のユーザーを取得する。不在の場合は (nil, nil) を返す。
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := toModelUser(row)
	return &user, nil
}

// Search は名前またはメールアドレスへの部分一致で検索する（大文字小文字を区別しない）。
func (r *userRepository) Search(ctx context.Context, query string) ([]model.User, error) {
	pattern := escapeLikePattern(query)
	rows, err := r.q.SearchUsers(ctx, sqlcgen.SearchUsersParams{
		Name:  pattern,
		Limit: int32Param(userSearchLimit),
	})
	if err != nil {
		return nil, err
	}
	users := make([]model.User, len(rows))
	for i, row := range rows {
		users[i] = toModelUser(row)
	}
	return users, nil
}

// Update はユーザー情報を更新する（GORMのSave＝全カラム上書きに相当）。
// passwordは対象外（UpdateUserPasswordを使う。user_credentials側にありusersの列でも
// ないため、そもそも対象にできない。DEVSYNC-159）。
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	row, err := r.q.UpdateUser(ctx, sqlcgen.UpdateUserParams{
		ID:                  int64(user.ID),
		Username:            user.Username,
		Name:                user.Name,
		Email:               user.Email,
		AvatarUrl:           &user.AvatarURL,
		Bio:                 &user.Bio,
		GitHubID:            &user.GitHubID,
		GitHubUsername:      &user.GitHubUsername,
		GitHubToken:         &user.GitHubToken,
		GitHubConnected:     user.GitHubConnected,
		SpotifyConnected:    user.SpotifyConnected,
		SpotifyToken:        &user.SpotifyToken,
		SpotifyRefreshToken: &user.SpotifyRefreshToken,
		SpotifyTokenExpiry:  toTimestamptz(&user.SpotifyTokenExpiry),
		ZennUsername:        &user.ZennUsername,
		QiitaUsername:       &user.QiitaUsername,
		AtCoderUsername:     &user.AtCoderUsername,
		PaizaRank:           &user.PaizaRank,
		SkillsLanguages:     &user.SkillsLanguages,
		SkillsFrameworks:    &user.SkillsFrameworks,
		OnboardingCompleted: user.OnboardingCompleted,
		EmailWeeklyReport:   user.EmailWeeklyReport,
		EmailLanguage:       &user.EmailLanguage,
	})
	if err != nil {
		return err
	}
	*user = toModelUser(row)
	return nil
}

// FindByEmail はメールアドレスでユーザーを取得する。パスワードハッシュも含めて返す
// （ログイン処理での照合に使うため）。不在の場合は (nil, nil) を返す。
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := toModelUser(row.User)
	user.Password = fromStringPtr(row.PasswordHash)
	return &user, nil
}

// FindByGitHubID は GitHub の ID でユーザーを取得する。パスワードハッシュも含めて返す。
// 0 は「GitHub 未連携」を表す値で特定のユーザーを指さないため、常に不在として扱う。
func (r *userRepository) FindByGitHubID(ctx context.Context, githubID int64) (*model.User, error) {
	if githubID == 0 {
		return nil, nil
	}
	row, err := r.q.GetUserByGitHubID(ctx, &githubID)
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := toModelUser(row.User)
	user.Password = fromStringPtr(row.PasswordHash)
	return &user, nil
}

// FindByIDWithPassword はIDでユーザーを取得する。パスワードハッシュも含めて返す
// （退会時の本人確認に使うため）。不在の場合は (nil, nil) を返す。
func (r *userRepository) FindByIDWithPassword(ctx context.Context, id uint) (*model.User, error) {
	row, err := r.q.GetUserByIDWithPassword(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := toModelUser(row.User)
	user.Password = fromStringPtr(row.PasswordHash)
	return &user, nil
}

// Create はユーザーを作成する。
// EmailWeeklyReport/EmailLanguageはGORMの `gorm:"default:..."` に相当し、
// Goのゼロ値のときはDBデフォルト（true / "ja"）を補う。
// パスワードが空文字でなければuser_credentialsへも同一SQL文で登録する（DEVSYNC-159）。
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	emailLanguage := user.EmailLanguage
	if emailLanguage == "" {
		emailLanguage = "ja"
	}
	emailWeeklyReport := true
	password := user.Password

	row, err := r.q.CreateUser(ctx, sqlcgen.CreateUserParams{
		Username:            user.Username,
		Name:                user.Name,
		Email:               user.Email,
		AvatarUrl:           &user.AvatarURL,
		Bio:                 &user.Bio,
		GitHubID:            &user.GitHubID,
		GitHubUsername:      &user.GitHubUsername,
		GitHubToken:         &user.GitHubToken,
		GitHubConnected:     user.GitHubConnected,
		SpotifyConnected:    user.SpotifyConnected,
		SpotifyToken:        &user.SpotifyToken,
		SpotifyRefreshToken: &user.SpotifyRefreshToken,
		SpotifyTokenExpiry:  toTimestamptz(&user.SpotifyTokenExpiry),
		ZennUsername:        &user.ZennUsername,
		QiitaUsername:       &user.QiitaUsername,
		AtCoderUsername:     &user.AtCoderUsername,
		PaizaRank:           &user.PaizaRank,
		SkillsLanguages:     &user.SkillsLanguages,
		SkillsFrameworks:    &user.SkillsFrameworks,
		OnboardingCompleted: user.OnboardingCompleted,
		EmailWeeklyReport:   emailWeeklyReport,
		EmailLanguage:       &emailLanguage,
		PasswordHash:        password,
	})
	if err != nil {
		return err
	}
	*user = toModelUser(sqlcgen.User(row))
	user.Password = password
	return nil
}

// UpdatePassword はパスワードハッシュだけを更新する（user_credentialsへupsertする。
// GitHubのみで登録したユーザーが後からパスワードを設定する場合、行がまだ無いため）。
func (r *userRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	return r.q.UpdateUserPassword(ctx, sqlcgen.UpdateUserPasswordParams{
		UserID:       int64(userID),
		PasswordHash: hashedPassword,
	})
}

// DeleteWithRelatedData はユーザーを削除する。関連データの削除はFKの
// ON DELETE CASCADE宣言（internal/infra/database/schema/schema.hcl）に委ねる
// （DEVSYNC-156でFKを全テーブルへ投入済み）。かつては削除順序を守るため
// 約20本のクエリを手書きしていたが、新しいユーザー所有テーブルを追加するたびに
// ここへの追記が必要という壊れやすい運用だった。追記漏れがあってもFKが無い側では
// 検出できず孤児行が残るため、削除順序をコードで再現するのをやめ宣言に一本化する。
func (r *userRepository) DeleteWithRelatedData(ctx context.Context, id uint) error {
	return r.q.DeleteUser(ctx, int64(id))
}
