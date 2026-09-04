package persistence

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/assert"
)

func TestToModelPostTemplate(t *testing.T) {
	createdAt := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 8, 17, 12, 30, 0, 0, time.UTC)

	t.Run("title_template が設定されている場合はそのまま写す", func(t *testing.T) {
		title := "週報 {date}"
		row := sqlcgen.PostTemplate{
			ID:              1,
			UserID:          2,
			Name:            "週報テンプレ",
			TitleTemplate:   &title,
			ContentTemplate: "content",
			CreatedAt:       pgtype.Timestamptz{Time: createdAt, Valid: true},
			UpdatedAt:       pgtype.Timestamptz{Time: updatedAt, Valid: true},
		}

		got := toModelPostTemplate(row)

		assert.Equal(t, uint(1), got.ID)
		assert.Equal(t, uint(2), got.UserID)
		assert.Equal(t, "週報テンプレ", got.Name)
		assert.Equal(t, "週報 {date}", got.TitleTemplate)
		assert.Equal(t, "content", got.ContentTemplate)
		assert.Equal(t, createdAt, got.CreatedAt)
		assert.Equal(t, updatedAt, got.UpdatedAt)
	})

	t.Run("title_template が NULL の場合は空文字になる（GORM の string フィールドと同じ既定値）", func(t *testing.T) {
		row := sqlcgen.PostTemplate{
			ID:              1,
			UserID:          2,
			Name:            "名前のみ",
			TitleTemplate:   nil,
			ContentTemplate: "content",
			CreatedAt:       pgtype.Timestamptz{Time: createdAt, Valid: true},
			UpdatedAt:       pgtype.Timestamptz{Time: updatedAt, Valid: true},
		}

		got := toModelPostTemplate(row)

		assert.Equal(t, "", got.TitleTemplate)
	})
}
