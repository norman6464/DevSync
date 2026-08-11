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
