package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// UserRepository はユーザーデータへのアクセスを提供するリポジトリ実装。
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository は新しいUserRepositoryインスタンスを生成する。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAll は全ユーザーを取得する。
func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	result := r.db.Find(&users)
	return users, result.Error
}

// FindByID は指定IDのユーザーを取得する。
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail は指定メールアドレスのユーザーを取得する。
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByUsername は指定ユーザー名のユーザーを取得する。
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Search は名前またはメールアドレスでユーザーを検索する（最大50件）。
func (r *UserRepository) Search(query string) ([]model.User, error) {
	var users []model.User
	searchPattern := EscapeLikePattern(query)
	result := r.db.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern).Limit(50).Find(&users)
	return users, result.Error
}

// FindByGitHubID はGitHub IDでユーザーを検索する。
func (r *UserRepository) FindByGitHubID(githubID int64) (*model.User, error) {
	var user model.User
	result := r.db.Where("git_hub_id = ?", githubID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create は新しいユーザーをデータベースに作成する。
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update は既存のユーザー情報を更新する。
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete は指定IDのユーザーを削除する。
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// DeleteWithRelatedData はユーザーと全ての関連データをトランザクション内で削除する。
// 通知、メッセージ、コメント、いいね、投稿、フォロー、GitHub連携データ、
// パスワードリセットトークンを順に削除してからユーザー本体を削除する。
func (r *UserRepository) DeleteWithRelatedData(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 通知を削除（受信者または実行者として）
		if err := tx.Where("user_id = ? OR actor_id = ?", id, id).Delete(&model.Notification{}).Error; err != nil {
			return err
		}

		// メッセージを削除（送信・受信の両方）
		if err := tx.Where("sender_id = ? OR receiver_id = ?", id, id).Delete(&model.Message{}).Error; err != nil {
			return err
		}

		// コメントを削除
		if err := tx.Where("user_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
			return err
		}

		// いいねを削除
		if err := tx.Where("user_id = ?", id).Delete(&model.Like{}).Error; err != nil {
			return err
		}

		// 投稿を削除
		if err := tx.Where("user_id = ?", id).Delete(&model.Post{}).Error; err != nil {
			return err
		}

		// フォロー関係を削除（両方向）
		if err := tx.Where("follower_id = ? OR followee_id = ?", id, id).Delete(&model.Follow{}).Error; err != nil {
			return err
		}

		// GitHub連携データを削除
		if err := tx.Where("user_id = ?", id).Delete(&model.GitHubContribution{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.GitHubLanguageStat{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.GitHubRepository{}).Error; err != nil {
			return err
		}

		// パスワードリセットトークンを削除
		if err := tx.Where("user_id = ?", id).Delete(&model.PasswordResetToken{}).Error; err != nil {
			return err
		}

		// 最後にユーザー本体を削除
		if err := tx.Delete(&model.User{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdatePassword は指定ユーザーのパスワードハッシュを更新する。
func (r *UserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}
