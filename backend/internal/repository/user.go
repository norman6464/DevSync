package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Search(query string) ([]model.User, error) {
	var users []model.User
	result := r.db.Where("name ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").Limit(50).Find(&users)
	return users, result.Error
}

func (r *UserRepository) FindByGitHubID(githubID int64) (*model.User, error) {
	var user model.User
	result := r.db.Where("github_id = ?", githubID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// DeleteWithRelatedData deletes a user and all their related data
func (r *UserRepository) DeleteWithRelatedData(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete notifications (where user is recipient or actor)
		if err := tx.Where("user_id = ? OR actor_id = ?", id, id).Delete(&model.Notification{}).Error; err != nil {
			return err
		}

		// Delete messages (sent or received)
		if err := tx.Where("sender_id = ? OR receiver_id = ?", id, id).Delete(&model.Message{}).Error; err != nil {
			return err
		}

		// Delete comments
		if err := tx.Where("user_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
			return err
		}

		// Delete likes
		if err := tx.Where("user_id = ?", id).Delete(&model.Like{}).Error; err != nil {
			return err
		}

		// Delete posts
		if err := tx.Where("user_id = ?", id).Delete(&model.Post{}).Error; err != nil {
			return err
		}

		// Delete follows (both directions)
		if err := tx.Where("follower_id = ? OR followee_id = ?", id, id).Delete(&model.Follow{}).Error; err != nil {
			return err
		}

		// Delete GitHub data
		if err := tx.Where("user_id = ?", id).Delete(&model.GitHubContribution{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.GitHubLanguageStat{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.GitHubRepository{}).Error; err != nil {
			return err
		}

		// Delete password reset tokens
		if err := tx.Where("user_id = ?", id).Delete(&model.PasswordResetToken{}).Error; err != nil {
			return err
		}

		// Finally delete the user
		if err := tx.Delete(&model.User{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}
