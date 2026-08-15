package persistence

import (
	"sync"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/schema"
)

// relationsOf はモデルが持つ関連名を返す。
// Preload は関連が無いと実行時に "unsupported relations" で失敗し、
// クエリ自体が成功していても結果が捨てられるため、宣言の有無をここで固定する。
func relationsOf(t *testing.T, model any) map[string]bool {
	t.Helper()
	parsed, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	require.NoError(t, err)

	names := make(map[string]bool, len(parsed.Relationships.Relations))
	for name := range parsed.Relationships.Relations {
		names[name] = true
	}
	return names
}

// TestCodeSnippetHasNoUserRelation はコードスニペットが投稿者の関連を持たないことを確認する。
// 検索クエリはこの前提で Preload を外している。関連を足すなら Preload を戻してよい。
func TestCodeSnippetHasNoUserRelation(t *testing.T) {
	relations := relationsOf(t, &model.CodeSnippet{})

	assert.False(t, relations["User"],
		"User 関連が無いモデルを Preload すると検索が 500 になる。関連を追加したなら検索の Preload も見直すこと")
}

// TestSnippetCommentHasUserRelation はインラインコメントが投稿者の関連を持つことを確認する。
// こちらは Preload("User") を使っているため、関連が消えると同じ壊れ方をする。
func TestSnippetCommentHasUserRelation(t *testing.T) {
	relations := relationsOf(t, &model.SnippetComment{})

	assert.True(t, relations["User"], "Preload(\"User\") を使っているため関連が必要")
}
