package persistence

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/assert"
)

func TestToModelWeeklyChallenge(t *testing.T) {
	createdAt := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)

	t.Run("完了済み（completed_atあり）を変換する", func(t *testing.T) {
		completedAt := time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC)
		currentValue := int64(300)
		isCompleted := true
		row := sqlcgen.WeeklyChallenge{
			ID:            1,
			UserID:        2,
			Year:          2026,
			Week:          33,
			ChallengeType: "duration_total",
			Description:   "weekly_duration_300",
			TargetValue:   300,
			CurrentValue:  &currentValue,
			IsCompleted:   &isCompleted,
			CompletedAt:   pgtype.Timestamptz{Time: completedAt, Valid: true},
			CreatedAt:     pgtype.Timestamptz{Time: createdAt, Valid: true},
			UpdatedAt:     pgtype.Timestamptz{Time: updatedAt, Valid: true},
		}

		got := toModelWeeklyChallenge(row)

		assert.Equal(t, uint(1), got.ID)
		assert.Equal(t, uint(2), got.UserID)
		assert.Equal(t, 2026, got.Year)
		assert.Equal(t, 33, got.Week)
		assert.EqualValues(t, "duration_total", got.ChallengeType)
		assert.Equal(t, 300, got.CurrentValue)
		assert.True(t, got.IsCompleted)
		assert.Equal(t, completedAt, *got.CompletedAt)
		assert.Equal(t, createdAt, got.CreatedAt)
		assert.Equal(t, updatedAt, got.UpdatedAt)
	})

	t.Run("未完了（completed_at・current_value・is_completedがNULL）を変換する", func(t *testing.T) {
		row := sqlcgen.WeeklyChallenge{
			ID:            1,
			UserID:        2,
			ChallengeType: "streak_days",
			CurrentValue:  nil,
			IsCompleted:   nil,
			CompletedAt:   pgtype.Timestamptz{},
			CreatedAt:     pgtype.Timestamptz{Time: createdAt, Valid: true},
			UpdatedAt:     pgtype.Timestamptz{Time: updatedAt, Valid: true},
		}

		got := toModelWeeklyChallenge(row)

		assert.Equal(t, 0, got.CurrentValue)
		assert.False(t, got.IsCompleted)
		assert.Nil(t, got.CompletedAt)
	})
}
