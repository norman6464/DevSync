package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// errAIAdviceNotFound は既読にする対象のアドバイスが無いときに返すエラー。
var errAIAdviceNotFound = errors.New("アドバイスが見つかりません")

// aiAdviceRepository は [repository.AIAdviceRepository] の GORM 実装。
type aiAdviceRepository struct {
	db *gorm.DB
}

// NewAIAdviceRepository は AIAdviceRepository の GORM 実装を返す。
func NewAIAdviceRepository(db *gorm.DB) repository.AIAdviceRepository {
	return &aiAdviceRepository{db: db}
}

var _ repository.AIAdviceRepository = (*aiAdviceRepository)(nil)

// CreateBatch は複数のアドバイスを一括作成する。
func (r *aiAdviceRepository) CreateBatch(ctx context.Context, advices []*model.AIAdvice) error {
	if len(advices) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&advices).Error
}

// FindByUserID は優先度の高い順・作成の新しい順にアドバイスを返す。
func (r *aiAdviceRepository) FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AIAdvice, error) {
	var advices []model.AIAdvice
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("priority ASC, created_at DESC").
		Limit(limit).
		Find(&advices).Error
	return advices, err
}

// FindUnreadByUserID は未読のアドバイスを優先度順で返す。
func (r *aiAdviceRepository) FindUnreadByUserID(ctx context.Context, userID uint) ([]model.AIAdvice, error) {
	var advices []model.AIAdvice
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_read = ?", userID, false).
		Order("priority ASC, created_at DESC").
		Find(&advices).Error
	return advices, err
}

// MarkAsRead は本人のアドバイス 1 件を既読にする。対象が無ければエラーを返す。
func (r *aiAdviceRepository) MarkAsRead(ctx context.Context, id, userID uint) error {
	result := r.db.WithContext(ctx).Model(&model.AIAdvice{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errAIAdviceNotFound
	}
	return nil
}

// DeleteByUserID は指定ユーザーのアドバイスをすべて削除する。
func (r *aiAdviceRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.AIAdvice{}).Error
}
