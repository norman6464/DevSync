package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testModule = "github.com/example/app"

func testPrefix() string { return testModule + "/internal/" }

// --- classifyRel ---

func TestClassifyRel(t *testing.T) {
	tests := []struct {
		rel  string
		want string
	}{
		{"domain", layerDomain},
		{"domain/validator", layerDomain},
		{"model", layerModel},
		{"dto", layerDTO},
		{"usecase", layerUsecase},
		{"usecase/repository", layerUsecaseRepo},
		{"usecase/repository/sub", layerUsecaseRepo},
		{"handler", layerHandler},
		{"handler/middleware", layerHandler},
		{"adapter/persistence", layerAdapter},
		{"adapter/external", layerAdapter},
		{"infra", layerInfra},
		{"infra/ws", layerInfra},
		// 配線と対象外
		{"di", ""},
		{"router", ""},
		{"config", ""},
		{"middleware", ""},
		// 前方一致で誤分類しないこと
		{"models", ""},
		{"usecases", ""},
		{"domainx", ""},
	}
	for _, tt := range tests {
		t.Run(tt.rel, func(t *testing.T) {
			assert.Equal(t, tt.want, classifyRel(tt.rel))
		})
	}
}

// TestClassifyRel_UsecaseRepoBeforeUsecase は port が usecase 本体に吸収されないことを確認する。
// 前方一致の判定順を間違えると usecase/repository が usecase 扱いになり、DIP 違反を見逃す。
func TestClassifyRel_UsecaseRepoBeforeUsecase(t *testing.T) {
	assert.Equal(t, layerUsecaseRepo, classifyRel("usecase/repository"))
	assert.NotEqual(t, layerUsecase, classifyRel("usecase/repository"))
}

// --- classifyImport ---

func TestClassifyImport(t *testing.T) {
	tests := []struct {
		imp  string
		want string
	}{
		{"net/http", targetNetHTTP},
		{"github.com/gin-gonic/gin", targetGin},
		{"github.com/gin-contrib/cors", ""},
		{"gorm.io/gorm", targetGORM},
		{"gorm.io/driver/postgres", targetGORM},
		{testModule + "/internal/usecase", layerUsecase},
		{testModule + "/internal/usecase/repository", layerUsecaseRepo},
		{testModule + "/internal/adapter/persistence", layerAdapter},
		{testModule + "/internal/infra/ws", layerInfra},
		{testModule + "/internal/di", ""},
		// 別モジュールの同名パッケージを誤検出しない
		{"github.com/other/app/internal/usecase", ""},
		{"strings", ""},
	}
	for _, tt := range tests {
		t.Run(tt.imp, func(t *testing.T) {
			assert.Equal(t, tt.want, classifyImport(tt.imp, testPrefix()))
		})
	}
}

// --- violationsFor ---

func TestViolationsFor_DetectsForbiddenImport(t *testing.T) {
	vs := violationsFor(layerUsecase, "internal/usecase/x.go", []importRef{
		{path: testModule + "/internal/adapter/persistence", line: 5, target: layerAdapter},
	})
	require.Len(t, vs, 1)
	assert.Equal(t, 5, vs[0].line)
	assert.Contains(t, vs[0].msg, "usecase は adapter を import できません")
}

func TestViolationsFor_AllowsPermittedImport(t *testing.T) {
	vs := violationsFor(layerUsecase, "internal/usecase/x.go", []importRef{
		{path: testModule + "/internal/usecase/repository", line: 4, target: layerUsecaseRepo},
		{path: testModule + "/internal/model", line: 5, target: layerModel},
		{path: testModule + "/internal/domain", line: 6, target: layerDomain},
	})
	assert.Empty(t, vs)
}

func TestViolationsFor_SkipsSuppressed(t *testing.T) {
	vs := violationsFor(layerUsecase, "internal/usecase/x.go", []importRef{
		{path: "net/http", line: 3, target: targetNetHTTP, suppressed: true},
	})
	assert.Empty(t, vs)
}

// TestViolationsFor_PortMustNotDependOnImplementation は DIP の要である
// 「port が実装を知らない」ことを検出できるか確認する。
func TestViolationsFor_PortMustNotDependOnImplementation(t *testing.T) {
	vs := violationsFor(layerUsecaseRepo, "internal/usecase/repository/x.go", []importRef{
		{path: testModule + "/internal/adapter/persistence", line: 4, target: layerAdapter},
		{path: "gorm.io/gorm", line: 5, target: targetGORM},
	})
	assert.Len(t, vs, 2)
}

// TestViolationsFor_HandlerMustNotTouchAdapter は handler が永続化実装を直接触らないことを確認する。
func TestViolationsFor_HandlerMustNotTouchAdapter(t *testing.T) {
	vs := violationsFor(layerHandler, "internal/handler/x.go", []importRef{
		{path: testModule + "/internal/adapter/persistence", line: 4, target: layerAdapter},
	})
	require.Len(t, vs, 1)
	assert.Contains(t, vs[0].msg, "handler は adapter を直接 import できません")
}

// --- parseImports ---

func TestParseImports(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.go")
	src := `package usecase

import (
	"net/http" //archlint:allow
	"strings"

	"` + testModule + `/internal/adapter/persistence"
)

var _ = strings.TrimSpace
var _ = http.StatusOK
`
	require.NoError(t, os.WriteFile(path, []byte(src), 0o600))

	imports, ignore, err := parseImports(path, testPrefix())
	require.NoError(t, err)
	assert.False(t, ignore)
	require.Len(t, imports, 3)

	assert.True(t, imports[0].suppressed, "//archlint:allow が付いた import は抑制される")
	assert.Equal(t, targetNetHTTP, imports[0].target)
	assert.False(t, imports[2].suppressed)
	assert.Equal(t, layerAdapter, imports[2].target)
}

func TestParseImports_IgnoreFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.go")
	src := `//archlint:ignore-file
package usecase

import "` + testModule + `/internal/adapter/persistence"
`
	require.NoError(t, os.WriteFile(path, []byte(src), 0o600))

	imports, ignore, err := parseImports(path, testPrefix())
	require.NoError(t, err)
	assert.True(t, ignore)
	assert.Empty(t, imports)
}

// TestParseImports_IgnoreFileOnlyBeforePackage は package 宣言より後ろのコメントでは
// ファイル全体の抑制ができないことを確認する（ルールの抜け道を塞ぐ）。
func TestParseImports_IgnoreFileOnlyBeforePackage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.go")
	src := `package usecase

//archlint:ignore-file

import "` + testModule + `/internal/adapter/persistence"
`
	require.NoError(t, os.WriteFile(path, []byte(src), 0o600))

	_, ignore, err := parseImports(path, testPrefix())
	require.NoError(t, err)
	assert.False(t, ignore, "package 宣言より後ろの ignore-file は効かない")
}

// --- readModulePath ---

func TestReadModulePath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(path, []byte("module "+testModule+"\n\ngo 1.24\n"), 0o600))

	got, err := readModulePath(path)
	require.NoError(t, err)
	assert.Equal(t, testModule, got)
}

func TestReadModulePath_Missing(t *testing.T) {
	_, err := readModulePath(filepath.Join(t.TempDir(), "go.mod"))
	assert.Error(t, err)
}

// --- runCLI（end-to-end） ---

// writeModule はテスト用の最小モジュールを作る。
func writeModule(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module "+testModule+"\n"), 0o600))
	for rel, src := range files {
		path := filepath.Join(root, rel)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
		require.NoError(t, os.WriteFile(path, []byte(src), 0o600))
	}
	return root
}

func TestRunCLI_Clean(t *testing.T) {
	root := writeModule(t, map[string]string{
		"internal/usecase/x.go": `package usecase

import "` + testModule + `/internal/usecase/repository"

var _ = repository.Marker
`,
		"internal/usecase/repository/x.go": "package repository\n\nvar Marker = 1\n",
	})

	var stdout, stderr bytes.Buffer
	code := runCLI([]string{root}, &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Contains(t, stdout.String(), "OK")
	assert.Empty(t, stderr.String())
}

func TestRunCLI_ReportsViolations(t *testing.T) {
	root := writeModule(t, map[string]string{
		"internal/handler/x.go": `package handler

import "` + testModule + `/internal/adapter/persistence"

var _ = persistence.Marker
`,
		"internal/adapter/persistence/x.go": "package persistence\n\nvar Marker = 1\n",
	})

	var stdout, stderr bytes.Buffer
	code := runCLI([]string{root}, &stdout, &stderr)

	assert.Equal(t, 1, code)
	assert.Contains(t, stdout.String(), "handler は adapter を直接 import できません")
	assert.Contains(t, stderr.String(), "1 件")
}

// TestRunCLI_SkipsWiringAndTests は配線パッケージとテストファイルを対象外にすることを確認する。
func TestRunCLI_SkipsWiringAndTests(t *testing.T) {
	root := writeModule(t, map[string]string{
		// 配線は全層を組み立てるので対象外
		"internal/di/container.go": `package di

import "` + testModule + `/internal/adapter/persistence"

var _ = persistence.Marker
`,
		// テストは port モックのために port を import してよい
		"internal/handler/x_test.go": `package handler

import "` + testModule + `/internal/adapter/persistence"

var _ = persistence.Marker
`,
		"internal/adapter/persistence/x.go": "package persistence\n\nvar Marker = 1\n",
	})

	var stdout, stderr bytes.Buffer
	code := runCLI([]string{root}, &stdout, &stderr)

	assert.Equal(t, 0, code, stdout.String())
}

func TestRunCLI_MissingGoMod(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runCLI([]string{t.TempDir()}, &stdout, &stderr)

	assert.Equal(t, 2, code)
	assert.Contains(t, stderr.String(), "go.mod")
}

// TestRunCLI_OnThisRepository は実際のこのリポジトリに対して違反が無いことを確認する。
// 依存方向が崩れたらこのテストが落ちる（CI ステップと同じ検証をテストでも行う）。
func TestRunCLI_OnThisRepository(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runCLI([]string{"../.."}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("依存方向違反が見つかりました:\n%s%s", stdout.String(), stderr.String())
	}
	assert.True(t, strings.Contains(stdout.String(), "OK"))
}
