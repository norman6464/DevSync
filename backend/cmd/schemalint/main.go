// Command schemalint は、宣言的スキーマ（internal/infra/database/schema/schema.hcl）が
// 実際に適用された PostgreSQL に対して、DB設計規約からの逸脱を検出する。
//
// archlint が「Go の import 方向」を静的に検証するのに対し、schemalint は「DB の構造」を
// 実際に適用済みの PostgreSQL へ接続し pg_catalog / information_schema から検証する。
// Atlas OSS 自体はスキーマの適用はできても「FKにon_deleteが宣言されているか」
// 「子側に索引があるか」といった設計規約は検証しないため、適用後の実体を見る必要がある。
//
// 使い方:
//
//	ATLAS_DB_URL=postgres://devsync:devsync@localhost:5432/devsync?sslmode=disable \
//	  go run ./cmd/schemalint
//
// 検証対象のDBには internal/infra/database/schema/schema.hcl が適用済みであること
// （`make db-schema-apply` または CI の apply ステップの後に実行する）。
//
// 現時点（Phase 0, DEVSYNC-155）ではルールをすべて warn として出力し、
// exit code は常に 0 を返す。DEVSYNC-156以降でFK・索引・命名規約を実際に整えたあと、
// ルールを1本ずつ error（違反があれば exit code 1）へ昇格させる
// （このコメントと下の ruleSeverity を合わせて更新すること）。
package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/jackc/pgx/v5"
)

// severity はルールが違反を検出したときの扱い。
type severity int

const (
	severityWarn severity = iota
	severityError
)

func (s severity) String() string {
	if s == severityError {
		return "error"
	}
	return "warn"
}

// ruleSeverity は各ルールの現在の重大度。Phase 0 では全ルールが warn。
// ルールをerrorへ昇格させる際はここだけを変更する。
var ruleSeverity = map[string]severity{
	"fk-missing":       severityWarn,
	"fk-index-missing": severityWarn,
	"naming":           severityWarn,
}

type finding struct {
	rule    string
	table   string
	column  string
	message string
}

func main() {
	os.Exit(runCLI(os.Stdout, os.Stderr))
}

// runCLI は os.Exit を呼ばないことで end-to-end テストを可能にする。
func runCLI(stdout, stderr *os.File) int {
	dbURL := os.Getenv("ATLAS_DB_URL")
	if dbURL == "" {
		fmt.Fprintln(stderr, "schemalint: 環境変数 ATLAS_DB_URL が設定されていません")
		return 1
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(stderr, "schemalint: DB接続に失敗しました: %v\n", err)
		return 1
	}
	defer func() { _ = conn.Close(ctx) }()

	var findings []finding
	for _, rule := range rules {
		fs, err := rule.check(ctx, conn)
		if err != nil {
			fmt.Fprintf(stderr, "schemalint: ルール %q の実行に失敗しました: %v\n", rule.name, err)
			return 1
		}
		findings = append(findings, fs...)
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].table != findings[j].table {
			return findings[i].table < findings[j].table
		}
		if findings[i].column != findings[j].column {
			return findings[i].column < findings[j].column
		}
		return findings[i].rule < findings[j].rule
	})

	errorCount := 0
	for _, f := range findings {
		sev := ruleSeverity[f.rule]
		if sev == severityError {
			errorCount++
		}
		fmt.Fprintf(stdout, "[%s/%s] %s.%s: %s\n", sev, f.rule, f.table, f.column, f.message)
	}

	if len(findings) == 0 {
		fmt.Fprintln(stdout, "schemalint: OK — 検出されたルール逸脱なし")
		return 0
	}

	fmt.Fprintf(stdout, "\nschemalint: %d 件検出（うちerror %d件・warn %d件）\n",
		len(findings), errorCount, len(findings)-errorCount)

	if errorCount > 0 {
		return 1
	}
	return 0
}
