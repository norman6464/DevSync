package persistence

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// pgUniqueViolationCode は PostgreSQL の一意制約違反を表す SQLSTATE。
const pgUniqueViolationCode = "23505"

// isUniqueViolation はエラーが一意制約違反かを判定する。
// アプリ側の重複チェックは同時実行をすり抜けるため、DB の制約違反を
// 「重複」として意味づけする入り口をここに揃える。
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode
}

// isNoRows はエラーが「該当行なし」かを判定する（sqlc/pgx実装向け）。
// 不在は (nil, nil) に正規化する規約のため、sqlc移行済みリポジトリはこれで
// pgx.ErrNoRows を吸収してから usecase 側へ返す。
func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
