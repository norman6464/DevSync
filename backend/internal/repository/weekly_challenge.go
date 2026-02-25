package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// WeeklyChallengeRepository はウィークリーチャレンジのリポジトリ実装。
type WeeklyChallengeRepository struct {
	db *gorm.DB
}

// NewWeeklyChallengeRepository は新しいWeeklyChallengeRepositoryを生成する。
func NewWeeklyChallengeRepository(db *gorm.DB) *WeeklyChallengeRepository {
	return &WeeklyChallengeRepository{db: db}
}

// Create はウィークリーチャレンジを作成する。
func (r *WeeklyChallengeRepository) Create(challenge *model.WeeklyChallenge) error {
	return r.db.Create(challenge).Error
}

// FindByUserAndWeek は指定ユーザーの指定週のチャレンジを取得する。
func (r *WeeklyChallengeRepository) FindByUserAndWeek(userID uint, year, week int) (*model.WeeklyChallenge, error) {
	var challenge model.WeeklyChallenge
	err := r.db.Where("user_id = ? AND year = ? AND week = ?", userID, year, week).First(&challenge).Error
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

// Update はウィークリーチャレンジを更新する。
func (r *WeeklyChallengeRepository) Update(challenge *model.WeeklyChallenge) error {
	return r.db.Save(challenge).Error
}
