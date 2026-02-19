package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// 学習ログCSVエクスポートテスト
// ============================================================

func TestExportCSV_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	now := time.Now()
	logs := []model.LearningLog{
		{
			ID:        1,
			UserID:    1,
			Title:     "Goの基礎",
			Content:   "変数と型を学んだ",
			Category:  model.LogCategoryCoding,
			Duration:  60,
			CreatedAt: now,
		},
		{
			ID:        2,
			UserID:    1,
			Title:     "クリーンアーキテクチャ",
			Content:   "依存性逆転の原則",
			Category:  model.LogCategoryReading,
			Duration:  45,
			CreatedAt: now.Add(-24 * time.Hour),
		},
	}

	repo.On("GetByPeriod", uint(1), 30).Return(logs, nil)

	csvBytes, err := svc.ExportCSV(1, 30)
	assert.NoError(t, err)
	assert.NotNil(t, csvBytes)

	csvContent := string(csvBytes)
	// ヘッダー行の確認
	assert.Contains(t, csvContent, "日付")
	assert.Contains(t, csvContent, "カテゴリ")
	assert.Contains(t, csvContent, "タイトル")
	assert.Contains(t, csvContent, "学習時間(分)")
	// データ行の確認
	assert.Contains(t, csvContent, "Goの基礎")
	assert.Contains(t, csvContent, "クリーンアーキテクチャ")
	assert.Contains(t, csvContent, "60")
	assert.Contains(t, csvContent, "45")
}

func TestExportCSV_EmptyLogs(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetByPeriod", uint(1), 7).Return([]model.LearningLog{}, nil)

	csvBytes, err := svc.ExportCSV(1, 7)
	assert.NoError(t, err)
	// ヘッダー行のみが含まれる
	lines := strings.Split(strings.TrimSpace(string(csvBytes)), "\n")
	assert.Equal(t, 1, len(lines)) // ヘッダー行のみ
}

func TestExportCSV_InvalidPeriod(t *testing.T) {
	svc, _ := newTestLearningLogService()

	_, err := svc.ExportCSV(1, -1)
	assert.Error(t, err)
}

func TestExportCSV_AllPeriod(t *testing.T) {
	svc, repo := newTestLearningLogService()

	logs := []model.LearningLog{
		{Title: "過去ログ", Category: model.LogCategoryCourse, Duration: 30},
	}
	// period=0 は全期間
	repo.On("GetByPeriod", uint(1), 0).Return(logs, nil)

	csvBytes, err := svc.ExportCSV(1, 0)
	assert.NoError(t, err)
	assert.Contains(t, string(csvBytes), "過去ログ")
}

func TestExportCSV_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetByPeriod", uint(1), 30).Return([]model.LearningLog{}, errors.New("db error"))

	_, err := svc.ExportCSV(1, 30)
	assert.Error(t, err)
}
