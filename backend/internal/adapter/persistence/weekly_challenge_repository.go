package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// weeklyChallengeRepository は [repository.WeeklyChallengeRepository] の sqlc(pgx) 実装。
type weeklyChallengeRepository struct {
	q *sqlcgen.Queries
}

// NewWeeklyChallengeRepository は WeeklyChallengeRepository の sqlc(pgx) 実装を返す。
func NewWeeklyChallengeRepository(q *sqlcgen.Queries) repository.WeeklyChallengeRepository {
	return &weeklyChallengeRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.WeeklyChallengeRepository = (*weeklyChallengeRepository)(nil)

// toModelWeeklyChallenge は sqlc の生成行を model.WeeklyChallenge へ変換する。
func toModelWeeklyChallenge(row sqlcgen.WeeklyChallenge) model.WeeklyChallenge {
	return model.WeeklyChallenge{
		ID:            uint(row.ID),
		UserID:        uint(row.UserID),
		Year:          int(row.Year),
		Week:          int(row.Week),
		ChallengeType: model.ChallengeType(row.ChallengeType),
		Description:   row.Description,
		TargetValue:   int(row.TargetValue),
		CurrentValue:  fromInt64Ptr(row.CurrentValue),
		IsCompleted:   fromBoolPtr(row.IsCompleted),
		CompletedAt:   fromTimestamptz(row.CompletedAt),
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}

// Create はウィークリーチャレンジを作成する。
func (r *weeklyChallengeRepository) Create(ctx context.Context, challenge *model.WeeklyChallenge) error {
	row, err := r.q.CreateWeeklyChallenge(ctx, sqlcgen.CreateWeeklyChallengeParams{
		UserID:        int64(challenge.UserID),
		Year:          int64(challenge.Year),
		Week:          int64(challenge.Week),
		ChallengeType: string(challenge.ChallengeType),
		Description:   challenge.Description,
		TargetValue:   int64(challenge.TargetValue),
		CurrentValue:  toInt64Ptr(challenge.CurrentValue),
		IsCompleted:   &challenge.IsCompleted,
		CompletedAt:   toTimestamptz(challenge.CompletedAt),
	})
	if err != nil {
		return err
	}
	*challenge = toModelWeeklyChallenge(row)
	return nil
}

// FindByUserAndWeek は指定ユーザーの指定週のチャレンジを取得する。
// 未登録は「不在」として (nil, nil) に正規化し、pgx のエラー型を usecase へ漏らさない。
func (r *weeklyChallengeRepository) FindByUserAndWeek(ctx context.Context, userID uint, year, week int) (*model.WeeklyChallenge, error) {
	row, err := r.q.GetWeeklyChallengeByUserAndWeek(ctx, sqlcgen.GetWeeklyChallengeByUserAndWeekParams{
		UserID: int64(userID),
		Year:   int64(year),
		Week:   int64(week),
	})
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	challenge := toModelWeeklyChallenge(row)
	return &challenge, nil
}

// Update はウィークリーチャレンジを更新する。
func (r *weeklyChallengeRepository) Update(ctx context.Context, challenge *model.WeeklyChallenge) error {
	row, err := r.q.UpdateWeeklyChallenge(ctx, sqlcgen.UpdateWeeklyChallengeParams{
		ID:            int64(challenge.ID),
		UserID:        int64(challenge.UserID),
		Year:          int64(challenge.Year),
		Week:          int64(challenge.Week),
		ChallengeType: string(challenge.ChallengeType),
		Description:   challenge.Description,
		TargetValue:   int64(challenge.TargetValue),
		CurrentValue:  toInt64Ptr(challenge.CurrentValue),
		IsCompleted:   &challenge.IsCompleted,
		CompletedAt:   toTimestamptz(challenge.CompletedAt),
	})
	if err != nil {
		return err
	}
	*challenge = toModelWeeklyChallenge(row)
	return nil
}
