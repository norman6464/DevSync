package persistence

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/assert"
)

func TestToModelPasswordResetToken(t *testing.T) {
	expiresAt := time.Date(2026, 8, 16, 11, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)

	t.Run("used が true の場合を変換する", func(t *testing.T) {
		used := true
		row := sqlcgen.PasswordResetToken{
			ID:        1,
			UserID:    2,
			Token:     "hashed-token",
			ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
			Used:      &used,
			CreatedAt: pgtype.Timestamptz{Time: createdAt, Valid: true},
		}

		got := toModelPasswordResetToken(row)

		assert.Equal(t, uint(1), got.ID)
		assert.Equal(t, uint(2), got.UserID)
		assert.Equal(t, "hashed-token", got.Token)
		assert.Equal(t, expiresAt, got.ExpiresAt)
		assert.True(t, got.Used)
		assert.Equal(t, createdAt, got.CreatedAt)
	})

	t.Run("used が NULL の場合は false になる", func(t *testing.T) {
		row := sqlcgen.PasswordResetToken{
			Used: nil,
		}

		got := toModelPasswordResetToken(row)

		assert.False(t, got.Used)
	})
}
