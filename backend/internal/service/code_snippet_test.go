package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// newTestCodeSnippetService はテスト用のCodeSnippetServiceインスタンスを生成する。
func newTestCodeSnippetService() (*CodeSnippetService, *MockCodeSnippetRepository, *MockPostRepository) {
	snippetRepo := new(MockCodeSnippetRepository)
	postRepo := new(MockPostRepository)
	svc := NewCodeSnippetService(snippetRepo, postRepo)
	return svc, snippetRepo, postRepo
}

// ---------- スニペットCRUD ----------

func TestSnippetCreate_Success(t *testing.T) {
	svc, snippetRepo, postRepo := newTestCodeSnippetService()

	post := &model.Post{Title: "Test", Content: "Content", UserID: 1}
	post.ID = 5
	postRepo.On("FindByID", uint(5)).Return(post, nil)

	snippet := &model.CodeSnippet{
		PostID:   5,
		UserID:   1,
		Language: "go",
		FileName: "main.go",
		Code:     "package main",
	}
	snippetRepo.On("Create", snippet).Run(func(args mock.Arguments) {
		s := args.Get(0).(*model.CodeSnippet)
		s.ID = 10
	}).Return(nil)
	snippetRepo.On("FindByID", uint(10)).Return(&model.CodeSnippet{
		PostID: 5, UserID: 1, Language: "go", FileName: "main.go", Code: "package main",
	}, nil)

	result, err := svc.Create(snippet)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	snippetRepo.AssertCalled(t, "Create", snippet)
}

func TestSnippetCreate_PostNotFound(t *testing.T) {
	svc, _, postRepo := newTestCodeSnippetService()

	postRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	snippet := &model.CodeSnippet{
		PostID:   999,
		UserID:   1,
		Language: "go",
		Code:     "package main",
	}

	result, err := svc.Create(snippet)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSnippetUpdate_Success(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	existing := &model.CodeSnippet{
		PostID: 5, UserID: 1, Language: "go", FileName: "main.go", Code: "package main",
	}
	existing.ID = 10
	snippetRepo.On("FindByID", uint(10)).Return(existing, nil)
	snippetRepo.On("Update", existing).Return(nil)

	result, err := svc.Update(10, 1, "typescript", "index.ts", "console.log('hello')")
	assert.NoError(t, err)
	assert.Equal(t, "typescript", result.Language)
	assert.Equal(t, "index.ts", result.FileName)
	assert.Equal(t, "console.log('hello')", result.Code)
}

func TestSnippetUpdate_Forbidden(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	existing := &model.CodeSnippet{
		PostID: 5, UserID: 1, Language: "go", Code: "package main",
	}
	existing.ID = 10
	snippetRepo.On("FindByID", uint(10)).Return(existing, nil)

	result, err := svc.Update(10, 999, "typescript", "", "")
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
}

func TestSnippetUpdate_PartialUpdate(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	existing := &model.CodeSnippet{
		PostID: 5, UserID: 1, Language: "go", FileName: "main.go", Code: "package main",
	}
	existing.ID = 10
	snippetRepo.On("FindByID", uint(10)).Return(existing, nil)
	snippetRepo.On("Update", existing).Return(nil)

	// languageのみ更新（fileName, codeは空文字列→変更なし）
	result, err := svc.Update(10, 1, "python", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "python", result.Language)
	assert.Equal(t, "main.go", result.FileName)  // 変更されない
	assert.Equal(t, "package main", result.Code) // 変更されない
}

func TestSnippetDelete_Success(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	existing := &model.CodeSnippet{PostID: 5, UserID: 1, Language: "go", Code: "package main"}
	existing.ID = 10
	snippetRepo.On("FindByID", uint(10)).Return(existing, nil)
	snippetRepo.On("Delete", uint(10)).Return(nil)

	err := svc.Delete(10, 1)
	assert.NoError(t, err)
	snippetRepo.AssertCalled(t, "Delete", uint(10))
}

func TestSnippetDelete_Forbidden(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	existing := &model.CodeSnippet{PostID: 5, UserID: 1, Language: "go", Code: "package main"}
	existing.ID = 10
	snippetRepo.On("FindByID", uint(10)).Return(existing, nil)

	err := svc.Delete(10, 999)
	assert.ErrorIs(t, err, ErrForbidden)
}

// ---------- インラインコメント ----------

func TestSnippetCommentCreate_Success(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	existing := &model.CodeSnippet{PostID: 5, UserID: 1, Language: "go", Code: "package main\nfunc main() {}"}
	existing.ID = 10
	snippetRepo.On("FindByID", uint(10)).Return(existing, nil)

	comment := &model.SnippetComment{
		SnippetID:  10,
		UserID:     2,
		LineNumber: 2,
		Content:    "この関数名をもっと具体的にしてください",
	}
	snippetRepo.On("CreateComment", comment).Return(nil)

	err := svc.CreateComment(comment)
	assert.NoError(t, err)
	snippetRepo.AssertCalled(t, "CreateComment", comment)
}

func TestSnippetCommentCreate_SnippetNotFound(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	snippetRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	comment := &model.SnippetComment{
		SnippetID:  999,
		UserID:     2,
		LineNumber: 1,
		Content:    "コメント",
	}

	err := svc.CreateComment(comment)
	assert.Error(t, err)
	snippetRepo.AssertExpectations(t)
}

func TestSnippetCommentDelete_Success(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	snippetRepo.On("DeleteComment", uint(5), uint(2)).Return(nil)

	err := svc.DeleteComment(5, 2)
	assert.NoError(t, err)
	snippetRepo.AssertCalled(t, "DeleteComment", uint(5), uint(2))
}

func TestSnippetCommentDelete_Forbidden(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	snippetRepo.On("DeleteComment", uint(5), uint(999)).Return(gorm.ErrRecordNotFound)

	err := svc.DeleteComment(5, 999)
	assert.Error(t, err)
}

// ---------- 投稿IDによるスニペット取得 ----------

func TestSnippetGetByPostID_Success(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	snippets := []model.CodeSnippet{
		{PostID: 5, UserID: 1, Language: "go", Code: "package main"},
		{PostID: 5, UserID: 1, Language: "typescript", Code: "console.log()"},
	}
	snippetRepo.On("FindByPostID", uint(5)).Return(snippets, nil)

	result, err := svc.GetByPostID(5)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "go", result[0].Language)
	snippetRepo.AssertExpectations(t)
}

func TestSnippetGetByPostID_Empty(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	snippetRepo.On("FindByPostID", uint(99)).Return([]model.CodeSnippet{}, nil)

	result, err := svc.GetByPostID(99)
	assert.NoError(t, err)
	assert.Empty(t, result)
	snippetRepo.AssertExpectations(t)
}

// ---------- スニペットコメント取得 ----------

func TestSnippetGetComments_Success(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	comments := []model.SnippetComment{
		{SnippetID: 10, UserID: 1, LineNumber: 1, Content: "コメント1"},
		{SnippetID: 10, UserID: 2, LineNumber: 5, Content: "コメント2"},
	}
	snippetRepo.On("GetComments", uint(10)).Return(comments, nil)

	result, err := svc.GetComments(10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "コメント1", result[0].Content)
	snippetRepo.AssertExpectations(t)
}

func TestSnippetGetComments_Empty(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()

	snippetRepo.On("GetComments", uint(99)).Return([]model.SnippetComment{}, nil)

	result, err := svc.GetComments(99)
	assert.NoError(t, err)
	assert.Empty(t, result)
	snippetRepo.AssertExpectations(t)
}

func TestSnippetCreate_RepoError(t *testing.T) {
	svc, snippetRepo, postRepo := newTestCodeSnippetService()
	post := &model.Post{Title: "Test", Content: "Content", UserID: 1}
	post.ID = 5
	postRepo.On("FindByID", uint(5)).Return(post, nil)
	snippet := &model.CodeSnippet{PostID: 5, UserID: 1, Language: "go", Code: "package main"}
	snippetRepo.On("Create", snippet).Return(gorm.ErrInvalidDB)
	result, err := svc.Create(snippet)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSnippetCreate_FindByIDAfterCreateFails(t *testing.T) {
	svc, snippetRepo, postRepo := newTestCodeSnippetService()
	post := &model.Post{Title: "Test", Content: "Content", UserID: 1}
	post.ID = 5
	postRepo.On("FindByID", uint(5)).Return(post, nil)
	snippet := &model.CodeSnippet{PostID: 5, UserID: 1, Language: "go", Code: "package main"}
	snippetRepo.On("Create", snippet).Run(func(args mock.Arguments) {
		s := args.Get(0).(*model.CodeSnippet)
		s.ID = 10
	}).Return(nil)
	snippetRepo.On("FindByID", uint(10)).Return(nil, gorm.ErrRecordNotFound)
	result, err := svc.Create(snippet)
	assert.NoError(t, err)
	assert.NotNil(t, result) // フォールバックでsnippetを返す
}

func TestSnippetDelete_NotFound(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()
	snippetRepo.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)
	err := svc.Delete(99, 1)
	assert.Error(t, err)
}

func TestSnippetUpdate_NotFound(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()
	snippetRepo.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)
	result, err := svc.Update(99, 1, "go", "", "")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSnippetUpdate_RepoError(t *testing.T) {
	svc, snippetRepo, _ := newTestCodeSnippetService()
	existing := &model.CodeSnippet{PostID: 5, UserID: 1, Language: "go", Code: "package main"}
	existing.ID = 10
	snippetRepo.On("FindByID", uint(10)).Return(existing, nil)
	snippetRepo.On("Update", existing).Return(gorm.ErrInvalidDB)
	result, err := svc.Update(10, 1, "python", "", "")
	assert.Error(t, err)
	assert.Nil(t, result)
}
