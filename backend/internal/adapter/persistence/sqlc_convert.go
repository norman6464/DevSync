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

// toInt64PtrFromUintPtr は model の *uint（NULL許容な外部キー等）を *int64 へ変換する。
func toInt64PtrFromUintPtr(p *uint) *int64 {
	if p == nil {
		return nil
	}
	v := int64(*p)
	return &v
}

// fromInt64PtrToUintPtr は *int64 を model の *uint へ変換する（NULLはnil）。
func fromInt64PtrToUintPtr(p *int64) *uint {
	if p == nil {
		return nil
	}
	v := uint(*p)
	return &v
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

// toTimestamptzNotNull は model の time.Time（NOT NULL カラム）をタイムスタンプへ変換する。
func toTimestamptzNotNull(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// fromTimestamptz はNULL許容タイムスタンプを model の *time.Time へ変換する（NULLはnil）。
func fromTimestamptz(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	tm := t.Time
	return &tm
}

// toDateNotNull は model の time.Time（NOT NULL カラム）を時刻情報を持たない日付へ変換する。
func toDateNotNull(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

// fromDate は日付カラムを model の time.Time へ変換する。
func fromDate(d pgtype.Date) time.Time {
	return d.Time
}

// dateStringLayout は model.StudyCircleCheckin.Date 等、APIでは "YYYY-MM-DD" の
// 文字列として扱う日付のレイアウト。
const dateStringLayout = "2006-01-02"

// toDateFromDateString は "YYYY-MM-DD" 形式の文字列（model側の表現）を日付カラムへ変換する。
func toDateFromDateString(s string) pgtype.Date {
	t, err := time.Parse(dateStringLayout, s)
	if err != nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: t, Valid: true}
}

// fromDateToDateString は日付カラムを "YYYY-MM-DD" 形式の文字列（model側の表現）へ変換する。
func fromDateToDateString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format(dateStringLayout)
}
