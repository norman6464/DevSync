package repository

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupReminderSettingsTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.ReminderSettings{}, &model.User{})
	assert.NoError(t, err)

	return db
}

func TestReminderSettingsRepository_CreateOrUpdate(t *testing.T) {
	db := setupReminderSettingsTestDB(t)
	repo := NewReminderSettingsRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	settings := &model.ReminderSettings{
		UserID:           user.ID,
		Enabled:          true,
		Frequency:        model.ReminderFrequencyDaily,
		NotificationTime: "09:00",
		InactiveDays:     3,
		EnableWeb:        true,
		EnableEmail:      false,
	}

	// Create
	err := repo.CreateOrUpdate(settings)
	assert.NoError(t, err)
	assert.NotZero(t, settings.ID)

	// Update
	settings.Frequency = model.ReminderFrequencyWeekly
	settings.EnableEmail = true
	err = repo.CreateOrUpdate(settings)
	assert.NoError(t, err)

	// Verify
	retrieved, err := repo.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.ReminderFrequencyWeekly, retrieved.Frequency)
	assert.Equal(t, true, retrieved.EnableEmail)
}

func TestReminderSettingsRepository_GetByUserID(t *testing.T) {
	db := setupReminderSettingsTestDB(t)
	repo := NewReminderSettingsRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	// 設定が存在しない場合
	settings, err := repo.GetByUserID(user.ID)
	assert.Error(t, err)
	assert.Nil(t, settings)

	// 設定を作成
	newSettings := &model.ReminderSettings{
		UserID:   user.ID,
		Enabled:  true,
		Frequency: model.ReminderFrequencyDaily,
	}
	db.Create(newSettings)

	// 設定が存在する場合
	settings, err = repo.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, user.ID, settings.UserID)
	assert.Equal(t, true, settings.Enabled)
}

func TestReminderSettingsRepository_GetOrCreateDefault(t *testing.T) {
	db := setupReminderSettingsTestDB(t)
	repo := NewReminderSettingsRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	// 設定が存在しない場合 → デフォルト設定を作成
	settings, err := repo.GetOrCreateDefault(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, user.ID, settings.UserID)
	assert.Equal(t, true, settings.Enabled)                                   // デフォルトtrue
	assert.Equal(t, model.ReminderFrequencyDaily, settings.Frequency)         // デフォルトdaily
	assert.Equal(t, "09:00", settings.NotificationTime)                       // デフォルト09:00
	assert.Equal(t, 3, settings.InactiveDays)                                 // デフォルト3日
	assert.Equal(t, true, settings.EnableWeb)                                 // デフォルトtrue
	assert.Equal(t, false, settings.EnableEmail)                              // デフォルトfalse

	// 既に設定が存在する場合 → 既存設定を返す
	settings.EnableEmail = true
	db.Save(settings)

	settings2, err := repo.GetOrCreateDefault(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, settings.ID, settings2.ID)
	assert.Equal(t, true, settings2.EnableEmail)
}
