package persistence

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/assert"
)

func TestToModelUserActivity(t *testing.T) {
	createdAt := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)

	t.Run("metadataが設定されている場合はそのまま写す", func(t *testing.T) {
		meta := `{"key":"value"}`
		row := sqlcgen.UserActivity{
			ID:           1,
			UserID:       2,
			ActivityType: "post_created",
			TargetType:   "post",
			TargetID:     3,
			Metadata:     &meta,
			CreatedAt:    pgtype.Timestamptz{Time: createdAt, Valid: true},
		}

		got := toModelUserActivity(row)

		assert.Equal(t, uint(1), got.ID)
		assert.Equal(t, uint(2), got.UserID)
		assert.EqualValues(t, "post_created", got.ActivityType)
		assert.Equal(t, "post", got.TargetType)
		assert.Equal(t, uint(3), got.TargetID)
		assert.Equal(t, meta, got.Metadata)
		assert.Equal(t, createdAt, got.CreatedAt)
	})

	t.Run("metadataがNULLの場合は空文字になる", func(t *testing.T) {
		row := sqlcgen.UserActivity{
			Metadata: nil,
		}

		got := toModelUserActivity(row)

		assert.Equal(t, "", got.Metadata)
	})
}
