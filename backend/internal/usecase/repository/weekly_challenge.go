package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// WeeklyChallengeRepository はウィークリーチャレンジの永続化に対する、usecase 側が要求する契約。
type WeeklyChallengeRepository interface {
	Create(ctx context.Context, challenge *model.WeeklyChallenge) error
	// FindByUserAndWeek は該当週のチャレンジを返す。
	// 未登録の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	// usecase 側が永続化技術のエラー型（gorm.ErrRecordNotFound 等）を知らずに済むようにするための契約。
	FindByUserAndWeek(ctx context.Context, userID uint, year, week int) (*model.WeeklyChallenge, error)
	Update(ctx context.Context, challenge *model.WeeklyChallenge) error
}
