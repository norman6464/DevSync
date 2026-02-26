package service

import (
	"encoding/json"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// defaultWidgetSettings はダッシュボードウィジェットのデフォルト設定。
const defaultWidgetSettings = `[{"key":"userProfile","visible":true,"order":0},{"key":"level","visible":true,"order":1},{"key":"streak","visible":true,"order":2},{"key":"dailyChallenge","visible":true,"order":3},{"key":"weeklyChallenge","visible":true,"order":4},{"key":"studyCircle","visible":true,"order":5},{"key":"quickEntry","visible":true,"order":6},{"key":"quickActions","visible":true,"order":7},{"key":"recommendedUsers","visible":true,"order":8},{"key":"trending","visible":true,"order":9},{"key":"aiAdvice","visible":true,"order":10},{"key":"goalsProgress","visible":true,"order":11},{"key":"recentNotifications","visible":true,"order":12},{"key":"quickStats","visible":true,"order":13}]`

// WidgetSettingsService はダッシュボードウィジェット設定のビジネスロジックを提供する。
type WidgetSettingsService struct {
	repo repository.WidgetSettingsRepositoryInterface
}

// NewWidgetSettingsService は新しいWidgetSettingsServiceインスタンスを生成する。
func NewWidgetSettingsService(repo repository.WidgetSettingsRepositoryInterface) *WidgetSettingsService {
	return &WidgetSettingsService{repo: repo}
}

// GetSettings は指定ユーザーのウィジェット設定を取得する。
// 設定が未登録の場合はデフォルト設定を返す。
func (s *WidgetSettingsService) GetSettings(userID uint) (*model.WidgetSettings, error) {
	settings, err := s.repo.FindByUserID(userID)
	if err != nil {
		// レコード未登録の場合はデフォルト設定を返す
		return &model.WidgetSettings{
			UserID:   userID,
			Settings: defaultWidgetSettings,
		}, nil
	}
	return settings, nil
}

// UpdateSettings はウィジェット設定を更新する。
// 設定はJSON配列形式である必要がある。
func (s *WidgetSettingsService) UpdateSettings(userID uint, settings string) error {
	if strings.TrimSpace(settings) == "" {
		return domain.NewError(domain.ErrCodeBadRequest, "設定は必須です", nil)
	}

	if !json.Valid([]byte(settings)) {
		return domain.NewError(domain.ErrCodeBadRequest, "設定は有効なJSONである必要があります", nil)
	}

	// JSON配列であることを検証
	var arr []json.RawMessage
	if err := json.Unmarshal([]byte(settings), &arr); err != nil {
		return domain.NewError(domain.ErrCodeBadRequest, "設定はJSON配列である必要があります", nil)
	}

	ws := &model.WidgetSettings{
		UserID:   userID,
		Settings: settings,
	}
	return s.repo.Upsert(ws)
}
