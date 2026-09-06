// Package database はPostgreSQLスキーマの起動時自己適用を提供する。
package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// schemaGenSQL はAtlasの宣言的スキーマ（schema/schema.hcl）から
// `make db-schema-sql`（backend/Makefile）が機械生成したDDL（schema/schema.gen.sql）。
// バイナリに埋め込んで空のDBへ流すため、アプリとスキーマ定義が必ず同じ版になる。
// 正本はschema.hcl。このファイルはDO NOT EDIT（make db-schema-sqlで作り直す）。
//
//go:embed schema/schema.gen.sql
var schemaGenSQL string

// ApplySchema はschema.gen.sqlを「まだ何も無い空のDB」へ適用する。
// CREATE文のみで構成されておりDOブロック・IF NOT EXISTSを持たないため、
// 既存のテーブルがあるDBへ素で流すと衝突する。usersテーブルの有無で
// 「初めて適用するか」を判定し、既に在れば何もしない。
//
// 既存データのある本番DB・書き換え済みのローカルDBへ適用するのはこの関数の役割ではない。
// そちらはschema.hclを正本にした `make db-schema-apply`（Atlasの宣言的apply）を使う。
func ApplySchema(ctx context.Context, pool *pgxpool.Pool) error {
	applied, err := hasCoreSchema(ctx, pool)
	if err != nil {
		return fmt.Errorf("スキーマの適用状況の確認に失敗: %w", err)
	}
	if applied {
		return nil
	}
	if _, err := pool.Exec(ctx, schemaGenSQL); err != nil {
		return fmt.Errorf("スキーマの適用に失敗: %w", err)
	}
	return nil
}

// hasCoreSchema はusersテーブルの有無でスキーマ適用済みかを判定する。
func hasCoreSchema(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	var name *string
	if err := pool.QueryRow(ctx, `SELECT to_regclass('public.users')::text`).Scan(&name); err != nil {
		return false, err
	}
	return name != nil, nil
}
