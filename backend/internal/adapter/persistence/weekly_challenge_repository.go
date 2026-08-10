package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// weeklyChallengeRepository は [repository.WeeklyChallengeRepository] の GORM 実装。
type weeklyChallengeRepository struct {
	db *gorm.DB
}

// NewWeeklyChallengeRepository は WeeklyChallengeRepository の GORM 実装を返す。
func NewWeeklyChallengeRepository(db *gorm.DB) repository.WeeklyChallengeRepository {
	return &weeklyChallengeRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.WeeklyChallengeRepository = (*weeklyChallengeRepository)(nil)

// Create はウィークリーチャレンジを作成する。
func (r *weeklyChallengeRepository) Create(ctx context.Context, challenge *model.WeeklyChallenge) error {
	return r.db.WithContext(ctx).Create(challenge).Error
}

// FindByUserAndWeek は指定ユーザーの指定週のチャレンジを取得する。
// 未登録は「不在」として (nil, nil) に正規化し、GORM のエラー型を usecase へ漏らさない。
func (r *weeklyChallengeRepository) FindByUserAndWeek(ctx context.Context, userID uint, year, week int) (*model.WeeklyChallenge, error) {
	var challenge model.WeeklyChallenge
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND year = ? AND week = ?", userID, year, week).
		First(&challenge).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

// Update はウィークリーチャレンジを更新する。
func (r *weeklyChallengeRepository) Update(ctx context.Context, challenge *model.WeeklyChallenge) error {
	return r.db.WithContext(ctx).Save(challenge).Error
}
