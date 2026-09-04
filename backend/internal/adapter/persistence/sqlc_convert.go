package persistence

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// GORM時代のmodelは非ポインタのstring/int/boolフィールドにNULL許容カラムを載せていたため、
// sqlc(pgx)側のポインタ/pgtype表現との間で変換が繰り返し必要になる。ここに集約する。

// fromStringPtr はNULL許容カラムの *string を model の string へ変換する（NULLは空文字）。
func fromStringPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// toInt64Ptr は model の int をNULL許容カラム用の *int64 へ変換する。
func toInt64Ptr(n int) *int64 {
	v := int64(n)
	return &v
}

// fromInt64Ptr はNULL許容カラムの *int64 を model の int へ変換する（NULLは0）。
func fromInt64Ptr(p *int64) int {
	if p == nil {
		return 0
	}
	return int(*p)
}

// fromBoolPtr はNULL許容カラムの *bool を model の bool へ変換する（NULLはfalse）。
func fromBoolPtr(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

// toTimestamptz は model の *time.Time をNULL許容タイムスタンプへ変換する（nilはNULL）。
func toTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// fromTimestamptz はNULL許容タイムスタンプを model の *time.Time へ変換する（NULLはnil）。
func fromTimestamptz(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	tm := t.Time
	return &tm
}
