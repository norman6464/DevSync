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
