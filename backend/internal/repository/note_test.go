package repository

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupNoteTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.Note{}, &model.NoteFolder{}, &model.User{})
	assert.NoError(t, err)

	return db
}

func TestNoteRepository_Create(t *testing.T) {
	db := setupNoteTestDB(t)
	repo := NewNoteRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	note := &model.Note{
		UserID:  user.ID,
		Title:   "テストノート",
		Content: "これはテストノートです",
		Tags:    "Go,TDD",
	}

	err := repo.Create(note)
	assert.NoError(t, err)
	assert.NotZero(t, note.ID)
}

func TestNoteRepository_FindByID(t *testing.T) {
	db := setupNoteTestDB(t)
	repo := NewNoteRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	note := &model.Note{
		UserID:  user.ID,
		Title:   "テストノート",
		Content: "これはテストノートです",
	}
	db.Create(note)

	// 取得
	retrieved, err := repo.FindByID(note.ID)
	assert.NoError(t, err)
	assert.Equal(t, note.ID, retrieved.ID)
	assert.Equal(t, "テストノート", retrieved.Title)
}

func TestNoteRepository_FindByUserID(t *testing.T) {
	db := setupNoteTestDB(t)
	repo := NewNoteRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	// 複数ノート作成
	db.Create(&model.Note{UserID: user.ID, Title: "ノート1", Content: "内容1"})
	db.Create(&model.Note{UserID: user.ID, Title: "ノート2", Content: "内容2"})

	// 取得
	notes, err := repo.FindByUserID(user.ID, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(notes))
}

func TestNoteRepository_Update(t *testing.T) {
	db := setupNoteTestDB(t)
	repo := NewNoteRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	note := &model.Note{
		UserID:  user.ID,
		Title:   "旧タイトル",
		Content: "旧内容",
	}
	db.Create(note)

	// 更新
	note.Title = "新タイトル"
	note.Content = "新内容"
	err := repo.Update(note)
	assert.NoError(t, err)

	// 確認
	retrieved, _ := repo.FindByID(note.ID)
	assert.Equal(t, "新タイトル", retrieved.Title)
	assert.Equal(t, "新内容", retrieved.Content)
}

func TestNoteRepository_Delete(t *testing.T) {
	db := setupNoteTestDB(t)
	repo := NewNoteRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	note := &model.Note{
		UserID:  user.ID,
		Title:   "削除予定",
		Content: "削除されます",
	}
	db.Create(note)

	// 削除
	err := repo.Delete(note.ID)
	assert.NoError(t, err)

	// 確認（存在しないはず）
	_, err = repo.FindByID(note.ID)
	assert.Error(t, err)
}

func TestNoteRepository_Search(t *testing.T) {
	db := setupNoteTestDB(t)
	repo := NewNoteRepository(db)

	// テスト用ユーザー作成
	user := &model.User{Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	// テストノート作成
	db.Create(&model.Note{UserID: user.ID, Title: "Goの学習", Content: "Goを学んでいます"})
	db.Create(&model.Note{UserID: user.ID, Title: "Reactの学習", Content: "Reactを学んでいます"})

	// 検索
	notes, total, err := repo.Search(user.ID, "Go", 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Goの学習", notes[0].Title)
}
