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
// DeleteWithRelatedData は多数のテーブルにまたがる削除を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type userRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewUserRepository は UserRepository の sqlc(pgx) 実装を返す。
func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.UserRepository = (*userRepository)(nil)

// 最小 port としても使えることを保証する（おすすめユーザーの算出はこちらに依存する）。
var _ repository.UserSkillsReader = (*userRepository)(nil)

// 認証が必要とする操作（作成・削除・パスワード更新）も同じ実装で満たす。
var _ repository.AuthUserRepository = (*userRepository)(nil)

// NewAuthUserRepository は AuthUserRepository の sqlc(pgx) 実装を返す。
func NewAuthUserRepository(pool *pgxpool.Pool) repository.AuthUserRepository {
	return &userRepository{pool: pool, q: sqlcgen.New(pool)}
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
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	row, err := r.q.UpdateUser(ctx, sqlcgen.UpdateUserParams{
		ID:                  int64(user.ID),
		Username:            user.Username,
		Name:                user.Name,
		Email:               user.Email,
		Password:            &user.Password,
		AvatarUrl:           &user.AvatarURL,
		Bio:                 &user.Bio,
		GitHubID:            &user.GitHubID,
		GitHubUsername:      &user.GitHubUsername,
		GitHubToken:         &user.GitHubToken,
		GitHubConnected:     &user.GitHubConnected,
		SpotifyConnected:    &user.SpotifyConnected,
		SpotifyToken:        &user.SpotifyToken,
		SpotifyRefreshToken: &user.SpotifyRefreshToken,
		SpotifyTokenExpiry:  toTimestamptz(&user.SpotifyTokenExpiry),
		ZennUsername:        &user.ZennUsername,
		QiitaUsername:       &user.QiitaUsername,
		AtCoderUsername:     &user.AtCoderUsername,
		PaizaRank:           &user.PaizaRank,
		SkillsLanguages:     &user.SkillsLanguages,
		SkillsFrameworks:    &user.SkillsFrameworks,
		OnboardingCompleted: &user.OnboardingCompleted,
		EmailWeeklyReport:   &user.EmailWeeklyReport,
		EmailLanguage:       &user.EmailLanguage,
	})
	if err != nil {
		return err
	}
	*user = toModelUser(row)
	return nil
}

// FindByEmail はメールアドレスでユーザーを取得する。不在の場合は (nil, nil) を返す。
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := toModelUser(row)
	return &user, nil
}

// FindByGitHubID は GitHub の ID でユーザーを取得する。不在の場合は (nil, nil) を返す。
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
	user := toModelUser(row)
	return &user, nil
}

// Create はユーザーを作成する。
// EmailWeeklyReport/EmailLanguageはGORMの `gorm:"default:..."` に相当し、
// Goのゼロ値のときはDBデフォルト（true / "ja"）を補う。
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	emailLanguage := user.EmailLanguage
	if emailLanguage == "" {
		emailLanguage = "ja"
	}
	emailWeeklyReport := true

	row, err := r.q.CreateUser(ctx, sqlcgen.CreateUserParams{
		Username:            user.Username,
		Name:                user.Name,
		Email:               user.Email,
		Password:            &user.Password,
		AvatarUrl:           &user.AvatarURL,
		Bio:                 &user.Bio,
		GitHubID:            &user.GitHubID,
		GitHubUsername:      &user.GitHubUsername,
		GitHubToken:         &user.GitHubToken,
		GitHubConnected:     &user.GitHubConnected,
		SpotifyConnected:    &user.SpotifyConnected,
		SpotifyToken:        &user.SpotifyToken,
		SpotifyRefreshToken: &user.SpotifyRefreshToken,
		SpotifyTokenExpiry:  toTimestamptz(&user.SpotifyTokenExpiry),
		ZennUsername:        &user.ZennUsername,
		QiitaUsername:       &user.QiitaUsername,
		AtCoderUsername:     &user.AtCoderUsername,
		PaizaRank:           &user.PaizaRank,
		SkillsLanguages:     &user.SkillsLanguages,
		SkillsFrameworks:    &user.SkillsFrameworks,
		OnboardingCompleted: &user.OnboardingCompleted,
		EmailWeeklyReport:   &emailWeeklyReport,
		EmailLanguage:       &emailLanguage,
	})
	if err != nil {
		return err
	}
	*user = toModelUser(row)
	return nil
}

// UpdatePassword はパスワードハッシュだけを更新する。
func (r *userRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	return r.q.UpdateUserPassword(ctx, sqlcgen.UpdateUserPasswordParams{
		ID:       int64(userID),
		Password: &hashedPassword,
	})
}

// DeleteWithRelatedData はユーザーと関連データをトランザクション内で削除する。
// 退会ユーザーの投稿に紐づく従属データ（他ユーザーのコメント・いいね等を含む）→
// 通知・メッセージ・本人のコメント・いいね・投稿・フォロー・GitHub 連携データ・
// パスワードリセットトークンの順に削除してから、最後にユーザー本体を削除する。
func (r *userRepository) DeleteWithRelatedData(ctx context.Context, id uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	userID := int64(id)

	// 退会ユーザーの投稿に紐づく従属データを先に消す。user_id 条件だけでは
	// 他ユーザーが付けたコメント・いいね等が削除済み投稿を参照したまま残る。
	// 対象テーブルは postRepository.Delete と同じ集合を退会者の全投稿へ広げたもの。

	// コメントに従属する行（コメントいいね・コメント由来のメンション）を先に消す。
	// 投稿配下のコメントに加え、退会者が他の投稿に書いたコメントの従属行も対象にする。
	if err := q.DeleteCommentLikesByUserCommentsSet(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteMentionsByUserCommentsSet(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteCommentsByUserPostsAndSelf(ctx, userID); err != nil {
		return err
	}

	// スニペットに従属する行を先に消す
	if err := q.DeleteSnippetCommentsByUserPostSnippets(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteCodeSnippetsByUserPosts(ctx, userID); err != nil {
		return err
	}

	// 投稿を直接参照する行を消す（本人分は user_id でも消し、他ユーザー分は post_id で消す）
	if err := q.DeleteLikesByUserPosts(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteReactionsByUserPosts(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteBookmarksByUserPosts(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteBookmarkCollectionItemsByUserPosts(ctx, userID); err != nil {
		return err
	}
	if err := q.DeletePostSeriesItemsByUserPosts(ctx, userID); err != nil {
		return err
	}
	if err := q.DeletePostCollectionItemsByUserPosts(ctx, userID); err != nil {
		return err
	}
	if err := q.DeletePostTagsByUserPosts(ctx, userID); err != nil {
		return err
	}
	if err := q.DeletePostPinsByUserPosts(ctx, userID); err != nil {
		return err
	}
	if err := q.DeletePostViewsByUserPosts(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteMentionsByUserPosts(ctx, userID); err != nil {
		return err
	}

	if err := q.DeleteNotificationsForUserDeletion(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteMessagesByUser(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteLikesByUser(ctx, userID); err != nil {
		return err
	}
	if err := q.DeletePostsByUser(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteFollowsByUser(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteGitHubContributionsByUser(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteGitHubLanguageStatsByUser(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteGitHubReposByUser(ctx, userID); err != nil {
		return err
	}
	if err := q.DeletePasswordResetTokensByUser(ctx, userID); err != nil {
		return err
	}
	if err := q.DeleteUser(ctx, userID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
