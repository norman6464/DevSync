package repository

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupNotificationSettingsTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.NotificationSettings{}, &model.User{})
	assert.NoError(t, err)

	return db
}

func TestNotificationSettingsRepository_CreateOrUpdate(t *testing.T) {
	db := setupNotificationSettingsTestDB(t)
	repo := NewNotificationSettingsRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	settings := &model.NotificationSettings{
		UserID:          user.ID,
		EnableLikes:     true,
		EnableComments:  true,
		EnableFollows:   false,
		EnableMessages:  true,
		EnableMentions:  true,
		EnableWebPush:   true,
		EnableEmail:     false,
		EnableSound:     true,
	}

	// Create
	err := repo.CreateOrUpdate(settings)
	assert.NoError(t, err)
	assert.NotZero(t, settings.ID)

	// Update
	settings.EnableFollows = true
	settings.EnableEmail = true
	err = repo.CreateOrUpdate(settings)
	assert.NoError(t, err)

	// Verify
	retrieved, err := repo.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, retrieved.EnableFollows)
	assert.Equal(t, true, retrieved.EnableEmail)
}

func TestNotificationSettingsRepository_GetByUserID(t *testing.T) {
	db := setupNotificationSettingsTestDB(t)
	repo := NewNotificationSettingsRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	// 設定が存在しない場合
	settings, err := repo.GetByUserID(user.ID)
	assert.Error(t, err)
	assert.Nil(t, settings)

	// 設定を作成
	newSettings := &model.NotificationSettings{
		UserID:         user.ID,
		EnableLikes:    true,
		EnableComments: true,
	}
	db.Create(newSettings)

	// 設定が存在する場合
	settings, err = repo.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, user.ID, settings.UserID)
	assert.Equal(t, true, settings.EnableLikes)
}

func TestNotificationSettingsRepository_GetOrCreateDefault(t *testing.T) {
	db := setupNotificationSettingsTestDB(t)
	repo := NewNotificationSettingsRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	// 設定が存在しない場合 → デフォルト設定を作成
	settings, err := repo.GetOrCreateDefault(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, user.ID, settings.UserID)
	assert.Equal(t, true, settings.EnableLikes)     // デフォルトtrue
	assert.Equal(t, true, settings.EnableComments)  // デフォルトtrue
	assert.Equal(t, true, settings.EnableFollows)   // デフォルトtrue
	assert.Equal(t, true, settings.EnableMessages)  // デフォルトtrue
	assert.Equal(t, true, settings.EnableMentions)  // デフォルトtrue
	assert.Equal(t, true, settings.EnableWebPush)   // デフォルトtrue
	assert.Equal(t, true, settings.EnableEmail)     // デフォルトtrue
	assert.Equal(t, true, settings.EnableSound)     // デフォルトtrue

	// 既に設定が存在する場合 → 既存設定を返す
	settings.EnableSound = false
	db.Save(settings)

	settings2, err := repo.GetOrCreateDefault(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, settings.ID, settings2.ID)
	assert.Equal(t, false, settings2.EnableSound)
}
