package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// rule は1つの検証観点。check は違反(finding)の一覧を返す（エラーが無ければ len=0）。
type rule struct {
	name  string
	check func(ctx context.Context, conn *pgx.Conn) ([]finding, error)
}

var rules = []rule{
	{name: "fk-missing", check: checkFKMissing},
	{name: "fk-index-missing", check: checkFKIndexMissing},
	{name: "naming", check: checkNaming},
}

// externalIDColumns は内部テーブルではなく外部サービス（GitHub/Qiita/Zenn等）のIDを保持する
// 列で、FKを張れないことが設計上正しいため fk-missing の対象から除外する。
// ここに載っていない `%_id` bigint 列は「内部テーブルを指すはずなのにFKが無い」候補として
// 検出される（target_id のようなポリモーフィック参照も意図的に除外しない — それ自体が
// 設計上の問題として検出されるべきもの）。
var externalIDColumns = map[string]bool{
	"git_hub_id":      true, // users.git_hub_id — GitHub側のユーザー数値ID
	"git_hub_repo_id": true,
	"github_repo_id":  true,
	"zenn_id":         true,
	"qiita_id":        true,
}

// checkFKMissing は、bigint型で `%_id` という名前を持つが FK 制約が張られていない列を検出する。
func checkFKMissing(ctx context.Context, conn *pgx.Conn) ([]finding, error) {
	rows, err := conn.Query(ctx, `
		SELECT c.table_name, c.column_name
		FROM information_schema.columns c
		WHERE c.table_schema = 'public'
		  AND c.data_type = 'bigint'
		  AND c.column_name LIKE '%\_id' ESCAPE '\'
		  AND NOT EXISTS (
		    SELECT 1
		    FROM information_schema.key_column_usage kcu
		    JOIN information_schema.table_constraints tc
		      ON tc.constraint_name = kcu.constraint_name
		     AND tc.constraint_schema = kcu.constraint_schema
		    WHERE tc.constraint_type = 'FOREIGN KEY'
		      AND kcu.table_schema = c.table_schema
		      AND kcu.table_name = c.table_name
		      AND kcu.column_name = c.column_name
		  )
		ORDER BY c.table_name, c.column_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []finding
	for rows.Next() {
		var table, column string
		if err := rows.Scan(&table, &column); err != nil {
			return nil, err
		}
		if externalIDColumns[column] {
			continue
		}
		out = append(out, finding{
			rule:    "fk-missing",
			table:   table,
			column:  column,
			message: "bigintの%_id列だがFK制約が無い（外部サービスのIDならexternalIDColumnsへ追加する）",
		})
	}
	return out, rows.Err()
}

// checkFKIndexMissing は、FK制約の子側（先頭列）に索引が張られていないケースを検出する。
// PostgreSQLはFKを張っても子側の索引を自動生成しないため、親の1行削除・更新のたびに
// 子テーブル全体が走査される。
func checkFKIndexMissing(ctx context.Context, conn *pgx.Conn) ([]finding, error) {
	rows, err := conn.Query(ctx, `
		SELECT tc.table_name, kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
		  ON tc.constraint_name = kcu.constraint_name
		 AND tc.constraint_schema = kcu.constraint_schema
		WHERE tc.constraint_type = 'FOREIGN KEY'
		  AND tc.table_schema = 'public'
		  AND kcu.position_in_unique_constraint IS NULL
		  AND NOT EXISTS (
		    SELECT 1
		    FROM pg_index i
		    JOIN pg_class t ON t.oid = i.indrelid
		    JOIN pg_namespace n ON n.oid = t.relnamespace
		    JOIN pg_attribute a
		      ON a.attrelid = t.oid
		     AND a.attnum = i.indkey[0]
		    WHERE n.nspname = 'public'
		      AND t.relname = tc.table_name
		      AND a.attname = kcu.column_name
		  )
		ORDER BY tc.table_name, kcu.column_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []finding
	for rows.Next() {
		var table, column string
		if err := rows.Scan(&table, &column); err != nil {
			return nil, err
		}
		out = append(out, finding{
			rule:    "fk-index-missing",
			table:   table,
			column:  column,
			message: "FK列だが、その列を先頭に持つ索引が無い（親の削除・更新のたびに全表走査になる）",
		})
	}
	return out, rows.Err()
}

// checkNaming は索引・制約の命名規約（FreStyle準拠: idx_/uq_/fk_/ck_ プレフィックス）を検証する。
func checkNaming(ctx context.Context, conn *pgx.Conn) ([]finding, error) {
	rows, err := conn.Query(ctx, `
		SELECT c.conname,
		       t.relname,
		       CASE c.contype
		         WHEN 'f' THEN 'fk_'
		         WHEN 'c' THEN 'ck_'
		         WHEN 'u' THEN 'uq_'
		         ELSE NULL
		       END AS expected_prefix
		FROM pg_constraint c
		JOIN pg_class t ON t.oid = c.conrelid
		JOIN pg_namespace n ON n.oid = t.relnamespace
		WHERE n.nspname = 'public'
		  AND c.contype IN ('f', 'c', 'u')

		UNION ALL

		SELECT ic.relname AS conname,
		       t.relname,
		       CASE WHEN i.indisunique THEN 'uq_' ELSE 'idx_' END AS expected_prefix
		FROM pg_index i
		JOIN pg_class ic ON ic.oid = i.indexrelid
		JOIN pg_class t ON t.oid = i.indrelid
		JOIN pg_namespace n ON n.oid = t.relnamespace
		WHERE n.nspname = 'public'
		  AND NOT i.indisprimary
		  -- UNIQUE制約が自動生成する索引は上のpg_constraint側で既に検査しているため除く
		  AND NOT EXISTS (
		    SELECT 1 FROM pg_constraint c2
		    WHERE c2.conindid = i.indexrelid
		  )
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []finding
	for rows.Next() {
		var conname, table string
		var expectedPrefix *string
		if err := rows.Scan(&conname, &table, &expectedPrefix); err != nil {
			return nil, err
		}
		if expectedPrefix == nil {
			continue
		}
		if hasPrefix(conname, *expectedPrefix) {
			continue
		}
		out = append(out, finding{
			rule:    "naming",
			table:   table,
			column:  conname,
			message: "命名規約に反する（" + *expectedPrefix + "接頭辞を期待）",
		})
	}
	return out, rows.Err()
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
