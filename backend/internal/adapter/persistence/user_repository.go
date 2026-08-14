package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// userSearchLimit はユーザー検索の最大件数。
const userSearchLimit = 50

// userRepository は [repository.UserRepository] の GORM 実装。
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository は UserRepository の GORM 実装を返す。
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.UserRepository = (*userRepository)(nil)

// 最小 port としても使えることを保証する（おすすめユーザーの算出はこちらに依存する）。
var _ repository.UserSkillsReader = (*userRepository)(nil)

// 認証が必要とする操作（作成・削除・パスワード更新）も同じ実装で満たす。
var _ repository.AuthUserRepository = (*userRepository)(nil)

// NewAuthUserRepository は AuthUserRepository の GORM 実装を返す。
func NewAuthUserRepository(db *gorm.DB) repository.AuthUserRepository {
	return &userRepository{db: db}
}

// FindAll は全ユーザーを取得する。
func (r *userRepository) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

// FindByID は指定 ID のユーザーを取得する。不在の場合は (nil, nil) を返す。
func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	return firstUser(r.db.WithContext(ctx).First(&user, id), &user)
}

// FindByUsername は指定ユーザー名のユーザーを取得する。不在の場合は (nil, nil) を返す。
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	return firstUser(r.db.WithContext(ctx).Where("username = ?", username).First(&user), &user)
}

// firstUser は 1 件取得の結果を「不在は (nil, nil)」の契約へ変換する。
func firstUser(tx *gorm.DB, user *model.User) (*model.User, error) {
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

// Search は名前またはメールアドレスへの部分一致で検索する（大文字小文字を区別しない）。
func (r *userRepository) Search(ctx context.Context, query string) ([]model.User, error) {
	pattern := escapeLikePattern(query)
	var users []model.User
	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR email ILIKE ?", pattern, pattern).
		Limit(userSearchLimit).
		Find(&users).Error
	return users, err
}

// Update はユーザー情報を更新する。
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// FindByEmail はメールアドレスでユーザーを取得する。不在の場合は (nil, nil) を返す。
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	return firstUser(r.db.WithContext(ctx).Where("email = ?", email).First(&user), &user)
}

// FindByGitHubID は GitHub の ID でユーザーを取得する。不在の場合は (nil, nil) を返す。
// 0 は「GitHub 未連携」を表す値で特定のユーザーを指さないため、常に不在として扱う。
func (r *userRepository) FindByGitHubID(ctx context.Context, githubID int64) (*model.User, error) {
	if githubID == 0 {
		return nil, nil
	}
	var user model.User
	return firstUser(r.db.WithContext(ctx).Where("git_hub_id = ?", githubID).First(&user), &user)
}

// Create はユーザーを作成する。
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// UpdatePassword はパスワードハッシュだけを更新する。
func (r *userRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).Update("password", hashedPassword).Error
}

// DeleteWithRelatedData はユーザーと関連データをトランザクション内で削除する。
// 通知・メッセージ・コメント・いいね・投稿・フォロー・GitHub 連携データ・
// パスワードリセットトークンを削除してから、最後にユーザー本体を削除する。
func (r *userRepository) DeleteWithRelatedData(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? OR actor_id = ?", id, id).Delete(&model.Notification{}).Error; err != nil {
			return err
		}
		if err := tx.Where("sender_id = ? OR receiver_id = ?", id, id).Delete(&model.Message{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.Like{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.Post{}).Error; err != nil {
			return err
		}
		if err := tx.Where("follower_id = ? OR followee_id = ?", id, id).Delete(&model.Follow{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.GitHubContribution{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.GitHubLanguageStat{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.GitHubRepository{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.PasswordResetToken{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.User{}, id).Error
	})
}
