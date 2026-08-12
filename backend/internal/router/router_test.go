package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetupWithContainer_RegistersAllRoutes は全ルートの登録が通ることを確認する。
// gin は同じ位置のパスパラメータに別名を使うと登録時に panic するため、
// このテストが無いと「起動しないサーバー」をマージできてしまう。
// ハンドラーは呼ばないので、コンテナは空で良い。
func TestSetupWithContainer_RegistersAllRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var r *gin.Engine
	require.NotPanics(t, func() {
		r = SetupWithContainer(&di.Container{}, &config.Config{CORSOrigins: "http://localhost:5173"})
	}, "ルート登録で panic してはいけない（サーバーが起動できなくなる）")

	require.NotNil(t, r)
	routes := r.Routes()
	assert.NotEmpty(t, routes)

	// 代表的なエンドポイントが登録されていること
	registered := make(map[string]bool, len(routes))
	for _, route := range routes {
		registered[route.Method+" "+route.Path] = true
	}
	for _, want := range []string{
		"GET /health",
		"GET /ws",
		"POST /api/v1/auth/register",
		"POST /api/v1/auth/login",
		"GET /api/v1/auth/me",
		"GET /api/v1/users/:id",
		"GET /api/v1/users/:id/activity",
		"GET /api/v1/users/by-username/:username",
	} {
		assert.True(t, registered[want], "%s が登録されていない", want)
	}
}

// TestSetupWithContainer_NoConflictingParamNames は同じ階層のパスパラメータ名が
// 揃っていることを確認する。gin の panic は最初の 1 件で止まるため、
// 登録済みのルート全体を突き合わせて残りの衝突も検出する。
func TestSetupWithContainer_NoConflictingParamNames(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := SetupWithContainer(&di.Container{}, &config.Config{CORSOrigins: "http://localhost:5173"})

	// 「その位置までのパス」→ 使われているパラメータ名
	paramAt := make(map[string]string)
	for _, route := range r.Routes() {
		var prefix string
		for _, segment := range splitPath(route.Path) {
			if len(segment) > 0 && segment[0] == ':' {
				if existing, ok := paramAt[prefix]; ok {
					assert.Equal(t, existing, segment,
						"%s の同じ位置で別のパラメータ名が使われている（path=%s）", prefix, route.Path)
				} else {
					paramAt[prefix] = segment
				}
			}
			prefix += "/" + segment
		}
	}
}

// splitPath は "/a/b/:id" を ["a","b",":id"] に分解する。
func splitPath(path string) []string {
	var segments []string
	current := ""
	for _, ch := range path {
		if ch == '/' {
			if current != "" {
				segments = append(segments, current)
				current = ""
			}
			continue
		}
		current += string(ch)
	}
	if current != "" {
		segments = append(segments, current)
	}
	return segments
}
