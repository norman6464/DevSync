package repository

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupNoteFolderTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.NoteFolder{}, &model.Note{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// テスト用ユーザー作成
	user := &model.User{
		Username: "testuser",
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	db.Create(user)

	return db
}

func TestNoteFolderRepository_Create(t *testing.T) {
	db := setupNoteFolderTestDB(t)
	repo := NewNoteFolderRepository(db)

	folder := &model.NoteFolder{
		UserID: 1,
		Name:   "マイフォルダ",
	}

	err := repo.Create(folder)
	assert.NoError(t, err)
	assert.NotZero(t, folder.ID)
	assert.NotZero(t, folder.CreatedAt)
}

func TestNoteFolderRepository_FindByID(t *testing.T) {
	db := setupNoteFolderTestDB(t)
	repo := NewNoteFolderRepository(db)

	// テストデータ作成
	folder := &model.NoteFolder{
		UserID: 1,
		Name:   "テストフォルダ",
	}
	db.Create(folder)

	// 取得テスト
	result, err := repo.FindByID(folder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "テストフォルダ", result.Name)
	assert.Equal(t, uint(1), result.UserID)
}

func TestNoteFolderRepository_FindByUserID(t *testing.T) {
	db := setupNoteFolderTestDB(t)
	repo := NewNoteFolderRepository(db)

	// テストデータ作成
	folders := []*model.NoteFolder{
		{UserID: 1, Name: "フォルダ1"},
		{UserID: 1, Name: "フォルダ2"},
		{UserID: 2, Name: "他ユーザーのフォルダ"},
	}
	for _, f := range folders {
		db.Create(f)
	}

	// 取得テスト
	results, total, err := repo.FindByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(2), total)
}

func TestNoteFolderRepository_FindByParentID(t *testing.T) {
	db := setupNoteFolderTestDB(t)
	repo := NewNoteFolderRepository(db)

	// 親フォルダ作成
	parent := &model.NoteFolder{
		UserID: 1,
		Name:   "親フォルダ",
	}
	db.Create(parent)

	// 子フォルダ作成
	child1 := &model.NoteFolder{
		UserID:   1,
		ParentID: &parent.ID,
		Name:     "子フォルダ1",
	}
	child2 := &model.NoteFolder{
		UserID:   1,
		ParentID: &parent.ID,
		Name:     "子フォルダ2",
	}
	db.Create(child1)
	db.Create(child2)

	// 取得テスト
	results, err := repo.FindByParentID(parent.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestNoteFolderRepository_Update(t *testing.T) {
	db := setupNoteFolderTestDB(t)
	repo := NewNoteFolderRepository(db)

	// テストデータ作成
	folder := &model.NoteFolder{
		UserID: 1,
		Name:   "元の名前",
	}
	db.Create(folder)

	// 更新
	folder.Name = "新しい名前"
	err := repo.Update(folder)
	assert.NoError(t, err)

	// 確認
	result, _ := repo.FindByID(folder.ID)
	assert.Equal(t, "新しい名前", result.Name)
}

func TestNoteFolderRepository_Delete(t *testing.T) {
	db := setupNoteFolderTestDB(t)
	repo := NewNoteFolderRepository(db)

	// テストデータ作成
	folder := &model.NoteFolder{
		UserID: 1,
		Name:   "削除対象フォルダ",
	}
	db.Create(folder)

	// 削除
	err := repo.Delete(folder.ID)
	assert.NoError(t, err)

	// 確認
	_, err = repo.FindByID(folder.ID)
	assert.Error(t, err)
}

func TestNoteFolderRepository_GetRootFolders(t *testing.T) {
	db := setupNoteFolderTestDB(t)
	repo := NewNoteFolderRepository(db)

	// テストデータ作成
	root1 := &model.NoteFolder{
		UserID:   1,
		ParentID: nil, // ルートフォルダ
		Name:     "ルート1",
	}
	root2 := &model.NoteFolder{
		UserID:   1,
		ParentID: nil, // ルートフォルダ
		Name:     "ルート2",
	}
	db.Create(root1)
	db.Create(root2)

	// 子フォルダ作成
	child := &model.NoteFolder{
		UserID:   1,
		ParentID: &root1.ID,
		Name:     "子フォルダ",
	}
	db.Create(child)

	// ルートフォルダのみ取得
	results, err := repo.GetRootFolders(1)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}
