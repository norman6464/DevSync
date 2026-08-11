package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// userSkillsReader は [repository.UserSkillsReader] の GORM 実装。
type userSkillsReader struct {
	db *gorm.DB
}

// NewUserSkillsReader は UserSkillsReader の GORM 実装を返す。
func NewUserSkillsReader(db *gorm.DB) repository.UserSkillsReader {
	return &userSkillsReader{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.UserSkillsReader = (*userSkillsReader)(nil)

// FindByID は指定 ID のユーザーを取得する。不在の場合は (nil, nil) を返す。
func (r *userSkillsReader) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
