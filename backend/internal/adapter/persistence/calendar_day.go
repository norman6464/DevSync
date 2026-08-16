package persistence

import "time"

// normalizeToCalendarDay は時刻をその値が持つロケーションの暦日（00:00 UTC 表現）へ正規化する。
// time.Truncate(24h) は UTC エポック基準の 24 時間境界で丸めるため、DB セッションの
// タイムゾーンやサマータイムによって暦日とずれるが、こちらは年月日そのものを取り出すためずれない。
// DB の DATE() が返す値（その暦日の 00:00）はセッションのタイムゾーンによらず同じ暦日に写る。
func normalizeToCalendarDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// isNextCalendarDay は newer が older のちょうど翌暦日かを返す。
// 経過時間の割り算と違い、サマータイムで 23/25 時間になる日でも正しく判定できる。
// 引数は normalizeToCalendarDay で正規化済みであること。
func isNextCalendarDay(newer, older time.Time) bool {
	return newer.Equal(older.AddDate(0, 0, 1))
}

// isTodayOrYesterday は day が today 当日または前日かを返す。
// 引数は normalizeToCalendarDay で正規化済みであること。
func isTodayOrYesterday(day, today time.Time) bool {
	return day.Equal(today) || day.Equal(today.AddDate(0, 0, -1))
}
