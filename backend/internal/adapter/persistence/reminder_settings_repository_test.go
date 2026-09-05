package persistence

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestToModelReminderSettings(t *testing.T) {
	createdAt := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 8, 17, 10, 0, 0, 0, time.UTC)

	t.Run("last_reminded_atありの設定を変換する", func(t *testing.T) {
		enabled := true
		frequency := "weekly"
		notificationTime := "21:30"
		inactiveDays := int64(7)
		enableWeb := false
		enableEmail := true
		lastRemindedAt := time.Date(2026, 8, 15, 9, 0, 0, 0, time.UTC)

		row := sqlcgen.ReminderSetting{
			ID:               10,
			UserID:           1,
			Enabled:          &enabled,
			Frequency:        &frequency,
			NotificationTime: &notificationTime,
			InactiveDays:     &inactiveDays,
			EnableWeb:        &enableWeb,
			EnableEmail:      &enableEmail,
			LastRemindedAt:   pgtype.Timestamptz{Time: lastRemindedAt, Valid: true},
			CreatedAt:        pgtype.Timestamptz{Time: createdAt, Valid: true},
			UpdatedAt:        pgtype.Timestamptz{Time: updatedAt, Valid: true},
		}

		got := toModelReminderSettings(row)

		assert.Equal(t, uint(10), got.ID)
		assert.Equal(t, uint(1), got.UserID)
		assert.True(t, got.Enabled)
		assert.Equal(t, model.ReminderFrequencyWeekly, got.Frequency)
		assert.Equal(t, "21:30", got.NotificationTime)
		assert.Equal(t, 7, got.InactiveDays)
		assert.False(t, got.EnableWeb)
		assert.True(t, got.EnableEmail)
		assert.Equal(t, lastRemindedAt, *got.LastRemindedAt)
		assert.Equal(t, createdAt, got.CreatedAt)
		assert.Equal(t, updatedAt, got.UpdatedAt)
	})

	t.Run("last_reminded_atがNULLの場合はnilになる", func(t *testing.T) {
		row := sqlcgen.ReminderSetting{
			LastRemindedAt: pgtype.Timestamptz{},
		}

		got := toModelReminderSettings(row)

		assert.Nil(t, got.LastRemindedAt)
	})
}
